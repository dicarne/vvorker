<script setup lang="ts">
import { defineProps, h } from 'vue'
import { NButton, useMessage, useNotification, NIcon } from 'naive-ui'
import { Play24Regular as RunIcon } from '@vicons/fluent'
import { runWorker } from '@/api/workers'
import { decodeBase64 } from '@/utils/utils'

// 定义组件接收的属性
const props = defineProps<{
  uid: string
}>()

// 定义消息和通知实例
const message = useMessage()
const notification = useNotification()

// 运行 Worker 的方法
const handleRun = async () => {
  try {
    const resp = await runWorker(props.uid)
    const decodedResp = decodeBase64(resp?.data?.run_resp)
    notification.info({
      title: 'Worker Run Result',
      content: () => h('div', {}, decodedResp),
    })
  } catch (error) {
    console.error('runWorker Error', error)
    message.error('运行 Worker 失败')
  }
}
</script>
<template>
  <NButton quaternary type="primary" @click="handleRun">
    <NIcon><RunIcon /></NIcon>
    运行
  </NButton>
</template>
