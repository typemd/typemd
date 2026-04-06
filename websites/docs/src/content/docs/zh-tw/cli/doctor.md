---
title: tmd doctor
description: 全面的 vault 健康檢查，提供診斷與自動修復。
sidebar:
  order: 19
---

對 vault 進行全面的健康檢查，涵蓋 8 個類別。是 [`tmd type validate`](/zh-tw/cli/validate/) 的超集，額外包含結構完整性檢查。

```bash
tmd doctor
```

## 檢查類別

### 1. Schemas

驗證所有 type schema（同 `tmd type validate` 階段 1）。

### 2. Objects

依據 type schema 驗證 object 屬性（同階段 2）。

### 3. Relations

檢查所有 relation 端點是否參照存在的 object（同階段 3）。

### 4. Wiki-links

偵測壞掉的 `[[target]]` 參照（同階段 4）。

### 5. Uniqueness

檢查設定 `unique: true` 的 type 是否有重複的 name（同階段 5）。

### 6. Files

偵測損壞的 object 檔案——`objects/` 下的檔案若有無法解析的 YAML frontmatter 或缺少 `---` 分隔符號，在一般操作中會被靜默跳過，doctor 會主動報告。

### 7. Index

檢查 SQLite index 是否與磁碟上的檔案同步。**自動修復**：若 index 過時，會自動重建。

### 8. Orphans

偵測 `objects/` 或 `templates/` 下沒有對應 type schema 的目錄。以警告方式報告。

## 輸出

結果按類別分組，以 `✓`/`✗` 標示：

```
  ✓ Schemas
  ✓ Objects
  ✓ Relations
  ✓ Wiki-links
  ✓ Uniqueness
  ✗ Files
    [error] book/bad-file.md: yaml: did not find expected key
  ✓ Index (auto-fixed)
  ✗ Orphans
    [warn] object directory objects/ghost has no type schema

2 issue(s) found, 1 auto-fixed.
```

## 結束代碼

- `0`——沒有發現問題（自動修復的項目不算在內）
- `1`——發現一個或多個問題
