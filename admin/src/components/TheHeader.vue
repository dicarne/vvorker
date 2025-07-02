<script setup lang="ts">
import { NAvatar, NButton, useMessage } from 'naive-ui'
import FuncIcon from '@/components/icons/FuncIcon.vue'
import { onMounted, ref } from 'vue'
import type { UserInfo } from '@/types/auth'
import { getUserInfo, logout } from '@/api/auth'
import { useNavigate } from '@/composables/useNavigate'
const userInfo = ref<UserInfo>()
const message = useMessage()
const { navigate } = useNavigate()
const handleLogout = async () => {
  try {
    await logout()
    message.success('已退出登录')
    navigate('/login')
  } catch (error) {
    console.error(error)
    message.error('退出登录失败: ' + error)
  }
}
onMounted(async () => {
  try {
    userInfo.value = await getUserInfo()
  } catch (error) {
    console.error(error)
    message.error('获取用户信息失败: ' + error)
  }
})
</script>

<template>
  <div class="header v-flex-between-center">
    <div class="v-flex-center">
      <FuncIcon class="v-item" />
      <span class="v-item" style="font-size: 24px">VVorker</span>
    </div>
    <div class="v-flex-center">
      <NAvatar round>{{ userInfo?.userName.slice(0, 2).toLocaleUpperCase() }}</NAvatar>
      <NButton secondary type="primary" class="v-item" @click="handleLogout">登出</NButton>
    </div>
  </div>
</template>

<style scoped>
.header {
  width: 100vw;
  height: 60px;
  border-bottom: 1px solid #e5e5e5;
}
</style>
