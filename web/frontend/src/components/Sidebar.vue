<script setup>
import { ref, watch, inject, onMounted } from "vue";
import vault from "../lib/vault";

const props = defineProps({
  types: Array,
  selectedId: String,
});

const theme = inject("theme");
const cycleTheme = inject("cycleTheme");

const themeLabels = { warm: "Warm", dark: "Dark", light: "Light" };
const themeIcons = { warm: "◐", dark: "●", light: "○" };

const emit = defineEmits(["select", "create"]);

const expanded = ref({});
const objectsByType = ref({});
const fetched = new Set();

onMounted(() => {
  if (props.types.length > 0 && Object.keys(expanded.value).length === 0) {
    expanded.value = { [props.types[0].name]: true };
  }
});

watch(() => props.types, (types) => {
  if (types.length > 0 && Object.keys(expanded.value).length === 0) {
    expanded.value = { [types[0].name]: true };
  }
});

watch(expanded, (exp) => {
  for (const name of Object.keys(exp)) {
    if (exp[name] && !fetched.has(name)) {
      fetched.add(name);
      vault.listObjects(name).then((objs) => {
        objectsByType.value = { ...objectsByType.value, [name]: objs || [] };
      });
    }
  }
}, { deep: true });

function toggle(name) {
  expanded.value = { ...expanded.value, [name]: !expanded.value[name] };
}
</script>

<template>
  <aside class="sidebar">
    <!-- Brand -->
    <div class="sidebar-header">
      <div class="brand">TypeMD</div>
      <button class="new-object-btn" @click="emit('create', types[0]?.name || '')">
        <span class="new-object-icon">+</span>
        New Object
      </button>
    </div>

    <!-- Type list -->
    <div class="sidebar-scroll flex-1 overflow-y-auto px-3 py-2">
      <div v-for="type in types" :key="type.name" class="type-group">
        <!-- Type group header -->
        <button class="type-header" @click="toggle(type.name)">
          <span class="type-arrow" :class="{ 'type-arrow--open': expanded[type.name] }">▶</span>
          <span v-if="type.emoji" class="type-emoji">{{ type.emoji }}</span>
          <span v-else class="type-dot">●</span>
          <span class="type-label">{{ type.plural || type.name }}</span>
          <span class="type-count">{{ type.count }}</span>
        </button>

        <!-- Objects under this type -->
        <div v-if="expanded[type.name]" class="object-list">
          <button
            v-for="obj in (objectsByType[type.name] || [])"
            :key="obj.id"
            class="object-item"
            :class="{ 'object-item--selected': selectedId === obj.id }"
            @click="emit('select', obj.id)"
          >
            <span v-if="obj.locked" class="lock-icon">🔒</span>
            <span class="truncate">{{ obj.name }}</span>
          </button>

          <button class="new-type-btn" @click="emit('create', type.name)">
            + New {{ type.name }}
          </button>
        </div>
      </div>
    </div>

    <!-- Bottom bar -->
    <div class="sidebar-footer">
      <span class="footer-text">Local-first knowledge</span>
      <button class="theme-btn" @click="cycleTheme" :title="`Theme: ${themeLabels[theme]}`">
        <span class="theme-icon">{{ themeIcons[theme] }}</span>
        <span class="theme-label">{{ themeLabels[theme] }}</span>
      </button>
    </div>
  </aside>
</template>

<style scoped>
.sidebar {
  position: relative;
  width: 260px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
  background: var(--color-sidebar-bg);
}

.sidebar-header {
  padding: 1.5rem 1.25rem 1rem;
}

.brand {
  font-family: var(--font-title);
  font-size: 17px;
  font-weight: 600;
  letter-spacing: -0.02em;
  color: var(--color-sidebar-text-bright);
  margin-bottom: 1.25rem;
}

