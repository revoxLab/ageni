import {supportChainOnPancake, supportTokenOnPancake} from '@/contract/constant'
import {
  ChainId,
  CurrencyAmount,
  Native,
  Percent,
  Token,
  TradeType,
} from '@pancakeswap/sdk'
import {
  OnChainProvider,
  SMART_ROUTER_ADDRESSES,
  SmartRouter,
  SmartRouterTrade,
  SwapRouter,
} from '@pancakeswap/smart-router'
import {ethers} from 'ethers'
import {GraphQLClient} from 'graphql-request'
import {
  Address,
  createPublicClient,
  EstimateGasParameters,
  hexToBigInt,
  http,
} from 'viem'
import {
  calculateGasMargin,
  currentWallet,
  fromReadableAmount,
  getChainId,
  getPureWalletProvider,
} from './basics'

function getPublicClient(chainId: number) {
  const chain = supportChainOnPancake.find(({id}) => id === chainId)
  if (chain) {
    return createPublicClient({
      chain,
      transport: http(chain.rpcUrls.default.http[0]),
      batch: {
        multicall: {
          batchSize: 1024 * 200,
        },
      },
    })
  }
}

const v3SubgraphClient = (chainId: number) => {
  const support = supportChainOnPancake.some(({id}) => id === chainId)
  if (support) {
    let url =
      'https://thegraph.com/hosted-service/subgraph/pancakeswap/exchange-v3-bsc'
    if (chainId === 56) {
      url =
        'https://thegraph.com/hosted-service/subgraph/pancakeswap/exchange-v3-bsc'
    } else if (chainId === 204) {
      url =
        'https://opbnb-mainnet-graph.nodereal.io/subgraphs/name/pancakeswap/exchange-v3'
    } else if (chainId === 8453) {
      url =
        'https://api.studio.thegraph.com/query/45376/exchange-v3-base/version/latest'
    } else if (chainId === 59144) {
      url =
        'https://graph-query.linea.build/subgraphs/name/pancakeswap/exchange-v3-linea'
    } else if (chainId === 324) {
      url =
        'https://api.studio.thegraph.com/query/45376/exchange-v3-zksync/version/latest'
    } else if (chainId === 1) {
      url =
        'https://thegraph.com/hosted-service/subgraph/pancakeswap/exchange-v3-eth'
    }
    return new GraphQLClient(url)
  }
  return null
}

const v2SubgraphClient = (chainId: number) => {
  const support = supportChainOnPancake.some(({id}) => id === chainId)
  if (support) {
    let url = 'https://nodereal.io/meganode/api-marketplace/pancakeswap-graphql'
    if (chainId === 56) {
      url = 'https://nodereal.io/meganode/api-marketplace/pancakeswap-graphql'
    } else if (chainId === 204) {
      url =
        'https://opbnb-mainnet-graph.nodereal.io/subgraphs/name/pancakeswap/exchange-v2'
    } else if (chainId === 8453) {
      url =
        'https://api.studio.thegraph.com/query/45376/exchange-v2-base/version/latest'
    } else if (chainId === 59144) {
      url =
        'https://graph-query.linea.build/subgraphs/name/pancakeswap/exhange-v2'
    } else if (chainId === 324) {
      url =
        'https://api.studio.thegraph.com/query/45376/exchange-v2-zksync/version/latest'
    } else if (chainId === 1) {
      url = 'https://api.thegraph.com/subgraphs/name/pancakeswap/exhange-eth'
    }
    return new GraphQLClient(url)
  }
  return null
}

