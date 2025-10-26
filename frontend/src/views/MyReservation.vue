<template>
    <Layout>
        <el-main class="content">
                <!-- 账号管理卡片 -->
                <el-card class="account-card" shadow="hover">
                    <template #header>
                        <div class="card-header">
                            <div class="header-title">
                                <el-icon><User /></el-icon>
                                <span>账号管理</span>
                            </div>
                            <el-button
                                :icon="Refresh"
                                circle
                                @click="fetchAccountList"
                                :loading="accountLoading"
                            />
                        </div>
                    </template>

                    <div class="accounts-list">
                        <div
                            v-for="account in accounts"
                            :key="account.ID"
                            class="account-item"
                        >
                            <el-input
                                v-model="account.Lable"
                                placeholder="标签"
                                :readonly="!account.isEditing"
                                size="small"
                                clearable
                                style="width: 140px"
                                @keyup.enter="handleModify(account)"
                            />
                            <el-input
                                v-model="account.Username"
                                placeholder="用户名"
                                :readonly="!account.isEditing"
                                size="small"
                                clearable
                                style="width: 140px"
                                @keyup.enter="handleModify(account)"
                            />
                            <el-input
                                v-model="account.Password"
                                :type="
                                    account.showPassword ? 'text' : 'password'
                                "
                                placeholder="密码"
                                :readonly="!account.isEditing"
                                size="small"
                                style="width: 140px"
                                show-password
                                @keyup.enter="handleModify(account)"
                            />
                            <el-button
                                :type="
                                    account.isNew
                                        ? 'success'
                                        : account.isEditing
                                          ? 'success'
                                          : 'primary'
                                "
                                :icon="
                                    account.isNew
                                        ? Plus
                                        : account.isEditing
                                          ? Check
                                          : Edit
                                "
                                size="small"
                                @click="handleModify(account)"
                            >
                                {{
                                    account.isNew
                                        ? "添加"
                                        : account.isEditing
                                          ? "更新"
                                          : "修改"
                                }}
                            </el-button>
                            <el-button
                                type="danger"
                                :icon="Delete"
                                size="small"
                                @click="handleDelete(account)"
                            >
                                删除
                            </el-button>
                        </div>
                    </div>

                    <el-button
                        class="add-btn"
                        type="success"
                        :icon="Plus"
                        circle
                        size="default"
                        @click="addNewAccount"
                    />
                </el-card>

                <!-- 统计卡片（更紧凑 + 更多指标） -->
                <el-row :gutter="16" class="stats-row">
                    <el-col :xs="12" :sm="12" :md="6" :lg="6" :xl="6">
                        <el-card
                            shadow="hover"
                            class="stat-card success-card compact"
                        >
                            <el-statistic
                                title="预约成功率"
                                :value="successRate"
                                suffix="%"
                            />
                        </el-card>
                    </el-col>
                    <el-col :xs="12" :sm="12" :md="6" :lg="6" :xl="6">
                        <el-card
                            shadow="hover"
                            class="stat-card total-card compact"
                        >
                            <el-statistic
                                title="总预约数"
                                :value="reservations.length"
                            />
                        </el-card>
                    </el-col>
                    <el-col :xs="12" :sm="12" :md="6" :lg="6" :xl="6">
                        <el-card
                            shadow="hover"
                            class="stat-card ok-card compact"
                        >
                            <el-statistic
                                title="成功次数"
                                :value="successCount"
                            />
                        </el-card>
                    </el-col>
                    <el-col :xs="12" :sm="12" :md="6" :lg="6" :xl="6">
                        <el-card
                            shadow="hover"
                            class="stat-card error-card compact"
                        >
                            <el-statistic title="失败次数" :value="failCount" />
                        </el-card>
                    </el-col>
                </el-row>

                <!-- 预约列表卡片 -->
                <el-card class="reservation-card" shadow="hover">
                    <template #header>
                        <div class="card-header">
                            <div class="header-title">
                                <el-icon><Calendar /></el-icon>
                                <span>预约列表</span>
                            </div>
                            <div class="header-actions">
                                <el-switch
                                    v-model="showAllReservations"
                                    active-text="显示所有预约"
                                    inactive-text="仅显示未来预约"
                                />
                                <el-button
                                    :icon="Refresh"
                                    circle
                                    @click="fetchReservationList"
                                    :loading="reservationLoading"
                                />
                            </div>
                        </div>
                    </template>

                    <ReservationTable
                        :reservations="filteredReservations"
                        :loading="reservationLoading"
                        @cancel="cancelReservation"
                    />
                </el-card>
        </el-main>
    </Layout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import {
    User,
    Refresh,
    Edit,
    Delete,
    Plus,
    Check,
    Calendar,
} from "@element-plus/icons-vue";
import Layout from "@/components/Layout.vue";
import ReservationTable from "@/components/ReservationTable.vue";
import { get, post, patch, del } from "@/utils/api";
import { getUsername } from "@/utils/auth";

