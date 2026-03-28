<script setup>
import { ref, watch, nextTick } from "vue";
import vault from "../lib/vault";

const props = defineProps({
  types: Array,
  initialType: String,
});

const emit = defineEmits(["created", "close"]);

const typeName = ref(props.initialType || "");
const name = ref("");
const templates = ref([]);
const template = ref("");
const creating = ref(false);
const error = ref(null);
const nameRef = ref(null);

watch(typeName, async (val) => {
  if (!val) {
    templates.value = [];
    template.value = "";
    return;
  }
  try {
    const t = await vault.listTemplates(val);
    templates.value = t || [];
    template.value = "";
  } catch {
    templates.value = [];
  }
  await nextTick();
  nameRef.value?.focus();
});

async function handleSubmit() {
  if (!typeName.value || !name.value.trim() || creating.value) return;
  creating.value = true;
  error.value = null;
  try {
    const obj = await vault.createObject(typeName.value, name.value.trim(), template.value);
    emit("created", obj);
  } catch (err) {
    error.value = err.message || String(err);
    creating.value = false;
  }
}

function onOverlayKeydown(e) {
  if (e.key === "Escape") emit("close");
}
</script>

<template>
  <div class="overlay" @click="emit('close')" @keydown="onOverlayKeydown">
    <div class="backdrop" />

    <div class="dialog" @click.stop>
      <div class="dialog-header">
        <h2 class="dialog-title">New Object</h2>
      </div>

      <form class="dialog-form" @submit.prevent="handleSubmit">
        <div class="field">
          <label class="field-label">Type</label>
          <select v-model="typeName" class="field-input">
            <option value="">Choose a type…</option>
            <option v-for="t in types" :key="t.name" :value="t.name">
              {{ t.emoji ? `${t.emoji}  ` : "" }}{{ t.name }}
            </option>
          </select>
        </div>

        <div v-if="typeName" class="field field--animated">
          <label class="field-label">Name</label>
          <input
            ref="nameRef"
            v-model="name"
            type="text"
            :placeholder="`Enter ${typeName} name…`"
            class="field-input"
          />
        </div>

        <div v-if="typeName && templates.length > 0" class="field field--animated">
          <label class="field-label">Template</label>
          <select v-model="template" class="field-input">
            <option value="">None</option>
            <option v-for="t in templates" :key="t" :value="t">{{ t }}</option>
          </select>
        </div>

        <p v-if="error" class="error">{{ error }}</p>

        <div class="actions">
          <button type="button" class="btn-cancel" @click="emit('close')">Cancel</button>
          <button
            type="submit"
            class="btn-create"
            :disabled="!typeName || !name.trim() || creating"
          >
            {{ creating ? "Creating…" : "Create" }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
.overlay {
  position: fixed;
  inset: 0;
  z-index: 50;
  display: flex;
  align-items: center;
  justify-content: center;
}

.backdrop {
  position: absolute;
  inset: 0;
  background: var(--color-overlay);
  backdrop-filter: blur(4px);
  animation: overlayIn 200ms ease-out;
}

.dialog {
  position: relative;
  width: 460px;
  border-radius: 1rem;
  border: 1px solid var(--color-border);
  background: var(--color-bg-pure);
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.12), 0 4px 12px rgba(0, 0, 0, 0.06);
  animation: dialogIn 280ms cubic-bezier(0.16, 1, 0.3, 1);
}

.dialog-header {
  padding: 1.75rem 1.75rem 0.5rem;
}

.dialog-title {
  font-family: var(--font-title);
  font-size: 24px;
  font-weight: 600;
  letter-spacing: -0.03em;
  color: var(--color-text);
}

.dialog-form {
  padding: 1.25rem 1.75rem 1.75rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.field-label {
  display: block;
  margin-bottom: 0.375rem;
  font-size: 12px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--color-text-muted);
}

.field--animated {
  animation: slideUp 180ms ease-out;
}

.field-input {
  width: 100%;
  border-radius: 0.5rem;
  border: 1px solid var(--color-border);
  padding: 0.625rem 1rem;
  font-size: 14px;
  color: var(--color-text);
  background: var(--color-bg-pure);
  outline: none;
  transition: all 0.15s;
}

.error {
  border-radius: 0.5rem;
  padding: 0.625rem 1rem;
  font-size: 13px;
  background: #fef2f2;
  color: var(--color-red);
  border: 1px solid #fee2e2;
}

.actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  padding-top: 1rem;
}

.btn-cancel {
  border-radius: 0.5rem;
  padding: 0.625rem 1.25rem;
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-secondary);
  background: transparent;
  border: none;
  cursor: pointer;
  transition: all 0.15s;
}
.btn-cancel:hover {
  background: var(--color-surface);
}

.btn-create {
  border-radius: 0.5rem;
  padding: 0.625rem 1.25rem;
  font-size: 13px;
  font-weight: 500;
  color: #fff;
  background: var(--color-accent);
  border: none;
  cursor: pointer;
  transition: all 0.15s;
}
.btn-create:hover {
  background: var(--color-accent-hover);
}
.btn-create:disabled {
  opacity: 0.3;
  cursor: default;
}
</style>
