import PropertyRow from "./PropertyRow";

export default function Properties({ object, displayProps, onSave }) {
  const pinned = displayProps.filter((p) => p.pin && p.pin > 0).sort((a, b) => a.pin - b.pin);
  const unpinned = displayProps.filter((p) => !p.pin || p.pin === 0);

  const Section = ({ title, items }) => (
    <div className="mb-3">
      {title && (
        <div className="px-3 pb-1 text-[10px] font-semibold uppercase tracking-widest"
          style={{ color: "var(--color-muted)" }}>{title}</div>
      )}
      <div className="mx-2 rounded-lg border overflow-hidden"
        style={{ borderColor: "var(--color-divider-light)", background: "rgba(255,255,255,0.6)" }}>
        {items.map((prop, i) => (
          <PropertyRow key={prop.key} prop={prop} objectId={object.id}
            locked={object.locked} onSave={onSave} showDivider={i < items.length - 1} />
        ))}
      </div>
    </div>
  );

  return (
    <aside className="flex w-72 shrink-0 flex-col overflow-y-auto"
      style={{ background: "var(--color-props)" }}>
      <div className="px-3 pt-3 pb-2 text-[10.5px] font-semibold uppercase tracking-[0.08em]"
        style={{ color: "var(--color-muted)" }}>
        Properties
      </div>
      {pinned.length > 0 && <Section items={pinned} />}
      {unpinned.length > 0 && <Section items={unpinned} title={pinned.length > 0 ? "Details" : undefined} />}
    </aside>
  );
}
