package e2e

import (
	"os"
	"testing"

	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	apimachineryresource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cloud-provider/volume/helpers"
	"sigs.k8s.io/e2e-framework/klient/decoder"
)

const (
	initSize   = 1 * helpers.GiB
	expandSize = 2 * helpers.GiB
)

func newSecret(namespace, name string, data map[string]string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		StringData: data,
	}
}

func newStorageClass(namespace string, volumeType string, allowExpansion bool, t *testing.T) *storagev1.StorageClass {
	if storageClassMap[volumeType] == "" {
		t.Fatalf("The cinder volume type %s don't support.", volumeType)
	}
	return &storagev1.StorageClass{
		ObjectMeta:           metav1.ObjectMeta{Name: storageClassMap[volumeType], Namespace: namespace},
		Provisioner:          "cinder.metal.csi",
		AllowVolumeExpansion: &allowExpansion,
		Parameters: map[string]string{
			"cinderVolumeType":                               volumeType,
			"csi.storage.k8s.io/node-stage-secret-name":      rbdVolumeSecret,
			"csi.storage.k8s.io/node-stage-secret-namespace": namespace,
		},
	}
}

func newPresistentVolumeClaim(namespace string, name string, sc string, size int64, t *testing.T) *corev1.PersistentVolumeClaim {
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
			StorageClassName: &sc,
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					"storage": *apimachineryresource.NewQuantity(size, apimachineryresource.BinarySI),
				},
			},
		},
	}
}

func newPod(namespace, name, pvc string, t *testing.T) *corev1.Pod {
	f, err := os.Open("./pod.yaml")
	if err != nil {
		t.Fatal(err)
	}
	obj, err := decoder.DecodeAny(f)
	if err != nil {
		t.Fatal(err)
	}
	podObj := obj.(*corev1.Pod)
	podObj.SetNamespace(namespace)
	podObj.SetName(name)
	podObj.Spec.Volumes[0].PersistentVolumeClaim.ClaimName = pvc
	return podObj
}

func podInNode(pod *corev1.Pod, node *corev1.Node) bool {
	for _, addr := range node.Status.Addresses {
		if pod.Status.HostIP == addr.Address {
			return true
		}
	}
	return false
}
