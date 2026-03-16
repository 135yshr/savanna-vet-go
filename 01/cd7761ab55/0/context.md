# Session Context

## User Prompts

### Prompt 1

Implement the following plan:

# リリース自動化の追加

## Context

savanna-vet-go プロジェクトにリリース自動化機能がない。meow プロジェクト (`/Users/135yshr/go/src/github.com/135yshr/meow`) の構成を参考に、semantic-release + gitmoji + GoReleaser によるリリースパイプラインを構築する。

main ブランチへのマージ時に、gitmoji コミットメッセージからバージョンを自動決定し、GitHub Release とクロスコンパイルバイナリを自動生成する。

## 作成するファイル

### 1. `.releaserc.json` — semantic-release 設定

meow の `.releaserc.json` をベースに作成。

- branches: `["main"]`
- plugins: `semantic-release-gitmoji`, `@semantic-release/changelog`,...

### Prompt 2

# Smart Commit with Gitmoji

Execute the following steps non-interactively:

## Branch Management

- If currently on `main` or `master` branch, create and checkout a new feature branch with a descriptive name
- If on any other branch, proceed with commit on current branch (no branch creation)
- If no changes are detected, exit without doing anything

## Change Analysis & Commit

1. **Review all changes** using `git status` and `git diff --staged` (or `git diff` if nothing staged)
2. **Stage c...

