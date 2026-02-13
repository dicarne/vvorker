<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { NPagination, NTag, NList, NListItem, NAffix } from 'naive-ui'
import type { WorkerLog } from '@/types/workers'
import { formatDate } from '@/utils/utils'
import { getWorkerLogs } from '@/api/workers'

const props = defineProps<{
  uid: string
}>()

const containerRef = ref<HTMLElement | undefined>(undefined)

const DEFAULT_PAGE_SIZE = 50
const totalCount = ref<number>(100)
const curPage = ref<number>(1)

const logs = ref<WorkerLog[]>([])

const fetchLogs = async () => {
  const res = await getWorkerLogs(props.uid, curPage.value, DEFAULT_PAGE_SIZE)
  logs.value = res.data.data.logs
  totalCount.value = res.data.data.total
}

let intervalId: number | null = null
onMounted(async () => {
  // 立即获取一次日志
  await fetchLogs()
  // 启动定时任务，每秒获取一次日志
  intervalId = window.setInterval(async () => {
    await fetchLogs()
  }, 1000)
})
onUnmounted(() => {
  // 组件销毁前清除定时任务
  if (intervalId) {
    window.clearInterval(intervalId)
  }
})

function renderTag(type: string) {
  if (type === 'warn') return "warning"
  if (type === 'success') return "success"
  if (type === 'error') return "error"
  if (type === "stderr") return "error"
  return "default"
}
</script>
<template>
  <div ref="containerRef">
    <NList class="log-container" bordered>
      <NListItem v-for="item in logs" :key="item.log_uid ? item.log_uid : item.uid + item.time">
        <template #prefix>
          <NTag size="small" :type="renderTag(item.type)">{{ formatDate(new Date(item.time)) }}
          </NTag>
        </template>
        <div class="log-content">{{ item.output }}</div>
      </NListItem>
    </NList>
    <NPagination class="v-item-column" v-model:page="curPage" :page-size="DEFAULT_PAGE_SIZE" :item-count="totalCount"
      @update:page="fetchLogs" />
  </div>
</template>

<style scoped>
.log-container {
  /* 运算符前后添加空格 */
  height: calc(100vh - 310px);
  overflow-y: auto;
}

.log-content {
  word-wrap: break-word;
  word-break: break-word;
  white-space: pre-wrap;
  overflow-wrap: break-word;
}
</style>
