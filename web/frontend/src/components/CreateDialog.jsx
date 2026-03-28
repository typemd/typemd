import { useState, useEffect, useRef } from "react";
import vault from "../lib/vault";

export default function CreateDialog({ types, initialType, onCreated, onClose }) {
  const [typeName, setTypeName] = useState(initialType || "");
  const [name, setName] = useState("");
  const [templates, setTemplates] = useState([]);
  const [template, setTemplate] = useState("");
  const [creating, setCreating] = useState(false);
  const [error, setError] = useState(null);
  const nameRef = useRef(null);

  useEffect(() => {
    if (!typeName) { setTemplates([]); setTemplate(""); return; }
    vault.listTemplates(typeName).then((t) => { setTemplates(t || []); setTemplate(""); }).catch(() => setTemplates([]));
  }, [typeName]);

  useEffect(() => { if (nameRef.current) nameRef.current.focus(); }, [typeName]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!typeName || !name.trim() || creating) return;
    setCreating(true); setError(null);
    try {
      const obj = await vault.createObject(typeName, name.trim(), template);
      onCreated(obj);
    } catch (err) { setError(err.message || String(err)); setCreating(false); }
  };

  const inputCls = "w-full rounded-lg border px-4 py-3 text-[15px] outline-none transition-colors";
  const inputSt = { borderColor: "var(--color-border)", color: "var(--color-text)" };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center"
      onClick={onClose} onKeyDown={(e) => e.key === "Escape" && onClose()}>
      <div className="absolute inset-0 bg-black/20" />

      <div className="relative w-[480px] rounded-2xl border bg-white shadow-2xl"
        style={{ borderColor: "var(--color-border)", animation: "slideUp 200ms cubic-bezier(0.16,1,0.3,1)" }}
        onClick={(e) => e.stopPropagation()}>

        <div className="px-7 pt-7 pb-1">
          <h2 className="text-[22px] font-semibold tracking-[-0.02em]" style={{ fontFamily: "var(--font-title)" }}>
            New Object
          </h2>
        </div>

        <form onSubmit={handleSubmit} className="px-7 py-5 space-y-5">
          <div>
            <label className="mb-2 block text-[13px] font-medium"
              style={{ color: "var(--color-text-secondary)" }}>Type</label>
            <select value={typeName} onChange={(e) => setTypeName(e.target.value)} className={inputCls} style={inputSt}>
              <option value="">Choose a type…</option>
              {types.map((t) => (
                <option key={t.name} value={t.name}>{t.emoji ? `${t.emoji}  ` : ""}{t.name}</option>
              ))}
            </select>
          </div>

          {typeName && (
            <div>
              <label className="mb-2 block text-[13px] font-medium"
                style={{ color: "var(--color-text-secondary)" }}>Name</label>
              <input ref={nameRef} type="text" value={name} onChange={(e) => setName(e.target.value)}
                placeholder={`Enter ${typeName} name…`} className={inputCls} style={inputSt} />
            </div>
          )}

          {typeName && templates.length > 0 && (
            <div>
              <label className="mb-2 block text-[13px] font-medium"
                style={{ color: "var(--color-text-secondary)" }}>Template</label>
              <select value={template} onChange={(e) => setTemplate(e.target.value)} className={inputCls} style={inputSt}>
                <option value="">None</option>
                {templates.map((t) => (<option key={t} value={t}>{t}</option>))}
              </select>
            </div>
          )}

          {error && (
            <p className="rounded-lg px-4 py-2.5 text-[14px]"
              style={{ background: "#fef2f2", color: "var(--color-red)", border: "1px solid #fee2e2" }}>
              {error}
            </p>
          )}

          <div className="flex justify-end gap-3 pt-3 pb-1">
            <button type="button" onClick={onClose}
              className="rounded-lg px-5 py-2.5 text-[14px] font-medium transition-colors hover:bg-[--color-surface-hover]"
              style={{ color: "var(--color-text-secondary)" }}>
              Cancel
            </button>
            <button type="submit" disabled={!typeName || !name.trim() || creating}
              className="rounded-lg px-5 py-2.5 text-[14px] font-medium text-white transition-colors disabled:opacity-30"
              style={{ background: "var(--color-accent)" }}>
              {creating ? "Creating…" : "Create"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
