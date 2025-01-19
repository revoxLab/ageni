import {ref} from 'vue'

const timer = ref<NodeJS.Timeout>()

const toastVariant = ref<'base' | 'check'>()
const toastVisible = ref(false)
const toastContent = ref('')

export const showToast = (
  message: string,
  variant: 'base' | 'check' = 'base'
) => {
  clearTimeout(timer.value)
  toastVisible.value = true
  toastContent.value = message
  toastVariant.value = variant
  timer.value = setTimeout(() => {
    toastVisible.value = false
  }, 5000)
}

export const toastModal = {
  toastVisible,
  toastContent,
  toastVariant,
}
