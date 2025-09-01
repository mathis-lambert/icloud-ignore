# Contributing to icignore

Thanks for your interest in contributing! This project aims to be simple, safe, and friendly to improvements. Please take a moment to read this guide.

## Development setup

- Prerequisites: Go 1.21+, Make
- Clone and build:
  - `git clone https://github.com/mathis-lambert/icloud-ignore.git`
  - `cd icloud-ignore`
  - `make build` (binary at `bin/icignore`)
- Run checks locally:
  - `make fmt vet test`

## Branches and commits

- Use feature branches from `main`.
- Keep commits focused and with clear messages.
- Conventional Commits are welcome but not required.

## Pull requests

- Include a clear description of the change and rationale.
- Update documentation and README when behavior changes.
- Ensure `make build` and `make vet` succeed; add tests when practical.

## Code style

- Keep the code small and readable, prefer standard library.
- Cross-check for safe filesystem operations; avoid destructive defaults.

## Reporting issues

- Provide steps to reproduce, expected vs actual behavior, and environment details.
- If it’s a question or idea, label it accordingly.

## Security

If you discover a security or privacy issue, please follow the process in `SECURITY.md` instead of filing a public issue.

Thanks for contributing!

