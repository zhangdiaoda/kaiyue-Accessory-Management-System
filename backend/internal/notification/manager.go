package notification

import (
	"encoding/json"
	"log"
	"sync"
	"time"
	"warehouse/internal/model"

	"gorm.io/gorm"
)

// Manager notification manager
type Manager struct {
	providers map[ProviderType]NotificationProvider
	db        *gorm.DB
	queue     chan QueueTask
	wg        sync.WaitGroup
	stopChan  chan bool
	lock      sync.RWMutex
}

// NewManager create notification manager
func NewManager(db *gorm.DB) *Manager {
	return &Manager{
		providers: make(map[ProviderType]NotificationProvider),
		db:        db,
		queue:     make(chan QueueTask, 1000),
		stopChan:  make(chan bool),
	}
}

// RegisterProvider register notification provider
func (m *Manager) RegisterProvider(provider NotificationProvider) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.providers[provider.GetType()] = provider
}

// ClearProviders clear all registered providers
func (m *Manager) ClearProviders() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.providers = make(map[ProviderType]NotificationProvider)
}

// SendNotification send notification (sync)
func (m *Manager) SendNotification(notif *Notification, providers []ProviderType) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	var lastErr error
	for _, providerType := range providers {
		provider, ok := m.providers[providerType]
		if !ok {
			log.Printf("Provider not found: %s", providerType)
			continue
		}

		logEntry := m.createLogEntry(notif, providerType)
		m.db.Create(logEntry)

		err := provider.Send(notif)
		if err != nil {
			logEntry.Status = "failed"
			logEntry.ErrorMsg = err.Error()
			lastErr = err
		} else {
			logEntry.Status = "success"
			now := time.Now()
			logEntry.SentAt = &now
		}
		m.db.Save(logEntry)
	}
	return lastErr
}

// SendNotificationAsync send notification (async)
func (m *Manager) SendNotificationAsync(notif *Notification, providers []ProviderType) {
	for _, providerType := range providers {
		task := QueueTask{
			Provider:     providerType,
			Notification: notif,
			RetryCount:   0,
			ScheduledAt:  time.Now(),
		}
		m.queue <- task
	}
}

// StartWorker start worker threads
func (m *Manager) StartWorker(workerCount int) {
	for i := 0; i < workerCount; i++ {
		m.wg.Add(1)
		go m.worker()
	}
	log.Printf("Notification system: started %d workers", workerCount)
}

// worker worker thread
func (m *Manager) worker() {
	defer m.wg.Done()
	for {
		select {
		case task := <-m.queue:
			m.processTask(task)
		case <-m.stopChan:
			return
		}
	}
}

// processTask process task
func (m *Manager) processTask(task QueueTask) {
	m.lock.RLock()
	provider, ok := m.providers[task.Provider]
	m.lock.RUnlock()

	if !ok {
		log.Printf("Provider not found: %s", task.Provider)
		return
	}

	logEntry := m.createLogEntry(task.Notification, task.Provider)
	logEntry.RetryCount = task.RetryCount
	m.db.Create(logEntry)

	err := provider.Send(task.Notification)
	if err != nil {
		logEntry.Status = "failed"
		logEntry.ErrorMsg = err.Error()
		log.Printf("Send failed [%s]: %v", task.Provider, err)

		if task.RetryCount < 3 {
			task.RetryCount++
			task.ScheduledAt = time.Now().Add(5 * time.Minute)

			go func() {
				time.Sleep(5 * time.Minute)
				m.queue <- task
			}()
			log.Printf("Will retry in 5 minutes (attempt %d)", task.RetryCount)
		}
	} else {
		logEntry.Status = "success"
		now := time.Now()
		logEntry.SentAt = &now
	}

	m.db.Save(logEntry)
}

