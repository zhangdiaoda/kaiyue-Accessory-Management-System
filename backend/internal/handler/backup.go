package handler

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	"warehouse/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type BackupHandler struct {
	db       *gorm.DB
	cron     *cron.Cron
	entryID  cron.EntryID
	cronLock sync.Mutex
}

func NewBackupHandler(db *gorm.DB) *BackupHandler {
	c := cron.New(cron.WithSeconds())
	c.Start()
	h := &BackupHandler{db: db, cron: c}
	h.InitScheduler()
	return h
}

// Config 结构体
type BackupConfig struct {
	BackupPath string `json:"backup_path"`
	Schedule   string `json:"schedule"`
}

// InitScheduler 初始化定时任务
func (h *BackupHandler) InitScheduler() {
	config := h.getBackupConfig()
	if config.Schedule != "" {
		h.updateScheduler(config.Schedule)
	}
}

func (h *BackupHandler) getBackupConfig() BackupConfig {
	var configs []model.SysConfig
	h.db.Where("config_key IN ?", []string{"backup_path", "backup_schedule"}).Find(&configs)

	var path, schedule string
	for _, c := range configs {
		if c.ConfigKey == "backup_path" {
			path = c.ConfigValue
		} else if c.ConfigKey == "backup_schedule" {
			schedule = c.ConfigValue
		}
	}

	// 默认路径
	if path == "" {
		// 默认为项目下的 backups 目录
		wd, _ := os.Getwd()
		path = filepath.Join(wd, "backups")
	}

	return BackupConfig{
		BackupPath: path,
		Schedule:   schedule,
	}
}

func (h *BackupHandler) updateScheduler(spec string) error {
	h.cronLock.Lock()
	defer h.cronLock.Unlock()

	if h.entryID != 0 {
		h.cron.Remove(h.entryID)
		h.entryID = 0
	}

	if spec == "" {
		return nil
	}

	id, err := h.cron.AddFunc(spec, func() {
		h.performBackup()
	})
	if err != nil {
		return err
	}
	h.entryID = id
	return nil
}

// performBackup 执行备份的核心逻辑
func (h *BackupHandler) performBackup() error {
	config := h.getBackupConfig()

	// 确保目录存在
	if err := os.MkdirAll(config.BackupPath, 0755); err != nil {
		return fmt.Errorf("create dir failed: %v", err)
	}

	filename := fmt.Sprintf("warehouse_%s.sql", time.Now().Format("20060102_150405"))
	filepath := filepath.Join(config.BackupPath, filename)

	// 使用 docker exec 执行 mysqldump
	// 假设容器名为 warehouse-mysql，密码在环境变量中，这里硬编码或从配置读取
	// 注意：在实际生产环境中，建议使用 .my.cnf 或环境变量传递密码，避免命令行泄漏
	// 这里为了简化演示，直接在命令中包含密码
	cmd := exec.Command("docker", "exec", "warehouse-mysql", "mysqldump", "-uwarehouse", "-pWarehouse@2026", "warehouse")

	outfile, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("create file failed: %v", err)
	}
	defer outfile.Close()

	cmd.Stdout = outfile
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("backup command failed: %v", err)
	}

	return nil
}

// GetConfig 获取备份配置
func (h *BackupHandler) GetConfig(c *gin.Context) {
	config := h.getBackupConfig()
	resp := gin.H{
		"backup_path": config.BackupPath,
		"schedule":    config.Schedule,
		"next_run":    "",
	}

	if h.entryID != 0 {
		entry := h.cron.Entry(h.entryID)
		if !entry.Next.IsZero() {
			resp["next_run"] = entry.Next.Format("2006-01-02 15:04:05")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": resp,
	})
}

