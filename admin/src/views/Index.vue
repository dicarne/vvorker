<script setup lang="ts">
import { ref } from 'vue'
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

const collapsed = ref<boolean>(true)
const activeKey = ref<string>('')
const menuOptions: MenuOption[] = [
  {
    label: renderMenuRouterLink('Workers', 'Workers'),
    key: 'workers',
    icon: renderIcon(WorkersIcon),
  },
  {
    label: renderMenuRouterLink('Status', 'Status'),
    key: 'status',
    icon: renderIcon(StatusIcon),
  },
  {
    label: renderMenuRouterLink('Task', 'Task'),
    key: 'task',
    icon: renderIcon(TaskIcon),
  },
  {
    label: renderMenuRouterLink('PGSQL', 'PGSQL'),
    key: 'pgsql',
    icon: renderIcon(PGSQLIcon),
  },
  {
    label: renderMenuRouterLink('MySQL', 'MySQL'),
    key: 'mysql',
    icon: renderIcon(MYSQLIcon),
  },
  {
    label: renderMenuRouterLink('OSS', 'OSS'),
    key: 'oss',
    icon: renderIcon(OSSIcon),
  },
  {
    label: renderMenuRouterLink('KV', 'KV'),  
    key: 'kv',
    icon: renderIcon(KVIcon),
  },
  {
    label: renderMenuRouterLink('User', 'User'),
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
