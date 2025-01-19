<script lang="ts" setup>
import {px2vw} from '@/utils'
import {computed} from 'vue'

const props = defineProps<{
  color?: string
  size?: number
}>()
const background = computed(() => props.color ?? '#ffa516')
const offset = px2vw((props.size ?? 30) / 6)
const size = px2vw(props.size ?? 30)
</script>

<template>
  <div class="loader"></div>
</template>

<style scoped>
.loader {
  width: v-bind(size);
  aspect-ratio: 1;
  border-radius: 50%;
  background:
    radial-gradient(farthest-side, v-bind(background) 94%, #0000)
      top/v-bind(offset) v-bind(offset) no-repeat,
    conic-gradient(#0000 30%, v-bind(background));
  mask: radial-gradient(
    farthest-side,
    #0000 calc(100% - v-bind(offset)),
    #000 0
  );
  animation: animation 1s infinite linear;
}
@keyframes animation {
  100% {
    transform: rotate(1turn);
  }
}
</style>
