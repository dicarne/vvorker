<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useMessage } from 'naive-ui'
import NodeItem from '@/components/NodeItem.vue'
import type { Node, PingMapList } from '@/types/nodes'
import { getNodes } from '@/api/nodes'

const message = useMessage()
const nodes = ref<Node[]>([])
const pingMapList = ref<PingMapList>({})

// 加载所有 node
const loadNodes = async () => {
  try {
    const resp = await getNodes()
    nodes.value = resp.data.nodes
    const v = Object.entries(resp.data.ping).map(([k, v]) => {
      let t = pingMapList.value[k] || Array.from({ length: 50 }, () => -1)
      if (t.length > 50) {
        t.shift()
      }
      let s = [k, [...t, resp?.data.ping[k] || 0]]
      return s
    })
    const a = Object.fromEntries(v) as PingMapList
    pingMapList.value = a
  } catch (error) {
    console.error('getNodes Error', error)
    message.error('获取 Node 列表失败')
  }
}

let intervalId: number | null = null
onMounted(async () => {
  await loadNodes()
})
onMounted(async () => {
  await loadNodes()
  // 启动定时任务，每 10 秒获取一次
  intervalId = window.setInterval(async () => {
    await loadNodes()
  }, 10000)
})
onUnmounted(() => {
  // 组件销毁前清除定时任务
  if (intervalId) {
    window.clearInterval(intervalId)
  }
})
</script>
<template>
  <div class="v-main v-flex-start-center" style="flex-wrap: wrap; gap: 16px">
    <NodeItem v-for="node in nodes" :key="node.UID" :node="node" :ping="pingMapList[node.Name]" />
  </div>
</template>
<style scoped></style>
