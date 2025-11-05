<template>
    <Layout>
        <el-main class="content">
            <!-- 账户信息卡片 -->
            <el-card shadow="hover" class="info-card">
                <template #header>
                    <div class="card-header">
                        <div class="header-title">
                            <el-icon><Setting /></el-icon>
                            <span>账户设置</span>
                        </div>
                    </div>
                </template>

                <div class="account-info-row">
                    <el-tag type="primary" size="large" effect="plain">
                        <el-icon><User /></el-icon>
                        <span style="margin-left: 8px"
                            >用户名: {{ username }}</span
                        >
                    </el-tag>
                </div>
            </el-card>

            <!-- 修改密码卡片 -->
            <el-card shadow="hover" class="password-card">
                <template #header>
                    <div class="card-header">
                        <div class="header-title">
                            <el-icon><Key /></el-icon>
                            <span>修改密码</span>
                        </div>
                    </div>
                </template>

                <el-form
                    ref="passwordFormRef"
                    :model="passwordForm"
                    :rules="passwordRules"
                    label-width="100px"
                >
                    <el-form-item label="当前密码" prop="oldPassword">
                        <el-input
                            v-model="passwordForm.oldPassword"
                            type="password"
                            placeholder="请输入当前密码"
                            show-password
                            style="max-width: 400px"
                        />
                    </el-form-item>

                    <el-form-item label="新密码" prop="newPassword">
                        <el-input
                            v-model="passwordForm.newPassword"
                            type="password"
                            placeholder="请输入新密码"
                            show-password
                            style="max-width: 400px"
                        />
                    </el-form-item>

                    <el-form-item label="确认新密码" prop="confirmPassword">
                        <el-input
                            v-model="passwordForm.confirmPassword"
                            type="password"
                            placeholder="请再次输入新密码"
                            show-password
                            style="max-width: 400px"
                        />
                    </el-form-item>

                    <el-form-item>
                        <el-button
                            type="primary"
                            :icon="Check"
                            @click="handlePasswordChange"
                        >
                            修改密码
                        </el-button>
                    </el-form-item>
                </el-form>
            </el-card>

            <!-- API KEY 卡片 -->
            <el-card shadow="hover" class="api-card">
                <template #header>
                    <div class="card-header">
                        <div class="header-title">
                            <el-icon><Key /></el-icon>
                            <span>API KEY</span>
                        </div>
                    </div>
                </template>

                <div class="api-key-group">
                    <el-input
                        v-model="apiKey"
                        :readonly="!isEditingApiKey"
                        placeholder="请输入 API Key"
                        style="max-width: 600px"
                    />
                    <el-button
                        :type="isEditingApiKey ? 'success' : 'primary'"
                        :icon="isEditingApiKey ? Check : Edit"
                        @click="toggleApiKeyEdit"
                    >
                        {{ isEditingApiKey ? "保存" : "修改" }}
                    </el-button>
                </div>
            </el-card>

            <!-- 退出登录（页面底部） -->
            <el-card shadow="hover" class="logout-card">
                <div class="logout-row">
                    <el-button
                        type="danger"
                        :icon="SwitchButton"
                        @click="handleLogout"
                        size="large"
                        style="width: 100%"
                    >
                        退出登录
                    </el-button>
                </div>
            </el-card>
        </el-main>
    </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import {
    User,
    Key,
    SwitchButton,
    Check,
    Edit,
    Setting,
} from "@element-plus/icons-vue";
import Layout from "@/components/Layout.vue";
import { patch, logout } from "@/utils/api";
import {
    getUsername,
    isAdmin as isAdminUser,
    getCaptchaApi,
    updateCaptchaApi,
} from "@/utils/auth";

const router = useRouter();
const passwordFormRef = ref(null);

const username = ref("");
const isAdmin = ref(false);
const apiKey = ref(getCaptchaApi() || "");
const isEditingApiKey = ref(false);

const passwordForm = ref({
    oldPassword: "",
    newPassword: "",
    confirmPassword: "",
});

const validateConfirmPassword = (rule, value, callback) => {
    if (value !== passwordForm.value.newPassword) {
        callback(new Error("两次输入的密码不一致"));
    } else {
        callback();
    }
};

const passwordRules = {
    oldPassword: [
        { required: true, message: "请输入当前密码", trigger: "blur" },
    ],
    newPassword: [
        { required: true, message: "请输入新密码", trigger: "blur" },
        { min: 6, message: "密码长度至少 6 个字符", trigger: "blur" },
    ],
    confirmPassword: [
        { required: true, message: "请再次输入新密码", trigger: "blur" },
        { validator: validateConfirmPassword, trigger: "blur" },
    ],
};

const loadUserInfo = () => {
    username.value = getUsername() || "";
    isAdmin.value = isAdminUser();
    apiKey.value = getCaptchaApi() || "";
};

const handlePasswordChange = async () => {
    if (!passwordFormRef.value) return;

    try {
        await passwordFormRef.value.validate();

        const data = await patch("/user/password", {
            username: username.value,
            password: passwordForm.value.oldPassword,
            new_password: passwordForm.value.newPassword,
        });
        if (data.message === "success") {
            ElMessage.success("密码修改成功");
            clearPasswordFields();
        } else {
            throw new Error(data.data || "密码修改失败");
        }
    } catch (error) {
        if (error.message) {
            console.error("Error:", error);
            ElMessage.error(error.message || "密码修改失败");
        }
    }
};

const clearPasswordFields = () => {
    passwordForm.value = {
        oldPassword: "",
        newPassword: "",
        confirmPassword: "",
    };
    if (passwordFormRef.value) {
        passwordFormRef.value.resetFields();
    }
};

const toggleApiKeyEdit = async () => {
    if (isEditingApiKey.value) {
        try {
            const data = await patch("/user/captcha_api", {
                user: username.value,
                captcha_api: apiKey.value,
            });
            if (data.message === "success") {
                updateCaptchaApi(apiKey.value);
                ElMessage.success("API Key 更新成功");
            } else {
                throw new Error(data.data || "API Key 更新失败");
            }
        } catch (error) {
            console.error("Error:", error);
            ElMessage.error(error.message || "API Key 更新失败");
            apiKey.value = getCaptchaApi() || "";
        }
    }
    isEditingApiKey.value = !isEditingApiKey.value;
};

const handleLogout = async () => {
    try {
        await ElMessageBox.confirm("确定要退出登录吗？", "提示", {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "warning",
        });

        await logout();
        ElMessage.success("登出成功");
        router.push("/login");
    } catch (error) {
        if (error !== "cancel") {
            console.error("Logout error:", error);
            ElMessage.error("登出失败");
        }
    }
};

onMounted(() => {
    loadUserInfo();
});
</script>

<style scoped>
.content {
    padding: 24px;
    background-color: var(--background);
}

.info-card,
.password-card,
.api-card {
    margin-bottom: 24px;
    max-width: 900px;
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

.account-info-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 0;
}

.api-key-group {
    display: flex;
    align-items: center;
    gap: 16px;
}

:deep(.el-card__header) {
    padding: 18px 20px;
    border-bottom: 1px solid var(--el-border-color-light);
}

:deep(.el-form-item__label) {
    font-weight: 500;
}

:deep(.el-tag) {
    padding: 8px 16px;
    font-size: 14px;
}

.logout-card {
    max-width: 900px;
}

.logout-row {
    display: flex;
    align-items: center;
}
</style>
