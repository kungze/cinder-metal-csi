package rbd

import (
	"encoding/json"
	"fmt"
	"time"
)

// rbdDeviceInfo strongly typed JSON spec for rbd device list output (of type krbd).
type rbdDeviceInfo struct {
	ID             string `json:"id"`
	Pool           string `json:"pool"`
	RadosNamespace string `json:"namespace"`
	Name           string `json:"name"`
	Device         string `json:"device"`
}

// rbdGetDeviceList queries rbd about mapped devices and returns a list of rbdDeviceInfo
// It will selectively list devices mapped using krbd or nbd as specified by accessType.
func rbdGetDeviceList() ([]rbdDeviceInfo, error) {
	// rbd device list --format json
	var rbdDeviceList []rbdDeviceInfo

	stdout, err := utilsExecute("rbd", "device", "list", "--format="+"json")
	if err != nil {
		return nil, fmt.Errorf("error getting device list from rbd for devices of type (%s): %w", "dddddddddd", err)
	}
	err = json.Unmarshal([]byte(stdout), &rbdDeviceList)
	if err != nil {
		return nil, fmt.Errorf(
			"error to parse JSON output of device list for devices of type (%s): %w", "ddddddddddddddddddd", err)
	}

	return rbdDeviceList, nil
}

// findDeviceMappingImage finds a devicePath, if available, based on image spec (pool/{namespace/}image) on the node.
func findDeviceMappingImage(pool, image string) (string, bool) {

	// imageSpec := fmt.Sprintf("%s/%s", pool, image)
	// if namespace != "" {
	// 	imageSpec = fmt.Sprintf("%s/%s/%s", pool, namespace, image)
	// }

	rbdDeviceList, err := rbdGetDeviceList()
	if err != nil {
		//	log.WarningLog(ctx, "failed to determine if image (%s) is mapped to a device (%v)", imageSpec, err)

		return "", false
	}

	for _, device := range rbdDeviceList {
		if device.Name == image && device.Pool == pool {
			return device.Device, true
		}
	}

	return "", false
}

// Stat a path, if it doesn't exist, retry maxRetries times.
func waitForPath(pool, image string, maxRetries int) (string, bool) {
	for i := 0; i < maxRetries; i++ {
		if i != 0 {
			time.Sleep(time.Second)
		}

		device, found := findDeviceMappingImage(pool, image)
		if found {
			return device, found
		}
	}

	return "", false
}
