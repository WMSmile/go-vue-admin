<template>
  <el-menu
    :default-active="activeMenu"
    router
    :collapse="collapse"
    :collapse-transition="false"
    class="el-menu-vertical"
  >
    <sidebar-menu :items="menus" />
  </el-menu>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useUserStore } from '@/store/user'
import SidebarMenu from './SidebarMenu.vue'

defineProps({
  collapse: { type: Boolean, default: false }
})

const route = useRoute()
const userStore = useUserStore()
const { menus } = storeToRefs(userStore)
const activeMenu = computed(() => route.path)
</script>

<style scoped>
.el-menu-vertical { border-right: none; height: 100%; }
</style>
