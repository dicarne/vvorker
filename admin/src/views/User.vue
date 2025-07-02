<script setup lang="ts">
import { inject, onMounted, ref, type Ref } from 'vue'
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
import type { UserInfo } from '@/types/auth'
import type { AccessKey } from '@/types/user'
import { changePassword } from '@/api/users'
import { createAccessKey, deleteAccessKey, getAccessKeys } from '@/api/auth'

const userInfo = inject<Ref<UserInfo>>('userInfo')!
const message = useMessage()

// Config 相关
// 修改密码
const showChangePasswordModal = ref<boolean>(false)
const isChangingPassword = ref<boolean>(false)
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
    isChangingPassword.value = true
    const res = await changePassword(userInfo.value.id, pswdForm.value.password)
    message.success('修改密码成功')
    handleChangePasswordClose()
  } catch (error) {
    console.error(error)
    message.error('修改密码失败: ' + error)
  } finally {
    isChangingPassword.value = false
  }
}
const handleChangePasswordClose = () => {
  showChangePasswordModal.value = false
  pswdForm.value.password = ''
}

// Access Key 相关
const accessKeys = ref<AccessKey[]>([])
// 加载所有 Access Key
const loadAccessKeys = async () => {
  try {
    const data = await getAccessKeys()
    accessKeys.value = data.data.data
  } catch (error) {
    console.error(error)
    message.error('获取 Access Key 列表失败: ' + error)
  }
}

// 创建 Access Key
const showCreateAccessKeyModal = ref<boolean>(false)
const IsCreatingAccessKey = ref<boolean>(false)
const createAccessKeyForm = ref({
  accessKeyName: '',
})
const createAccessKeyFormRef = ref<FormInst | null>(null)
const createAccessKeyRules: FormRules = {
  accessKeyName: {
    required: true,
    message: '请输入 Access Key 名称',
  },
}
const handleCreateAccessKeyConfirm = async () => {
  if (!createAccessKeyFormRef.value) return
  try {
    // 校验表单
    await createAccessKeyFormRef.value.validate()
  } catch (error) {
    console.error(error)
    return
  }
  try {
    // 调用创建 Access Key 接口
    IsCreatingAccessKey.value = true
    const newKey = await createAccessKey(createAccessKeyForm.value.accessKeyName)
    accessKeys.value.push(newKey.data.data as AccessKey)
    message.success('创建 Access Key 成功')
    handleCreateAccessKeyClose()
  } catch (error) {
    console.error(error)
    message.error('创建 Access Key 失败: ' + error)
  } finally {
    IsCreatingAccessKey.value = false
  }
}
const handleCreateAccessKeyClose = () => {
  showCreateAccessKeyModal.value = false
  createAccessKeyForm.value.accessKeyName = ''
}

// 删除 Access Key
const showDeleteAccessKeyModal = ref<boolean>(false)
const IsDeletingAccessKey = ref<boolean>(false)
const accessKeyToDelete = ref<string>('')
const handleDeleteAccessKeyClick = (id: string) => {
  accessKeyToDelete.value = id
  showDeleteAccessKeyModal.value = true
}
const handleDeleteAccessKeyConfirm = async () => {
  try {
    // 调用创建 Access Key 接口
    IsDeletingAccessKey.value = true
    await deleteAccessKey(accessKeyToDelete.value)
    accessKeys.value = accessKeys.value.filter((key) => key.key !== accessKeyToDelete.value)
    message.success('删除 Access Key 成功')
    handleDeleteAccessKeyClose()
  } catch (error) {
    console.error(error)
    message.error('删除 Access Key 失败: ' + error)
  } finally {
    IsDeletingAccessKey.value = false
  }
}
const handleDeleteAccessKeyClose = () => {
  showDeleteAccessKeyModal.value = false
  accessKeyToDelete.value = ''
}

onMounted(async () => {
  loadAccessKeys()
})
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
          :loading="isChangingPassword"
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
        <NButton type="primary" secondary @click="showCreateAccessKeyModal = true">创建</NButton>
        <NModal
          v-model:show="showCreateAccessKeyModal"
          preset="dialog"
          title="创建 Access Key"
          positive-text="确认"
          negative-text="取消"
          :loading="IsCreatingAccessKey"
          :mask-closable="false"
          @positive-click="handleCreateAccessKeyConfirm"
          @negative-click="handleCreateAccessKeyClose"
        >
          <NForm
            :model="createAccessKeyForm"
            :rules="createAccessKeyRules"
            ref="createAccessKeyFormRef"
          >
            <NFormItem label="Access Key 名称" path="accessKeyName">
              <NInput
                v-model:value="createAccessKeyForm.accessKeyName"
                placeholder="请输入 Access Key 名称"
              />
            </NFormItem>
          </NForm>
        </NModal>
      </template>
      <NList>
        <NListItem v-for="item in accessKeys" :key="item.key">
          <template #prefix>
            <div style="min-width: 400px">Key Name: {{ item.name }}</div>
          </template>
          <template #suffix>
            <NButton type="error" secondary @click="handleDeleteAccessKeyClick(item.key)"
              >删除</NButton
            >
          </template>
          <div class="v-item" style="min-width: 400px">Key ID: {{ item.key }}</div>
        </NListItem>
      </NList>
    </NCard>
    <NModal
      v-model:show="showDeleteAccessKeyModal"
      preset="dialog"
      title="删除 Access Key"
      positive-text="确认"
      negative-text="取消"
      :loading="IsDeletingAccessKey"
      :mask-closable="false"
      @positive-click="handleDeleteAccessKeyConfirm"
      @negative-click="handleDeleteAccessKeyClose"
    >
      <div>确认要删除这个 Access Key ？</div>
    </NModal>
  </div>
</template>
