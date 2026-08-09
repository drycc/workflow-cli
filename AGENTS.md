# AGENTS.md — workflow-cli

> **When you change code patterns, conventions, or architecture, update this file
> in the same commit.** AGENTS.md is the single source of truth for project
> conventions — if it's stale, the next agent will follow wrong patterns.

`drycc` is the command line utility for the Drycc PaaS. This Go project is the
**user-facing tier** of a three-repo stack — it talks to a controller over HTTP
through `controller-sdk-go`. Read this file before editing anything.

## The three-repo stack (read first)

```
controller            (Python/Django, the API server)  ← defines the contract
   ↓ HTTP REST
controller-sdk-go     (Go client SDK)                  ← wraps the HTTP calls
   ↓ imported
workflow-cli          (this repo, the `drycc` binary)  ← parses args, renders output
```

- The **controller is the source of truth** for URL shape, request/response
  payloads, and error semantics. When adding a feature, look at the controller's
  `rootfs/api/urls.py` and `rootfs/api/views/` first.
- A new backend endpoint normally touches **all three repos**: add the SDK
  function, then the CLI command, then (if needed) the controller test already
  exists or is added by the controller author.
- These repos live as siblings on this machine under `~/Sources/`. Paths used in
  examples below assume that layout.

## Tech stack

- Go 1.26, module path `github.com/drycc/workflow-cli`
- CLI framework: `github.com/spf13/cobra`
- HTTP client: `github.com/drycc/controller-sdk-go`
- Output: `github.com/olekukonko/tablewriter` (tables), `pkg/coder` (flat YAML)
- i18n: `github.com/chai2010/gettext-go` + `.po`/`.mo` catalogs
- `vendor/` is **gitignored** — it is regenerated locally, never committed

## Directory structure

```
cmd/root.go            Binary entrypoint; builds the root cobra command
internal/
  parser/              cobra command definitions (one file per domain)
  commands/            command implementations (DryccCmd methods) + Commander interface
  completion/          shell-completion helpers (AppCompletion, AddonCompletion, ...)
  loader/              LoadAppSettings: resolves appID + settings + SDK client
  template/            cobra example templating
  plugins/             plugin loader
pkg/
  coder/               YAML encode/decode (flat snake_case, no K8s Manifest wrapper)
  i18n/                translation loader + catalogs (en_US, zh_CN)
  settings/            client config file (~/.drycc/client.json)
  testutil/            test HTTP server + helpers
  git, ssh, logging, webbrowser   small supporting packages
scripts/
  build                release build script
  update-translations.sh   i18n extraction (needs go-xgettext + gettext)
version/               version constants
```

## Command architecture — the 5 layers of a feature

Every domain (addons, apps, certs, gateways, ...) is spread across **five**
files. To add a subcommand you touch all of them. Example for `addons conn`:

| Layer | File | Responsibility |
|-------|------|----------------|
| 1. SDK | `controller-sdk-go/<domain>/<domain>.go` | `Connection(c, appID, name)` → HTTP call |
| 2. command | `internal/commands/addons.go` | `(d *DryccCmd) AddonsConnection(...)` — loads settings, calls SDK, renders output |
| 3. cobra wiring | `internal/parser/addons.go` | `addonsConnectionCommand(cmdr)` — `Use`/`Short`/`Example`/flags; registered in `NewAddonsCommand` via `cmd.AddCommand(...)` |
| 4. interface | `internal/commands/commands.go` | add the method signature to the `Commander` interface |
| 5. root | `cmd/root.go` | `rootCmd.AddCommand(parser.NewAddonsCommand(&cmdr))` (only for new top-level domains, not subcommands) |
| 6. completion | `internal/completion/completion.go` | `<Domain>Completion` struct with `CompletionFunc` (shell tab-complete) |
| 7. coder | `pkg/coder/<domain>.go` | only when the command reads/writes K8s-style YAML manifests |

