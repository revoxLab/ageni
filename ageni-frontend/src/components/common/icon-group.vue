<script lang="ts" setup>
import {ElTooltip} from 'element-plus'
import Icon from './icon.vue'

defineProps<{data: {name: string; image: string}[]}>()
</script>

<template>
  <ElTooltip>
    <div class="fbh" v-bind="$attrs">
      <Icon
        v-for="plugin in data.slice(0, 4)"
        style="margin-left: -5px"
        :name="plugin.image"
        :size="18"
      />
      <Icon
        v-if="data.length > 4"
        style="margin-left: -5px"
        name="more3"
        :size="18"
      />
    </div>
    <template #content>
      <div class="icon-group">
        <div v-for="plugin in data.slice(0, 8)" class="fbh fbac g4">
          <Icon :name="plugin.image" :size="18" />
          <p class="f14 label1">{{ plugin.name }}</p>
        </div>
        <p v-if="data.length > 8" class="f14 label1">...</p>
      </div>
    </template>
  </ElTooltip>
</template>

<style lang="less" scoped>
.icon-group {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  column-gap: 24px;
  row-gap: 12px;
}
img {
  border-radius: 4px;
}
</style>
