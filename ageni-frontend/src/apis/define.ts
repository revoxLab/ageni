import axios, {AxiosRequestConfig} from 'axios'
import {isFunction} from 'lodash-es'
import {ref} from 'vue'
import {logout, sleep, storage} from '../utils'

const debugging = false

const request = axios.create({
  withCredentials: true,
  baseURL: import.meta.env.VITE_PROXY,
})

const skipToken = [
  '/web/wallet_login',
  '/v1/studio/bot_detail',
  '/v1/studio/bot_list',
  '/v1/studio/bot_tabs',
]

async function fetchData<Data = null>(
  props: AxiosRequestConfig<Data>,
  mockData: Computable<Data, any>
) {
  const {method, params, url} = props
  const skipAuth = skipToken.includes(url!)
  const token = storage.get('token') ?? null

  if (!skipAuth && !token) {
    return
  }

  return new Promise<Data>((resolve, reject) => {
    if (debugging) {
      setTimeout(() => {
        const data = isFunction(mockData) ? mockData(params) : mockData
        console.log(`require for ${url}`, params, data)
        resolve(data)
      }, 500)
    } else {
      request({
        ...props,
        data: method === 'post' ? params : undefined,
        params: method === 'get' ? params : undefined,
        headers: {authorization: token},
      })
        .then((res) => {
          if (res.data.code === 20007) {
            logout()
          } else {
            resolve(res.data)
          }
        })
        .catch(reject)
    }
  })
}

type DefineProps<Params, Data> = {
  method?: 'get' | 'post'
  mockData?: Computable<Data, Params>
  mockParams?: Params
  url: string
}

export function defineMutation<Params, Data>({
  method = 'post',
  mockData,
  url,
}: DefineProps<Params, Data>) {
  return (props?: {params?: Params}) =>
    fetchData({...props, url, method}, mockData)
}

export function defineQuery<Params, Data>({
  method = 'get',
  mockData,
  url,
}: DefineProps<Params, Data>) {
  return (props?: {params?: Params; skip?: boolean; interval?: number}) => {
    const error = ref()
    const timeout = ref()
    const data = ref<Data>()
    const loading = ref(false)
    const {params: initialParams, skip, interval} = props ?? {}
    const refetch = async (
      params?: Partial<Params>
    ): Promise<Data | undefined> => {
      if (interval) {
        clearTimeout(timeout.value)
        timeout.value = setTimeout(refetch, interval)
      }
      if (loading.value) {
        return
      }
      await sleep(0)
      loading.value = true
      return fetchData(
        {params: {...initialParams, ...params}, url, method},
        mockData
      )
        .finally(() => (loading.value = false))
        .then((res) => (data.value = res))
        .catch((err) => (error.value = err))
    }
    if (skip !== true) {
      refetch()
    }
    return {data, error, loading, refetch}
  }
}

export function defineApi<Params, Data>({
  method = 'post',
  ...rest
}: DefineProps<Params, Data>) {
  return {
    useRequest: defineQuery({method, ...rest}),
    request: defineMutation({method, ...rest}),
  }
}