**Don't forget layer 4.** The `Commander` interface enumerates every command
method; mocks and tests rely on it. Adding a `DryccCmd` method without declaring
it in the interface is the #1 missed step.

## Conventions

### Naming
- Top-level cobra `Use` strings are **full words** unless they are an industry
  term (`tls`, `ps`, `ptypes`, `certs`). `addons`, `gateways`, `volumes`,
  `domains` are all full words. Don't abbreviate casually.
- Go identifiers mirror the SDK/controller exactly: SDK `Connection()` → CLI
  `AddonsConnection()`. Never invent a divergent name (e.g. don't call the SDK
  `Connection` something else in the CLI).
- `addon` / `AddOn` / `addon` — the project spells it **lowercase "addon"** in
  prose and identifiers (`addons`, `AddonInstance`, `AddonClass`). Keep it.
- **All API-related field names and CLI parameters use snake_case.** K8s
  camelCase names (e.g., `livenessProbe`, `httpGet`, `postStart`, `backendRefs`)
  are converted to snake_case by the controller's `CamelCaseJSONField`. CLI
  command arguments and display labels must also use snake_case:
  `liveness_probe`, `http_get`, `post_start`, `backend_refs`, etc.

### i18n — every user-facing string goes through `i18n.T()`
All `Short`/`Long`/flag descriptions/`<name>` arg descriptions must be wrapped:
```go
Short: i18n.T("Show connection details for an addon"),
```
After adding/changing strings, regenerate catalogs:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"   # go-xgettext lives here
go install github.com/gosexy/gettext/go-xgettext@latest   # one-time
bash scripts/update-translations.sh -xkg     # extract + fix .po + build .mo
```
Then translate the new `msgstr ""` entries in
`pkg/i18n/translations/drycc/zh_CN/LC_MESSAGES/cli.po` (English is the source —
`en_US` entries mirror the msgid). Re-run `-g` after editing `.po`. Keep the
technical term "addon" untranslated in zh_CN to match the identifiers.

### Output rendering
- **Lists** -> `d.getDefaultFormatTable([]string{"COL1","COL2"})`, render with
  `table.Append`/`table.Render`. It is whitespace-packed (`NoWhiteSpace`); assert
  with `assert.Contains` in tests, not exact-string equals.
- **Single-object details** -> `pkg/coder` flat YAML (snake_case, no K8s Manifest
  wrapper) via a `*Coder.Encode()`. Coder `Decode` is just `json.Unmarshal`,
  `Encode` is just `marshalYAML` - no field renaming.
- **Connection/secret-style key->value blobs** -> KEY/VALUE table, keys sorted.
- **Mutating ops** (apply/remove/upsert) -> print `"<verb> <obj> ... "`,
  call `progress(d.WOut)` for a spinner, then `done`.

### Shared parser state
`internal/parser/flags.go` holds **package-level** `app`, `version`, `limit`
vars reused across command constructors. Set via `PersistentFlags` / `Flags`.
`commands.ResponseLimit(limit)` resolves `0` → user's configured default.

### Testing
- `pkg/testutil.NewTestServerAndClient()` spins up an `httptest` server + writes
  a temp client config. Register handlers with `server.Mux.HandleFunc(...)`.
- `internal/commands/main_test.go` sets `LANG=LC_ALL=LC_MESSAGES=en_US.UTF-8`
  so `i18n.T()` returns English and assertions are deterministic. **Don't delete
  it** and don't set conflicting locales in command tests.
- Mark new tests `t.Parallel()`.
- Use `github.com/stretchr/testify/assert` (already imported everywhere).
- The whole `commands` package is tested via the `DryccCmd` struct wired to the
  fake server; never hit a real controller in unit tests.

## SDK dependency management (important)

`workflow-cli` depends on `github.com/drycc/controller-sdk-go` by **pseudo-version
commit**. When you change the SDK locally and want to test here:

1. Commit (and usually force-push, since feature commits get amended) the SDK
   change to `controller-sdk-go` first.
2. In workflow-cli: `go get github.com/drycc/controller-sdk-go@<new-sha>`,
   then `go mod vendor`, then `go build ./...`.
3. Commit `go.mod` + `go.sum`. `vendor/` stays untracked (gitignored).

Do **not** leave a `replace github.com/drycc/controller-sdk-go => ../controller-sdk-go`
directive in `go.mod` — it's a dev-only shortcut that breaks everyone else's
build. If you add it temporarily, remove it before committing.

## Build / test / vet commands

```bash
go build ./...                              # compile everything
go vet ./...                                # static checks — run before commit
go test ./...                               # all tests
go test ./internal/commands/... ./pkg/...   # just the layers you touched
go test -run TestAddons -v ./...            # one feature's tests
go run . addons conn --help                 # smoke-test a command's help
LANG=zh_CN.UTF-8 LC_ALL=zh_CN.UTF-8 go run . addons conn --help   # zh check
```

There is no `go generate` step in normal builds. CI (`Makefile`) runs inside a
`go-dev` container and calls `update-translations.sh -g` during `build-binary`.

## Git conventions

- Commit subject format: `type(scope): subject`, e.g.
  `feat(addons): add addons support`, `chore(api): add uid for app`.
  Common types: `feat`, `chore`, `fix`. Scope is the domain (`addons`, `gateway`,
  `cli`, `resources`, ...).
- A single feature spans both `controller-sdk-go` and `workflow-cli`. Land it as
  **one commit per repo** with the **same subject** in both. If you iterate on
  the feature after the initial commit, `git commit --amend --no-edit` to fold
  follow-ups into that one commit rather than stacking new ones.
- Force-pushing an amended feature commit to `controller-sdk-go` is expected and
  accepted here — confirm the target branch first (`main`) and use
  `--force-with-lease`.
- Never commit secrets, `vendor/`, or the dev `replace` directive.

## Adding a new domain end-to-end (checklist)

1. **Controller** — confirm the URL + viewset exist in `controller/rootfs/api/`.
2. **SDK** (`controller-sdk-go`): add types in `api/<domain>.go`, functions in
   `<domain>/<domain>.go`, fixtures+tests in `<domain>/<domain>_test.go`.
   `go build && go test ./<domain>/...` there, commit, force-push.
3. **CLI SDK bump**: `go get ...@<sha> && go mod vendor`.
4. **CLI command** (`internal/commands/<domain>.go`): `DryccCmd` methods, load
   settings via `loader.LoadAppSettings`, render output.
5. **CLI parser** (`internal/parser/<domain>.go`): cobra commands, i18n-wrapped
   strings, flag wiring, register subcommands in the domain's `New<Domain>Command`.
6. **Commander interface** (`internal/commands/commands.go`): add every new
   method signature.
7. **Root** (`cmd/root.go`): `rootCmd.AddCommand(parser.New<Domain>Command(&cmdr))`.
8. **Completion** (`internal/completion/completion.go`): if the command takes a
   `<name>` arg that maps to existing resources, add a `<Domain>Completion`.
9. **Coder** (`pkg/coder/<domain>.go`): only if the command reads/writes YAML
   manifests.
10. **i18n**: `bash scripts/update-translations.sh -xkg`, translate zh_CN, `-g`.
11. **Tests**: `internal/commands/<domain>_test.go` (+ `pkg/coder/<domain>_test.go`).
12. **Verify**: `go vet ./... && go test ./...` before committing.

## Common mistakes

- Forgetting to add the method to the `Commander` interface (step 6 above).
- Leaving the `replace` directive in `go.mod` (breaks CI / other clones).
- Hardcoding `LANG` in a test that conflicts with `main_test.go`.
- Committing `vendor/` (it's gitignored — if `git status` shows vendor files,
  your `.gitignore` is wrong).
- Asserting exact table output strings — the tablewriter pads unpredictably; use
  `assert.Contains`.
- Abbreviating a `Use` name that should be a full word (see Naming).
- Adding a `msgstr` but forgetting to re-run `-g` to rebuild the `.mo`.
