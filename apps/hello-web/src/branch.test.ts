import { describe, expect, it } from "vitest"

import { clearRoute, selectRoute, selectedRoute, validRoute } from "./branch"

describe("branch cookie", () => {
  it("accepts route keys and rejects header-shaped input", () => {
    expect(validRoute("alice-feature-a1b2c3d4")).toBe(true)
    expect(validRoute("Feature X")).toBe(false)
    expect(validRoute("anything\r\nx-dev-branch: stolen")).toBe(false)
  })

  it("selects and clears a scoped cookie", () => {
    selectRoute("alice-feature-a1b2c3d4")
    expect(selectedRoute()).toBe("alice-feature-a1b2c3d4")
    clearRoute()
    expect(selectedRoute()).toBe("")
  })
})
