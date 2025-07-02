<script setup lang="ts">
import { inject, ref, type Ref } from 'vue'
import {
  NCard,
  NList,
  NListItem,
  NButton,
  NModal,
  NInput,
  useMessage,
  NForm,
  NFormItem,
} from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui' // 导入 FormInst 类型
import { passwordRules } from '@/constant/formrules'
import { changePassword } from '@/api/users'
import type { UserInfo } from '@/types/auth'

const userInfo = inject<Ref<UserInfo>>('userInfo')!
const message = useMessage()

const showChangePasswordModal = ref<boolean>(false)
const pswdForm = ref({
  password: '',
})
const pswdFormRef = ref<FormInst | null>(null)
const pswdRules: FormRules = {
  password: passwordRules,
}
const handleChangePasswordConfirm = async () => {
  if (!pswdFormRef.value) return
  try {
    // 校验表单
    await pswdFormRef.value.validate()
  } catch (error) {
    console.error(error)
    return
  }
  try {
    // 调用修改密码接口
    const res = await changePassword(userInfo.value.id, pswdForm.value.password)
    message.success('修改密码成功')
    handleChangePasswordClose()
  } catch (error) {
    console.error(error)
    message.error('修改密码失败: ' + error)
  }
}
const handleChangePasswordClose = () => {
  showChangePasswordModal.value = false
  pswdForm.value.password = ''
}
</script>
<template>
  <div class="v-main">
    <NCard title="Config">
      <div class="v-flex-between-center">
        <div>用户密码</div>
        <NButton secondary type="primary" @click="showChangePasswordModal = true">修改密码</NButton>
        <NModal
          v-model:show="showChangePasswordModal"
          preset="dialog"
          title="修改密码"
          positive-text="确认"
          negative-text="取消"
          :mask-closable="false"
          @positive-click="handleChangePasswordConfirm"
          @negative-click="handleChangePasswordClose"
        >
          <NForm :model="pswdForm" :rules="pswdRules" ref="pswdFormRef">
            <NFormItem label="新密码" path="password">
              <NInput
                v-model:value="pswdForm.password"
                type="password"
                placeholder="请输入新密码"
              />
            </NFormItem>
          </NForm>
        </NModal>
      </div>
    </NCard>
    <NCard title="Access Key" style="margin-top: 16px">
      <template #header-extra>
        <NButton type="primary" secondary>创建</NButton>
      </template>
      <NList>
        <NListItem>
          <template #prefix>
            <div style="min-width: 400px">Name:</div>
          </template>
          <template #suffix>
            <NButton type="error" secondary>删除</NButton>
          </template>
          <div class="v-item" style="min-width: 400px">ID:</div>
        </NListItem>
      </NList>
    </NCard>
  </div>
</template>
