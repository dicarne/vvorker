<script setup lang="ts">
import { computed, inject, onMounted, ref, type Ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  NButton,
  useMessage,
  NIcon,
  NBreadcrumb,
  NBreadcrumbItem,
  NTag,
  NTabs,
  NTabPane,
  NLayoutSider,
  NInput,
  NInputNumber,
  NSelect,
  NLayout,
  NLayoutContent,
  NNotificationProvider,
} from 'naive-ui'
import { Copy24Regular as CopyIcon } from '@vicons/fluent'
import WorkerRun from '@/components/WorkerRun.vue'
import WorkerLogs from '@/components/WorkerLogs.vue'
import WorkerRules from '@/components/WorkerRules.vue'
import WorkerAuth from '@/components/WorkerAuth.vue'
import {
  DEFAULT_WORKER_ITEM,
  type VorkerSettingsProperties,
  type WorkerItem,
} from '@/types/workers'
import type { Node } from '@/types/nodes'
import { getWorker, updateWorker } from '@/api/workers'
import { getNodes } from '@/api/nodes'
import { useNavigate } from '@/composables/useNavigate'
import { useCopyContent } from '@/composables/useUtils'
import { genWorkerUrl } from '@/utils/utils'
const router = useRouter()
const uid = router.currentRoute.value.query.uid as string
const { navigate } = useNavigate()
// 定义消息和通知实例
const message = useMessage()
const { copyContent } = useCopyContent()
const appConfig = inject<Ref<VorkerSettingsProperties>>('appConfig')!
const worker = ref<WorkerItem>(DEFAULT_WORKER_ITEM)
const workerURL = computed(() => {
  return genWorkerUrl(appConfig.value, worker.value.Name)
})
// 保存 Worker
const handleSaveWorkerClick = async () => {
  if (!worker.value) {
    message.error('Worker 不存在')
    return
  }
  try {
    await updateWorker(worker.value)
    message.success('保存 Worker 成功')
  } catch (error) {
    console.error('updateWorker Error', error)
    message.error('保存 Worker 失败')
  }
}

const nodes = ref<Node[]>([])
onMounted(async () => {
  try {
    worker.value = await getWorker(uid)
  } catch (error) {
    console.error('getWorker Error', error)
    message.error('获取 Worker 失败')
  }

  try {
    const res = await getNodes()
    if (res.code === 0) {
      nodes.value = res.data.nodes
    }
  } catch (error) {
    console.error('getNodes Error', error)
    message.error('获取节点列表失败')
  }
})

const handleOpenWorkerClick = async () => {
  const workerUrl = genWorkerUrl(appConfig.value, worker.value.Name)
  if (workerUrl) {
    window.open(workerUrl, '_blank')
  } else {
    message.error('Worker URL 无效')
  }
}
</script>
<template>
  <div class="v-main">
    <div class="v-flex-between-center v-item-column">
      <NBreadcrumb>
        <NBreadcrumbItem @click="navigate('/workers')"> Workers </NBreadcrumbItem>
        <NBreadcrumbItem>
          {{ worker?.Name }}
        </NBreadcrumbItem>
      </NBreadcrumb>
      <div>
        <NButton class="v-item" type="primary" secondary @click="navigate('/workers')">
          返回
        </NButton>
        <NButton class="v-item" type="primary" secondary @click="handleSaveWorkerClick"> 保存 </NButton>
        <NButton type="primary" secondary @click="handleOpenWorkerClick"> 打开 </NButton>
      </div>
    </div>
    <div class="v-flex-start-center v-item-column">
      <div>
        UID: <NTag class="v-item">{{ worker.UID }}</NTag>
        <NButton quaternary type="primary" @click="copyContent(worker.UID)">
          <NIcon>
            <CopyIcon />
          </NIcon>
        </NButton>
      </div>
      <div class="v-item">
        URL: <NTag class="v-item">{{ workerURL }}</NTag>
        <NButton quaternary type="primary" @click="copyContent(workerURL)">
          <NIcon>
            <CopyIcon />
          </NIcon>
        </NButton>
      </div>
    </div>
    <NTabs type="line" animated>
      <NTabPane name="property" tab="属性">
        <NLayout has-sider class="v-item-column">
          <NLayoutSider> 名称 </NLayoutSider>
          <NLayoutContent>
            <NInput style="min-width: 200px; max-width: 400px;" v-model:value="worker.Name" />
          </NLayoutContent>
        </NLayout>
        <NLayout has-sider class="v-item-column">
          <NLayoutSider> 节点 </NLayoutSider>
          <NLayoutContent>
            <NSelect style="min-width: 200px; max-width: 400px;" v-model:value="worker.NodeName" :options="nodes.map((node) => ({
              label: node.Name,
              value: node.Name,
            }))
              " />
          </NLayoutContent>
        </NLayout>
        <NLayout has-sider class="v-item-column">
          <NLayoutSider> 实例 </NLayoutSider>
          <NLayoutContent>
            <NInputNumber min="1" max="20" style="min-width: 200px; max-width: 400px;"
              v-model:value="worker.MaxCount" />
          </NLayoutContent>
        </NLayout>
      </NTabPane>
      <NTabPane name="logs" tab="日志">
        <WorkerLogs :uid="worker.UID" />
      </NTabPane>
      <NTabPane name="rules" tab="规则">
        <WorkerRules :uid="worker.UID" />
      </NTabPane>
      <NTabPane name="auth" tab="鉴权">
        <WorkerAuth :uid="worker.UID" />
      </NTabPane>
      <template #suffix>
        <!-- 使用 WorkerRun 组件 -->
        <NNotificationProvider placement="bottom-right">
          <WorkerRun :uid="worker.UID" />
        </NNotificationProvider>
      </template>
    </NTabs>
  </div>
</template>
