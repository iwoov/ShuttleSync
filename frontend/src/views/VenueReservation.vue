<template>
    <Layout>
        <el-main class="content">
            <el-card class="main-card" shadow="hover">
                <template #header>
                    <div class="card-header">
                        <div class="header-title">
                            <el-icon><OfficeBuilding /></el-icon>
                            <span>普通预约</span>
                        </div>
                        <div class="button-group">
                            <el-button
                                type="primary"
                                :icon="Plus"
                                @click="showCreateDialog = true"
                            >
                                创建预约
                            </el-button>
                            <el-button
                                :icon="Refresh"
                                @click="fetchTasks"
                                :loading="loading"
                            >
                                刷新
                            </el-button>
                        </div>
                    </div>
                </template>

                <!-- 预约任务列表 -->
                <el-table :data="tasks" style="width: 100%" v-loading="loading">
                    <el-table-column prop="task_id" label="任务ID" width="200">
                        <template #default="{ row }">
                            <el-text truncated>{{ row.TaskID }}</el-text>
                        </template>
                    </el-table-column>

                    <el-table-column
                        prop="username"
                        label="预约账号"
                        width="150"
                    />

                    <el-table-column label="场馆/场地" width="150">
                        <template #default="{ row }">
                            <div>{{ getVenueName(row.VenueSiteID) }}</div>
                            <el-text type="info" size="small">{{
                                row.SiteName
                            }}</el-text>
                        </template>
                    </el-table-column>

                    <el-table-column label="预约时间" width="180">
                        <template #default="{ row }">
                            <div>{{ row.ReservationDate }}</div>
                            <el-text type="info" size="small">{{
                                row.ReservationTime
                            }}</el-text>
                        </template>
                    </el-table-column>

                    <el-table-column label="状态" width="100">
                        <template #default="{ row }">
                            <el-tag
                                :type="getStatusType(row.IsFinished)"
                                size="small"
                            >
                                {{ row.IsFinished ? "已完成" : "进行中" }}
                            </el-tag>
                        </template>
                    </el-table-column>

                    <el-table-column
                        prop="CreateTime"
                        label="创建时间"
                        width="180"
                    >
                        <template #default="{ row }">
                            {{ formatTime(row.CreateTime) }}
                        </template>
                    </el-table-column>

                    <el-table-column label="操作" width="150" fixed="right">
                        <template #default="{ row }">
                            <el-button
                                type="primary"
                                size="small"
                                :icon="View"
                                @click="viewTaskDetail(row)"
                            >
                                详情
                            </el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </el-card>

            <!-- 创建预约对话框 -->
            <el-dialog
                v-model="showCreateDialog"
                title="创建普通预约"
                width="600px"
                :close-on-click-modal="false"
            >
                <el-form
                    :model="newReservation"
                    :rules="rules"
                    ref="reservationFormRef"
                    label-width="120px"
                >
                    <!-- 预约账号选择 -->
                    <el-form-item label="预约账号" prop="accountId">
                        <el-select
                            v-model="newReservation.accountId"
                            placeholder="选择预约账号"
                            style="width: 100%"
                            @change="onAccountChange"
                        >
                            <el-option
                                v-for="account in accounts"
                                :key="account.ID"
                                :label="`${account.Lable} (${account.Username})`"
                                :value="account.Username"
                            />
                        </el-select>
                    </el-form-item>

                    <!-- 账号信息显示 -->
                    <el-form-item label="账号信息" v-if="selectedAccount">
                        <div class="account-info-display">
                            <el-input
                                :model-value="selectedAccount.Username"
                                readonly
                                style="width: 150px"
                            />
                            <el-input
                                :model-value="selectedAccount.Password"
                                :type="showPassword ? 'text' : 'password'"
                                readonly
                                style="width: 150px"
                            >
                                <template #suffix>
                                    <el-icon
                                        @click="togglePassword"
                                        style="cursor: pointer"
                                    >
                                        <component
                                            :is="showPassword ? 'Hide' : 'View'"
                                        />
                                    </el-icon>
                                </template>
                            </el-input>
                            <el-button
                                size="small"
                                :type="
                                    getLoginButtonType(
                                        selectedAccount.loginStatus,
                                    )
                                "
                                :icon="
                                    getLoginButtonIcon(
                                        selectedAccount.loginStatus,
                                    )
                                "
                                @click="testLogin(false)"
                            >
                                {{
                                    getLoginButtonText(
                                        selectedAccount.loginStatus,
                                    )
                                }}
                            </el-button>
                        </div>
                    </el-form-item>

                    <!-- 同伴账号选择 -->
                    <el-form-item label="同伴账号" prop="partnerId">
                        <el-select
                            v-model="newReservation.partnerId"
                            placeholder="选择同伴账号"
                            style="width: 100%"
                            @change="onPartnerChange"
                        >
                            <el-option
                                v-for="partner in filteredPartners"
                                :key="partner.ID"
                                :label="`${partner.Lable} (${partner.Username})`"
                                :value="partner.Username"
                            />
                        </el-select>
                    </el-form-item>

                    <!-- 同伴账号信息显示 -->
                    <el-form-item label="同伴信息" v-if="selectedPartner">
                        <div class="account-info-display">
                            <el-input
                                :model-value="selectedPartner.Username"
                                readonly
                                style="width: 150px"
                            />
                            <el-input
                                :model-value="selectedPartner.Password"
                                :type="
                                    showPartnerPassword ? 'text' : 'password'
                                "
                                readonly
                                style="width: 150px"
                            >
                                <template #suffix>
                                    <el-icon
                                        @click="togglePartnerPassword"
                                        style="cursor: pointer"
                                    >
                                        <component
                                            :is="
                                                showPartnerPassword
                                                    ? 'Hide'
                                                    : 'View'
                                            "
                                        />
                                    </el-icon>
                                </template>
                            </el-input>
                            <el-button
                                size="small"
                                :type="
                                    getLoginButtonType(
                                        selectedPartner.loginStatus,
                                    )
                                "
                                :icon="
                                    getLoginButtonIcon(
                                        selectedPartner.loginStatus,
                                    )
                                "
                                @click="testLogin(true)"
                            >
                                {{
                                    getLoginButtonText(
                                        selectedPartner.loginStatus,
                                    )
                                }}
                            </el-button>
                        </div>
                    </el-form-item>

                    <!-- 场馆选择 -->
                    <el-form-item label="场馆" prop="venueId">
                        <el-select
                            v-model="newReservation.venueId"
                            placeholder="选择场馆"
                            style="width: 100%"
                            @change="onVenueChange"
                        >
                            <el-option label="体育馆" :value="1" />
                            <el-option label="风雨操场" :value="2" />
                        </el-select>
                    </el-form-item>

                    <!-- 场地选择 -->
                    <el-form-item label="场地号" prop="locationId">
                        <el-select
                            v-model="newReservation.locationId"
                            placeholder="选择场地"
                            style="width: 100%"
                        >
                            <el-option
                                v-for="location in availableLocations"
                                :key="location.value"
                                :label="location.label"
                                :value="location.value"
                            />
                        </el-select>
                    </el-form-item>

                    <!-- 预约日期 -->
                    <el-form-item label="预约日期" prop="date">
                        <el-date-picker
                            v-model="newReservation.date"
                            type="date"
                            placeholder="选择日期"
                            format="YYYY-MM-DD"
                            value-format="YYYY-MM-DD"
                            :disabled-date="disabledDate"
                            style="width: 100%"
                        />
                    </el-form-item>

                    <!-- 时间段 -->
                    <el-form-item label="时间段" prop="timeSlot">
                        <el-select
                            v-model="newReservation.timeSlot"
                            placeholder="选择时间段"
                            style="width: 100%"
                        >
                            <el-option
                                v-for="slot in timeSlots"
                                :key="slot.value"
                                :label="slot.label"
                                :value="slot.value"
                            />
                        </el-select>
                    </el-form-item>
                </el-form>

                <template #footer>
                    <el-button @click="showCreateDialog = false"
                        >取消</el-button
                    >
                    <el-button
                        type="primary"
                        @click="handleCreateReservation"
                        :loading="creating"
                        :disabled="!canSubmit"
                    >
                        创建预约
                    </el-button>
                </template>
            </el-dialog>

            <!-- 任务详情对话框 -->
            <el-dialog
                v-model="showDetailDialog"
                title="预约详情"
                width="700px"
            >
                <el-descriptions :column="2" border v-if="currentTask">
                    <el-descriptions-item label="任务ID" :span="2">
                        {{ currentTask.TaskID }}
                    </el-descriptions-item>
                    <el-descriptions-item label="预约账号">
                        {{ currentTask.Username }}
                    </el-descriptions-item>
                    <el-descriptions-item label="状态">
                        <el-tag :type="getStatusType(currentTask.IsFinished)">
                            {{ currentTask.IsFinished ? "已完成" : "进行中" }}
                        </el-tag>
                    </el-descriptions-item>
                    <el-descriptions-item label="场馆">
                        {{ getVenueName(currentTask.VenueSiteID) }}
                    </el-descriptions-item>
                    <el-descriptions-item label="场地">
                        {{ currentTask.SiteName }}
                    </el-descriptions-item>
                    <el-descriptions-item label="预约日期">
                        {{ currentTask.ReservationDate }}
                    </el-descriptions-item>
                    <el-descriptions-item label="预约时间">
                        {{ currentTask.ReservationTime }}
                    </el-descriptions-item>
                    <el-descriptions-item label="同伴ID">
                        {{ currentTask.BuddyUserID }}
                    </el-descriptions-item>
                    <el-descriptions-item label="同伴码">
                        {{ currentTask.BuddyNum }}
                    </el-descriptions-item>
                    <el-descriptions-item label="创建时间" :span="2">
                        {{ formatTime(currentTask.CreateTime) }}
                    </el-descriptions-item>
                </el-descriptions>
            </el-dialog>
        </el-main>
    </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { ElMessage } from "element-plus";
