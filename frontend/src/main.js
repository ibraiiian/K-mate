import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'

const app = createApp(App)

// Register Pinia (State Management)
app.use(createPinia())

// Register Router
app.use(router)

// Mount App
app.mount('#app')
