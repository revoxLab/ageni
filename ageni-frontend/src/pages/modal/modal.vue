<script lang="ts" setup>
import {Icon} from '@/components'

defineProps<{open?: boolean}>()
defineEmits<{close: []}>()
</script>

<template>
  <Teleport to="#modal">
    <div
      v-if="open"
      class="modal wh100p fbv fbjc fbac pa"
      @click="$emit('close')"
    >
      <div class="card pr" v-bind="$attrs" @click.stop="">
        <div class="border1"></div>
        <div class="border2"></div>
        <Icon @click="$emit('close')" name="close" class="icon" :size="14" />
        <slot></slot>
      </div>
    </div>
  </Teleport>
</template>

<style lang="less" scoped>
.modal {
  background: #12121299;
  z-index: 999;
  left: 0;
  top: 0;
}
.icon {
  position: absolute;
  z-index: 1;
  right: 18px;
  top: 18px;
}
.card {
  .border1 {
    position: absolute;
    border-left: 1px solid var(--secondary-border-color);
    height: 30px * sqrt(2);
    rotate: 45deg;
    left: 16px;
    top: -6px;
  }
  .border2 {
    position: absolute;
    border-left: 1px solid var(--secondary-border-color);
    height: 30px * sqrt(2);
    rotate: 45deg;
    right: 12px;
    bottom: -8px;
  }
  &::before {
    z-index: -1;
    border: 1px solid var(--secondary-border-color);
    background: var(--modal-background-color);
    position: absolute;
    content: '';
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    clip-path: polygon(
      0 30px,
      30px 0,
      100% 0,
      100% calc(100% - 30px),
      calc(100% - 30px) 100%,
      0 100%
    );
  }
}
</style>
