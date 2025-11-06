<template>
    <Layout>
        <el-main class="content">
            <el-card shadow="hover" class="guide-card">
                <template #header>
                    <div class="card-header">
                        <div class="header-title">
                            <el-icon><Document /></el-icon>
                            <span>使用说明</span>
                        </div>
                    </div>
                </template>

                <div class="markdown-content">
                    <h1>ShuttleSync 场馆预约系统使用指南</h1>

                    <h2>一、前置准备</h2>

                    <h3>1.1 配置云码 API Key</h3>
                    <p>
                        在使用预约功能之前，您需要先配置云码的 API
                        Key，用于自动识别和处理预约过程中的验证码。
                    </p>

                    <div class="step-box">
                        <p><strong>操作步骤：</strong></p>
                        <ol>
                            <li>
                                点击左侧导航栏的
                                <el-tag size="small" type="primary"
                                    >账户设置</el-tag
                                >
                            </li>
                            <li>
                                在「API KEY」区域，点击
                                <el-tag size="small">修改</el-tag> 按钮
                            </li>
                            <li>输入您的云码 API Key</li>
                            <li>
                                点击 <el-tag size="small" type="success"
                                    >保存</el-tag
                                >
                                完成配置
                            </li>
                        </ol>
                    </div>

                    <el-alert
                        title="重要提示"
                        type="warning"
                        :closable="false"
                        show-icon
                        class="alert-box"
                    >
                        <p>
                            没有配置 API Key
                            将无法自动处理验证码，可能导致预约失败。
                        </p>
                    </el-alert>

                    <h3>1.2 添加预约账户</h3>
                    <p>
                        系统需要<strong>至少 2 个账户</strong
                        >才能正常进行预约：
                    </p>
                    <ul>
                        <li>
                            <strong>主账户</strong
                            >：用于发起预约并完成场地预订
                        </li>
                        <li>
                            <strong>同伴账户</strong
                            >：提供同伴码，协助主账户完成预约
                        </li>
                    </ul>

                    <div class="step-box">
                        <p><strong>操作步骤：</strong></p>
                        <ol>
                            <li>
                                点击左侧导航栏的
                                <el-tag size="small" type="primary"
                                    >账号管理</el-tag
                                >
                            </li>
                            <li>
                                点击 <el-tag size="small" type="success"
                                    >添加账号</el-tag
                                >
                                按钮
                            </li>
                            <li>依次填写账户的登录信息（用户名、密码等）</li>
                            <li>重复以上步骤，至少添加 2 个账户</li>
                        </ol>
                    </div>

                    <el-alert
                        title="最佳实践"
                        type="info"
                        :closable="false"
                        show-icon
                        class="alert-box"
                    >
                        <p>
                            建议添加
                            <strong>3 个或更多账户</strong
                            >，以提高预约成功率。当某个账户出现问题时，系统可以自动切换到其他可用账户。
                        </p>
                    </el-alert>

                    <h2>二、预约模式介绍</h2>

                    <p>
                        本系统提供两种预约模式，满足不同的预约需求：<strong
                            >普通模式</strong
                        >
                        和 <strong>捡漏模式</strong>。
                    </p>

                    <h3>2.1 普通预约模式</h3>

                    <h4>适用场景</h4>
                    <p>适合提前规划，抢订新开放的场地时段。</p>

                    <h4>工作原理</h4>
                    <ul>
                        <li>
                            场馆通常会<strong>提前 2 天</strong
                            >开放预约，比如周一早上 9 点开放周三的场地
                        </li>
                        <li>
                            您需要至少<strong>提前 2 天</strong
                            >在系统中提交预约任务
                        </li>
                        <li>
                            系统会在目标日期前 2 天的<strong
                                >早上 9:00</strong
                            >
                            准时开始预约
                        </li>
                        <li>
                            按照您提交的场馆、场地、时间等条件进行自动预约
                        </li>
                    </ul>

                    <div class="step-box">
                        <p><strong>操作步骤：</strong></p>
                        <ol>
                            <li>
                                点击左侧导航栏的
                                <el-tag size="small" type="primary"
                                    >场馆预约</el-tag
                                >
                            </li>
                            <li>选择您想预约的场馆</li>
                            <li>选择日期（需要是 2 天后或更远的日期）</li>
                            <li>选择场地编号</li>
                            <li>选择时间段</li>
                            <li>
                                点击
                                <el-tag size="small" type="success"
                                    >提交预约</el-tag
                                >
                            </li>
                        </ol>
                    </div>

                    <div class="example-box">
                        <p><strong>示例：</strong></p>
                        <p>
                            假设今天是<strong>周一（1月6日）</strong
                            >，您想预约<strong>周三（1月8日）</strong
                            >下午 2-3 点的场地：
                        </p>
                        <ul>
                            <li>
                                您需要在<strong>周一当天</strong
                                >提交预约任务
                            </li>
                            <li>
                                系统会在<strong>周一早上 9:00</strong
                                >自动为您抢订周三的场地
                            </li>
                            <li>预约结果会在提交后显示在「我的预约」中</li>
                        </ul>
                    </div>

                    <h3>2.2 捡漏预约模式</h3>

                    <h4>适用场景</h4>
                    <p>
                        适合预约<strong>已经开放</strong
                        >但尚未被预订满的场地，或者监控已被预订场地的退订情况。
                    </p>

                    <h4>工作原理</h4>
                    <ul>
                        <li>
                            针对<strong>当前已开放</strong
                            >的预约时间段（通常是未来 2
                            天内或已经开放的时段）
                        </li>
                        <li>系统会根据您设置的条件，定时自动查询场地状态</li>
                        <li>
                            一旦发现有符合条件的<strong>空闲场地</strong
                            >，立即自动预约
                        </li>
                        <li>
                            即使场地已满，也会持续监控，捕捉其他用户的退订机会
                        </li>
                    </ul>

                    <div class="step-box">
                        <p><strong>操作步骤：</strong></p>
                        <ol>
                            <li>
                                点击左侧导航栏的
                                <el-tag size="small" type="primary"
                                    >捡漏模式</el-tag
                                >
                            </li>
                            <li>选择目标场馆</li>
                            <li>
                                选择日期（通常是今天、明天或已经开放预约的日期）
                            </li>
                            <li>
                                选择希望预约的场地编号（可选择多个或全部）
                            </li>
                            <li>选择时间段</li>
                            <li>
                                点击
                                <el-tag size="small" type="success"
                                    >开启捡漏</el-tag
                                >
                            </li>
                        </ol>
                    </div>

                    <div class="example-box">
                        <p><strong>示例：</strong></p>
                        <p>
                            假设今天是<strong>周一</strong
                            >，您想预约<strong>今天晚上 7-8 点</strong
                            >的场地，但发现已经满了：
                        </p>
                        <ul>
                            <li>
                                您可以在<strong>捡漏模式</strong
                                >中提交预约任务
                            </li>
                            <li>
                                系统会每隔一定时间（如每 1-5
                                分钟）自动检查场地状态
                            </li>
                            <li>
                                如果有人取消预约，系统会<strong
                                    >立即为您预订</strong
                                >
                            </li>
                        </ul>
                    </div>

                    <el-alert
                        title="注意"
                        type="warning"
                        :closable="false"
                        show-icon
                        class="alert-box"
                    >
                        <p>
                            捡漏模式会持续运行直到成功预约或您手动停止，建议定期查看「我的预约」确认预约状态。
                        </p>
                    </el-alert>

                    <h2>三、查看预约结果</h2>

                    <div class="step-box">
                        <p><strong>操作步骤：</strong></p>
                        <ol>
                            <li>
                                点击左侧导航栏的
                                <el-tag size="small" type="primary"
                                    >我的预约</el-tag
                                >
                            </li>
                            <li>查看所有预约任务的状态</li>
                            <li>
                                可以查看预约详情、取消预约或删除已完成的任务
                            </li>
                        </ol>
                    </div>

                    <h2>四、常见问题</h2>

                    <div class="faq-item">
                        <h4>Q1: 预约失败怎么办？</h4>
                        <p>
                            <strong>A:</strong> 请检查以下几点：
                        </p>
                        <ul>
                            <li>是否已配置云码 API Key</li>
                            <li>账户数量是否至少有 2 个</li>
                            <li>账户登录信息是否正确</li>
                            <li>预约时间是否已经过期</li>
                            <li>场地是否已经被预订满</li>
                        </ul>
                    </div>

                    <div class="faq-item">
                        <h4>Q2: 普通模式和捡漏模式可以同时使用吗？</h4>
                        <p>
                            <strong>A:</strong>
                            可以。两种模式互不冲突，您可以同时提交多个预约任务。
                        </p>
                    </div>

                    <div class="faq-item">
                        <h4>Q3: 捡漏模式会一直运行吗？</h4>
                        <p>
                            <strong>A:</strong>
                            捡漏模式会持续运行直到预约成功或手动停止。您可以在「我的预约」中随时停止或删除捡漏任务。
                        </p>
                    </div>

                    <div class="faq-item">
                        <h4>Q4: 如何提高预约成功率？</h4>
                        <p>
                            <strong>A:</strong> 建议：
                        </p>
                        <ul>
                            <li>添加 3 个或更多账户</li>
                            <li>确保账户状态正常且余额充足</li>
                            <li>提前测试账户登录是否正常</li>
                            <li>合理设置预约时间，避免过于热门的时段</li>
                        </ul>
                    </div>

                    <h2>五、技术支持</h2>

                    <p>
                        如果您在使用过程中遇到任何问题，请及时联系系统管理员获取帮助。
                    </p>

                    <div class="footer-tip">
                        <el-icon><SuccessFilled /></el-icon>
                        <span>祝您预约顺利！</span>
                    </div>
                </div>
            </el-card>
        </el-main>
    </Layout>
