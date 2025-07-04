// type ListResourceRequest struct {
// 	Page     int    `json:"page"`
// 	PageSize int    `json:"page_size"`
// 	RType    string `json:"type"`
// }

// type ResourceData struct {
// 	UID  string `json:"uid"`
// 	Name string `json:"name"`
// 	Type string `json:"type"`
// }

// type ListResourceResponse struct {
// 	Total int64          `json:"total"`
// 	Data  []ResourceData `json:"data"`
// }

export interface ResourceData {
    uid: string
    name: string
    type: string
}
export interface ListResourceResponse {
    total: number
    data: ResourceData[]
}
export interface ListResourceRequest {
    page: number
    page_size: number
    type: string
}

// type CreateNewResourcesRequest struct {
// 	Name   string `json:"name"`
// }

// type DeleteResourcesReq struct {
// 	UID string `json:"uid"`
// }

// type DeleteResourcesResp struct {
// 	Status int `json:"status"` // 0: success, 1: failed
// }

// type CreateNewResourcesResponse struct {
// 	UID    string `json:"uid"`
// 	Status int    `json:"status"` // 0: success, 1: failed
// 	Name   string `json:"name"`
// 	Type   string `json:"type"`
// }


export interface CreateNewResourcesRequest {
    name: string
}

export interface DeleteResourcesReq {
    uid: string
}

export interface DeleteResourcesResp {
    // 0: success, 1: failed
    status: number
}

export interface CreateNewResourcesResponse {
    uid: string
    status: number
    name: string
    type: string
}