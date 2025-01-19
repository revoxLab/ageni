<script lang="ts" setup>
import {pluginList, StudioMessage} from '@/apis'
import {LoadingDot} from '@/components'
import Icon from '@/components/common/icon.vue'

const props = defineProps<{
  data: StudioMessage['tool_calls']
  running: boolean
}>()
const plugins = props.data
  ?.map(
    ({plugin_id}) =>
      pluginList.value?.data.list.find(({id}) => id === plugin_id)!
  )
  .filter(Boolean)
</script>

<template>
  <div class="fbh g16">
    <span v-if="running">
      <LoadingDot :size="24" color="var(--primary-text-color)" />
    </span>
    <div class="fbv g12">
      <div v-for="(item, i) in plugins" class="fbh fbac g8">
        <Icon class="br4" :name="item.image" :size="20" />
        <p class="f16">
          {{ item.name }} -
          <span class="chip br4 px8 py2 code">
            {{ data?.[i].function.name }}
          </span>
          {{ running ? 'is being called' : 'has been called' }}
        </p>
      </div>
    </div>
  </div>
</template>

<style lang="less" scoped>
.chip {
  border: 1px solid var(--secondary-border-color);
  background: var(--secondary-background-color);
}
.code {
  font-family: Consola;
}
</style>
