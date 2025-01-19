<script lang="ts" setup>
import {toastModal} from '../toast'
import Icon from './icon.vue'
const {toastVisible, toastContent, toastVariant} = toastModal
</script>

<template>
  <Teleport to="#modal">
    <div
      v-if="toastVariant === 'check'"
      class="toast-check animation pf py13 px16 fbh fbac g12 f14 br8"
      :class="toastVisible ? 'ne show' : ''"
    >
      <Icon name="check2" :size="20" />
      {{ toastContent }}
      <Icon
        name="close"
        class="p5"
        :size="10"
        style="pointer-events: auto"
        @click="toastVisible = false"
      />
    </div>
    <div
      v-else
      class="toast-base animation pf center f14 py12 px24 br8"
      :class="toastVisible ? 'ne show' : ''"
    >
      {{ toastContent }}
    </div>
  </Teleport>
</template>

<style lang="css" scoped>
.animation {
  left: 50%;
  z-index: 99999;
  transform: translate(-50%, -50%);
  transition: opacity 0.5s;
  opacity: 0;
  &.show {
    opacity: 1;
  }
}
.toast-base {
  top: 50%;
  max-width: 300px;
  background: rgba(0, 0, 0, 0.8);
  color: #ffffff;
}
.toast-check {
  border: 1px solid var(--secondary-border-color);
  background-color: var(--reverse-text-color);
  top: 80px;
}
</style>
