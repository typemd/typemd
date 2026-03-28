import { useState, useEffect, useCallback, useRef } from "react";
import vault from "./lib/vault";
import Sidebar from "./components/Sidebar";
import Body from "./components/Body";
import Properties from "./components/Properties";
import CreateDialog from "./components/CreateDialog";

export default function App() {
  const [types, setTypes] = useState([]);
  const [selectedId, setSelectedId] = useState(null);
  const [object, setObject] = useState(null);
  const [displayProps, setDisplayProps] = useState([]);
  const [showProps, setShowProps] = useState(true);
  const [focusMode, setFocusMode] = useState(false);
  const [showCreate, setShowCreate] = useState(false);
  const [createType, setCreateType] = useState(null);
  const [sidebarKey, setSidebarKey] = useState(0);

  useEffect(() => {
    vault.listTypes().then(setTypes).catch(console.error);
  }, []);

  const loadObject = useCallback((id) => {
    if (!id) { setObject(null); setDisplayProps([]); return; }
    Promise.all([
      vault.getObject(id),
      vault.getDisplayProperties(id),
    ]).then(([obj, props]) => {
      setObject(obj);
      setDisplayProps(props || []);
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
      if (e.key === ".") { e.preventDefault(); setFocusMode((f) => !f); }
      if (e.key === "p") { e.preventDefault(); setShowProps((s) => !s); }
      if (e.key === "n") { e.preventDefault(); setShowCreate(true); }
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, []);

  return (
    <div className="flex h-screen" style={{ background: "var(--color-bg)" }}>
      {!focusMode && (
        <Sidebar
          key={sidebarKey}
          types={types}
          selectedId={selectedId}
          onSelect={setSelectedId}
          onCreate={(t) => { setCreateType(t); setShowCreate(true); }}
        />
      )}

      <Body object={object} onSave={refreshObject} />

      {!focusMode && showProps && object && (
        <Properties
          object={object}
          displayProps={displayProps}
          onSave={refreshObject}
        />
      )}

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
