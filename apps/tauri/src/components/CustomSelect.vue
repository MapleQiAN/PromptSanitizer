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
          fill="var(--color-primary)"
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
  padding: calc(var(--unit) * 2.5);
  background: var(--color-bg-secondary);
  border: 3px solid var(--color-border);
  border-radius: var(--radius-md);
  font-family: var(--font-body);
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-primary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: space-between;
  transition: all 0.3s var(--ease-smooth);
  outline: none;
  user-select: none;
  min-height: 44px;
}

.custom-select__trigger:hover {
  border-color: var(--color-primary-light);
  background-color: var(--color-surface);
}

.custom-select--open .custom-select__trigger,
.custom-select__trigger:focus {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px var(--color-primary-glow), var(--shadow-md);
  background-color: var(--color-surface);
  transform: scale(1.02);
}

.custom-select__value {
  flex: 1;
  text-align: left;
}

.custom-select__arrow {
  transition: transform 0.3s var(--ease-smooth);
  flex-shrink: 0;
  margin-left: calc(var(--unit) * 2);
}

.custom-select__arrow--open {
  transform: rotate(180deg);
}

.custom-select__dropdown {
  position: absolute;
  top: calc(100% + 8px);
  left: 0;
  right: 0;
  background: var(--color-surface);
  border: 3px solid var(--color-border);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-lg), 0 0 20px var(--color-primary-glow);
  z-index: 1000;
  overflow: hidden;
  margin-top: 4px;
  max-height: 300px;
  overflow-y: auto;
  backdrop-filter: blur(10px);
}

/* Custom scrollbar for dropdown */
.custom-select__dropdown::-webkit-scrollbar {
  width: 8px;
}

.custom-select__dropdown::-webkit-scrollbar-track {
  background: var(--color-bg-secondary);
  border-radius: var(--radius-sm);
}

.custom-select__dropdown::-webkit-scrollbar-thumb {
  background: var(--color-primary-light);
  border-radius: var(--radius-sm);
}

.custom-select__dropdown::-webkit-scrollbar-thumb:hover {
  background: var(--color-primary);
}

.custom-select__option {
  padding: calc(var(--unit) * 3) calc(var(--unit) * 3.5);
  font-family: var(--font-body);
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-primary);
  cursor: pointer;
  transition: all 0.2s var(--ease-smooth);
  position: relative;
  border-bottom: 1px solid var(--color-border-light);
  min-height: 44px;
  display: flex;
  align-items: center;
}

.custom-select__option:last-child {
  border-bottom: none;
}

.custom-select__option:hover {
  background: linear-gradient(90deg, var(--color-primary-light), var(--color-primary));
  color: #FFFFFF;
  transform: translateX(4px);
  border-left: 4px solid var(--color-primary);
  padding-left: calc(var(--unit) * 2.5);
}

.custom-select__option--selected {
  background: linear-gradient(90deg, var(--color-primary), var(--color-primary-dark));
  color: #FFFFFF;
  font-weight: 600;
  border-left: 4px solid var(--color-primary-dark);
  padding-left: calc(var(--unit) * 2.5);
}

.custom-select__option--selected::before {
  content: "✓";
  position: absolute;
  left: calc(var(--unit) * 2);
  font-weight: 700;
  font-size: 16px;
  color: #FFFFFF;
}

.custom-select__option--selected:hover {
  background: linear-gradient(90deg, var(--color-primary-dark), var(--color-primary));
  color: #FFFFFF;
  transform: translateX(4px);
}

/* Dropdown animation */
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.3s var(--ease-spring);
}

.dropdown-enter-from {
  opacity: 0;
  transform: translateY(-10px) scale(0.95);
}

.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-10px) scale(0.95);
}
</style>
