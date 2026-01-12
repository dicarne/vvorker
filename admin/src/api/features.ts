import http from './http'

export interface Feature {
  name: string
  enable: boolean
}

export async function getFeatures() {
  const res = await http.get('/api/features/list')
  return res.data.data as Feature[]
}
