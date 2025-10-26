<template>
  <div class="layout-container">
    <div id="particles-js"></div>

    <!-- Desktop fixed sidebar -->
    <Sidebar v-if="showNav" class="desktop-sidebar" />

    <!-- Mobile header -->
    <header class="mobile-header" v-if="showNav">
      <el-button type="primary" :icon="Menu" circle @click="drawerOpen = true" />
      <div class="app-title">预约系统</div>
    </header>

    <!-- Mobile drawer menu -->
    <el-drawer
      v-if="showNav"
      v-model="drawerOpen"
      size="250px"
      direction="ltr"
      :with-header="false"
      class="mobile-drawer"
    >
      <NavMenu @navigate="drawerOpen = false" />
    </el-drawer>

    <div class="layout-content" :class="{ 'with-sidebar': showNav }">
      <slot />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { Menu } from '@element-plus/icons-vue'
import Sidebar from '@/components/Sidebar.vue'
import NavMenu from '@/components/NavMenu.vue'

const drawerOpen = ref(false)
const route = useRoute()
const showNav = computed(() => route.path !== '/login')

const initParticles = () => {
  if ((window as any).particlesJS) {
    (window as any).particlesJS('particles-js', {
      particles: {
        number: { value: 80, density: { enable: true, value_area: 800 } },
        color: { value: '#4a90e2' },
        shape: { type: 'circle' },
        opacity: { value: 0.5, random: true },
        size: { value: 3, random: true },
        line_linked: { enable: true, distance: 150, color: '#4a90e2', opacity: 0.4, width: 1 },
        move: { enable: true, speed: 2, direction: 'none', random: true, out_mode: 'out' }
      },
      interactivity: {
        detect_on: 'canvas',
        events: { onhover: { enable: true, mode: 'repulse' }, onclick: { enable: true, mode: 'push' }, resize: true },
        modes: { repulse: { distance: 100, duration: 0.4 }, push: { particles_nb: 4 } }
      },
      retina_detect: true
    })
  }
}

onMounted(() => initParticles())
</script>

<style scoped>
.layout-container {
  min-height: 100vh;
  position: relative;
}

#particles-js {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
  pointer-events: none;
  background-color: #f0f4f8;
}

/* content wrapper: leave room for desktop sidebar */
.layout-content {
  position: relative;
  z-index: 1;
  margin-left: 0;
}

/* mobile header only on small screens */
.mobile-header {
  display: none;
}

.layout-content.with-sidebar {
  margin-left: 250px;
}

@media (max-width: 768px) {
  .layout-content.with-sidebar {
    margin-left: 0;
  }
  .layout-content {
    margin-left: 0;
  }

  .mobile-header {
    position: sticky;
    top: 0;
    z-index: 10;
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px 12px;
    background: rgba(255, 255, 255, 0.8);
    backdrop-filter: saturate(180%) blur(8px);
    border-bottom: 1px solid var(--el-border-color-lighter);
  }

  .app-title {
    font-size: 16px;
    font-weight: 600;
    color: var(--text-primary);
  }
}

/* Make drawer body full-height so NavMenu footer can stick bottom */
:deep(.mobile-drawer .el-drawer__body) {
  padding: 0;
  height: 100%;
  display: flex;
}
</style>
