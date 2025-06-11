import api from './http'
import { Task, TaskLog, VorkerSettingsProperties, WorkerItem, WorkerLog } from '@/types/workers'
import {
  AccessTokenCreateRequest,
  AccessTokenListRequest,
  AccessTokenDeleteRequest,
  InternalWhiteListCreateRequest,
  InternalWhiteListListRequest,
  InternalWhiteListUpdateRequest,
  InternalWhiteListDeleteRequest,
  InternalServerWhiteList,
  ExternalServerToken,
  AccessRule,
  ListAccessRuleRequest,
  DeleteAccessRuleRequest,
  AccessControlRequest,
  EnableAccessControlRequest
} from '@/types/access' // 假设存在对应的类型定义

export const getWorker = (uid: string) => {
  return api
    .get<{ data: WorkerItem[] }>(`/api/worker/${uid}`)
    .then((res) => res.data.data?.[0])
}

export const getWorkerCursor = (offset: number, limit: number) => {
  return api
    .get<{ data: WorkerItem[] }>(`/api/workers/${offset}/${limit}`)
    .then((res) => res.data.data)
}

export const getAllWorkers = () => {
  return api
    .get<{ data: WorkerItem[] }>('/api/allworkers')
    .then((res) => res.data.data)
}

export const createWorker = (worker: WorkerItem) => {
  return api.post('/api/worker/create', worker).then((res) => res.data)
}

export const deleteWorker = (uid: string) => {
  return api.delete(`/api/worker/${uid}`, {}).then((res) => res.data)
}

export const updateWorker = (uid: string, worker: WorkerItem) => {
  return api.post(`/api/worker/${uid}`, worker).then((res) => res.data)
}

export const flushWorker = (uid: string) => {
  return api.get(`/api/worker/flush/${uid}`, {}).then((res) => res.data)
}

export const flushAllWorkers = () => {
  return api.get(`/api/workers/flush`, {}).then((res) => res.data)
}

export const getAppConfig = () => {
  return api
    .get<{ data: VorkerSettingsProperties }>(`/api/vvorker/config`, {})
    .then((res) => res.data.data)
}

export const runWorker = (uid: string) => {
  return api.get(`/api/worker/run/${uid}`, {}).then((res) => res.data)
}

export const getWorkersStatus = (uids: string[]) => {
  return api
    .post(`/api/workers/status`, { uids })
    .then((res) => res.data.data)
}


// 访问令牌相关 API
export const createAccessToken = (request: AccessTokenCreateRequest) => {
  return api.post('/api/worker/access/token/create', request).then((res) => res.data)
}

export const listAccessTokens = async (request: AccessTokenListRequest) => {
  return (await api.post<CommonResponse<{ access_tokens: ExternalServerToken[] }>>('/api/worker/access/token/list', request)).data
}

export const deleteAccessToken = (request: AccessTokenDeleteRequest) => {
  return api.post('/api/worker/access/token/delete', request).then((res) => res.data)
}

// 定义通用响应类型
interface CommonResponse<T> {
  code: number;
  msg: string;
  data: T;
}

// 内部白名单相关 API
export const createInternalWhiteList = async (request: InternalWhiteListCreateRequest) => {
  const res = await api.post<CommonResponse<{ internal_white_list: InternalServerWhiteList }>>('/api/worker/access/whitelist/create', request);
  return res.data;
};

export const listInternalWhiteLists = async (request: InternalWhiteListListRequest) => {
  const res = await api.post<CommonResponse<{ internal_white_lists: InternalServerWhiteList[] }>>('/api/worker/access/whitelist/list', request);
  return res.data;
};

export const updateInternalWhiteList = async (request: InternalWhiteListUpdateRequest) => {
  const res = await api.post<CommonResponse<null>>('/api/worker/access/whitelist/update', request);
  return res.data;
};

export const deleteInternalWhiteList = async (request: InternalWhiteListDeleteRequest) => {
  const res = await api.post<CommonResponse<null>>('/api/worker/access/whitelist/delete', request);
  return res.data;
};

// 访问控制相关 API
export const updateEnableAccessControl = async (request: EnableAccessControlRequest) => {
  const res = await api.post<CommonResponse<null>>('/api/worker/access/control/update-control', request);
  return res.data;
};

export const getAccessControl = async (request: AccessControlRequest) => {
  const res = await api.post<CommonResponse<{ EnableAccessControl: boolean }>>('/api/worker/access/control/get-control', request);
  return res.data;
};

export const addAccessRule = async (request: AccessRule) => {
  const res = await api.post<CommonResponse<null>>('/api/worker/access/control/create-rule', request);
  return res.data;
};

export const deleteAccessRule = async (request: DeleteAccessRuleRequest) => {
  const res = await api.post<CommonResponse<null>>('/api/worker/access/control/delete-rule', request);
  return res.data;
};

export const listAccessRules = async (request: ListAccessRuleRequest) => {
  const res = await api.post<CommonResponse<{ total: number, access_rules: AccessRule[] }>>('/api/worker/access/control/list-rules', request);
  return res.data;
};

export const listTasks = (page: number, page_size: number) => {
  return api.post<{
    data: {
      total: number,
      tasks: Task[]
    }
  }>('/api/ext/task/list', { page, page_size })
}

export const getTaskLogs = (worker_uid: string, trace_id: string, page: number, page_size: number) => {
  return api.post<{
    data: {
      total: number,
      logs: TaskLog[]
    }
  }>('/api/ext/task/logs', { trace_id, page, page_size, worker_uid })
}

export const interruptTask = (trace_id: string, worker_uid: string) => {
  return api.post('/api/ext/task/cancel', { trace_id, worker_uid })
}

export const getWorkerLogs = (uid: string, page: number, page_size: number) => {
  return api.post<{ data: { total: number, logs: WorkerLog[] } }>(`/api/worker/logs/${uid}`, { page, page_size })
}


