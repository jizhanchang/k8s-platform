package settings

import (
	"log"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

var Conf = new(Config)

type Config struct {
	*AppConfig        `mapstructure:"app"`
	*LogConfig        `mapstructure:"log"`
	*MySQLConfig      `mapstructure:"mysql"`
	*RedisConfig      `mapstructure:"redis"`
	*KubernetesConfig `mapstructure:"kubernetes"`
}

type KubernetesConfig struct {
	KubeConfig *string `mapstructure:"kube_config"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Mode    string `mapstructure:"mode"`
	Version string `mapstructure:"version"`
	Port    int    `mapstructure:"port"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	Host        string `mapstructure:"host"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure"password"`
	DBName      string `mapstructure:"dbname"`
	Port        int    `mapstructure:"port"`
	MaxOpenConn int    `mapstructure:"max_open_conns"`
	MaxIdleConn int    `mapstructure:"max_idle_conns"`
	Enable      bool   `mapstructure:"enable"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	Db       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
	Enable   bool   `mapstructure:"enable"`
}

func Init(cfg *string) (err error) {

	viper.SetConfigName(*cfg)   // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml") // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(".")    // 还可以在工作目录中查找配置
	err = viper.ReadInConfig()  // 查找并读取配置文件
	if err != nil {             // 处理读取配置文件的错误
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return
		}
	}
	if err := viper.Unmarshal(Conf); err != nil {
		return err
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("重新加载配置文件...")
		if err := viper.Unmarshal(Conf); err != nil {
			return
		}
	})
	return
}
