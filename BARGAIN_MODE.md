# 捡漏模式功能设计文档

## 功能概述

捡漏模式是一个自动化预约系统，允许用户设置定时任务来监控场馆的空余场地，一旦发现符合条件的场地，系统会自动为指定的两个账号进行预约。

## 功能特性

### 1. 核心功能
- ✅ 支持两个账号预约（账号1主预约 + 账号2提供同伴码）
- ✅ 自定义扫描间隔（1-60分钟）
- ✅ 可选的场地号和时间段过滤
- ✅ 自动重试和错误处理
- ✅ 详细的扫描日志记录
- ✅ 服务器重启后自动恢复活跃任务

### 2. 与普通模式的区别
**普通模式：**
- 定时在固定时间（如9:00）执行预约
- 预约指定的场地和时间

**捡漏模式：**
- 在可预约时间段内，定期扫描未被预约的空场地
- 发现符合条件的场地后立即预约
- 可选过滤场地号和时间段

### 3. 任务状态管理
- `active`: 任务运行中
- `paused`: 任务已暂停
- `completed`: 任务完成（成功预约）
- `cancelled`: 任务已取消

## 数据库设计

### BargainTaskDb 表（捡漏任务）
```go
{
    ID              uint      // 主键
    User            string    // 任务所属用户
    TaskID          string    // 任务唯一标识（UUID）
    AccountID1      uint      // 第一个预约账号ID
    AccountID2      uint      // 第二个预约账号ID
    VenueSiteID     string    // 场馆ID（必填）
    ReservationDate string    // 预约日期 YYYY-MM-DD（必填）
    SiteName        string    // 场地号（可选，空则任意场地）
    ReservationTime string    // 时间段（可选，空则任意时间）
    ScanInterval    int       // 扫描间隔（分钟，1-60）
    Status          string    // 任务状态
    SuccessCount    int       // 成功预约次数
    ScanCount       int       // 总扫描次数
    LastScanTime    time.Time // 最后扫描时间
    CreatedAt       time.Time // 创建时间
    UpdatedAt       time.Time // 更新时间
}
```

### BargainLogDb 表（扫描日志）
```go
{
    ID             uint      // 主键
    TaskID         string    // 关联任务ID
    ScanTime       time.Time // 扫描时间
    AvailableSlots int       // 发现的可用场地数
    Success        bool      // 是否成功预约
    Message        string    // 日志消息
    Details        string    // 详细信息（JSON格式）
}
```

## API 接口

### 1. 创建捡漏任务
**POST** `/api/bargain/create`

**请求头：**
```
Authorization: Bearer <access_token>
```

**请求体：**
```json
{
    "account_id_1": 1,           // 第一个账号ID（必填）
    "account_id_2": 2,           // 第二个账号ID（必填）
    "venue_site_id": "xxx",      // 场馆ID（必填）
    "reservation_date": "2025-11-10",  // 预约日期（必填）
    "site_name": "1号场",        // 场地号（可选）
    "reservation_time": "19:00", // 时间段（可选）
    "scan_interval": 10          // 扫描间隔（分钟，1-60）
}
```

**响应示例：**
```json
{
    "message": "捡漏任务创建成功",
    "data": {
        "id": 1,
        "task_id": "uuid-xxx-xxx",
        "status": "active",
        "created_at": "2025-11-05T10:00:00Z"
    }
}
```

---

### 2. 获取用户任务列表
**GET** `/api/bargain/list`

**请求头：**
```
Authorization: Bearer <access_token>
```

**响应示例：**
```json
{
    "message": "success",
    "data": [
        {
            "id": 1,
            "task_id": "uuid-xxx-xxx",
            "venue_site_id": "xxx",
            "reservation_date": "2025-11-10",
            "site_name": "1号场",
            "reservation_time": "19:00",
            "scan_interval": 10,
            "status": "active",
            "success_count": 0,
            "scan_count": 15,
            "last_scan_time": "2025-11-05T10:30:00Z",
            "created_at": "2025-11-05T10:00:00Z"
        }
    ]
}
```

---

### 3. 获取任务详情
**GET** `/api/bargain/:id`

**请求头：**
```
Authorization: Bearer <access_token>
```

