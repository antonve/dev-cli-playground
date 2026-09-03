import { createRoot } from "react-dom/client"
import {
  Outlet,
  RouterProvider,
  createRootRoute,
  createRoute,
  createRouter,
} from "@tanstack/react-router"
import { StrictMode, useEffect, useState } from "react"

import { clearRoute, selectRoute, selectedRoute, validRoute } from "./branch"
import "./styles.css"

type Diagnostic = Record<string, unknown>

function App() {
  const [selected, setSelected] = useState(() => selectedRoute())
  const [candidate, setCandidate] = useState(() => selectedRoute())
  const [diagnostic, setDiagnostic] = useState<Diagnostic>({ loading: true })
  const [error, setError] = useState("")

  useEffect(() => {
    Promise.all([
      fetch("/__web").then((response) => response.json()),
      fetch("/api/diagnostic").then((response) => response.json()),
    ])
      .then(([web, request]) => setDiagnostic({ web, request }))
      .catch((cause: unknown) => setDiagnostic({ error: String(cause) }))
  }, [selected])

  const choose = () => {
    setError("")
    if (!validRoute(candidate)) {
      setError("Use the DNS-safe route key printed by dev up or dev status.")
      return
    }
    selectRoute(candidate)
    setSelected(candidate)
    window.location.reload()
  }

  const clear = () => {
    clearRoute()
    setSelected("")
    setCandidate("")
    window.location.reload()
  }

  return (
    <main>
      <header>
        <span className="eyebrow">Bazel × Kubernetes × TanStack</span>
        <h1>Dev CLI Playground</h1>
        <p>Each service independently selects this branch overlay or its always-running base.</p>
      </header>
      <section className="menu">
        <label htmlFor="route">Branch route key</label>
        <input
          id="route"
          value={candidate}
          onChange={(event) => setCandidate(event.target.value.trim())}
          placeholder="owner-feature-hash"
          autoComplete="off"
        />
        <div className="actions">
          <button type="button" onClick={choose}>Select branch</button>
          <button type="button" className="secondary" onClick={clear}>Clear</button>
        </div>
        {error && <p className="error">{error}</p>}
        <p>Selected: <code>{selected || "base"}</code></p>
      </section>
      <section>
        <h2>Request diagnostics</h2>
        <pre>{JSON.stringify(diagnostic, null, 2)}</pre>
      </section>
      <Outlet />
    </main>
  )
}

const rootRoute = createRootRoute({ component: App })
const indexRoute = createRoute({ getParentRoute: () => rootRoute, path: "/" })
const routeTree = rootRoute.addChildren([indexRoute])
const router = createRouter({ routeTree })

declare module "@tanstack/react-router" {
  interface Register { router: typeof router }
}

const root = document.getElementById("root")
if (!root) throw new Error("root element missing")
createRoot(root).render(<StrictMode><RouterProvider router={router} /></StrictMode>)
