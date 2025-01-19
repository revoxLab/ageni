<script lang="ts" setup>
import {PluginShape, usePluginListQuery, usePluginTabsQuery} from '@/apis'
import {CustomButton, Icon, IconGroup, Image, Scroll} from '@/components'
import {jump, prepare} from '@/utils'
import {ElTooltip} from 'element-plus'
import {computed, ref, watch} from 'vue'
import {useRouter} from 'vue-router'
import Search from './search.vue'

const page = ref(1)
const activeTab = ref()
const {push} = useRouter()
const list = ref<PluginShape[]>([])
const {data: tabs} = usePluginTabsQuery()
const {data, loading, refetch} = usePluginListQuery({
  params: {page: 1, page_size: 30, tab: '', keywords: ''},
  skip: true,
})
const tabData = computed(() =>
  tabs.value?.data.map((key) => ({key, title: key}))
)
const ready = prepare(1000)

watch(tabData, () => {
  activeTab.value = tabData.value?.[0].key
})

watch(activeTab, () => {
  list.value = []
  refetch({tab: activeTab.value, page: (page.value = 1)})
})

watch(data, () => {
  list.value.push(...(data.value?.data.list ?? []))
})
</script>

<template>
  <div class="m24">
    <div class="fbh fbac fbjsb py8 mb8">
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
        <ElTooltip content="Coming soon">
          <CustomButton class="f14" :width="138" :height="36">
            + Create Plugin
          </CustomButton>
        </ElTooltip>
      </div>
    </div>
    <ElSkeleton v-if="loading && !list.length" :rows="5" animated />
    <div v-else-if="ready && !list.length" class="fbv empty center">
      <Image name="skeleton1" :width="268" :height="72" />
      <p class="f18 bold pt24">Plugin not found</p>
      <p class="f14 label2 pt16">Click button to create a plugin</p>
      <ElTooltip content="Coming Soon">
        <CustomButton :width="133" :height="36" class="f14 pt20">
          + Create Plugin
        </CustomButton>
      </ElTooltip>
    </div>
    <Scroll
      v-else
      class="card-group oa"
      @on-attach-bottom="refetch({page: ++page, tab: activeTab})"
    >
      <div
        v-for="item in list"
        class="p20 pb14 section-card fbv g12 hand"
        @click="push(`/plugin?id=${item.id}`)"
      >
        <Icon :name="item.image" :size="56" class="br6" />
        <div class="fbv g8">
          <p class="f16 bold">{{ item.name }}</p>
          <div class="fbh fbac g4">
            <Icon :name="item.creator.head_pic" :size="14" class="br100" />
            <p class="label1 f13">@{{ item.creator.name }}</p>
          </div>
        </div>
        <p class="f13 label1 omit1">{{ item.desc }}</p>

        <div v-if="item.linked_agent" class="fbh fbac g4">
          <IconGroup :data="item.linked_agent" />
          <p class="label1 f12">
            Used by {{ item.linked_agent.length }} Agents
          </p>
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
</style>
