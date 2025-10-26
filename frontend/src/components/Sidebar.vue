<template>
  <el-aside width="250px" class="sidebar-container">
    <NavMenu />
  </el-aside>
</template>

<script setup lang="ts">
import NavMenu from '@/components/NavMenu.vue'
</script>

<style scoped>
.sidebar-container {
    background-color: white;
    box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
    height: 100vh;
    position: fixed;
    left: 0;
    top: 0;
    z-index: 100;
    display: flex;
    flex-direction: column;
}

/* 小屏下隐藏固定侧边栏（使用 Drawer） */
@media (max-width: 768px) {
  .sidebar-container {
    display: none;
  }
}
</style>
