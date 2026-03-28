// Vault API adapter — all data access goes through here.
// Swap this implementation for different backends (HTTP, static, Wails).

const API = "/api";

async function request(path, options = {}) {
  const headers = options.method && options.method !== "GET"
    ? { "Content-Type": "application/json", ...options.headers }
    : { ...options.headers };
  const res = await fetch(`${API}${path}`, { ...options, headers });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(err.error || res.statusText);
  }
  if (res.status === 204) return null;
  return res.json();
}

const vault = {
  listTypes: () => request("/types"),

  getType: (name) => request(`/types/${name}`),

  listObjects: (type) =>
    request(`/objects${type ? `?type=${encodeURIComponent(type)}` : ""}`),

  // id format: "type/slug" e.g. "book/clean-code-01xxx"
  getObject: (id) => request(`/objects/${id}`),

  getDisplayProperties: (id) => request(`/properties/${id}`),

  updateObjectBody: (id, body) =>
    request(`/objects/${id}`, {
      method: "PUT",
      body: JSON.stringify({ body }),
    }),

  updateProperty: (id, key, value) =>
    request(`/properties/${id}/${key}`, {
      method: "PUT",
      body: JSON.stringify({ value }),
    }),

  createObject: (type, name, template) =>
    request("/objects", {
      method: "POST",
      body: JSON.stringify({ type, name, template }),
    }),

  listTemplates: (type) => request(`/templates/${type}`),

  getConfig: () => request("/config"),
};

export default vault;