.new-object-btn {
  display: flex;
  width: 100%;
  align-items: center;
  gap: 0.625rem;
  border-radius: 0.5rem;
  padding: 0.625rem 0.75rem;
  font-size: 13px;
  font-weight: 500;
  color: var(--color-sidebar-accent);
  background: var(--color-sidebar-btn-bg);
  border: 1px solid var(--color-sidebar-btn-border);
  cursor: pointer;
  transition: all 0.2s;
}
.new-object-btn:hover {
  background: var(--color-sidebar-btn-bg-hover);
  border-color: var(--color-sidebar-btn-border-hover);
}

.new-object-icon {
  font-size: 15px;
}

.type-group {
  margin-bottom: 2px;
}

.type-header {
  display: flex;
  width: 100%;
  align-items: center;
  gap: 0.625rem;
  border-radius: 0.5rem;
  padding: 9px 0.75rem;
  text-align: left;
  font-size: 13px;
  color: var(--color-sidebar-text);
  background: transparent;
  border: none;
  cursor: pointer;
  transition: all 0.15s;
}
.type-header:hover {
  background: var(--color-sidebar-hover);
}

.type-arrow {
  width: 1rem;
  text-align: center;
  flex-shrink: 0;
  font-size: 10px;
  color: var(--color-sidebar-text-muted);
  transition: transform 0.2s;
}
.type-arrow--open {
  transform: rotate(90deg);
}

.type-emoji {
  font-size: 15px;
  width: 20px;
  text-align: center;
  flex-shrink: 0;
}

.type-dot {
  width: 20px;
  text-align: center;
  flex-shrink: 0;
  font-size: 8px;
  color: var(--color-sidebar-text-muted);
}

.type-label {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-weight: 500;
  text-transform: uppercase;
  font-size: 11px;
  letter-spacing: 0.08em;
  color: var(--color-sidebar-text-muted);
}

.type-count {
  font-size: 11px;
  font-variant-numeric: tabular-nums;
  color: var(--color-sidebar-text-muted);
  opacity: 0;
  transition: opacity 0.2s;
}
.type-header:hover .type-count {
  opacity: 1;
}

.object-list {
  margin-left: 20px;
  padding-left: 1rem;
  margin-bottom: 0.5rem;
  border-left: 1px solid var(--color-sidebar-border);
  animation: slideIn 180ms ease-out;
}

.object-item {
  display: flex;
  width: 100%;
  align-items: center;
  gap: 0.5rem;
  padding: 7px 0.5rem;
  text-align: left;
  font-size: 14px;
  color: var(--color-sidebar-text);
  background: transparent;
  border: none;
  cursor: pointer;
  transition: all 0.15s;
}
.object-item:hover {
  background: var(--color-sidebar-hover);
}
.object-item--selected {
  background: var(--color-sidebar-active);
  color: var(--color-sidebar-text-bright);
  border-radius: 4px;
}

.lock-icon {
  font-size: 10px;
  opacity: 0.4;
}

.new-type-btn {
  display: flex;
  width: 100%;
  align-items: center;
  gap: 0.5rem;
  padding: 7px 0.5rem;
  font-size: 12px;
  color: var(--color-sidebar-text-muted);
  background: transparent;
  border: none;
  cursor: pointer;
  transition: all 0.15s;
}
.new-type-btn:hover {
  color: var(--color-sidebar-accent);
}

.sidebar-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1.25rem;
  font-size: 11px;
  letter-spacing: 0.04em;
  color: var(--color-sidebar-text-muted);
  border-top: 1px solid var(--color-sidebar-border);
}

.footer-text {
  font-family: var(--font-title);
  font-style: italic;
  opacity: 0.6;
}

.theme-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 11px;
  color: var(--color-sidebar-text-muted);
  background: transparent;
  border: none;
  cursor: pointer;
  transition: all 0.15s;
}
.theme-btn:hover {
  color: var(--color-sidebar-text);
  background: var(--color-sidebar-hover);
}

.theme-icon {
  font-size: 12px;
}

.theme-label {
  letter-spacing: 0.02em;
}
</style>
