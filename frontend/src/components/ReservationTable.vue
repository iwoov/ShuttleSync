<template>
    <div class="table-wrapper">
        <el-table
            :data="reservations"
            stripe
            border
            style="width: 100%"
            :empty-text="'暂无预约数据'"
            v-loading="loading"
        >
            <el-table-column
                prop="TaskID"
                label="任务ID"
                width="180"
                align="center"
            />

            <el-table-column label="创建时间" width="180" align="center">
                <template #default="scope">
                    {{ formatDateTime(scope.row.createTime) }}
                </template>
            </el-table-column>

            <el-table-column
                prop="taskStatus"
                label="预约状态"
                width="120"
                align="center"
            >
                <template #default="scope">
                    <el-tag
                        :type="getStatusType(scope.row.taskStatus)"
                        effect="dark"
                        size="small"
                    >
                        {{ scope.row.taskStatus }}
                    </el-tag>
                </template>
            </el-table-column>

            <el-table-column
                prop="username"
                label="预约账号"
                width="120"
                align="center"
            />

            <el-table-column
                v-if="showOwner"
                prop="owner"
                label="所属用户"
                width="140"
                align="center"
            />

            <el-table-column
                prop="date"
                label="预约日期"
                width="120"
                align="center"
            >
                <template #default="scope">
                    {{ scope.row.date?.split("T")[0] || "-" }}
                </template>
            </el-table-column>

            <el-table-column
                prop="venue"
                label="预约场馆"
                width="180"
                align="center"
            />

            <el-table-column
                prop="site"
                label="预约场地"
                width="100"
                align="center"
            />

            <el-table-column
                prop="time"
                label="预约时间"
                width="120"
                align="center"
            />

            <el-table-column label="预约码" width="150" align="center">
                <template #default="scope">
                    <div
                        v-if="
                            scope.row.orderStatus === '预约成功' &&
                            scope.row.ReservationStatus
                        "
                    >
                        <div
                            v-if="
                                isReservationExpired(
                                    scope.row.date,
                                    scope.row.time,
                                )
                            "
                        >
                            <el-tag type="info" effect="plain" size="small"
                                >已过期</el-tag
                            >
                        </div>
                        <div v-else>
                            <el-button
                                v-if="!scope.row.showCode"
                                type="primary"
                                size="small"
                                :icon="View"
                                @click="fetchQRCode(scope.row)"
                            >
                                点击显示
                            </el-button>
                            <el-image
                                v-else
                                :src="
                                    'data:image/png;base64,' +
                                    scope.row.qrCodeBase64
                                "
                                style="width: 120px; height: 120px"
                                :preview-src-list="[
                                    'data:image/png;base64,' +
                                        scope.row.qrCodeBase64,
                                ]"
                                fit="contain"
                            />
                        </div>
                    </div>
                    <span v-else>-</span>
                </template>
            </el-table-column>

            <el-table-column
                prop="orderStatus"
                label="订单状态"
                width="120"
                align="center"
            >
                <template #default="scope">
                    <el-tag
                        :type="getOrderStatusType(scope.row.orderStatus)"
                        effect="plain"
                        size="small"
                    >
                        {{ scope.row.orderStatus }}
                    </el-tag>
                </template>
            </el-table-column>

            <el-table-column
                label="订单操作"
                width="100"
                align="center"
                fixed="right"
            >
                <template #default="scope">
                    <div
                        v-if="
                            scope.row.orderStatus === '预约成功' &&
                            scope.row.ReservationStatus
                        "
                    >
                        <div
                            v-if="
                                isReservationExpired(
                                    scope.row.date,
                                    scope.row.time,
                                )
                            "
                        >
                            <span>-</span>
                        </div>
                        <el-button
                            v-else
                            type="danger"
                            size="small"
                            :icon="Delete"
                            @click="handleCancel(scope.row.TaskID)"
                        >
                            取消
                        </el-button>
                    </div>
                    <span v-else>-</span>
                </template>
            </el-table-column>
        </el-table>
    </div>
</template>

<script setup lang="ts">
import { ElMessage } from "element-plus";
import { View, Delete } from "@element-plus/icons-vue";
import { get } from "@/utils/api";

defineProps({
    reservations: {
        type: Array,
        required: true,
        default: () => [],
    },
    loading: {
        type: Boolean,
        default: false,
    },
    showOwner: {
        type: Boolean,
        default: false,
    },
});

const emit = defineEmits(["cancel"]);

const formatDateTime = (dateTimeStr) => {
    if (!dateTimeStr) return "-";

    try {
        const date = new Date(dateTimeStr);
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, "0");
        const day = String(date.getDate()).padStart(2, "0");
        const hours = String(date.getHours()).padStart(2, "0");
        const minutes = String(date.getMinutes()).padStart(2, "0");
        const seconds = String(date.getSeconds()).padStart(2, "0");

        return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    } catch (error) {
        return dateTimeStr;
    }
};

const isReservationExpired = (date, time) => {
    const now = new Date();
    const [reservationDate] = date.split("T");
    const [hours, minutes] = time.split(":");
    const reservationDateTime = new Date(
        `${reservationDate}T${hours}:${minutes}:00`,
    );
    return now > reservationDateTime;
};

const getStatusType = (status) => {
    const statusMap = {
        等待中: "info",
        执行中: "warning",
        执行完成: "success",
        失败: "danger",
    };
    return statusMap[status] || "info";
};

const getOrderStatusType = (status) => {
    const statusMap = {
        预约成功: "success",
        预约失败: "danger",
        等待预约: "warning",
        已取消: "info",
    };
    return statusMap[status] || "info";
};

const fetchQRCode = async (reservation) => {
    try {
        const data = await get("/tyys/qr_code", {
            username: reservation.username,
            password: reservation.password,
            order_id: reservation.OrderId,
        });

        if (data.message === "success") {
            reservation.qrCodeBase64 = data.data;
            reservation.showCode = true;
        } else {
            throw new Error(data.message || "获取二维码失败");
        }
    } catch (error) {
        console.error("Fetch QR code error:", error);
        ElMessage.error("获取二维码失败: " + (error.message || "请稍后重试"));
    }
};

const handleCancel = (taskId) => {
    emit("cancel", taskId);
};
</script>

<style scoped>
.table-wrapper {
    width: 100%;
    overflow-x: auto;
}

:deep(.el-table) {
    font-size: 14px;
}

:deep(.el-table th) {
    background-color: #f5f7fa;
    font-weight: 600;
}

:deep(.el-table td) {
    padding: 12px 0;
}

:deep(.el-table__empty-text) {
    color: var(--el-text-color-secondary);
}
</style>
