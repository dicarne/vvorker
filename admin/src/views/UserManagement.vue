<script setup lang="ts">
import { inject, onMounted, ref, type Ref, h } from 'vue'
import {
  NCard,
  NDataTable,
  NButton,
  NSelect,
  useMessage,
  NSpace,
  NTag,
  NModal,
  NForm,
  NFormItem,
  NInput,
} from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui'
import { passwordRules } from '@/constant/formrules'
import type { UserInfo } from '@/types/auth'
import { getUsers, updateUserRole, createUser } from '@/api/users'

interface User {
  ID: number
  user_name: string
  email: string
  role: string
  status: number
}

const userInfo = inject<Ref<UserInfo>>('userInfo')!
const message = useMessage()
console.log(userInfo)

const users = ref<User[]>([])
const loading = ref(false)

const loadUsers = async () => {
  loading.value = true
  try {
    const res = await getUsers()
    // 过滤掉自己
    users.value = res.data.data.users.filter((user: User) => user.ID !== userInfo.value.id)
  } catch (error) {
    console.error('getUsers Error', error)
    message.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

const handleRoleChange = async (userId: number, newRole: string) => {
  try {
    await updateUserRole(userId, newRole)
    message.success('角色更新成功')
    // 更新本地数据
    const user = users.value.find(u => u.ID === userId)
    if (user) {
      user.role = newRole
    }
  } catch (error) {
    console.error('updateUserRole Error', error)
    message.error('角色更新失败')
  }
}

// 新建用户相关
const showCreateModal = ref(false)
const createForm = ref({
  username: '',
  password: '',
})
const createFormRef = ref<FormInst | null>(null)
const createRules: FormRules = {
  username: {
    required: true,
    message: '请输入用户名',
  },
  password: passwordRules,
}

const handleCreateUser = async () => {
  if (!createFormRef.value) return
  try {
    await createFormRef.value.validate()
    await createUser(createForm.value)
    message.success('用户创建成功')
    showCreateModal.value = false
    createForm.value = { username: '', password: '' }
    loadUsers()
  } catch (error) {
    console.error('createUser Error', error)
    message.error('用户创建失败')
  }
}

const columns = [
  {
    title: '用户名',
    key: 'user_name',
  },
  {
    title: '角色',
    key: 'role',
    render: (row: User) => {
      return h(NSelect, {
        value: row.role,
        options: [
          { label: '普通用户', value: 'normal' },
          { label: '管理员', value: 'admin' },
        ],
        'onUpdate:value': (value: string) => handleRoleChange(row.ID, value),
        style: { width: '120px' },
      })
    },
  },
  {
    title: '状态',
    key: 'status',
    render: (row: User) => {
      const statusMap = {
        0: { text: '未知', type: 'default' },
        1: { text: '禁用', type: 'error' },
        2: { text: '正常', type: 'success' },
      }
      const status = statusMap[row.status as keyof typeof statusMap] || statusMap[0]
      return h(NTag, { type: status.type as any }, { default: () => status.text })
    },
  },
]

onMounted(() => {
  loadUsers()
})
</script>

<template>
  <div>
    <NCard title="用户管理">
      <NSpace justify="end" style="margin-bottom: 16px;">
        <NButton type="primary" @click="showCreateModal = true">
          新建用户
        </NButton>
      </NSpace>
      <NDataTable
        :columns="columns"
        :data="users"
        :loading="loading"
        :pagination="false"
      />
    </NCard>

    <NModal
      v-model:show="showCreateModal"
      preset="card"
      title="新建用户"
      size="huge"
      :bordered="false"
      :segmented="false"
    >
      <NForm
        ref="createFormRef"
        :model="createForm"
        :rules="createRules"
        label-placement="top"
      >
        <NFormItem label="用户名" path="username">
          <NInput v-model:value="createForm.username" placeholder="请输入用户名" />
        </NFormItem>
        <NFormItem label="密码" path="password">
          <NInput
            v-model:value="createForm.password"
            type="password"
            placeholder="请输入密码"
          />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="showCreateModal = false">取消</NButton>
          <NButton type="primary" @click="handleCreateUser">创建</NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>