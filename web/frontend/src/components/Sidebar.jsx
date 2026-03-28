import { useState, useEffect, useRef } from "react";
import vault from "../lib/vault";

export default function Sidebar({ types, selectedId, onSelect, onCreate }) {
  const [expanded, setExpanded] = useState({});
  const [objectsByType, setObjectsByType] = useState({});
  const fetchedRef = useRef(new Set());

  useEffect(() => {
    if (types.length > 0 && Object.keys(expanded).length === 0) {
      setExpanded({ [types[0].name]: true });
    }
  }, [types]);

  useEffect(() => {
    for (const name of Object.keys(expanded)) {
      if (expanded[name] && !fetchedRef.current.has(name)) {
        fetchedRef.current.add(name);
        vault.listObjects(name).then((objs) => {
          setObjectsByType((prev) => ({ ...prev, [name]: objs || [] }));
        });
      }
    }
  }, [expanded]);

  const toggle = (name) => setExpanded((prev) => ({ ...prev, [name]: !prev[name] }));

  return (
    <aside className="w-[272px] shrink-0 flex flex-col h-full border-r overflow-hidden"
      style={{ background: "var(--color-sidebar-bg)", borderColor: "var(--color-border)" }}>

      {/* Top actions */}
      <div className="px-5 pt-6 pb-3">
        <button
          onClick={() => onCreate(types[0]?.name || "")}
          className="flex w-full items-center gap-2.5 rounded-lg px-3 py-2.5 text-[14px] font-medium transition-colors hover:bg-[--color-surface-hover]"
          style={{ color: "var(--color-text-secondary)" }}>
          <span className="text-[16px]">＋</span>
          New Object
        </button>
      </div>

      {/* Divider */}
      <div className="mx-5 my-1" style={{ borderTop: "1px solid var(--color-border)" }} />

      {/* Type list */}
      <div className="flex-1 overflow-y-auto px-4 py-3">
        {types.map((type) => (
          <div key={type.name} className="mb-1">
            {/* Type group header */}
            <button onClick={() => toggle(type.name)}
              className="group flex w-full items-center gap-2.5 rounded-lg px-3 py-[9px] text-left text-[14px] transition-colors hover:bg-[--color-surface-hover]"
              style={{ color: "var(--color-text)" }}>
              {type.emoji
                ? <span className="text-[16px] w-[22px] text-center shrink-0">{type.emoji}</span>
                : <span className="w-[22px] text-center shrink-0 text-[13px]" style={{ color: "var(--color-text-muted)" }}>●</span>
              }
              <span className="flex-1 truncate font-medium">{type.plural || type.name}</span>
              <span className="text-[12px] tabular-nums opacity-0 group-hover:opacity-100 transition-opacity"
                style={{ color: "var(--color-text-muted)" }}>{type.count}</span>
            </button>

            {/* Objects under this type */}
            {expanded[type.name] && (
              <div className="ml-[22px] pl-3 mb-3 border-l" style={{ borderColor: "var(--color-border-subtle)" }}>
                {(objectsByType[type.name] || []).map((obj) => {
                  const sel = selectedId === obj.id;
                  return (
                    <button key={obj.id} onClick={() => onSelect(obj.id)}
                      className={`flex w-full items-center gap-2 rounded-lg px-3 py-[8px] text-left text-[14px] transition-colors ${
                        sel ? "font-medium" : "hover:bg-[--color-surface-hover]"
                      }`}
                      style={{
                        background: sel ? "var(--color-accent-light)" : undefined,
                        color: sel ? "var(--color-accent-text)" : "var(--color-text-secondary)",
                      }}>
                      {obj.locked && <span className="text-[11px] opacity-50">🔒</span>}
                      <span className="truncate">{obj.name}</span>
                    </button>
                  );
                })}

                <button onClick={() => onCreate(type.name)}
                  className="flex w-full items-center gap-2 rounded-lg px-3 py-[8px] text-[13px] transition-colors hover:bg-[--color-surface-hover]"
                  style={{ color: "var(--color-text-muted)" }}>
                  <span>+ New {type.name}</span>
                </button>
              </div>
            )}
          </div>
        ))}
      </div>
    </aside>
  );
}
