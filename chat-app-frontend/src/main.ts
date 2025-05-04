import { createApp } from 'vue'
import "@/app/styles/index.css"
import App from './app/App.vue'

import { router } from '@/pages/router/Router'
import { createPinia } from 'pinia'

createApp(App)
.use(createPinia())
.use(router)
.mount('#app')
