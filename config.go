package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// Dashboard 配置结构
type Dashboard struct {
	Name         string
	DashboardUID string
	Slug         string
	OrgID        int `mapstructure:"orgID"`
}

// Config 配置文件结构
type Config struct {
	Token        string
	BaseURL      string
	OrgID        int
	CronTime     string
	DevMode      bool
	Dashboards   []Dashboard
	WeChatBotKey string
}

// LoadConfig 加载 config.yaml
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // 文件名（不带扩展名）
	viper.SetConfigType("yaml")   // 文件类型
	viper.AddConfigPath(".")      // 当前目录

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &cfg, nil
}
