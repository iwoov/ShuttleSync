<template>
    <Layout>
        <el-main class="content">
                <el-card class="main-card" shadow="hover">
                    <template #header>
                        <div class="card-header">
                            <div class="header-title">
                                <el-icon><OfficeBuilding /></el-icon>
                                <span>场馆预约</span>
                            </div>
                            <div class="button-group">
                                <el-button
                                    type="primary"
                                    :icon="Plus"
                                    @click="addNewReservation"
                                >
                                    普通模式
                                </el-button>
                                <el-button
                                    type="warning"
                                    :icon="Clock"
                                    @click="showDevelopingMessage"
                                >
                                    定时模式
                                </el-button>
                                <el-button
                                    type="info"
                                    :icon="Search"
                                    @click="showDevelopingMessage"
                                >
                                    捡漏模式
                                </el-button>
                            </div>
                        </div>
                    </template>

                    <div class="cards-grid">
                        <el-card
                            v-for="(reservation, index) in reservations"
                            :key="index"
                            class="reservation-card"
                            shadow="hover"
                        >
                            <el-form
                                :model="reservation"
                                label-width="60px"
                                @submit.prevent="submitReservation(index)"
                            >
                                <!-- 预约账号选择 -->
                                <el-form-item label="预约账号">
                                    <el-select
                                        v-if="!reservation.accountId"
                                        v-model="reservation.accountId"
                                        placeholder="选择预约账号"
                                        style="width: 100%"
                                    >
                                        <el-option
                                            v-for="account in accounts"
                                            :key="account.ID"
                                            :label="`${account.Lable} (${account.Username})`"
                                            :value="account.Username"
                                        />
                                    </el-select>
                                    <div v-else class="account-info-group">
                                        <el-input
                                            :model-value="
                                                getSelectedAccount(
                                                    reservation.accountId,
                                                ).Username
                                            "
                                            readonly
                                            style="width: 120px"
                                        />
                                        <el-input
                                            :model-value="
                                                getSelectedAccount(
                                                    reservation.accountId,
                                                ).Password
                                            "
                                            :type="
                                                showPassword
                                                    ? 'text'
                                                    : 'password'
                                            "
                                            readonly
                                            style="width: 120px"
                                        >
                                            <template #suffix>
                                                <el-icon
                                                    @click="togglePassword"
                                                    style="cursor: pointer"
                                                >
                                                    <component
                                                        :is="
                                                            showPassword
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
                                                    getSelectedAccount(
                                                        reservation.accountId,
                                                    ).loginStatus,
                                                )
                                            "
                                            :icon="
                                                getLoginButtonIcon(
                                                    getSelectedAccount(
                                                        reservation.accountId,
                                                    ).loginStatus,
                                                )
                                            "
                                            @click="
                                                testLogin(
                                                    reservation.accountId,
                                                    false,
                                                )
                                            "
                                            style="flex-shrink: 0; width: 100px; white-space: nowrap"
                                        >
                                            {{
                                                getLoginButtonText(
                                                    getSelectedAccount(
                                                        reservation.accountId,
                                                    ).loginStatus,
                                                )
                                            }}
                                        </el-button>
                                    </div>
                                </el-form-item>

                                <!-- 同伴账号选择 -->
                                <el-form-item label="同伴账号">
                                    <el-select
                                        v-if="!reservation.partnerId"
                                        v-model="reservation.partnerId"
                                        placeholder="选择同伴账号"
                                        style="width: 100%"
                                    >
                                        <el-option
                                            v-for="partner in partners"
                                            :key="partner.ID"
                                            :label="`${partner.Lable} (${partner.Username})`"
                                            :value="partner.Username"
                                        />
                                    </el-select>
                                    <div v-else class="account-info-group">
                                        <el-input
                                            :model-value="
                                                getSelectedAccount(
                                                    reservation.partnerId,
                                                ).Username
                                            "
                                            readonly
                                            style="width: 120px"
                                        />
                                        <el-input
                                            :model-value="
                                                getSelectedAccount(
                                                    reservation.partnerId,
                                                ).Password
                                            "
                                            :type="
                                                showPartnerPassword
                                                    ? 'text'
                                                    : 'password'
                                            "
                                            readonly
                                            style="width: 120px"
                                        >
                                            <template #suffix>
                                                <el-icon
                                                    @click="
                                                        togglePartnerPassword
                                                    "
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
                                                    getSelectedAccount(
                                                        reservation.partnerId,
                                                    ).loginStatus,
                                                )
                                            "
                                            :icon="
                                                getLoginButtonIcon(
                                                    getSelectedAccount(
                                                        reservation.partnerId,
                                                    ).loginStatus,
                                                )
                                            "
                                            @click="
                                                testLogin(
                                                    reservation.partnerId,
                                                    true,
                                                )
                                            "
                                            style="flex-shrink: 0; width: 100px; white-space: nowrap"
                                        >
                                            {{
                                                getLoginButtonText(
                                                    getSelectedAccount(
                                                        reservation.partnerId,
                                                    ).loginStatus,
                                                )
                                            }}
                                        </el-button>
                                    </div>
                                </el-form-item>

                                <!-- 场馆场地选择 -->
                                <el-form-item label="场馆场地">
                                    <div class="venue-group">
                                        <el-select
                                            v-model="reservation.venueId"
                                            placeholder="选择场馆"
                                            @change="fetchLocations(index)"
                                            style="flex: 1"
                                        >
                                            <el-option
                                                v-for="venue in venues"
                                                :key="venue.id"
                                                :label="venue.name"
                                                :value="venue.id"
                                            />
                                        </el-select>
                                        <el-select
                                            v-model="reservation.locationId"
                                            placeholder="选择场地"
                                            style="flex: 1"
                                        >
                                            <el-option
                                                v-for="location in locations[
                                                    index
                                                ]"
                                                :key="location.id"
                                                :label="location.name"
                                                :value="location.id"
                                            />
                                        </el-select>
                                    </div>
                                </el-form-item>

                                <!-- 日期时间选择 -->
                                <el-form-item label="预约时间">
                                    <div class="datetime-group">
                                        <el-date-picker
                                            v-model="reservation.date"
                                            type="date"
                                            placeholder="选择日期"
                                            format="YYYY-MM-DD"
                                            value-format="YYYY-MM-DD"
                                            :disabled-date="disabledDate"
                                            @change="fetchTimeSlots(index)"
                                            style="flex: 2"
                                        />
                                        <el-select
                                            v-model="reservation.timeSlot"
                                            placeholder="选择时间"
                                            style="flex: 3"
                                        >
                                            <el-option
                                                v-for="slot in timeSlots[index]"
                                                :key="slot.value"
                                                :label="slot.label"
                                                :value="slot.value"
                                            />
                                        </el-select>
                                    </div>
                                </el-form-item>

                                <!-- 提交按钮 -->
                                <el-form-item>
                                    <el-button
                                        type="primary"
                                        :icon="Check"
                                        style="width: 100%"
                                        @click="submitReservation(index)"
                                        :disabled="!canSubmit(index)"
                                        size="large"
                                    >
                                        立即提交预约
                                    </el-button>
                                    <el-text
                                        v-if="!canSubmit(index)"
                                        type="info"
                                        size="small"
                                        style="margin-top: 8px; display: block"
                                    >
                                        {{ getSubmitButtonTitle(index) }}
                                    </el-text>
                                </el-form-item>
                            </el-form>
                        </el-card>
                    </div>
                </el-card>
        </el-main>
    </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { ElMessage } from "element-plus";
