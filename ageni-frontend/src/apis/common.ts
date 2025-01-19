import {defineMutation} from './define'

export const walletLoginMutation = defineMutation({
  url: '/web/wallet_login',
  mockParams: {
    wallet_address: '0x1111111',
    signature: '111111',
    from: 'studio',
  },
  mockData: {
    code: 0,
    message: 'success',
    data: {
      token: '1111',
    },
  },
})
