import { useState, useEffect, useCallback } from "react";
import vault from "./lib/vault";
import Sidebar from "./components/Sidebar";
import ObjectPage from "./components/ObjectPage";
import CreateDialog from "./components/CreateDialog";

export default function App() {
  const [types, setTypes] = useState([]);
  const [selectedId, setSelectedId] = useState(null);
  const [object, setObject] = useState(null);
  const [displayProps, setDisplayProps] = useState([]);
  const [typeSchema, setTypeSchema] = useState(null);
  const [showCreate, setShowCreate] = useState(false);
  const [createType, setCreateType] = useState(null);
  const [sidebarKey, setSidebarKey] = useState(0);

  useEffect(() => {
    vault.listTypes().then(setTypes).catch(console.error);
  }, []);

  const loadObject = useCallback((id) => {
    if (!id) { setObject(null); setDisplayProps([]); setTypeSchema(null); return; }
    const typeName = id.split("/")[0];
    Promise.all([
      vault.getObject(id),
      vault.getDisplayProperties(id),
      vault.getType(typeName).catch(() => null),
    ]).then(([obj, props, schema]) => {
      setObject(obj);
      setDisplayProps(props || []);
      setTypeSchema(schema);
    }).catch(console.error);
  }, []);

  useEffect(() => { loadObject(selectedId); }, [selectedId, loadObject]);
  const refreshObject = useCallback(() => { loadObject(selectedId); }, [selectedId, loadObject]);
  const refreshTypes = useCallback(() => {
    vault.listTypes().then(setTypes).catch(console.error);
    setSidebarKey((k) => k + 1);
  }, []);

  useEffect(() => {
    const handler = (e) => {
      const tag = e.target.tagName;
      if (tag === "INPUT" || tag === "TEXTAREA" || tag === "SELECT") return;
      if (e.key === "n") { e.preventDefault(); setShowCreate(true); }
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, []);

  return (
    <div className="flex h-screen bg-white">
      <Sidebar
        key={sidebarKey}
        types={types}
        selectedId={selectedId}
        onSelect={setSelectedId}
        onCreate={(t) => { setCreateType(t); setShowCreate(true); }}
      />

      <main className="flex-1 min-w-0 overflow-y-auto">
        {object ? (
          <ObjectPage
            object={object}
            displayProps={displayProps}
            typeSchema={typeSchema}
            onSave={refreshObject}
          />
        ) : (
          <div className="flex h-full items-center justify-center">
            <div className="text-center select-none" style={{ animation: "fadeIn 300ms ease-out" }}>
              <p className="text-[15px]" style={{ color: "var(--color-text-muted)" }}>
                Select an object to get started
              </p>
            </div>
          </div>
        )}
      </main>

      {showCreate && (
        <CreateDialog
          types={types}
          initialType={createType}
          onCreated={(obj) => { setShowCreate(false); setCreateType(null); refreshTypes(); setSelectedId(obj.id); }}
          onClose={() => { setShowCreate(false); setCreateType(null); }}
        />
      )}
    </div>
  );
}
