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
    <aside className="flex w-56 shrink-0 flex-col overflow-y-auto border-r"
      style={{ background: "var(--color-sidebar)", borderColor: "var(--color-divider)" }}>
      <div className="px-4 py-3 text-[13px] font-semibold" style={{ color: "var(--color-secondary)" }}>
        TypeMD
      </div>

      <div className="flex-1 px-1.5 pb-3">
        {types.map((type) => (
          <div key={type.name} className="mb-px">
            <button
              onClick={() => toggle(type.name)}
              className="group flex w-full items-center gap-1.5 rounded-md px-2 py-[5px] text-left text-[12px] font-medium hover:bg-[--color-hover] transition-colors"
              style={{ color: "var(--color-secondary)" }}
            >
              <svg width="8" height="8" viewBox="0 0 8 8" className="shrink-0 transition-transform duration-150"
                style={{ transform: expanded[type.name] ? "rotate(90deg)" : "rotate(0)", color: "var(--color-muted)" }}>
                <path d="M2 1L6 4L2 7" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/>
              </svg>
              {type.emoji && <span className="text-sm leading-none">{type.emoji}</span>}
              <span className="flex-1 truncate">{type.plural || type.name}</span>
              <span className="text-[11px] tabular-nums opacity-0 group-hover:opacity-100" style={{ color: "var(--color-muted)" }}>
                {type.count}
              </span>
            </button>

            {expanded[type.name] && (
              <div className="ml-[14px] mt-px mb-1 pl-2.5 border-l" style={{ borderColor: "var(--color-divider-light)" }}>
                {(objectsByType[type.name] || []).map((obj) => {
                  const sel = selectedId === obj.id;
                  return (
                    <button key={obj.id} onClick={() => onSelect(obj.id)}
                      className={`flex w-full items-center gap-1.5 rounded-md px-2 py-[4px] text-left text-[12.5px] transition-colors ${
                        sel ? "bg-[--color-selected] font-medium" : "hover:bg-[--color-hover]"
                      }`}
                      style={{ color: sel ? "var(--color-selected-text)" : "var(--color-text)" }}
                    >
                      {obj.locked && <span className="text-[10px] opacity-40">🔒</span>}
                      <span className="truncate">{obj.name}</span>
                    </button>
                  );
                })}
                <button onClick={() => onCreate(type.name)}
                  className="flex w-full items-center gap-1 rounded-md px-2 py-[4px] text-[12px] hover:bg-[--color-hover] transition-colors"
                  style={{ color: "var(--color-muted)" }}
                >
                  <span className="text-[11px]">＋</span>
                  <span>New {type.name}</span>
                </button>
              </div>
            )}
          </div>
        ))}
      </div>
    </aside>
  );
}
