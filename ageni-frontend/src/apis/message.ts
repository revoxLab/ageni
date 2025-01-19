import {processToolCalls} from '@/contract/studio'
import {sleep} from '@/utils'
import {defineMutation} from './define'
import {Response, StudioMessage, StudioMessageInput, ToolResult} from './type'

const studioMessageQuery = defineMutation<
  StudioMessageInput,
  Response<StudioMessage>
>({
  method: 'post',
  url: '/v1/studio/get_message',
  mockData: {
    code: 0,
    data: {
      id: 111,
      conversation_id: 111,
      message_type: 'assistant',
      content: 'xxx',
      status: 1,
      next_message_ids: [5781, 5782, 5783],
      created_at: {
        seconds: 121212,
      },
    },
    message: 'process success',
  },
})

export const sendStudioMessageMutation = defineMutation({
  url: '/v1/studio/send_message',
  mockParams: {
    conversation_id: 1,
    content: '1212',
    tool_results: [
      {tool_call_id: 'call_Jsjd9vo9TfKZNkRmRZFjBrbi', result: 'result'},
    ] as ToolResult[] | undefined,
  },
  mockData: {
    code: 0,
    data: {
      resp_message_id: 10,
    },
    message: 'process success',
  },
})

export async function getStudioMessage(
  conversationId: number,
  params: StudioMessageInput,
  onSuccess: (data: StudioMessage) => void,
  onFailed: (data?: Response<unknown>) => void
) {
  return new Promise<void>(async (resolve) => {
    const getMessage = async () => {
      const res = await studioMessageQuery({params})
      if (res?.code === 100003) {
        await sleep(1000)
        return getMessage()
      } else if (res?.code === 0) {
        onSuccess(res.data)
        const results = []
        const toolCalls = res.data.tool_calls ?? []
        if (toolCalls.length > 0) {
          for (let i = 0; i < toolCalls.length; i++) {
            const result = await processToolCalls(toolCalls[i])
            if (result?.tool_call_id) {
              results.push(result)
            }
          }
          const res = await sendStudioMessageMutation({
            params: {
              content: '',
              conversation_id: conversationId,
              tool_results: results,
            },
          })
          if (res?.data) {
            await getStudioMessage(
              conversationId,
              {message_id: res.data.resp_message_id},
              onSuccess,
              onFailed
            )
          }
        } else if (res.data.next_message_ids?.length) {
          for (let id of res.data.next_message_ids) {
            await getStudioMessage(
              conversationId,
              {message_id: id},
              onSuccess,
              onFailed
            )
          }
        }
      } else {
        onFailed(res)
      }
    }
    await getMessage()
    resolve()
  })
}

export function mockUserMessage(query: string): StudioMessage {
  return {
    created_at: {seconds: Math.round(Date.now() / 1000)},
    conversation_id: 1111,
    message_type: 'user',
    content: query,
    status: 1,
    id: 0,
  }
}

export function mockErrorMessage(): StudioMessage {
  return {
    created_at: {seconds: Math.round(Date.now() / 1000)},
    message_type: 'assistant',
    conversation_id: 1111,
    content: 'error',
    status: 3,
    id: 0,
  }
}
