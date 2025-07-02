<script setup lang="ts">
import { inject, ref, type Ref } from 'vue'
import { useMessage, NCard, NForm, NFormItem, NInput, NButton } from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui' // 导入 FormInst 类型
import type { LoginRequest, UserInfo } from '@/types/auth'
import { passwordRules, usernameRules } from '@/constant/formrules'
import { getUserInfo, login } from '@/api/auth'
import { useNavigate } from '@/composables/useNavigate'

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

// 处理登录逻辑的函数
const handleLogin = async () => {
  if (!formRef.value) return
  try {
    // 校验表单
    await formRef.value.validate()
  } catch (error) {
    console.error(error)
    return
  }
  try {
    // 调用登录接口
    const res = await login(form.value)
    if (res.status === 0) {
      // 登录成功获取用户信息
      try {
        userInfo.value = await getUserInfo()
      } catch (error) {
        console.error(error)
        message.error('登录成功后获取用户信息失败: ' + error)
      }
      message.success('登录成功')
      navigate('/workers')
    }
  } catch (error) {
    console.error(error)
    message.error('登录失败: ' + error)
  }
}
</script>
<template>
  <div class="v-base-page v-flex-center">
    <NCard style="width: 400px">
      <template #header>
        <div class="v-card-header">登录</div>
      </template>
      <NForm :model="form" :rules="rules" ref="formRef">
        <NFormItem label="用户名" path="userName">
          <NInput v-model:value="form.userName" placeholder="请输入用户名" />
        </NFormItem>
        <NFormItem label="密码" path="password">
          <NInput v-model:value="form.password" type="password" placeholder="请输入密码" />
        </NFormItem>
        <NFormItem>
          <NButton type="primary" native-type="submit" block @click="handleLogin"> 登录 </NButton>
        </NFormItem>
      </NForm>
    </NCard>
  </div>
</template>
