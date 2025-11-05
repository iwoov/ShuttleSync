<template>
  <div class="nav-menu">
    <div class="sidebar-header">
      <el-icon :size="32" color="#409eff">
        <Calendar />
      </el-icon>
      <h2>预约系统</h2>
    </div>

    <el-menu
      :default-active="activeRoute"
      class="sidebar-menu"
      :router="true"
      @select="handleSelect"
    >
      <el-menu-item index="/my-reservation">
        <el-icon><Calendar /></el-icon>
        <span>我的预约</span>
      </el-menu-item>

      <el-menu-item index="/venue-reservation">
        <el-icon><OfficeBuilding /></el-icon>
        <span>场馆预约</span>
      </el-menu-item>

      <el-menu-item index="/bargain-mode">
        <el-icon><Search /></el-icon>
        <span>捡漏模式</span>
      </el-menu-item>

      <el-menu-item index="/user-management" v-if="isAdmin">
        <el-icon><UserFilled /></el-icon>
        <span>用户管理</span>
      </el-menu-item>

      <el-menu-item index="/all-reservations" v-if="isAdmin">
        <el-icon><Tickets /></el-icon>
        <span>全部预约</span>
      </el-menu-item>

      <el-menu-item index="/my-account">
        <el-icon><User /></el-icon>
        <span>我的账户</span>
      </el-menu-item>
    </el-menu>

    <div class="sidebar-footer">
      <el-button
        type="danger"
        :icon="SwitchButton"
        @click="handleLogout"
        style="width: 100%"
      >
        退出登录
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Calendar, OfficeBuilding, UserFilled, User, SwitchButton, Tickets, Search } from '@element-plus/icons-vue'
import { logout } from '@/utils/api'
import { isAdmin as checkIsAdmin } from '@/utils/auth'

const route = useRoute()
const router = useRouter()
const emit = defineEmits(['navigate'])

const isAdmin = computed(() => checkIsAdmin())
const activeRoute = ref(route.path)

watch(
  () => route.path,
  (newPath) => {
    activeRoute.value = newPath
  }
)

const handleSelect = (index: string) => {
  activeRoute.value = index
  // 选择路由后通知上层（用于关闭 Drawer）
  emit('navigate')
}

const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    await logout()

    ElMessage.success('已退出登录')
    emit('navigate')
    router.push('/login')
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('登出错误:', error)
    }
  }
}
</script>

<style scoped>
.nav-menu {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 100%;
}

.sidebar-header {
  padding: 30px 20px;
  text-align: center;
  border-bottom: 1px solid var(--el-border-color-light);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.sidebar-header h2 {
  color: var(--el-color-primary);
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.sidebar-menu {
  border-right: none;
  padding: 10px 0;
  flex: 1;
  overflow-y: auto;
}

:deep(.el-menu-item) {
  margin: 4px 12px;
  border-radius: 8px;
  transition: all 0.3s;
}

:deep(.el-menu-item:hover) {
  background-color: var(--el-color-primary-light-9);
}

:deep(.el-menu-item.is-active) {
  background-color: var(--el-color-primary);
  color: white;
}

:deep(.el-menu-item.is-active .el-icon) {
  color: white;
}

.sidebar-footer {
  padding: 20px;
  border-top: 1px solid var(--el-border-color-light);
  margin-top: auto;
}
</style>
