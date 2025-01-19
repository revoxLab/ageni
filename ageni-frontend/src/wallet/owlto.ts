import {ethers} from 'ethers'
import * as owlto from 'owlto-sdk'
import {getChainId, getPureWalletProvider} from './basics'
import {FunctionCallError} from './error'

export async function crossChain(
  tokenName: string,
  fromChainName: string,
  toChainName: string,
  addressFrom: string,
  addressTo: string,
  uiValue: number,
  valueIncludeGasFee?: boolean
) {
  const {allPairInfos, supportedFromChainNames, supportedToChainNames} =
    await getAllPair()
  if (!supportedFromChainNames.has(fromChainName)) {
    throw new FunctionCallError(`${fromChainName} is not supported`)
  }
  if (!supportedToChainNames.has(toChainName)) {
    throw new FunctionCallError(`${toChainName} is not supported`)
  }
  const chainId = await getChainId()
  if (chainId === null) {
    throw new FunctionCallError('current chain is null')
  }

  for (const pairInfo of allPairInfos) {
    if (pairInfo.fromChainName.toLowerCase() === fromChainName.toLowerCase()) {
      if (pairInfo.fromChainId !== `${chainId}`) {
        throw new FunctionCallError(
          `current chain isn't ${fromChainName},please switch chain first`
        )
      }
      break
    }
  }

  const options: owlto.BridgeOptions = {}
  const bridge = new owlto.Bridge(options)
  const result = await bridge.getBuildTx(
    tokenName, //token name
    fromChainName, //from chain name
    toChainName, // to chain name
    `${uiValue}`, // value
    addressFrom, // from address
    addressTo, // to address
    valueIncludeGasFee
  )
  const provider = new ethers.BrowserProvider(await getPureWalletProvider())
  // type of the from chain, only ethereum is supported
  if (result.networkType != owlto.NetworkType.NetworkTypeEthereum) {
    return
  }
  const signer = await provider.getSigner()
  // if need approve, Send approve transaction first
  if (result.txs.approveBody) {
    const tx = await signer.sendTransaction(
      result.txs.approveBody as ethers.TransactionRequest
    )
    await tx.wait() // Wait for the transaction to be mined
  }
  // Send the transfer transaction
  const tx = await signer.sendTransaction(
    result.txs.transferBody as ethers.TransactionRequest
  )
  await tx.wait() // Wait for the transaction to be mined
  try {
    const receipt = await bridge.waitReceipt(tx.hash)
    return {
      ...receipt,
      fromChainHash: tx.hash,
    }
  } catch (error) {}
  return {
    fromChainHash: tx.hash,
  }
}

async function getAllPair(category?: string, valueIncludeGasFee?: boolean) {
  const options: owlto.BridgeOptions = {}
  const bridge = new owlto.Bridge(options)
  const pairInfos = (await bridge.getAllPairInfos(category, valueIncludeGasFee))
    .pairInfos
  const supportedFromChainIds = new Set()
  const supportedToChainIds = new Set()
  const supportedFromChainNames = new Set()
  const supportedToChainNames = new Set()
  for (const pairInfo of pairInfos) {
    supportedFromChainIds.add(pairInfo.fromChainId)
    supportedFromChainNames.add(pairInfo.fromChainName)
    supportedToChainIds.add(pairInfo.toChainId)
    supportedToChainNames.add(pairInfo.toChainName)
  }
  return {
    allPairInfos: pairInfos,
    supportedFromChainIds: supportedFromChainIds,
    supportedFromChainNames: supportedFromChainNames,
    supportedToChainIds: supportedToChainIds,
    supportedToChainNames: supportedToChainNames,
  }
}

export async function getAllPairInfos(
  category?: string,
  valueIncludeGasFee?: boolean
) {
  const {allPairInfos} = await getAllPair(category, valueIncludeGasFee)
  return allPairInfos
}

export async function listSupportedChains(
  category?: string,
  valueIncludeGasFee?: boolean
) {
  const {allPairInfos} = await getAllPair(category, valueIncludeGasFee)
  const supportedChains = new Set()
  for (const pairInfo of allPairInfos) {
    supportedChains.add(pairInfo.fromChainName)
  }
  return Array.from(supportedChains)
}

export async function listSupportedTokens(
  fromChainName: string,
  toChainName: string
) {
  const {allPairInfos, supportedFromChainNames, supportedToChainNames} =
    await getAllPair()
  if (!supportedFromChainNames.has(fromChainName)) {
    throw new FunctionCallError(`${fromChainName} is not supported`)
  }
  if (!supportedToChainNames.has(toChainName)) {
    throw new FunctionCallError(`${toChainName} is not supported`)
  }
  const tokens = new Set()
  for (const pairInfo of allPairInfos) {
    if (
      pairInfo.fromChainName.toLowerCase() === fromChainName.toLowerCase() &&
      pairInfo.toChainName.toLowerCase() === toChainName.toLowerCase()
    ) {
      tokens.add(pairInfo.tokenName)
    }
  }
  return Array.from(tokens)
}

export async function getAllPairInfosForCurrentChain(
  category?: string,
  valueIncludeGasFee?: boolean
) {
  const {allPairInfos} = await getAllPair(category, valueIncludeGasFee)
  const pairInfos = []
  const chainId = await getChainId()
  if (chainId == null) {
    return []
  }
  for (const pairInfo of allPairInfos) {
    if (pairInfo?.fromChainId.toLowerCase() === `${chainId}`) {
      pairInfos.push(pairInfo)
    }
  }
  return pairInfos
}

export async function getPairInfo(
  tokenName: string,
  fromChainName: string,
  toChainName: string,
  valueIncludeGasFee?: boolean
) {
  const {supportedFromChainNames, supportedToChainNames} = await getAllPair()
  if (!supportedFromChainNames.has(fromChainName)) {
    throw new FunctionCallError(`${fromChainName} is not supported`)
  }
  if (!supportedToChainNames.has(toChainName)) {
    throw new FunctionCallError(`${toChainName} is not supported`)
  }
  const options: owlto.BridgeOptions = {}
  const bridge = new owlto.Bridge(options)
  const pairInfo = await bridge.getPairInfo(
    tokenName,
    fromChainName,
    toChainName,
    valueIncludeGasFee
  )
  return pairInfo
}

export async function getFeeInfo(
  tokenName: string,
  fromChainName: string,
  toChainName: string,
  uiValue: number
) {
  const {supportedFromChainNames, supportedToChainNames} = await getAllPair()
  if (!supportedFromChainNames.has(fromChainName)) {
    throw new FunctionCallError(`${fromChainName} is not supported`)
  }
  if (!supportedToChainNames.has(toChainName)) {
    throw new FunctionCallError(`${toChainName} is not supported`)
  }
  const options: owlto.BridgeOptions = {}
  const bridge = new owlto.Bridge(options)
  const pairInfo = await bridge.getFeeInfo(
    tokenName,
    fromChainName,
    toChainName,
    uiValue
  )
  return pairInfo
}

export async function getReceipt(fromChainHash: string) {
  const options: owlto.BridgeOptions = {}
  const bridge = new owlto.Bridge(options)
  const receipt = await bridge.getReceipt(fromChainHash)
  return receipt
}
