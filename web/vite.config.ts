import { defineConfig, loadEnv, type Plugin } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import { fileURLToPath, URL } from 'node:url'
import path from 'node:path'

const src = (p: string) => fileURLToPath(new URL(`./src/${p}`, import.meta.url))

function mockApiPlugin(isMock: boolean): Plugin {
  const realId = path.resolve(src('api/real.ts'))
  const mockId = path.resolve(src('api/mock.ts'))
  return {
    name: 'mock-api',
    enforce: 'pre',
    resolveId(id, importer) {
      if (!isMock) return
      const resolved = importer ? path.resolve(path.dirname(importer), id) : id
      if (resolved === realId || id.endsWith('api/real.ts') || id.endsWith('api/index.ts')) {
        return mockId
      }
    },
  }
}

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const isMock = env.VITE_MOCK === 'true'

  if (isMock) console.log('\x1b[33m[mock] API mock enabled — no real router needed\x1b[0m')

  return {
    plugins: [mockApiPlugin(isMock), vue(), tailwindcss()],
    resolve: {
      alias: { '@': src('') },
    },
    server: {
      proxy: isMock ? {} : {
        '/api': { target: 'http://localhost:8080', changeOrigin: true },
      },
    },
  }
})
