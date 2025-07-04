<script setup lang="ts">
import { onMounted, ref } from 'vue'
import {
  NCard,
  NList,
  NListItem,
  NButton,
  useMessage,
  NTag,
  NInput,
  NModal,
  NForm,
  NFormItem,
  type FormInst,
  type FormRules,
} from 'naive-ui'
import type { ResourceData } from '@/types/resources'
import { createResource, deleteResource, getResourceList } from '@/api/resources'

const message = useMessage()
const kvs = ref<ResourceData[]>([])

// 加载所有 KV
const loadKvs = async () => {
  try {
    const response = await getResourceList(0, 10000, 'kv')
    kvs.value = response.data
  } catch (error) {
    console.error(error)
    message.error('获取 KV 列表失败: ' + error)
  }
}

// 创建 KV
const showCreateKVModal = ref<boolean>(false)
const IsCreatingKV = ref<boolean>(false)
const createKVForm = ref({
  name: '',
})
const createKVFormRef = ref<FormInst | null>(null)
const createKVRules: FormRules = {
  name: {
    required: true,
    message: '请输入 KV 名称',
  },
}
const handleCreateKVConfirm = async () => {
  if (!createKVFormRef.value) return
  try {
    // 校验表单
    await createKVFormRef.value.validate()
  } catch (error) {
    console.error(error)
    return
  }
  try {
    // 调用创建 KV 接口
    IsCreatingKV.value = true
    const newKV = await createResource(createKVForm.value.name, 'kv')
    kvs.value.push(newKV)
    message.success('创建 KV 成功')
    handleCreateKVClose()
  } catch (error) {
    console.error(error)
    message.error('创建 KV 失败: ' + error)
  } finally {
    IsCreatingKV.value = false
  }
}
const handleCreateKVClose = () => {
  showCreateKVModal.value = false
  createKVForm.value.name = ''
}

// 删除 KV
const showDeleteKVModal = ref<boolean>(false)
const IsDeletingKV = ref<boolean>(false)
const kvToDelete = ref<string>('')
const handleDeleteKVClick = async (uid: string) => {
  kvToDelete.value = uid
  showDeleteKVModal.value = true
}
const handleDeleteKVConfirm = async () => {
  if (!kvToDelete.value) return
  try {
    IsDeletingKV.value = true
    await deleteResource(kvToDelete.value, 'kv')
    kvs.value = kvs.value.filter((kv) => kv.uid !== kvToDelete.value)
    message.success('删除 KV 成功')
    handleDeleteKVClose()
  } catch (error) {
    console.error(error)
    message.error('删除 KV 失败: ' + error)
  } finally {
    IsDeletingKV.value = false
  }
}
const handleDeleteKVClose = () => {
  showDeleteKVModal.value = false
  kvToDelete.value = ''
}

onMounted(async () => {
  loadKvs()
})
</script>
<template>
  <div class="v-main">
    <NCard title="KV" :bordered="false">
      <template #header-extra>
        <NButton type="primary" secondary @click="showCreateKVModal = true">创建</NButton>
      </template>
      <NList bordered>
        <NListItem v-for="item in kvs" :key="item.uid">
          <template #prefix>
            <div style="width: 200px">
              <NTag size="large">{{ item.name }}</NTag>
            </div>
          </template>
          <template #suffix>
            <NButton quaternary type="primary" @click="handleDeleteKVClick(item.uid)">
              删除
            </NButton>
          </template>
          <div class="v-item">ID: {{ item.uid }}</div>
        </NListItem>
      </NList>
    </NCard>
    <NModal
      v-model:show="showCreateKVModal"
      preset="dialog"
      title="创建 KV"
      positive-text="确认"
      negative-text="取消"
      :loading="IsCreatingKV"
      :mask-closable="false"
      @positive-click="handleCreateKVConfirm"
      @negative-click="handleCreateKVClose"
    >
      <NForm :model="createKVForm" :rules="createKVRules" ref="createKVFormRef">
        <NFormItem label="KV 名称" path="name">
          <NInput v-model:value="createKVForm.name" placeholder="请输入 KV 名称" />
        </NFormItem>
      </NForm>
    </NModal>
    <NModal
      v-model:show="showDeleteKVModal"
      preset="dialog"
      title="删除 KV"
      positive-text="确认"
      negative-text="取消"
      :loading="IsDeletingKV"
      :mask-closable="false"
      @positive-click="handleDeleteKVConfirm"
      @negative-click="handleDeleteKVClose"
    >
      <div>确认要删除 KV {{ kvToDelete }}？</div>
    </NModal>
  </div>
</template>
