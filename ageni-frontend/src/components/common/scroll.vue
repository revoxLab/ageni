<script lang="ts" setup>
import {holdScrollPosition} from '@/utils'
import {onMounted, onUnmounted, ref} from 'vue'

const containerRef = ref<HTMLDivElement>()
const props = defineProps<{
  keepPosition?: boolean
}>()
const emit = defineEmits<{
  onAttachTop: []
  onAttachBottom: []
}>()

defineExpose({
  onScrollToTop: () => {
    setTimeout(() => {
      containerRef.value?.scrollTo({top: 10, behavior: 'smooth'})
    })
  },
  onScrollToBottom: () => {
    setTimeout(() => {
      containerRef.value?.scrollTo({top: 1000000, behavior: 'smooth'})
    }, 100)
  },
})

const onScroll = () => {
  if (containerRef.value) {
    const {scrollTop, scrollHeight, clientHeight} = containerRef.value
    if (scrollTop === 0) {
      emit('onAttachTop')
    } else if (Math.abs(scrollHeight - scrollTop - clientHeight) < 2) {
      emit('onAttachBottom')
    }
  }
}

if (props.keepPosition) {
  let disconnect: Function | null
  onMounted(() => (disconnect = holdScrollPosition(containerRef.value!)))
  onUnmounted(() => disconnect?.())
}
</script>

<template>
  <div ref="containerRef" @scroll="onScroll"><slot></slot></div>
</template>
