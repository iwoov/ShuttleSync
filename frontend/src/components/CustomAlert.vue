<template>
  <transition name="fade">
    <div v-if="visible" class="alert-overlay">
      <div :class="['custom-alert', type]">
        <div class="alert-icon">
          <i class="fas" :class="iconClass"></i>
        </div>
        <span class="alert-message">{{ message }}</span>
        <button class="close-button" @click="close">&times;</button>
      </div>
    </div>
  </transition>
</template>

<script lang="ts">
export default {
  name: 'CustomAlert',
  data() {
    return {
      visible: false,
      message: '',
      type: 'info',
      timeout: null
    }
  },
  computed: {
    iconClass() {
      switch(this.type) {
        case 'success':
          return 'fa-check-circle'
        case 'warning':
          return 'fa-exclamation-triangle'
        case 'error':
          return 'fa-times-circle'
        default:
          return 'fa-info-circle'
      }
    }
  },
  methods: {
    show(message, type = 'info') {
      this.message = message
      this.type = type
      this.visible = true

      if (this.timeout) {
        clearTimeout(this.timeout)
      }

      this.timeout = setTimeout(() => {
        this.close()
      }, 3000)
    },
    close() {
      this.visible = false
    }
  }
}
</script>

<style scoped>
.alert-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding-top: 100px;
  background-color: rgba(0, 0, 0, 0.2);
  z-index: 1000;
}

.custom-alert {
  padding: 16px 24px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 300px;
  max-width: 500px;
  background-color: white;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.alert-icon {
  font-size: 24px;
  flex-shrink: 0;
}

.alert-message {
  flex-grow: 1;
  font-size: 16px;
}

.close-button {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  padding: 0 5px;
  flex-shrink: 0;
  opacity: 0.6;
  transition: opacity 0.3s;
}

.close-button:hover {
  opacity: 1;
}

/* 成功状态 */
.custom-alert.success {
  border-left: 4px solid #4caf50;
}
.custom-alert.success .alert-icon {
  color: #4caf50;
}
.custom-alert.success .close-button {
  color: #4caf50;
}

/* 警告状态 */
.custom-alert.warning {
  border-left: 4px solid #ff9800;
}
.custom-alert.warning .alert-icon {
  color: #ff9800;
}
.custom-alert.warning .close-button {
  color: #ff9800;
}

/* 错误状态 */
.custom-alert.error {
  border-left: 4px solid #f44336;
}
.custom-alert.error .alert-icon {
  color: #f44336;
}
.custom-alert.error .close-button {
  color: #f44336;
}

/* 淡入淡出动画 */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s, transform 0.3s;
}

.fade-enter-from, .fade-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}
</style>
