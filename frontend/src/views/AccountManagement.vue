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
                        <div class="header-actions">
                            <el-button
                                type="success"
                                :icon="Plus"
                                @click="addNewAccount"
                            >
                                添加账号
                            </el-button>
                            <el-button
                                :icon="Refresh"
                                @click="fetchAccountList"
                                :loading="accountLoading"
                            >
                                刷新
                            </el-button>
                        </div>
                    </div>
                </template>

                <el-table :data="accounts" style="width: 100%" v-loading="accountLoading">
                    <el-table-column prop="Lable" label="标签" width="150">
                        <template #default="{ row }">
                            <el-input
                                v-if="row.isEditing"
                                v-model="row.Lable"
                                placeholder="标签"
                                size="small"
                            />
                            <span v-else>{{ row.Lable }}</span>
                        </template>
                    </el-table-column>

                    <el-table-column prop="Username" label="用户名" width="180">
                        <template #default="{ row }">
                            <el-input
                                v-if="row.isEditing"
                                v-model="row.Username"
                                placeholder="用户名"
                                size="small"
                            />
                            <span v-else>{{ row.Username }}</span>
                        </template>
                    </el-table-column>

                    <el-table-column prop="Password" label="密码" width="180">
                        <template #default="{ row }">
                            <el-input
                                v-if="row.isEditing"
                                v-model="row.Password"
                                type="password"
                                placeholder="密码"
                                size="small"
                                show-password
                            />
                            <el-input
                                v-else
                                :model-value="row.Password"
                                type="password"
                                readonly
                                size="small"
                                show-password
                            />
                        </template>
                    </el-table-column>

                    <el-table-column prop="Phone" label="手机号" width="150">
                        <template #default="{ row }">
                            {{ row.Phone || '-' }}
                        </template>
                    </el-table-column>

                    <el-table-column label="创建时间" width="180">
                        <template #default="{ row }">
                            {{ row.CreatedAt ? formatTime(row.CreatedAt) : '-' }}
                        </template>
                    </el-table-column>

                    <el-table-column label="操作" width="200" fixed="right">
                        <template #default="{ row }">
                            <el-button
                                v-if="!row.isEditing"
                                type="primary"
                                size="small"
                                :icon="Edit"
                                @click="startEdit(row)"
                            >
                                编辑
                            </el-button>
                            <el-button
                                v-else
                                type="success"
                                size="small"
                                :icon="Check"
                                @click="handleModify(row)"
                            >
                                保存
                            </el-button>
                            <el-button
                                v-if="row.isEditing"
                                size="small"
                                @click="cancelEdit(row)"
                            >
                                取消
                            </el-button>
                            <el-button
                                v-if="!row.isEditing"
                                type="danger"
                                size="small"
                                :icon="Delete"
                                @click="handleDelete(row)"
                            >
                                删除
                            </el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </el-card>
        </el-main>
    </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import {
    User,
    Refresh,
    Edit,
    Delete,
    Plus,
    Check,
} from "@element-plus/icons-vue";
import Layout from "@/components/Layout.vue";
import { get, post, patch, del } from "@/utils/api";
import { getUsername } from "@/utils/auth";

const accounts = ref<any[]>([]);
const accountLoading = ref(false);
const originalAccount = ref<any>(null);

const fetchAccountList = async () => {
    accountLoading.value = true;
    try {
        const username = getUsername();
        if (!username) {
            ElMessage.error("请先登录");
            return;
        }
        const data = await get("/account/list", { user: username });
        if (data.message === "success") {
            accounts.value = data.data.map((account: any) => ({
                ...account,
                isEditing: false,
                showPassword: false,
            }));
        }
    } catch (error) {
        console.error("获取账号列表失败:", error);
        ElMessage.error("获取账号列表失败");
    } finally {
        accountLoading.value = false;
    }
};

const startEdit = (account: any) => {
    // 保存原始数据用于取消编辑
    originalAccount.value = { ...account };
    account.isEditing = true;
};

const cancelEdit = (account: any) => {
    if (account.isNew) {
        // 如果是新添加的账号，直接删除
        const index = accounts.value.indexOf(account);
        if (index > -1) {
            accounts.value.splice(index, 1);
        }
    } else {
        // 恢复原始数据
        if (originalAccount.value) {
            Object.assign(account, originalAccount.value);
        }
        account.isEditing = false;
    }
};

const handleModify = async (account: any) => {
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
        } else {
            throw new Error(data.data || "操作失败");
        }
    } catch (error: any) {
        console.error("操作失败:", error);
        ElMessage.error(error.message || (account.isNew ? "添加账号失败" : "更新账号失败"));
    }
};

const handleDelete = async (account: any) => {
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
            console.error("删除失败:", error);
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
        Phone: "",
        CreatedAt: "",
        isEditing: true,
        showPassword: false,
        isNew: true,
    };
    accounts.value.unshift(newAccount);
};

const formatTime = (time: string) => {
    if (!time) return "";
    return new Date(time).toLocaleString("zh-CN");
};

onMounted(() => {
    fetchAccountList();
});
</script>

<style scoped>
.content {
    padding: 20px;
    background-color: #f5f7fa;
    min-height: calc(100vh - 80px);
}

.account-card {
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

.header-actions {
    display: flex;
    gap: 10px;
}

:deep(.el-table) {
    font-size: 14px;
}

:deep(.el-card__header) {
    padding: 18px 20px;
    border-bottom: 1px solid var(--el-border-color-light);
}
</style>
