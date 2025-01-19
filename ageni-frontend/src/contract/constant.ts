import {ERC20Token} from '@pancakeswap/sdk'
import {
  arbitrumTokens,
  bscTokens,
  ethereumTokens,
  lineaTokens,
  opBnbTokens,
  polygonZkEvmTokens,
} from '@pancakeswap/tokens'
import {arbitrum, bsc, linea, mainnet, opBNB, polygonZkEvm} from 'viem/chains'
import arbitumJSON from './token/arbitrum.json'
import bscJSON from './token/bsc.json'
import ethJSON from './token/eth.json'
import lineaJSON from './token/linea.json'
import opbnbJSON from './token/opbnb.json'
import zkevmJSON from './token/zkevm.json'

export const supportChainOnPancake = [
  bsc,
  opBNB,
  mainnet,
  linea,
  polygonZkEvm,
  arbitrum,
]

function getIntersectionToken(
  tokensFromSDK: ERC20Token[],
  tokensFromJSON: {address: string; logoURI: string}[]
) {
  return tokensFromJSON
    .map(({address, logoURI}) => {
      const target = tokensFromSDK.find((item) => item.address === address)
      return target! && Object.assign(target, {logoURI})
    })
    .filter(Boolean)
}

export const supportTokenOnPancake: Record<
  number,
  ReturnType<typeof getIntersectionToken>
> = {
  [bsc.id]: getIntersectionToken(Object.values(bscTokens), bscJSON.tokens),
  [opBNB.id]: getIntersectionToken(
    Object.values(opBnbTokens),
    opbnbJSON.tokens
  ),
  [mainnet.id]: getIntersectionToken(
    Object.values(ethereumTokens),
    ethJSON.tokens
  ),
  [linea.id]: getIntersectionToken(
    Object.values(lineaTokens),
    lineaJSON.tokens
  ),
  [polygonZkEvm.id]: getIntersectionToken(
    Object.values(polygonZkEvmTokens),
    zkevmJSON.tokens
  ),
  [arbitrum.id]: getIntersectionToken(
    Object.values(arbitrumTokens),
    arbitumJSON.tokens
  ),
}
