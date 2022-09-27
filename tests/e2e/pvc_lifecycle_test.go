package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/kungze/cinder-metal-csi/pkg/openstack"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	apimachineryresource "k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cloud-provider/volume/helpers"
	"sigs.k8s.io/e2e-framework/klient/k8s"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

const rbdVolumeSecret = "cinder-metal-csi-e2e-rbd"

func TestPVCLifecycle(t *testing.T) {
	f := features.New("PresistentVolumeClain Lifecycle").
		Setup(func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			r := c.Client().Resources()
			secret := newSecret(c.Namespace(), rbdVolumeSecret, map[string]string{"cephClientUser": cephUser, "cephClientKey": cepyKeyring})
			err := r.Create(ctx, secret)
			if err != nil {
				t.Fatal(err)
			}
			sc := newStorageClass(c.Namespace(), volumeType, true, t)
			err = r.Create(ctx, sc)
			if err != nil {
				t.Fatal(err)
			}
			return context.WithValue(ctx, storageClassKey, sc)
		}).
		Assess("Create PresistentVolumeClain", func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			r := c.Client().Resources()
			pvc := newPresistentVolumeClaim(c.Namespace(), "e2e-test-pvc", ctx.Value(storageClassKey).(*storagev1.StorageClass).GetName(), initSize, t)
			err := r.Create(ctx, pvc)
			if err != nil {
				t.Fatal(err)
			}
			err = wait.For(conditions.New(r).ResourceMatch(pvc, func(object k8s.Object) bool {
				return object.(*corev1.PersistentVolumeClaim).Status.Phase == "Bound"
			}))
			if err != nil {
				t.Fatal(err)
			}
			osClient := *(ctx.Value(osClientKey).(*openstack.IOpenstack))
			for i := 0; i <= waitNum; i++ {
				if i == waitNum {
					t.Fatal("The related cinder volume not fount.")
				}
				vols, err := osClient.GetVolumeByName(pvc.Spec.VolumeName)
				if err != nil {
					t.Fatal(err)
				}
				if len(vols) > 0 {
					vol := vols[0]
					if vol.Size != initSize/helpers.GiB {
						t.Errorf("The cinder volume size %d not equal to pvc %s request storage size %d", vol.Size, pvc.GetName(), initSize/helpers.GiB)
					}
					if vol.VolumeType != volumeType {
						t.Errorf("The cinder volume's type %s not equal to %s", vol.VolumeType, volumeType)
					}
					ctx = context.WithValue(ctx, cinderVolumeIdKey, vol.ID)
					break
				}
				time.Sleep(waitInterval)
			}
			return context.WithValue(ctx, presistentVolumeClaimKey, pvc)
		}).
		Assess("Expand PresistentVolumeClain", func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			r := c.Client().Resources()
			pvc := ctx.Value(presistentVolumeClaimKey).(*corev1.PersistentVolumeClaim)
			patch, err := json.Marshal(map[string]any{
				"spec": map[string]any{
					"resources": map[string]any{
						"requests": map[string]string{
							"storage": fmt.Sprintf("%dGi", expandSize/helpers.GiB),
						},
					},
				},
			})
			if err != nil {
				t.Fatal(err)
			}
			err = r.Patch(ctx, pvc, k8s.Patch{PatchType: types.StrategicMergePatchType, Data: patch})
			if err != nil {
				t.Fatal(err)
			}
			err = wait.For(conditions.New(r).ResourceMatch(pvc, func(object k8s.Object) bool {
				pvc := object.(*corev1.PersistentVolumeClaim)
				for _, cond := range pvc.Status.Conditions {
					if cond.Type == "FileSystemResizePending" {
						return true
					}
				}
				return false
			}))
			if err != nil {
				t.Error(err)
			}
			osClient := *(ctx.Value(osClientKey).(*openstack.IOpenstack))
			vol, err := osClient.GetVolumeByID(ctx.Value(cinderVolumeIdKey).(string))
			if err != nil {
				t.Fatal(err)
			}
			if vol.Size != int(expandSize/helpers.GiB) {
				t.Errorf("The cinder volume size don't expand to expected size %d", expandSize)
			}
			return ctx
		}).
		Assess("Mount PresistentVolume", func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			r := c.Client().Resources()
			pvc := ctx.Value(presistentVolumeClaimKey).(*corev1.PersistentVolumeClaim)
			podObj := newPod(c.Namespace(), "e2e-test-pod", pvc.GetName(), t)
			err := r.Create(ctx, podObj)
			if err != nil {
				t.Error(err)
				return ctx
			}
			ctx = context.WithValue(ctx, podKey, podObj)
			err = wait.For(conditions.New(r).PodRunning(podObj))
			if err != nil {
				t.Error(err)
				return ctx
			}
			err = r.Get(ctx, pvc.GetName(), pvc.GetNamespace(), pvc)
			if err != nil {
				t.Error(err)
				return ctx
			}
			if !pvc.Status.Capacity.Storage().Equal(*apimachineryresource.NewQuantity(expandSize, apimachineryresource.BinarySI)) {
				t.Errorf("The pvc's capacity storage don't change to expect size: %d", expandSize)
				return ctx
			}
			osClient := *(ctx.Value(osClientKey).(*openstack.IOpenstack))
			volID := ctx.Value(cinderVolumeIdKey).(string)
			vol, err := osClient.GetVolumeByID(volID)
			if err != nil {
				t.Error(err)
				return ctx
			}
			if vol.Status != "in-use" {
				t.Errorf("The cinder volume's status is '%s', but the expected status is 'in-use", vol.Status)
				return ctx
			}

			attach, err := osClient.GetAttachmentByVolumeID(volID)
			if err != nil {
				t.Error(err)
				return ctx
			}

			csinodes := storagev1.CSINodeList{}
			err = r.List(ctx, &csinodes)
			if err != nil {
				t.Error(err)
				return ctx
			}
			for _, csinode := range csinodes.Items {
				node := &corev1.Node{}
				err = r.Get(ctx, csinode.GetName(), csinode.GetNamespace(), node)
				if err != nil {
					t.Error(err)
					return ctx
				}
				if podInNode(podObj, node) {
					if csinode.GetUID() == types.UID(attach.Instance) {
						return ctx
					}
				}
			}
			t.Error("Cinder volume attachment not match k8s pod's node.")
			return ctx
		}).
		Assess("Unmount PresistentVolume", func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			r := c.Client().Resources()
			pod := ctx.Value(podKey).(*corev1.Pod)
			err := r.Delete(ctx, pod)
			if err != nil {
				t.Error(err)
			}
			err = wait.For(conditions.New(r).ResourceDeleted(pod))
			if err != nil {
				t.Error(err)
			}
			osClient := *(ctx.Value(osClientKey).(*openstack.IOpenstack))
			volID := ctx.Value(cinderVolumeIdKey).(string)
			for i := 0; i <= waitNum; i++ {
				vol, err := osClient.GetVolumeByID(volID)
				if err != nil {
					t.Error(err)
					return ctx
				}
				if vol.Status == "available" {
					break
				}
				if i == waitNum {
					t.Errorf("The expected cinder volume status is available, but now it is %s.", vol.Status)
				}
				time.Sleep(waitInterval)
			}
			return ctx
		}).
		Assess("Delete PresistentVolumeClain", func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			r := c.Client().Resources()
			pvc := ctx.Value(presistentVolumeClaimKey).(*corev1.PersistentVolumeClaim)
			if err := r.Delete(ctx, pvc); err != nil {
				t.Error(err)
			}
			if err := wait.For(conditions.New(r).ResourceDeleted(pvc)); err != nil {
				t.Fatal(err)
			}
			osClient := *(ctx.Value(osClientKey).(*openstack.IOpenstack))
			volID := ctx.Value(cinderVolumeIdKey).(string)
			for i := 0; i <= waitNum; i++ {
				vol, err := osClient.GetVolumeByID(volID)
				if err != nil {
					t.Error(err)
					return ctx
				}
				if vol == nil {
					break

				}
				if i == waitNum {
					t.Error("Related cinder volume don't be deleted.")
				}
				time.Sleep(waitInterval)
			}
			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			r := c.Client().Resources()
			err := r.Delete(ctx, ctx.Value(storageClassKey).(*storagev1.StorageClass))
			if err != nil {
				t.Error(err)
			}
			return ctx
		}).Feature()
	testenv.Test(t, f)
}