**响应示例：**
```json
{
    "message": "success",
    "data": {
        "id": 1,
        "task_id": "uuid-xxx-xxx",
        "account_1_label": "主账号",
        "account_2_label": "副账号",
        "venue_site_id": "xxx",
        "reservation_date": "2025-11-10",
        "site_name": "1号场",
        "reservation_time": "19:00",
        "scan_interval": 10,
        "status": "active",
        "success_count": 0,
        "scan_count": 15,
        "last_scan_time": "2025-11-05T10:30:00Z"
    }
}
```

---

### 4. 取消任务
**DELETE** `/api/bargain/:id`

**请求头：**
```
Authorization: Bearer <access_token>
```

**响应示例：**
```json
{
    "message": "任务已取消"
}
```

---

### 5. 获取任务日志
**GET** `/api/bargain/:id/logs`

**请求头：**
```
Authorization: Bearer <access_token>
```

**响应示例：**
```json
{
    "message": "success",
    "data": [
        {
            "id": 1,
            "task_id": "uuid-xxx-xxx",
            "scan_time": "2025-11-05T10:30:00Z",
            "available_slots": 2,
            "success": false,
            "message": "发现2个可用场地",
            "details": ""
        },
        {
            "id": 2,
            "task_id": "uuid-xxx-xxx",
            "scan_time": "2025-11-05T10:40:00Z",
            "available_slots": 0,
            "success": false,
            "message": "没有可用场地",
            "details": ""
        }
    ]
}
```

---

### 6. 获取所有任务（管理员）
**GET** `/api/bargain/all`

**请求头：**
```
Authorization: Bearer <access_token>
```

**权限要求：** 管理员

**响应示例：** 同任务列表格式

---

## 核心逻辑流程

### 1. 任务创建流程
```
1. 用户提交创建请求
2. 验证两个账号是否存在且属于当前用户
3. 验证预约日期必须在当前日期之后
4. 验证两个账号不能相同
5. 创建任务记录到数据库
6. 启动定时调度器
7. 返回任务信息
```

### 2. 定时扫描流程
```
1. 根据扫描间隔定时执行
2. 检查任务状态（非active则停止）
3. 使用账号1登录获取场馆信息
4. 解析可用场地和时间段
5. 根据过滤条件筛选符合要求的场地
6. 如果有可用场地：
   - 选择第一个可用场地
   - 登录账号2获取同伴码（BuddyNum和UserId）
   - 使用账号1作为主预约账号
   - 使用账号2提供同伴信息
   - 两个账号一起预约同一个场地
   - 更新成功次数
   - 标记任务为completed
   - 停止定时任务
7. 如果没有可用场地：
   - 记录日志
   - 等待下次扫描
8. 更新扫描次数和时间
```

### 3. 场地匹配逻辑
```go
// 遍历所有场地
for each space in venue_info {
    // 如果指定了场地名称，只匹配该场地
    if siteName != "" && space.name != siteName {
        continue  // 未指定时，不过滤，返回所有场地
    }

    // 遍历时间段
    for each time in space.timeSlots {
        // 只选择可预约的时间段（status == 1）
        if time.status == 1 {
            // 如果指定了时间，只匹配该时间
            if reservationTime != "" && time.beginTime != reservationTime {
                continue  // 未指定时，不过滤，返回所有时间
            }

            // 添加到可用列表
            availableSlots.add(space, time)
        }
    }
}

// 从可用列表中随机选择一个
randomIndex = rand.Intn(len(availableSlots))
selectedSlot = availableSlots[randomIndex]
```

**场地选择策略：**
- 如果**指定**了场地号和时间：只匹配完全符合的场地
- 如果**未指定**场地号：返回所有可用场地
- 如果**未指定**时间：返回所有可用时间段
- 发现多个可用场地时：**随机选择一个**进行预约

## 可复用的现有接口

### 1. ✅ 场地信息获取
- `fetchVenueInfo(client, date, venueSiteId)`
- 返回指定日期的场馆可用时间段信息

### 2. ✅ 用户认证
- `NewClient(username, password)`
- 创建已认证的HTTP客户端

### 3. ✅ 预约执行
- `tyysReserveTask(taskInfo, isInstant)`
- 完整的预约流程（包含验证码处理）

### 4. ✅ 数据库操作
- GORM ORM 框架
- 已有的用户和账号表

### 5. ✅ 定时任务
- `robfig/cron/v3` 库
- 支持分钟级精度的定时任务

## 优势特性

### 1. 高可用性
- ✅ 服务器重启后自动恢复活跃任务
- ✅ 任务状态持久化到数据库
- ✅ 支持并发处理多个任务