import {
    OfficeBuilding,
    Plus,
    Refresh,
    View,
    Select,
    SuccessFilled,
    CircleClose,
} from "@element-plus/icons-vue";
import Layout from "@/components/Layout.vue";
import { get, post } from "@/utils/api";
import { getUsername, getCaptchaApi } from "@/utils/auth";

// 状态
const loading = ref(false);
const creating = ref(false);
const showCreateDialog = ref(false);
const showDetailDialog = ref(false);
const showPassword = ref(false);
const showPartnerPassword = ref(false);

const tasks = ref<any[]>([]);
const accounts = ref<any[]>([]);
const currentTask = ref<any>(null);
const reservationFormRef = ref();
const apiKey = ref(getCaptchaApi() || "");

// 新预约表单
const newReservation = ref({
    accountId: "",
    partnerId: "",
    venueId: null as number | null,
    locationId: "",
    date: "",
    timeSlot: "",
    accountInfo: null as any,
    partnerInfo: null as any,
});

// 场地和时间段选项
const availableLocations = ref<any[]>([]);
const timeSlots = ref<any[]>([]);

// 表单验证规则
const rules = {
    accountId: [
        { required: true, message: "请选择预约账号", trigger: "change" },
    ],
    partnerId: [
        { required: true, message: "请选择同伴账号", trigger: "change" },
    ],
    venueId: [{ required: true, message: "请选择场馆", trigger: "change" }],
    locationId: [{ required: true, message: "请选择场地", trigger: "change" }],
    date: [{ required: true, message: "请选择预约日期", trigger: "change" }],
    timeSlot: [{ required: true, message: "请选择时间段", trigger: "change" }],
};

