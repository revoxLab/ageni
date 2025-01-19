import {NodeGlobalsPolyfillPlugin} from '@esbuild-plugins/node-globals-polyfill'
import legacy from '@vitejs/plugin-legacy'
import vue from '@vitejs/plugin-vue'
import {fileURLToPath, URL} from 'node:url'
import nodePolyfills from 'rollup-plugin-node-polyfills'
import {defineConfig} from 'vite'

// https://vitejs.dev/config/
export default defineConfig({
  // server: {
  //   port: 80,
  //   host: true,
  // },
  optimizeDeps: {
    esbuildOptions: {
      // Fix global is not defined error
      define: {
        global: 'globalThis',
      },
      supported: {
        bigint: true,
      },
      plugins: [
        // Without this, npm run dev will output Buffer or process is not defined error
        NodeGlobalsPolyfillPlugin({
          buffer: true,
          process: true,
        }),
      ],
    },
  },
  build: {
    rollupOptions: {
      plugins: [nodePolyfills() as any],
    },
    commonjsOptions: {
      transformMixedEsModules: true,
    },
  },
  plugins: [
    vue(),
    legacy({
      targets: ['defaults', 'not IE 11'],
      renderLegacyChunks: false,
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
})
