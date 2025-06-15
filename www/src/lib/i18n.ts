import { atom } from 'nanostores'

const zh = {
  login: '登陆',
  logout: '登出',
  register: '注册',
  username: '用户名',
  password: '密码',
  email: '邮箱',
  submit: '提交',
  code: '代码',
  editor: '编辑器',
  edit: '编辑',
  run: '运行',
  sync: '同步',
  create: '创建',
  delete: '删除',
  openWeb: '打开',
  refresh: '刷新',
  deleteWorker: '删除函数',
  deleteNode: '删除节点',
  resourceName: '资源名称',
  cancel: '取消',
  id: 'ID',
  task: "任务",
  completed: "完成",
  running: "运行中",
  canceled: "取消",
  failed: "失败",
  log: "日志",
  look: "查看",
  node: "节点",
  workerEntry: "函数入口",
  add: "新增",

  tokenOnce: "此 token 仅显示一次，请妥善保存。",
  deleteConfirm: "确认删除吗？",
  important: "重要提示",
  addAccessKey: "添加访问密钥",
  addInternalAccess: "添加内部访问",
  internalAccess: "内部访问",
  accessKey: "访问密钥",
  enableAccessControl: "启用访问控制",
  addRule: "添加规则",
  prefix: "路由前缀",

  workerConfirmDelete: '确认删除函数',
  workerDeleteSuccess: '删除成功',
  workerCreateSuccess: '创建成功',
  workerSaveSuccess: '保存成功',
  workerSyncSuccess: '同步成功',
  nodeDeleteSuccess: '删除成功',
  nodeSyncSuccess: '同步成功',
  backToList: '返回列表',
  noWorkerPrompt: '空空如也',

  notLoggedInPrompt: '没有登陆，正在跳转...',
  loggingOutPrompt: '正在退出登陆...',
  loginSuccess: '登陆成功',
  loginFailed: '用户名或密码错误',
  registerSuccess: '注册成功',
  registerFailed: '注册失败',

  warnDeleteResource: "确定要删除这个资源吗？",

  back: '返回',
  save: '保存',
  property: '属性',
  config: '配置',
  logs: '日志',
  rules: '规则',
  auth: '鉴权',


  "internal": "内部访问",
  "aksk": "AccessKey 与 AccessSecret",
  "token": "AccessKey",
  "sso": "SSO",
  "open": "开放"
}

const en = {
  login: 'Login',
  logout: 'Logout',
  register: 'Register',
  username: 'Username',
  password: 'Password',
  email: 'Email',
  submit: 'Submit',
  code: 'Code',
  editor: 'Editor',
  edit: 'Edit',
  run: 'Run',
  sync: 'Sync',
  create: 'New',
  delete: 'Delete',
  openWeb: 'Open',
  refresh: 'Refresh',
  deleteWorker: 'Delete worker',
  deleteNode: 'Delete Node',
  resourceName: 'Resource name',
  cancel: 'Cancel',
  id: 'ID',
  task: "Task",
  completed: "Completed",
  running: "Running",
  canceled: "Canceled",
  failed: "Failed",
  log: "Log",
  look: "Look",
  tokenOnce: "This token will only be displayed once. Please keep it safe.",
  deleteConfirm: "Are you sure you want to delete this item?",
  important: "Important Notice",
  addAccessKey: "Add Access Key",
  addInternalAccess: "Add Internal Access",
  internalAccess: "Internal Access",
  accessKey: "Access Key",
  enableAccessControl: "Enable Access Control",
  addRule: "Add Rule",
  prefix: "Prefix",
  node: "Node",
  workerEntry: "Worker Entry",
  add: "Add",

  workerConfirmDelete: 'Confirm delete worker',
  workerDeleteSuccess: 'Delete worker success',
  workerCreateSuccess: 'Create worker success',
  workerSaveSuccess: 'Save worker success',
  workerSyncSuccess: 'Sync worker success',
  nodeDeleteSuccess: 'Delete node success',
  nodeSyncSuccess: 'Sync node success',
  backToList: 'Back',
  noWorkerPrompt: 'No workers here.',

  notLoggedInPrompt: 'Not logged in, redirecting...',
  loggingOutPrompt: 'Logging out...',
  loginSuccess: '',
  loginFailed: 'Incorrect username or password',
  registerSuccess: 'Successfully registered',
  registerFailed: 'Failed to register',

  warnDeleteResource: "Are you sure you want to delete this resource?",
  back: 'Back',
  save: 'Save',
  property: 'Property',
  config: 'Config',
  logs: 'Logs',
  rules: 'Rules',
  auth: 'Auth',

  "internal": "Internal access",
  "aksk": "AccessKey with AccessSecret",
  "token": "AccessKey",
  "sso": "SSO",
  "open": "Open",
}

export type Key = keyof typeof zh & keyof typeof en

export const i18n = (k: Key) => {
  return { zh, en }[$lang.get().split('-')[0] || 'en']![k] || k
}

export const t = new Proxy(
  {},
  {
    get(_, k) {
      return i18n(k as Key)
    },
  }
) as Record<Key, string>

export const $lang = atom(
  typeof navigator !== 'undefined' ? navigator.language : 'en'
)
