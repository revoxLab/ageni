<script lang="ts" setup>
import {windowHeight} from '@/utils'
import {computed, ref} from 'vue'
import {Icon} from '..'

const show = ref(false)
const listRef = ref<HTMLDivElement>()
const height = computed(() => {
  if (listRef.value) {
    const {y} = listRef.value.getBoundingClientRect()
    return `${windowHeight.value - y - 48}px`
  }
  return 0
})

defineProps<{
  label?: string
  selected?: SelectOption
  useCustomInput?: boolean
  options: SelectOption[]
}>()
defineEmits<{
  onSelectedChange: [option: SelectOption]
}>()
</script>

<template>
  <div class="pr">
    <div v-if="show" class="interactive" @click="show = false"></div>
    <template v-if="!useCustomInput">
      <p v-if="label" class="label">{{ label }}</p>
      <div class="select-input fbh fbjsb fbac g4 hand" @click="show = true">
        <Icon v-if="selected?.icon" :name="selected.icon" :size="24" />
        <p class="fb1">{{ selected?.label }}</p>
        <Icon name="arrow-down" :size="12" />
      </div>
    </template>
    <div v-else class="hand" @click="show = true">
      <slot :show="show"></slot>
    </div>
    <div
      v-if="show"
      ref="listRef"
      class="list pa oa fbv p16"
      @focusout="show = false"
    >
      <template v-for="(option, index) in options">
        <div v-if="index !== 0" class="divider mt8 mb8"></div>
        <div
          class="fbh fbac g8 hand"
          @click="$emit('onSelectedChange', option), (show = false)"
        >
          <Icon v-if="option.icon" :name="option.icon" :size="24" />
          <div
            :class="[
              'list-item fb1 fbh fbac fbjsb py4 hand',
              {
                disabled: option.disabled,
                selected: option.key === selected?.key,
              },
            ]"
          >
            <p class="fb1 omit">{{ option.label }}</p>
            <Icon v-if="option.key === selected?.key" name="check" :size="18" />
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<style lang="less" scoped>
.interactive {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}
.label {
  font-size: 14px;
  color: var(--primary-label-color);
  padding-bottom: 8px;
}
.select-input {
  padding: 8px 16px;
  border-radius: 8px;
  border: 1px solid var(--secondary-border-color);
  background-color: transparent;
  font-size: 16px;
}
.list {
  width: 100%;
  z-index: 99;
  box-sizing: border-box;
  max-height: v-bind(height);
  background: var(--modal-background-color);
  border: 1px solid var(--secondary-border-color);
  border-radius: 8px;
  font-size: 16px;
  margin-top: 12px;
  .list-item {
    opacity: 0.6;
    &:hover {
      opacity: 1;
    }
    &.selected {
      opacity: 1;
    }
    &.disabled {
      opacity: 0.5;
      pointer-events: none;
    }
  }
}
</style>
