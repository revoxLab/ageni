import {ChainId, create, supportedChains} from 'kiloex-ts-sdk'
import {Address} from 'viem'
import {currentWallet, getChainId} from './basics'

let currentChainId = ChainId.BSCTEST
let kiloClient: ReturnType<typeof create>

setTimeout(() => {
  try {
    kiloClient = create(currentChainId)
  } catch {}
})

export function supportChainsOnKiloex() {
  return supportChains()
}

function supportChains() {
  const result = []
  const chains = supportedChains()
  for (const chain of chains) {
    result.push({chainName: chain.name, chainId: chain.id})
  }
  return result
}

export async function supportedProductsOnKiloex(chainId: number) {
  await switchChain()
  const products = await kiloClient.supportedProductsByChain(chainId)
  const res = []
  for (const product of products) {
    res.push({
      name: product.name,
      productId: product.productId,
      isHot: product.isHot,
      base: product.base,
      quote: product.quote,
    })
  }
  return res
}

export async function getProductsInfoOnKiloex(productIds: number[]) {
  await switchChain()
  const allProducts = await kiloClient.getProductsInfo(productIds)
  return JSON.stringify(allProducts)
}

export async function getProductPricesOnKiloex(productIds: number[]) {
  await switchChain()
  const prices = await kiloClient.getProductPrices(productIds)
  return prices
}

async function switchChain() {
  const chainId = await getChainId()
  if (chainId == null) return
  if (currentChainId !== chainId) {
    const chains = supportChains()
    for (const chain of chains) {
      if (chain.chainId === chainId) {
        kiloClient = create(chainId)
        currentChainId = chainId
        break
      }
    }
  }
}

export async function increasePositionOnKiloex({
  openType,
  ...position
}: {
  openType: number
  productId: number
  leverage: number
  isLong: boolean
  margin: number
  point: number
  tickerPrice: string
  stopLessPrice?: number
  takeProfitPrice?: number
  entryPrice?: number
}) {
  const address = currentWallet()
  await switchChain()
  console.log('increasePosition', {
    address,
    position,
    openType,
  })
  const res = await kiloClient.increasePosition(
    address as Address,
    openType,
    position
  )
  return {txHash: res.transactionHash}
}

export async function closePositionOnKiloex(
  productId: number,
  margin: number,
  isLong: boolean,
  tickerPrice: number
) {
  await switchChain()
  const address = currentWallet()
  console.log('closePosition', {
    address,
    position: {
      tickerPrice,
      productId,
      isLong,
      margin,
    },
  })
  const res = await kiloClient.closePosition(address as Address, {
    tickerPrice: tickerPrice,
    productId: productId,
    isLong: isLong,
    margin: margin,
  })
  return {txHash: res.transactionHash}
}

export async function updatePositionMarginOnKiloex(
  productId: number,
  margin: number,
  isLong: boolean
) {
  await switchChain()
  const address = currentWallet()
  console.log('updatePositionMargin', {
    address,
    position: {
      productId,
      isLong,
      margin,
    },
  })
  const res = await kiloClient.updatePositionMargin(address as Address, {
    productId: productId,
    isLong: isLong,
    margin: margin,
  })
  return {txHash: res.transactionHash}
}

export async function cancelOrderOnKiloex(
  type: 'Increase' | 'Decrease',
  orderIndex: number
) {
  await switchChain()
  const address = currentWallet()
  console.log('cancelOrder', {address, type, orderIndex})
  const res = await kiloClient.cancelOrder(address as Address, type, orderIndex)
  return {txHash: res.transactionHash}
}

export async function updateOrderOnKiloex(
  type: 'STOP_LOSS' | 'TAKE_PROFIT',
  orderIndex: number,
  margin: number,
  leverage: number,
  limitPrice: number
) {
  await switchChain()
  const address = currentWallet()
  console.log('updateOrder', {
    address,
    updateData: {
      type,
      orderIndex,
      margin,
      leverage,
      limitPrice,
    },
  })
  const res = await kiloClient.updateOrder(address as Address, {
    type,
    orderIndex,
    margin,
    leverage,
    limitPrice,
  })
  return {txHash: res.transactionHash}
}

export async function getAllPositionsOnKiloex() {
  await switchChain()
  const address = currentWallet()
  const allPositions = await kiloClient.getAllPositions(address as Address)
  const res = []
  for (const position of allPositions) {
    res.push({
      borrowing: position.borrowing,
      funding: position.funding,
      isLong: position.isLong,
      leverage: position.leverage,
      margin: position.margin,
      oraclePrice: position.oraclePrice,
      price: position.price,
      productId: position.productId,
      size: position.size,
      symbol: position.symbol,
      timestamp: position.timestamp,
    })
  }
  console.log('getAllPositions', res)
  return res
}

export async function getAllOrdersOnKiloex() {
  await switchChain()
  const address = currentWallet()
  const res = await kiloClient.getAllOrders(address as Address)
  console.log(res)
  return JSON.stringify(res)
}

export async function getTradesHistoryOnKiloex() {
  const address = currentWallet()
  await switchChain()
  const res = await kiloClient.getTradesHistory(address as Address)
  return JSON.stringify(res)
}

export async function setApproveOnKiloex(spender: string) {
  const address = currentWallet()
  await switchChain()
  const res = await kiloClient.setApprove(
    address as Address,
    spender as Address
  )
  return res
}
