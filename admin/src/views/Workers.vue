<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { CH } from '@/lib/color'
import { NCard, NList, NListItem, NButton, useMessage, NAvatar, NTag } from 'naive-ui'
import { DEFAULT_WORKER_ITEM, type WorkerItem } from '@/types/workers'
import { createWorker, getAllWorkers } from '@/api/workers'

const message = useMessage()

const workers = ref<WorkerItem[]>([])
// 加载所有 Worker
const loadWorkers = async () => {
  try {
    workers.value = await getAllWorkers()
  } catch (error) {
    console.error(error)
    message.error('获取 Worker 列表失败: ' + error)
  }
}

// 创建 Worker
const handleCreateWorkerClick = async () => {
  try {
    await createWorker({ ...DEFAULT_WORKER_ITEM })
    await loadWorkers()
    message.success('创建 Worker 成功')
  } catch (error) {
    console.error(error)
    message.error('创建 Worker 失败: ' + error)
  }
}

onMounted(async () => {
  loadWorkers()
})
</script>
<template>
  <div>
    <NCard title="Workers">
      <template #header-extra>
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
            <!-- <NButton type="error" secondary @click="handleDeleteAccessKeyClick(item.key)">删除</NButton> -->
          </template>
          <div class="v-flex-center-start-column">
            <div class="v-item">{{ item.Name }}</div>
            <div class="v-item">
              Node:
              <NTag size="small" :style="{ color: CH.hex(item.NodeName) }">
                {{ item.NodeName }}
              </NTag>
            </div>
          </div>
        </NListItem>
      </NList>
    </NCard>
  </div>
</template>
