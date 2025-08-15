<script setup lang="ts">
import { onMounted, ref } from 'vue'
import {
  NCard,
  NForm,
  NSelect,
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
import { h } from 'vue'
import type {
  AccessRule,
  AccessRuleType,
  UpdateAccessRuleRequest,
  DeleteAccessRuleRequest,
  EnableAccessControlRequest,
  SwitchAccessRuleRequest,
} from '@/types/access'
import {
  addAccessRule,
  updateAccessRule,
  deleteAccessRule,
  getAccessControl,
  listAccessRules,
  switchAccessRule,
  updateEnableAccessControl,
} from '@/api/workers'

const props = defineProps<{
  uid: string
}>()

const message = useMessage()

// 访问控制相关
const isAccessControlEnabled = ref<boolean>(false)
const handleSwitchChange = async (checked: boolean) => {
  try {
    const request: EnableAccessControlRequest = {
      enable: checked,
      worker_uid: props.uid,
    }
    await updateEnableAccessControl(request)
    isAccessControlEnabled.value = checked
    message.success('更新访问控制状态成功')
  } catch (error) {
    console.error('updateEnableAccessControl Error', error)
    message.error('更新访问控制状态失败')
  }
}
const renderAccessControl = () =>
  h('div', { class: 'v-flex-start-center' }, [
    h('span', '启用访问控制'),
    h(NSwitch, {
      class: 'v-item',
      // 将 value 作为属性传入
      value: isAccessControlEnabled.value,
      onUpdateValue: async (value: boolean) => {
        await handleSwitchChange(value)
      },
    }),
  ])

const rules = ref<AccessRule[]>([])
const fetchRules = async () => {
  try {
    const response1 = await getAccessControl({ worker_uid: props.uid })
    isAccessControlEnabled.value = response1.data.EnableAccessControl
    const response = await listAccessRules({
      worker_uid: props.uid,
      page: 1,
      page_size: 100,
    })
    rules.value = response.data.access_rules
  } catch (error) {
    console.error('listAccessRules Error', error)
    message.error('获取访问规则失败')
  }
}

// 添加 rule
const ruleTypeOptions = [
  {
    label: '内部规则',
    value: 'internal',
  },
  {
    label: '外部TOKEN',
    value: 'token',
  },
  {
    label: '开放',
    value: 'open',
  },
  {
    label: 'SSO',
    value: 'sso',
  },
]
const showCreateRuleModal = ref<boolean>(false)
const IsCreatingRule = ref<boolean>(false)
const createRuleForm = ref({
  path: '/',
  description: '',
  ruleType: 'internal',
  data: '',
})
const createRuleFormRef = ref<FormInst | null>(null)
const workerRuleRules: FormRules = {
  path: {
    required: true,
    message: '请输入路径',
  },
  description: {
    required: true,
    message: '请输入描述',
  },
  rule_type: {
    required: true,
    message: '请选择规则类型',
  },
}
const handleCreateRuleConfirm = async () => {
  if (!createRuleFormRef.value) return
  try {
    // 校验表单
    await createRuleFormRef.value.validate()
  } catch (error) {
    console.error('createRuleFormRef validate Error', error)
    return
  }
  try {
    // 调用创建规则接口
    IsCreatingRule.value = true
    const request: AccessRule = {
      worker_uid: props.uid,
      rule_type: createRuleForm.value.ruleType as AccessRuleType,
      path: createRuleForm.value.path,
      description: createRuleForm.value.description,
      rule_uid: '',
      data: createRuleForm.value.data,
    }
    await addAccessRule(request)
    await fetchRules()
    message.success('创建规则成功')
    handleCreateRuleClose()
  } catch (error) {
    console.error('addAccessRule Error', error)
    message.error('创建规则失败')
  } finally {
    IsCreatingRule.value = false
  }
}
const handleCreateRuleClose = () => {
  showCreateRuleModal.value = false
  createRuleForm.value.path = '/'
  createRuleForm.value.description = ''
  createRuleForm.value.ruleType = 'internal'
  createRuleForm.value.data = ''
}

// 启用/禁用每条规则
const handleRuleSwitchChange = async (item: AccessRule) => {
  try {
    const request: SwitchAccessRuleRequest = {
      worker_uid: props.uid,
      rule_uid: item.rule_uid,
      disable: !!item.disabled,
    }
    await switchAccessRule(request)
    await fetchRules()
    message.success('更新规则状态成功')
  } catch (error) {
    console.error('switchAccessRule Error', error)
    message.error('更新规则状态失败')
  }
}

// 编辑 rule
const showEditRuleModal = ref<boolean>(false)
const IsEditingRule = ref<boolean>(false)
const editRuleForm = ref({
  path: '/',
  description: '',
  ruleType: 'internal',
  data: '',
})
const editRuleFormRef = ref<FormInst | null>(null)
const ruleUidToEdit = ref<string>('')
const handleEditRuleClick = async (item: AccessRule) => {
  ruleUidToEdit.value = item.rule_uid
  editRuleForm.value.path = item.path
  editRuleForm.value.description = item.description
  editRuleForm.value.ruleType = item.rule_type
  editRuleForm.value.data = item.data

  showEditRuleModal.value = true
}
const handleEditRuleConfirm = async () => {
  if (!ruleUidToEdit.value) return
  try {
    IsEditingRule.value = true
    const request: UpdateAccessRuleRequest = {
      worker_uid: props.uid,
      rule_uid: ruleUidToEdit.value,
      description: editRuleForm.value.description,
      path: editRuleForm.value.path,
      rule_type: editRuleForm.value.ruleType as AccessRuleType,
      data: editRuleForm.value.data,
    }
    await updateAccessRule(request)
    await fetchRules()
    message.success('编辑规则成功')
    handleEditRuleClose()
  } catch (error) {
    console.error('updateAccessRule Error', error)
    message.error('编辑规则失败')
  } finally {
    IsEditingRule.value = false
  }
}
const handleEditRuleClose = () => {
  showEditRuleModal.value = false
  ruleUidToEdit.value = ''
}

// 删除 rule
const showDeleteRuleModal = ref<boolean>(false)
const IsDeletingRule = ref<boolean>(false)
const ruleUidToDelete = ref<string>('')
const handleDeleteRuleClick = async (uid: string) => {
  ruleUidToDelete.value = uid
  showDeleteRuleModal.value = true
}
const handleDeleteRuleConfirm = async () => {
  if (!ruleUidToDelete.value) return
  try {
    IsDeletingRule.value = true
    const request: DeleteAccessRuleRequest = {
      worker_uid: props.uid,
      rule_uid: ruleUidToDelete.value,
    }
    await deleteAccessRule(request)
    await fetchRules()
    message.success('删除规则成功')
    handleDeleteRuleClose()
  } catch (error) {
    console.error('deleteAccessRule Error', error)
    message.error('删除规则失败')
  } finally {
    IsDeletingRule.value = false
  }
}
const handleDeleteRuleClose = () => {
  showDeleteRuleModal.value = false
  ruleUidToDelete.value = ''
}

onMounted(async () => {
  await fetchRules()
})
</script>

<template>
  <div>
    <NCard :title="renderAccessControl" :bordered="false">
      <template #header-extra>
        <NButton type="primary" secondary @click="showCreateRuleModal = true">添加规则</NButton>
      </template>
      <NTable :bordered="false" :single-line="false">
        <thead>
          <tr>
            <th>路由前缀</th>
            <th>控制类型</th>
            <th>描述</th>
            <th>权限</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in rules" :key="item.rule_uid">
            <td>{{ item.path }}</td>
            <td>
              {{ ruleTypeOptions.find((option) => option.value === item.rule_type)?.label }}
              ({{ item.rule_type }})
            </td>
            <td>{{ item.description }}</td>
            <td>{{ item.data }}</td>
            <td>
              <NSwitch class="v-item" :round="false" :value="!item.disabled" @update:value="handleRuleSwitchChange(item)"/>
              <NButton quaternary type="primary" @click="handleEditRuleClick(item)"> 编辑 </NButton>
              <NButton quaternary type="primary" @click="handleDeleteRuleClick(item.rule_uid)">
                删除
              </NButton>
            </td>
          </tr>
        </tbody>
      </NTable>
    </NCard>
    <NModal
      v-model:show="showCreateRuleModal"
      preset="dialog"
      title="添加规则"
      positive-text="确认"
      negative-text="取消"
      :loading="IsCreatingRule"
      :mask-closable="false"
      @positive-click="handleCreateRuleConfirm"
      @negative-click="handleCreateRuleClose"
    >
      <NForm ref="createRuleFormRef" :rules="workerRuleRules" :model="createRuleForm">
        <NFormItem label="路由前缀">
          <NInput v-model:value="createRuleForm.path" placeholder="请输入路由前缀" />
        </NFormItem>
        <NFormItem label="描述">
          <NInput v-model:value="createRuleForm.description" placeholder="请输入描述" />
        </NFormItem>
        <NFormItem label="规则类型">
          <NSelect v-model:value="createRuleForm.ruleType" :options="ruleTypeOptions" />
        </NFormItem>
        <NFormItem label="权限">
          <NInput v-model:value="createRuleForm.data" placeholder="请输入权限" />
        </NFormItem>
      </NForm>
    </NModal>
    <NModal
      v-model:show="showEditRuleModal"
      preset="dialog"
      title="编辑规则"
      positive-text="确认"
      negative-text="取消"
      :loading="IsEditingRule"
      :mask-closable="false"
      @positive-click="handleEditRuleConfirm"
      @negative-click="handleEditRuleClose"
    >
      <NForm ref="editRuleFormRef" :rules="workerRuleRules" :model="editRuleForm">
        <NFormItem label="路由前缀">
          <NInput v-model:value="editRuleForm.path" placeholder="请输入路由前缀" />
        </NFormItem>
        <NFormItem label="描述">
          <NInput v-model:value="editRuleForm.description" placeholder="请输入描述" />
        </NFormItem>
        <NFormItem label="规则类型">
          <NSelect v-model:value="editRuleForm.ruleType" :options="ruleTypeOptions" />
        </NFormItem>
        <NFormItem label="权限">
          <NInput v-model:value="editRuleForm.data" placeholder="请输入权限" />
        </NFormItem>
      </NForm>
    </NModal>
    <NModal
      v-model:show="showDeleteRuleModal"
      preset="dialog"
      title="删除规则"
      positive-text="确认"
      negative-text="取消"
      :loading="IsDeletingRule"
      :mask-closable="false"
      @positive-click="handleDeleteRuleConfirm"
      @negative-click="handleDeleteRuleClose"
    >
      <div>确认要删除此规则？</div>
    </NModal>
  </div>
</template>

<style scoped></style>
