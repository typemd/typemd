---
title: 翻譯指南
description: 如何為 TypeMD 文件貢獻翻譯。
sidebar:
  order: 5
---

TypeMD 文件使用 [Starlight 的 i18n 系統](https://starlight.astro.build/guides/i18n/) 支援多語言。

## 支援的語言

| 語言 | 代碼 | 目錄 |
|------|------|------|
| English（預設） | `en` | `docs/src/content/docs/` |
| 繁體中文 | `zh-TW` | `docs/src/content/docs/zh-tw/` |

## 新增翻譯

1. 在 `docs/src/content/docs/` 底下找到英文原始檔案
2. 在 `docs/src/content/docs/zh-tw/` 底下建立相同路徑的對應檔案
3. 翻譯內容，保留程式碼區塊、指令和技術術語為英文
4. 更新 YAML frontmatter 中的 `title` 和 `description`

## 新增語言

1. 在 `docs/astro.config.mjs` 的 `locales` 中新增語言設定
2. 在每個側邊欄群組的 `translations` 欄位中新增翻譯
3. 在 `docs/src/content/docs/<locale>/` 底下建立目錄
4. 以翻譯內容鏡像英文的檔案結構
5. 在專案根目錄建立 `README.<LOCALE>.md`

## 翻譯準則

- **保留英文：** TypeMD、Object、Type、Relation、CLI 指令、屬性名稱、程式碼區塊、YAML 鍵值、檔案路徑
- **翻譯：** 描述性文字、章節標題、表格說明、frontmatter 中的 `title` 和 `description`
- **排版：** 全型中文字和半型英數字中間要留空格
