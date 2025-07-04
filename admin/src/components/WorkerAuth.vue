<script setup lang="ts">
import { onMounted, ref } from 'vue'
import {
  NCard,
  NForm,
  NButton,
  useMessage,
  NFormItem,
  NInput,
  NModal,
  NTag,
  NList,
  NListItem,
  type FormInst,
  type FormRules,
} from 'naive-ui'
import type {
  AccessTokenCreateRequest,
  AccessTokenDeleteRequest,
  ExternalServerToken,
  InternalServerWhiteList,
  InternalWhiteListCreateRequest,
  InternalWhiteListDeleteRequest,
} from '@/types/access'
import {
  createAccessToken,
  createInternalWhiteList,
  deleteAccessToken,
  deleteInternalWhiteList,
  listAccessTokens,
  listInternalWhiteLists,
} from '@/api/workers'

const props = defineProps<{
  uid: string
}>()

const message = useMessage()

// 内部白名单列表
const internalWhiteLists = ref<InternalServerWhiteList[]>([])
const fetchInternalWhiteLists = async () => {
  try {
    const response = await listInternalWhiteLists({
      worker_uid: props.uid,
      page: 1,
      page_size: 100,
    })
    internalWhiteLists.value = response.data.internal_white_lists
  } catch (error) {
    console.error('Failed to fetch internal white lists', error)
    message.error('获取内部白名单列表失败')
  }
}

// 添加内部访问白名单
const showCreateInternalModal = ref<boolean>(false)
const IsCreatingInternal = ref<boolean>(false)
const createInternalForm = ref({
  allowWorkerName: '',
  description: '',
})
const createInternalFormRef = ref<FormInst | null>(null)
const createInternalRules: FormRules = {
  workerName: {
    required: true,
    message: '请输入授权访问的 Worker 名称',
  },
  description: {
    required: true,
    message: '请输入描述',
  },
}
const handleCreateInternalConfirm = async () => {
  if (!createInternalFormRef.value) return
  try {
    // 校验表单
    await createInternalFormRef.value.validate()
  } catch (error) {
    console.error(error)
    return
  }
  try {
    // 调用创建内部访问白名单接口
    IsCreatingInternal.value = true
    const request: InternalWhiteListCreateRequest = {
      worker_uid: props.uid,
      name: createInternalForm.value.allowWorkerName,
      description: createInternalForm.value.description,
    }
    await createInternalWhiteList(request)
    await fetchInternalWhiteLists()
    message.success('创建内部访问白名单成功')
    handleCreateInternalClose()
  } catch (error) {
    console.error(error)
    message.error('创建内部访问白名单失败: ' + error)
  } finally {
    IsCreatingInternal.value = false
  }
}
const handleCreateInternalClose = () => {
  showCreateInternalModal.value = false
  createInternalForm.value.allowWorkerName = ''
  createInternalForm.value.description = ''
}

// 访问密钥列表
const accessTokens = ref<ExternalServerToken[]>([])
const fetchAccessTokens = async () => {
  try {
    const response = await listAccessTokens({
      worker_uid: props.uid,
      page: 1,
      page_size: 100,
    })
    accessTokens.value = response.data.access_tokens
  } catch (error) {
    console.error('Failed to fetch access tokens', error)
    message.error('获取访问密钥列表失败')
  }
}

// 添加访问密钥
const showCreateTokenModal = ref<boolean>(false)
const IsCreatingToken = ref<boolean>(false)
const createTokenForm = ref({
  description: '',
})
const newToken = ref<string>('')
const showTokenModal = ref<boolean>(false)
const handleTokenModalClose = () => {
  console.log('close')
  showTokenModal.value = false
  newToken.value = ''
}
const createTokenFormRef = ref<FormInst | null>(null)
const createTokenRules: FormRules = {
  description: {
    required: true,
    message: '请输入描述',
  },
}
const handleCreateTokenConfirm = async () => {
  if (!createTokenFormRef.value) return
  try {
    // 校验表单
    await createTokenFormRef.value.validate()
  } catch (error) {
    console.error(error)
    return
  }
  try {
    // 调用创建内部访问白名单接口
    IsCreatingToken.value = true
    const request: AccessTokenCreateRequest = {
      worker_uid: props.uid,
      description: createTokenForm.value.description,
    }
    const response = await createAccessToken(request)
    newToken.value = response.data.access_token
    await fetchAccessTokens()
    message.success('创建Token成功')
    handleCreateTokenClose()
  } catch (error) {
    console.error(error)
    message.error('创建Token失败: ' + error)
  } finally {
    IsCreatingToken.value = false
  }
  showTokenModal.value = true
}
const handleCreateTokenClose = () => {
  showCreateTokenModal.value = false
  createTokenForm.value.description = ''
}