export async function getQuotes(
  dex: string,
  amountIn: number,
  tokenFrom: string | null,
  tokenTo: string | null
) {
  const address = currentWallet()
  if (address == null) return
  const chainId = await getChainId()
  if (chainId == null) {
    throw new Error('current chain is null')
  }
  if (dex.toLocaleLowerCase() === 'pancakeswap.finance'.toLowerCase()) {
    const support = supportChainOnPancake.some((item) => item.id === chainId)
    if (support) {
      const from = toTokenOnPancake(chainId, tokenFrom)
      const to = toTokenOnPancake(chainId, tokenTo)
      console.log('start quotes,params is ' + amountIn, from, to)
      const outputAmount = await getQuotesFromPancake(amountIn, from, to)
      return outputAmount
    } else {
      throw new Error('current chain is not supported!!!')
    }
  } else {
    throw new Error('dex is not support')
  }
}

function findTokenByAddress(chainId: number, address: string) {
  const tokens = supportTokenOnPancake[chainId]
  return (
    tokens.find(
      (item) => item.address.toLowerCase() === address.toLowerCase()
    ) ?? null
  )
}

function toTokenOnPancake(chainId: number, token: string | null): Token | null {
  if (
    token == null ||
    token.trim() === '' ||
    token.trim() === '0x0000000000000000000000000000000000000000'
  ) {
    return null
  } else {
    return findTokenByAddress(chainId, token)
  }
}

export async function swapTokens(
  dex: string,
  amountIn: number,
  tokenFrom: string | null,
  tokenTo: string | null,
  slippageTolerance: number | null
) {
  const address = currentWallet()
  const chainId = await getChainId()
  if (chainId == null) {
    throw new Error('current chain is null')
  }
  if (dex.toLocaleLowerCase() === 'pancakeswap.finance'.toLowerCase()) {
    const support = supportChainOnPancake.some((item) => item.id === chainId)
    if (support) {
      const from = toTokenOnPancake(chainId, tokenFrom)
      const to = toTokenOnPancake(chainId, tokenTo)
      console.log(
        'start swap,params is ',
        amountIn,
        from,
        to,
        slippageTolerance
      )
      return await swapTokensFromPancake(
        address!,
        amountIn,
        from,
        to,
        slippageTolerance
      )
    } else {
      throw new Error('current chain is not supported!!!')
    }
  } else {
    throw new Error('des is not supported!!!')
  }
}

export async function getQuotesFromPancake(
  amountIn: number,
  tokenFrom: Token | null,
  tokenTo: Token | null
): Promise<string | undefined> {
  try {
    if (tokenFrom?.address === tokenTo?.address) {
      return amountIn.toString()
    }
    const trade = await createTrade(amountIn, tokenFrom, tokenTo)
    return trade?.outputAmount?.toExact()
  } catch (e) {
    console.error(e)
    return 'failed'
  }
}

async function swapTokensFromPancake(
  address: string,
  amountIn: number,
  tokenFrom: Token | null,
  tokenTo: Token | null,
  slippageTolerance: number | null
): Promise<string | null> {
  console.log('start swap token!!!')
  const trade = await createTrade(amountIn, tokenFrom, tokenTo)
  if (trade == null) {
    throw new Error('trade create failed,while swap token!!!')
  }
  console.log(
    'final price is ',
    trade.outputAmount.toExact(),
    ' and slippage tolerance is ',
    new Percent(
      slippageTolerance == null ? 1 : slippageTolerance * 10000,
      slippageTolerance == null ? 100 : 100 * 10000
    ).toFixed()
  )
  const chainId = await getChainId()
  const web3Provider = new ethers.BrowserProvider(await getPureWalletProvider())
  const {value, calldata} = SwapRouter.swapCallParameters(trade!, {
    recipient: address as Address,
    slippageTolerance: new Percent(
      slippageTolerance == null ? 1 : slippageTolerance * 10000,
      slippageTolerance == null ? 100 : 100 * 10000
    ),
  })
  const tx = {
    account: address,
    to: SMART_ROUTER_ADDRESSES[chainId as ChainId],
    data: calldata,
    value: hexToBigInt(value),
  }
  const gasEstimate = await getPublicClient(chainId!)?.estimateGas(
    tx as EstimateGasParameters
  )
  const signer = await web3Provider.getSigner()
  const hash = await signer.sendTransaction({
    chainId,
    to: SMART_ROUTER_ADDRESSES[chainId as ChainId],
    data: calldata,
    value: hexToBigInt(value),
    gasLimit: calculateGasMargin(gasEstimate!),
  })
  console.log('swap waiting...')
  await hash.wait()
  console.log('swap finish,hash:' + hash.hash)
  return hash.hash
}

