import {storage} from '@/utils'
import {SMART_ROUTER_ADDRESSES} from '@pancakeswap/smart-router'
import axios from 'axios'
import {Eip1193Provider, ethers} from 'ethers'
import {Address, Chain} from 'viem'
import {ref} from 'vue'
import {useAccount} from 'wagmi'
import {evmLogin, web3Modal as wallet} from './ethers'

export const account = ref<ReturnType<typeof useAccount>>()
export const openWalletRef = ref<Function>()
export const gasBalanceRef = ref<String>()
export const chainRef = ref<Chain>()

const useWagmiConnector = false

export const openWallet = () => {
  if (useWagmiConnector) {
    openWalletRef.value?.()
  } else {
    wallet.open()
  }
}

export function getChainId() {
  if (useWagmiConnector) {
    return account.value?.connector?.getChainId()
  } else {
    return wallet.getChainId()
  }
}

export function getWalletAddress() {
  if (useWagmiConnector) {
    return account.value?.address
  } else {
    return wallet.getAddress()
  }
}

export function switchNetwork(chainId: number) {
  if (useWagmiConnector) {
    return account.value?.connector?.switchChain?.({chainId})
  } else {
    return wallet.switchNetwork(chainId)
  }
}

export function disconnectWallet() {
  if (useWagmiConnector) {
    return account.value?.connector?.disconnect()
  } else {
    return wallet.disconnect()
  }
}

export function getPureWalletProvider() {
  return new Promise<Eip1193Provider>((resolve) => {
    const run = async () => {
      const provider = useWagmiConnector
        ? await account.value?.connector?.getProvider()
        : wallet.getWalletProvider()
      if (provider) {
        resolve(provider as Eip1193Provider)
      } else {
        console.warn('provider not available, please wait...')
        setTimeout(run, 500)
      }
    }
    run()
  })
}

export async function getProvider() {
  const address = storage.get('address')
  const walletAddress = getWalletAddress()
  /**
   * Re-login if the address is inconsistent with the cached address
   * caseA: Wallet type change
   * caseB: The wallet is the same but the account has changed
   * caseC: Locked or not logged in
   */
  if (address !== walletAddress) {
    walletAddress ? evmLogin(walletAddress) : openWallet()
  } else {
    return getPureWalletProvider()
  }
}

// Chains currently supported by web3model
export function supportChains(): Record<string, number> {
  return {
    opbnbmainnet: 204,
    bnbmainnet: 56,
    ethereummainnet: 1,
    xlayermainnet: 196,
    lineamainnet: 59144,
    loki: 12015,
    taiko: 16700,
    basemainnet: 8453,
    sei: 1329,
    polygonposmainnet: 137,
    optimismmainnet: 10,
    zksyncmainnet: 324,
    arbitrumonemainnet: 42161,
    bsctestnet: 97,
  }
}

// Manually maintained mapping table of chain IDs to chain names
const chainNameMap: Record<Meta, string> = {
  1: 'Ethereum',
  196: 'xLayer',
  59144: 'Linea',
  204: 'opBNB',
  12015: 'Loki',
  56: 'BSC',
  16700: 'taiko',
  8453: 'base',
  1329: 'sei',
  137: 'Polygon',
  10: 'OptimismMainnet',
  324: 'ZksyncMainnet',
  42161: 'ArbitrumOneMainnet',
  97: 'bsctestnet',
}

// Polling function
const pollForAddress = (): Promise<string> => {
  return new Promise((resolve, reject) => {
    const intervalId = window.setInterval(() => {
      try {
        const address = getWalletAddress()
        if (address) {
          clearInterval(intervalId)
          resolve(address)
        } else {
          console.log('Waiting for wallet connection...')
        }
      } catch (error) {
        console.error('Error getting wallet address: ', error)
        clearInterval(intervalId)
        reject(error)
      }
    }, 2000)
  })
}

export async function connectWallet(chainId?: number) {
  useWagmiConnector
    ? await account.value?.connector?.connect({chainId: chainId})
    : await wallet.open()
  return pollForAddress()
}

