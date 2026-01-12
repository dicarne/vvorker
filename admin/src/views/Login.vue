<script setup lang="ts">
import { inject, ref, type Ref } from 'vue'
import { useMessage, NCard, NForm, NFormItem, NInput, NButton, NModal, NInputOtp, NSpace } from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui' // 导入 FormInst 类型
import type { LoginRequest, UserInfo } from '@/types/auth'
import type { VorkerSettingsProperties } from '@/types/workers'
import { passwordRules, usernameRules } from '@/constant/formrules'
import { getUserInfo, login } from '@/api/auth'
import { useNavigate } from '@/composables/useNavigate'
import { APIError } from '@/api/http'

const appConfig = inject<Ref<VorkerSettingsProperties>>('appConfig')!
// 初始化表单数据
const form = ref<LoginRequest>({
  userName: '',
  password: '',
})

// 表单引用，明确指定类型为 FormInst 或 null
const formRef = ref<FormInst | null>(null)

const rules: FormRules = {
  userName: usernameRules,
  password: passwordRules,
}

// 消息提示实例
const message = useMessage()

const { navigate } = useNavigate()

const userInfo = inject<Ref<UserInfo>>('userInfo')!

const loading = ref(false)

// OTP相关
const showOtpModal = ref(false)
const otpCode = ref([''])

// 处理登录逻辑的函数
const handleLogin = async () => {
  if (!formRef.value) return
  loading.value = true
  try {
    // 校验表单
    await formRef.value.validate()
  } catch (error) {
    console.error('formRef validate Error', error)
    loading.value = false
    return
  }
  try {
    // 调用登录接口
    const res = await login(form.value)
    if (res.code === 12) {
      // 需要OTP验证
      showOtpModal.value = true
      loading.value = false
      return
    }
    if (res.code === 0) {
      // 登录成功获取用户信息
      try {
        userInfo.value = await getUserInfo()
      } catch (error) {
        console.error('getUserInfo Error', error)
        message.error('登录成功后获取用户信息失败')
      }
      loading.value = false
      message.success('登录成功')
      navigate('/workers')
    }
  } catch (error: any) {
    if (error instanceof APIError) {
      console.error('login Error', error.code)
      if (error.code === 6) {
        message.error("登陆失败")
      } else if (error.code === 12) {
        // 需要OTP验证
        showOtpModal.value = true
      } else {
        const errorMsg = error.data?.message || '登录失败'
        message.error(errorMsg)
      }
    } else {
      console.error('login Error', error)
      const errorMsg = error.response?.data?.message || '登录失败'
      message.error(errorMsg)
    }

    loading.value = false
  }
}

// 处理OTP验证登录
const handleOtpLogin = async (value: string[]) => {
  if (value.join('').length != 6) {
    return
  }
  loading.value = true
  try {
    form.value.otpcode = value.join('')
    await handleLogin()
  } catch (error: any) {
    console.error('login Error', error)
    const errorMsg = error.response?.data?.message || 'OTP验证失败'
    message.error(errorMsg)
    loading.value = false
  }
}

</script>
<template>
  <div class="v-base-page v-flex-center">
    <NCard style="width: 400px">
      <template #header>
        <div class="v-card-header">登录</div>
      </template>
      <NForm :model="form" :rules="rules" ref="formRef" @keyup.enter="handleLogin">
        <NFormItem label="用户名" path="userName">
          <NInput v-model:value="form.userName" placeholder="请输入用户名" />
        </NFormItem>
        <NFormItem label="密码" path="password">
          <NInput v-model:value="form.password" type="password" placeholder="请输入密码" />
        </NFormItem>
        <NFormItem>
          <NButton type="primary" native-type="submit" block :loading="loading" @click="handleLogin"> 登录 </NButton>
        </NFormItem>
      </NForm>
      <div v-if="appConfig.EnableRegister" class="v-flex-center">
        <span>没有账号？</span>
        <NButton quaternary type="info" @click="navigate('/register')">去注册</NButton>
      </div>
    </NCard>

    <!-- OTP验证弹窗 -->
    <NModal v-model:show="showOtpModal" preset="dialog" title="OTP验证" :mask-closable="false" :loading="loading">
      <NSpace justify="center">
        <div style="padding: 16px 0">
          <p style="margin-bottom: 16px">您的账号已启用OTP双因素认证</p>
          <NInputOtp v-model:value="otpCode" maxlength="6" @update:value="handleOtpLogin" />
        </div>
      </NSpace>
    </NModal>
  </div>
</template>