// 计算属性
const selectedAccount = computed(() => {
    return accounts.value.find(
        (acc) => acc.Username === newReservation.value.accountId,
    );
});

const selectedPartner = computed(() => {
    return accounts.value.find(
        (acc) => acc.Username === newReservation.value.partnerId,
    );
});

const filteredPartners = computed(() => {
    if (!newReservation.value.accountId) return accounts.value;
    return accounts.value.filter(
        (acc) => acc.Username !== newReservation.value.accountId,
    );
});

const canSubmit = computed(() => {
    return (
        selectedAccount.value?.loginStatus === "success" &&
        selectedPartner.value?.loginStatus === "success" &&
        apiKey.value
    );
});

// 禁用过去的日期
const disabledDate = (time: Date) => {
    return time.getTime() < Date.now() - 8.64e7;
};

// 切换密码显示
const togglePassword = () => {
    showPassword.value = !showPassword.value;
};

const togglePartnerPassword = () => {
    showPartnerPassword.value = !showPartnerPassword.value;
};

// 场馆切换
const onVenueChange = () => {
    const venueId = newReservation.value.venueId;
    if (!venueId) {
        availableLocations.value = [];
        newReservation.value.locationId = "";
        return;
    }

    const locationCount = venueId === 1 ? 12 : 20;
    availableLocations.value = Array.from(
        { length: locationCount },
        (_, i) => ({
            value: `${i + 1}号`,
            label: `${i + 1}号场地`,
        }),
    );

    newReservation.value.locationId = "";
};

