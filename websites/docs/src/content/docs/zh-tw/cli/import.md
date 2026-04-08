---
title: tmd import
description: 將外部 markdown 檔案匯入 vault。
sidebar:
  order: 25
---

使用三步驟工作流程將外部 markdown 檔案匯入 vault：掃描、規劃、執行。

## 掃描來源

```bash
tmd import scan <paths...>
```

掃描來源目錄或檔案中的 markdown 內容，擷取 frontmatter 模式並收集檔案統計資料。

### 範例

```bash
tmd import scan ~/notes ~/docs/blog
```

輸出 JSON `ScanResult`，包含：

- `sources` — 每個 markdown 檔案的路徑、大小和 frontmatter 鍵值
- `file_count` — 找到的 markdown 檔案總數
- `directories` — 目錄結構與每個目錄的檔案數量
- `patterns` — frontmatter 鍵值的彙總統計（出現頻率和範例值）
- `existing_types` — vault 中已定義的 type schema 及其屬性
- `no_frontmatter_count` — 沒有 YAML frontmatter 的檔案數量

## 產生規劃

```bash
tmd import plan <classifications-file>
tmd import plan classifications.json --output plan.json
```

從包含物件分類的 JSON 檔案產生匯入規劃。分類檔案是一個 `ObjectPlan` 的 JSON 陣列，通常由 AI 分析掃描結果後產生。

### Flags

| Flag | 縮寫 | 說明 |
|------|------|------|
| `--output` | `-o` | 將規劃寫入檔案而非標準輸出 |

### ObjectPlan 格式

```json
[
  {
    "source_path": "books/clean-code.md",
    "type_name": "book",
    "name": "Clean Code",
    "properties": {"author": "Robert C. Martin"},
    "body": "書籍內容...",
    "conflict": "none",
    "depends_on": []
  }
]
```

plan 指令會偵測需要建立的新 type、檢查與現有物件的衝突，並計算依賴順序的匯入序列（tag 優先，接著是獨立物件，最後是有依賴的物件）。

## 執行規劃

```bash
tmd import execute <plan-file>
```

執行匯入規劃，在 vault 中建立 type 和物件。

### 執行階段

1. **建立 type** — 規劃中列出的新 type schema
2. **建立物件** — 按依賴順序，遵循衝突標記（`skip` 或 `none`）
3. **調和** — 解析所有匯入物件中的 wiki-links

### 輸出

回傳 JSON `ImportReport`：

```json
{
  "types_created": 1,
  "objects_created": 8,
  "objects_skipped": 2,
  "objects_failed": 0,
  "details": [...],
  "unresolved_refs": [],
  "suggestions": []
}
```

## AI 輔助工作流程

import 指令設計用於透過 `onboarding` skill 進行 AI 編排：

1. AI 執行 `tmd import scan` 分析來源檔案
2. AI 將檔案分類為 type 並對應屬性
3. AI 建構分類 JSON 並產生規劃
4. 使用者審核並核准規劃
5. AI 執行 `tmd import execute` 使用核准的規劃

使用 `tmd instructions onboarding` 取得完整的 onboarding skill 指令及 vault context。
