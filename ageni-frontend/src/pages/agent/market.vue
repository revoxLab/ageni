<script lang="ts" setup>
import {AgentShape, useAgentListQuery, useAgentTabsQuery} from '@/apis'
import {
  CustomButton,
  Icon,
  IconGroup,
  Image,
  Scroll,
  Select,
  TabSwitch,
} from '@/components'
import {checkLogin, jump, prepare, shortAddress} from '@/utils'
import {computed, ref, watch} from 'vue'
import {useRouter} from 'vue-router'
import Search from '../plugin/search.vue'

const page = ref(1)
const activeTab = ref()
const {push} = useRouter()
const list = ref<AgentShape[]>([])
const {data: tabs} = useAgentTabsQuery()
const {data, loading, refetch} = useAgentListQuery({
  params: {pick_type: 1, page: 1, page_size: 27, tab: '', keywords: ''},
  skip: true,
})
const tabData = computed(() =>
  tabs.value?.data.map((key) => ({key, title: key}))
)
const ready = prepare(1000)
const filterOptions: SelectOption[] = [
  {key: 1, label: 'Top Picks', value: 1},
  {key: 0, label: 'Show All', value: 0},
]
const selected = ref(filterOptions[0])

watch(tabData, () => {
  activeTab.value = tabData.value?.[0].key
})

watch([activeTab, selected], () => {
  list.value = []
  refetch({
    tab: activeTab.value,
    page: (page.value = 1),
    pick_type: selected.value.value,
  })
})

watch(data, () => {
  list.value.push(...(data.value?.data.bots ?? []))
})
</script>

<template>
  <div class="m24">
    <div class="fbh fbac fbjsb py8">
      <div class="fbh fbac g32">
        <p class="f20 bold">Agent Marketplace</p>
        <Search />
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
          class="f14"
          @click="checkLogin(() => push('/space'))"
          :width="138"
          :height="36"
        >
          + Create Agent
        </CustomButton>
      </div>
    </div>
    <div class="fbh fbjsb">
      <TabSwitch
        :data="tabData ?? []"
        :active-key="activeTab"
        @on-active-key-change="(value) => (activeTab = value)"
        active-color="var(--theme-color)"
        class="pt8 pb16"
      />
      <Select
        :options="filterOptions"
        :selected="selected"
        @on-selected-change="(value) => (selected = value)"
        class="filter mt12"
      />
    </div>
    <ElSkeleton v-if="loading && !list.length" :rows="5" animated />
    <div v-else-if="ready && !list.length" class="fbv empty center">
      <Image name="skeleton" :width="268" :height="98" />
      <p class="f18 bold pt24">No creations here</p>
      <p class="f14 label2 pt16">
        Build an Al Agent with the power of LLM and plugins in minutes
      </p>
      <CustomButton
        :width="133"
        :height="36"
        class="f14 pt20"
        @click="checkLogin(() => push('/space'))"
      >
        + Create Agent
      </CustomButton>
    </div>
    <Scroll
      v-else
      class="card-group oa"
      @on-attach-bottom="
        refetch({page: ++page, tab: activeTab, pick_type: selected.value})
      "
    >
      <div
        v-for="item in list"
        class="p20 pb14 section-card fbv g16 hand"
        @click="checkLogin(() => jump(`/agent?id=${item.id}`))"
      >
        <div class="fbh fbac g16">
          <Icon :name="item.image ?? 'agent'" :size="56" class="br6" />
          <div class="fbv g8">
            <p class="f18 bold">{{ item.name }}</p>
            <div class="fbh fbac g4">
              <Icon :name="item.creator.head_pic" :size="14" class="br100" />
              <p class="label1 f13">@{{ shortAddress(item.creator.name) }}</p>
            </div>
          </div>
        </div>
        <p class="f13 label1 omit3">{{ item.desc }}</p>
        <div class="divider"></div>
        <div class="fbh fbac f12 label2 g4">
          <Icon name="user" style="opacity: 0.4" :size="16" />
          <p class="pr20">{{ item.users }}</p>
          <Icon name="comment" style="opacity: 0.4" :size="16" />
          <p>{{ item.conversations }}</p>
          <div v-if="item.linked_plugin" class="fb1 fbh fbje pr">
            <IconGroup :data="item.linked_plugin" />
          </div>
        </div>
      </div>
    </Scroll>
  </div>
</template>

<style lang="less" scoped>
.empty {
  width: 700px;
  height: 540px;
  background: url(@/assets/images/empty_bg.png);
  background-size: 100% 100%;
  margin: auto;
}
.filter {
  right: 0;
  min-width: 140px;
}
</style>
