<script lang="ts" setup>
defineProps<{
  min: number
  max: number
  step: number
}>()
const model = defineModel<number>({
  required: true,
  default: 0,
})
</script>

<template>
  <div class="fbh g24">
    <div class="pr w100p slider">
      <ElSlider
        v-model="model"
        :min="min"
        :max="max"
        :step="step"
        style="box-sizing: border-box"
      />
      <p class="pa label2 tick" style="right: 0">{{ max }}</p>
      <p class="pa label2 tick">{{ min }}</p>
    </div>
    <div class="number-input fbh fbac br8 px16 py6 g16 lh1">
      <p class="f14 fb1 text">{{ model.toLocaleString() }}</p>
      <p
        class="f20 label1 hand"
        :class="{disabled: model <= min}"
        @click="model -= step"
      >
        -
      </p>
      <p
        class="f20 label1 hand"
        :class="{disabled: model >= max}"
        @click="model += step"
      >
        +
      </p>
    </div>
  </div>
</template>

<style lang="less" scoped>
.number-input {
  box-sizing: border-box;
  border: 1px solid var(--tertiary-border-color);
  .text {
    width: 60px;
  }
}
.disabled {
  pointer-events: none;
  opacity: 0.4;
}
.slider {
  .tick {
    opacity: 0;
    bottom: -8px;
  }
  &:active {
    .tick {
      opacity: 1;
    }
  }
}
</style>
