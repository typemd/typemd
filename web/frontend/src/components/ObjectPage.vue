<script setup>
import { ref, watch, nextTick } from "vue";
import vault from "../lib/vault";
import PropertyRow from "./PropertyRow.vue";

const props = defineProps({
  object: Object,
  displayProps: Array,
  typeSchema: Object,
});

const emit = defineEmits(["save"]);

const editing = ref(false);
const body = ref("");
const saving = ref(false);
const textareaRef = ref(null);

watch(() => props.object?.id, () => {
  editing.value = false;
  body.value = props.object?.body || "";
});

watch(editing, async (val) => {
  if (val) {
    await nextTick();
    textareaRef.value?.focus();
  }
});

async function save() {
  if (!props.object || saving.value) return;
  saving.value = true;
  try {
    await vault.updateObjectBody(props.object.id, body.value);
    editing.value = false;
    emit("save");
  } catch (err) {
    console.error(err);
  } finally {
    saving.value = false;
  }
}

function onKeyDown(e) {
  if ((e.metaKey || e.ctrlKey) && e.key === "Enter") { e.preventDefault(); save(); }
  if (e.key === "Escape") { editing.value = false; body.value = props.object?.body || ""; }
}

const pinned = ref([]);
const unpinned = ref([]);

watch(() => props.displayProps, (dp) => {
  pinned.value = dp.filter((p) => p.pin && p.pin > 0).sort((a, b) => a.pin - b.pin);
  unpinned.value = dp.filter((p) => !p.pin || p.pin === 0);
}, { immediate: true });

const emoji = ref("");
const typeName = ref("");

watch(() => props.typeSchema, (ts) => {
  emoji.value = ts?.emoji || "";
  typeName.value = ts?.name || props.object?.type || "";
}, { immediate: true });
</script>

<template>
  <div class="object-page">
    <!-- Type breadcrumb -->
    <div class="breadcrumb">
      <span class="type-badge">
        <span v-if="emoji" class="type-badge-emoji">{{ emoji }}</span>
        {{ typeName }}
      </span>
      <span v-if="object.locked" class="locked-badge">🔒 Locked</span>
    </div>

    <!-- Title -->
    <h1 class="title">{{ object.name }}</h1>

    <!-- Properties -->
    <div v-if="displayProps.length > 0" class="properties-card">
      <PropertyRow
        v-for="(prop, i) in pinned"
        :key="prop.key"
        :prop="prop"
        :object-id="object.id"
        :locked="object.locked"
        :show-divider="i < pinned.length - 1 || unpinned.length > 0"
        @save="emit('save')"
      />

      <div v-if="pinned.length > 0 && unpinned.length > 0" class="properties-divider" />

      <PropertyRow
        v-for="(prop, i) in unpinned"
        :key="prop.key"
        :prop="prop"
        :object-id="object.id"
        :locked="object.locked"
        :show-divider="i < unpinned.length - 1"
        @save="emit('save')"
      />
    </div>

    <!-- Body section -->
    <div class="body-section">
      <div class="body-header">
        <span class="body-label">Content</span>
        <button
          v-if="!object.locked && !editing"
          class="edit-btn"
          @click="editing = true"
        >
          Edit
        </button>
        <div v-if="editing" class="edit-actions">
          <button class="cancel-btn" @click="editing = false; body = object?.body || ''">
            Cancel
          </button>
          <button class="save-btn" :disabled="saving" @click="save">
            {{ saving ? "Saving…" : "Save" }}
          </button>
        </div>
      </div>

      <div class="body-card">
        <textarea
          v-if="editing"
          ref="textareaRef"
          v-model="body"
          class="body-editor"
          placeholder="Write something…"
          :spellcheck="false"
          @keydown="onKeyDown"
        />
        <div v-else class="body-content">
          <template v-if="object.body">{{ object.body }}</template>
          <span v-else class="body-empty">No content yet</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.object-page {
  max-width: 740px;
  margin: 0 auto;
  padding: 5rem 4rem 8rem;
  animation: slideUp 280ms cubic-bezier(0.16, 1, 0.3, 1);
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 2rem;
}

.type-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  border-radius: 9999px;
  padding: 0.25rem 0.75rem;
  font-size: 12px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  background: var(--color-surface);
  color: var(--color-text-secondary);
}

.type-badge-emoji {
  font-size: 13px;
  margin-left: -2px;
}

.locked-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  border-radius: 9999px;
  padding: 0.25rem 0.625rem;
  font-size: 12px;
  background: var(--color-surface);
  color: var(--color-text-muted);
}

.title {
  font-family: var(--font-title);
  font-size: 38px;
  font-weight: 600;
  line-height: 1.18;
  letter-spacing: -0.03em;
  color: var(--color-text);
  margin-bottom: 2.5rem;
}

.properties-card {
  margin-bottom: 3rem;
  border-radius: 0.75rem;
  overflow: hidden;
  border: 1px solid var(--color-border);
  background: var(--color-bg-pure);
}

.properties-divider {
  border-top: 2px solid var(--color-border);
}

.body-section {
  /* container */
}

.body-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1rem;
}

.body-label {
  font-size: 12px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--color-text-muted);
}

.edit-btn {
  border-radius: 0.375rem;
  padding: 0.25rem 0.75rem;
  font-size: 13px;
  font-weight: 500;
  color: var(--color-accent);
  background: var(--color-accent-light);
  border: none;
  cursor: pointer;
  transition: all 0.15s;
}
.edit-btn:hover {
  background: var(--color-accent);
  color: #fff;
}

.edit-actions {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.cancel-btn {
  border-radius: 0.375rem;
  padding: 0.25rem 0.75rem;
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-secondary);
  background: transparent;
  border: none;
  cursor: pointer;
  transition: background 0.15s;
}
.cancel-btn:hover {
  background: var(--color-surface);
}

.save-btn {
  border-radius: 0.375rem;
  padding: 0.25rem 1rem;
  font-size: 13px;
  font-weight: 500;
  color: #fff;
  background: var(--color-accent);
  border: none;
  cursor: pointer;
  transition: all 0.15s;
}
.save-btn:hover {
  background: var(--color-accent-hover);
}
.save-btn:disabled {
  opacity: 0.4;
  cursor: default;
}

.body-card {
  border-radius: 0.75rem;
  overflow: hidden;
  border: 1px solid var(--color-border);
  background: var(--color-bg-pure);
}

.body-editor {
  width: 100%;
  min-height: 300px;
  resize: vertical;
  padding: 1.5rem;
  font-family: var(--font-mono);
  font-size: 14px;
  line-height: 1.85;
  color: var(--color-text);
  background: var(--color-bg-pure);
  border: none;
  outline: none;
}

.body-content {
  padding: 1.5rem;
  font-family: var(--font-mono);
  font-size: 14px;
  line-height: 1.85;
  white-space: pre-wrap;
  color: var(--color-text);
}

.body-empty {
  color: var(--color-text-muted);
  font-style: italic;
}
</style>