async function createTrade(
  amountIn: number,
  swapFrom: Token | null,
  swapTo: Token | null
): Promise<SmartRouterTrade<TradeType> | null> {
  const chainId = await getChainId()
  if (chainId == null) {
    return null
  }
  const tokenFrom = swapFrom ?? Native.onChain(chainId as ChainId)
  const tokenTo = swapTo ?? Native.onChain(chainId as ChainId)
  const amount = CurrencyAmount.fromRawAmount(
    tokenFrom,
    fromReadableAmount(amountIn, 18)
  )
  const currentChainId = chainId
  const [v2Pools, v3Pools] = await Promise.all([
    SmartRouter.getV2CandidatePools({
      onChainProvider: (({chainId}) =>
        getPublicClient(chainId ?? currentChainId)) as OnChainProvider,
      v2SubgraphProvider: ({chainId}) =>
        v2SubgraphClient(chainId ?? currentChainId)!,
      v3SubgraphProvider: ({chainId}) =>
        v3SubgraphClient(chainId ?? currentChainId!)!,
      currencyA: amount.currency,
      currencyB: tokenTo,
    }),
    SmartRouter.getV3CandidatePools({
      onChainProvider: (({chainId}) =>
        getPublicClient(chainId ?? currentChainId)) as OnChainProvider,
      subgraphProvider: ({chainId}) =>
        v3SubgraphClient(chainId ?? currentChainId)!,
      currencyA: amount.currency,
      currencyB: tokenTo,
      subgraphFallback: false,
    }),
  ])
  const pools = [...v2Pools, ...v3Pools]
  const trade = await SmartRouter.getBestTrade(
    amount,
    tokenTo,
    TradeType.EXACT_INPUT,
    {
      gasPriceWei: () => getPublicClient(chainId)?.getGasPrice(),
      maxHops: 2,
      maxSplits: 2,
      poolProvider: SmartRouter.createStaticPoolProvider(pools),
      quoteProvider,
      quoterOptimization: true,
    }
  )
  console.log('createTrade successed!!!')
  return trade
}
const quoteProvider = SmartRouter.createQuoteProvider({
  onChainProvider: (({chainId}) => {
    getPublicClient(chainId!)
  }) as OnChainProvider,
})

export async function estimateGasFromPancake(
  address: string,
  amountIn: number,
  tokenFrom: string | null,
  tokenTo: string | null,
  slippageTolerance: number | null
): Promise<BigInt | undefined> {
  try {
    const chainId = await getChainId()
    if (chainId == null) return
    const from = toTokenOnPancake(chainId, tokenFrom)
    const to = toTokenOnPancake(chainId, tokenTo)
    const trade = await createTrade(amountIn, from, to)
    if (trade == null) {
      throw new Error('trade create failed,while swap token!!!')
    }
    const {value, calldata} = SwapRouter.swapCallParameters(trade!, {
      recipient: address as Address,
      slippageTolerance: new Percent(
        slippageTolerance == null ? 1 : slippageTolerance * 10000,
        slippageTolerance == null ? 100 : 100 * 10000
      ),
    })
    const tx = {
      account: address,
      to: SMART_ROUTER_ADDRESSES[chainId as ChainId],
      data: calldata,
      value: hexToBigInt(value),
    } as EstimateGasParameters
    return await getPublicClient(chainId!)?.estimateGas(tx)
  } catch {
    return 0n
  }
}
