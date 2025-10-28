<template>
  <Layout>
    <el-main class="content">
      <el-card shadow="hover" class="reservation-card">
        <template #header>
          <div class="card-header">
            <div class="header-title">
              <el-icon><Tickets /></el-icon>
              <span>全部预约</span>
            </div>
            <div class="header-actions">
              <el-switch
                v-model="showAllReservations"
                active-text="显示全部"
                inactive-text="仅未过期"
                size="small"
              />
              <el-button :icon="Refresh" circle @click="fetchReservationList" :loading="reservationLoading" />
            </div>
          </div>
        </template>

        <ReservationTable
          :reservations="filteredReservations"
          :loading="reservationLoading"
          :showOwner="true"
          @cancel="cancelReservation"
        />
      </el-card>
    </el-main>
  </Layout>
  
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Tickets } from '@element-plus/icons-vue'
import Layout from '@/components/Layout.vue'
import ReservationTable from '@/components/ReservationTable.vue'
import { get } from '@/utils/api'

const reservations = ref<any[]>([])
const showAllReservations = ref(false)
const reservationLoading = ref(false)

const filteredReservations = computed(() => {
  if (!Array.isArray(reservations.value)) return []

  const list = reservations.value
  if (showAllReservations.value) {
    return [...list].sort((a: any, b: any) => new Date(b.date).getTime() - new Date(a.date).getTime())
  }

  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return list
    .filter((r: any) => new Date(r.date).getTime() >= today.getTime())
    .sort((a: any, b: any) => new Date(b.date).getTime() - new Date(a.date).getTime())
})

const fetchReservationList = async () => {
  reservationLoading.value = true
  try {
    const data = await get('/task/all')
    if (data.message === 'success' && data.data && data.data.taskInfos) {
      const mapped = data.data.taskInfos.map((task: any) => {
        let taskStatus = task.IsFinished ? '已完成' : '进行中'
        let orderStatus = '预约等待'

        if (task.IsFinished) {
          if (!task.ReservationStatus) {
            orderStatus = '订单取消'
          } else {
            orderStatus = task.TradeNo && task.TradeNo !== '' && task.OrderId && task.OrderId !== '' ? '预约成功' : '预约失败'
          }
        }

        return {
          owner: task.User,
          TaskID: task.TaskID,
          createTime: task.CreateTime,
          taskStatus,
          username: task.Username,
          password: task.Password,
          date: task.ReservationDate,
          venue: task.VenueSiteID === '23' ? '风雨操场' : task.VenueSiteID === '143' ? '体育馆' : (task.VenueName || '-'),
          site: task.SiteName,
          time: task.ReservationTime,
          OrderId: task.OrderId,
          orderStatus,
          ReservationStatus: task.ReservationStatus,
        }
      })
      reservations.value = mapped
    } else {
      throw new Error(data.message || '获取预约列表失败')
    }
  } catch (error: any) {
    console.error('Fetch all reservations error:', error)
    ElMessage.error(error.message || '获取预约列表失败')
    reservations.value = []
  } finally {
    reservationLoading.value = false
  }
}

const cancelReservation = async (taskId: string) => {
  try {
    await ElMessageBox.confirm('确定要取消此预约吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    const data = await get('/task/cancel', { task_id: taskId })
    if (data.message === 'success') {
      ElMessage.success(data.data || '预约取消成功')
      await fetchReservationList()
    } else {
      ElMessage.error(data.data || '取消预约失败')
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Cancel reservation error:', error)
      ElMessage.error('取消预约失败')
    }
  }
}

onMounted(() => {
  fetchReservationList()
})
</script>

<style scoped>
.content {
  padding: 24px;
  background-color: var(--background);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.reservation-card {
  margin-bottom: 24px;
}
</style>

