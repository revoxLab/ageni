<script lang="ts" setup>
import {px2vw} from '@/utils/px2vw'
import {requireIcon, requireIconRef} from '@/utils/require'
import {ref, watch} from 'vue'

const props = defineProps<{
  name: string
  size?: number
  width?: number
  height?: number
}>()
const src = props.name?.startsWith('http')
  ? ref(props.name)
  : requireIconRef(props.name, 'svg')
const [scaledWidth, scaledHeight] = [
  px2vw(props.width ?? props.size ?? 16),
  px2vw(props.height ?? props.size ?? 16),
]

watch(
  props,
  async () =>
    (src.value = props.name?.startsWith('http')
      ? props.name
      : await requireIcon(props.name, 'svg'))
)
</script>

<template>
  <img
    :src="src"
    :style="{width: scaledWidth, height: scaledHeight}"
    :alt="name"
    class="hand"
  />
</template>
