import react from "@vitejs/plugin-react"
import { defineConfig } from "vitest/config"

export default defineConfig({
  resolve: { preserveSymlinks: true },
  plugins: [react()],
  test: { environment: "jsdom" },
})
