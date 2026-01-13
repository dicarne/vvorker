<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterView } from 'vue-router'
import type { MenuOption } from 'naive-ui'
import { NLayout, NLayoutSider, NMenu } from 'naive-ui'
import {
  CalendarWorkWeek24Regular as WorkersIcon,
  Status24Regular as StatusIcon,
  TaskListSquareLtr24Regular as TaskIcon,
  DatabaseLink24Regular as PGSQLIcon,
  Database24Regular as MYSQLIcon,
  Cloud24Regular as OSSIcon,
  Braces24Filled as KVIcon,
  Person24Regular as UserIcon,
} from '@vicons/fluent'
import TheHeader from '@/components/TheHeader.vue'
import { renderIcon, renderMenuRouterLink } from '@/utils/render'
import { getFeatures } from '@/api/features'
import type { Feature } from '@/types/features'

const collapsed = ref<boolean>(true)
const activeKey = ref<string>('')
const features = ref<Feature[]>([])

type MenuOptionWithFeature = MenuOption & {
  feature?: string
}

const allMenuOptions: MenuOptionWithFeature[] = [
  {
    label: renderMenuRouterLink('Workers', 'Workers'),
    key: 'workers',
    icon: renderIcon(WorkersIcon),
  },
  {
    label: renderMenuRouterLink('状态', 'Status'),
    key: 'status',
    icon: renderIcon(StatusIcon),
  },
  // {
  //   label: renderMenuRouterLink('Task', 'Task'),
  //   key: 'task',
  //   icon: renderIcon(TaskIcon),
  // },
  {
    label: renderMenuRouterLink('PGSQL', 'PGSQL'),
    key: 'pgsql',
    feature: 'pgsql',
    icon: renderIcon(PGSQLIcon),
  },
  {
    label: renderMenuRouterLink('MySQL', 'MySQL'),
    key: 'mysql',
    feature: 'mysql',
    icon: renderIcon(MYSQLIcon),
  },
  {
    label: renderMenuRouterLink('OSS', 'OSS'),
    key: 'oss',
    feature: 'minio',
    icon: renderIcon(OSSIcon),
  },
  {
    label: renderMenuRouterLink('KV', 'KV'),
    key: 'kv',
    feature: 'redis',
    icon: renderIcon(KVIcon),
  },
  {
    label: renderMenuRouterLink('用户', 'User'),
    key: 'user',
    icon: renderIcon(UserIcon),
  },
]

const menuOptions = computed<MenuOptionWithFeature[]>(() => {
  const featureMap = new Map(features.value.map(f => [f.name, f.enable]))

  return allMenuOptions.filter(option => {
    if (!option.feature) {
      return true
    }
    return featureMap.get(option.feature) === true
  })
})

onMounted(async () => {
  try {
    features.value = await getFeatures()
  } catch (error) {
    console.error('Failed to load features:', error)
  }
})
</script>

<template>
  <div>
    <div><TheHeader /></div>
    <NLayout has-sider>
      <NLayoutSider
        bordered
        collapse-mode="width"
        :collapsed-width="64"
        :width="240"
        :collapsed="collapsed"
        show-trigger
        @collapse="collapsed = true"
        @expand="collapsed = false"
      >
        <NMenu
          v-model:value="activeKey"
          :collapsed="collapsed"
          :collapsed-width="64"
          :collapsed-icon-size="22"
          :options="menuOptions"
        />
      </NLayoutSider>
      <NLayout>
        <RouterView />
      </NLayout>
    </NLayout>
  </div>
</template>
