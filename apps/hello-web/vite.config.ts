import react from "@vitejs/plugin-react"
import type { IncomingMessage, ServerResponse } from "node:http"
import { defineConfig, type Plugin } from "vite"

function identity() {
  return {
    service: "hello-web",
    branch: process.env.DEV_BRANCH || "base",
    revision: process.env.DEV_REVISION || "base",
    pod: process.env.HOSTNAME || "local",
  }
}

function diagnostics(): Plugin {
  const handle = (req: IncomingMessage, res: ServerResponse, next: () => void) => {
    if (req.url === "/readyz") {
      res.statusCode = 200
      res.end("ok\n")
      return
    }
    if (req.url === "/__web") {
      res.setHeader("content-type", "application/json")
      res.end(JSON.stringify(identity()))
      return
    }
    next()
  }
  return {
    name: "playground-diagnostics",
    configureServer(server) {
      server.middlewares.use(handle)
    },
    configurePreviewServer(server) {
      server.middlewares.use(handle)
    },
  }
}

export default defineConfig({
  resolve: { preserveSymlinks: true },
  plugins: [react(), diagnostics()],
  server: {
    allowedHosts: true,
    watch: { usePolling: true, interval: 100 },
  },
})
