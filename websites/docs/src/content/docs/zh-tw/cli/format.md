---
title: tmd format
description: 正規化 object frontmatter 和 schema YAML 格式。
sidebar:
  order: 10.5
---

格式化所有 object 和 schema 檔案，統一 property 排序和 YAML 風格。類似 `gofmt`，強制使用一致的檔案格式。

```bash
tmd format                    # 格式化所有 objects 和 schemas
tmd format --type book        # 只格式化 book 的 objects 和 schema
tmd format --dry-run          # 列出需要格式化的檔案（CI 模式）
```

## 旗標

| 旗標 | 說明 |
|------|------|
| `--type <name>` | 只格式化指定 type 的 objects 和 schema |
| `--dry-run` | 列出需要格式化的檔案，不實際修改 |

## 運作方式

### Property 排序

Frontmatter 的 property 會按照標準順序重寫：

1. **系統屬性** — `name`、`description`、`created_at`、`updated_at`、`tags`
2. **Schema 定義的屬性** — 依 type schema 中定義的順序
3. **額外屬性** — 不在 schema 中的屬性，按字母順序排列

### YAML 正規化

YAML 格式會統一為 `yaml.v3` 預設值——一致的引號風格、縮排和值的表示方式。

### Schema 格式化

Type schema 檔案（`.typemd/types/<name>/schema.yaml`）也會透過標準序列化器重新格式化。

### 保留不變的內容

- **Body 內容**不會被修改
- **`updated_at`** 不會被更改——格式化是純粹的排版變更，不改變語意

## Dry-run 模式

使用 `--dry-run` 時，命令會列出需要變更的檔案，並在有檔案需要格式化時以結束代碼 1 退出。適用於 CI 流程：

```bash
tmd format --dry-run || echo "有檔案需要格式化"
```

### 結束代碼

- `0`——所有檔案已經格式化完成
- `1`——有一個或多個檔案需要格式化（dry-run），或發生錯誤