export async function currentNetwork() {
  const chainId = await getChainId()
  return {chainId, chainName: chainNameMap[chainId!] ?? chainId}
}

export function currentWallet() {
  const address = getWalletAddress()
  if (address == null) {
    throw new Error('address is null')
  }
  return address
}

export function isConnected() {
  if (useWagmiConnector) {
    return account.value?.isConnected
  } else {
    return wallet.getIsConnected()
  }
}

export async function currentConnectState() {
  if (isConnected()) {
    return {
      ...(await currentNetwork()),
      walletAddress: currentWallet(),
    }
  }
}

export async function switchChain(chainId: number) {
  if (chainId in chainNameMap) {
    return switchNetwork(chainId)
  } else {
    throw new Error(`chainId=${chainId} is not supported now`)
  }
}

export async function getBalanceOfGas() {
  if (useWagmiConnector) {
    return gasBalanceRef.value
  } else {
    const web3Provider = new ethers.BrowserProvider(
      await getPureWalletProvider()
    )
    const signer = await web3Provider.getSigner()
    const address = await signer.getAddress()
    const balance = await web3Provider.getBalance(address)
    return ethers.formatEther(balance)
  }
}

export async function getTokenBalanceForErc20(contractAddress?: string) {
  if (contractAddress) {
    return getTokenBalanceOfErc20ForWeb3Modal(contractAddress)
  } else {
    return await getAllTokenBalanceOfErc20()
  }
}

const erc20Abi = [
  'function balanceOf(address) view returns (uint256)',
  'function name() view returns (string)',
  'function symbol() view returns (string)',
]

async function getTokenBalanceOfErc20ForWeb3Modal(contractAddress: string) {
  const walletAddress = currentWallet()
  if (walletAddress == null || walletAddress.length < 1)
    throw new Error('walletAddress is null')
  if (contractAddress.length < 1) {
    return await getAllTokenBalanceOfErc20()
  }
  const tokenInfo = await getTokenBalancesAndInfo(walletAddress, [
    contractAddress,
  ])
  return tokenInfo[contractAddress].balance
}

async function getTokenTransactions(walletAddress: string, apiKey: string) {
  const url = `https://api.bscscan.com/api?module=account&action=tokentx&address=${walletAddress}&startblock=0&endblock=99999999&sort=asc&apikey=${apiKey}`
  try {
    const response = await axios.get(url)
    const data = response.data
    return data?.result ?? []
  } catch (error) {
    throw new Error(`from bscscan get data failed,error is ${error}`)
  }
}

async function getAllTokenBalanceOfErc20() {
  const walletAddress = currentWallet()
  if (walletAddress == null || walletAddress.length < 1)
    throw new Error('walletAddress is null')
  const transactions = await getTokenTransactions(
    walletAddress,
    'N8C48PM8D1DXQNXGWJ85C2HPJT94JG4MC6'
  )
  const tokenAddresses = new Set<string>()
  transactions.forEach((tx: any) => {
    if (tx.to.toLowerCase() === walletAddress.toLowerCase()) {
      tokenAddresses.add(tx.contractAddress)
    }
  })
  const uniqueTokenAddresses = [...tokenAddresses]
  const tokenInfos = await getTokenBalancesAndInfo(
    walletAddress,
    uniqueTokenAddresses
  )
  const result = []
  for (const address of tokenAddresses) {
    result.push(tokenInfos[address])
  }
  return result
}

async function getTokenBalancesAndInfo(
  walletAddress: string,
  tokenAddresses: any
) {
  if (walletAddress == null || walletAddress.length < 1) return {}
  const tokensInfo: AnyObject = {}
  const web3Provider = new ethers.BrowserProvider(await getPureWalletProvider())
  for (const address of tokenAddresses) {
    const tokenContract = new ethers.Contract(address, erc20Abi, web3Provider)
    const balance = await tokenContract.balanceOf(walletAddress)
    const name = await tokenContract.name()
    const symbol = await tokenContract.symbol()

    tokensInfo[address] = {
      balance: ethers.formatUnits(balance, 18),
      name: name,
      symbol: symbol,
      address: address,
    }
  }
  return tokensInfo
}

