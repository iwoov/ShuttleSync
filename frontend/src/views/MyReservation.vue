<template>
    <Layout>
        <el-main class="content">
            <!-- 顶部统计卡片 -->
            <div class="stats-cards">
                <el-card class="stat-card" shadow="hover">
                    <div class="stat-content">
                        <div class="stat-icon">
                            <el-icon :size="32" color="#409eff"
                                ><Calendar
                            /></el-icon>
                        </div>
                        <div class="stat-info">
                            <div class="stat-label">总预约数</div>
                            <div class="stat-value">
                                {{ reservations.length }}
                            </div>
                        </div>
                    </div>
                </el-card>

                <el-card class="stat-card" shadow="hover">
                    <div class="stat-content">
                        <div class="stat-icon">
                            <el-icon :size="32" color="#67c23a"
                                ><CircleCheck
                            /></el-icon>
                        </div>
                        <div class="stat-info">
                            <div class="stat-label">成功率</div>
                            <div class="stat-value">{{ successRate }}%</div>
                        </div>
                    </div>
                </el-card>

                <el-card class="stat-card" shadow="hover">
                    <div class="stat-content">
                        <div class="stat-icon">
                            <el-icon :size="32" color="#5cb85c"
                                ><SuccessFilled
                            /></el-icon>
                        </div>
                        <div class="stat-info">
                            <div class="stat-label">成功次数</div>
                            <div class="stat-value">{{ successCount }}</div>
                        </div>
                    </div>
                </el-card>

                <el-card class="stat-card" shadow="hover">
                    <div class="stat-content">
                        <div class="stat-icon">
                            <el-icon :size="32" color="#f56c6c"
                                ><CircleClose
                            /></el-icon>
                        </div>
                        <div class="stat-info">
                            <div class="stat-label">失败次数</div>
                            <div class="stat-value">{{ failCount }}</div>
                        </div>
                    </div>
                </el-card>
            </div>

            <!-- 图表区域 -->
            <div class="charts-section">
                <!-- 预约趋势 -->
                <el-card class="chart-card chart-large" shadow="hover">
                    <template #header>
                        <div class="chart-header">
                            <span class="chart-title">预约趋势</span>
                        </div>
                    </template>
                    <div class="chart-legend">
                        <span class="legend-item">
                            <span class="legend-dot total"></span>总预约
                        </span>
                        <span class="legend-item">
                            <span class="legend-dot success"></span>成功
                        </span>
                        <span class="legend-item">
                            <span class="legend-dot fail"></span>失败
                        </span>
                    </div>
                    <div class="line-chart-container">
                        <svg viewBox="0 0 800 300" class="line-chart">
                            <!-- 网格线 -->
                            <g class="grid">
                                <line
                                    x1="60"
                                    y1="250"
                                    x2="760"
                                    y2="250"
                                    stroke="#eee"
                                    stroke-width="1"
                                />
                                <line
                                    x1="60"
                                    y1="200"
                                    x2="760"
                                    y2="200"
                                    stroke="#eee"
                                    stroke-width="1"
                                />
                                <line
                                    x1="60"
                                    y1="150"
                                    x2="760"
                                    y2="150"
                                    stroke="#eee"
                                    stroke-width="1"
                                />
                                <line
                                    x1="60"
                                    y1="100"
                                    x2="760"
                                    y2="100"
                                    stroke="#eee"
                                    stroke-width="1"
                                />
                                <line
                                    x1="60"
                                    y1="50"
                                    x2="760"
                                    y2="50"
                                    stroke="#eee"
                                    stroke-width="1"
                                />
                            </g>
                            <!-- Y轴标签 -->
                            <text x="45" y="255" class="axis-label">0</text>
                            <text x="35" y="205" class="axis-label">
                                {{ maxValue * 0.25 }}
                            </text>
                            <text x="35" y="155" class="axis-label">
                                {{ maxValue * 0.5 }}
                            </text>
                            <text x="35" y="105" class="axis-label">
                                {{ maxValue * 0.75 }}
                            </text>
                            <text x="35" y="55" class="axis-label">
                                {{ maxValue }}
                            </text>
                            <!-- 折线 -->
                            <path
                                :d="totalPath"
                                class="line-path total-line"
                                fill="none"
                            />
                            <path
                                :d="successPath"
                                class="line-path success-line"
                                fill="none"
                            />
                            <path
                                :d="failPath"
                                class="line-path fail-line"
                                fill="none"
                            />
                            <!-- X轴日期标签 -->
                            <text
                                v-for="(label, index) in dateLabels"
                                :key="index"
                                :x="
                                    60 + index * (700 / (dateLabels.length - 1))
                                "
                                y="275"
                                class="axis-label date-label"
                            >
                                {{ label }}
                            </text>
                        </svg>
                    </div>
                </el-card>

                <!-- 预约分布 -->
                <el-card class="chart-card chart-small" shadow="hover">
                    <template #header>
                        <div class="chart-header">
                            <span class="chart-title">预约分布</span>
                        </div>
                    </template>
                    <div class="donut-chart-container">
                        <svg viewBox="0 0 200 200" class="donut-chart">
                            <!-- 成功 -->
                            <path
                                v-if="successCount > 0"
                                :d="getSuccessArc()"
                                class="donut-slice success-slice"
                            />
                            <!-- 失败 -->
                            <path
                                v-if="failCount > 0"
                                :d="getFailArc()"
                                class="donut-slice fail-slice"
                            />
                            <!-- 中心圆 -->
                            <circle
                                cx="100"
                                cy="100"
                                r="50"
                                class="donut-center"
                            />
                            <!-- 百分比显示 -->
                            <text x="100" y="95" class="donut-percent-label">
                                成功率
                            </text>
                            <text x="100" y="115" class="donut-percent-value">
                                {{ successRate }}%
                            </text>
                        </svg>
                        <div class="donut-stats">
                            <div class="donut-stat-item">
                                <div class="donut-stat-label">成功预约</div>
                                <div class="donut-stat-value success-text">
                                    {{ successCount }}
                                </div>
                            </div>
                            <div class="donut-stat-item">
                                <div class="donut-stat-label">失败预约</div>
                                <div class="donut-stat-value fail-text">
                                    {{ failCount }}
                                </div>
                            </div>
                        </div>
                    </div>
                </el-card>
            </div>

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
    Refresh,
    Calendar,
    CircleCheck,
    CircleClose,
    SuccessFilled,
} from "@element-plus/icons-vue";
import Layout from "@/components/Layout.vue";
import ReservationTable from "@/components/ReservationTable.vue";
import { get } from "@/utils/api";

