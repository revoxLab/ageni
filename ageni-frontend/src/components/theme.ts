import {ref, watch} from 'vue'

export const themeMode = ref<'light' | 'dark'>('dark')

export function themeStyleFilter() {
  return themeMode.value === 'light'
    ? 'filter: grayscale(1) brightness(2);'
    : 'filter: brightness(0);'
}

watch(
  themeMode,
  () => {
    if (themeMode.value === 'light') {
      if (document.body.classList.contains('dark-theme')) {
        document.body.classList.replace('dark-theme', 'light-theme')
      } else {
        document.body.classList.add('light-theme')
      }
      localStorage.setItem('theme', 'light')
    } else if (themeMode.value === 'dark') {
      if (document.body.classList.contains('light-theme')) {
        document.body.classList.replace('light-theme', 'dark-theme')
      } else {
        document.body.classList.add('dark-theme')
      }
      localStorage.setItem('theme', 'dark')
    }
  },
  {immediate: true}
)
