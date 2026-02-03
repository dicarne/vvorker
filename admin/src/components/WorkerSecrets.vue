<script setup lang="ts">
import { onMounted, ref } from 'vue'
import {
  NForm,
  NButton,
  useMessage,
  NFormItem,
  NInput,
  NModal,
  NSwitch,
  NTable,
  type FormInst,
  type FormRules,
} from 'naive-ui'
import {
  listWorkerSecrets,
  createWorkerSecret,
  updateWorkerSecret,
  deleteWorkerSecret,
} from '@/api/workers'
import type {
  DeleteSecretRequest,
  Secret,
  SecretCreateRequest,
  UpdateSecretRequest,
} from '@/types/access'

const props = defineProps<{
  uid: string
}>()

const message = useMessage()

const secrets = ref<Secret[]>([])
const fetchSecrets = async () => {
  try {
    const response = await listWorkerSecrets(props.uid)
    secrets.value = response.data.secrets
  } catch (error) {
    console.error('listWorkerSecrets Error', error)
    message.error('获取变量失败')
  }
}

// 添加 secret
const showCreateSecretModal = ref<boolean>(false)
const IsCreatingSecret = ref<boolean>(false)
const createSecretForm = ref({
  key: '',
  value: '',
})
const createSecretFormRef = ref<FormInst | null>(null)
const workerSecretRules: FormRules = {
  key: {
    required: true,
    message: '请输入变量名',
  },
  value: {
    required: true,
    message: '请输入变量值',
  },
}
const handleCreateSecretConfirm = async () => {
  if (!createSecretFormRef.value) return
  try {
    // 校验表单
    await createSecretFormRef.value.validate()
  } catch (error) {
    console.error('createSecretFormRef validate Error', error)
    return
  }
  try {
    // 调用创建变量接口
    IsCreatingSecret.value = true
    const request: SecretCreateRequest = {
      worker_uid: props.uid,
      key: createSecretForm.value.key,
      value: createSecretForm.value.value,
    }
    await createWorkerSecret(request)
    await fetchSecrets()
    message.success('创建变量成功')
    handleCreateSecretClose()
  } catch (error) {
    console.error('createWorkerSecret Error', error)
    message.error('创建变量失败')
  } finally {
    IsCreatingSecret.value = false
  }
}
const handleCreateSecretClose = () => {
  showCreateSecretModal.value = false
  createSecretForm.value.key = ''
  createSecretForm.value.value = ''
}

// 编辑 secret
const showEditSecretModal = ref<boolean>(false)
const IsEditingSecret = ref<boolean>(false)
const editSecretForm = ref({
  key: '',
  value: '',
})
const editSecretFormRef = ref<FormInst | null>(null)
const secretUidToEdit = ref<number>(0)
const handleEditSecretClick = async (item: Secret) => {
  secretUidToEdit.value = item.ID
  editSecretForm.value.key = item.Key
  editSecretForm.value.value = item.Value
  showEditSecretModal.value = true
}
const handleEditSecretConfirm = async () => {
  if (!secretUidToEdit.value) return
  try {
    IsEditingSecret.value = true
    const request: UpdateSecretRequest = {
      worker_uid: props.uid,
      id: secretUidToEdit.value,
      key: editSecretForm.value.key,
      value: editSecretForm.value.value,
    }
    await updateWorkerSecret(request)
    await fetchSecrets()
    message.success('编辑变量成功')
    handleEditSecretClose()
  } catch (error) {
    console.error('updateWorkerSecret Error', error)
    message.error('编辑变量失败')
  } finally {
    IsEditingSecret.value = false
  }
}
const handleEditSecretClose = () => {
  showEditSecretModal.value = false
  secretUidToEdit.value = 0
}

// 删除 secret
const showDeleteSecretModal = ref<boolean>(false)
const IsDeletingSecret = ref<boolean>(false)
const secretUidToDelete = ref<number>(0)
const handleDeleteSecretClick = async (id: number) => {
  secretUidToDelete.value = id
  showDeleteSecretModal.value = true
}
const handleDeleteSecretConfirm = async () => {
  if (!secretUidToDelete.value) return
  try {
    IsDeletingSecret.value = true
    const request: DeleteSecretRequest = {
      worker_uid: props.uid,
      id: secretUidToDelete.value,
    }
    await deleteWorkerSecret(request)
    await fetchSecrets()
    message.success('删除变量成功')
    handleDeleteSecretClose()
  } catch (error) {
    console.error('deleteWorkerSecret Error', error)
    message.error('删除变量失败')
  } finally {
    IsDeletingSecret.value = false
  }
}
const handleDeleteSecretClose = () => {
  showDeleteSecretModal.value = false
  secretUidToDelete.value = 0
}

onMounted(async () => {
  await fetchSecrets()
})
</script>

<template>
  <div>
    <NButton style="margin: 8px 0" type="primary" secondary @click="showCreateSecretModal = true">
      添加变量
    </NButton>
    <NTable :bordered="false" :single-line="false">
      <thead>
        <tr>
          <th>变量名</th>
          <th>变量值</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in secrets" :key="item.ID">
          <td>{{ item.Key }}</td>
          <td>{{ item.Value }}</td>
          <td>
            <NButton quaternary type="primary" @click="handleEditSecretClick(item)"> 编辑 </NButton>
            <NButton quaternary type="primary" @click="handleDeleteSecretClick(item.ID)">
              删除
            </NButton>
          </td>
        </tr>
      </tbody>
    </NTable>
    <NModal
      v-model:show="showCreateSecretModal"
      preset="dialog"
      title="添加变量"
      positive-text="确认"
      negative-text="取消"
      :loading="IsCreatingSecret"
      :mask-closable="false"
      @positive-click="handleCreateSecretConfirm"
      @negative-click="handleCreateSecretClose"
    >
      <NForm ref="createSecretFormRef" :rules="workerSecretRules" :model="createSecretForm">
        <NFormItem label="Key">
          <NInput v-model:value="createSecretForm.key" placeholder="请输入Key" />
        </NFormItem>
        <NFormItem label="Value">
          <NInput v-model:value="createSecretForm.value" placeholder="请输入Value" />
        </NFormItem>
      </NForm>
    </NModal>
    <NModal
      v-model:show="showEditSecretModal"
      preset="dialog"
      title="编辑变量"
      positive-text="确认"
      negative-text="取消"
      :loading="IsEditingSecret"
      :mask-closable="false"
      @positive-click="handleEditSecretConfirm"
      @negative-click="handleEditSecretClose"
    >
      <NForm ref="editSecretFormRef" :rules="workerSecretRules" :model="editSecretForm">
        <NFormItem label="Key">
          <NInput v-model:value="editSecretForm.key" placeholder="请输入Key" />
        </NFormItem>
        <NFormItem label="Value">
          <NInput v-model:value="editSecretForm.value" placeholder="请输入Value" />
        </NFormItem>
      </NForm>
    </NModal>
    <NModal
      v-model:show="showDeleteSecretModal"
      preset="dialog"
      title="删除变量"
      positive-text="确认"
      negative-text="取消"
      :loading="IsDeletingSecret"
      :mask-closable="false"
      @positive-click="handleDeleteSecretConfirm"
      @negative-click="handleDeleteSecretClose"
    >
      <div>确认要删除此变量？</div>
    </NModal>
  </div>
</template>

<style scoped></style>
