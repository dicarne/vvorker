<script setup lang="ts">
import { ref } from 'vue'
import { RouterView } from 'vue-router'
import type { MenuOption } from 'naive-ui'
import { NLayout, NLayoutSider, NMenu } from 'naive-ui'
import {
  CalendarWorkWeek24Regular as WorkersIcon,
  Status24Regular as StatusIcon,
  TaskListSquareLtr24Regular as TaskIcon,
  Database24Regular as SQLIcon,
  Cloud24Regular as OSSIcon,
  Braces24Filled as KVIcon,
  Person24Regular as UserIcon,
} from '@vicons/fluent'
import TheHeader from '@/components/TheHeader.vue'
import { renderIcon } from '@/utils/render'
const collapsed = ref<boolean>(true)
const activeKey = ref<string>('')
const menuOptions: MenuOption[] = [
  {
    label: 'Workers',
    key: 'workers',
    icon: renderIcon(WorkersIcon),
  },
  {
    label: 'Status',
    key: 'status',
    icon: renderIcon(StatusIcon),
  },
  {
    label: 'Task',
    key: 'task',
    icon: renderIcon(TaskIcon),
  },
  {
    label: 'SQL',
    key: 'sql',
    icon: renderIcon(SQLIcon),
  },
  {
    label: 'OSS',
    key: 'oss',
    icon: renderIcon(OSSIcon),
  },
  {
    label: 'KV',
    key: 'kv',
    icon: renderIcon(KVIcon),
  },
  {
    label: 'User',
    key: 'user',
    icon: renderIcon(UserIcon),
  },
]
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
