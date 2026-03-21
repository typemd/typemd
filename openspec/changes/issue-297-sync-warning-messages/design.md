## Context

The TUI's `refreshData()` method converts `SyncResult.Unresolved` into toast items using a single group key `"unresolved refs"`. The `UnresolvedRelation` struct already carries `Reason` ("not_found" or "ambiguous") and `Matches` (candidate IDs for ambiguous cases), but this information is discarded during the toast conversion. The toast widget already supports multiple groups in a single `Show()` call, rendering each group as a separate summary line.

## Goals / Non-Goals

**Goals:**
- Display distinct toast messages for "not found" and "ambiguous" unresolved references
- Leverage existing `UnresolvedRelation.Reason` field — no changes to core sync logic

**Non-Goals:**
- Changing the toast widget itself (multi-group already works)
- Adding per-item detail messages (e.g., showing each object ID individually)
- Adding success toasts for expanded references (can be a separate enhancement)

## Decisions

### Use Reason as group key

Map `UnresolvedRelation.Reason` to human-readable group labels:
- `"not_found"` → group `"not found"`
- `"ambiguous"` → group `"ambiguous"`

This produces toast output like:
```
⚠ 2 not found
⚠ 1 ambiguous
```

**Alternative considered**: Use a single group but include reason in the message text (e.g., `"3 unresolved refs (2 not found, 1 ambiguous)"`). Rejected because the toast widget's group aggregation already handles multi-line display cleanly, and separate groups are more scannable.

**Alternative considered**: Show individual item messages (e.g., `"⚠ person/nobody: no match found"`). Rejected because with many unresolved references the toast would be too tall. The aggregated count is a better fit for transient notifications; users can run `tmd validate` for full details.

## Risks / Trade-offs

- [Risk] Unknown `Reason` value from future code changes → Mitigation: use a fallback group `"unresolved"` for any reason that isn't `"not_found"` or `"ambiguous"`.
