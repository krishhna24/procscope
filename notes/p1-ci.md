# P1.T3.S3.M2 — proving the CI gate catches red

## Prediction (written before the first push)

- Branch: `ci/red-test`
- Deliberate failure: `internal/sample/red_test.go`, `TestIntentionalRed`, which calls `t.Fatal`.
- Job I expect to fail: **`test`** (`make test-race`). The test body fails, so `go test -race ./...` exits non-zero.
- Jobs I expect to stay green:
  - **`lint`**: golangci-lint compiles the tests but does not execute them, and the new file is lint-clean.
  - **`vuln`**: govulncheck only looks at the dependency graph and reachable code, not test results.
- Predicted time-to-signal: **about 2–4 minutes** from `git push` to the red check appearing on the PR (checkout, `setup-go`, `make tools` on the lint job, build/test, then GitHub reporting the check).

## Actual red run

- Red run: _pending_
- Job that failed: _pending_
- Failing log line: _pending_
- Actual time-to-signal: _pending_

## Actual green run (failing test removed)

- Green run: _pending_

## `make vuln` and the three govulncheck counts

Commands run locally:

```text
$ GOTOOLCHAIN=go1.26.0 make vuln
...
Your code is affected by 0 vulnerabilities.
This scan also found 7 vulnerabilities in packages you import and 40
vulnerabilities in modules you require, but your code doesn't appear to call
these vulnerabilities.

$ make vuln
No vulnerabilities found.
```

The three counts, in plain terms:

1. **"Your code is affected by …"** — the call-graph result. A known-vulnerable
   symbol is actually reachable from code this module builds: some function we
   call (directly or through our own call chain) is the vulnerable one. This is
   the only count that describes a real, reachable risk in this program.
2. **"Packages you import"** — packages that appear in our import graph and
   contain *some* known vulnerability, even if the vulnerable symbol is never
   called. The package is present, but the dangerous function may be unreachable.
3. **"Modules you require"** — every module in the module graph (including
   transitive requirements and the standard library pseudo-module) that has a
   known vulnerability anywhere in it. This is the broadest, least actionable
   bucket: a module can be required without the vulnerable package ever being
   imported or called.

**Which one should fail CI:** only **"your code is affected by"** (reachable
vulnerabilities). That is the one that says an attacker could actually reach a
known-bad function through this binary. Failing on "packages you import" or
"modules you require" would block on code the program never calls, producing
constant false alarms and pressure to pin/ignore unrelated dependencies. The
smaller counts are still useful as an inventory — upgrade when practical — but
they should not be a merge gate by themselves.

## What deleting the `toolchain` line would change

`go.mod` currently has both `go 1.26.0` and `toolchain go1.26.10`. CI reads the
version through `actions/setup-go` with `go-version-file: go.mod`. If the
`toolchain` line were deleted, CI would fall back to the `go` directive and use
**go1.26.0** instead of go1.26.10.

For `make vuln`, the toolchain decides which standard-library version
govulncheck measures. Under go1.26.0 the scan found **7 vulnerabilities in
packages we import and 40 in modules we require** — all of them standard-library
issues already fixed in later 1.26.x patch releases (the verbose output shows
`Found in: stdlib@go1.26`, `Fixed in: stdlib@go1.26.1`, etc.). Under the pinned
go1.26.10 toolchain the same scan prints **"No vulnerabilities found."**
Reachable vulnerabilities stayed at zero either way, so `make vuln` still exits
0; what changes is that the older toolchain surfaces a large patch-level stdlib
inventory that the pinned toolchain has already fixed.
