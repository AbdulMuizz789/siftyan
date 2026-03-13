# Siftyan

## 📌 Project Overview
**Siftyan** is a local **CLI (Command-Line Interface)** tool written in **Go** that scans a project's dependencies to detect **software license conflicts**.

It analyzes dependencies from package managers such as:

- Python (`pip`)
- Node.js (`npm`)

Siftyan identifies license types used by each dependency and detects **incompatible license combinations**. It explains conflicts in **simple English**, shows the **cause of the conflict**, and suggests **possible solutions**.

The tool outputs:

- A **terminal summary**
- A **detailed HTML report**
- An **interactive dependency graph**

---

## 🎯 Problem Statement
Modern software projects rely heavily on third-party libraries. These libraries come with different **software licenses** (MIT, GPL, Apache, BSD, etc.).

Some licenses have restrictions that may **conflict with other licenses or the project's license**, which can cause legal or distribution issues.

Developers often detect these conflicts **very late in development**.

**Siftyan solves this problem** by automatically scanning dependencies and clearly explaining conflicts.

---

## 🚀 Features

- 🔍 **Dependency Scanner**
  - Detects dependencies from project files such as:
    - `requirements.txt`
    - `package.json`

- 📜 **License Detection**
  - Identifies licenses used by each dependency.

- ⚠️ **Conflict Detection**
  - Finds incompatible licenses.

- 💬 **Plain English Explanation**
  - Explains license conflicts in simple language.

- 🔎 **Cause Identification**
  - Shows why the conflict happened.

- 🛠 **Suggested Solutions**
  - Suggests actions such as replacing dependencies.

- 🖥 **Terminal Summary**
  - Displays quick results directly in CLI.

- 🌐 **HTML Report**
  - Generates a detailed report with:
    - dependency list
    - license information
    - conflict explanation
    - interactive dependency graph

---

## 🏗 System Architecture

```
Project Directory
        │
        ▼
Dependency Scanner
(requirements.txt / package.json)
        │
        ▼
License Detector
        │
        ▼
Conflict Analyzer
        │
        ▼
Output Generator
   ├── CLI Summary
   └── HTML Report + Dependency Graph
```

---

## 🛠 Tech Stack

- **Language:** Go (Golang)
- **CLI Framework:** Cobra (optional)
- **Graph Visualization:** D3.js / Go HTML templates
- **Dependency Parsing:** JSON / Text parsing

---

## ⚙️ Installation

### 1️⃣ Clone the repository

```bash
git clone https://github.com/yourusername/siftyan.git
cd siftyan
```

### 2️⃣ Build the CLI tool

```bash
go build -o siftyan
```

### 3️⃣ Run the tool

```bash
./siftyan scan
```

---

## ▶️ Usage

Run Siftyan inside any project directory.

```bash
siftyan scan
```

Example output:

```
Scanning project dependencies...

Found 28 dependencies
Detected 2 license conflicts

Conflict 1:
Dependency: library-x
License: GPL-3.0
Project License: MIT

Explanation:
GPL requires derivative works to also use GPL license.

Suggested Fix:
Replace the dependency or change project license.
```

---

## 📊 Generate HTML Report

To generate a detailed HTML report:

```bash
siftyan scan --report
```

This report includes:

- Full dependency tree
- License details
- Conflict explanations
- Suggested fixes
- Interactive dependency graph

---

## 📂 Project Structure

```
siftyan/
│
├── cmd/
│   └── root.go
│
├── scanner/
│   ├── pip_scanner.go
│   └── npm_scanner.go
│
├── license/
│   └── detector.go
│
├── analyzer/
│   └── conflict_checker.go
│
├── report/
│   ├── cli_output.go
│   └── html_report.go
│
├── web/
│   └── graph_template.html
│
├── main.go
└── README.md
```

---

## 🔮 Future Improvements

- Support for more ecosystems:
  - Maven
  - Cargo
  - Go Modules
- CI/CD integration
- GitHub Action support
- Automatic dependency replacement suggestions
- Web dashboard

---

## 🤝 Contributing

Contributions are welcome.

Steps:

1. Fork the repository
2. Create a new branch
3. Commit your changes
4. Submit a Pull Request

---

## 📄 License
This project is licensed under the **MIT License**.

---

## 👨‍💻 Author
Siftyan is designed to help developers **detect license conflicts early** and maintain **safe open-source compliance**.
