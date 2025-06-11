// 访问令牌创建请求接口
export interface AccessTokenCreateRequest {
  // 关联的 Worker UID
  worker_uid: string;
  // 访问令牌描述信息
  description: string;
  // 令牌是否永久有效
  forever: boolean;
  // 令牌过期时间
  expiration_time: string;
}

// 访问令牌列表请求接口
export interface AccessTokenListRequest {
  // 关联的 Worker UID
  worker_uid: string;
  // 分页页码，从 1 开始
  page: number;
  // 每页显示的记录数量
  page_size: number;
}

// 访问令牌删除请求接口
export interface AccessTokenDeleteRequest {
  // 关联的 Worker UID
  worker_uid: string;
}

// 内部白名单创建请求接口
export interface InternalWhiteListCreateRequest {
  // 关联的 Worker UID
  worker_uid: string;
  // 允许访问的 Worker UID
  allow_worker_uid: string;
  // 白名单描述信息
  description: string;
}

// 内部白名单列表请求接口
export interface InternalWhiteListListRequest {
  // 关联的 Worker UID
  worker_uid: string;
  // 分页页码，从 1 开始
  page: number;
  // 每页显示的记录数量
  page_size: number;
}

// 内部白名单更新请求接口
export interface InternalWhiteListUpdateRequest {
  // 关联的 Worker UID
  worker_uid: string;
  // 白名单描述信息
  description: string;
}

// 内部白名单删除请求接口
export interface InternalWhiteListDeleteRequest {
  // 关联的 Worker UID
  worker_uid: string;
  // 允许访问的 Worker UID
  allow_worker_uid: string;
}
