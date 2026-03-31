# Siftyan

**Siftyan** is a local-first CLI tool written in Go that scans your project's dependency tree and detects software license conflicts — before they become legal problems.

It supports **npm** (`package-lock.json`) and **pip** (`requirements.txt`) ecosystems, identifies incompatible license combinations relative to your distribution model, and explains conflicts in plain English with actionable suggestions.

---

## Features

- **Automatic lockfile detection** — finds `package-lock.json` and `requirements.txt` in the current directory
- **License conflict detection** — identifies Copyleft Propagation, Network Copyleft, Linking Exceptions, and License Ambiguity
- **Distribution model awareness** — conflict severity adapts to how you ship software (`saas`, `binary`, or `internal`)
- **PyPI enrichment** — fetches live license data for pip packages via the PyPI API (concurrent, cached)
- **Interactive HTML report** — D3.js dependency graph with clickable, conflict-highlighted nodes
- **Terminal summary** — instant results directly in your shell
- **Dev dependency control** — optionally include or exclude dev dependencies from the scan

---

## Installation

**Prerequisites:** Go 1.21+

```bash
git clone https://github.com/yourusername/siftyan.git
cd siftyan
go build -o siftyan cmd/main.go
```

---

## Usage

Run inside any project directory:

```bash
./siftyan scan
```

### Options

| Flag | Default | Description |
|---|---|---|
| `--model`, `-m` | `internal` | Distribution model: `saas`, `binary`, or `internal` |
| `--report`, `-r` | _(none)_ | Output path for an HTML report (e.g. `report.html`) |
| `--include-dev` | `false` | Include development dependencies in the scan |

### Examples

```bash
# Scan with binary distribution model
./siftyan scan --model binary

# Generate an HTML report
./siftyan scan --report report.html

# Include dev dependencies
./siftyan scan --include-dev

# Full scan with all options
./siftyan scan --model saas --report report.html --include-dev
```

---

## Distribution Models

Siftyan adjusts conflict severity based on how you distribute your software:

| Model | Description |
|---|---|
| `internal` | Internal tooling only — most restrictive licenses are low risk |
| `saas` | Hosted service — AGPL is high risk even without binary distribution |
| `binary` | Shipped to end users — GPL and AGPL are high risk due to copyleft propagation |

---

## Conflict Types

| Type | Description |
|---|---|
| **Copyleft Propagation** | Strong copyleft license (GPL) may require your entire project to adopt the same license |
| **Network Copyleft** | AGPL extends copyleft to network-accessible software |
| **Copyleft with Linking Exception** | LGPL requires dynamic linking to use the linking exception |
| **License Ambiguity** | Package has no license information — legal risk cannot be assessed automatically |

---

## Example Output

```
Scanning package-lock.json...

WARNING: Found 2 license conflicts

CONFLICT 1 — Copyleft Propagation
Path: my-app → libA
What this means: Strong Copyleft (GPL) in a binary distribution triggers the 'viral' clause.
Impact: HIGH
Suggested actions:
  - Switch project license to GPL
  - Replace dependency

CONFLICT 2 — Network Copyleft
Path: my-app → libB
What this means: Network Copyleft (AGPL) is a high risk in distributed software.
Impact: HIGH
Suggested actions:
  - Replace with a permissive alternative
  - Consult legal
```

---

## HTML Report

The HTML report includes a full interactive dependency graph (powered by D3.js) with conflict nodes highlighted in red, plus a detailed breakdown of every conflict found.

```bash
./siftyan scan --report report.html
open report.html
```

---

## Project Structure

```
siftyan/
├── cmd/
│   └── main.go               # CLI entry point (Cobra)
├── internal/
│   ├── engine/
│   │   ├── conflict.go       # Conflict detection logic and observer pattern
│   │   ├── model.go          # Distribution model rules
│   │   └── spdx.go           # SPDX license registry (singleton)
│   ├── enricher/
│   │   └── pypi.go           # Concurrent PyPI metadata fetcher
│   ├── parser/
│   │   ├── factory.go        # Parser factory and lockfile detection
│   │   ├── npm.go            # npm lockfile parser (v3)
│   │   ├── pip.go            # pip requirements.txt parser
│   │   ├── normalize.go      # License string normalization to SPDX
│   │   └── types.go          # Dependency tree types and builder
│   └── report/
│       ├── terminal.go       # Terminal renderer
│       ├── html.go           # HTML renderer
│       └── report.html       # Embedded report template
└── go.mod
```

---

## Supported Licenses

Siftyan maps licenses to their SPDX categories for conflict analysis:

| Category | Examples |
|---|---|
| Permissive | MIT, Apache-2.0, BSD-2-Clause, BSD-3-Clause, ISC, CC0-1.0 |
| Weak Copyleft | LGPL-2.1, LGPL-3.0, MPL-2.0, EPL-2.0, EUPL-1.2 |
| Strong Copyleft | GPL-2.0, GPL-3.0 |
| Network Copyleft | AGPL-3.0 |

---

## Roadmap

- [ ] Support for Maven, Cargo, and Go Modules
- [ ] CI/CD and GitHub Actions integration
- [ ] Automatic dependency replacement suggestions
- [ ] Web dashboard

---

## License

This project is licensed under the **GNU General Public License v3.0** or later. See `COPYING` for details.
