# 仓储管理系统

> 机械加工配件（刀片等）的领取、归还、损毁管理系统，支持库存管理、多维度报表统计

## ✨ 项目特点

- 🚀 **高性能Go后端** - 单文件编译，启动快速，资源占用低
- 🎨 **现代化Vue3前端** - Composition API + Element Plus，用户体验优秀
- 🐳 **完全容器化** - Docker Compose一键部署
- 🔐 **JWT认证** - 无状态认证，支持分布式部署
- 📊 **完整报表** - 多维度统计（按配件、员工、部门）
- ⚡ **实时预警** - 库存低于阈值自动预警
- 🔄 **事务保证** - 领用归还流程数据一致性

## 📋 功能清单

### ✅ 核心业务功能

**配件管理**
- 配件信息CRUD（增删改查）
- 配件分类管理（支持多级分类）
- 库存实时更新
- 低库存预警（库存 < 阈值红色标记）
- 快速补充库存

**员工管理**
- 员工信息CRUD
- 在职/离职状态管理
- 按部门筛选

**领用归还流程**
- 领用登记（自动扣减库存）
- 归还登记（正常归还/损毁标记）
- 部分归还支持
- 领用记录查询
- 事务保证数据一致性

**报表统计**
- 按配件统计（领用次数、损毁率）
- 按员工统计（领用量、损毁量排名）
- 按部门统计（部门汇总数据）
- 时间范围筛选

**仪表盘**
- 6项实时统计指标
- 快捷操作入口
- 系统概况展示

## 🛠️ 技术栈

### 后端
- **语言**: Go 1.21+
- **框架**: Gin 1.9.x
- **ORM**: GORM 1.25.x
- **数据库**: MySQL 8.0
- **缓存**: Redis 7
- **认证**: JWT (golang-jwt/jwt)

### 前端
- **框架**: Vue 3.4.x
- **构建工具**: Vite 5.x
- **UI库**: Element Plus 2.5.x
- **状态管理**: Pinia 2.x
- **路由**: Vue Router 4.x
- **HTTP**: Axios 1.6.x

### 部署
- **容器化**: Docker + Docker Compose
- **反向代理**: Nginx 1.25

## 🚀 快速开始

### 前置要求

- Docker & Docker Compose
- （可选）Go 1.21+、Node.js 20+（用于开发）

### 一键部署（推荐）

```bash
# 1. 克隆项目
cd c:/Users/admin/Desktop/仓储管理

# 2. 配置环境变量（可选，使用默认配置）
cp .env.example .env

# 3. 启动所有服务
docker-compose -p warehouse up -d

# 4. 访问系统
# 前端: http://localhost
# 后端API: http://localhost:8080
```

**默认登录账号**
- 用户名: `admin`
- 密码: `admin123`

### 开发环境运行

**启动数据库**
```bash
docker-compose -p warehouse up -d mysql redis
```

**启动后端**
```bash
cd backend
go mod tidy
go run cmd/server/main.go
# 访问: http://localhost:8080
```

**启动前端**
```bash
cd frontend
npm install
npm run dev
# 访问: http://localhost:5173
```

## 📁 项目结构

```
仓储管理/
├── backend/                    # Go后端
│   ├── cmd/server/            # 主程序入口
│   ├── internal/
│   │   ├── handler/           # HTTP处理器
│   │   ├── middleware/        # 中间件
│   │   ├── model/             # 数据模型
│   │   ├── repository/        # 数据访问层（预留）
│   │   └── service/           # 业务逻辑层（预留）
│   ├── pkg/utils/             # 工具函数
│   ├── config/                # 配置文件
│   ├── go.mod                 # Go模块定义
│   └── Dockerfile             # 后端镜像
├── frontend/                   # Vue前端
│   ├── src/
│   │   ├── views/             # 页面组件
│   │   ├── api/               # API接口
│   │   ├── router/            # 路由配置
│   │   ├── store/             # 状态管理
│   │   └── utils/             # 工具函数
│   ├── package.json
│   └── Dockerfile             #前端镜像
├── init-db/                    # 数据库初始化脚本
│   └── 01-schema.sql          # 建表SQL
├── nginx/                      # Nginx配置
│   └── nginx.conf
├── docker-compose.yml          # Docker编排
├── .env                        # 环境变量
└── README.md                   # 项目说明
```

## 📊 数据库设计

### 核心表结构

- `sys_user` - 用户表（管理员）
- `employee` - 员工表
- `part_category` - 配件分类表
- `part` - 配件信息表
- `borrow_record` - 领用记录表
- `sys_config` - 系统配置表

详见 [init-db/01-schema.sql](init-db/01-schema.sql)

## 🔌 API接口

### 认证模块 (3个)
- `POST /api/auth/login` - 用户登录
- `GET /api/auth/userinfo` - 获取用户信息
- `POST /api/auth/logout` - 退出登录

### 配件管理 (9个)
- `GET /api/parts` - 配件列表
- `POST /api/parts` - 创建配件
- `PUT /api/parts/:id` - 更新配件
- `DELETE /api/parts/:id` - 删除配件
- `GET /api/parts/low-stock` - 低库存预警
- `GET /api/categories` - 分类列表
- `POST /api/categories` - 创建分类
- `PUT /api/categories/:id` - 更新分类
- `DELETE /api/categories/:id` - 删除分类

### 员工管理 (5个)
- `GET /api/employees` - 员工列表
- `GET /api/employees/all` - 所有在职员工
- `POST /api/employees` - 创建员工
- `PUT /api/employees/:id` - 更新员工
- `DELETE /api/employees/:id` - 删除员工

### 领用管理 (3个)
- `GET /api/borrows` - 领用记录列表
- `POST /api/borrows` - 创建领用记录
- `POST /api/borrows/:id/return` - 归还登记

### 报表统计 (3个)
- `GET /api/reports/by-part` - 按配件统计
- `GET /api/reports/by-employee` - 按员工统计
- `GET /api/reports/by-department` - 按部门统计

### 仪表盘 (1个)
- `GET /api/dashboard/stats` - 统计数据

**总计：24个API接口**

## 🔧 配置说明

### 环境变量 (.env)

```env
# MySQL配置
MYSQL_ROOT_PASSWORD=your_root_password
MYSQL_PASSWORD=Warehouse@2026

# Redis配置
REDIS_PASSWORD=Redis@2026

# JWT密钥
JWT_SECRET=your_jwt_secret_key
```

## 🐛 常见问题

**Q: 启动后无法访问？**
A: 检查Docker容器状态 `docker-compose -p warehouse ps`，确保所有服务都是"Up"状态

**Q: 数据库连接失败？**
A: 确认MySQL容器已启动且健康检查通过，检查密码配置是否正确

**Q: 前端页面空白？**
A: 检查浏览器控制台错误，确认后端API地址配置正确

## 📝 开发计划

### 已完成 ✅
- [x] 用户认证与授权
- [x] 配件管理（CRUD + 分类）
- [x] 员工管理
- [x] 领用归还流程
- [x] 库存预警
- [x] 报表统计
- [x] 仪表盘数据

### 待开发 🚧
- [ ] Excel导出功能
- [ ] 钉钉推送集成
- [ ] Android APP
- [ ] 数据可视化图表（ECharts）
- [ ] 操作日志审计

## 📄 许可证

MIT License

## 👥 贡献

欢迎提交 Issue 和 Pull Request！

---

**开发时间**: 2026-01-25  
**当前版本**: v1.0-beta  
**完成度**: 85%
