<template>
    <Layout>
        <el-main class="content">
            <!-- 页面标题和操作按钮 -->
            <el-card class="main-card" shadow="hover">
                <template #header>
                    <div class="card-header">
                        <div class="header-title">
                            <el-icon><Search /></el-icon>
                            <span>捡漏模式</span>
                        </div>
                        <div class="button-group">
                            <el-button
                                type="primary"
                                :icon="Plus"
                                @click="showCreateDialog = true"
                            >
                                创建捡漏任务
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

                <!-- 任务列表 -->
                <el-table :data="tasks" style="width: 100%" v-loading="loading">
                    <el-table-column prop="task_id" label="任务ID" width="280">
                        <template #default="{ row }">
                            <el-text truncated>{{ row.task_id }}</el-text>
                        </template>
                    </el-table-column>

                    <el-table-column
                        prop="reservation_date"
                        label="预约日期"
                        width="120"
                    />

                    <el-table-column label="场地/时间" width="150">
                        <template #default="{ row }">
                            <div>{{ row.site_name || "任意" }}</div>
                            <el-text type="info" size="small">{{
                                row.reservation_time || "任意"
                            }}</el-text>
                        </template>
                    </el-table-column>

                    <el-table-column
                        prop="scan_interval"
                        label="扫描间隔"
                        width="100"
                    >
                        <template #default="{ row }">
                            {{ row.scan_interval }} 分钟
                        </template>
                    </el-table-column>

                    <el-table-column label="状态" width="100">
                        <template #default="{ row }">
                            <el-tag
                                :type="getStatusType(row.status)"
                                size="small"
                            >
                                {{ getStatusText(row.status) }}
                            </el-tag>
                        </template>
                    </el-table-column>

                    <el-table-column label="统计" width="120">
                        <template #default="{ row }">
                            <div>扫描: {{ row.scan_count }}</div>
                            <div>成功: {{ row.success_count }}</div>
                        </template>
                    </el-table-column>

                    <el-table-column label="截止时间" width="160">
                        <template #default="{ row }">
                            {{ row.deadline ? formatTime(row.deadline) : "无" }}
                        </template>
                    </el-table-column>

                    <el-table-column
                        prop="created_at"
                        label="创建时间"
                        width="160"
                    >
                        <template #default="{ row }">
                            {{ formatTime(row.created_at) }}
                        </template>
                    </el-table-column>

                    <el-table-column label="操作" width="260" fixed="right">
                        <template #default="{ row }">
                            <el-button
                                type="primary"
                                size="small"
                                :icon="View"
                                @click="viewTaskDetail(row)"
                            >
                                详情
                            </el-button>
                            <el-button
                                v-if="row.status === 'active'"
                                type="warning"
                                size="small"
                                :icon="Edit"
                                @click="handleEditTask(row)"
                            >
                                编辑
                            </el-button>
                            <el-button
                                v-if="row.status === 'active'"
                                type="danger"
                                size="small"
                                :icon="Close"
                                @click="handleCancelTask(row)"
                            >
                                取消
                            </el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </el-card>

            <!-- 创建任务对话框 -->
            <el-dialog
                v-model="showCreateDialog"
                title="创建捡漏任务"
                width="600px"
                :close-on-click-modal="false"
            >
                <el-form
                    :model="newTask"
                    :rules="rules"
                    ref="taskFormRef"
                    label-width="120px"
                >
                    <el-form-item label="主预约账号" prop="account_id_1">
                        <el-select
                            v-model="newTask.account_id_1"
                            placeholder="选择主预约账号"
                            style="width: 100%"
                            @change="onAccount1Change"
                        >
                            <el-option
                                v-for="account in accounts"
                                :key="account.ID"
                                :label="`${account.Lable} (${account.Username})`"
                                :value="account.ID"
                            />
                        </el-select>
                        <el-text type="info" size="small"
                            >账号1作为主预约账号</el-text
                        >
                    </el-form-item>

                    <!-- 主账号信息显示 -->
                    <el-form-item label="主账号信息" v-if="selectedAccount1">
                        <div class="account-info-display">
                            <el-input
                                :model-value="selectedAccount1.Username"
                                readonly
                                style="width: 150px"
                            />
                            <el-input
                                :model-value="selectedAccount1.Password"
                                :type="showPassword1 ? 'text' : 'password'"
                                readonly
                                style="width: 150px"
                            >
                                <template #suffix>
                                    <el-icon
                                        @click="showPassword1 = !showPassword1"
                                        style="cursor: pointer"
                                    >
                                        <component
                                            :is="
                                                showPassword1 ? 'Hide' : 'View'
                                            "
                                        />
                                    </el-icon>
                                </template>
                            </el-input>
                            <el-button
                                size="small"
                                :type="
                                    getLoginButtonType(
                                        selectedAccount1.loginStatus,
                                    )
                                "
                                :icon="
                                    getLoginButtonIcon(
                                        selectedAccount1.loginStatus,
                                    )
                                "
                                @click="testAccount1Login"
                            >
                                {{
                                    getLoginButtonText(
                                        selectedAccount1.loginStatus,
                                    )
                                }}
                            </el-button>
                        </div>
                    </el-form-item>

                    <el-form-item label="同伴账号" prop="account_id_2">
                        <el-select
                            v-model="newTask.account_id_2"
                            placeholder="选择同伴账号"
                            style="width: 100%"
                            @change="onAccount2Change"
                        >
                            <el-option
                                v-for="account in filteredSecondAccounts"
                                :key="account.ID"
                                :label="`${account.Lable} (${account.Username})`"
                                :value="account.ID"
                            />
                        </el-select>
                        <el-text type="info" size="small"
                            >账号2提供同伴码</el-text
                        >
                    </el-form-item>

                    <!-- 同伴账号信息显示 -->
                    <el-form-item label="同伴信息" v-if="selectedAccount2">
                        <div class="account-info-display">
                            <el-input
                                :model-value="selectedAccount2.Username"
                                readonly
                                style="width: 150px"
                            />
                            <el-input
                                :model-value="selectedAccount2.Password"
                                :type="showPassword2 ? 'text' : 'password'"
                                readonly
                                style="width: 150px"
                            >
                                <template #suffix>
                                    <el-icon
                                        @click="showPassword2 = !showPassword2"
                                        style="cursor: pointer"
                                    >
                                        <component
                                            :is="
                                                showPassword2 ? 'Hide' : 'View'
                                            "
                                        />
                                    </el-icon>
                                </template>
                            </el-input>
                            <el-button
                                size="small"
                                :type="
                                    getLoginButtonType(
                                        selectedAccount2.loginStatus,
                                    )
                                "
                                :icon="
                                    getLoginButtonIcon(
                                        selectedAccount2.loginStatus,
                                    )
                                "
                                @click="testAccount2Login"
                            >
                                {{
                                    getLoginButtonText(
                                        selectedAccount2.loginStatus,
                                    )
                                }}
                            </el-button>
                        </div>
                    </el-form-item>

                    <el-form-item label="场馆" prop="venue_id">
                        <el-select
                            v-model="newTask.venue_id"
                            placeholder="选择场馆"
                            style="width: 100%"
                            @change="onVenueChange"
                        >
                            <el-option label="体育馆" :value="1" />
                            <el-option label="风雨操场" :value="2" />
                        </el-select>
                    </el-form-item>

                    <el-form-item label="预约日期" prop="reservation_date">
                        <el-date-picker
                            v-model="newTask.reservation_date"
                            type="date"
                            placeholder="选择预约日期"
                            format="YYYY-MM-DD"
                            value-format="YYYY-MM-DD"
                            :disabledDate="disabledDate"
                            style="width: 100%"
                        />
                    </el-form-item>

                    <el-form-item label="场地号">
                        <el-select
                            v-model="newTask.site_name"
                            placeholder="可选，留空表示任意场地"
                            style="width: 100%"
                            clearable
                        >
                            <el-option
                                v-for="location in availableLocations"
                                :key="location.value"
                                :label="location.label"
                                :value="location.value"
                            />
                        </el-select>
                        <el-text type="info" size="small"
                            >留空则随机选择可用场地</el-text
                        >
                    </el-form-item>

                    <el-form-item label="时间段">
                        <el-select
                            v-model="newTask.reservation_time"
                            placeholder="可选，留空表示任意时间"
                            style="width: 100%"
                            clearable
                        >
                            <el-option
                                v-for="slot in timeSlots"
                                :key="slot.value"
                                :label="slot.label"
                                :value="slot.value"
                            />
                        </el-select>
                        <el-text type="info" size="small"
                            >留空则随机选择可用时间</el-text
                        >
                    </el-form-item>

                    <el-form-item label="扫描间隔" prop="scan_interval">
                        <el-input-number
                            v-model="newTask.scan_interval"
                            :min="1"
                            :max="60"
                            placeholder="分钟"
                        />
                        <el-text type="info" size="small"
                            >1-60分钟，建议5-10分钟</el-text
                        >
                    </el-form-item>

                    <el-form-item label="预约截止时间">
                        <el-date-picker
                            v-model="newTask.deadline"
                            type="datetime"
                            placeholder="可选，超过此时间未预约到则失败"
                            format="YYYY-MM-DD HH:mm:ss"
                            value-format="YYYY-MM-DD HH:mm:ss"
                            :disabledDate="disabledDeadlineDate"
                            :disabledHours="disabledDeadlineHours"
                            style="width: 100%"
                        />
                        <el-text type="info" size="small"
                            >留空则在预约日期结束前一直尝试</el-text
                        >
                    </el-form-item>
                </el-form>

                <template #footer>
                    <div style="width: 100%">
                        <el-text
                            v-if="!canCreateTask"
                            type="warning"
                            size="small"
                            style="display: block; margin-bottom: 12px"
                        >
                            请确保两个账号都完成登录测试后再创建任务
                        </el-text>
                        <div
                            style="
                                display: flex;
                                justify-content: flex-end;
                                gap: 10px;
                            "
                        >
                            <el-button @click="showCreateDialog = false">
                                取消
                            </el-button>
                            <el-button
                                type="primary"
                                @click="handleCreateTask"
                                :loading="creating"
                                :disabled="!canCreateTask"
                            >
                                创建
                            </el-button>
                        </div>
                    </div>
                </template>
            </el-dialog>

            <!-- 编辑任务对话框 -->
            <el-dialog
                v-model="showEditDialog"
                title="编辑捡漏任务"
                width="600px"
                :close-on-click-modal="false"
            >
                <el-form
                    :model="editingTask"
                    :rules="rules"
                    ref="editFormRef"
                    label-width="120px"
                    v-if="editingTask"
                >
                    <el-form-item label="主预约账号" prop="account_id_1">
                        <el-select
                            v-model="editingTask.account_id_1"
                            placeholder="选择主预约账号"
                            style="width: 100%"
                        >
                            <el-option
                                v-for="account in accounts"
                                :key="account.ID"
                                :label="`${account.Lable} (${account.Username})`"
                                :value="account.ID"
                            />
                        </el-select>
                    </el-form-item>

                    <el-form-item label="同伴账号" prop="account_id_2">
                        <el-select
                            v-model="editingTask.account_id_2"
                            placeholder="选择同伴账号"
                            style="width: 100%"
                        >
                            <el-option
                                v-for="account in accounts.filter(
                                    (acc) =>
                                        acc.ID !== editingTask.account_id_1,
                                )"
                                :key="account.ID"
                                :label="`${account.Lable} (${account.Username})`"
                                :value="account.ID"
                            />
                        </el-select>
                    </el-form-item>

                    <el-form-item label="场馆" prop="venue_id">
                        <el-select
                            v-model="editingTask.venue_id"
                            placeholder="选择场馆"
                            style="width: 100%"
                            @change="onEditVenueChange"
                        >
                            <el-option label="体育馆" :value="1" />
                            <el-option label="风雨操场" :value="2" />
                        </el-select>
                    </el-form-item>

                    <el-form-item label="预约日期" prop="reservation_date">
                        <el-date-picker
                            v-model="editingTask.reservation_date"
                            type="date"
                            placeholder="选择预约日期"
                            format="YYYY-MM-DD"
                            value-format="YYYY-MM-DD"
                            :disabledDate="disabledDate"
                            style="width: 100%"
                        />
                    </el-form-item>

                    <el-form-item label="场地号">
                        <el-select
                            v-model="editingTask.site_name"
                            placeholder="可选，留空表示任意场地"
                            style="width: 100%"
                            clearable
                        >
                            <el-option
                                v-for="location in editAvailableLocations"
                                :key="location.value"
                                :label="location.label"
                                :value="location.value"
                            />
                        </el-select>
                    </el-form-item>

                    <el-form-item label="时间段">
                        <el-select
                            v-model="editingTask.reservation_time"
                            placeholder="可选，留空表示任意时间"
                            style="width: 100%"
                            clearable
                        >
                            <el-option
                                v-for="slot in timeSlots"
                                :key="slot.value"
                                :label="slot.label"
                                :value="slot.value"
                            />
                        </el-select>
                    </el-form-item>

                    <el-form-item label="扫描间隔" prop="scan_interval">
                        <el-input-number
                            v-model="editingTask.scan_interval"
                            :min="1"
                            :max="60"
                            placeholder="分钟"
                        />
                    </el-form-item>

                    <el-form-item label="预约截止时间">
                        <el-date-picker
                            v-model="editingTask.deadline"
                            type="datetime"
                            placeholder="可选"
                            format="YYYY-MM-DD HH:mm:ss"
                            value-format="YYYY-MM-DD HH:mm:ss"
                            :disabledDate="disabledDeadlineDate"
                            style="width: 100%"
                        />
                    </el-form-item>
                </el-form>

                <template #footer>
                    <el-button @click="showEditDialog = false">取消</el-button>
                    <el-button
                        type="primary"
                        @click="handleUpdateTask"
                        :loading="updating"
                    >
                        更新
                    </el-button>
                </template>
            </el-dialog>

            <!-- 任务详情对话框 -->
            <el-dialog
                v-model="showDetailDialog"
                title="任务详情"
                width="800px"
            >
                <el-descriptions :column="2" border v-if="currentTask">
                    <el-descriptions-item label="任务ID" :span="2">
                        {{ currentTask.task_id }}
                    </el-descriptions-item>
                    <el-descriptions-item label="状态">
                        <el-tag :type="getStatusType(currentTask.status)">
                            {{ getStatusText(currentTask.status) }}
                        </el-tag>
                    </el-descriptions-item>
                    <el-descriptions-item label="预约日期">
                        {{ currentTask.reservation_date }}
                    </el-descriptions-item>
                    <el-descriptions-item label="场地号">
                        {{ currentTask.site_name || "任意" }}
                    </el-descriptions-item>
                    <el-descriptions-item label="时间段">
                        {{ currentTask.reservation_time || "任意" }}
                    </el-descriptions-item>
                    <el-descriptions-item label="扫描间隔">
                        {{ currentTask.scan_interval }} 分钟
                    </el-descriptions-item>
                    <el-descriptions-item label="扫描次数">
                        {{ currentTask.scan_count }}
                    </el-descriptions-item>
                    <el-descriptions-item label="成功次数">
                        {{ currentTask.success_count }}
                    </el-descriptions-item>
                    <el-descriptions-item label="最后扫描">
                        {{
                            currentTask.last_scan_time
                                ? formatTime(currentTask.last_scan_time)
                                : "未扫描"
                        }}
                    </el-descriptions-item>
                    <el-descriptions-item label="截止时间">
                        {{
                            currentTask.deadline
                                ? formatTime(currentTask.deadline)
                                : "无"
                        }}
                    </el-descriptions-item>
                    <el-descriptions-item label="创建时间">
                        {{ formatTime(currentTask.created_at) }}
                    </el-descriptions-item>
                    <el-descriptions-item
                        v-if="currentTask.status === 'failed'"
                        label="失败原因"
                        :span="2"
                    >
                        <el-text type="danger">{{
                            currentTask.failure_reason
                        }}</el-text>
                    </el-descriptions-item>
                </el-descriptions>

                <!-- 扫描日志 -->
                <div style="margin-top: 20px">
                    <div
                        style="
                            display: flex;
                            justify-content: space-between;
                            align-items: center;
                            margin-bottom: 10px;
                        "
                    >
                        <h3>扫描日志</h3>
                        <el-button
                            size="small"
                            :icon="Refresh"
                            @click="fetchLogs"
                            >刷新</el-button
                        >
                    </div>

                    <el-timeline v-if="logs.length > 0" v-loading="logsLoading">
                        <el-timeline-item
                            v-for="log in logs"
                            :key="log.id"
                            :timestamp="formatTime(log.scan_time)"
                            placement="top"
                        >
                            <el-card>
                                <div
                                    style="
                                        display: flex;
                                        justify-content: space-between;
                                        align-items: center;
                                    "
                                >
                                    <div>
                                        <el-tag
                                            v-if="log.success"
                                            type="success"
                                            size="small"
                                            >预约成功</el-tag
                                        >
                                        <el-tag v-else type="info" size="small"
                                            >扫描</el-tag
                                        >
                                        <span style="margin-left: 10px">{{
                                            log.message
                                        }}</span>
                                    </div>
                                    <div>
                                        <el-tag
                                            v-if="log.available_slots > 0"
                                            type="success"
                                            size="small"
                                        >
                                            发现
                                            {{ log.available_slots }} 个场地
                                        </el-tag>
                                    </div>
                                </div>
                                <div
                                    v-if="log.details"
                                    style="
                                        margin-top: 10px;
                                        color: #606266;
                                        font-size: 12px;
                                    "
                                >
                                    {{ log.details }}
                                </div>
                            </el-card>
                        </el-timeline-item>
                    </el-timeline>

                    <el-empty v-else description="暂无日志" />
                </div>
            </el-dialog>
        </el-main>
    </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import {
    Search,
    Plus,
    Refresh,
    View,
    Close,
    Select,
    SuccessFilled,
    CircleClose,
    Edit,
} from "@element-plus/icons-vue";
import Layout from "@/components/Layout.vue";
import { get, post } from "@/utils/api";
import {
    createBargainTask,
    getBargainTasks,
    updateBargainTask,
    cancelBargainTask,
    getBargainTaskLogs,
} from "@/utils/api";

