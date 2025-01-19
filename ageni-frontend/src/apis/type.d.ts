/// <reference types="vite/client" />

import {
  useAgentDetailQuery,
  useAgentListQuery,
  usePluginListQuery,
  useUserAgentList,
} from '.'

export type Response<T> = {
  data: T
  code: number
  message: string
}

export type StudioMessageInput = {message_id: number}

export type StudioMessage<T = {}> = T & {
  id: number
  conversation_id: number
  message_type: 'assistant' | 'user' | 'tool' | 'system'
  created_at: {seconds: number}
  content: string
  status: number
  next_message_ids?: number[]
  tool_calls?: [
    {
      id: string
      type: string
      function: {
        name: string
        arguments: string
      }
      plugin_id?: number
    },
  ]
}

export type ToolResult = {
  tool_call_id: string
  result: string
}

export type BotDetailShape = NonNullable<
  ReturnType<typeof useAgentDetailQuery>['data']['value']
>['data']

export type PluginShape = NonNullable<
  ReturnType<typeof usePluginListQuery>['data']['value']
>['data']['list'][number]

export type PluginMethodShape = NonNullable<PluginShape['methods']>[number]

export type AgentShape = NonNullable<
  ReturnType<typeof useAgentListQuery>['data']['value']
>['data']['bots'][number]

export type UserAgentShape = NonNullable<
  ReturnType<typeof useUserAgentList>['data']['value']
>['data']['bots'][number]

export type ModelSetting = {
  model: string
  temperature: number
  top_p: number
  rounds: number
  max_length: number
}
