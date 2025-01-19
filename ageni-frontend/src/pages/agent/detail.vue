<script lang="ts" setup>
import {
  createStudioConversationMutation,
  deleteConversationMutation,
  useAgentDetailQuery,
  useConversationListQuery,
} from '@/apis'
import {CustomButton, Icon, Image, Scroll, showToast} from '@/components'
import {copyToClipboard, sleep, withLogin} from '@/utils'
import {format} from 'date-fns'
import {ElTooltip} from 'element-plus'
import {ref, watch} from 'vue'
import Message from './message.vue'

const page = ref(1)
const loading = ref(false)
const expendInfo = ref(true)
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
const conversations = ref<{id: number; title: string}[]>([])

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
  if (loading.value) {
    return
  }
  loading.value = true
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
  loading.value = false
})

async function deleteConversation(id: number) {
  await deleteConversationMutation({params: {id}})
  if (id === selectConversationId.value) {
    selectConversationId.value = undefined
  }
  refetchHistory({page: (page.value = 1)})
}

function share() {
  copyToClipboard(window.location.href)
  showToast('The link has been copied. Share it with your friends!', 'check')
}
</script>

<template>
  <div v-if="!detail?.data" class="m24">
    <ElSkeleton :rows="5" animated />
  </div>
  <div v-else class="fbv wh100p">
    <div class="p20 fbh fbjsb header">
      <div class="fbh fbac g10">
        <Icon :name="detail.data.image ?? 'agent'" :size="32" class="br6" />
        <div class="fbv fbjsb">
          <p class="f15 bold">{{ detail.data.name }}</p>
          <div class="fbh fbac g4 f12 label2">
            Published by
            <Icon
              class="br100"
              :name="detail.data.creator.head_pic"
              :size="14"
            ></Icon>
            <p class="label1">{{ detail.data.creator.name }}</p>
            on
            {{ format(detail.data.created_at.seconds * 1000, 'yyyy.MM.dd') }}
          </div>
        </div>
      </div>
      <CustomButton
        @click="share"
        v-slot="{hovered}"
        default-stroke="var(--tertiary-border-color)"
        default-color="var(--primary-label-color)"
        active-fill="var(--primary-text-color)"
        active-color="var(--reverse-text-color)"
        :height="32"
        :width="90"
      >
        <div class="fbh fbac g6 f14">
          <Icon :name="hovered ? 'share1' : 'share'" :size="14" />
          <p>Share</p>
        </div>
      </CustomButton>
    </div>
    <div class="container fbh pr">
      <div ref="contentRef" class="fb1 left-card">
        <Message
          v-if="!selectConversationId"
          :conversation-id="Infinity"
          :create-conversation="createConversation"
          :bot="detail.data"
          welcome
        />
        <template v-for="item in conversations">
          <Message
            ref="messageRef"
            v-if="selectConversationId === item.id"
            :conversation-id="item.id"
            :create-conversation="createConversation"
            :bot="detail.data"
          />
        </template>
      </div>
      <div v-if="expendInfo" class="fbv oa right-content">
        <div class="p24 fbv g20">
          <div class="fbh fbac fbjsb">
            <div class="fbh fbac f16 g4">
              <Icon name="user" :size="16" />
              <p class="pr20">{{ detail.data.users }}</p>
              <Icon name="comment" :size="16" />
              <p>{{ detail.data.conversations }}</p>
            </div>
            <div class="p12 hide-icon center hand" @click="expendInfo = false">
              <Icon name="go1" :size="12" />
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
                <Image
                  :name="item.image"
                  :width="24"
                  :height="24"
                  class="br4"
                />
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
          class="fb1 fbv g10 p24 oa"
          @on-attach-bottom="refetchHistory({page: ++page})"
        >
          <p class="f16 bold">History</p>
          <div class="fbv g12">
            <div
              v-for="item in conversations"
              class="history br8 hand pr"
              @click="selectConversationId = item.id"
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
      <div v-else class="expend-icon center pa">
        <Icon
          name="go1"
          :size="12"
          class="py12 px16"
          style="rotate: 180deg"
          @click="expendInfo = true"
        />
      </div>
    </div>
  </div>
</template>

<style lang="less" scoped>
.container {
  height: calc(100% - 76px);
}
.header {
  background: var(--secondary-background-color);
  border-bottom: 1px solid var(--secondary-border-color);
}
.right-content {
  width: 280px;
  border-left: 1px solid var(--secondary-border-color);
}
.left-card {
  overflow: hidden;
  box-sizing: border-box;
}
.expend-icon {
  border-top-left-radius: 50%;
  border-bottom-left-radius: 50%;
  background: var(--tertiary-background-color);
  top: 28px;
  right: 0;
}
.hide-icon {
  border-radius: 100%;
  &:hover {
    background: var(--tertiary-background-color);
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
    .more {
      display: block;
    }
  }
}
.more {
  height: fit-content;
  display: none;
  right: 0;
  top: 0;
}
</style>
