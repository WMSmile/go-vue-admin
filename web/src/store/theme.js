import { defineStore } from 'pinia'

const STORAGE_KEY = 'theme-mode'

function systemPrefersDark() {
  return window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches
}

// Apply the chosen mode to <html class="dark">. Element Plus dark CSS vars react to it.
function apply(mode) {
  const isDark = mode === 'dark' || (mode === 'system' && systemPrefersDark())
  document.documentElement.classList.toggle('dark', isDark)
}

export const useThemeStore = defineStore('theme', {
  state: () => ({
    mode: localStorage.getItem(STORAGE_KEY) || 'system'
  }),
  actions: {
    init() {
      apply(this.mode)
      if (window.matchMedia) {
        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
          if (this.mode === 'system') apply('system')
        })
      }
    },
    setMode(mode) {
      this.mode = mode
      localStorage.setItem(STORAGE_KEY, mode)
      apply(mode)
    }
  }
})
