# Implementation Plan

## Phase 1 — Foundation

- `go mod init`, add `golang.org/x/crypto` dependency
- Create directory structure (`cmd/`, `internal/...`)
- Define `Options` struct in `internal/password`
- Define `Generate()` signature in `internal/htpasswd`
- `.gitignore`

**Success criterion:** `go build ./...` compiles (packages can be empty/stub).

---

## Phase 2 — `internal/password` (password generator)

- Implement `GenerateRandomPassword(opts Options) (string, error)`
- Tests: length, charset, each option individually
- **No internal package dependencies**

**Success criterion:** `go test ./internal/password/...` passes.

---

## Phase 3 — `internal/htpasswd` (bcrypt hash)

- Implement `Generate(username, password string) (string, error)`
- Tests: valid hash, correct line format, empty username, empty password, username with `:`
- **External dependency:** `golang.org/x/crypto/bcrypt`

**Success criterion:** `go test ./internal/htpasswd/...` passes.

---

## Phase 4 — `internal/cli` (argument parser)

- Implement flag parsing (`flag` package)
- Invoke `htpasswd.Generate()` and `password.GenerateRandomPassword()` based on args
- Test parsing (not business logic — just flag reading)
- **Dependencies:** `internal/htpasswd`, `internal/password`

**Success criterion:** `mkhtpasswd --user x --password y` prints `x:$2y$...` to stdout.

---

## Phase 5 — `internal/ui` (graphical interface)

- Implement Fyne window with all widgets
- Wire events to `htpasswd` and `password` packages
- **No business logic** — only call internal packages
- **Dependencies:** `internal/htpasswd`, `internal/password`, `fyne.io/fyne/v2`

**Success criterion:** `go run ./cmd/mkhtpasswd` opens the window and generates entries.

---

## Phase 6 — `cmd/mkhtpasswd/main.go` (entry point)

- Routing: no args → GUI, with args → CLI
- Handle exit codes and output
- Integrate `internal/cli` and `internal/ui`
- **Minimal code** — only decides what to run

**Success criterion:** both GUI and CLI modes work.

---

## Phase 7 — CI/CD and Build

- `.woodpecker.yml` with lint → test → build (fyne-cross)
- Verify cross-build manually with `fyne-cross`
- README with usage and compilation instructions
- License

**Success criterion:** green pipeline on Codeberg.