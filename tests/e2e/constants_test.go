package e2e

import "time"

type contextKey string

const (
	storageClassKey          contextKey = "storage-class"
	cinderVolumeIdKey        contextKey = "cinder-volume-id"
	presistentVolumeClaimKey contextKey = "presistent-volume-claim"
	podKey                   contextKey = "pod"
	osClientKey              contextKey = "os-client"
)

const (
	waitNum      = 10
	waitInterval = 30 * time.Second
)
