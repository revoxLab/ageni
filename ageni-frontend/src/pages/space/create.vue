<script lang="ts" setup>
import {AgentShape, updateAgentMutation, useAgentTabsQuery} from '@/apis'
import {CustomButton, Select, TextInput} from '@/components'
import {Modal} from '@/pages/modal'
import {computed, ref} from 'vue'

const {botId, data} = defineProps<{
  data: AgentShape
  botId: number
}>()
const emit = defineEmits<{
  confirm: []
  close: []
}>()
const name = ref(data.name)
const description = ref(data.desc)
const type = ref<SelectOption<string>>()
const {data: tabs} = useAgentTabsQuery()
const typeOptions = computed(() => {
  const options = (tabs.value?.data ?? [])
    .filter((tab) => tab !== 'Most used')
    .map((key) => ({
      key,
      label: key,
      value: key,
    }))
  type.value = options.find(({key}) => key === data.tab) ?? options[1]
  return options
})

async function onConfirm() {
  const res = await updateAgentMutation({
    params: {
      bot_id: botId,
      description: description.value,
      name: name.value,
      type: type.value!.value,
    },
  })
  if (res?.code === 0) {
    emit('confirm')
  }
}
</script>

<template>
  <Modal open @close="emit('close')">
    <div class="modal-root p20">
      <p class="f18 bold pb16 center">
        {{ data.status === 2 ? 'Edit the Agent' : 'Publish an Agent' }}
      </p>
      <div class="fbv g24">
        <div class="fbv g16 fbas">
          <p class="f14 bold label">
            Agent name <span class="require">*</span>
          </p>
          <TextInput
            v-model="name"
            placeholder="Enter a name for the agent"
            :max-length="30"
          />
        </div>
        <div class="fbv g16 fbas">
          <p class="f14 bold label">
            Agent type <span class="require">*</span>
          </p>
          <Select
            class="w100p"
            :selected="type"
            :options="typeOptions"
            @on-selected-change="(value) => (type = value)"
          />
        </div>
        <div class="fbv g16 fbas">
          <p class="f14 bold label">
            Agent description <span class="require">*</span>
          </p>
          <TextInput
            v-model="description"
            placeholder="Introducing the agent functions and it will be displayed to the agent users"
            :max-length="150"
            :min-rows="4"
          />
        </div>
        <div class="fbh fbje g16 f16">
          <CustomButton
            :width="109"
            :height="40"
            default-color="var(--primary-label-color)"
            default-stroke="var(--tertiary-border-color)"
            default-fill="transparent"
            @click="emit('close')"
          >
            Cancel
          </CustomButton>
          <CustomButton
            :width="109"
            :height="40"
            :disabled="!name || !type || !description"
            @click="onConfirm"
          >
            Confirm
          </CustomButton>
        </div>
      </div>
    </div>
  </Modal>
</template>

<style lang="less" scoped>
.modal-root {
  width: 540px;
  box-sizing: border-box;
}
.require {
  color: #e35461;
}
.active {
  color: var(--theme-color);
}
</style>
