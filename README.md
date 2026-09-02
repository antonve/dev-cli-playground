# dev-cli playground

Conformance fixture for [`antonve/dev-cli`](https://github.com/antonve/dev-cli).
It contains three intentionally small Bazel-owned deployables and all default
development-cluster manifests under `k8s/dev`.

```sh
bazel test --config=agent //...
bazel build --config=agent //...
kubectl kustomize k8s/dev
```

The browser reaches `dev-router`, which discards any client-supplied
`x-dev-branch` header and recreates it from the scoped `dev_branch` cookie.
Both the router and `hello-api` independently prefer a matching overlay Service
and fall back to the always-running base. See `docs/architecture.md` and
`docs/conformance.md`.
