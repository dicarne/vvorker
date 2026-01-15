<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  NCard,
  NButton,
  NInput,
  NList,
  NListItem,
  NPopconfirm,
  NText,
  NIcon,
  useMessage,
  NEmpty,
} from 'naive-ui'
import { Delete24Regular as DeleteIcon } from '@vicons/fluent'
import type { WorkerMember } from '@/types/workers'
import { addWorkerMember, removeWorkerMember, listWorkerMembers } from '@/api/workers'

const props = defineProps<{ uid: string; canManage: boolean }>()
const message = useMessage()
const members = ref<WorkerMember[]>([])
const newMemberName = ref('')
const isLoading = ref(false)

const loadMembers = async () => {
  try {
    isLoading.value = true
    members.value = await listWorkerMembers(props.uid)
  } catch (error) {
    console.error('loadMembers error', error)
    message.error('加载成员列表失败')
  } finally {
    isLoading.value = false
  }
}

const handleAddMember = async () => {
  if (!newMemberName.value.trim()) {
    message.warning('请输入用户名')
    return
  }
  try {
    await addWorkerMember(props.uid, newMemberName.value)
    message.success('添加成员成功')
    newMemberName.value = ''
    await loadMembers()
  } catch (error: any) {
    console.error('addMember error', error)
    message.error(error.response?.data?.msg || '添加成员失败')
  }
}

const handleRemoveMember = async (member: WorkerMember) => {
  try {
    await removeWorkerMember(props.uid, member.user_name)
    message.success('移除成员成功')
    await loadMembers()
  } catch (error: any) {
    console.error('removeMember error', error)
    message.error(error.response?.data?.msg || '移除成员失败')
  }
}

onMounted(() => {
  loadMembers()
})
</script>

<template>
  <NCard title="协作者管理">
    <div v-if="canManage" class="v-flex-center-start v-item-column" style="margin-bottom: 16px;">
      <NInput v-model:value="newMemberName" placeholder="输入用户名添加成员" style="width: 300px; margin-right: 8px;" />
      <NButton type="primary" @click="handleAddMember">添加</NButton>
    </div>

    <NList v-if="members.length > 0">
      <NListItem v-for="member in members" :key="member.ID">
        <div class="v-flex-between-center">
          <div>
            <NText style="margin-right: 8px; font-weight: 500;">{{ member.user_name }}</NText>
            <NText type="info" style="font-size: 12px;">
              添加于: {{ new Date(member.joined_at).toLocaleString() }}
            </NText>
          </div>
          <NPopconfirm v-if="canManage" @positive-click="() => handleRemoveMember(member)" positive-text="删除"
            negative-text="取消">
            <template #trigger>
              <NButton quaternary type="error">
                <NIcon>
                  <DeleteIcon />
                </NIcon>
              </NButton>
            </template>
            确定要移除 {{ member.user_name }} 吗？
          </NPopconfirm>
        </div>
      </NListItem>
    </NList>
    <NEmpty v-else-if="!isLoading" description="暂无协作者" />
  </NCard>
</template>
