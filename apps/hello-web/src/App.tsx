import { Outlet } from "@tanstack/react-router"
import { useEffect, useState } from "react"

type Diagnostic = Record<string, unknown>

export default function App() {
  const [diagnostic, setDiagnostic] = useState<Diagnostic>({ loading: true })

  useEffect(() => {
    Promise.all([
      fetch("/__web").then((response) => response.json()),
      fetch("/api/diagnostic").then((response) => response.json()),
    ])
      .then(([web, request]) => setDiagnostic({ web, request }))
      .catch((cause: unknown) => setDiagnostic({ error: String(cause) }))
  }, [])

  return (
    <main>
      <header>
        <span className="eyebrow">Bazel × Kubernetes × TanStack</span>
        <h1>Dev CLI Playground</h1>
        <p>Each service independently selects this branch overlay or its always-running base.</p>
      </header>
      <section>
        <p>Host: <code>{window.location.hostname}</code></p>
        <label htmlFor="route">Unsaved scratch note</label>
        <input id="route" placeholder="Type here to verify HMR preserves input" autoComplete="off" />
      </section>
      <section>
        <h2>Request diagnostics</h2>
        <pre>{JSON.stringify(diagnostic, null, 2)}</pre>
      </section>
      <Outlet />
    </main>
  )
}
