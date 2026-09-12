"""Read-only v0.3+ gateway deep-link acceptance, with two independent sessions.

Start two CLI environments with only hello-web and echo-api affected. This test
changes browser cookies only; it never creates, patches or deletes workloads.
"""

import argparse
import concurrent.futures
import http.cookiejar
import json
import urllib.parse
import urllib.request


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--host", default="https://dev-cli-playground.dev.lab")
    for key in ("a", "b"):
        parser.add_argument(f"--route-{key}", required=True)
        parser.add_argument(f"--branch-{key}", required=True)
    args = parser.parse_args()
    host = args.host.rstrip("/")

    def session():
        jar = http.cookiejar.CookieJar()
        return urllib.request.build_opener(urllib.request.HTTPCookieProcessor(jar)), jar

    def request(client, path, method="GET", **headers):
        req = urllib.request.Request(host + path, method=method, headers=headers)
        with client.open(req, timeout=10) as response:
            return response.headers, response.read(), response.url

    def diagnostic(client, route, branch, selection=None):
        path = "/api/diagnostic"
        if selection is not None:
            path += "?tab=profile&" + urllib.parse.urlencode({"dev-branch": selection})
        headers, body, url = request(client, path, **{"x-dev-branch": "forged"})
        result = json.loads(body)
        assert result["self"]["branch"] == "base", result
        assert result["echo"]["branch"] == branch, result
        for hop in (result["self"], result["echo"]):
            assert hop["observedRoute"] == route, result
            assert hop["pod"] and hop["revision"], result
        if selection is not None:
            # No redirect or path/query rewrite, including on the first request.
            assert url == host + path, url
            assert headers.get("Cache-Control") == "no-store", headers
            cookies = headers.get_all("Set-Cookie", [])
            assert len(cookies) == 1, cookies
            assert "Path=/" in cookies[0] and "Secure" in cookies[0] and "SameSite=Lax" in cookies[0], cookies
            assert "Domain=" not in cookies[0], cookies
            if selection == "base":
                assert "Max-Age=0" in cookies[0], cookies
            else:
                assert "dev_branch=" + route + ";" in cookies[0], cookies
        return result

    a, jar_a = session()
    b, jar_b = session()
    diagnostic(a, "", "base")
    diagnostic(b, "", "base")
    result_a = diagnostic(a, args.route_a, args.branch_a, args.route_a)
    result_b = diagnostic(b, args.route_b, args.branch_b, args.route_b)
    assert result_a["echo"]["pod"] != result_b["echo"]["pod"], "owners collided"
    diagnostic(a, args.route_a, args.branch_a)
    diagnostic(b, args.route_b, args.branch_b)

    # Switch A to B with a conflicting real A cookie, not just an unknown key.
    switched = diagnostic(a, args.route_b, args.branch_b, args.route_b)
    assert switched["echo"]["pod"] == result_b["echo"]["pod"]
    diagnostic(a, args.route_b, args.branch_b)
    diagnostic(a, "", "base", "base")
    assert not list(jar_a), "clear did not remove session cookie"
    diagnostic(a, "", "base")
    diagnostic(b, args.route_b, args.branch_b)

    # An ordinary deep application page sets the cookie without a menu/handler.
    path = "/settings?tab=profile&" + urllib.parse.urlencode({"dev-branch": args.route_a})
    headers, body, url = request(a, path, Accept="text/html")
    assert url == host + path and b"/@vite/client" in body, "first page not overlay HTML"
    assert headers.get_all("Set-Cookie"), headers
    _, body, _ = request(a, "/__web")
    assert json.loads(body)["branch"] == args.branch_a, body
    diagnostic(a, args.route_a, args.branch_a)

    # HEAD also selects/clears; mutation methods do not set selection cookies.
    request(a, "/api/diagnostic?dev-branch=base", method="HEAD")
    assert not list(jar_a)
    headers, _, _ = request(a, "/api/diagnostic?dev-branch=" + args.route_a, method="POST")
    assert not headers.get_all("Set-Cookie") and not list(jar_a)
    diagnostic(a, "", "base")
    diagnostic(a, args.route_a, args.branch_a, args.route_a)

    # Unknown selectors fall through to the current cookie, not another branch.
    _, body, _ = request(a, "/api/diagnostic?dev-branch=missing-environment")
    assert json.loads(body)["echo"]["pod"] == result_a["echo"]["pod"]

    def concurrent_check(which):
        client, _ = session()
        route, branch = (args.route_a, args.branch_a) if which == "a" else (args.route_b, args.branch_b)
        diagnostic(client, route, branch, route)
        diagnostic(client, route, branch)

    with concurrent.futures.ThreadPoolExecutor(max_workers=4) as pool:
        list(pool.map(concurrent_check, ["a", "b"] * 20))
    print(json.dumps({"passed": True, "concurrent_requests": 80,
                      "cases": ["base", "first-request-query", "forged-header", "cookie-persistence",
                                "conflicting-cookie-switch", "clear", "independent-session",
                                "deep-page", "HEAD", "POST-no-selection", "unknown-selection",
                                "base-api-to-overlay-echo"],
                      "pods": {"a": result_a["echo"]["pod"], "b": result_b["echo"]["pod"]}}))


if __name__ == "__main__":
    main()
