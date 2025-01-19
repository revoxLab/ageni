<script lang="ts" setup>
import {
  AgentShape,
  PluginShape,
  useAgentListQuery,
  usePluginListQuery,
} from '@/apis'
import {Icon} from '@/components'
import {jump, shortAddress} from '@/utils'
import {debounce} from 'lodash-es'
import {computed, ref, watch} from 'vue'

const search = ref('')
const focused = ref(false)
const agentList = ref<AgentShape[]>([])
const pluginList = ref<PluginShape[]>([])
const agentPage = ref(1)
const pluginPage = ref(1)
const {
  data: agentData,
  refetch: refetchAgent,
  loading: loadingAgent,
} = useAgentListQuery({
  params: {page: 1, page_size: 5, keywords: '', tab: ''},
  skip: true,
})
const {
  data: pluginData,
  refetch: refetchPlugin,
  loading: loadingPlugin,
} = usePluginListQuery({
  params: {page: 1, page_size: 5, keywords: '', tab: ''},
  skip: true,
})
const showResult = computed(
  () =>
    loadingPlugin.value ||
    loadingAgent.value ||
    pluginList.value.length ||
    agentList.value.length ||
    pluginData.value?.code === 0 ||
    agentData.value?.code === 0
)

watch(
  search,
  debounce(() => {
    agentData.value = undefined
    pluginData.value = undefined
    if (search.value) {
      refetchAgent({
        page: (agentPage.value = 1),
        keywords: search.value,
      })
      refetchPlugin({
        page: (pluginPage.value = 1),
        keywords: search.value,
      })
    }
  }, 500)
)

watch(agentData, () => {
  if (agentPage.value === 1) {
    agentList.value = agentData.value?.data.bots ?? []
  } else {
    agentList.value.push(...(agentData.value?.data.bots ?? []))
  }
})

watch(pluginData, () => {
  if (pluginPage.value === 1) {
    pluginList.value = pluginData.value?.data.list ?? []
  } else {
    pluginList.value.push(...(pluginData.value?.data.list ?? []))
  }
})
</script>

<template>
  <div class="pr">
    <input
      type="text"
      placeholder="Search"
      class="input br8 px16 py10 f14 bold"
      @focus="focused = true"
      v-model="search"
    />
    <Teleport v-if="focused" to="#modal">
      <div class="modal pa" @click.stop="focused = false"></div>
    </Teleport>
    <div v-if="focused && showResult" class="result py16 px20 pa oa fbv">
      <p class="f13 bold label1 pb12">Agent</p>
      <div class="fbv g20">
        <div
          v-for="item in agentList"
          @click="jump(`/agent?id=${item.id}`)"
          class="fbh fbac g10 hand"
        >
          <Icon :name="item.image" :size="32" />
          <div class="fbv g4">
            <div class="fbh fbac g12">
              <p class="f14 bold">{{ item.name }}</p>
              <div class="fbh fbac g4">
                <Icon :name="item.creator.head_pic" :size="14" class="br100" />
                <p class="label1 f13">@{{ shortAddress(item.creator.name) }}</p>
              </div>
              <div class="fbh fbac g4 f12 label2">
                <Icon name="user" style="opacity: 0.4" :size="12" />
                <p class="pr8">{{ item.users }}</p>
                <Icon name="comment" style="opacity: 0.4" :size="12" />
                <p>{{ item.conversations }}</p>
              </div>
            </div>
            <p class="f13 label1 omit1">{{ item.desc }}</p>
          </div>
        </div>
        <ElSkeleton v-if="loadingAgent" animated />
        <div
          v-if="agentData?.data.bots?.length === 5"
          class="fbh fbjc fbac g4 pb8 hand"
          @click="refetchAgent({page: ++agentPage, keywords: search})"
        >
          <p class="f14 label2">More</p>
          <Icon name="arrow-down" :size="10"></Icon>
        </div>
        <div v-else></div>
      </div>
      <div v-if="!loadingAgent && !agentList.length" class="fbv fbac g12 p30">
        <Icon name="empty1" :width="114" :height="92" />
        <p class="label1 f14">Not Found</p>
      </div>
      <p class="f13 bold label1 pb12">Plugin</p>
      <div class="fbv g20">
        <div
          v-for="item in pluginList"
          @click="jump(`/plugin?id=${item.id}`)"
          class="fbh fbac g10 hand"
        >
          <Icon :name="item.image" :size="32" />
          <div class="fbv g4">
            <div class="fbh fbac g12">
              <p class="f14 bold">{{ item.name }}</p>
              <div class="fbh fbac g4">
                <Icon :name="item.creator.head_pic" :size="14" class="br100" />
                <p class="label1 f13">@{{ item.creator.name }}</p>
              </div>
              <div class="f12 label2" v-if="item.linked_agent">
                Used by {{ item.linked_agent.length }} Agents
              </div>
            </div>
            <p class="f13 label1 omit1">{{ item.desc }}</p>
          </div>
        </div>
        <ElSkeleton v-if="loadingPlugin" animated />
        <div
          v-if="pluginData?.data.list?.length === 5"
          class="fbh fbjc fbac g4 pb8 hand"
          @click="refetchPlugin({page: ++pluginPage, keywords: search})"
        >
          <p class="f14 label2">More</p>
          <Icon name="arrow-down" :size="10"></Icon>
        </div>
      </div>
      <div v-if="!loadingPlugin && !pluginList.length" class="fbv fbac g12 p30">
        <Icon name="empty1" :width="114" :height="92" />
        <p class="label1 f14">Not Found</p>
      </div>
    </div>
  </div>
</template>

<style lang="less" scoped>
.input {
  width: 500px;
  height: 40px;
  outline: none;
  box-sizing: border-box;
  border: 1px solid var(--secondary-border-color);
  color: var(--primary-text-color);
  background: var(--tertiary-background-color);
  &::placeholder {
    color: var(--tertiary-label-color);
  }
}
.modal {
  left: 0;
  right: 0;
  bottom: 0;
  top: 80px;
}
.result {
  top: 50px;
  left: -100px;
  width: 700px;
  max-height: 500px;
  background: var(--reverse-text-color);
  border-radius: 8px;
  z-index: 2;
}
</style>
