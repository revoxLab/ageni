<script lang="ts" setup>
import markdown from 'markdown-it'
import {ref, watch} from 'vue'

const props = defineProps<{
  data: string
  animation?: boolean
}>()
const contentRef = ref<HTMLElement>()
const contentHtml = ref('')
const md = markdown()

let timeout: NodeJS.Timeout

watch(
  props,
  () => {
    clearTimeout(timeout)
    const result = md.render(props.data)
    const update = () => {
      contentHtml.value = result.slice(0, contentHtml.value.length + 5)
      if (contentHtml.value.length !== result.length) {
        timeout = setTimeout(update, 20)
        setTimeout(() => {
          contentRef.value?.querySelectorAll('a').forEach((el) => {
            el.setAttribute('target', '_blank')
          })
        }, 100)
      }
    }
    if (props.animation) {
      contentHtml.value = ''
      update()
    } else {
      contentHtml.value = result
      setTimeout(() => {
        contentRef.value?.querySelectorAll('a').forEach((el) => {
          el.setAttribute('target', '_blank')
        })
      }, 100)
    }
  },
  {immediate: true}
)
</script>

<template>
  <div ref="contentRef" class="fbv g12 h100p">
    <div class="summary fbv g12" v-html="contentHtml"></div>
  </div>
</template>

<style lang="less" scoped>
.summary {
  :deep(h1),
  :deep(h2),
  :deep(h3),
  :deep(h4),
  :deep(h5),
  :deep(h6) {
    margin-bottom: 6px;
  }
  :deep(*) {
    font-size: 16px;
    line-height: 24px;
  }
  :deep(img) {
    max-width: 800px;
    max-height: 300px;
  }
  :deep(li) {
    margin-left: 20px;
    margin-bottom: 6px;
  }
  :deep(a) {
    color: #17b09d;
  }
}
</style>
