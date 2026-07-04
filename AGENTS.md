# mkhtpasswd

## Goal

Create a small desktop application in Go, compilable for Linux, macOS and Windows, with a simple graphical interface to generate entries compatible with `.htpasswd` files.

The main goal is to let users generate an encrypted password locally, without having to send a plain-text password to the administrator.

The program must also support a command-line mode for automation and script integration.

---

## Scope

The program must:

* have a simple graphical interface;
* allow entering a username;
* allow entering a password;
* have a show/hide password button;
* generate a `.htpasswd`-compatible line;
* allow copying the result to the clipboard;
* have a random password generator;
* allow configuring basic random password parameters;
* work identically in GUI and CLI modes.

---

## Output format

The output must be:

```text
user:hash
```

Example:

```text
paulo:$2y$10$xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

---

## Algorithm

Use **bcrypt** (preferred).

Library:

```go
golang.org/x/crypto/bcrypt
```

Default cost: **10**.

Cost may be parameterized via CLI in the future.

---

## Graphical interface

The interface must contain only:

* Username field;
* Password field;
* Show/Hide password button;
* Generate random password button;
* Random password options:

  * length;
  * lowercase letters;
  * uppercase letters;
  * numbers;
  * symbols;
* Generate `.htpasswd` entry button;
* Read-only result field;
* Copy button.

The interface must be extremely simple.

The user must be able to use it without reading any documentation.

### Recommended UI details (simple, good UX)

- Use `widget.PasswordEntry` for the password field (native show/hide toggle included).
- Use `container.NewForm` or `layout.FormLayout` for clean alignment.
- Use `widget.SpinEntry` (or `widget.Entry` with a validator) for password length.
- Use `widget.Card` to group random password options.
- Disable the "Generate" button when required fields are empty.
- Show brief feedback ("Copied!") on the Copy button for ~2 seconds.

No business logic belongs in the UI package.

---

## GUI library

Use:

```text
Fyne v2.x
```

Reasons:

* written in Go;
* cross-platform;
* mature;
* simple;
* produces binaries for Linux, macOS and Windows.

---

## CLI mode

In addition to the GUI, the program must work entirely via CLI.

No arguments:

```bash
mkhtpasswd
```

Opens the GUI.

With enough arguments:

```bash
mkhtpasswd --user paulo --password "my password"
```

Runs in text-only mode.

Another example:

```bash
mkhtpasswd \
    --user paulo \
    --random \
    --length 20 \
    --symbols
```

---

## Flags

```text
--gui
    Force opening the GUI.

--user
    Username.

--password
    Password to encrypt.

--random
    Generate a random password.

--length
    Password length.

--lower
    Use lowercase letters.
    Default: true

--upper
    Use uppercase letters.
    Default: true

--numbers
    Use numbers.
    Default: true

--symbols
    Use symbols.
    Default: false

--show-password
    Display the generated password when used with --random.

--cost
    bcrypt cost.
    Default: 10

--version
    Show version.
```

---

## CLI security

When `--random` is used, the generated password **must not** be displayed, unless the user explicitly provides:

```text
--show-password
```

This prevents accidental password exposure in terminal history, logs and pipelines.

---

## Technical requirements

* Go 1.26 or later;
* official support for Linux;
* official support for macOS;
* official support for Windows;
* no Internet communication;
* no telemetry;
* no data collection;
* no password writing to disk;
* no logs containing passwords.

---

## Edge cases

- **Username containing `:`**: the colon character is the `.htpasswd` field separator. Reject or sanitize it with a clear error.
- **Password too long**: bcrypt silently truncates at 72 bytes. Document this behavior; optionally warn the user.
- **Non-UTF-8 characters in username**: handle gracefully. Prefer byte-level validation.

---

## Use of modern Go features

Always use modern language features when they make sense.

Prioritise:

* standard library;
* idiomatic code;
* explicit error handling;
* simple APIs;
* low dependency count.

Avoid compatibility with old Go versions.

The project must keep up with language evolution.

Whenever possible:

* use features introduced in recent Go releases;
* periodically review code to replace custom implementations with newly available standard library equivalents;
* remove legacy code when simpler, more idiomatic official alternatives exist;
* avoid maintaining custom solutions for problems already solved by the language or standard library.

Continuous modernisation is part of normal project maintenance.

---

## Testing

Every business rule must have automated tests.

Tests must cover:

* hash generation;
* `.htpasswd` line assembly;
* empty username;
* empty password;
* random password generation;
* password length;
* character set restrictions;
* each generation option.

Run regularly:

```bash
go test ./...
```

Also run:

```bash
go test -race ./...
```

and

```bash
go vet ./...
```

Whenever applicable.

---

## Directory structure

```text
mkhtpasswd/
├── AGENTS.md
├── README.md
├── LICENSE
├── .gitignore
├── .woodpecker.yml
├── go.mod
├── go.sum
├── cmd/
│   └── mkhtpasswd/
│       └── main.go
├── internal/
│   ├── cli/
│   │   └── cli.go
│   ├── ui/
│   │   └── app.go
│   ├── htpasswd/
│   │   ├── bcrypt.go
│   │   └── bcrypt_test.go
│   └── password/
│       ├── generator.go
│       └── generator_test.go
├── dist/
└── screenshots/
```

---

## Package organisation

### internal/htpasswd

Responsible for hash generation and output assembly.

```go
func Generate(username, password string) (string, error)
```

### internal/password

Responsible for password generation.

```go
type Options struct {
    Length int

    Lowercase bool
    Uppercase bool
    Numbers   bool
    Symbols   bool
}

