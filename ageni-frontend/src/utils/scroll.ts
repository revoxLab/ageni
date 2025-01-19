export function holdScrollPosition(target: Element) {
  let scrollTop: number
  const listener = () => {
    scrollTop = target.scrollTop
  }
  target.addEventListener('scroll', listener)
  const disconnect = watchScrollHeight(target, (scrollHeightOffset) => {
    target.scrollTop = scrollTop + scrollHeightOffset
  })
  return () => {
    disconnect()
    target.removeEventListener('scroll', listener)
  }
}

function watchScrollHeight(
  target: Element,
  onChange: (scrollHeightOffset: number) => void
) {
  let prevScrollHeight: number
  const observer = new MutationObserver(() => {
    const nextScrollHeight = target.scrollHeight
    if (nextScrollHeight !== prevScrollHeight) {
      onChange(nextScrollHeight - prevScrollHeight)
      prevScrollHeight = nextScrollHeight
    }
  })
  observer.observe(target, {
    attributes: true,
    childList: true,
    subtree: true,
  })
  return () => {
    observer.disconnect()
  }
}
