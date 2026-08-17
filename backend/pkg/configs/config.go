package configs

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppName          string
	Env              string
	Version          string
	Port             string
	JWTSecret        string
	JWTRefreshSecret string
	JWTExpiry        time.Duration
	JWTRefreshExpiry time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type EmailConfig struct {
	APIKey    string
	FromEmail string
	FromName  string
}

type MFAConfig struct {
	Issuer string
}

type PaymentConfig struct {
	EsewaMerchantID  string
	EsewaSecret      string
	KhaltiMerchantID string
	KhaltiSecret     string
	MockMode         bool
}

type CORSConfig struct {
	AllowedOrigins []string
}

type ServerConfig struct {
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxRequestSize int64
}

type RateLimitConfig struct {
	RequestsPerMinute int
	BurstSize         int
}

type Config struct {
	App       AppConfig
	DB        DatabaseConfig
	Redis     RedisConfig
	Email     EmailConfig
	MFA       MFAConfig
	Payment   PaymentConfig
	CORS      CORSConfig
	Server    ServerConfig
	RateLimit RateLimitConfig
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return fallback
}

func getEnvSlice(key string, fallback []string) []string {
	if value, ok := os.LookupEnv(key); ok {
		result := []string{}
		for _, s := range splitAndTrim(value, ",") {
			if s != "" {
				result = append(result, s)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return fallback
}

func splitAndTrim(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if string(s[i]) == sep {
			result = append(result, trimSpace(s[start:i]))
			start = i + 1
		}
	}
	result = append(result, trimSpace(s[start:]))
	return result
}

func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		App: AppConfig{
			AppName:          getEnv("APP_NAME", "HinduRitualPlatform"),
			Env:              getEnv("APP_ENV", "development"),
			Version:          getEnv("APP_VERSION", "1.0.0"),
			Port:             getEnv("APP_PORT", "8080"),
			JWTSecret:        getEnv("JWT_SECRET", "change-me-in-production"),
			JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "change-me-refresh-in-production"),
			JWTExpiry:        getEnvDuration("JWT_EXPIRY", 15*time.Minute),
			JWTRefreshExpiry: getEnvDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
		},
		DB: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "hindu_ritual_platform"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Email: EmailConfig{
			APIKey:    getEnv("EMAIL_API_KEY", ""),
			FromEmail: getEnv("FROM_EMAIL", "noreply@hinduplatform.com"),
			FromName:  getEnv("FROM_NAME", "Hindu Ritual Platform"),
		},
		MFA: MFAConfig{
			Issuer: getEnv("MFA_ISSUER", "HinduRitualPlatform"),
		},
		Payment: PaymentConfig{
			EsewaMerchantID:  getEnv("ESEWA_MERCHANT_ID", ""),
			EsewaSecret:      getEnv("ESEWA_SECRET", ""),
			KhaltiMerchantID: getEnv("KHALTI_MERCHANT_ID", ""),
			KhaltiSecret:     getEnv("KHALTI_SECRET", ""),
			MockMode:         getEnvBool("PAYMENT_MOCK_MODE", true),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:3001"}),
		},
		Server: ServerConfig{
			ReadTimeout:    getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout:   getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
			MaxRequestSize: int64(getEnvInt("SERVER_MAX_REQUEST_SIZE", 10*1024*1024)),
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: getEnvInt("RATE_LIMIT_REQUESTS", 100),
			BurstSize:         getEnvInt("RATE_LIMIT_BURST", 20),
		},
	}
}
