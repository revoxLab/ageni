<script lang="ts" setup>
import {toastModal} from '../toast'
const {toastVisible, toastContent} = toastModal
</script>

<template>
  <div class="toast" v-if="toastVisible">
    <div class="toast_txt">{{ toastContent }}</div>
  </div>
</template>

<style lang="css" scoped>
.toast {
  display: flex;
  position: fixed;
  left: 50%;
  top: 50%;
  padding: 12px 24px;
  max-width: 300px;
  justify-content: center;
  align-items: center;
  background: rgba(0, 0, 0, 0.8);
  border-radius: 8px;
  font-family: 'Roboto Condensed';
  font-style: normal;
  font-weight: 400;
  line-height: 120%;
  font-size: 14px;
  letter-spacing: 0.01em;
  color: #ffffff;
  transform: translate(-50%, -50%);
  z-index: 99999;
  opacity: 1;
  transition: 0.5s;
  &.end {
    opacity: 0;
  }
}
</style>
