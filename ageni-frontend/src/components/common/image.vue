<script lang="ts" setup>
import {px2vw, requireImage, requireImageRef} from '@/utils'
import {isNumber} from 'lodash-es'
import {ref, watch} from 'vue'

const props = defineProps<{
  name: string
  alt?: string
  width?: number | string
  height?: number | string
}>()
const src = props.name?.startsWith('http')
  ? ref(props.name)
  : requireImageRef(props.name, 'png')
const {width = '100%', height = '100%'} = props
const scaledWidth = isNumber(width) ? px2vw(width) : width
const scaledHeight = isNumber(height) ? px2vw(height) : height

watch(
  props,
  async () =>
    (src.value = props.name?.startsWith('http')
      ? props.name
      : await requireImage(props.name, 'png'))
)
</script>

<template>
  <img
    :src="src"
    :style="{width: scaledWidth, height: scaledHeight}"
    :alt="alt ?? name"
  />
</template>
