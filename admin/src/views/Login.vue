<script setup lang="ts">
import { ref } from 'vue'
import { useMessage, NCard, NForm, NFormItem, NInput, NButton } from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui' // 导入 FormInst 类型
import type { LoginRequest } from '@/types/auth'
import { passwordRules, usernameRules } from '@/constant/formrules'
import { login } from '@/api/auth'
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
    console.log(res)
    message.success('登录成功')
    navigate('/admin')
  } catch (error) {
    console.error(error)
    message.error('登录失败，请检查输入信息')
  }
}
</script>
<template>
  <div class="login-container">
    <NCard class="login-card">
      <template #header>
        <div class="login-title">登录</div>
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

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100vw;
  height: 100vh;
}

.login-card {
  width: 400px;
}

.login-title {
  text-align: center;
  font-size: 20px;
  font-weight: bold;
}
</style>
