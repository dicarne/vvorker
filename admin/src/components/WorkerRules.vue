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
  DeleteAccessRuleRequest,
  EnableAccessControlRequest,
} from '@/types/access'
import {
  addAccessRule,
  deleteAccessRule,
  listAccessRules,
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
    console.error('Failed to update access control status', error)
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
    const response = await listAccessRules({
      worker_uid: props.uid,
      page: 1,
      page_size: 100,
    })
    rules.value = response.data.access_rules
  } catch (error) {
    console.error('Failed to fetch access rules', error)
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
})
const createRuleFormRef = ref<FormInst | null>(null)
const createRuleRules: FormRules = {
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
    console.error(error)
    return
  }
  try {
    // 调用创建规则接口
    IsCreatingRule.value = true
    await addAccessRule({
      worker_uid: props.uid,
      rule_type: createRuleForm.value.ruleType as AccessRuleType,
      path: createRuleForm.value.path,
      description: createRuleForm.value.description,
      rule_uid: '',
    })
    await fetchRules()
    message.success('创建规则成功')
    handleCreateRuleClose()
  } catch (error) {
    console.error(error)
    message.error('创建规则失败: ' + error)
  } finally {
    IsCreatingRule.value = false
  }
}
const handleCreateRuleClose = () => {
  showCreateRuleModal.value = false
  createRuleForm.value.path = '/'
  createRuleForm.value.description = ''
  createRuleForm.value.ruleType = 'internal'
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
    console.error(error)
    message.error('删除规则失败: ' + error)
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
            <td>
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
      <NForm ref="createRuleFormRef" :rules="createRuleRules" :model="createRuleForm">
        <NFormItem label="路由前缀">
          <NInput v-model:value="createRuleForm.path" placeholder="请输入路由前缀" />
        </NFormItem>
        <NFormItem label="描述">
          <NInput v-model:value="createRuleForm.description" placeholder="请输入描述" />
        </NFormItem>
        <NFormItem label="规则类型">
          <NSelect v-model:value="createRuleForm.ruleType" :options="ruleTypeOptions" />
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
