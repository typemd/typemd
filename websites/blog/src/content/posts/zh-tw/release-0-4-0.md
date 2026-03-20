---
title: "TypeMD 0.4.0 — 先把地基打好"
description: "v0.4.0 專注在核心格式、設定、與建立體驗：內建 page 型別、vault 健康檢查、設定指令、入門範本，以及全面翻新的 TUI 建立流程。"
date: 2026-03-18
tags: [發布, 開發日誌]
---

如果說 0.3.0 是「幫你的 vault 定義規則」，那 0.4.0 就是「先把地基打好」。

我們沒有急著加新功能，而是回頭把核心的東西安頓好——資料格式要穩、設定要能存、建立物件要夠快。這些東西不炫，但少了任何一個，日常使用都會卡。

## 內建 page 型別：隨手就寫

不是每件事都適合放進 `book` 或 `person`。有時候你只是想記一段文字、一個想法、一個連結。

`page` 是新的內建型別（📄），跟 `tag` 一樣不需要 YAML 檔案就存在。沒有自訂屬性，就是最單純的自由格式容器：

```bash
tmd object create page "週末出遊清單"
```

搭配 vault 設定，你甚至可以把 `page` 設為預設型別，之後建物件連型別都不用打：

```bash
tmd config set cli.default_type page
tmd object create "隨手筆記"
```

## tmd doctor：幫你的 vault 做健康檢查

Vault 用久了難免會有些髒東西——型別刪了但目錄還在、物件的型別不存在了。以前只能自己翻目錄找，現在一行搞定：

```bash
tmd doctor
```

它會掃描孤立目錄、遺失型別等完整性問題，把結果分門別類列出來。

## tmd config：設定存起來

之前很多行為只能靠 flag 控制，用完就忘。現在 vault 有了自己的設定檔（`.typemd/config.yaml`），透過 `tmd config` 管理：

```bash
tmd config set cli.default_type idea
tmd config get cli.default_type
tmd config list
```

`tmd init` 也會自動建立這個檔案，預設 `default_type: page`。

## 入門型別範本：不再從零開始

0.3.0 拿掉了內建型別，給你一個乾淨的空 vault。但對新使用者來說，空白有時候比預設更讓人不知所措。

現在 `tmd init` 會問你要不要加入入門型別：

- 💡 **idea** — 帶 status 屬性（seed → growing → ready）
- 📝 **note** — 最簡單的筆記型別
- 📚 **book** — 完整的書籍追蹤（status、rating、作者 relation）

選了就加，不選就跳過。你的 vault，你決定。

## 型別 schema 新欄位

這一版我們在型別 schema 上加了三個實用欄位：

**version** — 用 semver 格式（`"1.0"`）追蹤 schema 版本，為未來的自動 migration 打基礎：

```yaml
name: book
version: "1.0"
```

**color** — 支援預設色名或 hex 色碼，讓型別在視覺上更容易辨識：

```yaml
name: book
color: blue
```

**description** — 型別和屬性都能加描述，不用再靠註解：

```yaml
name: book
description: "追蹤閱讀清單與書評"
properties:
  - name: rating
    type: number
    description: "1-5 星評分"
```

## TUI 建立體驗大翻新

建立物件和型別的流程整個重新設計了。

**建立物件**：不再跳出輸入框，而是直接在標題面板輸入名稱。輸入時右邊即時顯示範本預覽，讓你在確認前就知道會建出什麼。名稱支援自然語言，自動轉換成合法的 slug。

**建立型別**：新的多欄位精靈，用 Tab 切換 emoji、名稱、複數名稱三個欄位。右邊即時顯示 schema 預覽，確認後自動開啟型別編輯器。

**範本管理**：在型別編輯器裡新增了 Templates 區塊，可以直接瀏覽、編輯、建立和刪除物件範本。不用離開 TUI 切到檔案系統。

## 升級方式

```bash
brew upgrade typemd/tap/typemd-cli
```

或從 [GitHub Releases](https://github.com/typemd/typemd/releases/tag/v0.4.0) 下載預編譯的執行檔。

這個版本沒有破壞性變更，直接升級即可。

## 接下來

地基打好了，接下來要開始往上蓋。0.5.0 的方向還在規劃中，但 aliases、更完整的搜尋體驗、以及 Web UI 都在我們的雷達上。

有任何問題或想法，歡迎[開 issue](https://github.com/typemd/typemd/issues)，我們都會看。
