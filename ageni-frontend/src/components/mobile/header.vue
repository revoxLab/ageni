<script lang="ts" setup>
import {onMounted, onUnmounted, ref} from 'vue'

const observer = ref<ResizeObserver>()
const boxRef = ref<HTMLDivElement>()
const height = ref(0)

onUnmounted(() => {
  observer.value?.disconnect()
})
onMounted(() => {
  observer.value = new ResizeObserver((entries) => {
    for (let entry of entries) {
      if (entry.target === boxRef.value) {
        height.value = (boxRef.value?.clientHeight ?? 0) + 1
      }
    }
  })
  observer.value.observe(boxRef.value!)
})
</script>

<template>
  <div>
    <div v-bind="$attrs" class="title-bar" ref="boxRef">
      <slot></slot>
    </div>
    <div :style="`height: ${height}px`"></div>
  </div>
</template>

<style lang="less" scoped>
.title-bar {
  z-index: 1;
  background: var(--primary-background-color);
  position: fixed;
  right: 0;
  left: 0;
  top: 0;
}
</style>