### 2. 灵活配置
- ✅ 自定义扫描频率（1-60分钟）
- ✅ 可选的场地和时间过滤
- ✅ 支持任意已开放的预约日期

### 3. 可追溯性
- ✅ 详细的扫描日志
- ✅ 记录每次扫描的可用场地数
- ✅ 成功/失败状态追踪

### 4. 安全性
- ✅ JWT 身份验证
- ✅ 用户只能操作自己的任务
- ✅ 管理员权限控制

## 使用场景

### 场景1：抢热门时段
```
用户希望预约周五晚上19:00的场地，但该时段经常被秒抢或有人取消。
用户可以：
1. 配置两个预约账号（账号1主预约，账号2提供同伴码）
2. 创建捡漏任务，设置扫描间隔为5分钟
3. 指定venue_site_id和reservation_time为"19:00"
4. 不指定site_name（接受任何场地）
5. 系统每5分钟检查一次，发现有人取消立即用两个账号预约
```

### 场景2：等待特定场地
```
用户只想预约"1号场"，不接受其他场地。
用户可以：
1. 配置两个预约账号
2. 创建捡漏任务，指定site_name为"1号场"
3. 设置较长的扫描间隔（如30分钟）节省资源
4. 系统只在"1号场"可用时才使用两个账号进行预约
```

### 场景3：灵活预约
```
用户对时间和场地都不挑剔，只要有空位就预约。
用户可以：
1. 配置两个预约账号
2. 创建捡漏任务，不指定site_name和reservation_time
3. 系统扫描发现多个可用场地时，会随机选择一个
4. 随机选择可以避免多个捡漏任务竞争同一个场地
5. 提高预约成功率
```

## 技术实现细节

### 1. 并发安全
```go
var (
    bargainCronJobs = make(map[string]*cron.Cron)
    bargainMutex    sync.RWMutex  // 保护并发访问
)
```

### 2. 定时任务管理
```go
// 使用 cron 表达式
cronSpec := fmt.Sprintf("*/%d * * * *", scanInterval)
// 每个任务独立的 cron 实例
c := cron.New()
c.AddFunc(cronSpec, scanFunction)
c.Start()
```

### 3. 随机选择策略
```go
// 当发现多个可用场地时，随机选择一个
randomIndex := rand.Intn(len(availableSlots))
slot := availableSlots[randomIndex]
```
**优势：**
- 避免所有捡漏任务都抢同一个场地
- 分散预约压力，提高整体成功率
- 对于不指定场地/时间的用户更加灵活

### 4. 错误处理
- 登录失败：记录日志，等待下次扫描
- 网络错误：记录日志，等待下次扫描
- 预约失败：记录日志，等待下次扫描
- 任务不存在：停止调度器

### 5. 资源管理
- 每个任务有独立的 cron 实例
- 任务取消/完成时自动释放资源
- 限制日志最多返回100条

## 部署说明

### 1. 数据库迁移
系统启动时会自动创建新表：
- `bargain_task_dbs`
- `bargain_log_dbs`

### 2. 依赖包
需要添加的新依赖：
```bash
go get github.com/google/uuid
```

### 3. 启动流程
```go
func main() {
    InitDb()                       // 初始化数据库
    migrations.RunMigrations()     // 运行迁移
    restartActiveBargainTasks()    // 重启活跃任务
    router()                       // 启动路由
}
```

## 扩展建议

### 未来可能的功能扩展：
1. 支持更多账号（3个或更多）
2. 支持多个时间段选择
3. 支持多个场地选择
4. 添加通知功能（邮件/短信）
5. 添加任务优先级
6. 支持任务暂停/恢复
7. 添加任务执行时间限制
8. 支持按星期设置重复任务

## 总结

捡漏模式完全复用了现有的预约逻辑和API接口，只是在上层增加了定时扫描和自动触发的机制。

### 预约机制
- **一个预约任务 = 两个账号预约同一个场地**
  - 账号1：主预约账号
  - 账号2：提供同伴码（BuddyNum）
  - 与普通模式完全相同的预约流程

### 核心优势
- ✅ 无需前端人工监控
- ✅ 自动捕获取消的场地
- ✅ 灵活的过滤条件
- ✅ 完整的日志追踪
- ✅ 高可用性设计
- ✅ 完全复用现有预约逻辑
