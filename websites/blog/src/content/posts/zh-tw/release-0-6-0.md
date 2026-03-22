---
title: "TypeMD 0.6.0 — AI 走進你的 Vault"
description: "v0.6.0 帶來 AI 功能：自動產生描述、建議標籤、schema 探索，以及 tmd instructions 讓任何 AI 工具理解你的 vault。"
date: 2026-03-22
tags: [發布, 開發日誌]
---

0.5.0 打好了 View 系統的基礎，讓你用自己的方式看資料。這次我們往更大膽的方向走——**讓 AI 直接理解你的 vault**。

## AI 來了：描述、標籤、Schema 探索

TUI 新增三個 AI 動作，全都由 `claude` CLI 驅動：

- 在物件上按 `g`，選「Generate description」— AI 讀你的物件內容，自動產生一段描述
- 在物件上按 `g`，選「Suggest tags」— AI 根據內容建議適合的標籤
- 在型別上按 `Ctrl+E` — AI 抽樣分析該型別的物件，建議 schema 改進（缺少的屬性、更好的型別選擇）

這些功能需要啟用 AI：

```bash
tmd config set ai.enabled true
```

## tmd instructions：讓任何 AI 工具理解你的 vault

除了 TUI 裡的 AI 功能，我們更進一步——讓你的 vault 知識可以被**任何** AI 工具取用。

```bash
tmd instructions explore
```

這會輸出 explore skill 的完整指令，加上你的 vault 的型別摘要（名稱、emoji、描述、屬性）。把這段 JSON 餵給任何 AI，它就能立刻理解你的資料模型。

兩個嵌入的 skill：
- **explore** — 分析 markdown 檔案，建議型別 schema
- **importer** — 把現有 markdown 檔案轉成 vault 物件

想自訂 skill？放一個 `.typemd/instructions/<skill>.md` 就能覆寫。

## Marketplace：typemd plugin

我們也推出了 [typemd marketplace](https://github.com/typemd/typemd/tree/main/marketplace) 的第一個 plugin——**typemd**。安裝後 Claude Code 會自動學會如何操作你的 vault：

```bash
/plugin marketplace add typemd/marketplace
/plugin install typemd@typemd-marketplace
```

完整設定流程請看 [AI Setup 指南](https://docs.typemd.io/getting-started/ai-setup)。

## 更聰明的導航

這版也大幅提升了 CLI 的導航體驗：

- **Shell 自動補全** — 物件 ID、型別名稱、關聯名稱都能 tab 補全
- **互動式消歧義** — 前綴匹配到多個物件？彈出模糊選取器讓你選
- **簡寫 wiki-links** — 可以寫 `[[Clean Code]]` 或 `[[book/clean-code]]`，同步時自動解析為完整 ID

## 還有什麼

- **Stats 儀表板** — TUI 中的統計模式，看 vault 全貌
- **Toast 通知** — 同步警告和 AI 錯誤以浮動訊息呈現
- **`tmd format`** — 一個指令正規化所有物件的 frontmatter
- **物件重新命名** — 在 TUI 中按 `r` 直接改名

## 下一步

AI 的基礎打好了。接下來 0.7.0 我們要讓 TUI 的編輯體驗更完整——inline 屬性編輯、關聯選取器，讓你在 TUI 裡就能完成所有操作。

升級方式：

```bash
brew upgrade typemd/tap/typemd-cli
```

或從 [GitHub Releases](https://github.com/typemd/typemd/releases/tag/v0.6.0) 下載預編譯的執行檔。
