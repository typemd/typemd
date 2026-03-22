---
title: AI 設定
description: 設定 Claude Code 與 typemd marketplace plugin，啟用 AI 功能。
sidebar:
  order: 4
---

typemd 透過兩個管道整合 AI：

- **TUI AI 功能** — 自動產生描述、自動建議標籤、schema 探索，由 `claude` CLI 驅動
- **Skill instructions** — `tmd instructions` 輸出嵌入的 skill 指令並附帶 vault context，可供任何 AI 工具使用

## 前置條件

- 已安裝 [typemd](/zh-tw/getting-started/installation) 並初始化 vault
- 已安裝 [Claude Code](https://claude.com/code)

## 1. 安裝 typemd plugin

**typemd** plugin 讓 Claude Code 學會如何操作你的 vault。加入 typemd marketplace 並安裝：

```bash
claude
# 在 Claude Code 中：
/plugin marketplace add typemd/marketplace
/plugin install typemd@typemd-marketplace
```

這會安裝兩個 skill：

| Skill | 功能 |
|-------|------|
| **vault-guide** | 教 AI 如何設計與管理 vault — CLI 指令、type schema、物件格式、wiki-links、views、templates |
| **instructions-guide** | 教 AI 如何使用 `tmd instructions` 來取得附帶 vault context 的 skill 指令 |

vault-guide 會在 Claude Code 偵測到 `.typemd/` 目錄時自動啟用。

## 2. 啟用 TUI AI 功能

在 vault 設定中加入：

```yaml
# .typemd/config.yaml
ai:
  enabled: true
```

或使用 CLI：

```bash
tmd config set ai.enabled true
```

啟用後，TUI 中可以使用以下 AI 操作：

| 按鍵 | 動作 |
|------|------|
| `g` | 開啟 AI 動作選單（在物件上）— 產生描述 / 建議標籤 |
| `Ctrl+E` | Schema 探索模式（在 type 上）— AI 分析物件並建議 schema 改進 |

這些功能需要 `claude` CLI 已安裝且在 PATH 中可存取。

## 3. 使用 skill instructions

`tmd instructions` 輸出嵌入的 skill 指令，並附帶你的 vault type 定義。這讓任何 AI 工具都能理解你的 vault，不需要手動收集 context。

```bash
# 列出可用的 skills
tmd instructions

# 取得 explore skill 及 vault context（JSON 格式）
tmd instructions explore

# 儲存 skill 到 Claude Code skills 目錄
tmd instructions explore --skill > .claude/skills/explore/SKILL.md
```

JSON 輸出包含你的 vault type 摘要（名稱、emoji、描述、屬性），AI 可以立即了解你的資料模型。完整說明請參考 [`tmd instructions`](/zh-tw/cli/instructions)。

## 接下來

安裝 typemd plugin 並啟用 AI 後，你可以：

- 請 Claude Code **分析你的 markdown 檔案**並建議 type schema（explore skill）
- 請 Claude Code **匯入既有檔案**到你的 vault（importer skill）
- 使用 **TUI AI 功能**自動產生物件的描述和標籤
- 將 `tmd instructions` 的輸出導入**任何 AI 工具**，進行 vault 感知的工作流程
