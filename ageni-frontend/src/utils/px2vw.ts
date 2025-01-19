import {computed, ref} from 'vue'
import {isMobile} from './chaos'

export const maxScreenWidth = 1440
export const baseWebScreenWidth = 1440
export const baseMobileScreenWidth = 375
export const windowHeight = ref(window.innerHeight)
export const windowWidth = ref(window.innerWidth)
export const scale = ref()

function onResize() {
  scale.value = Math.max(1, window.innerWidth / maxScreenWidth)
  windowHeight.value = window.innerHeight
  windowWidth.value = window.innerWidth
}

onResize()
if (window?.addEventListener) {
  window.addEventListener('resize', onResize)
}

export function px2vw(value: number) {
  const baseWidth = isMobile() ? baseMobileScreenWidth : baseWebScreenWidth
  return computed(() => `${(value / (baseWidth * scale.value)) * 100}vw`)
}

export function px2vwSync(value: number) {
  const baseWidth = isMobile() ? baseMobileScreenWidth : baseWebScreenWidth
  return `${(value / (baseWidth * scale.value)) * 100}vw`
}
