package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// AppConfig 应用程序配置结构
type AppConfig struct {
	Gin struct {
		Port             string   `yaml:"port"`
		CorsAllowOrigins []string `yaml:"corsAllowOrigins"`
		Domain           string   `yaml:"domain"`
	} `yaml:"gin"`

	MySQL struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"mysql"`

	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
		Database int    `yaml:"database"`
	} `yaml:"redis"`

	JWT struct {
		SecretKey  string `yaml:"secretKey"`
		ExpireTime int    `yaml:"expireTime"`
	} `yaml:"jwt"`

	Mail struct {
		SenderName string `yaml:"senderName"`
		SmtpHost   string `yaml:"smtpHost"`
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
		SmtpPort   int    `yaml:"smtpPort"`
	} `yaml:"mail"`

	AzureTTS struct {
		BaseURL         string `yaml:"baseURL"`
		SSMLTemplate    string `yaml:"ssmlTemplate"`
		SubscriptionKey string `yaml:"subscriptionKey"`
		Region          string `yaml:"region"`
	} `yaml:"azureTTS"`

	AzureSpeech struct {
		SubscriptionKey string `yaml:"subscriptionKey"`
		Region          string `yaml:"region"`
		BaseURL         string `yaml:"baseURL"`
	} `yaml:"azureSpeech"`

	OpenAI struct {
		Key               string  `yaml:"key"`
		BaseURL           string  `yaml:"baseURL"`
		GlobalTemperature float64 `yaml:"globalTemperature"`
		ThinkModelName    string  `yaml:"thinkModelName"`
		AgentModelName    string  `yaml:"agentModelName"`
		UseModelName      string  `yaml:"useModelName"`
	} `yaml:"openAI"`

	AnythingLLM struct {
		BaseURL string `yaml:"baseURL"`
		APIKey  string `yaml:"apiKey"`
	} `yaml:"anythingLLM"`
}

// Config 全局配置变量
var Config AppConfig

// 初始化函数，在包被导入时自动执行
func ReadConfig() {
	// 获取可执行文件所在目录
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("获取可执行文件路径失败: %v", err)
	}

	execDir := filepath.Dir(execPath)
	configPath := filepath.Join(execDir, "app.yml")

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 尝试在当前工作目录查找
		workDir, _ := os.Getwd()
		configPath = filepath.Join(workDir, "app.yml")
	}

	// 读取配置文件
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 解析YAML
	if err := yaml.Unmarshal(data, &Config); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	fmt.Println("配置文件加载成功:", configPath)
}
