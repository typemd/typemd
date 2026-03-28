<script setup>
import { ref, watch, nextTick, computed } from "vue";
import vault from "../lib/vault";

const props = defineProps({
  prop: Object,
  objectId: String,
  locked: Boolean,
  showDivider: Boolean,
});

const emit = defineEmits(["save"]);

const editing = ref(false);
const value = ref("");
const inputRef = ref(null);
let saving = false;

const isEditable = computed(() =>
  !props.locked && !props.prop.isReverse && !props.prop.isBacklink &&
  props.prop.key !== "created_at" && props.prop.key !== "updated_at"
);

watch(editing, async (val) => {
  if (val) {
    await nextTick();
    inputRef.value?.focus();
    inputRef.value?.select?.();
  }
});

function startEdit() {
  if (!isEditable.value) return;
  value.value = props.prop.value != null ? String(props.prop.value) : "";
  editing.value = true;
}

async function save() {
  if (saving) return;
  saving = true;
  try {
    await vault.updateProperty(props.objectId, props.prop.key, value.value);
    editing.value = false;
    emit("save");
  } catch (err) {
    console.error(err);
  } finally {
    saving = false;
  }
}

function onKeyDown(e) {
  if (e.key === "Enter") { e.preventDefault(); save(); }
  if (e.key === "Escape") editing.value = false;
}

function toggleCheckbox() {
  if (!isEditable.value) return;
  const checked = props.prop.value === true || props.prop.value === "true";
  vault.updateProperty(props.objectId, props.prop.key, checked ? "false" : "true")
    .then(() => emit("save")).catch(console.error);
}

const inputType = computed(() => {
  if (props.prop.type === "number" || props.prop.type === "integer") return "number";
  if (props.prop.type === "date") return "date";
  return "text";
});

const isChecked = computed(() => props.prop.value === true || props.prop.value === "true");

const isRelation = computed(() => props.prop.isRelation || props.prop.isReverse || props.prop.isBacklink);

const relationArrow = computed(() => {
  if (props.prop.isBacklink) return "⟵";
  if (props.prop.isReverse) return "←";
  return "→";
});

const relationLabel = computed(() => props.prop.display?.replace(/^[→←⟵] /, "") || "");

const tagItems = computed(() => {
  if (props.prop.key !== "tags" && props.prop.type !== "multi_select") return [];
  return Array.isArray(props.prop.value) ? props.prop.value.map(String) : [];
});
</script>

<template>
  <div
    class="property-row"
    :class="{ 'property-row--editable': isEditable, 'property-row--divider': showDivider }"
    @dblclick="startEdit"
  >
    <!-- Key -->
    <div class="property-key">
      <span class="property-key-text" :title="prop.key">
        <span class="property-emoji">{{ prop.emoji || "" }}</span>
        {{ prop.key }}
      </span>
    </div>

    <!-- Value -->
    <div class="property-value">
      <!-- Editing -->
      <template v-if="editing">
        <template v-if="prop.type === 'checkbox'" />
        <input
          v-else
          ref="inputRef"
          v-model="value"
          :type="inputType"
          class="property-input"
          @blur="save"
          @keydown="onKeyDown"
        />
      </template>

      <!-- Display -->
      <template v-else>
        <!-- Checkbox -->
        <label v-if="prop.type === 'checkbox'" class="checkbox" :class="{ 'checkbox--checked': isChecked }">
          <input
            type="checkbox"
            :checked="isChecked"
            :disabled="!isEditable"
            @change="toggleCheckbox"
          />
          <svg v-if="isChecked" width="11" height="9" viewBox="0 0 11 9" fill="none">
            <path d="M1 4.5L4 7.5L10 1" stroke="white" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </label>

        <!-- Relation -->
        <span v-else-if="isRelation" class="relation-badge">
          <span class="relation-arrow">{{ relationArrow }}</span>
          {{ relationLabel }}
        </span>

        <!-- Tags / multi_select -->
        <template v-else-if="prop.key === 'tags' || prop.type === 'multi_select'">
          <span v-if="tagItems.length === 0" class="value-empty">—</span>
          <div v-else class="tag-list">
            <span v-for="(item, i) in tagItems" :key="i" class="tag">{{ item }}</span>
          </div>
        </template>

        <!-- Default value -->
        <span v-else class="value-text" :class="{ 'value-empty': !prop.display }">
          {{ prop.display || "—" }}
        </span>
      </template>
    </div>
  </div>
</template>

<style scoped>
.property-row {
  display: flex;
  align-items: center;
  min-height: 46px;
  transition: background 0.1s;
}
.property-row--editable {
  cursor: pointer;
}
.property-row--editable:hover {
  background: var(--color-surface-hover);
}
.property-row--divider {
  border-bottom: 1px solid var(--color-border-subtle);
}

.property-key {
  width: 180px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.625rem 1.25rem;
  align-self: stretch;
}

.property-key-text {
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.property-emoji {
  display: inline-block;
  width: 20px;
  text-align: center;
  margin-right: 0.25rem;
  flex-shrink: 0;
}

.property-value {
  flex: 1;
  min-width: 0;
  padding: 0.625rem 1.25rem;
}

.property-input {
  width: 100%;
  border-radius: 0.5rem;
  border: 1px solid var(--color-border);
  padding: 0.5rem 0.75rem;
  font-size: 14px;
  color: var(--color-text);
  background: var(--color-bg);
  outline: none;
  transition: all 0.15s;
}

.checkbox {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  border-radius: 5px;
  background: var(--color-bg-pure);
  border: 1.5px solid var(--color-border);
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
  cursor: pointer;
  transition: all 0.15s;
}
.checkbox--checked {
  background: var(--color-accent);
  border: none;
  box-shadow: none;
}
.checkbox input {
  position: absolute;
  opacity: 0;
  width: 100%;
  height: 100%;
  cursor: pointer;
}

.relation-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  border-radius: 0.375rem;
  padding: 0.25rem 0.625rem;
  font-size: 13px;
  font-weight: 500;
  background: var(--color-accent-light);
  color: var(--color-accent-text);
  border: 1px solid var(--color-border-subtle);
}

.relation-arrow {
  font-size: 10px;
  opacity: 0.4;
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.375rem;
}

.tag {
  display: inline-block;
  border-radius: 0.375rem;
  padding: 0.125rem 0.625rem;
  font-size: 13px;
  font-weight: 500;
  background: var(--color-tag-bg);
  color: var(--color-tag-text);
}

.value-text {
  font-size: 14px;
  color: var(--color-text);
}

.value-empty {
  color: var(--color-text-light);
}
</style>