export async function signMessageForAI(message: string) {
  const web3Provider = new ethers.BrowserProvider(await getPureWalletProvider())
  const signer = await web3Provider.getSigner()
  const signature = await signer.signMessage(message)
  return signature
}

export async function transfer(
  to: Address,
  amount: string,
  contractAddress: string,
  chainId?: number
) {
  if (
    contractAddress == null ||
    contractAddress.trim() === '' ||
    contractAddress.trim() === '0x0000000000000000000000000000000000000000' ||
    contractAddress.trim() === 'native'
  ) {
    return transferGasForWeb3Modal(to, amount, chainId)
  } else {
    return transferErc20ForWeb3Modal(
      contractAddress as Address,
      to as Address,
      amount,
      chainId
    )
  }
}

async function transferGasForWeb3Modal(
  to: string,
  amount: string,
  chainId?: number
) {
  const web3Provider = new ethers.BrowserProvider(await getPureWalletProvider())
  const signer = await web3Provider.getSigner()
  const tx = await signer.sendTransaction({
    to: to,
    value: ethers.parseEther(amount),
    chainId: chainId,
  })
  console.log('Transaction Hash:', tx.hash)
  await tx.wait()
  console.log('Transaction finish:', tx.hash)
  return tx.hash
}

async function transferErc20ForWeb3Modal(
  tokenAddress: Address,
  to: Address,
  amount: string,
  chainId?: number
) {
  const web3Provider = new ethers.BrowserProvider(await getPureWalletProvider())
  const signer = await web3Provider.getSigner()
  const erc20Contract = new ethers.Contract(
    tokenAddress,
    [
      'function transfer(address to, uint256 amount) public returns (bool)',
      'function decimals() view returns (uint8)',
    ],
    signer
  )
  const decimals = await erc20Contract.decimals()
  const amountInTokens = ethers.parseUnits(amount, decimals)
  const tx = await erc20Contract.transfer(to, amountInTokens, {
    chainId: chainId,
  })
  console.log('Transaction Hash:', tx.hash)
  await tx.wait()
  console.log('ERC-20 Token transfer confirmed')
  return tx.hash
}

export const ERC20_ABI_FROM_DEX = [
  // Read-Only Functions
  'function balanceOf(address owner) view returns (uint256)',
  'function decimals() view returns (uint8)',
  'function symbol() view returns (string)',

  // Authenticated Functions
  'function transfer(address to, uint amount) returns (bool)',
  'function approve(address _spender, uint256 _value) returns (bool)',

  // Events
  'event Transfer(address indexed from, address indexed to, uint amount)',
]

export async function getTokenTransferApproval(
  tokenAddress: string,
  amountIn: number
): Promise<string | null> {
  const provider = new ethers.BrowserProvider(await getPureWalletProvider())
  const address = currentWallet()
  if (!address) {
    console.log('No Provider Found')
    return null
  }
  const tokenContract = new ethers.Contract(
    tokenAddress,
    ERC20_ABI_FROM_DEX,
    provider
  )
  const chainId = await getChainId()
  const transaction = await tokenContract.approve.populateTransaction(
    SMART_ROUTER_ADDRESSES[chainId as 1],
    fromReadableAmount(amountIn, 18).toString()
  )
  const signer = await provider.getSigner()
  const result = await signer.sendTransaction({...transaction, from: address})
  await result.wait()
  return result.hash
}

export async function checkAllowance(tokenAddress: string, spender: string) {
  const provider = new ethers.BrowserProvider(await getPureWalletProvider())
  const abi = [
    'function allowance(address owner, address spender) view returns (uint256)',
  ]
  const tokenContract = new ethers.Contract(tokenAddress, abi, provider)
  const owner = currentWallet()
  const allowance = await tokenContract.allowance(owner, spender)
  return ethers.formatUnits(allowance, 18)
}

export function fromReadableAmount(amount: number, decimals: number) {
  return ethers.parseUnits(amount.toString(), decimals)
}

export function calculateGasMargin(value: bigint, margin = 1000n): bigint {
  return (value * (10000n + margin)) / 10000n
}
