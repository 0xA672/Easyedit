# Contributing to EasyEdit

Thanks for your interest in EasyEdit! Bug reports, reproducible examples, design discussions, testing, documentation feedback, and code contributions all help the project.

EasyEdit is also an experiment in AI-agent-driven software development. The repository's stated rule is **“No human-written code allowed.”** To preserve that experiment, source-code and test-code changes should be authored by an AI coding agent. Human participation is absolutely welcome: people can identify problems, write specifications, propose designs, run tests, review generated diffs, and direct their agents.

The goal is not to replace human judgment with generated code. Contributions still need to be understandable, tested, and reviewed.

## Ways to contribute

You do not need to start by changing Go code. Useful first contributions include:

- **Try EasyEdit:** use it on your normal files and report confusing behavior or compatibility problems.
- **Report a bug:** include a minimal reproduction, expected behavior, and actual behavior.
- **Test an edge case:** check Unicode, long lines, empty files, unusual line endings, terminal resizing, or other boundary conditions.
- **Benchmark:** measure a specific operation and include the command, input size, environment, and results.
- **Propose a design:** explain the user problem, expected behavior, alternatives, and trade-offs.
- **Improve documentation:** point out missing, misleading, or difficult-to-follow instructions.
- **Contribute through an AI agent:** ask your preferred coding agent to implement a narrowly scoped issue, then inspect and test the complete diff before opening a PR.

## Before you start

1. Search the [open issues](https://github.com/0xA672/Easyedit/issues) to see whether the problem or proposal is already being discussed.
2. For a bug, try to reproduce it on the latest version and record the smallest reliable reproduction.
3. For a feature or architectural change, open an issue first. Explain the use case and proposed behavior before investing in a large implementation.
4. For a small, self-contained fix, you can open a PR directly, but linking it to an issue is helpful.

Please avoid combining unrelated changes in one PR. Small, focused contributions are easier to review and safer to merge.

## Development setup

The project is written in Go. The module declares Go 1.19, while the current GitHub Actions workflow runs its checks with Go 1.21. Use Go 1.21 or newer for local development so your environment matches CI.

Clone the repository:

```sh
git clone https://github.com/0xA672/Easyedit.git
cd Easyedit
```

Build the application:

```sh
go build -o easyedit .
```

On Windows, you can build an executable with:

```powershell
go build -o easyedit.exe .
```

## Checks to run

Run the same core checks used by the current CI workflow before submitting a code PR:

```sh
# Format Go source files
go fmt ./...

# Run static analysis
go vet ./...

# Run tests with the race detector and coverage
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# Build the CLI
go build -o easyedit .
```

If your change affects terminal behavior, also launch EasyEdit and exercise the changed interaction manually. If it affects platform-specific behavior, test on the relevant operating system when possible and state which environments you actually tested.

Do not claim that a test passed unless you ran it. If a check cannot be run in your environment, say so in the PR description.

## Repository map

The repository is organized into the following main packages and files:

| Path | Area |
| --- | --- |
| `main.go` | CLI entry point and stream-mode dispatch |
| `ui/` | Interactive terminal editor |
| `document/` | Document and text-buffer functionality |
| `command/` | Editor command handling |
| `config/` | TOML configuration loading and defaults |
| `highlight/` | Syntax highlighting |
| `clipboard/` | Clipboard integration |
| `docs/` | User and contributor documentation |

When changing behavior, look for the existing package responsible for it and add tests close to that package where practical.

## AI-agent contribution requirements

For changes to executable Go source or test code:

1. Use an AI coding agent to author the implementation.
2. Review the complete generated diff yourself. Do not submit code you cannot explain at a high level.
3. Run the relevant tests and checks; generated code is not presumed correct.
4. In the PR description, disclose the AI agent/model used and summarize what it generated. Mention meaningful human direction, review, and manual changes as well.

You are welcome to use any coding agent. The project is not tied to a particular model or provider. If a new model contributes, include its name and provider so the maintainers can consider updating the Contributors Signature Wall in the README.

For non-code contributions such as bug reports, design proposals, and feedback, no AI agent is required.

## Bug reports

A useful bug report includes:

- EasyEdit version or commit hash.
- Operating system and architecture.
- Terminal emulator, if the issue concerns the TUI.
- Exact steps to reproduce the problem.
- Expected behavior and actual behavior.
- A small sample input or file when relevant. Remove private or sensitive data before sharing it.
- Relevant error output, screenshots, or a short recording when useful.

Please report one issue per independently actionable problem where possible.

## Pull request checklist

Before opening a PR, check that:

- [ ] The PR addresses one focused issue or clearly explained improvement.
- [ ] The implementation and tests follow the AI-agent contribution rule.
- [ ] Relevant tests have been added or updated.
- [ ] `go fmt ./...`, `go vet ./...`, and `go test -v -race -coverprofile=coverage.out -covermode=atomic ./...` have been run, or any skipped check is explained.
- [ ] The application builds with `go build -o easyedit .`.
- [ ] User-facing behavior and documentation are updated where necessary.
- [ ] The PR description identifies the AI agent/model and summarizes the change.

## Pull request description template

You can use this template when opening a PR:

```markdown
## Summary

What changed?

## Motivation

What problem does this solve?

## Implementation

How does the change work? Are there trade-offs?

## Validation

- [ ] go fmt ./...
- [ ] go vet ./...
- [ ] go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
- [ ] go build -o easyedit .
- [ ] Manual testing (describe environment and steps, if applicable)

## AI contribution

- Agent:
- Model/provider:
- What the agent generated:
- Human direction/review:
```

## Community expectations

Be constructive and respectful. Questions, disagreement, and critical reviews are welcome when they focus on the problem and the proposed solution. You do not need to be an expert or a maintainer to report a problem, challenge a design, or make a useful contribution.

Thank you for helping make EasyEdit better!
