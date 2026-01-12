import type {
  UserInfo,
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  RegisterResponse,
} from '@/types/auth'
import api from './http'

export const login = async (req: LoginRequest) => {
  const url = 'api/auth/login'
  const res = await api.post(url, req)
  return res.data
}

export const register = async (req: RegisterRequest) => {
  const res = await api.post('api/auth/register', req)
  return res.data.data as RegisterResponse
}

export const getUserInfo = async () => {
  const res = await api.get('api/user/info')
  return res.data.data as UserInfo
}

export const logout = () => {
  return api.get('api/auth/logout')
}

export const createAccessKey = (name: string) => {
  return api.post('api/user/create-access-key', {
    name: name
  })
}

export const getAccessKeys = () => {
  return api.post('api/user/access-keys')
}

export const deleteAccessKey = (access_key: string) => {
  return api.post('api/user/delete-access-key', {
    key: access_key
  })
}


export const enableOtp = () => {
  return api.post('api/otp/enable')
}

export const disableOtp = () => {
  return api.post('api/otp/disable')
}

export const validOtp = (code: string) => {
  return api.post('api/otp/valid-add?code=' + code)
}

export const isEnableOtp = () => {
  return api.post('api/otp/is-enable')
}