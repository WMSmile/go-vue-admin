<template>
  <div class="tags-view">
    <div class="tags-inner">
      <span
        v-for="tag in visitedViews"
        :key="tag.path"
        class="tag-item"
        :class="{ active: isActive(tag) }"
        @click="toTag(tag)"
        @contextmenu.prevent="openMenu(tag, $event)"
      >
        {{ tag.title }}
        <el-icon v-if="!tag.affix" class="close-icon" @click.stop="closeTag(tag)"><Close /></el-icon>
      </span>
    </div>

    <ul
      v-show="menuVisible"
      class="context-menu"
      :style="{ left: menuX + 'px', top: menuY + 'px' }"
    >
      <li @click="closeCurrent">关闭当前</li>
      <li @click="closeOthers">关闭其他</li>
      <li @click="closeLeft">关闭左侧</li>
      <li @click="closeRight">关闭右侧</li>
      <li @click="closeAll">全部关闭</li>
    </ul>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { Close } from '@element-plus/icons-vue'
import { useTagsViewStore } from '@/store/tagsView'
import { useUserStore } from '@/store/user'

const route = useRoute()
const router = useRouter()
const tagsStore = useTagsViewStore()
const userStore = useUserStore()
const { visitedViews } = storeToRefs(tagsStore)

const menuVisible = ref(false)
const menuX = ref(0)
const menuY = ref(0)
let selectedTag = null

const dashboardPath = ref('/')
const dashboardTitle = ref('首页')

const normalize = (p) => (!p ? '/' : p.startsWith('/') ? p : '/' + p)

function findDashboardMenu(menus) {
  let res = null
  const walk = (arr) => {
    for (const m of arr || []) {
      if (res) return
      if (m.type === 2 && (m.name === 'Dashboard' || normalize(m.path) === '/dashboard')) {
        res = { ...m, path: normalize(m.path) }
        return
      }
      if (m.children && m.children.length) walk(m.children)
    }
  }
  walk(menus)
  if (res) return res
  const walk2 = (arr) => {
    for (const m of arr || []) {
      if (res) return
      if (m.type === 2) {
        res = { ...m, path: normalize(m.path) }
        return
      }
      if (m.children && m.children.length) walk2(m.children)
    }
  }
  walk2(menus)
  return res
}

function computeDashboard() {
  const m = findDashboardMenu(userStore.menus)
  if (m) {
    dashboardPath.value = m.path
    dashboardTitle.value = m.title || '首页'
  }
}

function isActive(tag) {
  return tag.path === route.path
}

function addTag() {
  if (!route.meta || !route.meta.title) return
  tagsStore.addView({
    name: route.name,
    path: route.path,
    title: route.meta.title,
    affix: false
  })
}

function toTag(tag) {
  if (!isActive(tag)) router.push(tag.path)
}

function closeTag(tag) {
  if (tag.affix) return
  const list = [...tagsStore.visitedViews]
  const idx = list.findIndex((v) => v.path === tag.path)
  tagsStore.delView(tag)
  if (isActive(tag)) {
    const next = list[idx + 1] || list[idx - 1]
    router.push(next ? next.path : dashboardPath.value)
  }
}

function openMenu(tag, e) {
  selectedTag = tag
  const width = 120
  const height = 180
  menuX.value = e.clientX + width > window.innerWidth ? window.innerWidth - width : e.clientX
  menuY.value = e.clientY + height > window.innerHeight ? window.innerHeight - height : e.clientY
  menuVisible.value = true
}

function closeMenu() {
  menuVisible.value = false
}

function closeCurrent() {
  const tag = selectedTag
  if (tag && !tag.affix) {
    const list = [...tagsStore.visitedViews]
    const idx = list.findIndex((v) => v.path === tag.path)
    tagsStore.delView(tag)
    if (isActive(tag)) {
      const next = list[idx + 1] || list[idx - 1]
      router.push(next ? next.path : dashboardPath.value)
    }
  }
  closeMenu()
}

function closeOthers() {
  if (!selectedTag) return
  tagsStore.delOthersViews(selectedTag)
  if (!isActive(selectedTag)) router.push(selectedTag.path)
  closeMenu()
}

function closeLeft() {
  if (!selectedTag) return
  tagsStore.delLeftViews(selectedTag)
  if (!isActive(selectedTag)) router.push(selectedTag.path)
  closeMenu()
}

function closeRight() {
  if (!selectedTag) return
  tagsStore.delRightViews(selectedTag)
  if (!isActive(selectedTag)) router.push(selectedTag.path)
  closeMenu()
}

function closeAll() {
  tagsStore.delAllViews()
  router.push(dashboardPath.value)
  closeMenu()
}

onMounted(() => {
  computeDashboard()
  tagsStore.addView({ name: 'Dashboard', path: dashboardPath.value, title: dashboardTitle.value, affix: true })
  addTag()
  window.addEventListener('click', closeMenu)
})

onUnmounted(() => {
  window.removeEventListener('click', closeMenu)
})

watch(() => route.path, addTag)
</script>

<style scoped>
.tags-view {
  height: 40px;
  background: var(--el-bg-color, #fff);
  border-bottom: 1px solid var(--el-border-color-lighter, #ebeef5);
  display: flex;
  align-items: center;
  padding: 0 8px;
  position: relative;
}
.tags-inner {
  display: flex;
  align-items: center;
  gap: 6px;
  overflow-x: auto;
  white-space: nowrap;
  width: 100%;
  scrollbar-width: none;
}
.tags-inner::-webkit-scrollbar { display: none; }
.tag-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  height: 26px;
  line-height: 26px;
  padding: 0 8px;
  font-size: 13px;
  color: var(--el-text-color-regular, #606266);
  background: var(--el-fill-color-light, #f4f4f5);
  border: 1px solid var(--el-border-color-lighter, #ebeef5);
  border-radius: 3px;
  cursor: pointer;
  user-select: none;
}
.tag-item.active {
  color: #fff;
  background: var(--el-color-primary, #409eff);
  border-color: var(--el-color-primary, #409eff);
}
.close-icon {
  font-size: 12px;
  border-radius: 50%;
}
.close-icon:hover {
  background: rgba(0, 0, 0, 0.25);
  color: #fff;
}
.context-menu {
  position: fixed;
  z-index: 3000;
  margin: 0;
  padding: 5px 0;
  list-style: none;
  background: var(--el-bg-color, #fff);
  border: 1px solid var(--el-border-color-lighter, #ebeef5);
  border-radius: 4px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
  min-width: 110px;
}
.context-menu li {
  padding: 7px 16px;
  font-size: 13px;
  color: var(--el-text-color-regular, #606266);
  cursor: pointer;
}
.context-menu li:hover {
  background: var(--el-fill-color-light, #f5f7fa);
  color: var(--el-color-primary, #409eff);
}
</style>
