module.exports = {
  plugins: {
    autoprefixer: {
      overrideBrowserslist: [
        'Android 4.1',
        'iOS 7.1',
        'Chrome > 31',
        'ff > 31',
        'ie >= 8',
        'last 10 versions',
      ],
      grid: true,
    },
    'postcss-mobile-forever': {
      viewportWidth: 375,
      maxDisplayWidth: 1440,
      propList: [
        '*',
        '!border',
        '!border-left',
        '!border-right',
        '!border-top',
        '!border-bottom',
      ],
      exclude: /\/src\/assets\/style\//,
    },
    'postcss-px-to-viewport-update': {
      viewportWidth: 1440,
      exclude: /\/src\/mobile\/|\/src\/components\/mobile\//,
    },
  },
}