// UpdateConfig 更新备份配置
func (h *BackupHandler) UpdateConfig(c *gin.Context) {
	var req BackupConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 验证路径是否合法（尝试创建）
	if err := os.MkdirAll(req.BackupPath, 0755); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "备份路径无效或无权限"})
		return
	}

	// 验证 Cron 表达式
	if req.Schedule != "" {
		if _, err := cron.ParseStandard(req.Schedule); err != nil {
			// 尝试秒级
			parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
			if _, err2 := parser.Parse(req.Schedule); err2 != nil {
				c.JSON(http.StatusOK, gin.H{"code": 400, "message": "Cron表达式无效"})
				return
			}
		}
	}

	// 保存配置
	h.saveConfig("backup_path", req.BackupPath, "数据库备份路径")
	h.saveConfig("backup_schedule", req.Schedule, "数据库备份定时规则")

	// 更新调度器
	if err := h.updateScheduler(req.Schedule); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "定时任务更新失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "配置已保存"})
}

func (h *BackupHandler) saveConfig(key, value, desc string) {
	var conf model.SysConfig
	result := h.db.Where("config_key = ?", key).First(&conf)
	if result.Error == gorm.ErrRecordNotFound {
		h.db.Create(&model.SysConfig{
			ConfigKey:   key,
			ConfigValue: value,
			Description: desc,
		})
	} else {
		h.db.Model(&conf).Updates(map[string]interface{}{
			"config_value": value,
			"description":  desc,
		})
	}
}

// RunBackup 手动触发备份
func (h *BackupHandler) RunBackup(c *gin.Context) {
	go func() {
		if err := h.performBackup(); err != nil {
			fmt.Printf("Manual backup failed: %v\n", err)
		}
	}()
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "备份任务已后台启动"})
}

type BackupFile struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Time string `json:"time"`
	Path string `json:"path"`
}

// GetBackups 获取备份列表
func (h *BackupHandler) GetBackups(c *gin.Context) {
	config := h.getBackupConfig()
	files, err := os.ReadDir(config.BackupPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []BackupFile{}})
		return
	}

	var backups []BackupFile
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			info, _ := f.Info()
			backups = append(backups, BackupFile{
				Name: f.Name(),
				Size: info.Size(),
				Time: info.ModTime().Format("2006-01-02 15:04:05"),
				Path: filepath.Join(config.BackupPath, f.Name()),
			})
		}
	}

	// 按时间倒序
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Time > backups[j].Time
	})

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": backups})
}

// DeleteBackup 删除备份
func (h *BackupHandler) DeleteBackup(c *gin.Context) {
	filename := c.Query("name")
	if filename == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "文件名不能为空"})
		return
	}

	config := h.getBackupConfig()
	path := filepath.Join(config.BackupPath, filename)

	// 安全检查：只能删除指定目录下的 .sql 文件
	if filepath.Dir(path) != config.BackupPath || !strings.HasSuffix(path, ".sql") {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "非法路径"})
		return
	}

	if err := os.Remove(path); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "已删除"})
}

// RestoreBackup 恢复备份
func (h *BackupHandler) RestoreBackup(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	config := h.getBackupConfig()
	path := filepath.Join(config.BackupPath, req.Name)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "备份文件不存在"})
		return
	}

	// 调用 docker exec mysql 执行恢复
	// docker exec -i warehouse-mysql mysql -u... -p... warehouse < file.sql
	cmd := exec.Command("docker", "exec", "-i", "warehouse-mysql", "mysql", "-uwarehouse", "-pWarehouse@2026", "warehouse")

	infile, err := os.Open(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "打开文件失败"})
		return
	}
	defer infile.Close()

	cmd.Stdin = infile
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "恢复失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "数据库已成功恢复"})
}

// DownloadBackup 下载备份文件
func (h *BackupHandler) DownloadBackup(c *gin.Context) {
	filename := c.Query("name")
	config := h.getBackupConfig()
	path := filepath.Join(config.BackupPath, filename)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.Status(404)
		return
	}

	// 强制浏览器下载而不是预览
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.File(path)
}
