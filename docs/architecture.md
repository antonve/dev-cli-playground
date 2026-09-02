# Playground architecture

Argo CD owns the namespace, router, ingress, and base workloads from `k8s/dev`.
The CLI owns only labeled branch Deployments and Services. No base Deployment
is patched.

Ingress sends all browser traffic to `dev-router`. The router ignores a direct
`x-dev-branch` header, reads the HttpOnly `dev_branch` cookie, validates the
route key, and emits one normalized internal header. It independently probes
`hello-web-dev-<route>` or `hello-api-dev-<route>` and falls back to the base
Service. `hello-api` repeats the same decision for `echo-api` and propagates the
normalized header. This creates per-service fallback without a mesh or overlay
controller.

Each diagnostic response includes service, branch, source revision, pod,
observed route, target, and fallback state. Route keys include owner and branch
plus a hash, allowing the same branch name under different owners.

The frontend is a TypeScript React module served by a tiny Go development
server. The module is syntactically valid browser JavaScript and pins React ESM
imports. The server watches synchronized source content, emits an SSE HMR event,
and the bootstrap dynamically reloads and re-renders the module without a page
reload. Backend images contain the base service and `dev-supervisor`; a sync
sidecar stages successful Bazel binaries on a shared volume.
