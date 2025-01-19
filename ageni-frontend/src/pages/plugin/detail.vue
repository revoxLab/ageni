<script lang="ts" setup>
import {usePluginDetailQuery} from '@/apis'
import {Icon, IconGroup, TabSwitch} from '@/components'
import {format} from 'date-fns'
import {computed, ref, watch} from 'vue'
import {useRouter} from 'vue-router'

const {push} = useRouter()
const {id} = defineProps<{id: number}>()
const {data, loading} = usePluginDetailQuery({params: {id}})
const detail = computed(() => data.value?.data.plugin)
const functionTab = ref()
const activeMethod = computed(() =>
  detail.value?.methods?.find(({id}) => id === functionTab.value)
)
const methodType = ref<'request' | 'response'>('request')
const activeSchema = computed(() =>
  JSON.parse(
    (methodType.value === 'request'
      ? activeMethod.value?.input_schema
      : activeMethod.value?.output_schema) ?? 'null'
  )
)

watch(detail, () => {
  functionTab.value = detail.value?.methods?.[0].id
})
</script>

<template>
  <ElSkeleton v-if="loading" :rows="5" class="p24" animated />
  <div v-else-if="detail" class="m24">
    <Icon
      @click="push('/plugin')"
      name="arrow-left"
      style="opacity: 0.4"
      :size="12"
    />
    <div class="section-card mt30 p24 f14">
      <div class="fbh fbac g20">
        <Icon :name="detail.image" :size="88" class="br6" />
        <div class="fbv g8">
          <p class="f26 bold">{{ detail.name }}</p>
          <div class="fbh fbac g4">
            <Icon :name="detail.creator.head_pic" :size="14" class="br100" />
            <p class="label1">
              @{{ detail.creator.name }} · Updated at
              {{
                format(detail.created_at.seconds * 1000, 'yyy-MM-dd HH:mm:ss')
              }}
            </p>
            <template v-if="detail.linked_agent">
              <IconGroup class="ml24" :data="detail.linked_agent" />
              <p class="label1">
                Used by {{ detail.linked_agent.length }} Agents
              </p>
            </template>
          </div>
          <p class="label1 omit3 pt8">{{ detail.desc }}</p>
        </div>
      </div>
    </div>
    <div class="section-card bottom p24 fbv g16">
      <p class="f18 bold">Functions</p>
      <TabSwitch
        active-color="white"
        :active-key="functionTab"
        :data="
          detail.methods?.map(({id, name}) => ({key: id, title: name})) ?? []
        "
        @on-active-key-change="(value) => (functionTab = value)"
        class="oa"
      />
      <div class="f14 label1 p16 description br8">
        {{ activeMethod?.description }}
      </div>
      <div class="fbh f13">
        <div
          class="tab hand br4"
          :class="{active: methodType === 'request'}"
          @click="methodType = 'request'"
        >
          Request
        </div>
        <div
          class="tab hand br4"
          :class="{active: methodType === 'response'}"
          @click="methodType = 'response'"
        >
          Response
        </div>
      </div>
      <div class="table fbh br8 f14 label1">
        <div class="fb1 fbv">
          <div class="fbh header p16">
            <div class="fb1">Parameter</div>
            <div class="fb1">Description</div>
          </div>
          <template v-for="(item, i) in activeSchema?.parameters">
            <div v-if="i !== 0" class="divider mx16"></div>
            <div class="fbh">
              <div class="fb1 p16">{{ item.name }}</div>
              <div class="fb1 p16">{{ item.description }}</div>
            </div>
          </template>
        </div>
        <div class="fb1 example oa">
          <div class="p16 header w100p">Example</div>
          <pre class="p16 oa" style="text-wrap: wrap">{{
            methodType === 'request'
              ? activeMethod?.input_example
              : activeMethod?.output_example
          }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="less" scoped>
.bottom {
  margin-top: -0.5px;
}
.tab {
  padding: 5px 10px;
  color: var(--secondary-label-color);
  &.active {
    background: var(--tertiary-background-color);
    color: var(--primary-text-color);
  }
}
.description {
  border: 1px solid var(--secondary-border-color);
}
.table {
  border: 1px solid var(--secondary-border-color);
}
.header {
  background: var(--secondary-background-color);
  box-sizing: border-box;
}
.example {
  border-left: 1px solid var(--secondary-border-color);
}
</style>
