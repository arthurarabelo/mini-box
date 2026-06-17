import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './style.css'
import App from './App.vue'
import { createAppRouter } from './runtime/router/index.ts'
import { setupAuthGuard } from './runtime/router/guard.ts'
import { authService } from './runtime/services/authService'
import { useAppsStore } from './pkg/stores/apps'

import filesApp from './apps/files'

const router = createAppRouter();
const pinia = createPinia()
const app = createApp(App)

app.use(pinia)

const appsStore = useAppsStore()
for (const pluginApp of [filesApp]) {
  for (const route of pluginApp.routes) {
    router.addRoute('layout', route)
  }
  appsStore.register(pluginApp)
}

setupAuthGuard(router)
authService.init(router)

app.use(router)
app.mount('#app')