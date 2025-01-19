<script lang="ts" setup>
import {
  BotDetailShape,
  createStudioConversationMutation,
  getStudioMessage,
  mockErrorMessage,
  mockUserMessage,
  sendStudioMessageMutation,
  StudioMessage,
  useStudioConversationQuery,
} from '@/apis'
import {Icon, LoadingCircle, LoadingDot, Scroll} from '@/components'
import Plugin from '@/pages/agent/plugin.vue'
import {getAvatar} from '@/utils'
import {ElTooltip} from 'element-plus'
import {debounce, throttle} from 'lodash-es'
import {onMounted, ref, watch} from 'vue'
import {useRouter} from 'vue-router'
import Content from '../agent/content.vue'
import Sender from '../agent/sender.vue'

const {bot, conversationId: rawConversationId} = defineProps<{
  bot: Pick<
    BotDetailShape,
    'id' | 'name' | 'image' | 'welcome_msg' | 'guide_info'
  >
  conversationId?: number
}>()
const page = ref(1)
const conversationId = ref(rawConversationId ?? -1)
const {
  data: conversationData,
  loading: loadingTop,
  refetch: refetchConversation,
} = useStudioConversationQuery({
  skip: !rawConversationId,
  params: {
    conversation_id: conversationId.value,
    page: page.value,
    page_size: 30,
  },
})
const {push} = useRouter()
const messages = ref<
  StudioMessage<{
    animation?: boolean
  }>[]
>([])
const scrollRef = ref<ScrollShape>()
const isFirstLoading = ref(true)
const getMessageLoading = ref(false)
const isConversationBeginning = ref(false)
const onAttachTop = throttle(async () => {
  if (isConversationBeginning.value || loadingTop.value || !messages.value) {
    return
  }
  await refetchConversation({
    page: ++page.value,
    conversation_id: conversationId.value,
  })
  if (!conversationData.value?.data) {
    isConversationBeginning.value = true
    return
  }
  messages.value.unshift(
    ...(conversationData.value.data.messages ?? []).reverse()
  )
}, 1000)

watch(getMessageLoading, () => {
  scrollRef.value?.onScrollToBottom()
})

watch(loadingTop, () => {
  if (!loadingTop.value && isFirstLoading.value) {
    isFirstLoading.value = false
    scrollRef.value?.onScrollToBottom()
  }
})

watch(conversationData, () => {
  if (!messages.value.length) {
    messages.value.unshift(
      ...(conversationData.value?.data.messages ?? []).reverse()
    )
  }
})

watch(conversationData, async () => {
  const latestMessage =
    conversationData.value?.data.messages?.[
      conversationData.value.data.messages.length - 1
    ]
  if (latestMessage?.status === 0) {
    getMessageLoading.value = true
    await getStudioMessage(
      conversationId.value,
      {message_id: latestMessage.id},
      (data) => {
        if (data.conversation_id === conversationId.value) {
          if (data.content) {
            messages.value.push({
              ...data,
              tool_calls: undefined,
              animation: true,
            })
          }
          if (data.tool_calls?.some((item) => item.plugin_id)) {
            messages.value.push({...data, content: ''})
          }
        }
      },
      () => {
        messages.value.push(mockErrorMessage())
        scrollRef.value?.onScrollToBottom()
      }
    )
    getMessageLoading.value = false
  }
})

const onSend = debounce(async (query: string) => {
  const res = await sendStudioMessageMutation({
    params: {
      content: query,
      conversation_id: conversationId.value,
      tool_results: undefined,
    },
  })
  messages.value.push(mockUserMessage(query))
  getMessageLoading.value = true
  await getStudioMessage(
    conversationId.value,
    {message_id: res?.data.resp_message_id!},
    (data) => {
      if (data.conversation_id === conversationId.value) {
        if (data.content) {
          messages.value.push({
            ...data,
            tool_calls: undefined,
            animation: true,
          })
        }
        if (data.tool_calls?.some((item) => item.plugin_id)) {
          messages.value.push({...data, content: ''})
        }
      }
    },
    () => {
      messages.value.push(mockErrorMessage())
      scrollRef.value?.onScrollToBottom()
    }
  )
  getMessageLoading.value = false
}, 500)

const createConversation = throttle(async () => {
  const res = await createStudioConversationMutation({
    params: {bot_id: bot.id, title: 'new', type: 'draft'},
  })
  if (res?.data) {
    getMessageLoading.value = false
    conversationId.value = res.data.conversation_id
    messages.value.length = 0
    refetchConversation({
      conversation_id: conversationId.value,
      page: (page.value = 1),
    })
  }
}, 2000)

