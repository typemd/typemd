---
title: tmd object create
description: 從 Type schema 建立新 Object。
sidebar:
  order: 2.5
---

根據 Type schema 建立新的 Object 檔案（Markdown + YAML frontmatter）。

```bash
tmd object create [type] [name]
tmd object create book "Clean Code"
tmd object create book "Clean Code" -t review
tmd object create --type idea "Some Thought"
tmd object create "Some Thought"              # 使用 config 中的預設 type
tmd object create book
```

`type` 參數可以省略，前提是設定了 `--type` flag 或在 `.typemd/config.yaml` 中設定了 `cli.default_type`。當只提供一個參數且該參數是已知的 type 名稱時，會被視為 type。若名稱與 type 名稱相同，請使用 `--type` flag 明確指定。

名稱會自動轉換為 slug 作為檔名（例如 "Clean Code" 會變成 `clean-code`），而原始輸入會保留在 frontmatter 的 `name` 屬性中。

如果 type schema 定義了 name template，可以省略 `name` 參數——名稱會從 template 自動產生。若未提供名稱也沒有 name template，會回傳錯誤。

指令會在 `objects/<type>/` 下產生 `.md` 檔案，所有 schema 定義的屬性設為預設值（若無預設值則為 `null`）。Object 同時會被加入 SQLite 索引。

每個 Object 會被分配唯一的 ULID 後綴，因此同一 Type 下可以有相同名稱的物件（例如兩個 `book/clean-code-<ulid>` 但 ULID 不同）。如果 Type 不存在，會回傳錯誤。

## Flags

| Flag | 說明 |
|------|------|
| `--type` | Object type（覆蓋 config 預設值，無短旗標） |
| `-t`, `--template` | 要使用的 template 名稱（來自 `templates/<type>/`） |

## 預設 Type

如果在 `.typemd/config.yaml` 中設定了 `cli.default_type`，可以省略 type 參數：

```yaml
# .typemd/config.yaml
cli:
  default_type: page
```

```bash
# 當 default_type 是 "page" 時，以下兩種寫法等效：
tmd object create page "Some Thought"
tmd object create "Some Thought"
```

`--type` flag 永遠優先於 config 預設值。

## Template

如果該 type 在 `templates/<type>/` 下有定義 template，會自動解析：

- **沒有 template** — 不套用任何 template
- **單一 template** — 自動套用，不需提示
- **多個 template** — 互動式選擇提示（或使用 `-t` 跳過提示）

Template 的 frontmatter 屬性會覆蓋 schema 預設值，template 的內文會成為 object 的初始內文。自動管理的系統屬性（`created_at`、`updated_at`）在 template 中會被忽略。
