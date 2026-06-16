import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import { createAppRouter } from './runtime/router/index.ts'

const router = createAppRouter();
const app = createApp(App)

app.use(router)
app.mount('#app')