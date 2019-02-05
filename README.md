gitprompt
============

Parses `git status --porcelain=v2 --branch` and outputs a nicely formatted string, similar to vcprompt. It can also output line-delimited status info that can be parsed by your shell prompt.

The minimum git version for porcelain v2 with `--branch` is `v2.13.2`.

## Output explained:

 ` <branch>@<commit> [↑/↓ <ahead/behind count>][untracked][unmerged][modified][dirty/clean]`

- `?`  : untracked files
- `‼`  : unmerged : merge in process
- `Δ`  : modified : unstaged changes

Definitions taken from: https://www.kernel.org/pub/software/scm/git/docs/gitglossary.html#def_dirty
- `✘`  : dirty : working tree contains uncommited but staged changes
- `✔`  : clean : working tree corresponds to the revision referenced by HEAD

## Usage

Run `gitprompt` without any options to get the colorful output :)
For all supported options see `gitprompt -h`.
