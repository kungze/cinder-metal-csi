package local

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kungze/cinder-metal-csi/pkg/brick/utils"
	"k8s.io/klog/v2"
)

//ConnLocal A local volume type object
type ConnLocal struct {
	volumeID string
}

//NewLocalConnector Build a local volume type connection object
func NewLocalConnector(connInfo map[string]interface{}) *ConnLocal {
	conn := &ConnLocal{}
	conn.volumeID = utils.ToString(connInfo["volume_id"])
	return conn
}

//ConnectVolume Connect the local volume
func (c *ConnLocal) ConnectVolume() (map[string]string, error) {
	res := map[string]string{}
	globStr := fmt.Sprintf("/dev/*/*%s", c.volumeID)
	paths, err := filepath.Glob(globStr)
	if err != nil {
		return nil, err
	}
	if len(paths) != 1 {
		klog.Errorf("lvm volume path not found", err)
		return nil, err
	}
	klog.V(3).Infof("Get lvm path %s success", paths[0])
	res["path"] = paths[0]
	return res, nil
}

//DisConnectVolume DisConnect the local volume
func (c *ConnLocal) DisConnectVolume() error {
	klog.V(3).Info("local volume disconnect volume success")
	return nil
}

//ExtendVolume Extend the local volume
func (c *ConnLocal) ExtendVolume() (int64, error) {
	globStr := fmt.Sprintf("/dev/*/*%s", c.volumeID)
	paths, err := filepath.Glob(globStr)
	if err != nil {
		return 0, err
	}
	if len(paths) != 1 {
		klog.Errorf("lvm volume path not found", err)
		return 0, err
	}
	sizeCmd := fmt.Sprintf("lvdisplay --units B %s 2>&1 | grep 'LV Size' | awk '{print $3}'", globStr)
	out, err := utils.Execute(sizeCmd)
	if err != nil {
		klog.Errorf("Exec lvdisplay command failed", err)
		return 0, err
	}
	sizeStr := strings.Split(out, ".")[0]
	sizeInt, err := strconv.ParseInt(strings.TrimSpace(sizeStr), 10, 64)
	if err != nil {
		klog.Errorf("Parse lvm size failed", err)
		return 0, err
	}
	klog.V(3).Infof("Get lvm %s size %s success", c.volumeID, sizeInt)
	return sizeInt, nil
}

//GetDevicePath Get the volume device path
func (c *ConnLocal) GetDevicePath() string {
	globStr := fmt.Sprintf("/dev/*/*%s", c.volumeID)
	paths, err := filepath.Glob(globStr)
	if err != nil {
		return ""
	}
	if len(paths) != 1 {
		klog.Errorf("lvm volume path not found", err)
		return ""
	}
	klog.V(3).Infof("Get lvm path %s success", paths[0])
	return paths[0]
}