const accounts = ref([]);
const reservations = ref([]);
const showAllReservations = ref(false);
const accountLoading = ref(false);
const reservationLoading = ref(false);

const filteredReservations = computed(() => {
    if (!Array.isArray(reservations.value)) {
        return [];
    }

    if (showAllReservations.value) {
        return [...reservations.value].sort((a: any, b: any) =>
            new Date(b.date).getTime() - new Date(a.date).getTime(),
        );
    }

    const today = new Date();
    today.setHours(0, 0, 0, 0);
    return reservations.value
        .filter((reservation: any) => {
            const reservationDate = new Date(reservation.date);
            return reservationDate.getTime() >= today.getTime();
        })
        .sort((a: any, b: any) =>
            new Date(b.date).getTime() - new Date(a.date).getTime(),
        );
});

const successRate = computed(() => {
    if (!Array.isArray(reservations.value) || reservations.value.length === 0) {
        return 0;
    }
    const totalReservations = reservations.value.length;
    const failedReservations = reservations.value.filter(
        (r) => r.orderStatus === "预约失败",
    ).length;
    const rate =
        ((totalReservations - failedReservations) / totalReservations) * 100;
    return parseFloat(rate.toFixed(1));
});

const successCount = computed(() => {
    if (!Array.isArray(reservations.value)) return 0;
    return reservations.value.filter((r) => r.orderStatus === "预约成功")
        .length;
});

const failCount = computed(() => {
    if (!Array.isArray(reservations.value)) return 0;
    return reservations.value.filter((r) => r.orderStatus === "预约失败")
        .length;
});

const fetchAccountList = async () => {
    accountLoading.value = true;
    try {
        const data = await get("/account/list");
        if (data.message === "success") {
            accounts.value = data.data.map((account) => ({
                ...account,
                isEditing: false,
                showPassword: false,
            }));
        }
    } catch (error) {
        ElMessage.error("获取账号列表失败");
    } finally {
        accountLoading.value = false;
    }
};

const fetchReservationList = async () => {
    reservationLoading.value = true;
    try {
        const data = await get("/task/list");

        if (data.message === "success" && data.data && data.data.taskInfos) {
            const mappedData = data.data.taskInfos.map((task) => {
                let taskStatus = task.IsFinished ? "已完成" : "进行中";
                let orderStatus = "预约等待";

                if (task.IsFinished) {
                    if (!task.ReservationStatus) {
                        orderStatus = "订单取消";
                    } else {
                        orderStatus =
                            task.TradeNo &&
                            task.TradeNo !== "" &&
                            task.OrderId &&
                            task.OrderId !== ""
                                ? "预约成功"
                                : "预约失败";
                    }
                }

                return {
                    TaskID: task.TaskID,
                    createTime: task.CreateTime,
                    taskStatus: taskStatus,
                    username: task.Username,
                    password: task.Password,
                    date: task.ReservationDate,
                    venue:
                        task.VenueSiteID === "23"
                            ? "风雨操场"
                            : task.VenueSiteID === "143"
                              ? "体育馆"
                              : task.VenueName || "-",
                    site: task.SiteName,
                    time: task.ReservationTime,
                    OrderId: task.OrderId,
                    orderStatus: orderStatus,
                    ReservationStatus: task.ReservationStatus,
                };
            });
            reservations.value = mappedData;
        }
    } catch (error) {
        console.error("Fetch reservation error:", error);
        ElMessage.error("获取预约列表失败");
        reservations.value = [];
    } finally {
        reservationLoading.value = false;
    }
};

