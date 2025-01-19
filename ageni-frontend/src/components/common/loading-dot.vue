<script lang="ts" setup>
import {px2vw} from '@/utils/px2vw'
import {computed} from 'vue'

const props = defineProps<{
  color?: string
  size?: number
}>()
const size = px2vw(props.size ?? 40)
const color = computed(() => props.color ?? 'var(--primary-text-color)')
</script>

<template>
  <div class="loader"></div>
</template>

<style scoped>
.loader {
  width: v-bind(size);
  aspect-ratio: 2;
  --_g: no-repeat radial-gradient(circle closest-side, v-bind(color) 90%, #0000);
  background:
    var(--_g) 0% 50%,
    var(--_g) 50% 50%,
    var(--_g) 100% 50%;
  background-size: calc(100% / 3) 50%;
  animation: animation 1s infinite linear;
}
@keyframes animation {
  20% {
    background-position:
      0% 0%,
      50% 50%,
      100% 50%;
  }
  40% {
    background-position:
      0% 100%,
      50% 0%,
      100% 50%;
  }
  60% {
    background-position:
      0% 50%,
      50% 100%,
      100% 0%;
  }
  80% {
    background-position:
      0% 50%,
      50% 50%,
      100% 100%;
  }
}
</style>
