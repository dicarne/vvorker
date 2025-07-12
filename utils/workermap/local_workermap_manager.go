package workermap

import "vvorker/utils"

var localWorkerMapManager *utils.SyncMap[string, *WorkerMap]

func init() {
	localWorkerMapManager = &utils.SyncMap[string, *WorkerMap]{}
}

func GetWorkerMap(uid string, maxCount int) (*WorkerMap, bool) {
	return localWorkerMapManager.LoadOrStore(uid, NewWorkerMap(uid, maxCount))
}

func DeleteWorkerMap(uid string) {
	localWorkerMapManager.Delete(uid)
}
