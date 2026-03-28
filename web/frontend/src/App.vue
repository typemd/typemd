<script setup>
import { ref, watch, provide, onMounted, onUnmounted } from "vue";
import vault from "./lib/vault";
import Sidebar from "./components/Sidebar.vue";
import ObjectPage from "./components/ObjectPage.vue";
import CreateDialog from "./components/CreateDialog.vue";

const THEMES = ["warm", "dark", "light"];
const STORAGE_KEY = "typemd-theme";

const types = ref([]);
const selectedId = ref(null);
const object = ref(null);
const displayProps = ref([]);
const typeSchema = ref(null);
const showCreate = ref(false);
const createType = ref(null);
const sidebarKey = ref(0);
const theme = ref("warm");

function applyTheme(t) {
  document.documentElement.setAttribute("data-theme", t);
}

function setTheme(t) {
  theme.value = t;
  applyTheme(t);
  localStorage.setItem(STORAGE_KEY, t);
}

function cycleTheme() {
  const idx = THEMES.indexOf(theme.value);
  setTheme(THEMES[(idx + 1) % THEMES.length]);
}

provide("theme", theme);
provide("cycleTheme", cycleTheme);

onMounted(async () => {
  // 1. Apply stored theme immediately (no flash)
  const stored = localStorage.getItem(STORAGE_KEY);
  if (stored && THEMES.includes(stored)) {
    theme.value = stored;
    applyTheme(stored);
  }

  // 2. Load config from API
  try {
    const config = await vault.getConfig();
    // Only apply config theme if user hasn't overridden
    if (!stored && config.theme && THEMES.includes(config.theme)) {
      theme.value = config.theme;
      applyTheme(config.theme);
    }
  } catch {
    // Config endpoint may not exist — use default
  }

  vault.listTypes().then((t) => (types.value = t)).catch(console.error);
});

function loadObject(id) {
  if (!id) {
    object.value = null;
    displayProps.value = [];
    typeSchema.value = null;
    return;
  }
  const typeName = id.split("/")[0];
  Promise.all([
    vault.getObject(id),
    vault.getDisplayProperties(id),
    vault.getType(typeName).catch(() => null),
  ]).then(([obj, props, schema]) => {
    object.value = obj;
    displayProps.value = props || [];
    typeSchema.value = schema;
  }).catch(console.error);
}

watch(selectedId, (id) => loadObject(id));

function refreshObject() {
  loadObject(selectedId.value);
}

function refreshTypes() {
  vault.listTypes().then((t) => (types.value = t)).catch(console.error);
  sidebarKey.value++;
}

function openCreate(typeName) {
  createType.value = typeName;
  showCreate.value = true;
}

function onCreated(obj) {
  showCreate.value = false;
  createType.value = null;
  refreshTypes();
  selectedId.value = obj.id;
}

function closeCreate() {
  showCreate.value = false;
  createType.value = null;
}

function onKeydown(e) {
  const tag = e.target.tagName;
  if (tag === "INPUT" || tag === "TEXTAREA" || tag === "SELECT") return;
  if (e.key === "n") {
    e.preventDefault();
    showCreate.value = true;
  }
}

onMounted(() => window.addEventListener("keydown", onKeydown));
onUnmounted(() => window.removeEventListener("keydown", onKeydown));
</script>

<template>
  <div class="flex h-screen" :style="{ background: 'var(--color-bg)' }">
    <Sidebar
      :key="sidebarKey"
      :types="types"
      :selected-id="selectedId"
      @select="selectedId = $event"
      @create="openCreate"
    />

    <main class="flex-1 min-w-0 overflow-y-auto">
      <ObjectPage
        v-if="object"
        :object="object"
        :display-props="displayProps"
        :type-schema="typeSchema"
        @save="refreshObject"
      />
      <div v-else class="flex h-full items-center justify-center">
        <div class="text-center select-none empty-state">
          <p class="empty-title">TypeMD</p>
          <p class="empty-hint">Select an object to get started</p>
        </div>
      </div>
    </main>

    <CreateDialog
      v-if="showCreate"
      :types="types"
      :initial-type="createType"
      @created="onCreated"
      @close="closeCreate"
    />
  </div>
</template>

<style scoped>
.empty-state {
  animation: fadeIn 400ms ease-out;
}
.empty-title {
  font-family: var(--font-title);
  font-style: italic;
  font-size: 32px;
  color: var(--color-text-light);
  letter-spacing: -0.02em;
  margin-bottom: 0.75rem;
}
.empty-hint {
  font-size: 14px;
  color: var(--color-text-muted);
}
</style>
