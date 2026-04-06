---
title: tmd instructions
description: 輸出嵌入的 skill 指令並附帶 vault context。
sidebar:
  order: 24
---

輸出嵌入的 skill 指令並附帶 vault context。Skill 嵌入在 `tmd` binary 中，可透過 vault 層級檔案覆寫。

## 列出可用的 skills

```bash
tmd instructions
```

列出所有嵌入的 skill 及其描述。

### 範例輸出

```
explore   Explore existing markdown files and suggest typemd type schemas...
importer  Convert existing markdown files into typemd objects...
```

### JSON 輸出

```bash
tmd instructions --json
```

回傳 `{name, description}` 物件的 JSON 陣列。

## 取得附帶 vault context 的 skill

```bash
tmd instructions explore
```

以 JSON 格式輸出 skill 指令，並附帶 vault context（type 摘要，包含屬性資訊）。

### 範例輸出

```json
{
  "name": "explore",
  "description": "...",
  "instructions": "# Explore\n\n...",
  "context": {
    "types": [
      {
        "name": "book",
        "emoji": "📚",
        "description": "Track your reading",
        "properties": [
          {"name": "author", "type": "string"}
        ]
      }
    ]
  }
}
```

若在 vault 外執行，`context` 欄位會被省略。

## 原始 skill 輸出

```bash
tmd instructions explore --skill
```

輸出原始 SKILL.md 內容（含 YAML frontmatter），適合直接儲存到 skills 目錄中。不需要 vault 即可使用。

## Vault 覆寫

在 `.typemd/instructions/<skill>.md` 放置檔案即可覆寫嵌入的 skill。覆寫檔案使用相同的 SKILL.md 格式，可包含 `name` 和 `description` frontmatter 欄位。若沒有 frontmatter，覆寫內容會取代指令本體，但保留嵌入 skill 的 metadata。
