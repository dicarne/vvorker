import './assets/common.css'

import { createApp, ref } from 'vue'
import App from './App.vue'
import router from './router'
import type { UserInfo } from './types/auth'
import { getProviderData } from './provider/provider'

const providerData: any = await getProviderData()

const app = createApp(App)

for (const key in providerData) {
  app.provide(key, providerData[key])
}

app.use(router)

app.mount('#app')
