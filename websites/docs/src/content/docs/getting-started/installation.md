---
title: Installation
description: How to install TypeMD on your machine.
sidebar:
  order: 2
---

```bash
brew install typemd/tap/typemd-cli
```

Or from source (requires [Go](https://go.dev/)):

```bash
go install github.com/typemd/typemd/cmd/tmd@latest
```

Pre-built binaries are also available on [GitHub Releases](https://github.com/typemd/typemd/releases).

## Verify Installation

```bash
tmd --version
```
