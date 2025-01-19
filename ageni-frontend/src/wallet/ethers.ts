import {walletLoginMutation} from '@/apis'
import {showToast} from '@/components'
import {storage} from '@/utils'
import {createWeb3Modal, defaultConfig} from '@web3modal/ethers/vue'
import {BrowserProvider} from 'ethers'
import {debounce} from 'lodash-es'
import {disconnectWallet, getPureWalletProvider} from './basics'

const projectId = import.meta.env.VITE_WALLET_CONNECT_KEY
const eth = {
  chainId: 1,
  name: 'eth Mainnet',
  currency: 'eth',
  explorerUrl: '',
  rpcUrl: '',
}
const polygon = {
  chainId: 137,
  name: 'polygon Mainnet',
  currency: 'polygon',
  explorerUrl: 'https://polygonscan.com',
  rpcUrl: 'https://polygon-rpc.com/',
}
const xLayer = {
  chainId: 196,
  name: 'X Layer Mainnet',
  currency: 'OKB',
  explorerUrl: 'https://www.okx.com/explorer/xlayer',
  rpcUrl: 'https://xlayerrpc.okx.com',
}
const linear = {
  chainId: 59144,
  name: 'Linea',
  currency: 'ETH',
  explorerUrl: 'https://lineascan.build',
  rpcUrl: 'https://rpc.linea.build',
}
const opBNB = {
  chainId: 204,
  name: 'opBNB Mainnet',
  currency: 'BNB',
  explorerUrl: 'https://opbnb.bscscan.com',
  rpcUrl: 'https://opbnb-mainnet-rpc.bnbchain.org',
}
const readon = {
  chainId: 12015,
  name: 'ReadON Loki',
  currency: 'READ',
  explorerUrl: 'https://read.zkevm.lumoz.info',
  rpcUrl:
    'https://pre-alpha-zkrollup-rpc.opside.network/readon-content-test-chain',
}
const bnb = {
  chainId: 56,
  name: 'BNB Smart Chain Mainnet',
  currency: 'BNB',
  explorerUrl: 'https://bscscan.com/',
  rpcUrl: 'https://bsc-dataseed.binance.org/',
}
const taiko = {
  chainId: 167000,
  name: 'Taiko Mainnet',
  currency: 'ETH',
  explorerUrl: 'https://taikoscan.io',
  rpcUrl: 'https://rpc.mainnet.taiko.xyz',
}
const base = {
  chainId: 8453,
  name: 'Base Mainnet',
  currency: 'ETH',
  explorerUrl: 'https://basescan.org/',
  rpcUrl: 'https://mainnet.base.org/',
}
const sei = {
  chainId: 1329,
  name: 'Sei',
  currency: 'SEI',
  explorerUrl: 'https://seitrace.com/?chain=pacific-1',
  rpcUrl: 'https://evm-rpc.sei-apis.com',
}
const metadata = {
  name: 'REVOX Studio',
  description: 'Studio: The Web3 GPT',
  url: 'https://studio.revox.ai', // url must match your domain & subdomain
  icons: [],
}
const arbitrum = {
  chainId: 42161,
  name: 'ArbitrumOneMainnet',
  currency: 'ETH',
  explorerUrl: 'https://arbiscan.io/',
  rpcUrl: 'https://arb1.arbitrum.io/rpc',
}
const optimism = {
  chainId: 10,
  name: 'OptimismMainnet',
  currency: 'ETH',
  explorerUrl: 'https://optimistic.etherscan.io',
  rpcUrl: 'https://mainnet.optimism.io',
}
const zksync = {
  chainId: 324,
  name: 'ZksyncMainnet',
  currency: 'ETH',
  explorerUrl: 'https://explorer.zksync.io/',
  rpcUrl: 'https://mainnet.era.zksync.io',
}
const bscTest = {
  chainId: 97,
  name: 'BSCTestnet',
  currency: 'TBNB',
  explorerUrl: 'https://testnet.bscscan.com',
  rpcUrl: 'https://data-seed-prebsc-1-s1.bnbchain.org:8545',
}
const ethersConfig = defaultConfig({
  metadata,
  enableEIP6963: true,
  enableInjected: true,
  enableCoinbase: true,
})

export const web3Modal = createWeb3Modal({
  projectId,
  ethersConfig,
  chains: [
    eth,
    polygon,
    xLayer,
    linear,
    opBNB,
    readon,
    bnb,
    taiko,
    base,
    sei,
    arbitrum,
    optimism,
    zksync,
    bscTest,
  ],
  featuredWalletIds: [
    '8a0ee50d1f22f6651afcae7eb4253e52a3310b90af5daef78a8c4929a9bb99d4',
    '971e689d0a5be527bac79629b4ee9b925e82208e5168b733496a09c0faed0709',
    'c57ca95b47569778a828d19178114f4db188b89b763c899ba0be274e97267d96',
    '15c8b91ade1a4e58f3ce4e7a0dd7f42b47db0c8df7e0d84f63eb39bcb96c4e0f',
  ],
  allowUnsupportedChain: true,
  enableAnalytics: true,
  enableOnramp: true,
})

web3Modal.subscribeWalletInfo(() => {
  const address = web3Modal.getAddress()
  if (address && !storage.get('address')) {
    evmLogin(address)
  }
})

export const evmLogin = debounce(async (address: string) => {
  const signMessage = `Sign in to REVOX`
  const ethersProvider = new BrowserProvider(await getPureWalletProvider())
  const signer = await ethersProvider.getSigner()
  const signature = await signer.signMessage(signMessage).catch((error) => {
    showToast('Sign Failed')
    disconnectWallet()
    throw error
  })
  const res = await walletLoginMutation({
    params: {
      signature,
      wallet_address: address,
      from: 'studio',
    },
  })

  if (res?.code === 0) {
    storage.set('address', address)
    storage.set('token', res?.data.token)
    window.location.reload()
  } else {
    showToast('Login Failed')
    disconnectWallet()
  }
}, 500)
