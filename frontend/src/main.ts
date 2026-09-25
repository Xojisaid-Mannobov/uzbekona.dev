import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import { router } from './router'
import { installAnalytics } from './composables/analytics'
import './styles/base.css'

const app = createApp(App)
app.use(createPinia())
// Tracker router o'rnatilishidan oldin ulanadi — birinchi sahifa ko'rilishi ham hisoblanadi
installAnalytics(router)
app.use(router)
app.mount('#app')
