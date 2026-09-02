# Conformance and operations

## Build and render

```sh
bazel test --config=agent //...
bazel build --config=agent //...
kubectl kustomize k8s/dev >/dev/null
```

## Base and routing checks

```sh
curl --cacert path/to/lab-root-ca.crt https://dev-cli-playground.dev.lab/api/diagnostic
curl --cacert path/to/lab-root-ca.crt -c /tmp/a.cookie -b /tmp/a.cookie \
  'https://dev-cli-playground.dev.lab/__dev/select?branch=<route-key>'
curl --cacert path/to/lab-root-ca.crt -b /tmp/a.cookie \
  https://dev-cli-playground.dev.lab/api/diagnostic
```

Expected cases are: no cookie uses both base services; a hello-api-only overlay
uses overlay hello-api and base echo-api; adding echo-api makes both hops use
the overlay; a second cookie jar remains base or selects a different route.
Deleting one overlay causes that service alone to return to base on the next
request.

`dev status` reports route ownership and expiry. `dev down` deletes only the
current owner/route selector. `dev cleanup` removes expired managed routes.
