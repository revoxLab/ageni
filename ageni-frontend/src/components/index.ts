import {isMobile} from '@/utils'
import CustomButton from './common/button.vue'
import IconGroup from './common/icon-group.vue'
import Icon from './common/icon.vue'
import Image from './common/image.vue'
import TextInput from './common/input.vue'
import LoadingCircle from './common/loading-circle.vue'
import LoadingDot from './common/loading-dot.vue'
import Scroll from './common/scroll.vue'
import WebSelect from './common/select.vue'
import SliderInput from './common/slider-input.vue'
import WebTab from './common/tab.vue'
import WebToast from './common/toast.vue'
import MobileHeader from './mobile/header.vue'
import MobileSelect from './mobile/select.vue'
import MobileTab from './mobile/tab.vue'
import MobileToast from './mobile/toast.vue'

export * from './chaos'
export * from './theme'
export {showToast} from './toast'
export {
  CustomButton,
  Icon,
  IconGroup,
  Image,
  LoadingCircle,
  LoadingDot,
  MobileHeader,
  Scroll,
  SliderInput,
  TextInput,
}

export const Select = isMobile() ? MobileSelect : WebSelect
export const TabSwitch = isMobile() ? MobileTab : WebTab
export const Toast = isMobile() ? MobileToast : WebToast