const reservations = ref<any[]>([]);
const showAllReservations = ref(false);
const reservationLoading = ref(false);

const filteredReservations = computed(() => {
    if (!Array.isArray(reservations.value)) {
        return [];
    }

    if (showAllReservations.value) {
        return [...reservations.value].sort(
            (a: any, b: any) =>
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
        .sort(
            (a: any, b: any) =>
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

// 统计每天的预约数据
const reservationsByDate = computed(() => {
    const dateMap = new Map<
        string,
        { total: number; success: number; fail: number }
    >();

    reservations.value.forEach((r: any) => {
        const date = r.date;
        if (!dateMap.has(date)) {
            dateMap.set(date, { total: 0, success: 0, fail: 0 });
        }
        const stats = dateMap.get(date)!;
        stats.total++;
        if (r.orderStatus === "预约成功") {
            stats.success++;
        } else if (r.orderStatus === "预约失败") {
            stats.fail++;
        }
    });

    // 按日期排序
    const sortedDates = Array.from(dateMap.keys()).sort();
    return sortedDates.map((date) => ({
        date,
        ...dateMap.get(date)!,
    }));
});

// 获取最近7天的日期
const dateLabels = computed(() => {
    if (reservationsByDate.value.length === 0) {
        // 如果没有数据，显示最近7天
        const labels: string[] = [];
        for (let i = 6; i >= 0; i--) {
            const date = new Date();
            date.setDate(date.getDate() - i);
            labels.push(`${date.getMonth() + 1}/${date.getDate()}`);
        }
        return labels;
    }

    // 取最近的日期（最多7个）
    const dates = reservationsByDate.value.slice(-7);
    return dates.map((d) => {
        const [year, month, day] = d.date.split("-");
        return `${parseInt(month)}/${parseInt(day)}`;
    });
});

const maxValue = computed(() => {
    if (reservationsByDate.value.length === 0) return 10;
    const max = Math.max(...reservationsByDate.value.map((d) => d.total));
    // 向上取整到10的倍数
    return Math.ceil(max / 10) * 10 || 10;
});

// 生成平滑曲线图路径（使用贝塞尔曲线）
const generatePath = (dataKey: "total" | "success" | "fail") => {
    const data = reservationsByDate.value.slice(-7);
    if (data.length === 0) return "M 60 250";

    const points: Array<{ x: number; y: number }> = [];
    const xStep = 700 / Math.max(data.length - 1, 1);

    data.forEach((item, index) => {
        const x = 60 + index * xStep;
        const value = item[dataKey];
        const y = 250 - (value / maxValue.value) * 200;
        points.push({ x, y });
    });

    if (points.length === 1) {
        return `M ${points[0].x} ${points[0].y}`;
    }

    // 使用三次贝塞尔曲线创建平滑路径
    let path = `M ${points[0].x} ${points[0].y}`;

    for (let i = 0; i < points.length - 1; i++) {
        const current = points[i];
        const next = points[i + 1];
        const prev = i > 0 ? points[i - 1] : current;
        const afterNext = i < points.length - 2 ? points[i + 2] : next;

        // 计算控制点（使用 Catmull-Rom 算法）
        const tension = 0.3; // 张力系数，值越小曲线越平滑

        const cp1x = current.x + (next.x - prev.x) * tension;
        const cp1y = current.y + (next.y - prev.y) * tension;
        const cp2x = next.x - (afterNext.x - current.x) * tension;
        const cp2y = next.y - (afterNext.y - current.y) * tension;

        // 使用三次贝塞尔曲线
        path += ` C ${cp1x} ${cp1y}, ${cp2x} ${cp2y}, ${next.x} ${next.y}`;
    }

    return path;
};

const totalPath = computed(() => generatePath("total"));
const successPath = computed(() => generatePath("success"));
const failPath = computed(() => generatePath("fail"));

// 计算环形图路径
const getSuccessArc = () => {
    const total = reservations.value.length;
    if (total === 0) return "";
    const successPercentage = (successCount.value / total) * 100;
    return createDonutArc(0, successPercentage);
};

const getFailArc = () => {
    const total = reservations.value.length;
    if (total === 0) return "";
    const successPercentage = (successCount.value / total) * 100;
    const failPercentage = (failCount.value / total) * 100;
    return createDonutArc(
        successPercentage,
        successPercentage + failPercentage,
    );
};

// 创建环形图弧形路径
const createDonutArc = (startPercentage: number, endPercentage: number) => {
    const cx = 100;
    const cy = 100;
    const outerRadius = 80;
    const innerRadius = 50;

    const startAngle = (startPercentage / 100) * 360 - 90;
    const endAngle = (endPercentage / 100) * 360 - 90;

    const startRad = (startAngle * Math.PI) / 180;
    const endRad = (endAngle * Math.PI) / 180;

    const x1 = cx + outerRadius * Math.cos(startRad);
    const y1 = cy + outerRadius * Math.sin(startRad);
    const x2 = cx + outerRadius * Math.cos(endRad);
    const y2 = cy + outerRadius * Math.sin(endRad);
    const x3 = cx + innerRadius * Math.cos(endRad);
    const y3 = cy + innerRadius * Math.sin(endRad);
    const x4 = cx + innerRadius * Math.cos(startRad);
    const y4 = cy + innerRadius * Math.sin(startRad);

    const largeArcFlag = endPercentage - startPercentage > 50 ? 1 : 0;

    return `M ${x1} ${y1} A ${outerRadius} ${outerRadius} 0 ${largeArcFlag} 1 ${x2} ${y2} L ${x3} ${y3} A ${innerRadius} ${innerRadius} 0 ${largeArcFlag} 0 ${x4} ${y4} Z`;
};

const fetchReservationList = async () => {
    reservationLoading.value = true;
    try {
        const data = await get("/task/list");

        if (data.message === "success" && data.data && data.data.taskInfos) {
            const mappedData = data.data.taskInfos.map((task: any) => {
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

const cancelReservation = async (taskId: string) => {
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
    fetchReservationList();
});
</script>

<style scoped>
.content {
    padding: 24px;
    background-color: #f5f7fa;
}

/* 顶部统计卡片 */
.stats-cards {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
    gap: 16px;
    margin-bottom: 24px;
}

.stat-card {
    background: white;
    transition: all 0.3s;
}

.stat-card:hover {
    transform: translateY(-2px);
}

.stat-content {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 8px;
}

.stat-icon {
    flex-shrink: 0;
}

.stat-info {
    flex: 1;
}

.stat-label {
    font-size: 14px;
    color: var(--el-text-color-secondary);
    margin-bottom: 8px;
}

.stat-value {
    font-size: 28px;
    font-weight: 700;
    color: var(--el-text-color-primary);
    line-height: 1;
}

/* 图表区域 */
.charts-section {
    display: grid;
    grid-template-columns: 2fr 1fr;
    gap: 16px;
    margin-bottom: 24px;
}

.chart-card {
    background: white;
}

.chart-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.chart-title {
    font-size: 16px;
    font-weight: 600;
    color: var(--el-text-color-primary);
}

/* 折线图 */
.chart-legend {
    display: flex;
    gap: 24px;
    padding: 0 20px 16px;
    font-size: 14px;
}

.legend-item {
    display: flex;
    align-items: center;
    gap: 8px;
    color: var(--el-text-color-regular);
}

.legend-dot {
    width: 12px;
    height: 12px;
    border-radius: 50%;
}

.legend-dot.total {
    background: #409eff;
}

.legend-dot.success {
    background: #67c23a;
}

.legend-dot.fail {
    background: #f56c6c;
}

.line-chart-container {
    padding: 0 20px 20px;
}

.line-chart {
    width: 100%;
    height: auto;
}

.axis-label {
    font-size: 12px;
    fill: var(--el-text-color-secondary);
}

.date-label {
    font-size: 11px;
    text-anchor: middle;
}

.line-path {
    stroke-width: 2.5;
    stroke-linecap: round;
    stroke-linejoin: round;
}

.total-line {
    stroke: #409eff;
}

.success-line {
    stroke: #67c23a;
}

.fail-line {
    stroke: #f56c6c;
}

/* 环形图 */
.donut-chart-container {
    padding: 20px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 24px;
}

.donut-chart {
    width: 100%;
    max-width: 200px;
    height: auto;
}

.donut-slice {
    transition: all 0.3s;
    cursor: pointer;
}

.donut-slice:hover {
    opacity: 0.8;
}

.success-slice {
    fill: #67c23a;
}

.fail-slice {
    fill: #f56c6c;
}

.donut-center {
    fill: white;
}

.donut-percent-label {
    font-size: 14px;
    fill: var(--el-text-color-secondary);
    text-anchor: middle;
}

.donut-percent-value {
    font-size: 24px;
    font-weight: 700;
    fill: var(--el-text-color-primary);
    text-anchor: middle;
}

.donut-stats {
    display: flex;
    gap: 32px;
    width: 100%;
    justify-content: center;
}

.donut-stat-item {
    text-align: center;
}

.donut-stat-label {
    font-size: 12px;
    color: var(--el-text-color-secondary);
    margin-bottom: 4px;
}

.donut-stat-value {
    font-size: 20px;
    font-weight: 700;
}

.success-text {
    color: #67c23a;
}

.fail-text {
    color: #f56c6c;
}

/* 预约列表 */
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
    background: white;
}

:deep(.el-card__header) {
    padding: 18px 20px;
    border-bottom: 1px solid var(--el-border-color-light);
}

:deep(.el-card__body) {
    padding: 20px;
}

:deep(.el-switch__label) {
    font-size: 14px;
}

/* 响应式 */
@media (max-width: 1024px) {
    .charts-section {
        grid-template-columns: 1fr;
    }
}

@media (max-width: 768px) {
    .stats-cards {
        grid-template-columns: 1fr;
    }

    .header-actions {
        flex-direction: column;
        gap: 8px;
        align-items: flex-end;
    }
}
</style>
