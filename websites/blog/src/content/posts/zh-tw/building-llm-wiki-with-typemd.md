---
title: "用 typemd 打造你的 LLM Wiki"
description: "LLM Wiki 是一種讓 AI 持續維護個人知識庫的新模式。typemd 從第一天就是為這件事設計的。"
date: 2026-04-07
tags: [理念, 架構]
---

最近 Tobi Lütke 分享了一個叫 [LLM Wiki](https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f) 的概念。核心想法是：不要讓 AI 每次回答問題時從頭搜尋文件，而是讓 AI **持續維護一座結構化的知識庫**。每加入一份新資料，AI 就更新摘要、補充交叉引用、標記矛盾。知識是累積的，不是每次重新推導的。

讀完之後很有共鳴——這跟我們做 typemd 的初衷高度吻合。LLM Wiki 把我們一直在想的事情講得更清楚了，也讓我們看到了一些之前沒想到的面向。這篇文章來拆解這個概念，看看 typemd 目前走到哪裡，以及接下來要往哪裡去。

## RAG 的瓶頸

大部分人使用 AI 處理文件的方式是 RAG（Retrieval-Augmented Generation）：上傳一堆檔案，AI 在回答時搜尋相關片段，拼湊出答案。NotebookLM、ChatGPT 的檔案上傳，大多是這個模式。

問題在於：**每次都從零開始**。問一個需要綜合五份文件的問題，AI 每次都得重新找、重新拼。沒有任何東西被建立起來。

LLM Wiki 提出了不同的思路：與其讓 AI 每次即時搜尋，不如讓 AI **事先整理好**。把知識編譯一次，然後持續維護。交叉引用已經建好，矛盾已經被標記，綜合分析已經反映了所有讀過的東西。

## LLM Wiki 的三層架構

LLM Wiki 定義了三個層次：

1. **原始來源**（Raw Sources）——蒐集的文章、論文、筆記。不可變的事實來源。
2. **Wiki**——AI 產生並維護的結構化頁面。摘要、實體頁面、概念頁面、比較分析。
3. **Schema**——告訴 AI 如何組織 wiki 的約定：命名規則、頁面格式、工作流程。

三種主要操作：**Ingest**（加入新來源，AI 讀取後更新多個 wiki 頁面）、**Query**（提問，好的回答歸檔回 wiki）、**Lint**（定期健康檢查，找矛盾和孤立頁面）。

人負責策展、提問、思考意義。AI 負責苦工——摘要、交叉引用、歸檔、記帳。

## typemd 和 LLM Wiki 的交集

typemd 的核心設計跟 LLM Wiki 的三層架構有天然的對應：

| LLM Wiki | typemd | 說明 |
|---|---|---|
| Wiki Pages | Objects | 帶 YAML frontmatter 的結構化 Markdown |
| Schema | Type Schemas | `types/<name>/schema.yaml` 定義屬性、關聯、驗證規則 |
| 交叉引用 | Wiki-links | `[[type/name]]` 語法，自動展開為完整 ID 並建立 backlink 索引 |
| Typed Relations | Relation properties | 雙向、多值、型別約束的關聯系統 |
| Index | SQLite FTS5 | 全文搜尋 + 結構化屬性查詢 |
| Dataview 式查詢 | Views | per-type 的 filter/sort/group 規則，存為 YAML |
| Graph View | `tmd graph` | 輸出 DOT 格式的知識圖譜 |
| git 版本追蹤 | 原生支援 | vault 就是一個 git repo of markdown files |

AI 方面，typemd 有三個核心操作：**Describe**（根據內容自動產生摘要）、**SuggestTags**（分析內容推薦標籤）、**ExploreSchema**（分析物件建議 schema 調整）。這些透過可配置的 AI provider（Claude CLI 或 OpenAI-compatible API）執行。

Obsidian 使用者會覺得語法很熟悉——typemd 的 frontmatter 格式跟 Obsidian 相容，Web Clipper 擷取的 Markdown 可以直接被讀取。

## CLI-First：不是改造，是原生

