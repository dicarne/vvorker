import { ref } from 'vue'
import type { UserInfo } from '@/types/auth'
import { getAppConfig } from '@/api/workers'
import type { VorkerSettingsProperties } from '@/types/workers'

const userInfo = ref<UserInfo>({
  userName: '',
  email: '',
  role: '',
  id: -1,
})

export const getProviderData = async () => {
  const appConfig = ref<VorkerSettingsProperties>()
  appConfig.value = await getAppConfig()
  return {
    userInfo: userInfo,
    appConfig: appConfig,
  }
}
