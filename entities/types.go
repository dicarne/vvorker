package entities

type RegisterRequest struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Status int `json:"status"`
}

func (r *RegisterRequest) Validate() bool {
	if r == nil {
		return false
	}
	if (r.UserName == "" && r.Email == "") || r.Password == "" {
		return false
	}
	if len(r.UserName) > 32 || len(r.Email) > 64 || len(r.Password) > 64 {
		return false
	}
	return true
}

type LoginRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Status int    `json:"status"`
	Token  string `json:"token"`
}

func (l *LoginRequest) Validate() bool {
	if l == nil {
		return false
	}
	if l.UserName == "" || l.Password == "" {
		return false
	}
	if len(l.UserName) > 32 || len(l.Password) > 64 {
		return false
	}
	return true
}

type GetUserResponse struct {
	UserName string `json:"userName"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	ID       uint   `json:"id"`
	VK       string `json:"vk"`
}

type DeleteWorkerRequest struct {
	UID string `json:"uid"`
}

func (d *DeleteWorkerRequest) Validate() bool {
	if d == nil {
		return false
	}
	if d.UID == "" {
		return false
	}
	if len(d.UID) > 64 {
		return false
	}
	return true
}

type WorkerUIDVersion struct {
	UID     string `json:"uid"`
	Version string `json:"version"`
}

type AgentDiffSyncWorkersResp struct {
	WorkerUIDVersions []WorkerUIDVersion `json:"worker_uid_versions"`
}

type AgentSyncWorkersReq struct {
	WorkerNames []string `json:"worker_names"`
}

type AgentGetWorkerByUIDReq struct {
	UID string `json:"uid"`
}

type AgentSyncWorkersResp struct {
	WorkerList *WorkerList `json:"worker_list"`
}

type NotifyEventRequest struct {
	EventName string            `json:"event_name"`
	Extra     map[string][]byte `json:"extra"`
}

func (n *NotifyEventRequest) Validate() bool {
	if n == nil {
		return false
	}
	if n.EventName == "" {
		return false
	}
	if len(n.EventName) > 64 {
		return false
	}
	return true
}

type NotifyEventResponse struct {
	Status int `json:"status"` // 0: success, 1: failed
}

type SyncNodesResponse struct {
}

type RunWorkerResponse struct {
	Status  int    `json:"status"` // 0: success, 1: failed
	RunResp []byte `json:"run_resp"`
}

type AgentFillWorkerReq struct {
	UID string `json:"uid"`
}

type AgentFillWorkerResp struct {
	NewTemplate string `json:"new_template"`
}

type DeleteResourcesReq struct {
	UID string `json:"uid"`
}

type DeleteResourcesResp struct {
	Status int `json:"status"` // 0: success, 1: failed
}

func (d *DeleteResourcesReq) Validate() bool {
	if d == nil {
		return false
	}
	if d.UID == "" {
		return false
	}
	if len(d.UID) > 64 {
		return false
	}
	return true
}

type CreateNewResourcesRequest struct {
	Name string `json:"name"`
}

func (r *CreateNewResourcesRequest) Validate() bool {
	if r == nil {
		return false
	}
	if r.Name == "" {
		return false
	}
	if len(r.Name) > 30 {
		return false
	}
	return true
}

type CreateNewResourcesResponse struct {
	UID    string `json:"uid"`
	Status int    `json:"status"` // 0: success, 1: failed
	Name   string `json:"name"`
	Type   string `json:"type"`
}

type ExecuteSQLReq struct {
	Sql              string `json:"sql"`
	Params           []any  `json:"params"`
	Method           string `json:"method"`
	ConnectionString string `json:"connection_string"`
}

type ExecuteSQLResp struct {
	Rows []string `json:"rows"`
}

type ExecuteSQLRespAll struct {
	Rows  [][]string `json:"rows"`
	Types []string   `json:"types"`
}
