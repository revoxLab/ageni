<script lang="ts" setup>
import {
  ModelSetting,
  pluginListQuery,
  PluginShape,
  publishAgentMutation,
  saveAgentDraftMutation,
  useAgentDetailQuery,
  useAgentDraftQuery,
} from '@/apis'
import {
  CustomButton,
  Icon,
  LoadingDot,
  showToast,
  TextInput,
} from '@/components'
import {format} from 'date-fns'
import {ElInput, ElTooltip} from 'element-plus'
import {debounce} from 'lodash-es'
import {computed, ref, watch} from 'vue'
import CreateAgentModal from './create.vue'
import Message from './message.vue'
import Model from './model.vue'
import AddPluginModal from './plugin.vue'
import Section from './section.vue'
import {jump} from '@/utils'

defineEmits<{close: []}>()

const {id} = defineProps<{id: number}>()
const {data: detailData, refetch} = useAgentDetailQuery({
  params: {bot_id: id},
})
const {data: draftData} = useAgentDraftQuery({
  params: {bot_id: id},
})
const pluginModalOpen = ref(false)
const createModalOpen = ref(false)
const welcomeGenerating = ref(false)
const detail = computed(() => detailData.value?.data)
const draft = computed(() => draftData.value?.data.draft)
const activeTab = ref<'settings' | 'statistics'>('settings')
const modelSetting = ref<ModelSetting>({
  model: 'gpt-4o',
  max_length: 4095,
  temperature: 1,
  rounds: 10,
  top_p: 1,
})
const questionRef = ref<HTMLDivElement>()
const selectedPlugin = ref<PluginShape>()
const plugins = ref<PluginShape[]>([])
const guideInfo = ref<string[]>([])
const updateTime = ref(Date.now())
const questionIndex = ref(0)
const welcome = ref('')
const prompt = ref('')
const snapshot = computed(() => ({
  bot_id: id,
  guide_info: guideInfo.value,
  welcome_msg: welcome.value,
  prompt: prompt.value,
  plugins: plugins.value.map(({id}) => id),
  model_settings: modelSetting.value,
}))

const saveDraftMutation = debounce(() => {
  saveAgentDraftMutation({params: snapshot.value})
  updateTime.value = Date.now()
}, 1000)

watch([guideInfo, welcome, prompt, plugins, modelSetting], saveDraftMutation, {
  deep: true,
})

watch(draft, async (value) => {
  updateTime.value = (value?.updated_at.seconds ?? Date.now() / 1000) * 1000
  guideInfo.value = value?.guide_info ?? []
  welcome.value = value?.welcome_msg ?? ''
  prompt.value = value?.prompt ?? ''
  if (value?.model_settings) {
    modelSetting.value = Object.assign(
      {model: 'gpt-4o', max_length: 0, temperature: 0, rounds: 0, top_p: 0},
      value?.model_settings
    )
  }
  if (value?.plugins?.length) {
    const res = await pluginListQuery({params: {ids: value.plugins}})
    plugins.value = res?.data.list ?? []
  }
})

