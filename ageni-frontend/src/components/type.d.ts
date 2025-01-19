type SelectOption<T = any> = {
  key: Meta
  label: string
  value: T
  icon?: string
  disabled?: boolean
}

type ScrollShape = {
  onScrollToTop: () => void
  onScrollToBottom: () => void
}
