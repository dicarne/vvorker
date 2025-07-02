import './assets/common.css'

import { createApp, ref } from 'vue'
import App from './App.vue'
import router from './router'
import type { UserInfo } from './types/auth'

const userInfo = ref<UserInfo>({
  userName: '',
  email: '',
  role: '',
  id: -1,
})

const app = createApp(App)

app.provide('userInfo', userInfo)

app.use(router)

app.mount('#app')
