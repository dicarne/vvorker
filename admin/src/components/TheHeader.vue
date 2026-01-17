<script setup lang="ts">
import { inject, onMounted, ref, type Ref } from 'vue'
import { CH } from '@/lib/color'
import { NAvatar, NButton, NIcon, useMessage, NTag } from 'naive-ui'
import type { UserInfo } from '@/types/auth'
import { getUserInfo, logout } from '@/api/auth'
import { getAppConfig } from '@/api/workers'
import { useNavigate } from '@/composables/useNavigate'
import type { VorkerSettingsProperties } from '@/types/workers'
const userInfo = ref<UserInfo>()
const appConfig = inject<Ref<VorkerSettingsProperties>>("appConfig")
const message = useMessage()
const { navigate } = useNavigate()
const handleLogout = async () => {
  try {
    await logout()
    message.success('已退出登录')
    navigate('/login')
  } catch (error) {
    console.error('logout Error', error)
    message.error('退出登录失败')
  }
}

const userInfoInj = inject<Ref<UserInfo>>('userInfo')!

onMounted(async () => {
  try {
    userInfo.value = await getUserInfo()
    userInfoInj.value = userInfo.value
  } catch (error) {
    console.error('getUserInfo Error', error)
    message.error('获取用户信息失败')
  }
})
</script>

<template>
  <div class="header v-flex-between-center">
    <div class="v-flex-center">
      <img class="v-item" style="height: 30px" src="@/assets/logo.png" alt="" />
      <span class="v-item" style="font-size: 24px">VVorker</span>
      <NTag v-if="appConfig?.Version" class="v-item" size="small">{{ appConfig.Version }}</NTag>
    </div>
    <div class="v-flex-center">
      <NAvatar class="v-avatar" :style="{ background: userInfo?.userName ? CH.hex(userInfo.userName) : '#cccccc' }">
        {{ userInfo?.userName.slice(0, 2).toUpperCase() }}
      </NAvatar>
      <NButton secondary type="primary" class="v-item" @click="handleLogout">登出</NButton>
    </div>
  </div>
</template>

<style scoped>
.header {
  height: 60px;
  border-bottom: 1px solid #e5e5e5;
}
</style>