func GenerateRandomPassword(opts Options) (string, error)
```

### internal/ui

Exclusively responsible for the graphical interface.

No business logic in this package.

### internal/cli

Only responsible for argument parsing.

No business logic in this package.

---

## Architecture

All business logic must remain decoupled from the GUI.

The GUI and CLI must use exactly the same internal packages.

Never duplicate logic.

---

## Build

### Simple build

```bash
go build ./cmd/mkhtpasswd
```

### Cross-compilation with fyne-cross

[fyne-cross](https://github.com/fyne-io/fyne-cross) is the recommended tool for cross-platform builds. It uses Docker images with all required toolchains and system dependencies (OpenGL, X11, etc.).

Target platforms:

```text
linux/amd64
linux/arm64
darwin/amd64
darwin/arm64
windows/amd64
```

Example commands:

```bash
fyne-cross linux -output mkhtpasswd
fyne-cross darwin -output mkhtpasswd
fyne-cross windows -output mkhtpasswd
```

A single `fyne-cross` invocation can build multiple targets at once.

---

## CI/CD

The CI pipeline runs on [Codeberg's Woodpecker CI](https://ci.codeberg.org) (linux/amd64).

Pipeline file: `.woodpecker.yml`

Pipeline steps:

1. **Lint** — `go vet ./...` (and optionally `staticcheck`).
2. **Test** — `go test -race -cover ./...`.
3. **Build** — cross-compile all targets using `fyne-cross` via Docker.
4. **Release** — on tag push, attach binaries to the release.

All targets are built from a single linux/amd64 runner using `fyne-cross` Docker images.

> Access to Codeberg's Woodpecker instance requires manual approval via [this form](https://codeberg.org/Codeberg-e.V./requests/issues/new?template=ISSUE_TEMPLATE%2fWoodpecker-CI.yaml).

---

## Development workflow

```bash
# Run directly
go run ./cmd/mkhtpasswd

# Run tests
go test -race -cover ./...

# Lint
go vet ./...

# Build locally
go build ./cmd/mkhtpasswd
```

---

## .gitignore

```
/mkhtpasswd
/mkhtpasswd.exe
/dist/
/screenshots/
*.tar.gz
*.zip
```

---

## README

The README must explain:

* goal;
* installation;
* compilation;
* GUI mode;
* CLI mode;
* examples;
* security notice;
* license.

---

## Project principles

### Simplicity first

Always prefer the simplest solution.

Avoid:

* premature abstractions;
* generic interfaces;
* unnecessary frameworks;
* unnecessary design patterns.

### Standard library first

Before adding any external dependency, check whether the standard library solves the problem.

Dependencies must be justified.

### Idiomatic code

All code must follow Go community style.

Run regularly:

* gofmt;
* go test;
* go vet.

### Small APIs

Packages should expose only what is necessary.

Always prefer small, simple APIs.

### Testability

All logic must be easily testable.

The GUI must never contain business logic.

### Portability

The program must work identically on Linux, macOS and Windows.

Avoid OS-specific APIs.

### Security

Never log:

* passwords;
* temporary hashes;
* clipboard contents.

Use:

* crypto/rand;
* bcrypt (golang.org/x/crypto/bcrypt).

When possible, clear buffers containing passwords after use.

### Interface

The interface must remain extremely simple.

Do not add complex configuration screens.

### CLI

The CLI must be suitable for automation.

Expected behaviour:

* normal output to stdout;
* errors to stderr;
* exit code 0 on success;
* non-zero exit code on error.

### Evolution

The project must remain a small `.htpasswd` entry generator.

Avoid turning it into a full user management or authentication system.

New features must remain optional.

### Complexity control

Whenever a new feature:

* adds a significant dependency;
* introduces a new architecture;
* or requires roughly more than 200 lines of code,

the implementation must clearly justify that decision in documentation or the pull request.

Simplicity is a requirement, not a preference.

### Continuous refactoring

During changes to existing code, small related improvements are encouraged, provided they:

* do not change program behaviour;
* reduce complexity;
* improve readability;
* increase test coverage;
* remove dead or duplicated code;
* bring code closer to current Go practices.

Avoid large refactorings unrelated to the requested task. Refactoring should be done in separate commits, as a suggestion presented before the actual changes.

---

## Out of scope

Do not implement initially:

* `.htpasswd` file editing;
* multiple users;
* import/export;
* authentication management;
* LDAP;
* OAuth;
* HTTP server;
* password storage;
* Internet synchronisation.

---

## Success criteria

The project is considered complete when any user can:

1. enter a username and password;
2. generate a random password (optional);
3. obtain a `.htpasswd`-compatible line;
4. copy the result;
5. run it via both GUI and CLI;
6. compile without changes on Linux, macOS and Windows;
7. have test coverage for all business logic.