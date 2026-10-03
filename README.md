# dev-cli playground

Conformance fixture for [`antonve/dev-cli`](https://github.com/antonve/dev-cli).
It contains three intentionally small Bazel-owned deployables and all default
development-cluster manifests under `k8s/dev`.

```sh
pnpm install --frozen-lockfile
pnpm typecheck && pnpm build
bazel test --config=agent //...
bazel build --config=agent //apps/hello-web:image //apps/hello-api:image //apps/echo-api:image
kubectl kustomize k8s/dev
```

The pnpm/TanStack/Vite frontend supports filesystem sync and native HMR. Envoy
Gateway derives branch routing from the branch hostname, and the CLI
creates per-service HTTPRoutes with independent base fallback. Repository images
are registry-neutral in Bazel; CI publishes base images to GHCR and `dev up`
publishes affected branch images to the configured development registry. See
`docs/architecture.md` and `docs/conformance.md`.