// 状态
const loading = ref(false);
const creating = ref(false);
const updating = ref(false);
const logsLoading = ref(false);
const showCreateDialog = ref(false);
const showEditDialog = ref(false);
const showDetailDialog = ref(false);
const showPassword1 = ref(false);
const showPassword2 = ref(false);
const tasks = ref<any[]>([]);
const accounts = ref<any[]>([]);
const logs = ref<any[]>([]);
const currentTask = ref<any>(null);
const editingTask = ref<any>(null);
const taskFormRef = ref();
const editFormRef = ref();

// 新任务表单
const newTask = ref({
    account_id_1: null as number | null,
    account_id_2: null as number | null,
    venue_id: null as number | null,
    venue_site_id: "",
    reservation_date: "",
    site_name: "",
    reservation_time: "",
    scan_interval: 10,
    deadline: "", // 预约截止时间
});

// 场地选项
const availableLocations = ref<any[]>([]);
const editAvailableLocations = ref<any[]>([]);

// 时间段选项（与普通模式一致）
const timeSlots = ref<any[]>([]);

// 表单验证规则
const rules = {
    account_id_1: [
        { required: true, message: "请选择主预约账号", trigger: "change" },
    ],
    account_id_2: [
        { required: true, message: "请选择同伴账号", trigger: "change" },
    ],
    venue_id: [{ required: true, message: "请选择场馆", trigger: "change" }],
    reservation_date: [
        { required: true, message: "请选择预约日期", trigger: "change" },
    ],
    scan_interval: [
        { required: true, message: "请设置扫描间隔", trigger: "blur" },
    ],
};

