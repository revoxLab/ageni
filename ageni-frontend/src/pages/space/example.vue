<script lang="ts" setup>
import {PluginMethodShape} from '@/apis'
import {computed, ref} from 'vue'
import Modal from '../modal/modal.vue'

defineEmits<{close: []}>()

const props = defineProps<{data: PluginMethodShape}>()
const methodType = ref<'request' | 'response'>('request')
const activeSchema = computed(() =>
  JSON.parse(
    (methodType.value === 'request'
      ? props.data.input_schema
      : props.data.output_schema) ?? 'null'
  )
)
</script>

<template>
  <Modal open @close="$emit('close')">
    <div class="modal-root p24">
      <p class="center f18 bold">Functions</p>
      <div class="fbh f13 pb10">
        <div
          class="tab hand br4"
          :class="{active: methodType === 'request'}"
          @click="methodType = 'request'"
        >
          Request
        </div>
        <div
          class="tab hand br4"
          :class="{active: methodType === 'response'}"
          @click="methodType = 'response'"
        >
          Response
        </div>
      </div>
      <div class="table fbh br8 f14 label1">
        <div class="fb1 fbv">
          <div class="fbh header p16">
            <div class="fb1">Parameter</div>
            <div class="fb1">Description</div>
          </div>
          <template v-for="(item, i) in activeSchema?.parameters">
            <div v-if="i !== 0" class="divider mx16"></div>
            <div class="fbh">
              <div class="fb1 p16">{{ item.name }}</div>
              <div class="fb1 p16">{{ item.description }}</div>
            </div>
          </template>
        </div>
        <div class="fb1 example oa">
          <div class="p16 header w100p">Example</div>
          <pre class="p16 oa" style="text-wrap: wrap">{{
            methodType === 'request' ? data.input_example : data.output_example
          }}</pre>
        </div>
      </div>
    </div>
  </Modal>
</template>

<style lang="less" scoped>
.modal-root {
  width: 1000px;
  box-sizing: border-box;
}
.tab {
  padding: 5px 10px;
  color: var(--secondary-label-color);
  &.active {
    background: var(--tertiary-background-color);
    color: var(--primary-text-color);
  }
}
.table {
  border: 1px solid var(--secondary-border-color);
}
.header {
  background: var(--secondary-background-color);
  box-sizing: border-box;
}
.example {
  border-left: 1px solid var(--secondary-border-color);
}
</style>
