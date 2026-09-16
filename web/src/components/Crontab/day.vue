<template>
  <el-form>
    <el-form-item>
      <el-radio v-model="radioValue" :value="1">日，允许的通配符[, - * /]</el-radio>
    </el-form-item>
    <el-form-item>
      <el-radio v-model="radioValue" :value="2">不指定</el-radio>
    </el-form-item>
    <el-form-item>
      <el-radio v-model="radioValue" :value="3">
        周期从
        <el-input-number v-model="cycle01" :min="1" :max="30" /> -
        <el-input-number v-model="cycle02" :min="cycle01 + 1" :max="31" /> 日
      </el-radio>
    </el-form-item>
    <el-form-item>
      <el-radio v-model="radioValue" :value="4">
        从
        <el-input-number v-model="average01" :min="1" :max="30" /> 号开始，每
        <el-input-number v-model="average02" :min="1" :max="31 - average01" /> 日执行一次
      </el-radio>
    </el-form-item>
    <el-form-item>
      <el-radio v-model="radioValue" :value="5">
        指定
        <el-select clearable v-model="checkboxList" placeholder="可多选" multiple :multiple-limit="10">
          <el-option v-for="item in 31" :key="item" :label="item" :value="item" />
        </el-select>
      </el-radio>
    </el-form-item>
  </el-form>
</template>

<script setup>
import { computed, ref, watch } from 'vue'

const emit = defineEmits(['update'])
const props = defineProps({
  cron: { type: Object, default: () => ({ min: '*', hour: '*', day: '*', month: '*', week: '?' }) },
  check: { type: Function, default: () => {} }
})

const radioValue = ref(1)
const cycle01 = ref(1)
const cycle02 = ref(2)
const average01 = ref(1)
const average02 = ref(1)
const checkboxList = ref([])
const checkCopy = ref([1])

const cycleTotal = computed(() => {
  cycle01.value = props.check(cycle01.value, 1, 30)
  cycle02.value = props.check(cycle02.value, cycle01.value + 1, 31)
  return cycle01.value + '-' + cycle02.value
})
const averageTotal = computed(() => {
  average01.value = props.check(average01.value, 1, 30)
  average02.value = props.check(average02.value, 1, 31 - average01.value)
  return average01.value + '/' + average02.value
})
const checkboxString = computed(() => checkboxList.value.join(','))

watch(() => props.cron.day, (value) => changeRadioValue(value))
watch([radioValue, cycleTotal, averageTotal, checkboxString], () => onRadioChange())

function changeRadioValue(value) {
  if (value === '*') radioValue.value = 1
  else if (value === '?') radioValue.value = 2
  else if (value.indexOf('-') > -1) {
    const a = value.split('-')
    cycle01.value = Number(a[0])
    cycle02.value = Number(a[1])
    radioValue.value = 3
  } else if (value.indexOf('/') > -1) {
    const a = value.split('/')
    average01.value = Number(a[0])
    average02.value = Number(a[1])
    radioValue.value = 4
  } else {
    checkboxList.value = [...new Set(value.split(',').map((i) => Number(i)))]
    radioValue.value = 5
  }
}

function onRadioChange() {
  // 日/周互斥：设置一方时，另一方置为“不指定”
  if (radioValue.value === 2 && props.cron.week === '?') emit('update', 'week', '*')
  if (radioValue.value !== 2 && props.cron.week !== '?') emit('update', 'week', '?')

  switch (radioValue.value) {
    case 1: emit('update', 'day', '*'); break
    case 2: emit('update', 'day', '?'); break
    case 3: emit('update', 'day', cycleTotal.value); break
    case 4: emit('update', 'day', averageTotal.value); break
    case 5:
      if (!checkboxList.value.length) checkboxList.value.push(checkCopy.value[0])
      else checkCopy.value = checkboxList.value
      emit('update', 'day', checkboxString.value)
      break
  }
}
</script>
