import { useState, useEffect, useRef } from "react";
import vault from "../lib/vault";

export default function Body({ object, onSave }) {
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

  if (!object) {
    return (
      <main className="flex flex-1 items-center justify-center" style={{ background: "var(--color-body)" }}>
        <div className="text-center select-none">
          <div className="mb-3 text-4xl opacity-10">📄</div>
          <p className="text-[13px]" style={{ color: "var(--color-muted)" }}>Select an object to view</p>
          <p className="mt-2 text-[11px]" style={{ color: "var(--color-muted)", opacity: 0.6 }}>
            Press <kbd className="inline-block rounded border px-1.5 py-px text-[10px] font-mono"
              style={{ borderColor: "var(--color-divider)", color: "var(--color-secondary)", background: "var(--color-bg)" }}>N</kbd> to create
          </p>
        </div>
      </main>
    );
  }

  return (
    <main className="flex flex-1 flex-col min-w-0 border-r" style={{ background: "var(--color-body)", borderColor: "var(--color-divider)" }}>
      <div className="flex items-center gap-3 px-5 py-2.5 border-b" style={{ borderColor: "var(--color-divider-light)" }}>
        <h1 className="flex-1 min-w-0 truncate text-[15px] font-semibold tracking-tight">
          {object.name}
        </h1>
        <div className="flex items-center gap-1.5 shrink-0">
          {object.locked && <span className="text-xs opacity-40">🔒</span>}
          {!object.locked && !editing && (
            <button onClick={() => setEditing(true)}
              className="rounded-md border px-3 py-[3px] text-[11.5px] font-medium hover:bg-[--color-hover] transition-colors"
              style={{ borderColor: "var(--color-divider)", color: "var(--color-secondary)" }}>
              Edit
            </button>
          )}
          {editing && (
            <>
              <button onClick={() => { setEditing(false); setBody(object?.body || ""); }}
                className="rounded-md px-3 py-[3px] text-[11.5px] font-medium hover:bg-[--color-hover] transition-colors"
                style={{ color: "var(--color-secondary)" }}>
                Cancel
              </button>
              <button onClick={save} disabled={saving}
                className="rounded-md px-3 py-[3px] text-[11.5px] font-medium text-white hover:brightness-90 transition-colors disabled:opacity-40"
                style={{ background: "var(--color-accent)" }}>
                {saving ? "Saving..." : "Save"}
              </button>
            </>
          )}
        </div>
      </div>

      <div className="flex-1 overflow-y-auto">
        {editing ? (
          <textarea ref={ref} value={body} onChange={(e) => setBody(e.target.value)} onKeyDown={onKeyDown}
            className="h-full w-full resize-none border-none bg-transparent px-5 py-4 text-[13px] leading-[1.75] outline-none"
            style={{ fontFamily: "var(--font-mono)" }} placeholder="Write something..." spellCheck={false} />
        ) : (
          <div className="px-5 py-4 text-[13px] leading-[1.75] whitespace-pre-wrap"
            style={{ fontFamily: "var(--font-mono)" }}>
            {object.body || <span style={{ color: "var(--color-muted)", fontStyle: "italic" }}>No content</span>}
          </div>
        )}
      </div>
    </main>
  );
}
