package workermap

import "vvorker/utils"

type WorkerMap struct {
	UID      string
	MaxCount int

	portsManager *utils.SyncMap[string, int32]
}

func NewWorkerMap(uid string, maxCount int) *WorkerMap {
	return &WorkerMap{
		UID:          uid,
		MaxCount:     maxCount,
		portsManager: &utils.SyncMap[string, int32]{},
	}
}
