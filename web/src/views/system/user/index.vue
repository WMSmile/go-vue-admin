<template>
  <div>
    <div class="toolbar">
      <el-button type="primary" v-permission="'user:add'" @click="openDialog()">新增用户</el-button>
    </div>
    <el-table :data="list" border stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="username" label="用户名" />
      <el-table-column prop="nickname" label="昵称" />
      <el-table-column prop="email" label="邮箱" />
      <el-table-column label="状态" width="90">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180">
        <template #default="{ row }">
          <el-button size="small" v-permission="'user:edit'" @click="openDialog(row)">编辑</el-button>
          <el-button size="small" type="danger" v-permission="'user:del'" @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="visible" :title="form.id ? '编辑用户' : '新增用户'" width="460px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="用户名">
          <el-input v-model="form.username" :disabled="!!form.id" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" placeholder="留空则不修改" />
        </el-form-item>
        <el-form-item label="昵称"><el-input v-model="form.nickname" /></el-form-item>
        <el-form-item label="邮箱"><el-input v-model="form.email" /></el-form-item>
        <el-form-item label="角色">
          <el-select v-model="form.roleIds" multiple style="width: 100%">
            <el-option v-for="r in roles" :key="r.id" :label="r.name" :value="r.id" />
          </el-select>
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
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listUsers, createUser, updateUser, deleteUser, listRoles } from '@/api'

const list = ref([])
const roles = ref([])
const visible = ref(false)
const form = reactive({ id: null, username: '', password: '', nickname: '', email: '', roleIds: [] })

async function load() {
  try {
    const data = await listUsers()
    list.value = data.list || []
  } catch (e) {}
}
async function loadRoles() {
  try {
    const data = await listRoles()
    roles.value = data.list || []
  } catch (e) {}
}

function openDialog(row) {
  if (row) {
    form.id = row.id
    form.username = row.username
    form.password = ''
    form.nickname = row.nickname
    form.email = row.email
    form.roleIds = (row.roles || []).map((r) => r.id)
  } else {
    Object.assign(form, { id: null, username: '', password: '', nickname: '', email: '', roleIds: [] })
  }
  visible.value = true
}

async function submit() {
  const payload = { ...form }
  if (form.id) {
    await updateUser(form.id, payload)
  } else {
    await createUser(payload)
  }
  ElMessage.success('保存成功')
  visible.value = false
  load()
}

async function remove(row) {
  await ElMessageBox.confirm('确认删除该用户？', '提示', { type: 'warning' })
  await deleteUser(row.id)
  ElMessage.success('已删除')
  load()
}

onMounted(() => {
  load()
  loadRoles()
})
</script>

<style scoped>
.toolbar { margin-bottom: 12px; }
</style>
