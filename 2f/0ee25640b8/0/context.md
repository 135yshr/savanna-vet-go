# Session Context

## User Prompts

### Prompt 1

Implement the following plan:

# Plan: README.md の作成

## Context
savanna プロジェクトに README.md が存在しないため、Go OSS の慣習に沿った英語の README.md を新規作成する。元になった Java/Maven プロジェクト (kawasima/savanna-maven-plugin) へのクレジットも適切に記載する。

## 方針
- 言語: 英語
- インストール: `go install` を記載
- ライセンス: MIT License
- バッジ: Go Version, Go Report Card, License を冒頭に表示

## README.md の構成

1. **タイトル + バッジ**
   - プロジェクト名 `savanna`
   - バッジ: Go Version, Go Report Card, MIT License
   - 1行説明: A Go CLI tool that detects test smells in Go test fil...

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

