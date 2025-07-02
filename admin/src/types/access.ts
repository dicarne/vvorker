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
  worker_uid: string;
  id: number
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
  id: number;
}

// AccessKey 实体接口
export interface AccessKey {
  // 自增 ID
  id: number;
  // 创建时间
  created_at: string;
  // 更新时间
  updated_at: string;
  // 删除时间
  deleted_at: string | null;
  // 用户 ID
  user_id: number;
  // 访问密钥名称
  name: string;
  // 访问密钥
  key: string;
}

// InternalServerWhiteList 实体接口
export interface InternalServerWhiteList {
  // 自增 ID
  id: number;
  // 关联的 Worker UID
  worker_uid: string;
  // 允许访问的 Worker UID
  allow_worker_uid: string;
  // 白名单描述信息
  description: string;
  WorkerName: string;
}

// ExternalServerAKSK 实体接口
export interface ExternalServerAKSK {
  // 自增 ID
  id: number;
  // 创建时间
  created_at: string;
  // 更新时间
  updated_at: string;
  // 删除时间
  deleted_at: string | null;
  // 关联的 Worker UID
  worker_uid: string;
  // 访问密钥
  access_key: string;
  // 密钥
  secret_key: string;
  // 描述信息
  description: string;
  // 是否永久有效
  forever: boolean;
  // 过期时间
  expiration_time: string;
}

// ExternalServerToken 实体接口
export interface ExternalServerToken {
  // 自增 ID
  id: number;
  // 创建时间
  created_at: string;
  // 更新时间
  updated_at: string;
  // 删除时间
  deleted_at: string | null;
  // 关联的 Worker UID
  worker_uid: string;
  // 令牌
  token: string;
  // 描述信息
  description: string;
  // 是否永久有效
  forever: boolean;
  // 过期时间
  expiration_time: string;
}

// AccessRule 实体接口
export interface AccessRule {
  // 自增 ID
  id?: number;
  // 关联的 Worker UID
  worker_uid: string;
  // 规则类型
  rule_type: "internal" | "aksk" | "token" | "sso" | "open";
  // 规则描述信息
  description: string;
  // path
  path: boolean;
  rule_uid: string;
}


// 访问控制请求接口
export interface EnableAccessControlRequest {
  enable: boolean;
  worker_uid: string;
}

export interface AccessControlRequest {
  worker_uid: string;
}

export interface DeleteAccessRuleRequest {
  worker_uid: string;
  rule_uid: string;
}

export interface ListAccessRuleRequest {
  worker_uid: string;
  page: number;
  page_size: number;
}