// 账号切换
const onAccountChange = () => {
    if (newReservation.value.partnerId === newReservation.value.accountId) {
        newReservation.value.partnerId = "";
    }
};

const onPartnerChange = () => {
    // 可以添加额外的逻辑
};

// 初始化时间段
const initTimeSlots = () => {
    const slots = [];
    let hour = 8;
    let minute = 30;

    while (hour < 22 || (hour === 22 && minute === 30)) {
        const startTime = `${hour.toString().padStart(2, "0")}:${minute
            .toString()
            .padStart(2, "0")}`;
        let endHour = hour + 1;
        let endMinute = minute;
        const endTime = `${endHour.toString().padStart(2, "0")}:${endMinute
            .toString()
            .padStart(2, "0")}`;

        slots.push({
            value: startTime,
            label: `${startTime}-${endTime}`,
        });

        hour += 1;
    }

    timeSlots.value = slots;
};

// 获取账号列表
const fetchAccounts = async () => {
    try {
        const username = getUsername();
        if (!username) {
            ElMessage.error("获取账号列表失败，请重新登录");
            return;
        }
        const data = await get("/account/list", { user: username });
        if (data.message === "success") {
            accounts.value = data.data.map((account: any) => ({
                ...account,
                loginStatus: null,
            }));
        } else {
            throw new Error(data.message || "获取账号列表失败");
        }
    } catch (error: any) {
        console.error("获取账号列表失败", error);
        ElMessage.error(error.message || "获取账号列表失败");
    }
};

// 获取任务列表
const fetchTasks = async () => {
    loading.value = true;
    try {
        const username = getUsername();
        if (!username) {
            ElMessage.error("请先登录");
            return;
        }
        const data = await get("/task/list", { user: username });
        if (data.message === "success") {
            tasks.value = data.data || [];
        } else {
            throw new Error(data.message || "获取任务列表失败");
        }
    } catch (error: any) {
        console.error("获取任务列表失败", error);
        ElMessage.error(error.message || "获取任务列表失败");
    } finally {
        loading.value = false;
    }
};

// 测试登录
const testLogin = async (isPartner: boolean) => {
    const account = isPartner ? selectedPartner.value : selectedAccount.value;
    if (!account) return;

    try {
        const url = isPartner ? "/tyys/buddy_num" : "/tyys/login";

        const data = await post(url, {
            username: account.Username,
            password: account.Password,
        });

        if (data.message === "success") {
            account.loginStatus = "success";

            if (isPartner) {
                newReservation.value.partnerInfo = {
                    buddy_id: data.data.buddy_id,
                    buddy_num: data.data.buddy_num,
                };
            } else {
                newReservation.value.accountInfo = {
                    username: data.data.username,
                    password: data.data.password,
                    phone: data.data.phone,
                };
            }
            ElMessage.success("登录测试成功");
        } else {
            account.loginStatus = "error";
            ElMessage.error("登录测试失败");
        }
    } catch (error) {
        account.loginStatus = "error";
        ElMessage.error("登录测试失败");
    }
};