LLM Wiki 的實踐需要 AI 能**直接操作**知識庫——讀取、搜尋、建立、更新、連結。這就是為什麼 Obsidian 在今年二月發布了[官方 CLI](https://obsidian.md/cli)，Logseq 也推出了[內建 MCP server 的 CLI](https://www.npmjs.com/package/@logseq/cli)。這些工具都在為同一件事做準備：讓 AI agent 能程式化地操作知識庫。

差別在於，Obsidian 和 Logseq 是 GUI-first 的工具，CLI 和 MCP 是後來加上的橋接層。typemd 從第一天就是 CLI tool——`tmd` 就是主要介面。MCP server（`tmd mcp`）不是外掛，是內建的。每個操作都有對應的 CLI 命令，每個輸出都支援 JSON 格式。

這不是說 GUI 不重要——typemd 也有 TUI 和 Web UI。但核心操作全部都能在終端機裡完成，意味著 AI agent 操作 typemd vault 跟人操作的方式是一樣的，沒有額外的翻譯層。

## Pattern vs Platform

不過 typemd 和 LLM Wiki 的定位其實不一樣。

LLM Wiki 是一個 **pattern**——一套工作原則。Wiki 就是一堆 Markdown 檔案，schema 就是一份 CLAUDE.md，結構完全靠 LLM 和人的約定來維持。這很靈活，但也意味著 LLM 每次都要重新理解檔案之間的關係。

typemd 是一個 **platform**——有 type system、relation engine、reconciler、index。結構不是靠約定，而是靠 schema 約束和自動化維護。LLM 不需要猜一個物件有哪些屬性——schema 告訴它。不需要手動追蹤引用——reconciler 自動展開 wiki-link 並建立 backlink 索引。不需要每次搜尋全部檔案——SQLite FTS5 已經做好索引。

另一個差異是**誰維護 wiki**。LLM Wiki 的理念偏向「LLM 寫並維護全部」。typemd 的設計是**人和 LLM 共同維護**——你可以在 TUI 裡直接編輯屬性和內容，也可以讓 AI 透過 MCP 操作。`locked` 屬性讓你標記不希望被自動修改的物件。LLM 處理苦工，但人保留最終控制權。

## 我們還缺什麼

受到 LLM Wiki 的啟發，我們盤點了要成為完整 LLM Wiki 平台還缺少的東西，也開了對應的 issue：

**MCP 讀寫工具** — 目前 MCP server 只有 `search` 和 `get_object` 兩個唯讀工具。LLM 能讀不能寫。我們正在擴充完整的讀寫能力（[#377](https://github.com/typemd/typemd/issues/377)、[#382](https://github.com/typemd/typemd/issues/382)），讓任何支援 MCP 的 AI client 都能直接維護你的 vault。

**Onboarding 流程** — 最大的摩擦點是「怎麼把一堆現有的 Markdown 變成結構化的 vault」。我們在設計一個完整的流程（[#381](https://github.com/typemd/typemd/issues/381)）：掃描來源 → LLM 分類建議 schema → 批次匯入 → 驗證。支援多來源和追加匯入。

**Ingest 和 Lint** — LLM Wiki 的兩個核心操作。Ingest（[#378](https://github.com/typemd/typemd/issues/378)）讓 LLM 讀取新來源後自動更新相關頁面；Lint（[#379](https://github.com/typemd/typemd/issues/379)）讓 LLM 定期掃描孤立頁面、矛盾描述、缺少的概念頁面。目前 typemd 的驗證只做結構層面的檢查（schema 符合、引用存在），語意層面的健康檢查還是空白。

**衝突偵測** — 當人和 LLM 同時維護一個 vault，需要 optimistic locking（[#384](https://github.com/typemd/typemd/issues/384)）。

**操作日誌** — LLM 需要知道上次 session 做了什麼才能銜接（[#383](https://github.com/typemd/typemd/issues/383)）。

## 為什麼這條路值得走

LLM Wiki 說得好：人們放棄維護 wiki，是因為**維護成本增長得比價值快**。更新交叉引用、保持摘要同步、標記矛盾——這些苦工沒人想做。LLM 不會厭煩，一次可以觸動十幾個檔案。

但 LLM 需要**結構**才能高效工作。純粹的 Markdown 堆疊，LLM 必須每次重新理解檔案之間的關係。typemd 的 type schema、relation engine、reconciler 就是這個結構——它是 LLM 和你的知識之間的共同語言。

最重要的是，typemd vault 終究只是一個 git repo of markdown files。你可以用任何編輯器打開它，git 追蹤每一個變更。不會被鎖進任何工具。LLM 是維護者，不是守門人。

如果你也對 LLM 驅動的知識管理感興趣，歡迎到 [GitHub](https://github.com/typemd/typemd) 追蹤進度。
