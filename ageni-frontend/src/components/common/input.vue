<script lang="ts" setup>
import {ElInput} from 'element-plus'

const model = defineModel<string>()
const {minRows = 1} = defineProps<{
  minRows?: number
  maxLength?: number
  placeholder?: string
}>()
</script>

<template>
  <div class="pr w100p">
    <ElInput
      v-bind="$attrs"
      type="textarea"
      :maxlength="maxLength"
      :autosize="{minRows, maxRows: Infinity}"
      :placeholder="placeholder"
      v-model="model"
      resize="none"
    />
    <p v-if="maxLength" class="pa count f12 label2">
      {{ model?.length ?? 0 }}/{{ maxLength }}
    </p>
  </div>
</template>

<style lang="less" scoped>
.count {
  right: 16px;
  bottom: 8px;
}
</style>
