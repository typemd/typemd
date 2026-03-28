import { useState, useRef, useEffect } from "react";
import vault from "../lib/vault";

export default function PropertyRow({ prop, objectId, locked, onSave, showDivider }) {
  const [editing, setEditing] = useState(false);
  const [value, setValue] = useState("");
  const inputRef = useRef(null);

  const isEditable = !locked && !prop.isReverse && !prop.isBacklink &&
    prop.key !== "created_at" && prop.key !== "updated_at";

  useEffect(() => {
    if (editing && inputRef.current) {
      inputRef.current.focus();
      if (inputRef.current.select) inputRef.current.select();
    }
  }, [editing]);

  const startEdit = () => {
    if (!isEditable) return;
    setValue(prop.value != null ? String(prop.value) : "");
    setEditing(true);
  };

  const save = async () => {
    try {
      await vault.updateProperty(objectId, prop.key, value);
      setEditing(false);
      onSave();
    } catch (err) { console.error(err); }
  };

  const onKeyDown = (e) => {
    if (e.key === "Enter") { e.preventDefault(); save(); }
    if (e.key === "Escape") setEditing(false);
  };

  const toggleCheckbox = () => {
    if (!isEditable) return;
    const checked = prop.value === true || prop.value === "true";
    vault.updateProperty(objectId, prop.key, checked ? "false" : "true")
      .then(onSave).catch(console.error);
  };

  const renderEditor = () => {
    const cls = "w-full rounded-lg border px-3 py-2 text-[14px] outline-none";
    const st = { borderColor: "var(--color-border)", background: "white", color: "var(--color-text)" };

    if (prop.type === "checkbox") {
      return <input ref={inputRef} type="checkbox" checked={value === "true"}
        onChange={(e) => {
          vault.updateProperty(objectId, prop.key, e.target.checked ? "true" : "false")
            .then(onSave).catch(console.error);
          setEditing(false);
        }}
        className="h-5 w-5 rounded" style={{ accentColor: "var(--color-accent)" }} />;
    }
    if (prop.type === "date") {
      return <input ref={inputRef} type="date" value={value} onChange={(e) => setValue(e.target.value)}
        onBlur={save} onKeyDown={onKeyDown} className={cls} style={st} />;
    }
    return <input ref={inputRef} type={prop.type === "number" || prop.type === "integer" ? "number" : "text"}
      value={value} onChange={(e) => setValue(e.target.value)}
      onBlur={save} onKeyDown={onKeyDown} className={cls} style={st} />;
  };

  const renderValue = () => {
    if (prop.type === "checkbox") {
      const checked = prop.value === true || prop.value === "true";
      return <input type="checkbox" checked={checked} disabled={!isEditable}
        onChange={toggleCheckbox}
        className="h-5 w-5 rounded" style={{ accentColor: "var(--color-accent)", cursor: isEditable ? "pointer" : "default" }} />;
    }

    if (prop.isRelation || prop.isReverse || prop.isBacklink) {
      const arrow = prop.isBacklink ? "⟵" : prop.isReverse ? "←" : "→";
      const label = prop.display?.replace(/^[→←⟵] /, "") || "";
      return (
        <span className="inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-[13px] font-medium"
          style={{ background: "var(--color-accent-light)", color: "var(--color-accent-text)" }}>
          <span className="text-[11px] opacity-50">{arrow}</span>
          {label}
        </span>
      );
    }

    if (prop.key === "tags" || prop.type === "multi_select") {
      const items = Array.isArray(prop.value) ? prop.value.map(String) : [];
      if (items.length === 0) return <span style={{ color: "var(--color-text-muted)" }}>—</span>;
      return (
        <div className="flex flex-wrap gap-2">
          {items.map((item, i) => (
            <span key={i} className="inline-block rounded-full px-3 py-1 text-[13px] font-medium"
              style={{ background: "var(--color-tag-bg)", color: "var(--color-tag-text)" }}>
              {item}
            </span>
          ))}
        </div>
      );
    }

    return (
      <span className="text-[14px]" style={{ color: prop.display ? "var(--color-text)" : "var(--color-text-muted)" }}>
        {prop.display || "—"}
      </span>
    );
  };

  return (
    <div
      className={`flex items-center min-h-[52px] ${isEditable ? "cursor-pointer hover:bg-[--color-surface-hover]" : ""}`}
      style={{ borderBottom: showDivider ? "1px solid var(--color-border-subtle)" : "none" }}
      onDoubleClick={startEdit}
    >
      {/* Key — with left padding and right border */}
      <div className="w-[200px] shrink-0 flex items-center gap-2 px-6 py-3 self-stretch border-r"
        style={{ borderColor: "var(--color-border-subtle)" }}>
        <span className="text-[14px] truncate" title={prop.key}
          style={{ color: "var(--color-text-secondary)" }}>
          {prop.emoji && <span className="mr-1.5">{prop.emoji}</span>}
          {prop.key}
        </span>
      </div>

      {/* Value — with left padding */}
      <div className="flex-1 min-w-0 px-6 py-3">
        {editing ? renderEditor() : renderValue()}
      </div>
    </div>
  );
}
