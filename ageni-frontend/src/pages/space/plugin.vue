<script lang="ts" setup>
import {
  pluginListQuery,
  PluginMethodShape,
  PluginShape,
  usePluginDetailQuery,
  usePluginListQuery,
} from '@/apis'
import {CustomButton, Icon, IconGroup, Scroll} from '@/components'
import {Modal} from '@/pages/modal'
import {prepare} from '@/utils'
import {group, ungroup} from 'awesome-chart'
import {ElTooltip} from 'element-plus'
import {debounce} from 'lodash-es'
import {ref, watch} from 'vue'
import Example from './example.vue'

defineEmits<{
  addPlugin: [PluginShape]
  removePlugin: [PluginShape]
  close: []
}>()

const props = defineProps<{
  plugins: number[] | number
  variant: 'single' | 'multiple'
}>()
const page = ref(1)
const search = ref('')
const data = ref<PluginShape[]>([])
const expends = ref(new Set<number>())
const {
  data: pluginList,
  loading: loadingList,
  refetch,
} = usePluginListQuery({
  params: {page: 1, page_size: 10, keywords: '', tab: ''},
  skip: props.variant === 'single',
})
const {data: pluginDetail, loading: loadingDetail} = usePluginDetailQuery({
  params: {id: ungroup(props.plugins)!},
  skip: props.variant === 'multiple',
})
const activeMethod = ref<PluginMethodShape>()
const dependPlugins = ref(new Map<number, PluginShape>())
const ready = prepare(1000)

watch(data, async () => {
  const newIds = data.value.flatMap((item) =>
    (item.depend_plugins ?? []).filter((id) => !dependPlugins.value.has(id))
  )
  if (newIds.length) {
    const res = await pluginListQuery({
      params: {ids: newIds},
    })
    res?.data.list.forEach((item) => {
      dependPlugins.value.set(item.id, item)
    })
  }
})

watch(
  search,
  debounce(() => {
    data.value = []
    refetch({
      page: (page.value = 1),
      keywords: search.value,
    })
  }, 500)
)

watch(pluginList, () => {
  if (page.value === 1) {
    data.value = pluginList.value?.data.list ?? []
  } else {
    data.value.push(...(pluginList.value?.data.list ?? []))
  }
})

watch(pluginDetail, () => {
  if (pluginDetail.value) {
    data.value = [pluginDetail.value.data.plugin]
  }
})
</script>

<template>
  <Example
    v-if="activeMethod"
    :data="activeMethod"
    @close="activeMethod = undefined"
  />
  <Modal open @close="$emit('close')">
    <div class="modal-root p16">
      <div v-if="variant === 'multiple'" class="fbh fbac pl16">
        <p class="f18 bold pr100 mr20">Add plugin</p>
        <div class="pr center">
          <input
            type="text"
            ref="inputRef"
            placeholder="Search"
            class="input br8 f14"
            v-model="search"
          />
          <Icon name="search" class="search pa" :size="14" />
        </div>
      </div>
      <div v-else class="center f18 bold">Plugin</div>
      <Scroll
        class="scroll-container oa mt20 fbv g20"
        @on-attach-bottom="
          variant === 'multiple' && refetch({page: ++page, keywords: search})
        "
      >
        <div
          v-for="item in data"
          class="fbv g12 px12"
          :class="expends.has(item.id) && 'expend py20'"
        >
          <div class="fbh fbac g12">
            <Icon :name="item.image" :size="32" class="br6" />
            <div class="fbv g4 fb1">
              <p class="f14 bold">{{ item.name }}</p>
              <p class="f13 label1 omit1">{{ item.desc }}</p>
            </div>
            <CustomButton
              v-if="group(plugins).includes(item.id)"
              @click="$emit('removePlugin', item)"
              default-stroke="var(--tertiary-border-color)"
              default-color="var(--primary-label-color)"
              active-fill="var(--primary-text-color)"
              active-color="var(--reverse-text-color)"
              class="f14"
              :width="84"
              :height="30"
            >
              Remove
            </CustomButton>
            <CustomButton
              v-else
              @click="$emit('addPlugin', item)"
              class="f14"
              :width="84"
              :height="30"
            >
              Add
            </CustomButton>
            <Icon
              v-if="variant === 'multiple'"
              name="arrow-down1"
              :style="`rotate:${expends.has(item.id) ? 0 : -90}deg`"
              :size="12"
              @click="
                expends.has(item.id)
                  ? expends.delete(item.id)
                  : expends.add(item.id)
              "
            />
          </div>
          <div class="fbh fbac fbjsb">
            <div class="fbh fbac g8">
              <div class="chip label1 f12">
                {{ item.methods?.length ?? 0 }} Functions
              </div>
              <div
                v-if="item.depend_plugins"
                class="chip label1 f12 fbh fbac g4"
              >
                <ElTooltip>
                  <template #content>
                    Some plugins have scenarios where they depend on other
                    plugins. That is, if you want to add this plugin, you must
                    add its dependent plugins at the same time. Otherwise, the
                    plugin cannot work properly. In this scenario, when you add
                    a plugin with dependent plugins, its dependent plugins will
                    be automatically added.
                  </template>
                  <Icon name="query" :size="10" />
                </ElTooltip>
                Dependent Plugin:
                <ElTooltip
                  v-for="{name, image} in item.depend_plugins
                    .map((id) => dependPlugins.get(id)!)
                    .filter(Boolean)"
                  :content="name"
                >
                  <Icon :size="14" :name="image" />
                </ElTooltip>
              </div>
            </div>
            <div class="fbh fbac g4">
              <IconGroup v-if="item.linked_agent" :data="item.linked_agent" />
              <p class="label1 f14">
                Used by {{ item.linked_agent?.length }} Agents
              </p>
            </div>
          </div>
          <div class="divider mt8"></div>
          <div v-if="expends.has(item.id) || variant === 'single'">
            <div v-for="(method, i) in item.methods">
              <div class="divider my20" v-if="i"></div>
              <div class="fbh fbac fbjsb g24">
                <div class="fbv g4 method-content">
                  <p class="f16 bold">{{ method.name }}</p>
                  <p class="f13 label2 omit">{{ method.description }}</p>
                </div>
                <div
                  class="button f14 center br4 hand"
                  @click="activeMethod = method"
                >
                  Example
                </div>
              </div>
            </div>
          </div>
        </div>
        <template v-if="!data.length">
          <ElSkeleton v-if="loadingDetail || loadingList" animated />
          <div v-else-if="ready" class="fbv fbac g12 p100">
            <Icon name="empty1" :width="114" :height="92" />
            <p class="label1 f14">Not Found</p>
          </div>
        </template>
      </Scroll>
    </div>
  </Modal>
</template>

<style lang="less" scoped>
.modal-root {
  width: 800px;
}
.input {
  width: 300px;
  height: 32px;
  outline: none;
  box-sizing: border-box;
  padding: 10px 16px 10px 36px;
  border: 1px solid var(--secondary-border-color);
  color: var(--primary-text-color);
  background: var(--tertiary-background-color);
  &::placeholder {
    color: var(--tertiary-label-color);
  }
}
.search {
  left: 14px;
  bottom: 10px;
}
.chip {
  background: var(--tertiary-background-color);
  border-radius: 4px;
  padding: 2px 8px;
}
.scroll-container {
  height: 500px;
  min-height: 500px;
  .expend {
    background: var(--secondary-background-color);
    border-radius: 8px;
  }
}
.method-content {
  max-width: 600px;
}
.button {
  background: #44dac81f;
  color: #44dac8;
  width: 79px;
  height: 28px;
}
</style>
