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

  const inputClass = "w-full rounded-lg border px-3 py-2 text-[13px] outline-none";
  const inputStyle = { borderColor: "var(--color-divider)", color: "var(--color-text)" };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center"
      onClick={onClose} onKeyDown={(e) => e.key === "Escape" && onClose()}>
      <div className="absolute inset-0 bg-black/15" />
      <div className="relative w-96 rounded-xl border bg-white shadow-2xl"
        style={{ borderColor: "var(--color-divider)", animation: "dialogIn 200ms cubic-bezier(0.16,1,0.3,1)" }}
        onClick={(e) => e.stopPropagation()}>

        <div className="px-5 pt-5 pb-1">
          <h2 className="text-base font-semibold tracking-tight">Create Object</h2>
        </div>

        <form onSubmit={handleSubmit} className="px-5 py-3 space-y-3">
          <div>
            <label className="mb-1.5 block text-[11px] font-medium" style={{ color: "var(--color-secondary)" }}>Type</label>
            <select value={typeName} onChange={(e) => setTypeName(e.target.value)} className={inputClass} style={inputStyle}>
              <option value="">Select a type...</option>
              {types.map((t) => (
                <option key={t.name} value={t.name}>{t.emoji ? `${t.emoji} ` : ""}{t.name}</option>
              ))}
            </select>
          </div>

          {typeName && (
            <div>
              <label className="mb-1.5 block text-[11px] font-medium" style={{ color: "var(--color-secondary)" }}>Name</label>
              <input ref={nameRef} type="text" value={name} onChange={(e) => setName(e.target.value)}
                placeholder={`New ${typeName}...`} className={inputClass} style={inputStyle} />
            </div>
          )}

          {typeName && templates.length > 0 && (
            <div>
              <label className="mb-1.5 block text-[11px] font-medium" style={{ color: "var(--color-secondary)" }}>Template</label>
              <select value={template} onChange={(e) => setTemplate(e.target.value)} className={inputClass} style={inputStyle}>
                <option value="">(none)</option>
                {templates.map((t) => (<option key={t} value={t}>{t}</option>))}
              </select>
            </div>
          )}

          {error && (
            <p className="rounded-md px-3 py-1.5 text-[12px] text-red-600" style={{ background: "rgba(255,59,48,0.06)" }}>
              {error}
            </p>
          )}

          <div className="flex justify-end gap-2 pt-2 pb-2">
            <button type="button" onClick={onClose}
              className="rounded-lg px-4 py-1.5 text-[12.5px] font-medium transition-colors"
              style={{ color: "var(--color-secondary)" }}
              onMouseEnter={(e) => e.currentTarget.style.background = "var(--color-hover)"}
              onMouseLeave={(e) => e.currentTarget.style.background = "transparent"}>
              Cancel
            </button>
            <button type="submit" disabled={!typeName || !name.trim() || creating}
              className="rounded-lg px-4 py-1.5 text-[12.5px] font-medium text-white disabled:opacity-40"
              style={{ background: "var(--color-accent)" }}
              onMouseEnter={(e) => { if (!e.currentTarget.disabled) e.currentTarget.style.background = "var(--color-accent-hover)"; }}
              onMouseLeave={(e) => e.currentTarget.style.background = "var(--color-accent)"}>
              {creating ? "Creating..." : "Create"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
