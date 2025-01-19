import '@/assets/style/index.less'
import 'element-plus/dist/index.css'

import ElementPlus from 'element-plus'
import {createApp} from 'vue'
import {version} from '../package.json'
import App from './main.vue'
import {router} from './router'
import {isMobile} from './utils/chaos.js'

const app = createApp(App)

app.use(router)
app.use(ElementPlus)
console.log(`version: ${version}`)

if (isMobile()) {
  import('@/mobile/style/index.less').then(() => app.mount('#app'))
} else {
  import('@/pages/style/index.less').then(() => app.mount('#app'))
}
