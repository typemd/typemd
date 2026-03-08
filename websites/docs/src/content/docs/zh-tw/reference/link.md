---
title: tmd relation link
description: 在兩個 Object 之間建立 Relation。
sidebar:
  order: 6
---

在兩個 Object 之間建立 Relation。如果 schema 定義了 `bidirectional: true`，反向屬性會自動更新。

Object ID 支援前綴匹配 — 可以省略 ULID 後綴，只要前綴能唯一識別 Object 即可。

```bash
tmd relation link book/golang-in-action author person/alan-donovan
```

對於單值 relation（未設定 `multiple: true`），重新連結會覆蓋先前的值。對於多值 relation，新目標會附加到清單中。
