<template>
  <div class="custom-select" :class="{ 'custom-select--open': isOpen }">
    <div class="custom-select__trigger" @click="toggleOpen" @blur="handleBlur">
      <span class="custom-select__value">{{ displayValue }}</span>
      <svg
        class="custom-select__arrow"
        :class="{ 'custom-select__arrow--open': isOpen }"
        width="12"
        height="12"
        viewBox="0 0 12 12"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          d="M6 9L1 4h10z"
          fill="var(--color-text-secondary)"
        />
      </svg>
    </div>
    <transition name="dropdown">
      <div v-if="isOpen" class="custom-select__dropdown">
        <div
          v-for="option in options"
          :key="option.value"
          class="custom-select__option"
          :class="{ 'custom-select__option--selected': option.value === modelValue }"
          @click="selectOption(option.value)"
        >
          {{ option.label }}
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";

// 全局状态：跟踪当前打开的下拉框实例
const openSelectInstance = ref<symbol | null>(null);

interface Option {
  value: string;
  label: string;
}

interface Props {
  modelValue: string;
  options: Option[];
}

const props = defineProps<Props>();
const emit = defineEmits<{
  "update:modelValue": [value: string];
}>();

// 为每个组件实例创建唯一标识
const instanceId = Symbol();
const isOpen = computed(() => openSelectInstance.value === instanceId);

const displayValue = computed(() => {
  const option = props.options.find((opt) => opt.value === props.modelValue);
  return option ? option.label : props.modelValue;
});

const toggleOpen = () => {
  if (isOpen.value) {
    // 关闭当前下拉框
    openSelectInstance.value = null;
  } else {
    // 关闭其他所有下拉框，然后打开当前下拉框
    openSelectInstance.value = instanceId;
  }
};

const selectOption = (value: string) => {
  emit("update:modelValue", value);
  openSelectInstance.value = null;
};

const handleClickOutside = (e: MouseEvent) => {
  const target = e.target as HTMLElement;
  if (!target.closest(".custom-select")) {
    openSelectInstance.value = null;
  }
};

onMounted(() => {
  document.addEventListener("click", handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener("click", handleClickOutside);
  // 清理：如果当前实例是打开的，关闭它
  if (openSelectInstance.value === instanceId) {
    openSelectInstance.value = null;
  }
});
</script>

<style scoped>
.custom-select {
  position: relative;
  width: 100%;
}

.custom-select__trigger {
  width: 100%;
  padding: calc(var(--unit) * 2.5) calc(var(--unit) * 3);
  background: var(--color-bg-tertiary);
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-sm);
  font-family: var(--font-body);
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-primary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: space-between;
  transition: all 0.25s var(--ease-smooth);
  outline: none;
  user-select: none;
  min-height: 40px;
}

.custom-select__trigger:hover {
  border-color: var(--color-border);
  background-color: var(--color-bg-secondary);
}

.custom-select--open .custom-select__trigger {
  border-color: var(--color-primary-light);
  background-color: var(--color-surface);
  box-shadow: 0 2px 8px rgba(255, 107, 157, 0.12);
}

.custom-select__value {
  flex: 1;
  text-align: left;
  color: var(--color-text-primary);
}

.custom-select__arrow {
  transition: transform 0.25s var(--ease-smooth);
  flex-shrink: 0;
  margin-left: calc(var(--unit) * 2);
  opacity: 0.6;
}

.custom-select__arrow--open {
  transform: rotate(180deg);
  opacity: 1;
}

.custom-select__dropdown {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  right: 0;
  background: var(--color-surface);
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-sm);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08), 0 2px 8px rgba(255, 107, 157, 0.1);
  z-index: 1000;
  overflow: hidden;
  margin-top: 2px;
  max-height: 280px;
  overflow-y: auto;
}

/* Custom scrollbar for dropdown */
.custom-select__dropdown::-webkit-scrollbar {
  width: 6px;
}

.custom-select__dropdown::-webkit-scrollbar-track {
  background: transparent;
}

.custom-select__dropdown::-webkit-scrollbar-thumb {
  background: var(--color-border);
  border-radius: var(--radius-sm);
}

.custom-select__dropdown::-webkit-scrollbar-thumb:hover {
  background: var(--color-primary-light);
}

.custom-select__option {
  padding: calc(var(--unit) * 2.5) calc(var(--unit) * 3);
  font-family: var(--font-body);
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-primary);
  cursor: pointer;
  transition: all 0.2s var(--ease-smooth);
  position: relative;
  min-height: 40px;
  display: flex;
  align-items: center;
}

.custom-select__option:not(:last-child) {
  border-bottom: 1px solid var(--color-border-light);
}

.custom-select__option:hover {
  background: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

.custom-select__option--selected {
  background: var(--color-primary);
  color: #FFFFFF !important;
  font-weight: 600;
  padding-left: calc(var(--unit) * 5);
}

.custom-select__option--selected * {
  color: #FFFFFF !important;
}

.custom-select__option--selected::before {
  content: "✓";
  position: absolute;
  left: calc(var(--unit) * 2.5);
  font-weight: 600;
  font-size: 14px;
  color: #FFFFFF !important;
  opacity: 1;
}

.custom-select__option--selected:hover {
  background: var(--color-primary-dark);
  color: #FFFFFF !important;
}

.custom-select__option--selected:hover * {
  color: #FFFFFF !important;
}

/* Dropdown animation */
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s var(--ease-smooth);
}

.dropdown-enter-from {
  opacity: 0;
  transform: translateY(-4px);
}

.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
