<script lang="ts" setup>
import {
  createStudioConversationMutation,
  deleteConversationMutation,
  useAgentDetailQuery,
  useConversationListQuery,
} from '@/apis'
import {Icon, Image, MobileHeader, Scroll, showToast} from '@/components'
import {scale, sleep, windowHeight, withLogin} from '@/utils'
import {ElTooltip} from 'element-plus'
import {ref, watch} from 'vue'
import {useRouter} from 'vue-router'
import Message from './message.vue'

const page = ref(1)
const expendInfo = ref(false)
const expendText = ref(false)
const {id} = defineProps<{id: number}>()
const {data: detail} = useAgentDetailQuery({
  params: {bot_id: id},
})
const selectConversationId = ref<number>()
const messageRef = ref<{onSend: (query: string) => void}[]>([])
const {data: conversationData, refetch: refetchHistory} =
  useConversationListQuery({
    params: {bot_id: id, page_size: 30, page: page.value},
  })
const conversations = ref<
  NonNullable<(typeof conversationData)['value']>['data']
>([])
const {back} = useRouter()

watch(conversationData, () => {
  if (page.value === 1) {
    conversations.value = conversationData.value?.data ?? []
  } else {
    conversations.value.push(...(conversationData.value?.data ?? []))
  }
})

const createConversation = withLogin(async (query?: string) => {
  if (!query) {
    selectConversationId.value = undefined
    return
  }
  const res = await createStudioConversationMutation({
    params: {bot_id: id, title: query},
  })
  if (res?.data) {
    conversations.value.unshift({
      id: res.data.conversation_id,
      title: query,
    })
    refetchHistory({page: (page.value = 1)})
    selectConversationId.value = res.data.conversation_id
    await sleep(100)
    messageRef.value[0].onSend(query)
  } else {
    showToast(res?.message!)
  }
})

async function deleteConversation(id: number) {
  await deleteConversationMutation({params: {id}})
  if (id === selectConversationId.value) {
    selectConversationId.value = undefined
  }
  refetchHistory({page: (page.value = 1)})
}
</script>

<template>
  <div v-if="!detail?.data" class="wh100p p20">
    <ElSkeleton :rows="5" animated />
  </div>
  <div v-else class="fbv wh100p">
    <MobileHeader v-if="expendInfo" class="header fbh fbac g16 py8 px16">
      <Icon name="arrow-left" :size="16" @click="expendInfo = false" />
      <p class="f16 bold">{{ detail.data.name }}</p>
    </MobileHeader>
    <MobileHeader v-else class="header fbh fbac g16 py8 px16">
      <Icon name="arrow-left" :size="16" @click="back" />
      <p class="f16 bold fb1">{{ detail.data.name }}</p>
      <Icon name="more2" :size="16" @click="expendInfo = true" />
    </MobileHeader>
    <div v-if="expendInfo" class="fbv oa">
      <div class="divider"></div>
      <div class="py24 px16 fbv g20">
        <div class="fbh fbac fbjsb">
          <div class="fbh fbac f16 g4">
            <Icon name="user1" :size="16" />
            <p class="pr24">{{ detail.data.users }}</p>
            <Icon name="comment" :size="16" />
            <p>{{ detail.data.conversations }}</p>
          </div>
        </div>
        <div class="pr f14 label1">
          <p :class="{omit4: !expendText}">
            {{ detail.data.desc }}
          </p>
          <div
            v-if="false"
            class="expend-more pl30 fbh fbac pa hand"
            @click="expendText = true"
          >
            More
          </div>
        </div>
        <div class="fbv g12">
          <p class="f16 bold">Plugins</p>
          <div class="fbh fbac g10">
            <ElTooltip
              v-for="item in detail.data.linked_plugin"
              :content="item.name"
            >
              <Image :name="item.image" :width="24" :height="24" class="br4" />
            </ElTooltip>
          </div>
        </div>
        <div class="fbv g12">
          <p class="f16 bold">Workmode</p>
          <div class="fbh fbac g8 fbw">
            <div v-for="item in detail.data.work_modes" class="chip">
              {{ item.name }}
            </div>
          </div>
        </div>
        <div v-if="detail.data.configuration" class="fbv g12">
          <p class="f16 bold">Configuration item</p>
          <div class="fbh fbac g8 fbw">
            <div v-for="item in detail.data.configuration" class="chip">
              {{ item.name }}
            </div>
          </div>
        </div>
      </div>
      <div class="divider"></div>
      <Scroll
        class="fb1 fbv g10 p16 oa"
        @on-attach-bottom="refetchHistory({page: ++page})"
      >
        <p class="f16 bold">History</p>
        <div class="fbv g12">
          <div
            v-for="item in conversations"
            class="history br8 hand pr"
            @click="(selectConversationId = item.id), (expendInfo = false)"
            :key="item.id"
          >
            <p class="omit2">{{ item.title }}</p>
            <div class="pa more br100">
              <ElDropdown trigger="click">
                <Icon
                  name="more1"
                  class="p10"
                  :width="13"
                  :height="3"
                  @click.stop=""
                />
                <template #dropdown>
                  <div
                    class="fbh fbac py10 px16 g6 hand"
                    @click="deleteConversation(item.id)"
                  >
                    <Icon name="delete" :size="14" />
                    <p class="f14 label1">Delete</p>
                  </div>
                </template>
              </ElDropdown>
            </div>
          </div>
        </div>
      </Scroll>
    </div>
    <div v-show="!expendInfo">
      <Message
        v-if="!selectConversationId"
        :conversation-id="Infinity"
        :create-conversation="createConversation"
        :style="`height: ${windowHeight - 44 * scale}px`"
        :bot="detail.data"
        welcome
      />
      <template v-for="item in conversations">
        <Message
          ref="messageRef"
          v-if="selectConversationId === item.id"
          :conversation-id="item.id"
          :create-conversation="createConversation"
          :style="`height: ${windowHeight - 44 * scale}px`"
          :bot="detail.data"
        />
      </template>
    </div>
  </div>
</template>

<style lang="less" scoped>
.container {
  height: calc(100% - 74px);
}
:deep(.header) {
  box-sizing: border-box;
  height: 44px;
}
.chip {
  border-radius: 4px;
  color: var(--primary-label-color);
  background: var(--tertiary-background-color);
  padding: 2px 8px;
  font-size: 12px;
}
.history {
  color: var(--primary-label-color);
  border: 1px solid var(--secondary-border-color);
  font-size: 14px;
  padding: 12px;
  &:hover {
    color: var(--primary-text-color);
    border-color: var(--primary-text-color);
  }
}
.expend-more {
  right: 0;
  bottom: 0;
  color: var(--theme-color);
  background: var(--gradient);
  --gradient: linear-gradient(
    to right,
    transparent 0%,
    var(--primary-background-color) 20%
  );
}
.more {
  height: fit-content;
  right: 0;
  top: 0;
}
</style>
