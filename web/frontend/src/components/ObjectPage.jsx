import { useState, useEffect, useRef } from "react";
import vault from "../lib/vault";
import PropertyRow from "./PropertyRow";

export default function ObjectPage({ object, displayProps, typeSchema, onSave }) {
  const [editing, setEditing] = useState(false);
  const [body, setBody] = useState("");
  const [saving, setSaving] = useState(false);
  const ref = useRef(null);

  useEffect(() => { setEditing(false); setBody(object?.body || ""); }, [object?.id]);
  useEffect(() => { if (editing && ref.current) ref.current.focus(); }, [editing]);

  const save = async () => {
    if (!object || saving) return;
    setSaving(true);
    try {
      await vault.updateObjectBody(object.id, body);
      setEditing(false);
      onSave();
    } catch (err) { console.error(err); }
    finally { setSaving(false); }
  };

  const onKeyDown = (e) => {
    if ((e.metaKey || e.ctrlKey) && e.key === "Enter") { e.preventDefault(); save(); }
    if (e.key === "Escape") { setEditing(false); setBody(object?.body || ""); }
  };

  const pinned = displayProps.filter((p) => p.pin && p.pin > 0).sort((a, b) => a.pin - b.pin);
  const unpinned = displayProps.filter((p) => !p.pin || p.pin === 0);
  const emoji = typeSchema?.emoji;
  const typeName = typeSchema?.name || object.type;

  return (
    <div className="max-w-[720px] mx-auto px-16 pt-20 pb-32" style={{ animation: "slideUp 250ms ease-out" }}>

      {/* Type + meta line */}
      <div className="flex items-center gap-2.5 mb-6">
        <span className="text-[15px]" style={{ color: "var(--color-text-secondary)" }}>
          {emoji && <span className="mr-1.5">{emoji}</span>}
          {typeName}
        </span>
        {object.locked && (
          <>
            <span style={{ color: "var(--color-text-light)" }}>·</span>
            <span className="text-[14px]" style={{ color: "var(--color-text-muted)" }}>🔒 Locked</span>
          </>
        )}
      </div>

      {/* Title */}
      <h1 className="text-[40px] font-semibold leading-[1.2] mb-12 tracking-[-0.025em]"
        style={{ fontFamily: "var(--font-title)", color: "var(--color-text)" }}>
        {object.name}
      </h1>

      {/* Properties */}
      {displayProps.length > 0 && (
        <div className="mb-14">
          {pinned.map((prop, i) => (
            <PropertyRow key={prop.key} prop={prop} objectId={object.id}
              locked={object.locked} onSave={onSave}
              showDivider={i < pinned.length - 1 || unpinned.length > 0} />
          ))}

          {pinned.length > 0 && unpinned.length > 0 && (
            <div className="my-1" style={{ borderTop: "2px solid var(--color-border)" }} />
          )}

          {unpinned.map((prop, i) => (
            <PropertyRow key={prop.key} prop={prop} objectId={object.id}
              locked={object.locked} onSave={onSave}
              showDivider={i < unpinned.length - 1} />
          ))}
        </div>
      )}

      {/* Divider */}
      <div className="mb-10" style={{ borderTop: "1px solid var(--color-border)" }} />

      {/* Body */}
      <div>
        {!object.locked && !editing && (
          <div className="flex justify-end mb-3">
            <button onClick={() => setEditing(true)}
              className="rounded-lg px-3 py-1.5 text-[13px] font-medium transition-colors hover:bg-[--color-surface-hover]"
              style={{ color: "var(--color-text-secondary)" }}>
              Edit
            </button>
          </div>
        )}
        {editing && (
          <div className="flex justify-end gap-2 mb-3">
            <button onClick={() => { setEditing(false); setBody(object?.body || ""); }}
              className="rounded-lg px-3 py-1.5 text-[13px] font-medium transition-colors hover:bg-[--color-surface-hover]"
              style={{ color: "var(--color-text-secondary)" }}>
              Cancel
            </button>
            <button onClick={save} disabled={saving}
              className="rounded-lg px-4 py-1.5 text-[13px] font-medium text-white transition-colors disabled:opacity-40"
              style={{ background: "var(--color-accent)" }}>
              {saving ? "Saving…" : "Save"}
            </button>
          </div>
        )}

        {editing ? (
          <textarea ref={ref} value={body} onChange={(e) => setBody(e.target.value)} onKeyDown={onKeyDown}
            className="w-full min-h-[320px] resize-y rounded-xl border p-7 text-[15px] leading-[1.9] outline-none"
            style={{ fontFamily: "var(--font-mono)", color: "var(--color-text)", borderColor: "var(--color-border)", background: "var(--color-bg)" }}
            placeholder="Write something…" spellCheck={false} />
        ) : (
          <div className="text-[15px] leading-[1.9] whitespace-pre-wrap"
            style={{ fontFamily: "var(--font-mono)", color: "var(--color-text)" }}>
            {object.body || <span style={{ color: "var(--color-text-muted)", fontStyle: "italic" }}>No content yet</span>}
          </div>
        )}
      </div>
    </div>
  );
}
