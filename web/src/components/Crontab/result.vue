<template>
  <div class="result">
    <p class="title">最近 5 次运行时间</p>
    <ul class="result-scroll">
      <li v-for="item in resultList" :key="item">{{ item }}</li>
    </ul>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  ex: { type: String, default: '' }
})

const resultList = ref([])

// Turn a single cron field into the set of allowed integer values.
function parseField(field, min, max) {
  const set = new Set()
  if (field === '*' || field === '?') {
    for (let i = min; i <= max; i++) set.add(i)
    return set
  }
  for (const part of field.split(',')) {
    if (part === '') continue
    let step = 1
    let base = part
    if (base.includes('/')) {
      const sp = base.split('/')
      step = parseInt(sp[1], 10) || 1
      base = sp[0]
    }
    if (base.includes('-')) {
      const bounds = base.split('-').map((n) => parseInt(n, 10))
      const lo = Math.min(bounds[0], bounds[1])
      const hi = Math.max(bounds[0], bounds[1])
      for (let i = lo; i <= hi; i += step) set.add(i)
    } else if (base === '*') {
      for (let i = min; i <= max; i += step) set.add(i)
    } else {
      const n = parseInt(base, 10)
      if (!isNaN(n)) set.add(n)
    }
  }
  return set
}

function pad(n) {
  return n < 10 ? '0' + n : '' + n
}

function format(d) {
  return (
    d.getFullYear() + '-' + pad(d.getMonth() + 1) + '-' + pad(d.getDate()) +
    ' ' + pad(d.getHours()) + ':' + pad(d.getMinutes()) + ':' + pad(d.getSeconds())
  )
}

// Compute up to 5 next run times, stepping minute by minute (matches backend Match).
function computeNext(expr) {
  const arr = expr.trim().split(/\s+/)
  if (arr.length !== 5) return ['表达式格式不正确（需 5 段：分 时 日 月 周）']

  const sets = {
    min: parseField(arr[0], 0, 59),
    hour: parseField(arr[1], 0, 23),
    day: parseField(arr[2], 1, 31),
    month: parseField(arr[3], 1, 12),
    week: parseField(arr[4], 0, 6)
  }

  const res = []
  const cur = new Date()
  cur.setSeconds(0, 0)
  cur.setMinutes(cur.getMinutes() + 1)
  const limit = new Date(cur)
  limit.setFullYear(limit.getFullYear() + 5)

  let guard = 0
  while (cur <= limit && res.length < 5) {
    if (
      sets.min.has(cur.getMinutes()) &&
      sets.hour.has(cur.getHours()) &&
      sets.day.has(cur.getDate()) &&
      sets.month.has(cur.getMonth() + 1) &&
      sets.week.has(cur.getDay())
    ) {
      res.push(format(cur))
    }
    cur.setTime(cur.getTime() + 60000)
    if (++guard > 3000000) break
  }
  return res.length ? res : ['最近 5 年内无符合条件的结果']
}

watch(
  () => props.ex,
  () => {
    resultList.value = computeNext(props.ex)
  },
  { immediate: true }
)
</script>

<style scoped>
.result {
  box-sizing: border-box;
  margin: 8px auto;
  padding: 14px 10px 10px;
  border: 1px solid #dcdfe6;
  position: relative;
}
.result .title {
  position: absolute;
  top: -13px;
  left: 50%;
  width: 140px;
  font-size: 13px;
  margin-left: -70px;
  text-align: center;
  line-height: 24px;
  background: #fff;
}
.result-scroll {
  list-style: none;
  margin: 0;
  padding: 0;
  font-size: 12px;
  line-height: 24px;
  height: 10em;
  overflow-y: auto;
  font-family: arial;
}
</style>
