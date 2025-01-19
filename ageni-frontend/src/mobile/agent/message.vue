<script lang="ts" setup>
import {
  BotDetailShape,
  getStudioMessage,
  mockErrorMessage,
  mockUserMessage,
  sendStudioMessageMutation,
  StudioMessage,
  useStudioConversationQuery,
} from '@/apis'
import {Icon, LoadingCircle, LoadingDot, Scroll} from '@/components'
import {checkLogin, parseSmartWalletGuideInfo} from '@/utils'
import {ElTooltip} from 'element-plus'
import {debounce, throttle} from 'lodash-es'
import {computed, ref, watch} from 'vue'
import {useRouter} from 'vue-router'
import Content from './content.vue'
import Plugin from './plugin.vue'
import Sender from './sender.vue'

const props = defineProps<{
  welcome?: boolean
  bot: BotDetailShape
  conversationId: number
  createConversation: (query?: string) => void
}>()
const page = ref(1)
const {
  data: conversationData,
  loading: loadingTop,
  refetch: refetchConversation,
} = useStudioConversationQuery({
  skip: !!props.welcome,
  params: {
    conversation_id: props.conversationId,
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
const loadingBottom = computed(
  () =>
    getMessageLoading.value &&
    !messages.value[messages.value.length - 1].tool_calls
)
const onAttachTop = throttle(async () => {
  if (isConversationBeginning.value || loadingTop.value || !messages.value) {
    return
  }
  await refetchConversation({
    page: ++page.value,
  })
  if (!conversationData.value?.data) {
    isConversationBeginning.value = true
    return
  }
  messages.value.unshift(
    ...(conversationData.value.data.messages ?? []).reverse()
  )
}, 1000)

watch(loadingBottom, () => {
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
      props.conversationId,
      {message_id: latestMessage.id},
      (data) => {
        if (data.content) {
          messages.value.push({...data, tool_calls: undefined, animation: true})
        }
        if (data.tool_calls?.some((item) => item.plugin_id)) {
          messages.value.push({...data, content: ''})
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
      conversation_id: props.conversationId,
      tool_results: undefined,
    },
  })
  messages.value.push(mockUserMessage(query))
  getMessageLoading.value = true
  await getStudioMessage(
    props.conversationId,
    {message_id: res?.data.resp_message_id!},
    (data) => {
      if (data.content) {
        messages.value.push({...data, tool_calls: undefined, animation: true})
      }
      if (data.tool_calls?.some((item) => item.plugin_id)) {
        messages.value.push({...data, content: ''})
      }
    },
    () => {
      messages.value.push(mockErrorMessage())
      scrollRef.value?.onScrollToBottom()
    }
  )
  getMessageLoading.value = false
}, 500)

defineExpose({onSend})
</script>

<template>
  <Scroll
    ref="scrollRef"
    class="message-root fbv g20 pr"
    @on-attach-top="onAttachTop"
    keep-position
  >
    <div v-if="loadingTop" class="center loading">
      <LoadingCircle />
    </div>
    <div v-if="welcome">
      <div class="message">
        <Content :data="bot.welcome_msg ?? ''" />
      </div>
      <div class="fbh fbac pt12 g10 fbw">
        <div
          v-for="{text, icon} in bot.guide_info?.map(parseSmartWalletGuideInfo)"
          class="suggestion f14 hand px16 py8 br8 g6 fbh fbac"
          @click="checkLogin(() => createConversation(text))"
        >
          {{ text }}
          <Icon v-if="icon" :name="icon" :size="20" class="br4" />
        </div>
      </div>
    </div>
    <template
      v-for="(
        {message_type, content, status, animation, tool_calls}, index
      ) in messages"
    >
      <div
        v-if="(message_type !== 'tool' && content) || tool_calls"
        :class="[
          'message f16',
          {error: status === 3, user: message_type === 'user'},
        ]"
      >
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
      <div v-if="message_type === 'system'" class="fbh fbac g10 mt-8 fbw">
        <div
          v-for="{text, icon} in bot.guide_info?.map(parseSmartWalletGuideInfo)"
          class="suggestion f14 hand px16 py8 br8 g6 fbh fbac"
          @click="onSend(text)"
        >
          {{ text }}
          <Icon v-if="icon" :name="icon" :size="20" class="br4" />
        </div>
      </div>
    </template>
    <div v-if="loadingBottom" class="message">
      <LoadingDot />
    </div>
    <Sender
      :loading="loadingBottom"
      :on-send="welcome ? createConversation : onSend"
    >
      <ElTooltip content="New conversation">
        <Icon
          @click="() => createConversation()"
          :name="welcome ? 'chat' : 'chat1'"
          :class="{ne: welcome}"
          :size="32"
        />
      </ElTooltip>
    </Sender>
  </Scroll>
</template>

<style lang="less" scoped>
.message-root {
  padding: 16px 16px 100px;
  box-sizing: border-box;
  overflow-x: hidden;
  overflow-y: auto;
}
.message {
  max-width: 100%;
  width: fit-content;
  position: relative;
  word-break: break-word;
  background: var(--message-background-color);
  box-sizing: border-box;
  border-radius: 16px;
  border-bottom-left-radius: 0;
  padding: 16px;
  &.user {
    background: #17b09e;
    color: var(--reverse-text-color);
    border-bottom-left-radius: 16px;
    border-bottom-right-radius: 0;
    word-wrap: break-word;
    align-self: flex-end;
  }
  &.error {
    color: #e44646;
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
.suggestion {
  border: 1px solid var(--secondary-border-color);
  color: var(--primary-label-color);
  &:hover {
    background: var(--tertiary-background-color);
    color: var(--primary-text-color);
  }
}
</style>
