<script lang="ts" setup>
import {Toast} from '@/components'
import VConsole from 'vconsole'
import {onMounted} from 'vue'
import MobilePage from './mobile/index.vue'
import WebPage from './pages/index.vue'
import {isMobile, isTest} from './utils'

onMounted(() => {
  isMobile() && isTest() && new VConsole()
})
</script>

<template>
  <MobilePage v-if="isMobile()" />
  <WebPage v-else />
  <Toast />
</template>
