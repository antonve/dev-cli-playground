# Conformance and operations

## Clean build and render

```sh
pnpm install --frozen-lockfile
pnpm test
pnpm typecheck
pnpm build
bazel test --config=agent //...
bazel build --config=agent //apps/hello-web:image //apps/hello-api:image //apps/echo-api:image
kubectl kustomize k8s/dev >/dev/null
```

## Base and routing checks

Set `route` to the key printed by `dev up` or `dev status`:

```sh
host=https://dev-cli-playground.dev.lab
route=<route-key>

# No cookie: all hops are base. A forged browser header is removed.
curl --cacert path/to/lab-root-ca.crt "$host/api/diagnostic"
curl --cacert path/to/lab-root-ca.crt -H 'x-dev-branch: forged' "$host/api/diagnostic"

# Selected branch: each service independently uses its overlay or base.
curl --cacert path/to/lab-root-ca.crt -H "Cookie: dev_branch=$route" "$host/api/diagnostic"

# A separate session remains base.
curl --cacert path/to/lab-root-ca.crt "$host/api/diagnostic"
```

The required matrix is: no cookie uses all base services; a hello-api-only
overlay uses overlay hello-api and base echo-api; adding echo-api selects both
overlays; a second cookie selects another route without leakage; removing one
overlay falls back for only that service. The frontend menu writes or clears
the same Secure, SameSite=Lax, path-scoped cookie.

`dev status` reports route ownership, affected workloads, HTTPRoute count,
sync health, and expiry. `dev down` deletes only resources matching the current
owner and route. `dev cleanup` removes expired CLI-managed routes.