</template>

<script setup lang="ts">
import Layout from "@/components/Layout.vue";
import { Document, SuccessFilled } from "@element-plus/icons-vue";
</script>

<style scoped>
.content {
    padding: 24px;
    background-color: var(--background);
}

.guide-card {
    max-width: 1200px;
    margin: 0 auto;
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

.markdown-content {
    line-height: 1.8;
    color: var(--el-text-color-primary);
}

.markdown-content h1 {
    font-size: 28px;
    font-weight: 700;
    color: var(--el-color-primary);
    margin-bottom: 24px;
    padding-bottom: 12px;
    border-bottom: 2px solid var(--el-border-color);
}

.markdown-content h2 {
    font-size: 24px;
    font-weight: 600;
    color: var(--el-text-color-primary);
    margin-top: 32px;
    margin-bottom: 16px;
}

.markdown-content h3 {
    font-size: 20px;
    font-weight: 600;
    color: var(--el-text-color-regular);
    margin-top: 24px;
    margin-bottom: 12px;
}

.markdown-content h4 {
    font-size: 16px;
    font-weight: 600;
    color: var(--el-text-color-regular);
    margin-top: 16px;
    margin-bottom: 8px;
}

.markdown-content p {
    margin: 12px 0;
    font-size: 15px;
}

.markdown-content ul,
.markdown-content ol {
    margin: 12px 0;
    padding-left: 24px;
}

.markdown-content li {
    margin: 8px 0;
    font-size: 15px;
}

.step-box {
    background-color: var(--el-color-primary-light-9);
    border-left: 4px solid var(--el-color-primary);
    padding: 16px 20px;
    margin: 16px 0;
    border-radius: 4px;
}

.step-box p {
    margin: 0 0 12px 0;
}

.step-box ol,
.step-box ul {
    margin: 0;
}

.example-box {
    background-color: var(--el-color-success-light-9);
    border-left: 4px solid var(--el-color-success);
    padding: 16px 20px;
    margin: 16px 0;
    border-radius: 4px;
}

.example-box p {
    margin: 0 0 12px 0;
}

.example-box ul {
    margin: 0;
}

.alert-box {
    margin: 16px 0;
}

.alert-box p {
    margin: 4px 0;
}

.faq-item {
    margin: 20px 0;
    padding: 16px;
    background-color: var(--el-fill-color-light);
    border-radius: 8px;
}

.faq-item h4 {
    margin: 0 0 8px 0;
    color: var(--el-color-primary);
}

.faq-item p {
    margin: 8px 0;
}

.faq-item ul {
    margin: 8px 0;
}

.footer-tip {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    margin-top: 40px;
    padding: 20px;
    font-size: 18px;
    font-weight: 600;
    color: var(--el-color-success);
}

:deep(.el-card__header) {
    padding: 18px 20px;
    border-bottom: 1px solid var(--el-border-color-light);
}

:deep(.el-tag) {
    margin: 0 4px;
    font-weight: 500;
}
</style>
