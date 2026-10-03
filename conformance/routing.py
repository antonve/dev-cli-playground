"""Read-only live conformance: two independent branch hosts, both HTTP hops."""

import argparse
import concurrent.futures
import json
import urllib.request


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--host", default="https://dev-cli-playground.dev.lab")
    for key in ("a", "b"):
        parser.add_argument(f"--route-{key}", required=True)
        parser.add_argument(f"--branch-{key}", required=True)
        parser.add_argument(f"--services-{key}", default="")
    args = parser.parse_args()

    def check(case):
        name, route, forged, branch, services, normalized = case
        headers = {"Cache-Control": "no-cache"}
        host = args.host.rstrip("/")
        if route:
            host = host.replace("://", "://" + route + ".", 1)
        if name == "stale-cookie-base":
            headers["Cookie"] = "dev_branch=" + args.route_a
        if forged:
            headers["x-dev-branch"] = forged
        responses = []
        for path in ("/__web", "/api/diagnostic"):
            request = urllib.request.Request(host + path, headers=headers)
            with urllib.request.urlopen(request, timeout=5) as response:
                responses.append(json.load(response))
        web, api = responses
        for identity in (web, api["self"], api["echo"]):
            expected = branch if identity["service"] in services else "base"
            assert identity["branch"] == expected, (name, identity, expected)
            assert identity["revision"] and identity["pod"], (name, identity)
        for identity in (api["self"], api["echo"]):
            assert identity["observedRoute"] == normalized, (name, identity, normalized)
        return name

    cases = [
        ("base", "", "", "base", set(), ""),
        ("forged-only", "", args.route_a, "base", set(), ""),
        ("stale-cookie-base", "", args.route_b, "base", set(), ""),
        ("session-a", args.route_a, args.route_b, args.branch_a, set(args.services_a.split(",")), args.route_a),
        ("session-b", args.route_b, args.route_a, args.branch_b, set(args.services_b.split(",")), args.route_b),
    ]
    for case in cases:
        check(case)
    # Each request independently selects its hostname.
    with concurrent.futures.ThreadPoolExecutor(max_workers=4) as pool:
        results = list(pool.map(check, cases * 10))
    print(json.dumps({"passed": True, "cases": [c[0] for c in cases], "concurrent_requests": len(results) * 2}))


if __name__ == "__main__":
    main()
