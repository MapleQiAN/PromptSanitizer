import { createApp } from "vue";
import App from "./App.vue";
import "./styles.css";
import naive from "naive-ui";

// 创建并挂载应用
const app = createApp(App).use(naive);
app.mount("#root");

// 应用挂载后隐藏加载屏幕
const loadingScreen = document.getElementById("loading-screen");
if (loadingScreen) {
  // 等待一个 tick 确保 DOM 已更新
  setTimeout(() => {
    loadingScreen.classList.add("hidden");
    // 延迟移除元素以允许过渡动画
    setTimeout(() => {
      loadingScreen.remove();
    }, 300);
  }, 100);
}
