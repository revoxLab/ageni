<script lang="ts" setup>
import {Icon} from '@/components'
import {ElInput} from 'element-plus'
import {computed, ref, watch} from 'vue'

const query = ref('')
const props = defineProps<{
  loading: boolean
  onSend: (query: string) => unknown
}>()
const loading = ref(false)
const disabled = computed(() => {
  return !query.value || props.loading || loading.value
})

async function onSend() {
  loading.value = true
  await props.onSend(query.value)
  loading.value = false
}

watch(props, () => {
  if (props.loading) {
    query.value = ''
  }
})
</script>

<template>
  <div class="fbh fbac g10 root-container">
    <slot></slot>
    <div class="w100p input-container fbh fbae g16">
      <ElInput
        type="textarea"
        class="input"
        placeholder="Ask anything about agent"
        :autosize="{minRows: 1, maxRows: 20}"
        v-model="query"
        resize="none"
      />
      <div @click="onSend" :class="['button center hand g8', {disabled}]">
        <Icon :name="disabled ? 'send_disabled' : 'send'" :size="12" />
      </div>
    </div>
  </div>
</template>

<style lang="less" scoped>
.root-container {
  position: fixed;
  bottom: 24px;
  right: 16px;
  left: 16px;
}
.input-container {
  background: var(--message-sender-background-color);
  border: 1px solid var(--secondary-border-color);
  box-shadow: 0px 4px 12px 0px #00000029;
  border-radius: 16px;
  padding: 12px 14px;
}
.input {
  min-height: 26px;
}
.button {
  height: 28px;
  padding: 0 8px;
  border-radius: 16px;
  background: #17b09e;
  &.disabled {
    pointer-events: none;
    background: var(--tertiary-background-color);
    color: var(--primary-text-color);
  }
}
</style>
