# iCloud Ignore (`icignore`)

`icignore` is a small command-line utility for **selectively excluding folders from iCloud synchronization** on macOS.
It leverages the `.nosync` suffix recognized by macOS, while keeping folder access seamless through symbolic links, so apps and Finder keep working as before.

---

## ✨ Features

- Ignore: `icignore ignore <folder>` renames to `<name>.nosync` and creates a symlink back.
- Unignore: `icignore unignore <folder>` removes the symlink and renames back.
- Status: `icignore status <folder>` shows symlink + real directory status.
- Safe-by-default: refuses destructive overwrites and supports `--dry-run`.
- Homebrew-friendly: single static binary with formula and releases.

---

## 🚀 Installation

### Using Homebrew (recommended)

```bash
brew tap mathis-lambert/icloud-ignore
brew install icignore
```

---

## 🔧 Usage

### Exclude a folder

```bash
icignore ignore ~/Documents/Projects
```

Effect:

* The folder is renamed to `Projects.nosync` (ignored by iCloud).
* A symbolic link `Projects` is recreated at the same location → Finder and apps continue to work as before.

---

### Check folder status

```bash
icignore status ~/Documents/Projects
```

Example output:

```
SYMLINK: /Users/username/Documents/Projects -> Projects.nosync
REAL: /Users/username/Documents/Projects.nosync (excluded from iCloud via .nosync suffix)
```

---

### Restore sync

```bash
icignore unignore ~/Documents/Projects
```

Effect:

* Removes the symbolic link.
* Renames `Projects.nosync` back to `Projects`.
* The folder resumes syncing with iCloud.

---

## 🛠 Commands

```
icignore ignore <path>    # Exclude a folder from iCloud sync
icignore unignore <path>  # Restore sync for a folder
icignore status <path>    # Check folder status
```

Global flags:
  `--dry-run`    Print actions without changing anything
  `--verbose`    Extra logging

Version:
  `icignore version`

Exit codes: `0` on success; non-zero on error.

Note: A path with or without `.nosync` is accepted for all commands.

---

## 🧑‍💻 Local Development

Prerequisites:
- Go 1.21+
- Make

Clone and build:

```bash
git clone https://github.com/mathis-lambert/icloud-ignore.git
cd icloud-ignore
make build   # binary at bin/icignore
```

Run locally:

```bash
./bin/icignore --help
./bin/icignore status ~/Documents/Projects
```

Install locally:

```bash
make install   # to /usr/local/bin or /opt/homebrew/bin
```

Checks:

```bash
make fmt vet test
```

Project structure:

```
cmd/icignore/        # CLI entry
internal/icignore/   # Core operations (ignore/unignore/status)
.github/workflows/   # Release (GoReleaser)
```

---

## 📦 Release & Homebrew

Tagged releases are built and published via GoReleaser. The Homebrew tap repository is updated automatically.

- Source repo: `mathis-lambert/icloud-ignore`
- Tap repo: `mathis-lambert/homebrew-icloud-ignore`

Release process:
- Merge changes to `main` and update docs.
- Tag a version: `git tag v0.x.y && git push origin v0.x.y`
- GitHub Actions runs GoReleaser to publish archives and update the tap.

Maintainer setup (once):
- Create a Personal Access Token (classic) with `repo` scope and add it as a secret in this repo named `HOMEBREW_TAP_TOKEN`.
  - It must have write access to both repositories above.
- The workflow `.github/workflows/release.yml` uses this token to publish releases and push the formula to the tap.

Manual release (local):

```bash
GITHUB_TOKEN=<your_pat_with_repo_scope> make release
```

Tap name: `mathis-lambert/icloud-ignore` (backed by repo `homebrew-icloud-ignore`).

---

## 🤝 Contributing

Contributions are welcome! See `CONTRIBUTING.md`. Please follow the `CODE_OF_CONDUCT.md`.

---

## 🔒 Security

If you discover a security or privacy issue, please follow `SECURITY.md` instead of filing a public issue.

---

## ⚖ License

MIT License. See `LICENSE`.

---

## 💡 Quick Example

```bash
# Exclude the "Projects" folder in Documents
icignore ignore ~/Documents/Projects

# Verify
icignore status ~/Documents/Projects
```

Finder will still show `Projects`, but iCloud will no longer sync it 🚫☁️.
