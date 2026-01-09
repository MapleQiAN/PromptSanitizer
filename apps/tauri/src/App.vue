<template>
  <n-config-provider :theme="naiveTheme">
    <n-message-provider>
      <AppContent />
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from "vue";
import { darkTheme, type GlobalTheme } from "naive-ui";
import AppContent from "./components/AppContent.vue";

// 主题风格配置
type ThemeStyle = "kawaii" | "horizon";
const themeStyle = ref<ThemeStyle>("kawaii");
const theme = ref<"dark" | "light">("light");

const applyTheme = (colorMode: "dark" | "light", style: ThemeStyle) => {
  const body = document.body;
  body.classList.remove("theme-dark", "theme-light", "style-kawaii", "style-horizon");
  body.classList.add(colorMode === "dark" ? "theme-dark" : "theme-light");
  body.classList.add(style === "kawaii" ? "style-kawaii" : "style-horizon");
};

onMounted(() => {
  const prefersDark = window.matchMedia?.("(prefers-color-scheme: dark)")?.matches;
  theme.value = prefersDark ? "dark" : "light";
  applyTheme(theme.value, themeStyle.value);
});

watch(
  [theme, themeStyle],
  ([newTheme, newStyle]) => {
    applyTheme(newTheme, newStyle);
  },
  { flush: "post" }
);

const naiveTheme = computed<GlobalTheme | null>(() =>
  theme.value === "dark" ? darkTheme : null
);
</script>