function onNewQuestion(inputEvent: Event | KeyboardEvent) {
  let event = inputEvent as KeyboardEvent
  if (/^[a-zA-Z0-9\s`~!@#$%^&*()\-_=+\[\]{}|;:'",.<>/?\\]$/.test(event.key)) {
    guideInfo.value[questionIndex.value] =
      (guideInfo.value[questionIndex.value] ?? '') + event.key
  }
  requestAnimationFrame(() => {
    questionRef.value
      ?.querySelectorAll('textarea')
      [guideInfo.value.length - 1].focus()
  })
}

async function onPublish() {
  const res = await publishAgentMutation({
    params: snapshot.value,
  })
  if (res?.code === 0) {
    showToast('The agent is published', 'check')
  }
}

async function onUpdateAgent() {
  if (detail.value?.status === 2) {
    showToast('The agent is updated', 'check')
  } else {
    await onPublish()
  }
  createModalOpen.value = false
  refetch()
}
</script>

<template>
  <CreateAgentModal
    v-if="createModalOpen && detail"
    :data="detail"
    :bot-id="id"
    @close="createModalOpen = false"
    @confirm="onUpdateAgent"
  />
  <AddPluginModal
    v-if="pluginModalOpen"
    variant="multiple"
    :plugins="plugins.map(({id}) => id)"
    @close="pluginModalOpen = false"
    @add-plugin="(value) => plugins.push(value)"
    @remove-plugin="
      ({id}) =>
        plugins.splice(
          plugins.findIndex((value) => value.id === id),
          1
        )
    "
  />
  <AddPluginModal
    v-if="selectedPlugin"
    variant="single"
    :plugins="
      plugins.find(({id}) => id === selectedPlugin?.id)
        ? [selectedPlugin.id]
        : []
    "
    @close="selectedPlugin = undefined"
    @add-plugin="(value) => plugins.push(value)"
    @remove-plugin="
      ({id}) =>
        plugins.splice(
          plugins.findIndex((value) => value.id === id),
          1
        )
    "
  />
  <div v-if="!detail || !draft" class="modal wh100v">
    <ElSkeleton class="p20" style="box-sizing: border-box" animated />
  </div>
  <div v-else class="modal wh100v">
    <div class="fbh fbjsb fbac p20 header">
      <div class="fbh fbac g10">
        <Icon
          @click="$emit('close')"
          style="opacity: 0.4"
          name="arrow-left"
          :size="12"
        />
        <Icon :name="detail.image ?? 'agent'" :size="32" class="br6" />
        <div class="fbv fbjsb">
          <p class="f15 bold">
            {{ detail.name }}
            <Icon
              v-if="detail.status === 2"
              @click="createModalOpen = true"
              class="icon pl8"
              name="edit2"
              :size="12"
            />
          </p>
          <div
            class="br4 p2 pl0 fbh fbac g4 f12 lh1 label1"
            :class="{highlight: detail.status === 2}"
          >
            <Icon :name="detail.status === 2 ? 'check1' : 'time'" :size="10" />
            {{ detail.status === 2 ? 'Published' : 'Draft' }}
            <p class="pl4 label2">
              Last save:
              {{ format(updateTime, 'yyyy.MM.dd HH:mm') }}
            </p>
          </div>
        </div>
      </div>
      <div class="center g24 f18 bold hand">
        <p
          :class="{highlight2: activeTab === 'settings'}"
          @click="activeTab = 'settings'"
        >
          Settings
        </p>
      </div>
      <div class="fbh fbac g28">
        <div
          class="fbh fbac g4 hand"
          @click="jump('https://readon-me.gitbook.io/revox-studio-guild')"
        >
          <p class="f16" style="color: var(--theme-color)">Quick Guide</p>
          <Icon name="query2" :size="16" />
        </div>
        <CustomButton
          :width="96"
          :height="36"
          class="f14 fb1 fbh fbje"
          :disabled="!prompt"
          @click="detail.status === 2 ? onPublish() : (createModalOpen = true)"
        >
          + Publish
        </CustomButton>
      </div>
    </div>
    <div class="content">
      <div class="card oa">
        <div class="fbh fbjsb fbac px24 py14">
          <p class="f18 bold">Editor</p>
          <Model v-model="modelSetting" />
        </div>
        <div class="divider mx24"></div>
        <div class="fbh">
          <div class="fb1 py12 px24 fbv g12">
            <p class="f15 bold label1">Prompt</p>
            <ElInput
              type="textarea"
              placeholder="Fill in the character setting, functions, and workflow of the agent using natural language."
              :autosize="{minRows: 5, maxRows: Infinity}"
              v-model="prompt"
              resize="none"
            />
          </div>
          <div class="divider-v"></div>
          <div class="fb1 py12 px24 fbv g16">
            <p class="f15 bold label1">Skill</p>
            <Section title="Plugins" :initialOpen="true">
              <template #action>
                <ElTooltip content="Add plugin">
                  <Icon name="add" :size="12" @click="pluginModalOpen = true" />
                </ElTooltip>
              </template>
              <template #default>
                <p v-if="!plugins.length" class="f13 label2">
                  Plugins can enable the Agent to call external APIs, such as
                  searching for information, browsing web pages, generating
                  images, etc., expanding the capabilities and usage scenarios
                  of the Agent
                </p>
                <div v-else class="fbv g8">
                  <div
                    v-for="(item, i) in plugins"
                    class="fbh fbac g8 plugin py10 px12 br8"
                  >
                    <Icon :name="item.image" :size="36" class="br6" />
                    <div class="fbv fb1">
                      <p class="f14 bold">{{ item.name }}</p>
                      <p class="f12 label2 omit1">{{ item.desc }}</p>
                    </div>
                    <Icon
                      @click="selectedPlugin = item"
                      name="query"
                      :size="14"
                    />
                    <ElTooltip content="Remove">
                      <Icon
                        @click="plugins.splice(i, 1)"
                        name="delete"
                        :size="14"
                      />
                    </ElTooltip>
                  </div>
                </div>
              </template>
            </Section>
            <p class="f15 bold label1">Personalized Settings</p>
            <Section title="Opening remarks" :initialOpen="true">
              <template #default>
                <div class="fbv g10">
                  <p class="f12 label1">Opening remarks copywriting</p>
                  <div class="pr">
                    <TextInput
                      v-model="welcome"
                      placeholder="Opening remarks copywriting"
                      :max-length="150"
                      :min-rows="6"
                    />
                    <div
                      v-if="welcomeGenerating"
                      class="center wh100p pa fbv generate"
                    >
                      <div class="overlay pa wh100p"></div>
                      <LoadingDot :size="30" color="var(--theme-color)" />
                      <p class="f13">Generating</p>
                    </div>
                  </div>
                </div>
                <div class="fbv g10" ref="questionRef">
                  <p class="f12 label1">Guiding questions</p>
                  <div v-for="(_, i) in guideInfo" class="pr">
                    <ElInput
                      type="textarea"
                      :autosize="{minRows: 1, maxRows: Infinity}"
                      placeholder="Enter the guiding question"
                      v-model="guideInfo[i]"
                      resize="none"
                    />
                    <ElTooltip content="Delete question">
                      <Icon
                        :size="14"
                        name="delete"
                        class="delete p5 pa"
                        @click="guideInfo.splice(i, 1)"
                      />
                    </ElTooltip>
                  </div>
                  <ElInput
                    maxlength="0"
                    type="textarea"
                    :autosize="{minRows: 1, maxRows: Infinity}"
                    placeholder="Enter the guiding question"
                    @focus="questionIndex = guideInfo.length"
                    @keydown="onNewQuestion"
                    resize="none"
                  />
                </div>
              </template>
            </Section>
            <p class="f15 bold label1">Else</p>
            <Section title="Trigger" :initialOpen="false">
              <template #action>
                <Icon name="add" :size="12" />
              </template>
            </Section>
            <Section title="Workflow" :initialOpen="false">
              <template #action>
                <Icon name="add" :size="12" />
              </template>
            </Section>
            <Section title="Picture Flow" :initialOpen="false">
              <template #action>
                <Icon name="add" :size="12" />
              </template>
            </Section>
            <Section title="Problem Suggestions" :initialOpen="false">
              <template #action>
                <Icon name="add" :size="12" />
              </template>
            </Section>
            <Section title="Quick Commands" :initialOpen="false">
              <template #action>
                <Icon name="add" :size="12" />
              </template>
            </Section>
            <Section title="Background Picture" :initialOpen="false">
              <template #action>
                <Icon name="add" :size="12" />
              </template>
            </Section>
          </div>
        </div>
      </div>
      <div class="card">
        <p class="f18 bold px24 py23">Preview</p>
        <div class="divider mx24"></div>
        <Message
          :conversation-id="draft.debug_conversation_id"
          :bot="{
            id,
            name: detail.name,
            guide_info: guideInfo,
            image: detail.image ?? 'agent',
            welcome_msg: welcome,
          }"
        />
      </div>
    </div>
  </div>
</template>

<style lang="less" scoped>
.modal {
  background: var(--primary-background-color);
  position: absolute;
  z-index: 1;
  left: 0;
  top: 0;
}
.highlight {
  color: #00c087;
}
.highlight2 {
  color: #44dac8;
}
.header {
  height: 72px;
  box-sizing: border-box;
  .icon {
    opacity: 0.6;
    &:hover {
      opacity: 1;
    }
  }
}
.content {
  display: grid;
  margin-left: 1px;
  grid-template-columns: 930px 1fr;
  grid-template-rows: calc(100vh - 74px);
  background: var(--primary-border-color);
}
.card {
  border: 1px solid var(--tertiary-border-color);
  background: var(--primary-background-color);
  border-radius: 8px;
}
.select {
  width: 180px;
}
.delete {
  right: 11px;
  top: 6px;
}
.plugin {
  background: var(--tertiary-background-color);
}
.generate {
  top: 0;
  color: var(--theme-color);
  z-index: 1;
  .overlay {
    background: var(--primary-background-color);
    opacity: 0.8;
    z-index: -1;
  }
}
</style>
