<template>
  <div>
    <div class="toolbar">
      <el-button type="primary" @click="openCreate">创建密钥</el-button>
      <span class="tip">密钥供 skill / 脚本调用接口，仅展示一次，请妥善保存</span>
    </div>

    <el-table :data="list" border stripe>
      <el-table-column prop="name" label="名称" min-width="140" />
      <el-table-column label="密钥" min-width="200">
        <template #default="{ row }">
          <code class="masked">{{ row.maskedKey }}</code>
        </template>
      </el-table-column>
      <el-table-column prop="scope" label="权限范围" width="100">
        <template #default="{ row }">
          <el-tag :type="row.scope === 'readonly' ? 'info' : 'success'">
            {{ row.scope === 'readonly' ? '只读' : '全部' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="90">
        <template #default="{ row }">
          <el-tag v-if="row.expired" type="danger">已过期</el-tag>
          <el-tag v-else :type="row.status === 1 ? 'success' : 'danger'">
            {{ row.status === 1 ? '启用' : '已禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="expiresAt" label="过期时间" width="130">
        <template #default="{ row }">{{ row.expiresAt ? fmt(row.expiresAt) : '永久' }}</template>
      </el-table-column>
      <el-table-column prop="lastUsedAt" label="最近使用" width="170">
        <template #default="{ row }">{{ fmt(row.lastUsedAt) }}</template>
      </el-table-column>
      <el-table-column prop="createdAt" label="创建时间" width="170">
        <template #default="{ row }">{{ fmt(row.createdAt) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="100">
        <template #default="{ row }">
          <el-button size="small" type="danger" @click="revoke(row)">吊销</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- create dialog -->
    <el-dialog v-model="createVisible" title="创建 API 密钥" width="460px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="form.name" placeholder="便于识别，如 my-skill" maxlength="64" />
        </el-form-item>
        <el-form-item label="权限范围">
          <el-select v-model="form.scope" placeholder="选择权限范围">
            <el-option label="全部权限（增删改查）" value="all" />
            <el-option label="只读（仅 GET 查询）" value="readonly" />
          </el-select>
        </el-form-item>
        <el-form-item label="过期时间">
          <el-date-picker
            v-model="form.expiresAt"
            type="date"
            placeholder="不填则永不过期"
            value-format="YYYY-MM-DD"
            :disabled-date="disabledDate"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">取消</el-button>
        <el-button type="primary" @click="submitCreate">生成</el-button>
      </template>
    </el-dialog>

    <!-- show raw key once -->
    <el-dialog v-model="showVisible" title="密钥已生成" width="520px" :close-on-click-modal="false">
      <el-alert type="warning" :closable="false" title="请立即复制保存，密钥仅展示这一次，关闭后无法再次查看" />
      <div class="key-box">
        <code>{{ createdKey }}</code>
        <el-button size="small" @click="copy(createdKey)">复制</el-button>
      </div>
      <template #footer>
        <el-button type="primary" @click="showVisible = false">我已保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listApiKeys, createApiKey, deleteApiKey } from '@/api'

const list = ref([])
const createVisible = ref(false)
const showVisible = ref(false)
const createdKey = ref('')
const form = reactive({ name: '', scope: 'all', expiresAt: '' })

async function load() {
  const data = await listApiKeys()
  list.value = data.list || []
}

function openCreate() {
  form.name = ''
  form.scope = 'all'
  form.expiresAt = ''
  createVisible.value = true
}

async function submitCreate() {
  const payload = { name: form.name, scope: form.scope }
  if (form.expiresAt) payload.expiresAt = form.expiresAt
  const data = await createApiKey(payload)
  createdKey.value = data.key || data.secret
  createVisible.value = false
  showVisible.value = true
  load()
}

async function revoke(row) {
  await ElMessageBox.confirm(`确认吊销密钥「${row.name}」？吊销后该密钥立即失效。`, '提示', { type: 'warning' })
  await deleteApiKey(row.id)
  ElMessage.success('已吊销')
  load()
}

function fmt(t) {
  if (!t) return '—'
  const d = new Date(t)
  if (isNaN(d.getTime())) return '—'
  const p = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}`
}

function disabledDate(date) {
  return date.getTime() < Date.now() - 24 * 60 * 60 * 1000
}

async function copy(text) {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制')
  } catch {
    ElMessage.warning('复制失败，请手动复制')
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; align-items: center; gap: 12px; margin-bottom: 12px; }
.tip { color: #999; font-size: 13px; }
.masked { background: #f5f5f5; padding: 2px 6px; border-radius: 4px; letter-spacing: 1px; }
.key-box { margin: 14px 0; display: flex; gap: 8px; align-items: center; }
.key-box code { flex: 1; background: #f5f5f5; padding: 10px; border-radius: 4px; word-break: break-all; }
</style>