// 过滤第二个账号（不能和第一个账号相同）
const filteredSecondAccounts = computed(() => {
    if (!newTask.value.account_id_1) return accounts.value;
    return accounts.value.filter(
        (acc) => acc.ID !== newTask.value.account_id_1,
    );
});

// 计算属性：获取选中的账号
const selectedAccount1 = computed(() => {
    return accounts.value.find((acc) => acc.ID === newTask.value.account_id_1);
});

const selectedAccount2 = computed(() => {
    return accounts.value.find((acc) => acc.ID === newTask.value.account_id_2);
});

// 检查是否可以提交（两个账号都登录成功）
const canCreateTask = computed(() => {
    return (
        selectedAccount1.value?.loginStatus === "success" &&
        selectedAccount2.value?.loginStatus === "success"
    );
});

// 当第一个账号改变时，清空第二个账号的选择
const onAccount1Change = () => {
    if (newTask.value.account_id_2 === newTask.value.account_id_1) {
        newTask.value.account_id_2 = null;
    }
};

const onAccount2Change = () => {
    // 账号切换时可以添加额外逻辑
};

// 场馆切换时，加载场地列表并设置venue_site_id
const onVenueChange = () => {
    const venueId = newTask.value.venue_id;
    if (!venueId) {
        availableLocations.value = [];
        newTask.value.venue_site_id = "";
        newTask.value.site_name = "";
        return;
    }

    // 设置venue_site_id：体育馆=143，风雨操场=23
    newTask.value.venue_site_id = venueId === 1 ? "143" : "23";

    // 加载场地列表：体育馆12个场地，风雨操场20个场地
    const locationCount = venueId === 1 ? 12 : 20;
    availableLocations.value = Array.from(
        { length: locationCount },
        (_, i) => ({
            value: `${i + 1}号`,
            label: `${i + 1}号场地`,
        }),
    );

    // 清空已选择的场地
    newTask.value.site_name = "";
};