import {
    OfficeBuilding,
    Plus,
    Clock,
    Search,
    Check,
    Select,
    SuccessFilled,
    CircleClose,
} from "@element-plus/icons-vue";
import Layout from "@/components/Layout.vue";
import Sidebar from "@/components/Sidebar.vue";
import { get, post } from "@/utils/api";
import { getUsername, getCaptchaApi } from "@/utils/auth";

const accounts = ref([]);
const partners = ref([]);
const venues = ref([
    { id: 1, name: "体育馆" },
    { id: 2, name: "风雨操场" },
]);
const locations = ref({});
const timeSlots = ref({});
const reservations = ref([getEmptyReservation()]);
const showPassword = ref(false);
const showPartnerPassword = ref(false);
const apiKey = ref(getCaptchaApi() || "");

function getEmptyReservation() {
    return {
        accountId: "",
        partnerId: "",
        venueId: "",
        locationId: "",
        date: "",
        timeSlot: "",
        accountInfo: null,
        partnerInfo: null,
    };
}

const disabledDate = (time) => {
    return time.getTime() < Date.now() - 8.64e7;
};

const addNewReservation = () => {
    const newIndex = reservations.value.length;
    reservations.value.push(getEmptyReservation());
    locations.value[newIndex] = [];
    timeSlots.value[newIndex] = [];
};

const fetchAccounts = async () => {
    try {
        const username = getUsername();
        if (!username) {
            ElMessage.error("获取账号列表失败，请重新登录");
            return;
        }
        const data = await get("/account/list", { user: username });
        if (data.message === "success") {
            accounts.value = data.data.map((account) => ({
                ...account,
                loginStatus: null,
            }));
            partners.value = accounts.value;
        } else {
            throw new Error(data.message || "获取账号列表失败");
        }
    } catch (error) {
        console.error("获取账号列表失败", error);
        ElMessage.error(error.message || "获取账号列表失败");
    }
};

const fetchLocations = (index) => {
    const venueId = reservations.value[index].venueId;
    if (!venueId) return;

    const locationCount = Number(venueId) === 1 ? 12 : 20;
    locations.value[index] = Array.from({ length: locationCount }, (_, i) => ({
        id: i + 1,
        name: `${i + 1}号场地`,
    }));
};

const fetchTimeSlots = (index) => {
    if (
        !reservations.value[index].locationId ||
        !reservations.value[index].date
    )
        return;

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

    timeSlots.value[index] = slots;
};

