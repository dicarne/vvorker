<script setup lang="ts">
import { computed } from 'vue'
import { NCard, NButton, useMessage } from 'naive-ui'
import VVorkerTracker from '@/components/VVorkerTracker.vue'
import type { Node } from '@/types/nodes'
import { deleteNode, syncNodes } from '@/api/nodes'

const props = defineProps<{
  node: Node
  ping: number[]
}>()

const data = computed(() => {
  return props.ping.map((v, i) => {
    console.log()
    if (v >= 1000) {
      return { color: '#bf616a', tooltip: `${v}ms` }
    }
    if (v >= 100) {
      return { color: '#ebcb8b', tooltip: `${v}ms` }
    }
    if (v === -1) {
      return { color: '#4c566a', tooltip: 'N/A' }
    }
    return { color: '#a3be8c', tooltip: `${v}ms` }
  })
})
const sla = computed(() => {
  return parseFloat(
    ((1 - props.ping.filter((v) => v >= 500).length / props.ping.length) * 100).toFixed(2),
  )
})
const avg = computed(() => {
  const validValue = props.ping.filter((v) => v !== -1)
  return parseFloat((validValue.reduce((a, b) => a + b, 0) / validValue.length).toFixed(2))
})

const message = useMessage()

// 同步节点（非 Default）
const handleSyncNodeClick = async (nodeName: string) => {
  try {
    await syncNodes(nodeName)
    message.info('同步成功')
  } catch (error) {
    console.error('syncNodes Error', error)
    message.error('同步失败')
  }
}

// 删除节点（非 Default）
const handleDeleteNodeClick = async (nodeName: string) => {
  try {
    await deleteNode(nodeName)
    message.info('删除成功')
  } catch (error) {
    console.error('deleteNode Error', error)
    message.error('删除失败')
  }
}
</script>

<template>
  <div class="node-item">
    <NCard :title="node.Name">
      <template #header-extra>
        <NButton v-if="node.Name !== 'default'" type="primary" secondary @click="handleSyncNodeClick(node.Name)">同步</NButton>
        <NButton v-if="node.Name !== 'default'" type="primary" secondary @click="handleDeleteNodeClick(node.Name)">删除</NButton>
      </template>
      <div class="v-item-column">{{ node.UID }}</div>
      <div class="v-item-column" style="text-align: right">
        <span class="v-item">Uptime {{ sla }}%</span>Avg. {{ avg }}ms
      </div>
      <VVorkerTracker :data="data" />
    </NCard>
  </div>
</template>

<style scoped>
.node-item {
  width: fit-content;
}
</style>
