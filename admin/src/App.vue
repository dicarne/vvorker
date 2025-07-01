<script setup lang="ts">
import { RouterView } from 'vue-router'
import { NConfigProvider, NMessageProvider } from 'naive-ui'
import nordTheme from '@/assets/nord.json'
</script>

<template>
  <NConfigProvider :theme-overrides="nordTheme">
  <NMessageProvider>
    <RouterView />
  </NMessageProvider>
  </NConfigProvider>
</template>
