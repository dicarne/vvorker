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
  NQrCode,
  NSwitch,
  NSpace,
} from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui' // 导入 FormInst 类型
import { passwordRules } from '@/constant/formrules'
import type { UserInfo } from '@/types/auth'
import type { AccessKey } from '@/types/access'
import { Copy24Regular as CopyIcon } from '@vicons/fluent'
import { changePassword } from '@/api/users'
import { useCopyContent } from '@/composables/useUtils'
import {
  createAccessKey,
  deleteAccessKey,
  disableOtp,
  enableOtp,
  getAccessKeys,
  isEnableOtp,
  validOtp,
} from '@/api/auth'

const userInfo = inject<Ref<UserInfo>>('userInfo')!
const message = useMessage()
const { copyContent } = useCopyContent()

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
    console.error('pswdFormRef validate Error', error)
    return
  }
  try {
    // 调用修改密码接口
    isChangingPassword.value = true
    const res = await changePassword(userInfo.value.id, pswdForm.value.password)
    message.success('修改密码成功')
    handleChangePasswordClose()
  } catch (error) {
    console.error('changePassword Error', error)
    message.error('修改密码失败')
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
    console.error('getAccessKeys Error', error)
    message.error('获取 Access Key 列表失败')
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
    console.error('createAccessKeyFormRef validate Error', error)
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
    console.error('createAccessKey Error', error)
    message.error('创建 Access Key 失败')
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
    message.success('删除 API 密钥成功')
    handleDeleteAccessKeyClose()
  } catch (error) {
    console.error('deleteAccessKey Error', error)
    message.error('删除 API 密钥失败')
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
  const res = await isEnableOtp()
  otpEnabled.value = res.data.data.enabled
})

const otpUrl = ref<string>('')
const otpEnabled = ref<boolean>(false)
const handleOtpSwitchChange = async (value: boolean) => {
  if (value) {
    try {
      const res = await enableOtp()
      otpUrl.value = res.data.data.url
    } catch (error) {
      console.error('enableOtp Error', error)
      message.error('启用 OTP 失败')
      otpEnabled.value = false
    }
  } else {
    try {
      await disableOtp()
      otpUrl.value = ''
    } catch (error) {
      console.error('disableOtp Error', error)
      message.error('禁用 OTP 失败')
      otpEnabled.value = true
    }
  }
}

const optAddCode = ref<string>('')
const handleOptAddCodeConfirm = async () => {
  if (!optAddCode.value) {
    message.error('请输入 OTP 代码')
    return
  }
  try {
    const res = await validOtp(optAddCode.value)
    if (res.data.code === 0) {
      otpEnabled.value = true
      otpUrl.value = ''
      message.success('添加 OTP 成功')
      otpUrl.value = ''
    } else {
      message.error('添加 OTP 失败')
    }
  } catch (error) {
    console.error('validOtp Error', error)
    message.error('验证失败')
  }
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  try {
    const date = new Date(dateStr)
    if (isNaN(date.getTime())) {
      return dateStr
    }
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  } catch (error) {
    return dateStr
  }
}
</script>
<template>
  <div class="v-main">
    <NCard title="配置">
      <div class="v-flex-between-center">
        <div>用户密码</div>
        <NButton secondary type="primary" @click="showChangePasswordModal = true">修改密码</NButton>
        <NModal v-model:show="showChangePasswordModal" preset="dialog" title="修改密码" positive-text="确认"
          negative-text="取消" :mask-closable="false" :loading="isChangingPassword"
          @positive-click="handleChangePasswordConfirm" @negative-click="handleChangePasswordClose">
          <NForm :model="pswdForm" :rules="pswdRules" ref="pswdFormRef">
            <NFormItem label="新密码" path="password">
              <NInput v-model:value="pswdForm.password" type="password" placeholder="请输入新密码" />
            </NFormItem>
          </NForm>
        </NModal>
      </div>
    </NCard>
    <NCard title="API 密钥" style="margin-top: 16px">
      <template #header-extra>
        <NButton type="primary" secondary @click="showCreateAccessKeyModal = true">创建密钥</NButton>
        <NModal v-model:show="showCreateAccessKeyModal" preset="dialog" title="创建密钥" positive-text="确认"
          negative-text="取消" :loading="IsCreatingAccessKey" :mask-closable="false"
          @positive-click="handleCreateAccessKeyConfirm" @negative-click="handleCreateAccessKeyClose">
          <NForm :model="createAccessKeyForm" :rules="createAccessKeyRules" ref="createAccessKeyFormRef" class="mt-8">
            <NFormItem label="密钥名称" path="accessKeyName">
              <NInput v-model:value="createAccessKeyForm.accessKeyName" placeholder="请输入密钥名称" />
            </NFormItem>
          </NForm>
        </NModal>
      </template>
      <NList v-if="accessKeys.length > 0">
        <NListItem v-for="item in accessKeys" :key="item.key">
          <div class="access-key-item">
            <div class="access-key-content">
              <div class="access-key-name">
                {{ item.name }}
              </div>
              <div class="access-key-value">
                <NInput :value="item.key" readonly class="access-key-input" />
              </div>
              <div class="access-key-actions">
                <NButton quaternary type="primary" size="small" @click="copyContent(item.key)">
                  <template #icon>
                    <CopyIcon />
                  </template>
                  复制
                </NButton>
                <NButton type="error" secondary size="small" @click="handleDeleteAccessKeyClick(item.key)">删除</NButton>
              </div>
            </div>
            <div v-if="item.created_at" class="access-key-time">创建于 {{ formatDate(item.created_at) }}</div>
          </div>
        </NListItem>
      </NList>
      <NEmpty v-else description="暂无 API 密钥" />
    </NCard>
    <NCard title="OTP" style="margin-top: 16px">
      <NSpace>
        <div>启用OTP</div>
        <NSwitch v-model:value="otpEnabled" @update:value="handleOtpSwitchChange" />
      </NSpace>
      <NSpace v-if="otpUrl">
        <NQrCode :value="otpUrl" />
        <NInput v-model:value="optAddCode" placeholder="请输入 OTP 代码" />
        <NButton type="primary" secondary @click="handleOptAddCodeConfirm">验证并添加认证器</NButton>
      </NSpace>
    </NCard>
    <NModal v-model:show="showDeleteAccessKeyModal" preset="dialog" title="删除 API 密钥" positive-text="确认"
      negative-text="取消" :loading="IsDeletingAccessKey" :mask-closable="false"
      @positive-click="handleDeleteAccessKeyConfirm" @negative-click="handleDeleteAccessKeyClose">
      <div>确认要删除这个 API 密钥？</div>
    </NModal>
  </div>
</template>
<style scoped>
.access-key-item {
  width: 100%;
  padding: 8px 0;
}

.access-key-content {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.access-key-name {
  min-width: 150px;
  font-weight: 500;
  color: #333;
}

.access-key-value {
  flex: 1;
  min-width: 200px;
}

.access-key-input {
  width: 100%;
  font-family: 'Courier New', monospace;
}

.access-key-actions {
  display: flex;
  gap: 8px;
}

.access-key-time {
  font-size: 12px;
  color: #999;
  margin-top: 8px;
}
</style>
