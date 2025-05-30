// Config содержит конфигурацию приложения
type Config struct {
	// Server содержит настройки сервера
	Server struct {
		// Port порт, на котором будет запущен сервер
		Port string `env:"SERVER_PORT" envDefault:":3002"`
		// Host хост, на котором будет запущен сервер
		Host string `env:"SERVER_HOST" envDefault:"localhost"`
	}

	// Database содержит настройки базы данных
	Database struct {
		// Host хост базы данных
		Host string `env:"DB_HOST" envDefault:"localhost"`
		// Port порт базы данных
		Port string `env:"DB_PORT" envDefault:"5432"`
		// User пользователь базы данных
		User string `env:"DB_USER" envDefault:"postgres"`
		// Password пароль пользователя базы данных
		Password string `env:"DB_PASSWORD" envDefault:"postgres"`
		// Name название базы данных
		Name string `env:"DB_NAME" envDefault:"forum"`
	}

	// JWT содержит настройки JWT токенов
	JWT struct {
		// SecretKey секретный ключ для подписи токенов
		SecretKey string `env:"JWT_SECRET_KEY" envDefault:"your-secret-key"`
		// ExpirationTime время жизни токена в часах
		ExpirationTime int `env:"JWT_EXPIRATION_TIME" envDefault:"24"`
	}
} 