onMounted(() => !rawConversationId && createConversation())
</script>

<template>
  <Scroll
    ref="scrollRef"
    class="scroll-container fbv g20 pr"
    @on-attach-top="onAttachTop"
    keep-position
  >
    <div v-if="loadingTop" class="center loading">
      <LoadingCircle />
    </div>
    <div class="fbv center g16 pb4">
      <Icon :name="bot.image" :size="48" class="br6" />
      <p class="f15">{{ bot.name }}</p>
    </div>
    <div v-if="bot.welcome_msg" class="message-container">
      <div class="message">
        <Icon class="icon" :name="bot.image ?? 'agent'" :size="36" />
        <Content :data="bot.welcome_msg" />
      </div>
      <div class="fbh fbac g12 pt12 fbw">
        <div
          v-for="item in bot.guide_info"
          class="suggestion f14 hand px16 py8 br8"
          @click="onSend(item)"
        >
          {{ item }}
        </div>
      </div>
    </div>
    <template
      v-for="(
        {message_type, content, status, animation, tool_calls}, index
      ) in messages.filter(({message_type}) => message_type !== 'system')"
    >
      <div
        v-if="(message_type !== 'tool' && content) || tool_calls"
        :class="['message-container', message_type === 'user' && 'fbh fbje']"
      >
        <div :class="['message f16', {error: status === 3}]">
          <Icon
            :class="message_type === 'user' ? 'right-icon' : 'icon'"
            :name="message_type === 'user' ? getAvatar() : bot.image"
            :size="36"
          />
          <p v-if="status === 3">
            An error occurred. Either the engine you requested does not exist or
            there was another issue processing your request. Please try again.
          </p>
          <template v-else>
            <Plugin
              v-if="tool_calls"
              :running="index === messages.length - 1"
              :data="tool_calls"
            />
            <p v-else-if="message_type === 'user'">{{ content }}</p>
            <Content v-else :data="content" :animation="animation" />
          </template>
        </div>
      </div>
    </template>
    <div
      class="message-container"
      v-if="getMessageLoading && !messages[messages.length - 1]?.tool_calls"
    >
      <div class="message">
        <Icon class="icon" :name="bot.image" :size="36" />
        <LoadingDot :size="24" color="var(--primary-text-color)" />
      </div>
    </div>
    <div class="fb1 fbv fbje pr bottom-container">
      <Sender :loading="getMessageLoading" :on-send="onSend" />
      <ElTooltip content="Delete chat history">
        <Icon
          @click="createConversation"
          class="icon pa"
          name="clean"
          :size="32"
        />
      </ElTooltip>
    </div>
  </Scroll>
</template>

<style lang="less" scoped>
.scroll-container {
  height: calc(100% - 80px);
  padding: 24px 24px 0;
  font-family: Oxygen;
  box-sizing: border-box;
  overflow: auto;
}
.message-container {
  width: calc(100% - 80px);
  max-width: 770px;
  position: relative;
  margin: auto;
}
.message {
  padding: 12px 16px;
  border-radius: 8px;
  word-break: break-word;
  background: var(--message-background-color);
  max-width: 100%;
  width: fit-content;
  box-sizing: border-box;
  &.error {
    color: #e35461;
  }
  .icon {
    border-radius: 100%;
    position: absolute;
    left: -44px;
    top: 0;
  }
  .right-icon {
    border-radius: 100%;
    position: absolute;
    right: -44px;
    top: 0;
  }
}
.loading {
  left: 0;
  right: 0;
  z-index: 999;
  position: absolute;
  animation: slide 0.2s forwards;
  @keyframes slide {
    from {
      transform: translateY(-100%);
      opacity: 0;
    }
    to {
      transform: translateY(0);
      opacity: 1;
    }
  }
}
.bottom-container {
  background: var(--primary-background-color);
  padding-bottom: 16px;
  position: sticky;
  bottom: 0;
  .icon {
    background: var(--primary-background-color);
    padding: 10px 8px 30px 44px;
    left: -50px;
    bottom: 0;
  }
}
.suggestion {
  border: 1px solid var(--secondary-border-color);
  color: var(--primary-label-color);
  &:hover {
    background: var(--tertiary-background-color);
    color: var(--primary-text-color);
  }
}
</style>
