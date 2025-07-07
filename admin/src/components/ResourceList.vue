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

// 定义组件接收的属性
const props = defineProps<{
  rType: string
}>()

const message = useMessage()
const resources = ref<ResourceData[]>([])

// 加载所有资源
const loadResources = async () => {
  try {
    const response = await getResourceList(0, 10000, props.rType)
    resources.value = response.data
  } catch (error) {
    console.error('getResourceList Error', error)
    message.error(`获取 ${props.rType} 列表失败`)
  }
}

// 创建资源
const showCreateModal = ref<boolean>(false)
const IsCreating = ref<boolean>(false)
const createForm = ref({
  name: '',
})
const createFormRef = ref<FormInst | null>(null)
const createRules: FormRules = {
  name: {
    required: true,
    message: `请输入 ${props.rType} 名称`,
  },
}
const handleCreateConfirm = async () => {
  if (!createFormRef.value) return
  try {
    // 校验表单
    await createFormRef.value.validate()
  } catch (error) {
    console.error('createFormRef validate Error', error)
    return
  }
  try {
    // 调用创建接口
    IsCreating.value = true
    const newResource = await createResource(createForm.value.name, props.rType)
    resources.value.push(newResource)
    message.success(`创建 ${props.rType} 成功`)
    handleCreateClose()
  } catch (error) {
    console.error(`createResource ${props.rType} Error`, error)
    message.error(`创建 ${props.rType} 失败`)
  } finally {
    IsCreating.value = false
  }
}
const handleCreateClose = () => {
  showCreateModal.value = false
  createForm.value.name = ''
}

// 删除资源
const showDeleteModal = ref<boolean>(false)
const IsDeleting = ref<boolean>(false)
const resourceToDelete = ref<string>('')
const handleDeleteClick = async (uid: string) => {
  resourceToDelete.value = uid
  showDeleteModal.value = true
}
const handleDeleteConfirm = async () => {
  if (!resourceToDelete.value) return
  try {
    IsDeleting.value = true
    await deleteResource(resourceToDelete.value, props.rType)
    resources.value = resources.value.filter((kv) => kv.uid !== resourceToDelete.value)
    message.success(`删除 ${props.rType} 成功`)
    handleDeleteClose()
  } catch (error) {
    console.error(`deleteResource ${props.rType} Error`, error)
    message.error(`删除 ${props.rType} 失败`)
  } finally {
    IsDeleting.value = false
  }
}
const handleDeleteClose = () => {
  showDeleteModal.value = false
  resourceToDelete.value = ''
}

onMounted(async () => {
  await loadResources()
})
</script>
<template>
  <div>
    <NCard :title="`${props.rType} 列表`" :bordered="false">
      <template #header-extra>
        <NButton type="primary" secondary @click="showCreateModal = true">创建</NButton>
      </template>
      <NList bordered>
        <NListItem v-for="item in resources" :key="item.uid">
          <template #prefix>
            <div style="width: 200px">
              <NTag size="large">{{ item.name }}</NTag>
            </div>
          </template>
          <template #suffix>
            <NButton quaternary type="primary" @click="handleDeleteClick(item.uid)"> 删除 </NButton>
          </template>
          <div class="v-item">ID: {{ item.uid }}</div>
        </NListItem>
      </NList>
    </NCard>
    <NModal
      v-model:show="showCreateModal"
      preset="dialog"
      :title="`创建 ${props.rType}`"
      positive-text="确认"
      negative-text="取消"
      :loading="IsCreating"
      :mask-closable="false"
      @positive-click="handleCreateConfirm"
      @negative-click="handleCreateClose"
    >
      <NForm :model="createForm" :rules="createRules" ref="createFormRef">
        <NFormItem label="名称" path="name">
          <NInput v-model:value="createForm.name" placeholder="请输入名称" />
        </NFormItem>
      </NForm>
    </NModal>
    <NModal
      v-model:show="showDeleteModal"
      preset="dialog"
      :title="`删除 ${props.rType}`"
      positive-text="确认"
      negative-text="取消"
      :loading="IsDeleting"
      :mask-closable="false"
      @positive-click="handleDeleteConfirm"
      @negative-click="handleDeleteClose"
    >
      <div>确认要删除 {{ props.rType }} {{ resourceToDelete }}？</div>
    </NModal>
  </div>
</template>
