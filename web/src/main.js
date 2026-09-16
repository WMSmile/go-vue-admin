import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import router from './router'
import { permissionDirective } from '@/utils/permission'
import { useUserStore } from '@/store/user'
import { useSysConfigStore } from '@/store/sysConfig'

const app = createApp(App)
const pinia = createPinia()
app.use(pinia)

// Restore auth state AND register the CACHED dynamic routes BEFORE installing
// the router. `app.use(router)` triggers the router's very first navigation; if
// the dynamic routes (e.g. /system/role) are not registered yet, that first
// navigation falls through to the catch-all route and renders the 404 view even
// though the URL is correct. Registering the cached routes here guarantees a
// refresh / deep link resolves to the real page.
const userStore = useUserStore(pinia)
userStore.bootstrap()

app.use(router)
app.use(ElementPlus)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}
app.directive('permission', permissionDirective)

// Load branding parameters (site name, copyright, ...) for the login page / footer.
useSysConfigStore(pinia).load()

// Freshen the menu from the server and re-register routes BEFORE mounting, so
// menus added/seeded after the last login (e.g. 参数设置) are available on a
// refresh. The cached routes registered above already cover the first
// navigation, so this only updates/extends the route table.
if (localStorage.getItem('token')) {
  userStore.rehydrate().finally(() => app.mount('#app'))
} else {
  app.mount('#app')
}
