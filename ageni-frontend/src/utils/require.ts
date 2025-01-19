import {themeMode} from '@/components'
import {Ref, ref, watch} from 'vue'

export function preloadImage(src: string | Ref<string>) {
  const img = document.createElement('img')
  img.onload = () => document.body.removeChild(img)
  img.style.position = 'absolute'
  img.style.height = '0'
  img.style.width = '0'
  if (typeof src === 'string') {
    img.src = src
    document.body.append(img)
  } else {
    watch(src, (src) => {
      img.src = src
      document.body.append(img)
    })
  }
}

export async function requireIcon(name: string, extension = 'svg') {
  try {
    const result = await import(
      `@/assets/icons/${themeMode.value}/${name}.${extension}`
    )
    return result.default
  } catch {
    const result = await import(`@/assets/icons/${name}.${extension}`)
    return result.default
  }
}

export async function requireImage(name: string, extension = 'png') {
  try {
    const result = await import(
      `@/assets/images/${themeMode.value}/${name}.${extension}`
    )
    return result.default
  } catch {
    const result = await import(`@/assets/images/${name}.${extension}`)
    return result.default
  }
}

export function requireIconRef(name: string, extension = 'svg') {
  const src = ref<string>('')
  watch(
    themeMode,
    () => requireIcon(name, extension).then((path) => (src.value = path)),
    {immediate: true}
  )
  return src
}

export function requireImageRef(name: string, extension = 'png') {
  const src = ref<string>('')
  watch(
    themeMode,
    () => requireImage(name, extension).then((path) => (src.value = path)),
    {immediate: true}
  )
  return src
}

export function checkImage(url: string, callback: (url?: string) => void) {
  var img = new Image()
  img.onload = () => callback(url)
  img.onerror = () => callback()
  img.src = url
}
