<script lang="ts" setup>
const {data, activeColor} = defineProps<{
  data: {title: string; key: string}[]
  activeColor: string
  activeKey: string
}>()
const emit = defineEmits<{
  onActiveKeyChange: [key: any]
}>()
</script>

<template>
  <div v-bind="$attrs" class="tab-bar">
    <div
      v-for="{title, key} in data"
      @click="emit('onActiveKeyChange', key)"
      :class="['tab-item', {active: activeKey === key}]"
    >
      <p>{{ title }}</p>
      <div class="underline active" v-show="activeKey === key"></div>
    </div>
  </div>
  <slot name="default"></slot>
  <template v-for="{key} in data" v-show="activeKey === key">
    <slot v-if="activeKey === key" :name="key"></slot>
  </template>
</template>

<style lang="less" scoped>
.tab-bar {
  display: flex;
  align-items: center;
  line-height: 1.5;
  gap: 24px;
}
.tab-item {
  font-size: 16px;
  text-wrap: nowrap;
  color: var(--secondary-label-color);
  line-height: 2.2;
  font-weight: 700;
  position: relative;
  cursor: pointer;
  &.active {
    color: v-bind(activeColor);
  }
}
.underline {
  background: v-bind(activeColor);
  position: absolute;
  margin-left: 30%;
  height: 2px;
  width: 40%;
  bottom: 0;
}
</style>
