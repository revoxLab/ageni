<script lang="ts" setup>
import {AgentShape, useAgentListQuery, useAgentTabsQuery} from '@/apis'
import {
  Icon,
  MobileHeader,
  mobileMenuOpen,
  Scroll,
  TabSwitch,
} from '@/components'
import Image from '@/components/common/image.vue'
import {prepare, shortAddress, windowHeight} from '@/utils'
import {DropdownInstance, ElDropdown, ElTooltip} from 'element-plus'
import {debounce} from 'lodash-es'
import {computed, ref, watch} from 'vue'
import {useRouter} from 'vue-router'

const page = ref(1)
const search = ref('')
const activeTab = ref()
const {push} = useRouter()
const list = ref<AgentShape[]>([])
const {data: tabs} = useAgentTabsQuery()
const {data, loading, refetch} = useAgentListQuery({
  params: {pick_type: 1, page: 1, page_size: 30, tab: '', keywords: ''},
})
const tabData = computed(() =>
  tabs.value?.data.map((key) => ({key, title: key}))
)
const dropdownRef = ref<DropdownInstance>()
const isTopPicks = ref(true)
const ready = prepare(1000)

watch(
  search,
  debounce(() => {
    data.value = undefined
    refetch({
      page: (page.value = 1),
      keywords: search.value,
      pick_type: isTopPicks.value ? 1 : 0,
    })
  }, 500)
)

watch(tabData, () => {
  activeTab.value = tabData.value?.[0].key
})

watch([activeTab, isTopPicks], () => {
  refetch({
    tab: activeTab.value,
    page: (page.value = 1),
    pick_type: isTopPicks.value ? 1 : 0,
  })
})

watch(data, () => {
  if (page.value === 1) {
    list.value = data.value?.data.bots ?? []
  } else {
    list.value.push(...(data.value?.data.bots ?? []))
  }
})
</script>

<template>
  <Scroll
    class="px16 pb20 oa"
    :style="`height:${windowHeight}px`"
    @on-attach-bottom="refetch({page: ++page})"
  >
    <MobileHeader class="px16">
      <div class="fbh fbac g18 py16">
        <Icon name="menu" :size="18" @click="mobileMenuOpen = true" />
        <p class="f18 bold">Agent Marketplace</p>
      </div>
      <input
        type="text"
        ref="inputRef"
        placeholder="Search"
        class="input w100p br8 px16 py10 f14 bold"
        v-model="search"
      />
      <div class="fbh fbjsb">
        <TabSwitch
          :data="tabData ?? []"
          :active-key="activeTab"
          @on-active-key-change="(value) => (activeTab = value)"
          active-color="var(--theme-color)"
          class="pt8 pb16 oa pr60"
        />
        <ElDropdown v-if="tabData" ref="dropdownRef">
          <div class="filter pa">
            <Icon name="filter" class="icon pa p8" :size="16"></Icon>
            <div class="shadow pa"></div>
          </div>
          <template #dropdown>
            <div
              class="fbv fbac g14 p12 f12 label1 pr30"
              style="text-wrap: nowrap"
            >
              <div
                class="fbh fbac g6 hand"
                @click="(isTopPicks = true), dropdownRef?.handleClose()"
              >
                <Icon
                  :size="18"
                  :name="isTopPicks ? 'checkbox_select' : 'checkbox'"
                />
                <p>Top Picks</p>
              </div>
              <div
                class="fbh fbac g6 hand"
                @click="(isTopPicks = false), dropdownRef?.handleClose()"
              >
                <Icon
                  :size="18"
                  :name="!isTopPicks ? 'checkbox_select' : 'checkbox'"
                />
                <p>Show All</p>
              </div>
            </div>
          </template>
        </ElDropdown>
      </div>
    </MobileHeader>
    <ElSkeleton v-if="loading && !list.length" :rows="5" animated />
    <div v-else-if="ready && !list.length" class="fbv center p30">
      <Image name="skeleton" :width="200" :height="75" />
      <p class="f16 bold pt12">No creations here</p>
    </div>
    <div
      v-for="item in list"
      class="p16 card fbv g16 hand"
      @click="push(`/agent?id=${item.id}`)"
    >
      <div class="fbh fbac g16">
        <Icon :name="item.image" :size="56" class="br6" />
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
        <Icon name="user" style="opacity: 0.4" :size="12" />
        <p class="pr20">{{ item.users }}</p>
        <Icon name="comment" style="opacity: 0.4" :size="12" />
        <p>{{ item.conversations }}</p>
        <div v-if="item.linked_plugin" class="fb1 fbh fbje pr">
          <ElTooltip>
            <div class="fbh">
              <Icon
                v-for="plugin in item.linked_plugin.slice(0, 4)"
                style="margin-left: -5px"
                :name="plugin.image"
                :size="18"
              />
              <Icon
                v-if="item.linked_plugin.length > 4"
                style="margin-left: -5px"
                name="more3"
                :size="18"
              />
            </div>
            <template #content>
              <div class="icon-group">
                <div v-for="plugin in item.linked_plugin" class="fbh fbac g4">
                  <Icon :name="plugin.image" :size="18" />
                  <p class="f14 label1">{{ plugin.name }}</p>
                </div>
              </div>
            </template>
          </ElTooltip>
        </div>
      </div>
    </div>
  </Scroll>
</template>

<style lang="less" scoped>
.card {
  position: relative;
  margin-bottom: -1px;
  &::before {
    content: '';
    border-radius: 8px;
    background: var(--primary-background-color);
    border: 1px solid var(--primary-border-color);
    box-sizing: border-box;
    position: absolute;
    z-index: -1;
    height: 100%;
    width: 100%;
    left: 0;
    top: 0;
  }
  &::after {
    content: '';
    background: var(--primary-border-color);
    position: absolute;
    z-index: -2;
    height: 100%;
    width: 100%;
    left: 0;
    top: 0;
  }
}
.icon-group {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  column-gap: 24px;
  row-gap: 12px;
}
.input {
  height: 38px;
  outline: none;
  box-sizing: border-box;
  border: 1px solid var(--secondary-border-color);
  color: var(--primary-text-color);
  background: var(--tertiary-background-color);
  &::placeholder {
    color: var(--tertiary-label-color);
  }
}
.filter {
  right: 0;
  bottom: 16px;
  .icon {
    z-index: 1;
    right: 0;
    bottom: 0;
  }
  .shadow {
    right: 0;
    bottom: 0;
    width: 60px;
    height: 32px;
    background: linear-gradient(
      to right,
      transparent,
      var(--primary-background-color)
    );
  }
}
</style>
