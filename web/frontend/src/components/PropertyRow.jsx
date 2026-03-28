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

  const inputClass = "w-full rounded border px-1.5 py-0.5 text-[12px] outline-none";
  const inputStyle = { borderColor: "var(--color-divider)", background: "white", color: "var(--color-text)" };

  const toggleCheckbox = () => {
    if (!isEditable) return;
    const checked = prop.value === true || prop.value === "true";
    vault.updateProperty(objectId, prop.key, checked ? "false" : "true")
      .then(onSave).catch(console.error);
  };

  const renderEditor = () => {
    if (prop.type === "checkbox") {
      return <input ref={inputRef} type="checkbox" checked={value === "true"}
        onChange={(e) => {
          vault.updateProperty(objectId, prop.key, e.target.checked ? "true" : "false")
            .then(onSave).catch(console.error);
          setEditing(false);
        }}
        className="h-3.5 w-3.5 accent-accent" />;
    }
    if (prop.type === "date") {
      return <input ref={inputRef} type="date" value={value} onChange={(e) => setValue(e.target.value)}
        onBlur={save} onKeyDown={onKeyDown} className={inputClass} style={inputStyle} />;
    }
    return <input ref={inputRef} type={prop.type === "number" || prop.type === "integer" ? "number" : "text"}
      value={value} onChange={(e) => setValue(e.target.value)}
      onBlur={save} onKeyDown={onKeyDown} className={inputClass} style={inputStyle} />;
  };

  const renderValue = () => {
    if (prop.type === "checkbox") {
      const checked = prop.value === true || prop.value === "true";
      return <input type="checkbox" checked={checked} disabled={!isEditable}
        onChange={toggleCheckbox}
        className="h-3.5 w-3.5 accent-accent" style={{ cursor: isEditable ? "pointer" : "default" }} />;
    }
    return (
      <span className="text-[12.5px] break-words leading-snug"
        style={{ color: prop.display ? "var(--color-text)" : "var(--color-muted)" }}>
        {prop.display || "—"}
      </span>
    );
  };

  return (
    <div
      className={`px-3 py-[6px] ${isEditable ? "hover:bg-[--color-hover] cursor-pointer" : ""}`}
      style={{ borderBottom: showDivider ? "0.5px solid var(--color-divider-light)" : "none" }}
      onDoubleClick={startEdit}
    >
      <div className="mb-px truncate text-[11px]" style={{ color: "var(--color-muted)" }}>
        {prop.emoji && <span className="mr-0.5">{prop.emoji}</span>}{prop.key}
      </div>
      <div className="min-w-0">
        {editing ? renderEditor() : renderValue()}
      </div>
    </div>
  );
}
