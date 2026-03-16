---
title: tmd init
description: 初始化一個新的 vault。
sidebar:
  order: 2
---

初始化一個新的 vault。建立 `.typemd/` 目錄結構、SQLite 資料庫，以及一個排除 `index.db` 的 `.typemd/.gitignore`。

建立 vault 結構後，會顯示互動式 starter type 選擇器。Starter types 是常用的 type schemas（idea、note、book），幫助你快速上手。預設全部選取，使用鍵盤自訂你的選擇：

- **↑↓** 或 **j/k** — 移動游標
- **Space** — 切換勾選
- **a** — 全選
- **n** — 全不選
- **Enter** — 確認選擇
- **Esc** — 跳過（不選任何）

選取的 starter types 會寫入 `.typemd/types/*.yaml`，使用者擁有完整的編輯權。

```bash
tmd init
```

### Flags

| Flag | 說明 |
|------|------|
| `--no-starters` | 跳過 starter type 選擇器（空白 vault） |

```bash
# 非互動模式：跳過 starter types
tmd init --no-starters
```

### Vault 設定

選取 starter types 時，`tmd init` 也會建立 `.typemd/config.yaml`，設定快速建立 object 的預設 type：

- 選取 **idea** → `cli.default_type: idea`
- 選取 **note**（沒有 idea）→ `cli.default_type: note`
- 只選取 **book** → 不建立 config 檔

這讓你可以直接執行 `tmd object create "Some Thought"` 而不需指定 type。詳見 [tmd object create](/zh-tw/cli/create)。

在已初始化的 vault 上執行 `tmd init` 會回傳錯誤。