const handleModify = async (account: any) => {
    if (account.isEditing) {
        try {
            const baseBody = {
                lable: account.Lable,
                username: account.Username,
                password: account.Password,
            };

            let data;
            if (account.isNew) {
                const username = getUsername();
                if (!username) {
                    ElMessage.error("请先登录后再添加账号");
                    return;
                }
                const requestBody = { ...baseBody, user: username } as any;
                data = await post("/account/add", requestBody);
            } else {
                data = await patch("/account/update", baseBody);
            }
            if (data.message === "success") {
                ElMessage.success(
                    account.isNew ? "账号添加成功" : "账号更新成功",
                );
                account.isEditing = false;
                account.isNew = false;
                await fetchAccountList();
            }
        } catch (error) {
            ElMessage.error(account.isNew ? "添加账号失败" : "更新账号失败");
        }
    } else {
        account.isEditing = true;
    }
};

const handleDelete = async (account) => {
    if (account.isNew) {
        const index = accounts.value.indexOf(account);
        if (index > -1) {
            accounts.value.splice(index, 1);
        }
        return;
    }

    try {
        await ElMessageBox.confirm("确定要删除此账号吗？", "提示", {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "warning",
        });

        const data = await del("/account/delete", {
            username: account.Username,
        });
        if (data.message === "success") {
            ElMessage.success("账号删除成功");
            await fetchAccountList();
        } else {
            ElMessage.error(data.data || "删除账号失败");
        }
    } catch (error) {
        if (error !== "cancel") {
            ElMessage.error("删除账号失败");
        }
    }
};

const addNewAccount = () => {
    const newAccount = {
        ID: Date.now(),
        Lable: "",
        Username: "",
        Password: "",
        isEditing: true,
        showPassword: false,
        isNew: true,
    };
    accounts.value.push(newAccount);
};

const cancelReservation = async (taskId) => {
    try {
        await ElMessageBox.confirm("确定要取消此预约吗？", "提示", {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "warning",
        });

        const data = await get("/task/cancel", { task_id: taskId });
        if (data.message === "success") {
            ElMessage.success(data.data || "预约取消成功");
            await fetchReservationList();
        } else {
            ElMessage.error(data.data || "取消预约失败");
        }
    } catch (error) {
        if (error !== "cancel") {
            console.error("Cancel reservation error:", error);
            ElMessage.error("取消预约失败");
        }
    }
};

onMounted(() => {
    fetchAccountList();
    fetchReservationList();
});
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

.account-card {
    margin-bottom: 24px;
}

.accounts-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.account-item {
    display: flex;
    gap: 8px;
    align-items: center;
    padding: 10px 12px;
    background-color: #fff;
    border: 1px solid var(--el-border-color-lighter);
    border-radius: 8px;
    transition: all 0.3s;
}

.account-item:hover {
    border-color: var(--el-color-primary-light-5);
    box-shadow: 0 1px 6px rgba(0, 0, 0, 0.06);
}

.add-btn {
    margin-top: 16px;
}

.stats-row {
    margin-bottom: 24px;
}

.stat-card {
    text-align: center;
}

/* 更紧凑的统计卡片尺寸 */
.stat-card :deep(.el-card__body) {
    padding: 12px 16px;
}

.stat-card :deep(.el-statistic__head) {
    font-size: 12px;
}

.stat-card :deep(.el-statistic__content) {
    font-size: 18px;
    font-weight: 600;
}

.success-card {
    background: linear-gradient(135deg, #67c23a 0%, #85ce61 100%);
    color: white;
}

.success-card :deep(.el-statistic__head) {
    color: white;
}

.success-card :deep(.el-statistic__content) {
    color: white;
}

.total-card {
    background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
    color: white;
}

.total-card :deep(.el-statistic__head) {
    color: white;
}

.total-card :deep(.el-statistic__content) {
    color: white;
}

.reservation-card {
    margin-bottom: 24px;
}

/* 其他计数卡片配色 */
.ok-card {
    background: linear-gradient(135deg, #5cb85c 0%, #78d278 100%);
    color: white;
}

.ok-card :deep(.el-statistic__head) {
    color: white;
}

.ok-card :deep(.el-statistic__content) {
    color: white;
}

.error-card {
    background: linear-gradient(
        135deg,
        var(--el-color-danger) 0%,
        #f78989 100%
    );
    color: white;
}

.error-card :deep(.el-statistic__head) {
    color: white;
}

.error-card :deep(.el-statistic__content) {
    color: white;
}

:deep(.el-card__header) {
    padding: 18px 20px;
    border-bottom: 1px solid var(--el-border-color-light);
}

:deep(.el-switch__label) {
    font-size: 14px;
}
</style>