// createLogEntry create log entry
func (m *Manager) createLogEntry(notif *Notification, providerType ProviderType) *model.NotificationLog {
	var receiverID *uint
	if notif.ReceiverID > 0 {
		receiverID = &notif.ReceiverID
	}

	return &model.NotificationLog{
		ProviderType:   string(providerType),
		SceneType:      string(notif.Scene),
		Title:          notif.Title,
		Content:        notif.Content,
		ReceiverID:     receiverID,
		ReceiverOpenID: notif.OpenID,
		Status:         "pending",
		RetryCount:     0,
	}
}

// GetProviderConfig get provider config
func (m *Manager) GetProviderConfig(providerType ProviderType) (interface{}, error) {
	var config model.NotificationConfig
	err := m.db.Where("provider_type = ? AND is_enabled = ?", string(providerType), true).
		First(&config).Error
	if err != nil {
		return nil, err
	}

	var configData interface{}
	if err := json.Unmarshal([]byte(config.ConfigData), &configData); err != nil {
		return nil, err
	}

	return configData, nil
}

// UpdateProviderConfig update provider config
func (m *Manager) UpdateProviderConfig(providerType ProviderType, configName string, configData interface{}) error {
	jsonData, err := json.Marshal(configData)
	if err != nil {
		return err
	}

	var config model.NotificationConfig
	err = m.db.Where("provider_type = ?", string(providerType)).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		config = model.NotificationConfig{
			ConfigName:   configName,
			ProviderType: string(providerType),
			IsEnabled:    true,
			ConfigData:   string(jsonData),
		}
		return m.db.Create(&config).Error
	} else if err != nil {
		return err
	}

	config.ConfigData = string(jsonData)
	config.ConfigName = configName
	return m.db.Save(&config).Error
}

// Stop stop notification manager
func (m *Manager) Stop() {
	close(m.stopChan)
	m.wg.Wait()
	close(m.queue)
	log.Println("Notification system: stopped")
}

// GetRegisteredProviders get all registered provider types
func (m *Manager) GetRegisteredProviders() []ProviderType {
	m.lock.RLock()
	defer m.lock.RUnlock()
	types := make([]ProviderType, 0, len(m.providers))
	for t := range m.providers {
		types = append(types, t)
	}
	return types
}

// GetStats get statistics
func (m *Manager) GetStats() map[string]interface{} {
	var stats struct {
		TotalCount   int64
		SuccessCount int64
		FailedCount  int64
		PendingCount int64
	}

	m.db.Model(&model.NotificationLog{}).Count(&stats.TotalCount)
	m.db.Model(&model.NotificationLog{}).Where("status = ?", "success").Count(&stats.SuccessCount)
	m.db.Model(&model.NotificationLog{}).Where("status = ?", "failed").Count(&stats.FailedCount)
	m.db.Model(&model.NotificationLog{}).Where("status = ?", "pending").Count(&stats.PendingCount)

	return map[string]interface{}{
		"total":      stats.TotalCount,
		"success":    stats.SuccessCount,
		"failed":     stats.FailedCount,
		"pending":    stats.PendingCount,
		"queue_size": len(m.queue),
	}
}

// GetProvidersByScene 根据业务场景获取已启用的提供者列表
func (m *Manager) GetProvidersByScene(scene SceneType) []ProviderType {
	var configs []model.NotificationConfig
	if err := m.db.Where("is_enabled = ?", true).Find(&configs).Error; err != nil {
		return nil
	}

	result := make([]ProviderType, 0)
	for _, cfg := range configs {
		// 如果订阅场景为空，默认为“全订阅”（兼容旧版本）或者可以设为全场景开启
		if cfg.SubscribedScenes == "" {
			result = append(result, ProviderType(cfg.ProviderType))
			continue
		}

		// 检查场景是否在订阅列表中
		var scenes []string
		if err := json.Unmarshal([]byte(cfg.SubscribedScenes), &scenes); err != nil {
			// 如果解析失败，保守起见包含该渠道
			result = append(result, ProviderType(cfg.ProviderType))
			continue
		}

		for _, s := range scenes {
			if s == string(scene) {
				result = append(result, ProviderType(cfg.ProviderType))
				break
			}
		}
	}
	return result
}
