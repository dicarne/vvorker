<script setup lang="ts">
import { ref } from 'vue'
import { useMessage, NCard, NForm, NFormItem, NInput, NButton } from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui' // 导入 FormInst 类型
import type { RegisterRequest } from '@/types/auth'
import { emailRules, passwordRules, usernameRules } from '@/constant/formrules'
import { register } from '@/api/auth'
import { useNavigate } from '@/composables/useNavigate'

// 初始化表单数据
const form = ref<RegisterRequest>({
  userName: '',
  email: '',
  password: '',
})

// 表单引用，明确指定类型为 FormInst 或 null
const formRef = ref<FormInst | null>(null)

const rules: FormRules = {
  userName: usernameRules,
  password: passwordRules,
  email: emailRules,
}

// 消息提示实例
const message = useMessage()

const { navigate } = useNavigate()

// 处理注册逻辑的函数
const handleRegister = async () => {
  if (!formRef.value) return
  try {
    // 校验表单
    await formRef.value.validate()
  } catch (error) {
    console.error(error)
    return
  }
  try {
    // 调用注册接口
    const res = await register(form.value)
    if (res.status === 0) {
      message.success('注册成功')
      navigate('/login')
    }
  } catch (error) {
    console.error(error)
    message.error('注册失败: ' + error)
  }
}
</script>
<template>
  <div class="v-base-page v-flex-center">
    <NCard style="width: 400px">
      <template #header>
        <div class="v-card-header">注册</div>
      </template>
      <NForm :model="form" :rules="rules" ref="formRef">
        <NFormItem label="用户名" path="userName">
          <NInput v-model:value="form.userName" placeholder="请输入用户名" />
        </NFormItem>
        <NFormItem label="密码" path="password">
          <NInput v-model:value="form.password" type="password" placeholder="请输入密码" />
        </NFormItem>
        <NFormItem label="邮箱" path="email">
          <NInput v-model:value="form.email" placeholder="请输入邮箱" />
        </NFormItem>
        <NFormItem>
          <NButton type="primary" native-type="submit" block @click="handleRegister">
            注册
          </NButton>
        </NFormItem>
      </NForm>
    </NCard>
  </div>
</template>
