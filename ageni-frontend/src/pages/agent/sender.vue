<script lang="ts" setup>
import {CustomButton, Icon, themeMode} from '@/components'
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
  <div class="fb1 fbv fbje g10 root-container">
    <div class="input-container root-container fbh g12">
      <input
        type="text"
        class="input"
        @keydown="(e) => e.key === 'Enter' && !disabled && onSend()"
        placeholder="Ask anything about agent"
        v-model="query"
      />
      <CustomButton
        @click="onSend"
        :width="92"
        :height="30"
        keep-state="active"
        active-fill="var(--sender-background-color)"
        :class="disabled ? 'disabled-button' : 'normal-button'"
      >
        <div
          class="fbh fbac g12"
          :style="`color: ${
            (themeMode === 'light' && disabled) ||
            (themeMode === 'dark' && !disabled)
              ? '#121212'
              : '#ffffff'
          }`"
        >
          <Icon :name="disabled ? 'send_disabled' : 'send'" :size="16" />
        </div>
      </CustomButton>
    </div>
  </div>
</template>

<style lang="less" scoped>
.root-container {
  position: sticky;
  line-height: 1.3;
  padding-left: 40px;
  bottom: 0;
}
.input-container {
  background: var(--message-background-color);
  border: 1px solid var(--secondary-border-color);
  box-shadow: 0px 4px 8px 0px #0000000f;
  padding: 12px 16px;
  border-radius: 8px;
}
.input {
  width: 100%;
  border-radius: 8px;
  font-size: 16px;
  border: none;
  outline: none;
  box-sizing: border-box;
  background: transparent;
  color: var(--primary-text-color);
  &::placeholder {
    color: var(--secondary-label-color);
  }
}
.normal-button {
  --sender-background-color: var(--theme-color);
}
.disabled-button {
  --sender-background-color: var(--secondary-background-color);
  pointer-events: none;
}
.label {
  font-size: 12px;
  color: var(--primary-label-color);
}
.dark {
  filter: brightness(0);
}
</style>
