import { createApp } from "vue"
import { createPinia } from "pinia"
import axios from 'axios'

import App from "./App.vue"
import router from "./router"

import "./assets/main.scss"

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount("#app")

axios.defaults.baseURL = import.meta.env.VITE_API_BASE_URL;
