import { createReadStream, existsSync } from "node:fs"
import { extname, join, normalize } from "node:path"
import { createServer } from "node:http"

const root = new URL("./dist/", import.meta.url).pathname
const port = Number(process.env.PORT || 8080)
const types = new Map([
  [".css", "text/css; charset=utf-8"],
  [".html", "text/html; charset=utf-8"],
  [".js", "text/javascript; charset=utf-8"],
  [".json", "application/json"],
  [".svg", "image/svg+xml"],
])

function identity() {
  return {
    service: "hello-web",
    branch: process.env.DEV_BRANCH || "base",
    revision: process.env.DEV_REVISION || "base",
    pod: process.env.HOSTNAME || "local",
  }
}

createServer((request, response) => {
  const path = new URL(request.url || "/", "http://localhost").pathname
  if (path === "/readyz") {
    response.writeHead(200, { "content-type": "text/plain; charset=utf-8" })
    response.end("ok\n")
    return
  }
  if (path === "/__web") {
    response.writeHead(200, { "content-type": "application/json" })
    response.end(JSON.stringify(identity()))
    return
  }
  const relative = normalize(decodeURIComponent(path)).replace(/^[/\\]+/, "")
  const candidate = join(root, relative || "index.html")
  const file = candidate.startsWith(root) && existsSync(candidate) ? candidate : join(root, "index.html")
  response.writeHead(200, {
    "cache-control": file.endsWith("index.html") ? "no-store" : "public, max-age=31536000, immutable",
    "content-type": types.get(extname(file)) || "application/octet-stream",
  })
  createReadStream(file).pipe(response)
}).listen(port, "0.0.0.0", () => {
  console.log(`hello-web listening on :${port}`)
})
