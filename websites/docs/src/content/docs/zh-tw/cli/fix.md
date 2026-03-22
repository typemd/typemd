---
title: tmd fix
description: 修正 vault 中的常見問題。
sidebar:
  order: 10.7
---

`fix` 指令提供子指令來自動修正 vault 中的常見問題。

## tmd fix wikilinks

將簡寫的 wiki-link 目標展開為完整的 Object ID。

```bash
tmd fix wikilinks
```

### 做了什麼

遍歷 vault 中的所有物件，將簡寫的 wiki-link 目標替換為解析後的完整 ID（`type/name-ulid`）。

會被展開的簡寫格式：

| 格式 | 範例 | 解析方式 |
|------|------|---------|
| `[[type/name]]` | `[[book/clean-code]]` | 在指定 type 內依名稱解析 |
| `[[name]]` | `[[clean-code]]` | 在來源物件的同 type 內依名稱解析 |

完整 ID（`[[type/name-ulid]]`）已經是完整格式，不會被修改。

### 輸出

```
Expanded 3 wiki-link(s) to full IDs.
```

如果簡寫連結匹配到多個物件（有歧義），會回報但不展開：

```
  note/my-note-01abc: ambiguous [[golang]] — matches: book/golang-intro-01def, book/golang-guide-01ghi
```

如果所有 wiki-link 都已經是完整 ID：

```
All wiki-links are already full IDs. No changes needed.
```
