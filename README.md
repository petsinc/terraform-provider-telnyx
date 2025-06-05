# Telnyx Terraform Provider

This is a Terraform provider for Telnyx. It is currently in **alpha** status and **not recommended for production use**.

---

## 📖 Documentation

Official documentation is hosted on the [Terraform Registry](https://registry.terraform.io/providers/petsinc/telnyx/latest/docs).

---

## ⚙️ Prerequisites

Ensure the following tools are installed before starting development:

- [direnv](https://direnv.net/docs/installation.html) — for managing environment variables automatically.
- [asdf](https://asdf-vm.com/guide/getting-started.html) — for managing language versions.  
  ⚠️ Install using the official `git` method to avoid issues with pre-commit hooks locating `~/.asdf/bin`.

> **Note:** These tools require a UNIX-like shell (bash, zsh). On Windows, use WSL or Git Bash. If using PyCharm, configure the terminal to use Git Bash.

---

## 📚 Environment Configuration

- Create a `.env` file at the project root and set the following variable:

  ```bash
  TELNYX_API_KEY=your_api_key_here
  ```

- Run `direnv allow` in the project directory to load environment variables automatically.

---

## 👨‍💻 Development Patterns

### 🔧 Local Development

- Run tests:

  ```bash
  just test
  ```

- Format Go code:

  ```bash
  just format-go
  ```

- Format Terraform (HCL) code:

  ```bash
  just format-hcl
  ```

- Environment variables are managed using `.env` and `direnv`.

---

### 🐳 Docker-Based Development

Docker provides a consistent and isolated environment, particularly useful when local setups vary or when running CI pipelines.

- **Build Docker Image:**

  ```bash
  just build-docker
  ```

- **Run Tests in Docker:**

  ```bash
  just test-docker
  ```

- **Generate Documentation in Docker:**

  ```bash
  just gen-docs-docker
  ```

- **Open an Interactive Development Shell (Alpine-based):**

  ```bash
  just dev-docker
  ```

> ⚠️ Docker workflows are ephemeral; no build artifacts or state persist between runs.

---

## 📋 Common Commands Cheat Sheet

| Task              | Local Command     | Docker Command         |
| ----------------- | ----------------- | ---------------------- |
| Run Tests         | `just test`       | `just test-docker`     |
| Format Go Code    | `just format-go`  | N/A                    |
| Format HCL Code   | `just format-hcl` | N/A                    |
| Generate Docs     | N/A               | `just gen-docs-docker` |
| Interactive Shell | N/A               | `just dev-docker`      |

---

## 📝 Conventional Commits & Semantic Release

This project follows [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) to power [Semantic Release](https://semantic-release.gitbook.io/semantic-release/).

- **Commit Format:**

  ```
  <type>[optional scope]: <description>
  ```

  Examples:

  - `feat: Add new resource for SIP trunking`
  - `fix(provider): Resolve issue with API key handling`

- **Pull Requests:**
  - Name PRs using the conventional commit format. This ensures clean, semantically meaningful commits when merging.

---

## ✅ Pre-commit Hooks

- This project uses Husky and Lint-Staged to enforce formatting and validate commit messages.
- If a commit fails, ensure your message follows the Conventional Commits format and code formatting is correct.

---

## 📦 Project Automation with `just`

This project uses [`just`](https://github.com/casey/just) for repeatable development workflows. Review the `justfile` for additional commands.

---

## 👥 Author

This project is maintained by **Patient Engagement Technologies**.
