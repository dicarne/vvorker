<script setup lang="ts">
import { h, onMounted, onUnmounted, ref, watch } from 'vue'
import {
  NPagination,
  NTag,
  NButton,
  NDataTable,
  NEmpty,
} from 'naive-ui'
import type { Task, TaskLog } from '@/types/workers'
import { formatDate } from '@/utils/utils'
import { listTasks, getTaskLogs, interruptTask } from '@/api/workers'
import { useMessage } from 'naive-ui'

const props = defineProps<{
  uid: string
}>()

const message = useMessage()
const DEFAULT_PAGE_SIZE = 12
const totalCount = ref<number>(0)
const curPage = ref<number>(1)

const tasks = ref<Task[]>([])
const showLogs = ref<boolean>(false)
const currentTaskId = ref<string>('')
const logs = ref<TaskLog[]>([])
const logTotalCount = ref<number>(0)
const logCurPage = ref<number>(1)
const LOG_PAGE_SIZE = 50

// 表格最大高度
const tableMaxHeight = ref<number>(500)

// 获取任务列表
const fetchTasks = async () => {
  try {
    const res = await listTasks(curPage.value, DEFAULT_PAGE_SIZE, props.uid)
    tasks.value = res.data.data.tasks || []
    totalCount.value = res.data.data.total || 0
  } catch (error) {
    console.error('fetchTasks error', error)
  }
}

// 获取任务日志
const fetchLogs = async () => {
  if (!currentTaskId.value) return
  try {
    const res = await getTaskLogs(props.uid, currentTaskId.value, logCurPage.value, LOG_PAGE_SIZE)
    logs.value = res.data.data.logs || []
    logTotalCount.value = res.data.data.total || 0
  } catch (error) {
    console.error('fetchLogs error', error)
  }
}

// 取消任务
const handleCancelTask = async (traceId: string) => {
  try {
    await interruptTask(traceId, props.uid)
    message.success('任务已取消')
    await fetchTasks()
  } catch (error) {
    console.error('cancelTask error', error)
    message.error('取消任务失败')
  }
}

// 查看日志
const handleViewLogs = (traceId: string) => {
  currentTaskId.value = traceId
  showLogs.value = true
  logCurPage.value = 1
}

// 返回任务列表
const handleBackToList = () => {
  showLogs.value = false
  currentTaskId.value = ''
  logs.value = []
}

// 状态标签类型
const getStatusType = (status: string) => {
  switch (status) {
    case 'running':
      return 'info'
    case 'completed':
      return 'success'
    case 'canceled':
      return 'warning'
    case 'failed':
      return 'error'
    case 'interrupt':
      return 'error'
    default:
      return 'default'
  }
}

// 状态显示文本
const getStatusText = (status: string) => {
  switch (status) {
    case 'running':
      return '运行中'
    case 'completed':
      return '已完成'
    case 'canceled':
      return '已取消'
    case 'failed':
      return '失败'
    case 'interrupt':
      return '已中断'
    default:
      return status
  }
}

// 表格列定义
const columns = [
  {
    title: '任务ID',
    key: 'trace_id',
    width: 280,
    ellipsis: { tooltip: true },
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render(row: Task) {
      return h(
        NTag,
        { type: getStatusType(row.status), size: 'small' },
        { default: () => getStatusText(row.status) }
      )
    },
  },
  {
    title: '开始时间',
    key: 'start_time',
    width: 180,
    render(row: Task) {
      return row.start_time ? formatDate(new Date(row.start_time)) : '-'
    },
  },
  {
    title: '结束时间',
    key: 'end_time',
    width: 180,
    render(row: Task) {
      return row.end_time ? formatDate(new Date(row.end_time)) : '-'
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 160,
    render(row: Task) {
      const buttons = [
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            text: true,
            style: { marginRight: '12px' },
            onClick: () => handleViewLogs(row.trace_id),
          },
          { default: () => '查看日志' }
        ),
      ]
      if (row.status === 'running') {
        buttons.push(
          h(
            NButton,
            {
              size: 'small',
              type: 'error',
              text: true,
              onClick: () => handleCancelTask(row.trace_id),
            },
            { default: () => '终止' }
          )
        )
      }
      return buttons
    },
  },
]

let intervalId: number | null = null
let logIntervalId: number | null = null

// 计算表格高度
const updateTableHeight = () => {
  tableMaxHeight.value = window.innerHeight - 350
}

onMounted(async () => {
  updateTableHeight()
  window.addEventListener('resize', updateTableHeight)
  await fetchTasks()
  intervalId = window.setInterval(fetchTasks, 1000)
})

onUnmounted(() => {
  window.removeEventListener('resize', updateTableHeight)
  if (intervalId) {
    window.clearInterval(intervalId)
  }
  if (logIntervalId) {
    window.clearInterval(logIntervalId)
  }
})

// 监听 showLogs 变化，启动/停止日志轮询
watch(showLogs, (newVal) => {
  if (newVal) {
    // 进入日志视图，启动日志轮询
    fetchLogs()
    logIntervalId = window.setInterval(fetchLogs, 1000)
  } else {
    // 离开日志视图，停止日志轮询
    if (logIntervalId) {
      window.clearInterval(logIntervalId)
      logIntervalId = null
    }
  }
})
</script>

<template>
  <div class="worker-tasks">
    <!-- 日志视图 -->
    <div v-if="showLogs" class="logs-container">
      <div class="logs-header">
        <NButton text type="primary" @click="handleBackToList">
          ← 返回任务列表
        </NButton>
        <span class="task-id">任务ID: {{ currentTaskId }}</span>
      </div>
      <div class="log-list">
        <NEmpty v-if="logs.length === 0" description="暂无日志" />
        <div v-else class="log-items">
          <div v-for="(log, index) in logs" :key="index" class="log-item">
            <NTag size="small" type="default">{{ formatDate(new Date(log.time)) }}</NTag>
            <span class="log-content">{{ log.content }}</span>
          </div>
        </div>
      </div>
      <NPagination
        v-if="logTotalCount > LOG_PAGE_SIZE"
        v-model:page="logCurPage"
        :page-size="LOG_PAGE_SIZE"
        :item-count="logTotalCount"
        class="v-item-column"
        @update:page="fetchLogs"
      />
    </div>

    <!-- 任务列表视图 -->
    <div v-else class="tasks-container">
      <NDataTable
        :columns="columns"
        :data="tasks"
        :bordered="false"
        :single-line="false"
        :max-height="tableMaxHeight"
      />
      <NPagination
        v-model:page="curPage"
        :page-size="DEFAULT_PAGE_SIZE"
        :item-count="totalCount"
        class="v-item-column"
        @update:page="fetchTasks"
      />
    </div>
  </div>
</template>

<style scoped>
.worker-tasks {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.logs-container,
.tasks-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.tasks-container :deep(.n-data-table) {
  flex: 1;
  overflow: hidden;
}

.tasks-container :deep(.n-data-table-wrapper) {
  height: 100%;
}

.logs-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.task-id {
  color: #666;
  font-size: 14px;
}

.log-list {
  flex: 1;
  overflow-y: auto;
  background: #f9f9f9;
  border-radius: 4px;
  padding: 12px;
  margin-bottom: 12px;
}

.log-items {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.log-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 4px 0;
}

.log-content {
  word-break: break-all;
  white-space: pre-wrap;
  font-family: monospace;
  font-size: 13px;
}
</style>
