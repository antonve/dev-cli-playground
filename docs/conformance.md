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

With dev-cli v0.3+, use `dev url --owner <owner> '/desired/path?tab=profile'`
to open a branch directly; the gateway sets the cookie without any frontend
handler. `dev url --owner <owner> --clear '/desired/path'` selects base.
The older menu remains optional, not an integration requirement.

For two environments with only hello-web and echo-api changed, run the
read-only query-link test (each route must have been started with v0.3+):

```sh
python3 conformance/deep_links.py \
  --route-a "$route_a" --branch-a "$branch_a" \
  --route-b "$route_b" --branch-b "$branch_b"
```

It verifies first-request selection, old-cookie and forged-header precedence,
cookie persistence, switching, clearing, separate sessions, a deep page with
other query parameters, HEAD selection, POST non-selection, unknown selectors,
80 concurrent requests, and base hello-api calling overlay echo-api. It does not
mutate Kubernetes resources or claim browser HMR proof.

Browser acceptance: open the generated links in two isolated browser contexts,
check that the existing path/query/fragment survive, and switch/clear without
using the menu. With one foreground loop running, make a second App.tsx heading
edit: the DOM should update without changing `performance.timeOrigin` or unsaved
input. Change echo's displayed service label and fetch fresh diagnostics to prove
binary live sync; compare pod UIDs and image digests before and after. The other
owner and a no-cookie browser must remain unchanged.

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

After creating two branch overlays, run the read-only concurrent validator:

```sh
python3 conformance/routing.py \
  --route-a "$route_a" --branch-a acceptance/routing-a --services-a hello-api,echo-api \
  --route-b "$route_b" --branch-b acceptance/routing-b --services-b echo-api
```

It checks web and both backend identities, forged-header precedence, missing
cookies, and 100 interleaved HTTP requests from independent cookie sessions.
After deleting only A's labeled echo Service, rerun with `--services-a hello-api`:
A's API stays on its branch and echo must fall back independently. Wait for
health/DNS convergence; native failover does not promise zero transient errors.

For browser HMR, select the route using the menu, type an unsaved input value,
and edit `apps/hello-web/src/App.tsx` repeatedly. The heading must change while
the input value and `performance.timeOrigin` remain unchanged. Watch for the
Vite WebSocket connected/hot-updated messages. Editing bootstrap or dev-server
configuration may legitimately reload the page or restart the server.

`dev status` reports route ownership, affected workloads, HTTPRoute count,
sync health, and expiry. `dev down` deletes only resources matching the current
owner and route. `dev cleanup` removes expired CLI-managed routes.
