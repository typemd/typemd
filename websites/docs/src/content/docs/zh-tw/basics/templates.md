---
title: 模板
description: 模板如何為新 Object 提供預設內容和自動產生名稱。
sidebar:
  order: 3
---

模板（Template）讓你控制新 Object 的初始內容。TypeMD 支援兩種模板：**object template** 提供預設的 frontmatter 和內文，以及 **name template** 從樣式自動產生 Object 名稱。

## Object Template

每個 Type 可以在 `templates/<type>/<name>.md` 存放一個或多個 object template。Template 是一般的 Markdown 檔案，包含可選的 frontmatter 和內文，作為建立新 Object 時的預設內容。

```
vault/
├── templates/
│   └── book/
│       ├── review.md       # book 的「review」template
│       └── summary.md      # book 的「summary」template
└── objects/
    └── book/
        └── ...
```

### Template 內容

Template 檔案可以包含 frontmatter 屬性和內文：

```markdown
---
status: to-read
tags:
  - tag/to-review
---

## Summary

## Key Takeaways

## Notes
```

Template 的 frontmatter 屬性會覆蓋 schema 的預設值。Template 的內文會成為新 Object 的初始內文。

### Template 解析

使用 `tmd object create` 建立 Object 時，template 會自動解析：

- **沒有 template** — 不套用 template
- **單一 template** — 自動套用，不需要提示
- **多個 template** — 互動式選擇提示（或使用 `-t` 跳過提示）

```bash
# 如果 book 只有一個 template，會自動套用
tmd object create book clean-code

# 明確指定 template
tmd object create book clean-code -t review

# 如果有多個 template，會提示你選擇
tmd object create book clean-code
```

### 系統屬性處理

Template 可以覆蓋**使用者撰寫**的系統屬性（`name`、`description`、`tags`）。**自動管理**的系統屬性（`created_at`、`updated_at`）在 template 中會被忽略——它們永遠反映實際的建立時間。Template 中不屬於 type schema 定義的屬性會靜默忽略。

## Name Template

Object 名稱通常在建立時手動提供。但對於有可預測命名模式的 type（每日日記、會議筆記、週報），手動輸入既重複又容易出錯。Name template 讓 type schema 使用佔位符定義自動產生的名稱，實現一個指令就能建立 object。

### 定義 name template

在 type schema 的 `properties` 陣列中加入一個 `name` 項目，搭配 `template` 欄位：

```yaml
# types/journal/schema.yaml
name: journal
emoji: 📓
properties:
  - name: name
    template: "日記 {{ date:YYYY-MM-DD }}"
  - name: mood
    type: select
    options:
      - value: good
      - value: okay
      - value: bad
```

`properties` 中的 `name` 項目只允許 `template` 欄位——不能設定 `type`、`options`、`pin`、`emoji` 或其他屬性欄位。

### 日期佔位符

`{{ date:FORMAT }}` 佔位符接受以下格式代碼：

| 代碼 | 說明 | 範例 |
|------|------|------|
| `YYYY` | 四位數年份 | `2026` |
| `MM` | 兩位數月份 | `03` |
| `DD` | 兩位數日期 | `14` |
| `HH` | 兩位數小時（24 小時制） | `09` |
| `mm` | 兩位數分鐘 | `30` |
| `ss` | 兩位數秒數 | `05` |

範例：

| Template | 結果（2026-03-14 09:30 時） |
|----------|---------------------------|
| `日記 {{ date:YYYY-MM-DD }}` | `日記 2026-03-14` |
| `Meeting {{ date:YYYY-MM-DD HH:mm }}` | `Meeting 2026-03-14 09:30` |
| `Week {{ date:YYYY-MM }}` | `Week 2026-03` |
| `Weekly Review` | `Weekly Review`（無佔位符，作為字面值使用） |

### CLI 用法

當 type 有 name template 時，name 參數為選填：

```bash
# 從 template 自動產生名稱
tmd object create journal

# 明確指定名稱會覆蓋 template
tmd object create journal "my-journal"
```

如果 type 沒有 name template，則必須提供 name 參數。

## 相關頁面

- [屬性](/zh-tw/basics/properties#系統屬性) — 系統屬性的可變性
- [tmd object create](/zh-tw/tui/create) — 建立 object 的 CLI 參考
- [Type](/zh-tw/concepts/types) — type schema 結構
