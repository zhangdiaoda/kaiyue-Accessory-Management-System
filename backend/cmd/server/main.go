package main

import (
	"fmt"
	"log"
	"os"
	"warehouse/internal/handler"
	"warehouse/internal/middleware"
	"warehouse/internal/model"
	"warehouse/internal/notification"
	"warehouse/internal/notification/providers"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	// 连接数据库
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "warehouse:Warehouse@2026@tcp(localhost:3306)/warehouse?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	fmt.Println("✅ 数据库连接成功")

	// 自动迁移数据库（包含所有业务表）
	err = db.AutoMigrate(
		&model.User{},
		&model.Employee{},
		&model.PartCategory{},
		&model.Part{},
		&model.BorrowRecord{},
		&model.SysConfig{},
		&model.Announcement{},
		&model.InternalMessage{},
		&model.OldPartInventory{},
		&model.ScrapInventory{},
		&model.Permission{},
		&model.UserPermission{},
		&model.RolePermission{},
		&model.OperationLog{},
		&model.InboundRecord{},
		// 通知系统表
		&model.NotificationConfig{},
		&model.NotificationLog{},
		&model.UserNotificationSetting{},
		&model.WechatUserBinding{},
	)

	if err != nil {
		fmt.Printf("⚠️ 数据库部分迁移失败 (跳过已有冲突表): %v\n", err)
	} else {
		fmt.Println("📦 数据库迁移成功 (配置/公告/消息)")
	}

	// 初始化通知管理器
	notifManager := notification.NewManager(db)

	// 注册通知提供者
	providers.LoadAllProviders(db, notifManager)

	// 启动通知队列工作线程(3个)
	notifManager.StartWorker(3)
	fmt.Println("🔔 通知系统已启动")

	// 初始化通知集成助手
	notifIntegrator := handler.NewNotificationIntegrator(db, notifManager)

	// 启动通知定时任务调度器
	notifScheduler := handler.NewNotificationScheduler(notifIntegrator)
	notifScheduler.Start()

	// 初始化Gin
	router := gin.Default()

	// 使用CORS中间件
	router.Use(middleware.CORSMiddleware())

	// 创建Handler
	authHandler := handler.NewAuthHandler(db)
	partHandler := handler.NewPartHandler(db, notifIntegrator)
	employeeHandler := handler.NewEmployeeHandler(db)
	borrowHandler := handler.NewBorrowHandler(db, notifIntegrator)
	dashboardHandler := handler.NewDashboardHandler(db)
	reportHandler := handler.NewReportHandler(db, notifIntegrator)
	categoryHandler := handler.NewCategoryHandler(db)
	dingTalkHandler := handler.NewDingTalkHandler()
	systemHandler := handler.NewSystemHandler(db)
	messageHandler := handler.NewMessageHandler(db)
	notificationHandler := handler.NewNotificationHandler(db, notifManager, notifScheduler)
	permissionHandler := handler.NewPermissionHandler(db)
	auditHandler := handler.NewAuditHandler(db)

	// 公开路由
	api := router.Group("/api")
	{
		api.GET("/branding", systemHandler.GetBrandingConfig)

		// 认证路由管理
		auth := api.Group("/auth")
		{
			// 无需Token的公有接口
			auth.POST("/login", authHandler.Login)

			// 需要Token的私有接口
			authPrivate := auth.Group("")
			authPrivate.Use(middleware.AuthMiddleware(db))
			{
				authPrivate.GET("/userinfo", authHandler.GetUserInfo)
				authPrivate.GET("/users", authHandler.GetUserList)
				authPrivate.POST("/users", authHandler.CreateUser)
				authPrivate.PUT("/users/:id", authHandler.UpdateUser)
				authPrivate.DELETE("/users/:id", authHandler.DeleteUser)
				authPrivate.PUT("/users/:id/reset-password", authHandler.ResetPassword)
				authPrivate.PUT("/profile", authHandler.UpdateProfile)
				authPrivate.POST("/logout", authHandler.Logout)
			}
		}

		// 其他需要认证的业务路由
		authRequired := api.Group("")
		authRequired.Use(middleware.AuthMiddleware(db))
		// 启用审计日志（仓库管理员操作记录）
		authRequired.Use(middleware.AuditMiddleware(db))
		{
			// 系统配置与公告
			authRequired.GET("/system/config", systemHandler.GetConfig)
			authRequired.PUT("/system/config", systemHandler.UpdateConfig)

			announcements := authRequired.Group("/announcements")
			{
				announcements.GET("", systemHandler.GetAnnouncements)
				announcements.POST("", systemHandler.CreateAnnouncement)
				announcements.PUT("/:id", systemHandler.UpdateAnnouncement)
				announcements.DELETE("/:id", systemHandler.DeleteAnnouncement)
			}

			// 站内信
			messages := authRequired.Group("/messages")
			{
				messages.GET("", messageHandler.GetMessages)
				messages.POST("", messageHandler.SendMessage)
				messages.PUT("/:id/read", messageHandler.MarkAsRead)
				messages.GET("/unread-count", messageHandler.GetUnreadCount)
			}

			// 仪表盘
			authRequired.GET("/dashboard/stats", dashboardHandler.GetDashboardStats)

			// 配件管理
			parts := authRequired.Group("/parts")
			{
				parts.GET("", partHandler.GetPartList)
				parts.POST("", partHandler.CreatePart)
				parts.PUT("/:id", partHandler.UpdatePart)
				parts.DELETE("/:id", partHandler.DeletePart)
				parts.GET("/low-stock", partHandler.GetLowStockParts)
				parts.GET("/template", partHandler.DownloadTemplate)
				parts.POST("/import", partHandler.ImportParts)
				parts.GET("/export", partHandler.ExportParts)
			}

			// 配件分类
			categories := authRequired.Group("/categories")
			{
				categories.GET("", categoryHandler.GetCategoryList)
				categories.POST("", categoryHandler.CreateCategory)
				categories.PUT("/:id", categoryHandler.UpdateCategory)
				categories.DELETE("/:id", categoryHandler.DeleteCategory)
			}

			// 员工管理
			employees := authRequired.Group("/employees")
			{
				employees.GET("", employeeHandler.GetEmployeeList)
				employees.GET("/all", employeeHandler.GetAllEmployees) // 获取所有在职员工
				employees.POST("", employeeHandler.CreateEmployee)
				employees.PUT("/:id", employeeHandler.UpdateEmployee)
				employees.DELETE("/:id", employeeHandler.DeleteEmployee)
				employees.GET("/template", employeeHandler.DownloadTemplate)
				employees.POST("/import", employeeHandler.ImportEmployees)
				employees.GET("/export", employeeHandler.ExportEmployees)
			}

			// 领用管理
			borrows := authRequired.Group("/borrows")
			{
				borrows.GET("", borrowHandler.GetBorrowRecordList)
				borrows.POST("", borrowHandler.CreateBorrowRecord)
				borrows.POST("/:id/return", borrowHandler.ReturnBorrowRecord)
				borrows.GET("/check-unreturned", borrowHandler.CheckUnreturned)
				borrows.GET("/old-inventory", borrowHandler.GetOldPartInventory)
				borrows.GET("/scrap-inventory", borrowHandler.GetScrapInventory)
			}

			// 报表统计
			reports := authRequired.Group("/reports")
			{
				reports.GET("/by-part", reportHandler.GetPartReport)
				reports.GET("/by-employee", reportHandler.GetEmployeeReport)
				reports.GET("/by-department", reportHandler.GetDepartmentReport)
				reports.GET("/detailed", reportHandler.GetDetailedReport)
				reports.POST("/push", reportHandler.PushReport)
				reports.GET("/download/:filename", reportHandler.DownloadReport)
			}

			// 钉钉推送
			dingtalk := authRequired.Group("/dingtalk")
			{
				dingtalk.POST("/send", dingTalkHandler.SendReport)
				dingtalk.POST("/test", dingTalkHandler.TestWebhook)
			}

			// 补货管理
			inboundHandler := handler.NewInboundHandler(db, notifIntegrator)
			inbound := authRequired.Group("/inbound")
			{
				inbound.POST("", inboundHandler.RestockPart)
				inbound.GET("", inboundHandler.GetInboundList)
			}

			// 数据库备份
			backupHandler := handler.NewBackupHandler(db)
			backups := authRequired.Group("/system/backups")
			{
				backups.GET("/config", backupHandler.GetConfig)
				backups.PUT("/config", backupHandler.UpdateConfig)
				backups.POST("/run", backupHandler.RunBackup)
				backups.GET("", backupHandler.GetBackups)
				backups.DELETE("/", backupHandler.DeleteBackup)
				backups.POST("/restore", backupHandler.RestoreBackup)
				backups.GET("/download", backupHandler.DownloadBackup)
			}
			// 通知系统
			notifications := authRequired.Group("/notifications")
			{
				// 配置管理
				notifications.GET("/configs", notificationHandler.GetConfigs)
				notifications.GET("/config/:type", notificationHandler.GetConfig)
				notifications.POST("/config", notificationHandler.UpdateConfig)

				// 手动测试与运行
				notifications.POST("/test/dingtalk", notificationHandler.TestDingTalk)
				notifications.POST("/test/wechat", notificationHandler.TestWechat)
				notifications.POST("/run/daily-report", notificationHandler.RunDailyReportNow)
				notifications.POST("/run/overdue-check", notificationHandler.RunOverdueCheckNow)

				// 日志查询
				notifications.GET("/logs", notificationHandler.GetLogs)
				notifications.GET("/stats", notificationHandler.GetStats)

				// 手动发送
				notifications.POST("/send", notificationHandler.SendManualNotification)

				// 用户设置
				notifications.GET("/user/settings", notificationHandler.GetUserSettings)
				notifications.POST("/user/setting", notificationHandler.UpdateUserSetting)

				// 调度配置
				notifications.GET("/schedules", notificationHandler.GetScheduleConfigs)
				notifications.POST("/schedule", notificationHandler.UpdateScheduleConfig)

				// 微信绑定
				notifications.GET("/wechat/binding", notificationHandler.GetWechatBinding)
				notifications.POST("/wechat/bind", notificationHandler.BindWechatUser)
				notifications.POST("/wechat/subscribe", notificationHandler.UpdateSubscribeScenes)
			}

			// 权限管理（仅超级管理员）
			permissions := authRequired.Group("/permissions")
			{
				permissions.GET("", permissionHandler.GetAllPermissions)
				permissions.GET("/users/:id", permissionHandler.GetUserPermissions)
				permissions.PUT("/users/:id", permissionHandler.SetUserPermissions)
				permissions.GET("/roles/:role", permissionHandler.GetRolePermissions)
				permissions.PUT("/roles/:role", permissionHandler.SetRolePermissions)
			}

			// 操作日志（仅超级管理员）
			audit := authRequired.Group("/operation-logs")
			{
				audit.GET("", auditHandler.GetOperationLogs)
				audit.GET("/:id", auditHandler.GetLogDetail)
				audit.GET("/stats", auditHandler.GetOperationStats)
			}
		}
	}

	// 启动服务器
	fmt.Println("======================================")
	fmt.Println("🚀 仓储管理系统后端启动成功！")
	fmt.Println("📍 API地址: http://localhost:8080")
	fmt.Println("======================================")

	if err := router.Run(":8080"); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
