<script lang="ts" setup>
import {isMobile, px2vw, sleep} from '@/utils'
import {computed, ref} from 'vue'
import {themeMode, themeStyleFilter} from '../theme'

const props = defineProps<{
  width: number
  height: number
  disabled?: boolean
  themeStyle?: boolean
  keepState?: 'active' | 'inactive'
  defaultStroke?: string
  defaultColor?: string
  defaultFill?: string
  activeStroke?: string
  activeColor?: string
  activeFill?: string
  fullWidth?: boolean
}>()
const {
  width,
  height,
  fullWidth,
  defaultStroke = '#17b09e',
  defaultColor = '#17b09e',
  defaultFill = 'transparent',
  activeStroke = 'none',
  activeFill = '#17b09e',
} = props
const [scaledWidth, scaledHeight] = [
  fullWidth ? '100%' : px2vw(width),
  px2vw(height),
]
const hovered = ref(props.keepState === 'active' || false)
const mobile = isMobile()
const activeColor = computed(() =>
  props.activeColor
    ? props.activeColor
    : themeMode.value === 'light'
      ? 'white'
      : 'black'
)
const corner = Math.min(width, height) / 5

async function toggle(value: boolean) {
  if (value) {
    hovered.value = true
  } else {
    await sleep(isMobile() ? 200 : 50)
    hovered.value = false
  }
}
</script>

<template>
  <div :class="{disabled}">
    <div
      :class="['button-container center pr hand', {hovered}]"
      @mouseenter="() => !keepState && !mobile && toggle(true)"
      @mouseleave="() => !keepState && !mobile && toggle(false)"
      @touchstart="() => !keepState && mobile && toggle(true)"
      @touchend="() => !keepState && mobile && toggle(false)"
    >
      <svg
        class="svg"
        :viewBox="`0 0 ${width} ${height}`"
        preserveAspectRatio="none"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
      >
        <path
          :d="`M0.5 ${corner}L${corner} 0.5H${width - 1}V${height - corner - 1}L${width - corner - 1} ${height - 1}H0.5V${corner}Z`"
          :stroke="!hovered ? defaultStroke : activeStroke"
          :fill="!hovered ? defaultFill : activeFill"
        />
      </svg>
      <div
        :style="themeStyle && hovered ? themeStyleFilter() : ''"
        style="z-index: 0"
      >
        <slot :hovered="hovered"></slot>
      </div>
    </div>
  </div>
</template>

<style lang="less" scoped>
.button-container {
  width: v-bind(scaledWidth);
  height: v-bind(scaledHeight);
  color: v-bind(defaultColor);
  &.hovered {
    color: v-bind(activeColor);
  }
}
.disabled {
  pointer-events: none;
  filter: grayscale(1) opacity(0.6);
}
.svg {
  width: 100%;
  height: 100%;
  position: absolute;
  left: 0;
  top: 0;
}
</style>
