import {showToast} from '@/components'
import {openWallet} from '@/wallet'
import {ref} from 'vue'
import {storage} from './storage'

export function sleep(time: number) {
  return new Promise((resolve) => setTimeout(resolve, time))
}

export function displayAmount(amount: number, digital: number) {
  if (amount < 10 ** -digital) {
    return amount.toPrecision(digital)
  } else {
    return amount.toFixed(digital)
  }
}

export function getUrlLogo(urlString: string) {
  const url = new URL(urlString)
  const cleanURL = `${url.protocol}//${url.hostname}/favicon.ico`
  return cleanURL
}

export function getUrlDomain(urlString: string) {
  return urlString
    .replace(/^https?:\/\//, '')
    .split('/')[0]
    .replace(/^www\./, '')
}

export function copyToClipboard(data: string) {
  window.navigator.clipboard.writeText(data)
  showToast('copied')
}

export function jump(url: string) {
  if (url.startsWith('http')) {
    window.open(url, '_black')
  } else {
    window.open(new URL(location.origin + url), '_black')
  }
}

export function shortAddress(addr?: string) {
  if (addr && addr.length > 14) {
    return (
      addr.substring(0, 6) +
      '...' +
      addr.substring(addr.length - 4, addr.length)
    )
  } else {
    return addr
  }
}

export function isMobile() {
  const agents = ['Android', 'iPhone', 'Windows Phone', 'iPod']
  return agents.some((agent) => navigator.userAgent.includes(agent))
}

export function isTest() {
  return import.meta.env.MODE === 'test'
}

export function logout() {
  window.localStorage.clear()
  window.location.href = window.location.origin
}

export function getAvatar() {
  return `https://cdn.stamp.fyi/avatar/${storage.get('address')}`
}

export function parseSmartWalletGuideInfo(params: string): {
  text: string
  icon?: string
} {
  try {
    const data = JSON.parse(params)
    if (!data.icon || !data.text) {
      throw new Error()
    } else {
      return data
    }
  } catch {
    return {
      text: params,
      icon: undefined,
    }
  }
}

export function withLogin<T extends AnyFunction>(callback: T) {
  return (...params: Parameters<T>) => {
    if (storage.get('token')) {
      return callback(...params) as ReturnType<T>
    } else {
      openWallet()
    }
  }
}

export function checkLogin<T extends AnyFunction>(callback: T) {
  if (storage.get('token')) {
    return callback() as ReturnType<T>
  } else {
    openWallet()
  }
}

export function prepare(time: number) {
  const valueRef = ref(false)
  sleep(time).then(() => (valueRef.value = true))
  return valueRef
}
