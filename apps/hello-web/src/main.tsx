import { createRoot } from "react-dom/client"
import { RouterProvider, createRootRoute, createRoute, createRouter } from "@tanstack/react-router"
import { StrictMode } from "react"
import App from "./App"
import "./styles.css"

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