// 编辑对话框的场馆切换
const onEditVenueChange = () => {
    const venueId = editingTask.value?.venue_id;
    if (!venueId) {
        editAvailableLocations.value = [];
        return;
    }

    const locationCount = venueId === 1 ? 12 : 20;
    editAvailableLocations.value = Array.from(
        { length: locationCount },
        (_, i) => ({
            value: `${i + 1}号`,
            label: `${i + 1}号场地`,
        }),
    );
};

// 初始化时间段选项（与普通模式一致）
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

// 禁用过去的日期
const disabledDate = (time: Date) => {
    return time.getTime() < Date.now() - 8.64e7; // 禁用今天之前的日期
};

// 截止时间日期限制：不能早于当前时间，不能晚于预约日期
const disabledDeadlineDate = (time: Date) => {
    const now = Date.now();
    // 不能早于今天
    if (time.getTime() < now - 8.64e7) {
        return true;
    }
    // 如果已选择预约日期，不能晚于预约日期
    if (newTask.value.reservation_date) {
        const reservationDate = new Date(newTask.value.reservation_date);
        // 截止时间不能晚于预约日期当天
        if (time.getTime() > reservationDate.getTime() + 8.64e7 - 1) {
            return true;
        }
    }
    return false;
};

// 截止时间小时限制：如果是今天，不能早于当前小时
const disabledDeadlineHours = () => {
    const hours: number[] = [];
    const now = new Date();
    const currentHour = now.getHours();

    // 如果没选截止日期，不限制
    if (!newTask.value.deadline) return hours;

    const deadlineDate = new Date(newTask.value.deadline);
    // 如果是今天，禁用已过去的小时
    if (deadlineDate.toDateString() === now.toDateString()) {
        for (let i = 0; i < currentHour; i++) {
            hours.push(i);
        }
    }
    return hours;
};

