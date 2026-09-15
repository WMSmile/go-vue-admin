<template>
  <div>
    <div class="toolbar">
      <el-button type="primary" @click="openDialog()">新增菜单</el-button>
    </div>
    <el-table :data="list" border stripe row-key="id" :tree-props="{ children: 'children' }" default-expand-all>
      <el-table-column prop="title" label="名称" />
      <el-table-column prop="type" label="类型" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.type === 1">目录</el-tag>
          <el-tag v-else-if="row.type === 2" type="success">菜单</el-tag>
          <el-tag v-else type="warning">按钮</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="path" label="路径" />
      <el-table-column prop="component" label="组件" />
      <el-table-column prop="permission" label="权限标识" />
      <el-table-column prop="api" label="接口" />
      <el-table-column prop="method" label="方法" width="90" />
      <el-table-column prop="sort" label="排序" width="80" />
      <el-table-column label="操作" width="160">
        <template #default="{ row }">
          <el-button size="small" @click="openDialog(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="visible" :title="form.id ? '编辑菜单' : '新增菜单'" width="520px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="类型">
          <el-radio-group v-model="form.type">
            <el-radio :value="1">目录</el-radio>
            <el-radio :value="2">菜单</el-radio>
            <el-radio :value="3">按钮</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="上级菜单">
          <el-select v-model="form.parentId" clearable placeholder="顶级菜单" style="width: 100%">
            <el-option v-for="m in parentOptions" :key="m.id" :label="m.title" :value="m.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="名称">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="路由名">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="图标">
          <el-input v-model="form.icon" />
        </el-form-item>
        <el-form-item label="路径">
          <el-input v-model="form.path" placeholder="目录/菜单为路由路径，如 /system 或 user" />
        </el-form-item>
        <el-form-item label="组件" v-if="form.type !== 3">
          <el-input v-model="form.component" placeholder="如 system/user/index" />
        </el-form-item>
        <el-form-item label="权限标识" v-if="form.type === 3">
          <el-input v-model="form.permission" placeholder="如 user:add" />
        </el-form-item>
        <el-form-item label="接口" v-if="form.type !== 1">
          <el-input v-model="form.api" placeholder="如 /api/v1/users" />
        </el-form-item>
        <el-form-item label="方法" v-if="form.type !== 1">
          <el-select v-model="form.method" clearable style="width: 100%">
            <el-option label="GET" value="GET" />
            <el-option label="POST" value="POST" />
            <el-option label="PUT" value="PUT" />
            <el-option label="DELETE" value="DELETE" />
          </el-select>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" @click="submit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listMenus, createMenu, updateMenu, deleteMenu } from '@/api'

const list = ref([])
const all = ref([])
const visible = ref(false)
const form = reactive({
  id: null, parentId: null, name: '', title: '', icon: '', path: '',
  component: '', type: 2, permission: '', api: '', method: '', sort: 0, status: 1
})

const parentOptions = computed(() => all.value.filter((m) => m.type !== 3 && m.id !== form.id))

function toTree(flat) {
  const map = {}
  flat.forEach((m) => (map[m.id] = { ...m, children: [] }))
  const tree = []
  flat.forEach((m) => {
    if (m.parentId && map[m.parentId]) map[m.parentId].children.push(map[m.id])
    else tree.push(map[m.id])
  })
  return tree
}

async function load() {
  const flat = await listMenus()
  all.value = flat
  list.value = toTree(flat)
}

function reset() {
  Object.assign(form, {
    id: null, parentId: null, name: '', title: '', icon: '', path: '',
    component: '', type: 2, permission: '', api: '', method: '', sort: 0, status: 1
  })
}

function openDialog(row) {
  if (row) Object.assign(form, { ...row, parentId: row.parentId || null })
  else reset()
  visible.value = true
}

async function submit() {
  const payload = { ...form, parentId: form.parentId || 0 }
  if (form.id) await updateMenu(form.id, payload)
  else await createMenu(payload)
  ElMessage.success('保存成功')
  visible.value = false
  load()
}

async function remove(row) {
  await ElMessageBox.confirm('确认删除该菜单？', '提示', { type: 'warning' })
  await deleteMenu(row.id)
  ElMessage.success('已删除')
  load()
}

onMounted(load)
</script>

<style scoped>
.toolbar { margin-bottom: 12px; }
</style>
