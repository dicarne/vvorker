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
onMounted(() => {
  // 立即获取一次日志
  fetchLogs()
  // 启动定时任务，每秒获取一次日志
  intervalId = window.setInterval(() => {
    fetchLogs()
  }, 1000)
})
onUnmounted(() => {
  // 组件销毁前清除定时任务
  if (intervalId) {
    window.clearInterval(intervalId)
  }
})
</script>
<template>
  <div ref="containerRef">
    <NList class="log-container" bordered>
      <NListItem v-for="item in logs" :key="item.log_uid">
        <template #prefix>
          <NTag size="small">{{ formatDate(new Date(item.time)) }}</NTag>
        </template>
        <span>{{ item.output }}</span>
      </NListItem>
    </NList>
    <NPagination
      class="v-item-column"
      v-model:page="curPage"
      :page-size="DEFAULT_PAGE_SIZE"
      :item-count="totalCount"
      @update:page="fetchLogs"
    />
  </div>
</template>

<style scoped>
.log-container {
  /* 运算符前后添加空格 */
  height: calc(100vh - 300px);
  overflow-y: auto;
}
</style>
