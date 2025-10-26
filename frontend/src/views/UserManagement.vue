<template>
    <Layout>
        <el-main class="content">
                <el-card shadow="hover">
                    <template #header>
                        <div class="card-header">
                            <div class="header-title">
                                <el-icon><UserFilled /></el-icon>
                                <span>用户管理</span>
                            </div>
                            <el-button
                                :icon="Refresh"
                                circle
                                @click="fetchUserList"
                                :loading="loading"
                            />
                        </div>
                    </template>

                    <el-table
                        :data="users"
                        stripe
                        border
                        style="width: 100%"
                        v-loading="loading"
                    >
                        <el-table-column
                            prop="username"
                            label="用户名"
                            align="center"
                            width="180"
                        />

                        <el-table-column
                            label="创建时间"
                            align="center"
                            width="200"
                        >
                            <template #default="scope">
                                {{ formatDateTime(scope.row.createdAt) }}
                            </template>
                        </el-table-column>

                        <el-table-column
                            label="角色"
                            align="center"
                            width="120"
                        >
                            <template #default="scope">
                                <el-tag
                                    :type="
                                        scope.row.is_admin
                                            ? 'danger'
                                            : 'primary'
                                    "
                                    effect="dark"
                                    size="small"
                                >
                                    {{
                                        scope.row.is_admin
                                            ? "管理员"
                                            : "普通用户"
                                    }}
                                </el-tag>
                            </template>
                        </el-table-column>

                        <el-table-column
                            label="操作"
                            align="center"
                            fixed="right"
                            width="180"
                        >
                            <template #default="scope">
                                <el-button
                                    type="primary"
                                    size="small"
                                    :icon="Edit"
                                    @click="showEditModal(scope.row)"
                                >
                                    修改
                                </el-button>
                                <el-button
                                    type="danger"
                                    size="small"
                                    :icon="Delete"
                                    @click="handleDelete(scope.row)"
                                >
                                    删除
                                </el-button>
                            </template>
                        </el-table-column>
                    </el-table>

                    <el-button
                        class="add-btn"
                        type="success"
                        :icon="Plus"
                        @click="showAddModal"
                    >
                        添加用户
                    </el-button>
                </el-card>

                <!-- 添加/编辑用户对话框 -->
                <el-dialog
                    v-model="showModal"
                    :title="isEditing ? '编辑用户' : '添加用户'"
                    width="500px"
                    :close-on-click-modal="false"
                >
                    <el-form
                        ref="formRef"
                        :model="currentUser"
                        :rules="formRules"
                        label-width="100px"
                    >
                        <el-form-item label="用户名" prop="username">
                            <el-input
                                v-model="currentUser.username"
                                :readonly="isEditing"
                                placeholder="请输入用户名"
                            />
                        </el-form-item>

                        <el-form-item label="密码" prop="password">
                            <el-input
                                v-model="currentUser.password"
                                type="password"
                                placeholder="请输入密码"
                                show-password
                            />
                        </el-form-item>

                        <el-form-item label="管理员权限">
                            <el-switch
                                v-model="currentUser.is_admin"
                                active-text="是"
                                inactive-text="否"
                            />
                        </el-form-item>
                    </el-form>

                    <template #footer>
                        <el-button @click="closeModal">取消</el-button>
                        <el-button type="primary" @click="handleSubmit">
                            {{ isEditing ? "更新" : "添加" }}
                        </el-button>
                    </template>
                </el-dialog>
        </el-main>
    </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import {
    UserFilled,
    Refresh,
    Edit,
    Delete,
    Plus,
} from "@element-plus/icons-vue";
import Layout from "@/components/Layout.vue";
import { get, post, patch, del } from "@/utils/api";

const users = ref([]);
const showModal = ref(false);
const isEditing = ref(false);
const loading = ref(false);
const formRef = ref(null);

const currentUser = ref({
    username: "",
    password: "",
    is_admin: false,
});

const formRules = {
    username: [
        { required: true, message: "请输入用户名", trigger: "blur" },
        {
            min: 3,
            max: 20,
            message: "用户名长度在 3 到 20 个字符",
            trigger: "blur",
        },
    ],
    password: [
        { required: true, message: "请输入密码", trigger: "blur" },
        { min: 6, message: "密码长度至少 6 个字符", trigger: "blur" },
    ],
};

const fetchUserList = async () => {
    loading.value = true;
    try {
        const data = await get("/user/all");
        if (data.message === "success") {
            users.value = data.data.map((user) => ({
                username: user.username,
                password: user.password,
                is_admin: user.is_admin,
                createdAt: user.created_at,
                captcha_api: user.captcha_api,
            }));
        } else {
            throw new Error(data.message || "获取用户列表失败");
        }
    } catch (error) {
        console.error("Fetch user list error:", error);
        ElMessage.error(error.message || "获取用户列表失败");
    } finally {
        loading.value = false;
    }
};

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

const showAddModal = () => {
    isEditing.value = false;
    currentUser.value = {
        username: "",
        password: "",
        is_admin: false,
    };
    showModal.value = true;
};

const showEditModal = (user) => {
    isEditing.value = true;
    currentUser.value = { ...user, password: "" };
    showModal.value = true;
};

const closeModal = () => {
    showModal.value = false;
    if (formRef.value) {
        formRef.value.resetFields();
    }
};

const handleSubmit = async () => {
    if (!formRef.value) return;

    try {
        await formRef.value.validate();

        let data;
        if (isEditing.value) {
            data = await patch(
                `/users/${encodeURIComponent(currentUser.value.username)}`,
                currentUser.value,
            );
        } else {
            data = await post("/user/register", {
                username: currentUser.value.username,
                password: currentUser.value.password,
                is_admin: currentUser.value.is_admin,
            });
        }

        if (data.message === "success") {
            ElMessage.success(`用户${isEditing.value ? "更新" : "添加"}成功`);
            closeModal();
            await fetchUserList();
        } else {
            throw new Error(
                data.data || `用户${isEditing.value ? "更新" : "添加"}失败`,
            );
        }
    } catch (error) {
        if (error.message) {
            console.error("Error:", error);
            ElMessage.error(
                error.message || `${isEditing.value ? "更新" : "添加"}用户失败`,
            );
        }
    }
};

const handleDelete = async (user) => {
    try {
        await ElMessageBox.confirm(
            `确定要删除用户 "${user.username}" 吗？`,
            "提示",
            {
                confirmButtonText: "确定",
                cancelButtonText: "取消",
                type: "warning",
            },
        );

        const data = await del(`/users/${encodeURIComponent(user.username)}`);
        if (data.message === "success") {
            ElMessage.success("用户删除成功");
            await fetchUserList();
        } else {
            throw new Error(data.data || "删除用户失败");
        }
    } catch (error) {
        if (error !== "cancel") {
            console.error("Delete error:", error);
            ElMessage.error(error.message || "删除用户失败");
        }
    }
};

onMounted(() => {
    fetchUserList();
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

.add-btn {
    margin-top: 20px;
}

:deep(.el-card__header) {
    padding: 18px 20px;
    border-bottom: 1px solid var(--el-border-color-light);
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

:deep(.el-form-item__label) {
    font-weight: 500;
}

:deep(.el-switch__label) {
    font-size: 14px;
}
</style>
