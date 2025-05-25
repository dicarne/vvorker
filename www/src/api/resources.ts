//  /api/ext/list 
import { CreateNewResourcesRequest, CreateNewResourcesResponse, DeleteResourcesReq, DeleteResourcesResp, ResourceData } from '@/types/resources'
import api from './http'

export const getResourceList = (page: number, pageSize: number, rtype: string) => {
    return api
        .post<{ data: { total: Number, data: ResourceData[] } }>('/api/ext/list', {
            page,
            page_size: pageSize,
            type: rtype
        })
        .then((res) => res.data.data)
}


export const createResource = (name: string, rtype: string) => {
    return api
        .post<{ data: CreateNewResourcesResponse }>(`/api/ext/${rtype}/create-resource`, {
            name
        } as CreateNewResourcesRequest)
        .then((res) => res.data.data)
}

export const deleteResource = (uid: string, rtype: string) => {
    return api
        .post<{ data: DeleteResourcesResp }>(`/api/ext/${rtype}/delete-resource`, {
            uid,
        } as DeleteResourcesReq)
        .then((res) => res.data.data)
}