// 创建预约
const handleCreateReservation = async () => {
    if (!reservationFormRef.value) return;

    await reservationFormRef.value.validate(async (valid: boolean) => {
        if (!valid) return;

        if (
            !newReservation.value.accountInfo ||
            !newReservation.value.partnerInfo
        ) {
            ElMessage.warning("请先完成账号登录测试");
            return;
        }

        creating.value = true;
        try {
            const venue_site_id =
                newReservation.value.venueId === 1 ? "143" : "23";
            const currentUser = getUsername();
            if (!currentUser) {
                ElMessage.error("请先登录后再提交预约");
                return;
            }

            const requestData = {
                user: currentUser,
                username: newReservation.value.accountInfo.username,
                password: newReservation.value.accountInfo.password,
                user_phone: newReservation.value.accountInfo.phone,
                captcha_api: apiKey.value,
                buddy_user_id: newReservation.value.partnerInfo.buddy_id,
                buddy_num: newReservation.value.partnerInfo.buddy_num,
                venue_site_id: venue_site_id,
                reservation_date: newReservation.value.date,
                reservation_time: newReservation.value.timeSlot,
                site_name: newReservation.value.locationId,
            };

            const data = await post("/task/add", requestData);
            if (data.message === "success") {
                ElMessage.success("预约创建成功");
                showCreateDialog.value = false;
                resetForm();
                fetchTasks();
            } else {
                throw new Error(data.data || "预约创建失败");
            }
        } catch (error: any) {
            ElMessage.error(error.message || "预约创建失败");
        } finally {
            creating.value = false;
        }
    });
};

// 查看详情
const viewTaskDetail = (task: any) => {
    currentTask.value = task;
    showDetailDialog.value = true;
};

// 重置表单
const resetForm = () => {
    newReservation.value = {
        accountId: "",
        partnerId: "",
        venueId: null,
        locationId: "",
        date: "",
        timeSlot: "",
        accountInfo: null,
        partnerInfo: null,
    };
    availableLocations.value = [];

    // 重置账号登录状态
    accounts.value.forEach((acc) => {
        acc.loginStatus = null;
    });
};

// 获取登录按钮文本
const getLoginButtonText = (status: string) => {
    switch (status) {
        case "success":
            return "登录成功";
        case "error":
            return "登录失败";
        default:
            return "测试登录";
    }
};

const getLoginButtonType = (status: string) => {
    switch (status) {
        case "success":
            return "success";
        case "error":
            return "danger";
        default:
            return "primary";
    }
};

const getLoginButtonIcon = (status: string) => {
    switch (status) {
        case "success":
            return SuccessFilled;
        case "error":
            return CircleClose;
        default:
            return Select;
    }
};

// 获取场馆名称
const getVenueName = (venueId: string) => {
    return venueId === "143" ? "体育馆" : "风雨操场";
};

// 获取状态类型
const getStatusType = (isFinished: boolean) => {
    return isFinished ? "success" : "warning";
};

// 格式化时间
const formatTime = (time: string) => {
    if (!time) return "";
    return new Date(time).toLocaleString("zh-CN");
};

// 页面加载
onMounted(() => {
    fetchAccounts();
    fetchTasks();
    initTimeSlots();
});
</script>

<style scoped>
.content {
    padding: 20px;
    background-color: #f5f7fa;
    min-height: calc(100vh - 80px);
}

.main-card {
    margin-bottom: 20px;
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
    font-size: 18px;
    font-weight: 600;
}

.button-group {
    display: flex;
    gap: 10px;
}

.account-info-display {
    display: flex;
    gap: 8px;
    align-items: center;
    width: 100%;
    flex-wrap: wrap;
}

:deep(.el-table) {
    font-size: 14px;
}

:deep(.el-descriptions__label) {
    font-weight: 600;
}
</style>
