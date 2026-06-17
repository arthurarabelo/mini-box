import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './style.css'
import App from './App.vue'
import { createAppRouter } from './runtime/router/index.ts'
import { setupAuthGuard } from './runtime/router/guard.ts'

const router = createAppRouter();
const pinia = createPinia()
const app = createApp(App)

setupAuthGuard(router)
app.use(pinia)
app.use(router)
app.mount('#app')