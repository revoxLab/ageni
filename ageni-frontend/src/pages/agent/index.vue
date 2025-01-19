<script lang="ts" setup>
import {computed} from 'vue'
import {useRoute} from 'vue-router'
import Detail from './detail.vue'
import Market from './market.vue'

const route = useRoute()
const detailId = computed(() => Number(route.query.id))
</script>

<template>
  <Detail v-if="detailId" :id="detailId" />
  <Market v-else />
</template>
