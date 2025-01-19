<script lang="ts" setup>
import {Icon} from '@/components'
import {ElTooltip} from 'element-plus'
import {ref} from 'vue'

const {initialOpen} = defineProps<{
  title: string
  initialOpen: boolean
}>()
const open = ref(initialOpen)
</script>

<template>
  <div class="fbv g16">
    <div class="fbh fbac fbjsb">
      <div class="fbh fbac g8">
        <p class="f15 bold">{{ title }}</p>
        <Icon
          @click="open = !open"
          name="arrow-down1"
          :class="{ne: !initialOpen}"
          :style="`rotate:${open ? 0 : -90}deg`"
          :size="12"
        />
      </div>
      <ElTooltip :visible="open ? false : undefined" content="Coming Soon">
        <slot name="action"></slot>
      </ElTooltip>
    </div>
    <template v-if="open">
      <slot></slot>
      <div class="divider"></div>
    </template>
  </div>
</template>
