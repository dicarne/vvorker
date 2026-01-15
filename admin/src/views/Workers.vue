<script setup lang="ts">
import { inject, onMounted, ref, type Ref } from 'vue'
import { CH } from '@/lib/color'
import {
  NCard,
  NList,
  NListItem,
  NButton,
  useMessage,
  NAvatar,
  NTag,
  NDropdown,
  NModal,
  NIcon,
  NNotificationProvider,
  NText,
  NEllipsis
} from 'naive-ui'
import {
  MoreHorizontal24Regular as DropdownIcon,
  Edit24Regular as EditIcon,
  Link24Regular as LinkIcon,
  LockClosed24Regular as LockIcon,
} from '@vicons/fluent'
// 引入 WorkerRun 组件
import WorkerRun from '@/components/WorkerRun.vue'
import {
  DEFAULT_WORKER_ITEM,
  type VorkerSettingsProperties,
  type WorkerItem,
} from '@/types/workers'
import {
  createWorker,
  deleteWorker,
  flushAllWorkers,
  flushWorker,
  getAllWorkers,
} from '@/api/workers'
import { useNavigate } from '@/composables/useNavigate'
import { genWorkerUrl } from '@/utils/utils'

const message = useMessage()
const { navigate } = useNavigate()
const appConfig = inject<Ref<VorkerSettingsProperties>>('appConfig')!
const workers = ref<WorkerItem[]>([])

// 加载所有 Worker
const loadWorkers = async () => {
  try {
    workers.value = await getAllWorkers()
  } catch (error) {
    console.error('getAllWorkers Error', error)
    message.error('获取 Worker 列表失败')
  }
}

const handleReloadWorkersClick = async () => {
  await loadWorkers()
  message.success('刷新 Worker 列表成功')
}

// 同步所有 Worker
const handleFlushAllWorkersClick = async () => {
  try {
    await flushAllWorkers()
    message.success('同步 Workers 成功')
  } catch (error) {
    console.error('flushAllWorkers Error', error)
    message.error('同步 Workers 失败')
  }
}

// 创建 Worker
const handleCreateWorkerClick = async () => {
  try {
    await createWorker({ ...DEFAULT_WORKER_ITEM })
    await loadWorkers()
    message.success('创建 Worker 成功')
  } catch (error) {
    console.error('createWorker Error', error)
    message.error('创建 Worker 失败')
  }
}

// 打开 Worker
const handleOpenWorkerClick = async (worker: WorkerItem) => {
  const workerUrl = genWorkerUrl(appConfig.value, worker.Name)
  if (workerUrl) {
    window.open(workerUrl, '_blank')
  } else {
    message.error('Worker URL 无效')
  }
}

// 同步 Worker
const handleFlushWorkerClick = async (uid: string) => {
  try {
    await flushWorker(uid)
    message.success('同步 Worker 成功')
  } catch (error) {
    console.error('flushWorker Error', error)
    message.error('同步 Worker 失败')
  }
}
// 删除 Worker
const showDeleteWorkerModal = ref<boolean>(false)
const IsDeletingWorker = ref<boolean>(false)
const workerToDelete = ref<WorkerItem>()
const handleDeleteWorkerClick = async (worker: WorkerItem) => {
  workerToDelete.value = worker
  showDeleteWorkerModal.value = true
}
const handleDeleteWorkerConfirm = async () => {
  if (!workerToDelete.value) return
  try {
    IsDeletingWorker.value = true
    await deleteWorker(workerToDelete.value.UID)
    await loadWorkers()
    message.success('删除 Worker 成功')
    handleDeleteWorkerClose()
  } catch (error) {
    console.error('deleteWorker Error', error)
    message.error('删除 Worker 失败')
  } finally {
    IsDeletingWorker.value = false
  }
}
const handleDeleteWorkerClose = () => {
  showDeleteWorkerModal.value = false
  workerToDelete.value = undefined
}

// Worker 下拉菜单
const dropdownOptions = [
  {
    label: '同步',
    key: 'sync',
  },
  {
    label: '删除',
    key: 'delete',
  },
]

const handleDropdownSelect = (worker: WorkerItem, key: string) => {
  if (key === 'sync') {
    handleFlushWorkerClick(worker.UID)
    return
  }
  if (key == 'delete') {
    handleDeleteWorkerClick(worker)
    return
  }
}

onMounted(async () => {
  loadWorkers()
})
</script>
<template>
  <div class="v-main">
    <NCard title="Workers" :bordered="false">
      <template #header-extra>
        <NButton type="primary" secondary @click="handleReloadWorkersClick">刷新</NButton>
        <NButton class="v-item" type="primary" secondary @click="handleFlushAllWorkersClick">
          同步
        </NButton>
        <NButton type="primary" secondary @click="handleCreateWorkerClick">创建</NButton>
      </template>
      <NList>
        <NListItem v-for="item in workers" :key="item.UID">
          <template #prefix>
            <NAvatar class="v-avatar" :style="{ background: CH.hex(item.Name) }">
              {{ item.Name.slice(0, 2).toUpperCase() }}
            </NAvatar>
          </template>
          <template #suffix>
            <div class="v-flex-center">
              <NButton quaternary type="primary" @click="navigate(`/workeredit?uid=${item.UID}`)">
                <NIcon>
                  <EditIcon />
                </NIcon>编辑
              </NButton>
              <!-- 使用 WorkerRun 组件 -->
              <NNotificationProvider placement="bottom-right">
                <WorkerRun :uid="item.UID" />
              </NNotificationProvider>
              <NButton quaternary type="primary" @click="handleOpenWorkerClick(item)">
                <NIcon>
                  <LinkIcon />
                </NIcon>打开
              </NButton>
              <NDropdown trigger="hover" :options="dropdownOptions" @select="(key) => handleDropdownSelect(item, key)">
                <NButton quaternary>
                  <NIcon>
                    <DropdownIcon />
                  </NIcon>
                </NButton>
              </NDropdown>
            </div>
          </template>
          <div class="v-flex-center-start-column">
            <div class="v-item">
              <NText class="title-text">{{ item.Name }}</NText>
              <NIcon class="v-item" v-if="item.AccessControl">
                <LockIcon />
              </NIcon>
              <NTag size="small" :style="{ color: CH.hex(item.NodeName || '') }">
                {{ item.NodeName }}
              </NTag>
            </div>
            <div class="v-item desc-text">
              <NEllipsis style="max-width: 500px;" expand-trigger="click" line-clamp="1" :tooltip="false">{{ item.Description }}</NEllipsis>
            </div>
          </div>
        </NListItem>
      </NList>
    </NCard>
    <NModal v-model:show="showDeleteWorkerModal" preset="dialog" title="删除 Worker" positive-text="确认" negative-text="取消"
      :loading="IsDeletingWorker" :mask-closable="false" @positive-click="handleDeleteWorkerConfirm"
      @negative-click="handleDeleteWorkerClose">
      <div>确认要删除 {{ workerToDelete?.Name }}（ID: {{ workerToDelete?.UID }}）？</div>
    </NModal>
  </div>
</template>
<style lang="css" scoped>
.title-text {
  font-weight: 500;
  margin-right: 10px;
}

.desc-text {
  color: rgb(131, 131, 131);
}
</style>