// 获取账号列表
const fetchAccounts = async () => {
    try {
        const res = await get("/account/list");
        if (res.message === "success") {
            accounts.value = (res.data || []).map((account: any) => ({
                ...account,
                loginStatus: null,
            }));
        }
    } catch (error) {
        console.error("获取账号列表失败:", error);
        ElMessage.error("获取账号列表失败");
    }
};

// 测试主账号登录
const testAccount1Login = async () => {
    const account = selectedAccount1.value;
    if (!account) return;

    try {
        const data = await post("/tyys/login", {
            username: account.Username,
            password: account.Password,
        });

        if (data.message === "success") {
            account.loginStatus = "success";
            ElMessage.success("主账号登录测试成功");
        } else {
            account.loginStatus = "error";
            ElMessage.error("主账号登录测试失败");
        }
    } catch (error) {
        account.loginStatus = "error";
        ElMessage.error("主账号登录测试失败");
    }
};

// 测试同伴账号登录（获取同伴码）
const testAccount2Login = async () => {
    const account = selectedAccount2.value;
    if (!account) return;

    try {
        const data = await post("/tyys/buddy_num", {
            username: account.Username,
            password: account.Password,
        });

        if (data.message === "success") {
            account.loginStatus = "success";
            ElMessage.success("同伴账号登录测试成功");
        } else {
            account.loginStatus = "error";
            ElMessage.error("同伴账号登录测试失败");
        }
    } catch (error) {
        account.loginStatus = "error";
        ElMessage.error("同伴账号登录测试失败");
    }
};

