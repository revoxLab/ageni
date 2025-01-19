<script lang="ts" setup>
import {ModelSetting} from '@/apis'
import {Icon, Select, SliderInput} from '@/components'
import {ElDropdown, ElTooltip} from 'element-plus'
import {ref} from 'vue'

const aiSelectOptions: SelectOption<string>[] = [
  {key: 1, label: 'gpt-4o', icon: 'gpt', value: 'gpt-4o'},
]
const data = defineModel<ModelSetting>({required: true})
const selectAI = ref(
  aiSelectOptions.find(({value}) => value === data.value.model) ??
    aiSelectOptions[0]
)
</script>

<template>
  <ElDropdown trigger="click">
    <div class="select-input fbh fbjsb fbac g6 hand">
      <Icon :name="selectAI.icon!" :size="24" />
      <p class="fb1">{{ selectAI.label }}</p>
      <Icon name="arrow-down" :size="12" />
    </div>
    <template #dropdown>
      <div class="py20 px24 model-root">
        <p class="f18 bold pb16 center">Model settings</p>
        <p class="label1 f14 bold pb14">Model</p>
        <Select
          @on-selected-change="
            (value) => ((selectAI = value), (data.model = value.value))
          "
          :selected="selectAI"
          :options="aiSelectOptions"
          class="select"
        />
        <div class="fbv g14 pt24">
          <p class="label1 f14 bold">Model parameters</p>
          <div class="fbh fbac f16 g8">
            Temperature
            <ElTooltip>
              <Icon name="query" :size="14" />
              <template #content>
                <p class="tooltip">
                  Increasing the temperature will make the model's output more
                  diverse and innovative. Conversely, reducing the temperature
                  will make the output content more in line with the instruction
                  requirements but reduce diversity. It is recommended not to
                  adjust it simultaneously with "Top P".
                </p>
              </template>
            </ElTooltip>
          </div>
          <SliderInput
            v-model="data.temperature"
            :min="0"
            :max="2"
            :step="0.01"
          />
          <div class="fbh fbac f16 g8">
            Top P
            <ElTooltip>
              <Icon name="query" :size="14" />
              <template #content>
                <p class="tooltip">
                  Top P is cumulative probability: When generating output, the
                  model will start to select from the words with the highest
                  probability until the total probability of these words
                  accumulates to reach the Top P value. This can limit the model
                  to only select these high-probability words, thereby
                  controlling the diversity of the output content. It is
                  recommended not to adjust it simultaneously with
                  "Temperature".
                </p>
              </template>
            </ElTooltip>
          </div>
          <SliderInput v-model="data.top_p" :min="0" :max="1" :step="0.01" />
        </div>
        <div class="fbv g14 pt24">
          <p class="label1 f14 bold">Input and output settings</p>
          <div class="fbh fbac f16 g8">
            The number of rounds carrying context
          </div>
          <SliderInput v-model="data.rounds" :min="3" :max="100" :step="1" />
          <div class="fbh fbac f16 g8">Maximum reply length</div>
          <SliderInput
            v-model="data.max_length"
            :min="1"
            :max="4095"
            :step="1"
          />
        </div>
      </div>
    </template>
  </ElDropdown>
</template>

<style lang="less" scoped>
.model-root {
  min-width: 480px;
}
.select-input {
  background: var(--tertiary-background-color);
  color: var(--primary-text-color);
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 16px;
  min-width: 180px;
}
.tooltip {
  max-width: 400px;
}
</style>
