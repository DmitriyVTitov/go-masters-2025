package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

type Cfg struct {
	Port      string `mapstructure:"port"`
	DBConnStr string `mapstructure:"db_conn_str"`
}

var (
	once     sync.Once
	instance *Cfg
)

// Load загружает конфигурацию из файла и переменных окружения
func Load(path string) (*Cfg, error) {
	var err error
	once.Do(func() {
		// Устанавливаем имя конфигурационного файла (без расширения)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")

		// Устанавливаем путь для поиска конфигурационного файла
		viper.AddConfigPath(path)

		// Чтение конфигурационного файла
		if err = viper.ReadInConfig(); err != nil {
			log.Printf("Error reading config file, %s", err)
		}

		// Автоматическое связывание переменных окружения с конфигом
		viper.AutomaticEnv()

		// Устанавливаем префикс для переменных окружения (опционально)
		viper.SetEnvPrefix("MYAPP")

		// Устанавливаем альтернативные имена для переменных окружения
		viper.BindEnv("port", "MYAPP_PORT")

		// Устанавливаем значения по умолчанию
		viper.SetDefault("port", 10000)

		instance = &Cfg{}
		if err = viper.Unmarshal(instance); err != nil {
			log.Printf("Unable to decode into config struct, %v", err)
		}
	})

	return instance, err
}