// 获取登录按钮文本
const getLoginButtonText = (status: string | null) => {
    switch (status) {
        case "success":
            return "登录成功";
        case "error":
            return "登录失败";
        default:
            return "测试登录";
    }
};

const getLoginButtonType = (status: string | null) => {
    switch (status) {
        case "success":
            return "success";
        case "error":
            return "danger";
        default:
            return "primary";
    }
};

const getLoginButtonIcon = (status: string | null) => {
    switch (status) {
        case "success":
            return SuccessFilled;
        case "error":
            return CircleClose;
        default:
            return Select;
    }
};

// 获取任务列表
const fetchTasks = async () => {
    loading.value = true;
    try {
        const res = await getBargainTasks();
        if (res.message === "success") {
            tasks.value = res.data || [];
        }
    } catch (error) {
        console.error("获取任务列表失败:", error);
        ElMessage.error("获取任务列表失败");
    } finally {
        loading.value = false;
    }
};

// 创建任务
const handleCreateTask = async () => {
    if (!taskFormRef.value) return;

    await taskFormRef.value.validate(async (valid: boolean) => {
        if (!valid) return;

        creating.value = true;
        try {
            const data: any = {
                account_id_1: newTask.value.account_id_1!,
                account_id_2: newTask.value.account_id_2!,
                venue_site_id: newTask.value.venue_site_id,
                reservation_date: newTask.value.reservation_date,
                scan_interval: newTask.value.scan_interval,
            };

            // 只有在填写了可选字段时才传递
            if (newTask.value.site_name) {
                data.site_name = newTask.value.site_name;
            }
            if (newTask.value.reservation_time) {
                data.reservation_time = newTask.value.reservation_time;
            }
            if (newTask.value.deadline) {
                data.deadline = newTask.value.deadline;
            }

            const res = await createBargainTask(data);
            if (res.message === "捡漏任务创建成功") {
                ElMessage.success("创建成功");
                showCreateDialog.value = false;
                resetForm();
                fetchTasks();
            } else {
                ElMessage.error(res.message || "创建失败");
            }
        } catch (error: any) {
            console.error("创建任务失败:", error);
            ElMessage.error(error.message || "创建任务失败");
        } finally {
            creating.value = false;
        }
    });
};

