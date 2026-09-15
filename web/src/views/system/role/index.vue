<template>
  <div>
    <div class="toolbar">
      <el-button type="primary" @click="openDialog()">新增角色</el-button>
    </div>
    <el-table :data="list" border stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="keyword" label="角色标识" />
      <el-table-column prop="description" label="描述" />
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button size="small" @click="openAssign(row)">分配菜单</el-button>
          <el-button size="small" type="danger" @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="visible" :title="form.id ? '编辑角色' : '新增角色'" width="420px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="标识"><el-input v-model="form.keyword" :disabled="!!form.id" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" @click="submit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="assignVisible" title="分配菜单" width="420px">
      <el-tree
        ref="treeRef"
        :data="menuTree"
        :props="{ label: 'title', children: 'children' }"
        node-key="id"
        show-checkbox
        default-expand-all
      />
      <template #footer>
        <el-button @click="assignVisible = false">取消</el-button>
        <el-button type="primary" @click="saveAssign">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listRoles, createRole, deleteRole, listMenus, assignMenus } from '@/api'

const list = ref([])
const menuTree = ref([])
const visible = ref(false)
const assignVisible = ref(false)
const treeRef = ref(null)
const form = reactive({ id: null, name: '', keyword: '', description: '' })
const currentRole = ref(null)

async function load() {
  const data = await listRoles()
  list.value = data.list || []
}
async function loadMenus() {
  const flat = await listMenus()
  menuTree.value = buildTree(flat)
}

function buildTree(flat) {
  const map = {}
  flat.forEach((m) => (map[m.id] = { ...m, children: [] }))
  const tree = []
  flat.forEach((m) => {
    if (m.parentId && map[m.parentId]) map[m.parentId].children.push(map[m.id])
    else tree.push(map[m.id])
  })
  return tree
}

function openDialog(row) {
  if (row) Object.assign(form, { id: row.id, name: row.name, keyword: row.keyword, description: row.description })
  else Object.assign(form, { id: null, name: '', keyword: '', description: '' })
  visible.value = true
}

async function submit() {
  await createRole({ ...form })
  ElMessage.success('保存成功')
  visible.value = false
  load()
}

async function remove(row) {
  await ElMessageBox.confirm('确认删除该角色？', '提示', { type: 'warning' })
  await deleteRole(row.id)
  ElMessage.success('已删除')
  load()
}

function openAssign(row) {
  currentRole.value = row
  assignVisible.value = true
  const ids = (row.menus || []).map((m) => m.id)
  setTimeout(() => treeRef.value?.setCheckedKeys(ids), 50)
}

async function saveAssign() {
  const menuIds = treeRef.value?.getCheckedKeys() || []
  await assignMenus({ roleId: currentRole.value.id, menuIds })
  ElMessage.success('分配成功')
  assignVisible.value = false
  load()
}

onMounted(() => {
  load()
  loadMenus()
})
</script>

<style scoped>
.toolbar { margin-bottom: 12px; }
</style>
