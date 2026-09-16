import { defineStore } from 'pinia'
import { getConfigMap } from '@/api'

// Holds the system parameter key/value map, used for branding (site name,
// copyright, login title). Loaded once at app boot from the public map endpoint.
export const useSysConfigStore = defineStore('sysConfig', {
  state: () => ({
    map: JSON.parse(localStorage.getItem('sysConfig') || '{}')
  }),
  getters: {
    get: (state) => (key, def = '') =>
      state.map && state.map[key] != null && state.map[key] !== '' ? state.map[key] : def
  },
  actions: {
    async load() {
      try {
        const map = await getConfigMap()
        if (map && typeof map === 'object') {
          this.map = map
          localStorage.setItem('sysConfig', JSON.stringify(map))
          if (map['site.name']) document.title = map['site.name']
        }
      } catch (e) {
        // ignore: branding falls back to defaults
      }
    }
  }
})
