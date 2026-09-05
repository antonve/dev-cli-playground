# Playground architecture

The repository has only application code, application build rules, deployable
metadata, and default application manifests. It contains no development router,
restart supervisor, sync sidecar, or other CLI runtime.

## Delivery boundary

GitHub Actions builds the three application images on `ubuntu-latest` from the
Bazel graph and publishes both `latest` and `sha-<commit>` tags:

- `ghcr.io/antonve/dev-cli-playground/hello-web`
- `ghcr.io/antonve/dev-cli-playground/hello-api`
- `ghcr.io/antonve/dev-cli-playground/echo-api`

The same run retains Docker-compatible image archives for seven days. The
Homelab release gate promotes those exact CI outputs to the private development
registry as `registry.dev.lab/dev-cli-playground/<service>:latest`; the hosted
runner does not receive private-Lab network access. GHCR remains the CI build
record, while the cluster pulls anonymously from its dedicated Zot registry.

The workload specs use logical image names. The environment Kustomization is
the single injection point that maps those names to the configured Zot
repositories and `latest` tags. A real application should let Argo CD Image
Updater write immutable digests at that layer.

For branch work, Bazel metadata names an image and a registry-neutral push
target. `dev up` builds only affected image targets, publishes them to the
configured registry using a collision-resistant branch/revision tag, resolves
the registry digest, and injects that immutable reference into the CLI-owned
Deployment. Image locations are not embedded in Bazel rules or workload specs.

## Routing boundary

Homelab owns Envoy Gateway and a shared Gateway. The playground declares only
its base HTTPRoutes and an ingress pointing at the platform gateway Service.
The CLI declares temporary, labeled HTTPRoutes for each active branch.

For browser traffic, a branch route matches the scoped `dev_branch` cookie and
replaces `x-dev-branch` with the validated route key. Base routes remove any
browser-supplied routing header. Route precedence selects the matching branch
route when present and the base route otherwise. Every service gets its own
branch route, so its backend independently points to the overlay Service or the
always-running base Service.

`hello-api` forwards the normalized header through the same Gateway using the
internal host `echo-api.internal`. An internal branch route independently
chooses the echo overlay or base. No application probes Kubernetes, discovers
overlay Services, or implements routing policy.

## Development boundary

`hello-web` is a pnpm workspace application using React, TanStack Router, Vite,
and Vitest. Bazel invokes Vite for the production bundle and builds the OCI
image. For development, the CLI starts that image once with the repository's
declared pnpm/Vite command, copies changed source files into its writable
workspace, and leaves Vite responsible for filesystem watching and HMR.
Sync includes the root pnpm manifest, lockfile, and workspace configuration.
Metadata declares those dependency inputs and the frozen pnpm install command,
so the CLI reinstalls dependencies when they change, not on each source edit.

Backend images contain only Alpine and their statically linked Go binary. The
CLI injects its own versioned POSIX supervisor through a ConfigMap, builds the
changed Bazel binary, stages it in a writable volume, restarts the process, and
waits for readiness. If readiness fails it restores the previous binary.

Every diagnostic response includes service, branch, source revision, and pod
identity. Route keys combine normalized owner and branch values with a hash, so
simultaneous users and identically named branches do not collide.
