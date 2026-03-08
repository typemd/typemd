---
title: 安裝
description: 如何在你的機器上安裝 TypeMD。
sidebar:
  order: 2
---

## Homebrew（macOS）

```bash
brew install typemd/tap/typemd-cli
```

## 從原始碼安裝

確保你已安裝 [Go](https://go.dev/)，然後執行：

```bash
go install github.com/typemd/typemd/cmd/tmd@latest
```

## 驗證安裝

```bash
tmd --version
```