// 查看任务详情
const viewTaskDetail = async (task: any) => {
    currentTask.value = task;
    showDetailDialog.value = true;
    fetchLogs();
};

// 获取日志
const fetchLogs = async () => {
    if (!currentTask.value) return;

    logsLoading.value = true;
    try {
        const res = await getBargainTaskLogs(currentTask.value.task_id);
        if (res.message === "success") {
            logs.value = res.data || [];
        }
    } catch (error) {
        console.error("获取日志失败:", error);
    } finally {
        logsLoading.value = false;
    }
};

// 编辑任务
const handleEditTask = (task: any) => {
    // 复制任务数据到编辑表单
    editingTask.value = {
        task_id: task.task_id,
        account_id_1: task.account_id_1,
        account_id_2: task.account_id_2,
        venue_id: task.venue_site_id === "143" ? 1 : 2,
        venue_site_id: task.venue_site_id,
        reservation_date: task.reservation_date,
        site_name: task.site_name || "",
        reservation_time: task.reservation_time || "",
        scan_interval: task.scan_interval,
        deadline: task.deadline
            ? new Date(task.deadline)
                  .toISOString()
                  .slice(0, 19)
                  .replace("T", " ")
            : "",
    };

    // 加载场地列表
    const venueId = editingTask.value.venue_id;
    if (venueId) {
        const locationCount = venueId === 1 ? 12 : 20;
        editAvailableLocations.value = Array.from(
            { length: locationCount },
            (_, i) => ({
                value: `${i + 1}号`,
                label: `${i + 1}号场地`,
            }),
        );
    }

    showEditDialog.value = true;
};

