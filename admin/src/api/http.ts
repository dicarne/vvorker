import axios from 'axios'
// 引入 lodash 的 debounce 函数
import { debounce } from 'lodash'
import router from '@/router'

const instance = axios.create({})
const loginExpired = debounce(() => {
  console.log("loginExpired")
  router.push('/login')
}, 500)

instance.interceptors.response.use(
  (response) => {
    if (!!response.data.code) {
      throw response.data.msg
    }
    return response
  },
  (error) => {
    console.log(error)
    if (error.response.status === 403) {
      loginExpired()
      return Promise.reject("登录过期，请重新登录")
    } else {
      return Promise.reject(error)
    }
  },
)

export default instance
