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

The default configuration targets `homelab-talos-dev`. The Homelab Talos root
owns the playground Argo Application, which deploys `k8s/dev` using promoted
CI image digests. Before starting branch checks, confirm the Application is
Synced and Healthy and all three base Deployments are Ready. Base workloads and
routes stay Argo-owned; conformance environments use distinct CLI owners and
are removed with their matching `dev down --owner <owner>`.

With dev-cli v0.8+, `dev up` and `dev url --owner <owner>` print
`https://<route>.dev-cli-playground.dev.lab/`. Open that hostname in a fresh
browser profile with no selection cookies. Two owners can share one profile;
each host replaces forged routing headers and selects its own internal hop.
The Lab CA certificate must validate without bypass in curl.

Run the hostname matrix with two owners after their overlays are Ready:

```sh
python3 conformance/routing.py --branch-hosts \
  --route-a "$route_a" --branch-a "$branch_a" --services-a hello-web,echo-api \
  --route-b "$route_b" --branch-b "$branch_b" --services-b echo-api
```

The validator checks the base host, forged-header removal, both selected hosts,
independent service fallback and 80 interleaved requests. Start A with only
hello-web/echo-api selected and B with only echo-api selected. Wait for actual
backend identity after startup: Kubernetes Ready can precede Envoy health/DNS
convergence. To test deletion fallback, delete only A's labeled echo Service
and rerun with `--services-a hello-web`; B must remain selected.

Browser acceptance: open both branch hosts in the same context; assert no
selection cookie is needed. Type an unsaved input, edit the hello-web heading,
and require DOM HMR with unchanged `performance.timeOrigin` and input value.
Change echo's displayed service label, fetch diagnostics and compare pod UID
and image digest before/after to prove live binary sync. Base and B stay unchanged.

`dev status` reports owner, branch URLs, workloads, routing/sync health and expiry.
`dev down --owner <owner>` must remove only that owner's routing and workloads;
the deleted host returns 404 and no branch Ingress, Certificate or TLS Secret
remains. Expiry cleanup must satisfy the same checks and preserve the other owner.
Require a Synced/Healthy Argo base before and after the run.

## v0.8 compatibility checks

Cookie selection remains available through `dev url --cookie` during adoption.
Run `conformance/deep_links.py` and `conformance/routing.py` without
`--branch-hosts` to verify the previous selection behavior. The optional menu
exists for these compatibility checks and is unnecessary on branch hosts.
