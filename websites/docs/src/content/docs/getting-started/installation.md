---
title: Installation
description: How to install TypeMD on your machine.
sidebar:
  order: 2
---

## Homebrew (macOS)

```bash
brew install typemd/tap/typemd-cli
```

## From Source

Make sure you have [Go](https://go.dev/) installed, then:

```bash
go install github.com/typemd/typemd/cmd/tmd@latest
```

## Verify Installation

```bash
tmd --version
```
