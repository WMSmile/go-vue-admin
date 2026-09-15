<template>
  <div class="profile">
    <el-row :gutter="16">
      <el-col :span="8">
        <el-card>
          <div class="avatar-wrap">
            <el-avatar :size="80">{{ initial }}</el-avatar>
            <h3>{{ user.username }}</h3>
            <p class="role">{{ roleText }}</p>
          </div>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="昵称">{{ user.nickname || '-' }}</el-descriptions-item>
            <el-descriptions-item label="邮箱">{{ user.email || '-' }}</el-descriptions-item>
            <el-descriptions-item label="手机">{{ user.phone || '-' }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="user.status === 1 ? 'success' : 'danger'">
                {{ user.status === 1 ? '启用' : '禁用' }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
          <el-button class="mt" type="primary" @click="openEdit">编辑资料</el-button>
        </el-card>
      </el-col>
      <el-col :span="16">
        <el-card>
          <h3>安全设置</h3>
          <p class="tip">定期修改密码有助于保护账户安全。</p>
          <el-button @click="openPwd">修改密码</el-button>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="editVisible" title="编辑资料" width="460px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="昵称"><el-input v-model="editForm.nickname" /></el-form-item>
        <el-form-item label="邮箱"><el-input v-model="editForm.email" /></el-form-item>
        <el-form-item label="手机"><el-input v-model="editForm.phone" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" @click="saveEdit">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="pwdVisible" title="修改密码" width="460px">
      <el-form :model="pwdForm" label-width="90px">
        <el-form-item label="新密码">
          <el-input v-model="pwdForm.password" type="password" show-password placeholder="留空则不修改" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="pwdVisible = false">取消</el-button>
        <el-button type="primary" @click="savePwd">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'
import { updateUser } from '@/api'

const store = useUserStore()
const user = computed(() => store.userInfo || {})
const initial = computed(() => (user.value.username || '?').charAt(0).toUpperCase())
const roleText = computed(() => (user.value.roles || []).map((r) => r.name || r.keyword).join('、') || '-')

const editVisible = ref(false)
const pwdVisible = ref(false)
const editForm = reactive({ nickname: '', email: '', phone: '' })
const pwdForm = reactive({ password: '' })

function openEdit() {
  editForm.nickname = user.value.nickname || ''
  editForm.email = user.value.email || ''
  editForm.phone = user.value.phone || ''
  editVisible.value = true
}

async function saveEdit() {
  await updateUser(user.value.id, {
    username: user.value.username,
    nickname: editForm.nickname,
    email: editForm.email,
    phone: editForm.phone
  })
  store.userInfo = { ...store.userInfo, ...editForm }
  localStorage.setItem('userInfo', JSON.stringify(store.userInfo))
  ElMessage.success('保存成功')
  editVisible.value = false
}

function openPwd() {
  pwdForm.password = ''
  pwdVisible.value = true
}

async function savePwd() {
  if (!pwdForm.password) {
    ElMessage.warning('请输入新密码')
    return
  }
  await updateUser(user.value.id, { username: user.value.username, password: pwdForm.password })
  ElMessage.success('密码已更新')
  pwdVisible.value = false
}
</script>

<style scoped>
.avatar-wrap { text-align: center; padding: 12px 0; }
.avatar-wrap h3 { margin: 10px 0 4px; }
.role { color: #888; margin: 0; }
.mt { margin-top: 12px; }
.tip { color: #999; font-size: 13px; }
</style>
