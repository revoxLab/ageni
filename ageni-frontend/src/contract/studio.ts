import {StudioMessage} from '@/apis'
import {storage} from '@/utils'
import {
  cancelOrderOnKiloex,
  checkAllowance,
  closePositionOnKiloex,
  connectWallet,
  crossChain,
  currentConnectState,
  currentNetwork,
  currentWallet,
  estimateGasFromPancake,
  FunctionCallError,
  getAllOrdersOnKiloex,
  getAllPairInfos,
  getAllPositionsOnKiloex,
  getBalanceOfGas,
  getChainId,
  getFeeInfo,
  getPairInfo,
  getProductPricesOnKiloex,
  getQuotes,
  getReceipt,
  getTokenBalanceForErc20,
  getTokenTransferApproval,
  getTradesHistoryOnKiloex,
  increasePositionOnKiloex,
  isConnected,
  listSupportedChains,
  listSupportedTokens,
  signMessageForAI,
  supportChains,
  supportChainsOnKiloex,
  supportedProductsOnKiloex,
  swapTokens,
  switchChain,
  transfer,
  updateOrderOnKiloex,
  updatePositionMarginOnKiloex,
} from '@/wallet'
import {ApiError} from 'owlto-sdk'

export async function processToolCalls(
  toolCall: NonNullable<StudioMessage['tool_calls']>[number]
): Promise<any | null> {
  console.log('start tool call,content is ', toolCall)

  if (toolCall.type === 'function') {
    const functionName = toolCall.function.name
    const argsString = toolCall.function.arguments
    const result = {tool_call_id: toolCall.id, content: 'failed'}

    if (functionName === 'connect_wallet') {
      const connectStateOfWallet = isConnected()
      if (connectStateOfWallet !== true) {
        try {
          const address = await connectWallet()
          if (address) {
            result.content = buildResult('success')
          } else {
            result.content = buildResult('failed', 'address is null')
          }
        } catch (error) {
          result.content = processError(
            error,
            'connect_wallet_failed',
            argsString,
            'connect wallet failed'
          )
        }
      } else {
        result.content = buildResult('success')
      }
    } else if (functionName === 'current_network') {
      try {
        const {chainName} = await currentNetwork()
        result.content = buildResult(chainName)
      } catch (error) {
        result.content = processError(
          error,
          'current_network_failed',
          argsString,
          'get current network failed'
        )
      }
    } else if (functionName === 'current_wallet') {
      try {
        const address = currentWallet()
        result.content = address
          ? buildResult(address)
          : buildResult('failed', 'address is null')
      } catch (error) {
        result.content = processError(
          error,
          'current_wallet_failed',
          argsString,
          'get current wallet failed'
        )
      }
    } else if (functionName === 'current_connect_status') {
      try {
        const info = await currentConnectState()
        result.content = buildResult(info)
      } catch (error) {
        result.content = processError(
          error,
          'current_connect_status_failed',
          argsString,
          'get current connect status failed'
        )
      }
    } else if (functionName === 'switch_network') {
      if (argsString !== null && typeof argsString === 'string') {
        const args = JSON.parse(argsString)
        const network = args.network
        try {
          await switchChain(supportChains()[network.toLowerCase()])
          result.content = buildResult('success')
        } catch (error) {
          result.content = processError(
            error,
            'switch_network_failed',
            argsString,
            'switch network failed'
          )
        }
      }
    } else if (functionName === 'get_gas_balance') {
      try {
        const balance = await getBalanceOfGas()
        result.content = buildResult(balance)
      } catch (error) {
        result.content = processError(
          error,
          'get_gas_balance_failed',
          argsString,
          'get gas balance failed'
        )
      }
    } else if (functionName === 'get_token_balance') {
      if (argsString != null && typeof argsString === 'string') {
        if (argsString.length > 0) {
          const args = JSON.parse(argsString)
          const contractAddress = args.contract_address
          try {
            const balance = await getTokenBalanceForErc20(contractAddress)
            result.content = buildResult(balance)
          } catch (error) {
            result.content = processError(
              error,
              'get_token_balance_failed',
              argsString,
              'get token balance failed'
            )
          }
        } else {
          try {
            const balance = await getTokenBalanceForErc20()
            result.content = buildResult(balance)
          } catch (error) {
            result.content = processError(
              error,
              'get_token_balance_failed',
              argsString,
              'get token balance failed'
            )
          }
        }
      } else {
        try {
          const balance = await getTokenBalanceForErc20()
          result.content = buildResult(balance)
        } catch (error) {
          result.content = processError(
            error,
            'get_token_balance_failed',
            argsString,
            'get token balance failed'
          )
        }
      }
    } else if (functionName === 'sign_message') {
      try {
        if (argsString !== null && typeof argsString === 'string') {
          const args = JSON.parse(argsString)
          const message = args.message
          const info = await signMessageForAI(message)
          result.content = buildResult(info)
        }
      } catch (error) {
        result.content = processError(
          error,
          'sign_message_failed',
          argsString,
          'sign message failed'
        )
      }
    } else if (functionName === 'transfer') {
      if (argsString !== null && typeof argsString === 'string') {
        const args = JSON.parse(argsString)
        try {
          const hash = await transfer(
            args.receive_addr,
            args.amount,
            args.contract_address
          )
          result.content = buildResult(hash)
        } catch (error) {
          result.content = processError(
            error,
            'transfer_failed',
            argsString,
            'transfer failed'
          )
        }
      }
    } else if (functionName === 'get_quotes') {
      if (argsString !== null && typeof argsString === 'string') {
        try {
          const args = JSON.parse(argsString)
          const amount = Number(args.input_amount ?? args.output_amount)
          const [fromToken, toToken] = args.input_amount
            ? [args.consumed_token, args.desired_token]
            : [args.desired_token, args.consumed_token]
          const content =
            (await getQuotes(args.dex, amount, fromToken, toToken)) ?? 'failed'
          const price = Number(content) / amount
          const tradingFee = await estimateGasFromPancake(
            storage.get('address')!,
            amount,
            fromToken,
            toToken,
            0.5
          )
          result.content = buildResult({
            price,
            received: Number(content),
            tradingFee: tradingFee?.toString(),
            chainId: await getChainId(),
          })
        } catch (error) {
          result.content = processError(
            error,
            'get_quotes_failed',
            argsString,
            'get quotes failed'
          )
        }
      }
    } else if (functionName === 'swap_tokens') {
      if (argsString !== null && typeof argsString === 'string') {
        try {
          const args = JSON.parse(argsString)
          const content =
            (await swapTokens(
              args.dex,
              args.input_amount,
              args.consumed_token,
              args.desired_token,
              Number(args.slippage_tolerance)
            )) ?? 'failed'
          result.content = buildResult(content)
        } catch (error) {
          result.content = processError(
            error,
            'swap_failed',
            argsString,
            'swap tokens failed'
          )
        }
      } else {
        result.content = buildResult('failed', 'argsString is null')
      }
    } else if (functionName === 'allowance') {
      if (argsString !== null && typeof argsString === 'string') {
        try {
          const args = JSON.parse(argsString)
          const res = await checkAllowance(args.contract_address, args.spender)
          if (res != null) {
            result.content = buildResult(res)
          } else {
            result.content = buildResult('failed')
          }
        } catch (error) {
          result.content = processError(
            error,
            'allowance_failed',
            argsString,
            'allowance failed'
          )
        }
      }
    } else if (functionName === 'approve') {
      if (argsString !== null && typeof argsString === 'string') {
        try {
          const args = JSON.parse(argsString)
          const res = await getTokenTransferApproval(
            args.contract_address,
            Number(args.value)
          )
          result.content = buildResult(res ?? 'failed')
        } catch (error) {
          result.content = processError(
            error,
            'approve_failed',
            argsString,
            'approve failed'
          )
        }
      }
    } else if (functionName === 'execute_bridge_transfer') {
      if (argsString !== null && typeof argsString === 'string') {
        try {
          const args = JSON.parse(argsString)
          const res = await crossChain(
            args.tokenName,
            args.fromChainName,
            args.toChainName,
            args.fromAddress,
            args.toAddress,
            args.uiValue,
            args.valueIncludeGasFee
          )
          result.content = buildResult(res ?? 'failed')
        } catch (error) {
          result.content = processError(
            error,
            'crossChain_failed',
            argsString,
            'cross chain failed'
          )
        }
      }
    } else if (functionName === 'getAllPairInfos') {
      if (argsString !== null && typeof argsString === 'string') {
        try {
          const args = JSON.parse(argsString)
          const res = await getAllPairInfos(
            args.category,
            args.valueIncludeGasFee
          )
          result.content = buildResult(res ?? 'failed')
        } catch (error) {
          result.content = processError(
            error,
            'getAllPairInfos_failed',
            argsString,
            'get all pairInfos failed'
          )
        }
      }
    } else if (functionName === 'get_bridge_pair_info') {
      if (argsString !== null && typeof argsString === 'string') {
        try {
          const args = JSON.parse(argsString)
          const res = await getPairInfo(
            args.tokenName,
            args.fromChainName,
            args.toChainName,
            args.valueIncludeGasFee
          )
          result.content = buildResult(res ?? 'failed')
        } catch (error) {
          result.content = processError(
            error,
            'getPair_failed',
            argsString,
            'get pair failed'
          )
        }
      }
    } else if (functionName === 'estimate_bridge_fee') {
      if (argsString !== null && typeof argsString === 'string') {
        try {
          const args = JSON.parse(argsString)
          const res = await getFeeInfo(
            args.tokenName,
            args.fromChainName,
            args.toChainName,
            args.uiValue
          )
          result.content = buildResult(res ?? 'failed')
        } catch (error) {
          result.content = processError(
            error,
            'getFeeInfo_failed',
            argsString,
            'get Fee Info failed'
          )
        }
      }
    } else if (functionName === 'get_bridge_transaction_receipt') {
      if (argsString !== null && typeof argsString === 'string') {
        try {
          const args = JSON.parse(argsString)
          const res = await getReceipt(args.fromChainHash)
          result.content = buildResult(res ?? 'failed')
        } catch (error) {
          result.content = processError(
            error,
            'getReceipt_failed',
            argsString,
            'get Receipt failed'
          )
        }
      }
    } else if (functionName === 'list_supported_bridge_chains') {
      try {
        const res = await listSupportedChains()
        result.content = buildResult(res ?? 'failed')
      } catch (error) {
        result.content = processError(
          error,
          'listSupportedChains_failed',
          argsString,
          'list Supported Chains failed'
        )
      }
    } else if (functionName === 'get_bridge_compatible_tokens') {
      if (argsString !== null && typeof argsString === 'string') {
        try {
          const args = JSON.parse(argsString)
          const res = await listSupportedTokens(
            args.fromChainName,
            args.toChainName
          )
          result.content = buildResult(res ?? 'failed')
        } catch (error) {
          result.content = processError(
            error,
            'listSupportedTokens_failed',
            argsString,
            'list Supported Tokens failed'
          )
        }
      }
    } else if (functionName === 'list_supported_perpetual_chains') {
      result.content = buildResult(supportChainsOnKiloex())
    } else if (functionName === 'get_perpetual_products_by_chain') {
      if (argsString !== null && typeof argsString === 'string') {
        try {
          const args = JSON.parse(argsString)
          result.content = buildResult(
            await supportedProductsOnKiloex(args.chainId)
          )
        } catch (error) {
          result.content = processError(
            error,
            'supportedProductsByChain_failed',
            argsString,
            'get supported Products By Chain failed'
          )
        }
      }
    } else if (functionName === 'get_perpetual_product_price') {
      if (argsString !== null && typeof argsString === 'string') {
        try {
          const args = JSON.parse(argsString)
          result.content = buildResult(
            await getProductPricesOnKiloex([args.productId])
          )
        } catch (error) {
          result.content = processError(
            error,
            'getProductPrice_failed',
            argsString,
            'get Product Price failed'
          )
        }
      }
    } else if (
      [
        'increase_perpetual_position',
        'set_take_profit_stop_loss',
        'create_position_single_tpsl',
        'create_position_dual_tpsl',
      ].includes(functionName)
    ) {
      if (argsString !== null && typeof argsString === 'string') {
        try {
          const args = JSON.parse(argsString)
          const res = await increasePositionOnKiloex(args)
          result.content = buildResult(res ?? 'failed')
        } catch (error) {
          result.content = processError(
            error,
            `${functionName}_failed`,
            argsString,
            `${functionName} failed`
          )
        }
      }
    } else if (functionName === 'close_perpetual_position') {
      if (argsString !== null && typeof argsString === 'string') {
        try {
          const args = JSON.parse(argsString)
          const res = await closePositionOnKiloex(
            args.productId,
            args.margin,
            args.isLong,
            args.tickerPrice
          )
          result.content = buildResult(res ?? 'failed')
        } catch (error) {
          result.content = processError(
            error,
            'closePosition_failed',
            argsString,
            'close Position failed'
          )
        }
      }
    } else if (functionName === 'update_perpetual_position_margin') {
      if (argsString !== null && typeof argsString === 'string') {
        try {
          const args = JSON.parse(argsString)
          const res = await updatePositionMarginOnKiloex(
            args.productId,
            args.margin,
            args.isLong
          )
          result.content = buildResult(res ?? 'failed')
        } catch (error) {
          result.content = processError(
            error,
            'updatePositionMargin_failed',
            argsString,
            'update Position Margin failed'
          )
        }
      }
    } else if (functionName === 'cancel_perpetual_order') {
      if (argsString !== null && typeof argsString === 'string') {
        try {
          const args = JSON.parse(argsString)
          const res = await cancelOrderOnKiloex(args.type, args.orderIndex)
          result.content = buildResult(res ?? 'failed')
        } catch (error) {
          result.content = processError(
            error,
            'cancelOrder_failed',
            argsString,
            'cancel Order failed'
          )
        }
      }
    } else if (functionName === 'update_perpetual_order') {
      if (argsString !== null && typeof argsString === 'string') {
        try {
          const args = JSON.parse(argsString)
          const res = await updateOrderOnKiloex(
            args.orderType,
            args.orderIndex,
            args.margin,
            args.leverage,
            args.limitPrice
          )
          result.content = buildResult(res ?? 'failed')
        } catch (error) {
          result.content = processError(
            error,
            'updateOrder_failed',
            argsString,
            'update Order failed'
          )
        }
      }
    } else if (functionName === 'get_all_perpetual_positions') {
      if (argsString !== null && typeof argsString === 'string') {
        try {
          const res = await getAllPositionsOnKiloex()
          result.content = buildResult(res ?? 'failed')
        } catch (error) {
          result.content = processError(
            error,
            'getAllPositions_failed',
            argsString,
            'get All Positions failed'
          )
        }
      }
    } else if (functionName === 'get_all_perpetual_orders') {
      if (argsString !== null && typeof argsString === 'string') {
        try {
          const res = await getAllOrdersOnKiloex()
          result.content = buildResult(res ?? 'failed')
        } catch (error) {
          result.content = processError(
            error,
            'getAllOrders_failed',
            argsString,
            'Get All Orders failed'
          )
        }
      }
    } else if (functionName === 'get_trades_history') {
      try {
        const res = await getTradesHistoryOnKiloex()
        result.content = buildResult(res ?? 'failed')
      } catch (error) {
        result.content = processError(
          error,
          `${functionName}_failed`,
          argsString,
          `${functionName} failed`
        )
      }
    }
    return result
  }
}

function buildResult(result: any, reason: string = ''): string {
  return JSON.stringify({result, reason})
}

function processError(
  error: any,
  reportLabel: string,
  argsString: string,
  callbackFailedReason: string
) {
  let content
  if (error?.info?.error?.code === 4001) {
    content = buildResult('failed', error.info.error.message)
  } else if (error instanceof ApiError) {
    content = buildResult('failed', error.status.message)
  } else if (error instanceof FunctionCallError) {
    content = buildResult('failed', error.message)
  } else {
    content = buildResult(
      'failed',
      typeof error === 'string'
        ? error
        : (error?.info?.error?.data?.message ??
            error?.info?.error?.message ??
            error?.error?.message ??
            error?.message ??
            callbackFailedReason)
    )
  }
  console.error(error)
  return content
}