// 删除内部访问白名单或Token
const showDeleteModal = ref<boolean>(false)
const deleteType = ref<'internal' | 'token'>('internal')
const IsDeleting = ref<boolean>(false)
const idToDelete = ref<number>()
const handleDeleteClick = async (id: number, type: 'internal' | 'token') => {
  idToDelete.value = id
  deleteType.value = type
  showDeleteModal.value = true
}
const handleDeleteConfirm = async () => {
  if (idToDelete.value === undefined) return
  try {
    IsDeleting.value = true
    if (deleteType.value === 'internal') {
      const request: InternalWhiteListDeleteRequest = {
        worker_uid: props.uid,
        id: idToDelete.value,
      }
      await deleteInternalWhiteList(request)
      await fetchInternalWhiteLists()
      message.success('删除内部访问白名单成功')
    } else if (deleteType.value === 'token') {
      const request: AccessTokenDeleteRequest = {
        worker_uid: props.uid,
        id: idToDelete.value,
      }
      await deleteAccessToken(request)
      await fetchAccessTokens()
      message.success('删除Token成功')
    }
  } catch (error) {
    console.error(error)
    message.error('删除失败: ' + error)
  } finally {
    IsDeleting.value = false
  }
}
const handleDeleteClose = () => {
  showDeleteModal.value = false
  idToDelete.value = undefined
  deleteType.value = 'internal'
}

onMounted(async () => {
  await fetchInternalWhiteLists()
  await fetchAccessTokens()
})
</script>

<template>
  <div class="auth-container">
    <NCard title="内部访问" :bordered="false">
      <template #header-extra>
        <NButton type="primary" secondary @click="showCreateInternalModal = true">新增</NButton>
      </template>
      <NList bordered>
        <NListItem v-for="item in internalWhiteLists" :key="item.ID">
          <span>{{ item.WorkerName }}</span>
          <template #suffix>
            <NButton quaternary type="primary" @click="handleDeleteClick(item.ID, 'internal')">
              删除
            </NButton>
          </template>
        </NListItem>
      </NList>
    </NCard>
    <NCard title="访问密钥" :bordered="false">
      <template #header-extra>
        <NButton type="primary" secondary @click="showCreateTokenModal = true">新增</NButton>
      </template>
      <NList bordered>
        <NListItem v-for="item in accessTokens" :key="item.ID">
          <div>{{ item.token }}</div>
          <div>描述：{{ item.description }}</div>
          <template #suffix>
            <NButton quaternary type="primary" @click="handleDeleteClick(item.ID, 'token')">
              删除
            </NButton>
          </template>
        </NListItem>
      </NList>
    </NCard>
    <NModal
      v-model:show="showCreateInternalModal"
      preset="dialog"
      title="添加内部访问白名单"
      positive-text="确认"
      negative-text="取消"
      :loading="IsCreatingInternal"
      :mask-closable="false"
      @positive-click="handleCreateInternalConfirm"
      @negative-click="handleCreateInternalClose"
    >
      <NForm ref="createInternalFormRef" :rules="createInternalRules" :model="createInternalForm">
        <NFormItem label="授权访问的 Worker 名称">
          <NInput
            v-model:value="createInternalForm.allowWorkerName"
            placeholder="请输入授权访问的 Worker 名称"
          />
        </NFormItem>
        <NFormItem label="描述">
          <NInput v-model:value="createInternalForm.description" placeholder="请输入描述" />
        </NFormItem>
      </NForm>
    </NModal>
    <NModal
      v-model:show="showCreateTokenModal"
      preset="dialog"
      title="添加访问密钥"
      positive-text="确认"
      negative-text="取消"
      :loading="IsCreatingToken"
      :mask-closable="false"
      @positive-click="handleCreateTokenConfirm"
      @negative-click="handleCreateTokenClose"
    >
      <NForm ref="createTokenFormRef" :rules="createTokenRules" :model="createTokenForm">
        <NFormItem label="描述">
          <NInput v-model:value="createTokenForm.description" placeholder="请输入描述" />
        </NFormItem>
      </NForm>
    </NModal>
    <NModal
      v-model:show="showTokenModal"
      preset="dialog"
      title="重要提示"
      positive-text="确认"
      @positive-click="handleTokenModalClose"
      :mask-closable="false"
    >
      <div>此 token 仅显示一次，请妥善保存。</div>
      <NTag class="v-item-column" style="width: 100%">{{ newToken }}</NTag>
    </NModal>
    <NModal
      v-model:show="showDeleteModal"
      preset="dialog"
      :title="deleteType === 'internal' ? '删除内部访问白名单' : '删除Token'"
      positive-text="删除"
      negative-text="取消"
      :loading="IsDeleting"
      :mask-closable="false"
      @positive-click="handleDeleteConfirm"
      @negative-click="handleDeleteClose"
    >
      <div>确认要删除 {{ deleteType === 'internal' ? '内部访问白名单' : 'Token' }} 吗？</div>
    </NModal>
  </div>
</template>

<style scoped>
.auth-container {
  /* 运算符前后添加空格 */
  height: calc(100vh - 300px);
  overflow-y: auto;
}
</style>
