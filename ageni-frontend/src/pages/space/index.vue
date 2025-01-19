<script lang="ts" setup>
import {createAgentMutation, UserAgentShape, useUserAgentList} from '@/apis'
import {
  CustomButton,
  Icon,
  IconGroup,
  Image,
  Scroll,
  TabSwitch,
} from '@/components'
import {checkLogin, getAvatar, withLogin} from '@/utils'
import {uuid} from 'awesome-chart'
import {format} from 'date-fns'
import {ElTooltip} from 'element-plus'
import {ref, watch} from 'vue'
import AgentModal from './agent.vue'

const page = ref(1)
const list = ref<UserAgentShape[]>([])
const {data, loading, refetch} = useUserAgentList({
  params: {page: 1, page_size: 27},
})
const tabData = [
  {key: 'Agent', title: 'Agent'},
  {key: 'Plugin', title: 'Plugin'},
]
const activeTab = ref<'Plugin' | 'Agent'>(tabData[0].key as 'Agent')
const activeAgentId = ref<number | undefined>()

watch(data, () => {
  if (page.value === 1) {
    list.value = data.value?.data?.bots ?? []
  } else {
    list.value.push(...(data.value?.data?.bots ?? []))
  }
})

const createAgent = withLogin(async () => {
  const res = await createAgentMutation({
    params: {name: `Draft(#${uuid(6, '123456789')})`},
  })
  activeAgentId.value = res?.data.id
})
</script>

<template>
  <AgentModal
    v-if="activeAgentId"
    :id="activeAgentId"
    @close="(activeAgentId = undefined), refetch({page: (page = 1)})"
  />
  <div class="fbh fbac g10 py20 px24">
    <Icon :name="getAvatar()" :size="32" class="br100" />
    <p class="fb1 f18 bold">Personal Space</p>
    <ElTooltip v-if="activeTab === 'Plugin'" content="Coming Soon">
      <CustomButton :width="133" :height="36" class="f14">
        + Create Plugin
      </CustomButton>
    </ElTooltip>
    <CustomButton
      v-else
      :width="133"
      :height="36"
      class="f14"
      @click="createAgent"
    >
      + Create Agent
    </CustomButton>
  </div>
  <div class="divider"></div>
  <div class="px24 py8 fbv g12">
    <TabSwitch
      :data="tabData"
      :active-key="activeTab"
      @on-active-key-change="(value) => checkLogin(() => (activeTab = value))"
      active-color="var(--theme-color)"
    />
    <div v-if="activeTab === 'Plugin'" class="fbv empty center">
      <Image name="skeleton1" :width="268" :height="72" />
      <p class="f18 bold pt24">Plugin not found</p>
      <p class="f14 label2 pt16">Click button to create a plugin</p>
      <ElTooltip content="Coming Soon">
        <CustomButton :width="133" :height="36" class="f14 pt20">
          + Create Plugin
        </CustomButton>
      </ElTooltip>
    </div>
    <ElSkeleton v-else-if="loading && !list.length" :rows="5" animated />
    <div v-else-if="!list.length" class="fbv empty center">
      <Image name="skeleton" :width="268" :height="98" />
      <p class="f18 bold pt24">No creations here</p>
      <p class="f14 label2 pt16">
        Build an Al Agent with the power of LLM and plugins in minutes
      </p>
      <CustomButton
        :width="133"
        :height="36"
        class="f14 pt20"
        @click="createAgent"
      >
        + Create Agent
      </CustomButton>
    </div>
    <Scroll
      v-else
      class="card-group oa"
      @on-attach-bottom="refetch({page: ++page})"
    >
      <div
        v-for="item in list"
        class="p20 pb14 section-card fbv g8 hand"
        @click="activeAgentId = item.id"
      >
        <div class="fbh fbac fbjsb g16">
          <div class="fbv g6">
            <p class="f18 bold">{{ item.name }}</p>
            <p class="f13 label1 omit3">{{ item.desc }}</p>
          </div>
          <Icon :name="item.image ?? 'agent'" :size="56" class="br6" />
        </div>
        <p class="f12 label2 pt8">
          Last {{ item.status === 2 ? 'publish' : 'edit' }}
          {{ format(item.created_at.seconds * 1000, 'yyyy.MM.dd HH:mm') }}
        </p>
        <div class="fbh fbac fbjsb">
          <div
            class="chip br4 p8 fbh fbac g4 f12"
            :class="{publish: item.status === 2}"
          >
            <Icon :name="item.status === 2 ? 'check1' : 'time'" :size="10" />
            {{ item.status === 2 ? 'Published' : 'Draft' }}
          </div>
          <IconGroup v-if="item.linked_plugin" :data="item.linked_plugin" />
        </div>
      </div>
      <CustomButton
        :width="193"
        :height="46"
        class="f16 center pt40"
        style="grid-column: 1 / -1"
        @click="createAgent"
      >
        + Create Agent
      </CustomButton>
    </Scroll>
  </div>
</template>

<style lang="less" scoped>
.chip {
  color: var(--primary-label-color);
  background: var(--secondary-border-color);
  &.publish {
    background: #00c0871f;
    color: #00c087;
  }
}
.empty {
  width: 700px;
  height: 540px;
  background: url(@/assets/images/empty_bg.png);
  background-size: 100% 100%;
  margin: auto;
}
</style>
