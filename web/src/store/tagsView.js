import { defineStore } from 'pinia'

// Tracks the set of pages the user has opened, rendered as a tab bar.
export const useTagsViewStore = defineStore('tagsView', {
  state: () => ({
    visitedViews: []
  }),
  actions: {
    addView(view) {
      if (this.visitedViews.some((v) => v.path === view.path)) return
      this.visitedViews.push(
        Object.assign({}, view, { title: view.title || '未命名' })
      )
    },
    delView(view) {
      const i = this.visitedViews.findIndex((v) => v.path === view.path)
      if (i > -1) this.visitedViews.splice(i, 1)
    },
    delOthersViews(view) {
      this.visitedViews = this.visitedViews.filter((v) => v.affix || v.path === view.path)
    },
    delLeftViews(view) {
      const i = this.visitedViews.findIndex((v) => v.path === view.path)
      if (i > -1) {
        this.visitedViews = this.visitedViews.filter((v, idx) => v.affix || idx >= i)
      }
    },
    delRightViews(view) {
      const i = this.visitedViews.findIndex((v) => v.path === view.path)
      if (i > -1) {
        this.visitedViews = this.visitedViews.filter((v, idx) => v.affix || idx <= i)
      }
    },
    delAllViews() {
      this.visitedViews = this.visitedViews.filter((v) => v.affix)
    }
  }
})
