<script setup lang="ts">
import { NTooltip } from 'naive-ui'

const props = defineProps<{
  data: {
    color: string
    tooltip: string
  }[]
}>()
</script>
<template>
  <div class="v-flex-start-center" style="flex-wrap: nowrap; gap: 3px">
    <div v-for="(item, index) in data" :key="index">
      <NTooltip trigger="hover">
        <template #trigger>
          <div
            class="v-track-item"
            :class="
              index === 0
                ? 'v-track-item-left'
                : index === data.length - 1
                  ? 'v-track-item-right'
                  : ''
            "
            :style="{ backgroundColor: item.color }"
          ></div>
        </template>
        {{ item.tooltip }}
      </NTooltip>
    </div>
  </div>
</template>
<style scoped>
.v-track-item {
  height: 40px;
  width: 8px;
}
.v-track-item-left {
  border-radius: 3px 0 0 3px;
}
.v-track-item-right {
  border-radius: 0 3px 3px 0;
}
</style>
