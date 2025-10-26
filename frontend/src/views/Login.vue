<template>
    <Layout>
        <div class="login-container">
            <el-card class="login-card" shadow="hover">
                <template #header>
                    <div class="card-header">
                        <el-icon :size="32" color="#409eff">
                            <User />
                        </el-icon>
                        <h2>预约系统</h2>
                    </div>
                </template>

                <el-form
                    :model="loginForm"
                    :rules="rules"
                    ref="loginFormRef"
                    @submit.prevent="handleLogin"
                    label-width="80px"
                    size="large"
                >
                    <el-form-item label="用户名" prop="username">
                        <el-input
                            v-model="loginForm.username"
                            placeholder="请输入用户名"
                            clearable
                            :prefix-icon="User"
                        >
                        </el-input>
                    </el-form-item>

                    <el-form-item label="密码" prop="password">
                        <el-input
                            v-model="loginForm.password"
                            type="password"
                            placeholder="请输入密码"
                            show-password
                            :prefix-icon="Lock"
                        >
                        </el-input>
                    </el-form-item>

                    <el-form-item>
                        <el-button
                            type="primary"
                            style="width: 100%"
                            @click="handleLogin"
                            :loading="loading"
                            size="large"
                        >
                            <el-icon v-if="!loading"><Right /></el-icon>
                            登录
                        </el-button>
                    </el-form-item>
                </el-form>
            </el-card>
        </div>
    </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { User, Lock, Right } from "@element-plus/icons-vue";
import Layout from "@/components/Layout.vue";
import { login } from "@/utils/api";
import { saveTokenInfo, isAuthenticated } from "@/utils/auth";

const router = useRouter();
const loginFormRef = ref(null);
const loading = ref(false);

const loginForm = ref({
    username: "",
    password: "",
});

const rules = {
    username: [
        { required: true, message: "请输入用户名", trigger: "blur" },
        {
            min: 2,
            max: 20,
            message: "用户名长度在 2 到 20 个字符",
            trigger: "blur",
        },
    ],
    password: [
        { required: true, message: "请输入密码", trigger: "blur" },
        { min: 6, message: "密码长度至少 6 个字符", trigger: "blur" },
    ],
};

const handleLogin = async () => {
    if (!loginFormRef.value) return;

    await loginFormRef.value.validate(async (valid) => {
        if (valid) {
            loading.value = true;
            try {
                const data = await login(
                    loginForm.value.username,
                    loginForm.value.password,
                );

                if (data.message === "success" && data.data) {
                    // 保存 Token 信息
                    saveTokenInfo({
                        accessToken: data.data.access_token,
                        refreshToken: data.data.refresh_token,
                        username: data.data.username,
                        isAdmin: data.data.is_admin,
                        captchaApi: data.data.captcha_api || "",
                    });

                    ElMessage.success("登录成功");
                    router.push("/my-reservation");
                } else {
                    throw new Error(data.message || "登录失败");
                }
            } catch (error: any) {
                console.error("Login error:", error);
                ElMessage.error(
                    error.message || "登录失败，请检查用户名和密码",
                );
            } finally {
                loading.value = false;
            }
        }
    });
};

onMounted(() => {
    // 如果已经登录，直接跳转到我的预约页面
    if (isAuthenticated()) {
        router.push("/my-reservation");
    }
});
</script>

<style scoped>
.login-container {
    height: 100vh;
    display: flex;
    justify-content: center;
    align-items: center;
    background-color: var(--background);
}

.login-card {
    width: 450px;
    max-width: 90%;
    position: relative;
    z-index: 1;
}

.card-header {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 16px;
}

.card-header h2 {
    color: var(--el-color-primary);
    margin: 0;
    font-size: 24px;
    font-weight: 600;
}

:deep(.el-card__header) {
    padding: 30px 20px 20px;
    background: linear-gradient(135deg, #ecf5ff 0%, #ffffff 100%);
}

:deep(.el-card__body) {
    padding: 30px 40px 40px;
}

:deep(.el-form-item) {
    margin-bottom: 24px;
}

:deep(.el-form-item__label) {
    font-weight: 500;
}
</style>
