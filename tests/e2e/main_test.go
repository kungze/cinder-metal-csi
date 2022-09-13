package e2e

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/kungze/cinder-metal-csi/pkg/openstack"
	"sigs.k8s.io/e2e-framework/klient/conf"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"
)

var storageClassMap map[string]string = map[string]string{
	"rbd": "cinder-metal-csi-e2e-rbd",
	"lvm": "cinder-metal-csi-e2e-lvm",
}

var testenv env.Environment
var cloudConf string
var volumeType string
var cepyKeyring string
var cephUser string

func init() {
	flag.StringVar(&cloudConf, "cloud-conf", "./cloud.conf", "The openstack cinder auth config file.")
	flag.StringVar(&volumeType, "cinder-volume-type", "lvm", "The cinder volume type used to create StorageClass")
	flag.StringVar(&cephUser, "ceph-client-user", "cinder", "The ceph user used to access cinder volume pool")
	flag.StringVar(&cepyKeyring, "ceph-client-keyring", "", "The ceph client keyring corresponding to ceph-client-user")
}

func TestMain(m *testing.M) {
	flag.Parse()
	if storageClassMap[volumeType] == "" {
		panic(fmt.Errorf("The cinder volume type %s don't support.", volumeType))
	}
	osClient, err := openstack.CreateOpenstackClient(cloudConf)
	if err != nil {
		panic(err)
	}
	testenv = env.New()
	namespace := envconf.RandomName("sample-ns", 16)
	path := conf.ResolveKubeConfigFile()
	cfg := envconf.NewWithKubeConfig(path)
	testenv = env.NewWithConfig(cfg).WithContext(context.WithValue(context.Background(), osClientKey, &osClient))
	testenv.Setup(
		envfuncs.CreateNamespace(namespace),
	)
	testenv.Finish(
		envfuncs.DeleteNamespace(namespace),
	)
	os.Exit(testenv.Run(m))
}