// 更新任务
const handleUpdateTask = async () => {
    if (!editFormRef.value) return;

    await editFormRef.value.validate(async (valid: boolean) => {
        if (!valid) return;

        updating.value = true;
        try {
            const data: any = {
                account_id_1: editingTask.value.account_id_1,
                account_id_2: editingTask.value.account_id_2,
                venue_site_id: editingTask.value.venue_id === 1 ? "143" : "23",
                reservation_date: editingTask.value.reservation_date,
                scan_interval: editingTask.value.scan_interval,
            };

            // 只有在填写了可选字段时才传递
            if (editingTask.value.site_name) {
                data.site_name = editingTask.value.site_name;
            }
            if (editingTask.value.reservation_time) {
                data.reservation_time = editingTask.value.reservation_time;
            }
            if (editingTask.value.deadline) {
                data.deadline = editingTask.value.deadline;
            }

            const res = await updateBargainTask(
                editingTask.value.task_id,
                data,
            );
            if (res.message === "任务更新成功") {
                ElMessage.success("更新成功");
                showEditDialog.value = false;
                editingTask.value = null;
                fetchTasks();
            } else {
                ElMessage.error(res.message || "更新失败");
            }
        } catch (error: any) {
            console.error("更新任务失败:", error);
            ElMessage.error(error.message || "更新任务失败");
        } finally {
            updating.value = false;
        }
    });
};

// 取消任务
const handleCancelTask = async (task: any) => {
    try {
        await ElMessageBox.confirm("确定要取消这个捡漏任务吗？", "提示", {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "warning",
        });

        const res = await cancelBargainTask(task.task_id);
        if (res.message === "任务已取消") {
            ElMessage.success("取消成功");
            fetchTasks();
        }
    } catch (error: any) {
        if (error !== "cancel") {
            console.error("取消任务失败:", error);
            ElMessage.error("取消任务失败");
        }
    }
};

// 重置表单
const resetForm = () => {
    newTask.value = {
        account_id_1: null,
        account_id_2: null,
        venue_id: null,
        venue_site_id: "",
        reservation_date: "",
        site_name: "",
        reservation_time: "",
        scan_interval: 10,
        deadline: "",
    };
    availableLocations.value = [];

    // 重置账号登录状态
    accounts.value.forEach((acc) => {
        acc.loginStatus = null;
    });
};

// 格式化时间
const formatTime = (time: string) => {
    if (!time) return "";
    return new Date(time).toLocaleString("zh-CN");
};

// 获取状态类型
const getStatusType = (status: string) => {
    const map: Record<string, any> = {
        active: "success",
        completed: "info",
        cancelled: "warning",
        failed: "danger",
        paused: "info",
    };
    return map[status] || "info";
};

// 获取状态文本
const getStatusText = (status: string) => {
    const map: Record<string, string> = {
        active: "运行中",
        completed: "已完成",
        cancelled: "已取消",
        failed: "已失败",
        paused: "已暂停",
    };
    return map[status] || status;
};

// 页面加载时获取数据
onMounted(() => {
    fetchAccounts();
    fetchTasks();
    initTimeSlots(); // 初始化时间段选项

    // 每30秒自动刷新任务列表
    setInterval(() => {
        if (!showDetailDialog.value) {
            fetchTasks();
        }
    }, 30000);
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

:deep(.el-table) {
    font-size: 14px;
}

:deep(.el-descriptions__label) {
    font-weight: 600;
}

:deep(.el-timeline-item__timestamp) {
    color: #909399;
    font-size: 12px;
}

.account-info-display {
    display: flex;
    gap: 8px;
    align-items: center;
    width: 100%;
    flex-wrap: wrap;
}
</style>
