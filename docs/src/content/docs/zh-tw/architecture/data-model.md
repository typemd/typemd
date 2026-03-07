---
title: 資料模型
description: TypeMD 如何儲存和索引資料。
sidebar:
  order: 1
---

## 儲存

Object 以 Markdown 檔案搭配 YAML frontmatter 儲存在 `objects/<type>/` 底下。完整的 Object ID 格式為 `type/<slug>-<ulid>`，例如 `book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz`。透過 CLI 建立的 Object，檔名會自動附加 ULID（26 位小寫字元）以確保唯一性。

```
vault/
├── .typemd/
│   ├── types/              # Type schema 定義（YAML）
│   │   ├── book.yaml
│   │   └── person.yaml
│   └── index.db            # SQLite 索引（自動更新）
└── objects/
    ├── book/
    │   └── golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz.md
    └── person/
        └── alan-donovan-01jqr3k5mpbvn8e0f2g7h9txyz.md
```

## 索引

TypeMD 使用 SQLite 搭配 FTS5 進行索引。索引儲存在 `.typemd/index.db`，包含：

- Object 中繼資料（Type、檔名、屬性）
- 涵蓋檔名、屬性和內文的全文搜尋索引

開啟 vault 時，若資料庫為空或缺失（例如 clone 後首次開啟），索引會自動同步。使用 TUI 或 CLI 指令時也會保持更新。在 TypeMD 外部手動編輯檔案後，使用 `tmd reindex` 重建。

## 架構

TypeMD 是一個 monorepo，共用 Go 核心程式庫並提供多種介面：

```
typemd/
├── core/       # 核心程式庫——Object、Type、Relation、索引
├── cmd/        # CLI 指令（Cobra）
├── tui/        # 終端 UI（Bubble Tea）
├── mcp/        # MCP server，用於 AI 整合
├── web/        # Web UI API（規劃中）
├── app/        # 桌面應用程式（規劃中）
├── site/       # 官方網站（Astro）→ typemd.io
└── docs/       # 文件（Starlight）→ docs.typemd.io
```

所有介面共用相同的 `core` 程式庫。

## 技術堆疊

- **語言**：Go
- **TUI**：[Bubble Tea](https://github.com/charmbracelet/bubbletea) + [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **MCP**：[mcp-go](https://github.com/mark3labs/mcp-go) — Model Context Protocol server
- **索引**：SQLite 搭配 FTS5 全文搜尋
- **儲存**：Markdown + YAML frontmatter