const submitReservation = async (index) => {
    const reservation = reservations.value[index];

    if (!reservation.accountInfo || !reservation.partnerInfo) {
        ElMessage.warning("请先完成账号登录测试");
        return;
    }

    try {
        const venue_site_id = Number(reservation.venueId) === 1 ? "143" : "23";
        const currentUser = getUsername();
        if (!currentUser) {
            ElMessage.error("请先登录后再提交预约");
            return;
        }

        const requestData = {
            user: currentUser,
            username: reservation.accountInfo.username,
            password: reservation.accountInfo.password,
            user_phone: reservation.accountInfo.phone,
            captcha_api: apiKey.value,
            buddy_user_id: reservation.partnerInfo.buddy_id,
            buddy_num: reservation.partnerInfo.buddy_num,
            venue_site_id: venue_site_id,
            reservation_date: reservation.date,
            reservation_time: reservation.timeSlot,
            site_name: `${reservation.locationId}号`,
        };

        const data = await post("/task/add", requestData);
        if (data.message === "success") {
            ElMessage.success("预约提交成功");
            reservations.value[index] = getEmptyReservation();
        } else {
            throw new Error(data.data || "预约提交失败");
        }
    } catch (error) {
        ElMessage.error(error.message || "预约提交失败");
    }
};

const getSelectedAccount = (username) => {
    return (
        accounts.value.find((account) => account.Username === username) || {}
    );
};

const togglePassword = () => {
    showPassword.value = !showPassword.value;
};

const togglePartnerPassword = () => {
    showPartnerPassword.value = !showPartnerPassword.value;
};

const getLoginButtonText = (status) => {
    switch (status) {
        case "success":
            return "登录成功";
        case "error":
            return "登录失败";
        default:
            return "测试登录";
    }
};

const getLoginButtonType = (status) => {
    switch (status) {
        case "success":
            return "success";
        case "error":
            return "danger";
        default:
            return "primary";
    }
};

const getLoginButtonIcon = (status) => {
    switch (status) {
        case "success":
            return SuccessFilled;
        case "error":
            return CircleClose;
        default:
            return Select;
    }
};

const testLogin = async (username, isPartner = false) => {
    const account = getSelectedAccount(username);
    if (!account) return;

    try {
        const url = isPartner ? "/tyys/buddy_num" : "/tyys/login";

        const data = await post(url, {
            username: account.Username,
            password: account.Password,
        });

        if (data.message === "success") {
            account.loginStatus = "success";

            const reservationIndex = reservations.value.findIndex(
                (r) => (isPartner ? r.partnerId : r.accountId) === username,
            );

            if (reservationIndex !== -1) {
                if (isPartner) {
                    reservations.value[reservationIndex].partnerInfo = {
                        buddy_id: data.data.buddy_id,
                        buddy_num: data.data.buddy_num,
                    };
                } else {
                    reservations.value[reservationIndex].accountInfo = {
                        username: data.data.username,
                        password: data.data.password,
                        phone: data.data.phone,
                    };
                }
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

const canSubmit = (index) => {
    const reservation = reservations.value[index];
    const account = getSelectedAccount(reservation.accountId);
    const partner = getSelectedAccount(reservation.partnerId);

    return (
        account?.loginStatus === "success" &&
        partner?.loginStatus === "success" &&
        apiKey.value
    );
};

const getSubmitButtonTitle = (index) => {
    const reservation = reservations.value[index];
    const account = getSelectedAccount(reservation.accountId);
    const partner = getSelectedAccount(reservation.partnerId);

    if (account?.loginStatus !== "success") {
        return "请确保预约账号登录测试成功";
    }
    if (partner?.loginStatus !== "success") {
        return "请确保同伴账号登录测试成功";
    }
    if (!apiKey.value) {
        return "请先在个人中心设置API KEY";
    }
    return "提交预约";
};

const showDevelopingMessage = () => {
    ElMessage.info("正在开发中...");
};

onMounted(() => {
    fetchAccounts();
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
    flex-wrap: wrap;
    gap: 16px;
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
    gap: 12px;
    flex-wrap: wrap;
}

.cards-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(420px, 1fr));
    gap: 20px;
    margin-top: 20px;
}

.reservation-card {
    background: var(--el-fill-color-light);
}

.account-info-group {
    display: flex;
    gap: 8px;
    align-items: center;
    width: 100%;
    flex-wrap: nowrap;
}

.venue-group,
.datetime-group {
    display: flex;
    gap: 12px;
    width: 100%;
}

:deep(.el-form-item__label) {
    font-weight: 500;
    padding-right: 4px;
}

@media (max-width: 1200px) {
    .cards-grid {
        grid-template-columns: 1fr;
    }
}

@media (max-width: 768px) {
    .account-info-group {
        flex-direction: column;
        align-items: stretch;
    }

    .account-info-group > * {
        width: 100% !important;
    }
}
</style>
