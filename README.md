# Telnyx Terraform Provider

This is a Terraform provider for Telnyx. It is currently in alpha status and not recommended
for production use.

## Documentation

The provider [documentation can be found on the registry](https://registry.terraform.io/providers/petsinc/telnyx/latest/docs).

## Developing

### Initial Setup

#### Direct System Dependencies

First, you need a couple global dependencies installed, see their documentation for details:

- [direnv](https://direnv.net/docs/installation.html)
- [asdf](https://asdf-vm.com/guide/getting-started.html)
  - Be sure to use the official `git` installation method or you may have issues with
    pre-commit hooks finding `~/.asdf/bin`

Note that these tools require a UNIX-style shell, such as bash or zsh. If
you are on Windows, you can use WSL or Git Bash. If you are using Pycharm,
you can configure the built-in terminal to use Git Bash.

#### First Steps

Clone the repo and run `direnv allow`. This will take a while on the first time to install the remaining dependencies.

#### ENV Variables

You will need to set `TELNYX_API_KEY` in the `.env` file to run tests.

## Day-to-day Development

### Tests

To run all tests, run `just test`. If you add a new resource, be sure to add it
to the tests in `provider_test.go`.

### Conventional Commits & Semantic Release

This project uses [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/)
to power [semantic release](https://semantic-release.gitbook.io/semantic-release/). This means
that when you commit, you should use the following format:

```
<type>[optional scope]: <description>
```

For example, `feat: Add new feature` or `fix: Fix bug`.

When creating a PR, please name the PR in this way as well so that the squashed
commit from the PR will have a conventional commit message. There is a PR
check that enforces this.

### Pre-commit Hooks

This project uses Husky and Lint-staged to run pre-commit hooks. This means that
when you commit, it will format the files
you edited, and also check that your commit message is a conventional commit.

If you are not able to commit, it is likely because your commit message is not
in the conventional commit format.

### `justfile`

This project uses [`just`](https://github.com/casey/just) to automate various
project activities. Any new project commands should be added to the `justfile`.

There are commands such as `just test` and `just format-go`, see
the `justfile` for more details.

## Author

This project is authored by Patient Engagement Technologies.
