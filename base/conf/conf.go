package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

const (
	defaultLoglevel               = "info"
	defaultServerBind             = "0.0.0.0:8080"
	defaultJwtIssuer              = "api-server"
	defaultJwtExpireTime          = "1h"
	defaultLoadCasbinTimeDuration = 5 * time.Second
	defaultRedisExpireTime        = "30s"
)

func LoadConfig(configPath string) (err error) {
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("configuration file %s does not exist", configPath)
	}
	if err != nil {
		return fmt.Errorf("stat configuration file %s faild. err: %w", configPath, err)
	}
	zap.L().Info("start loading configuration file", zap.String("path", configPath))
	configDir := filepath.Dir(configPath)
	configBase := filepath.Base(configPath)
	viper.SetConfigName(configBase)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)
	if err = viper.ReadInConfig(); err != nil {
		return fmt.Errorf("reading configuration files %s faild. err: %w", configPath, err)
	}
	zap.S().Info("configuration file read successfully")
	return nil
}

func GetServerCompress() bool {
	compress := viper.GetBool("server.compress")
	return compress
}

func GetServerLogLevel() string {
	logLevel := viper.GetString("server.logLevel")
	if logLevel == "" {
		logLevel = defaultLoglevel
	}
	return logLevel
}

func GetServerBind() string {
	bind := viper.GetString("server.bind")
	if bind == "" {
		bind = defaultServerBind
	}
	return bind
}

func GetServerSalt() (string, error) {
	salt := viper.GetString("server.salt")
	if salt == "" {
		return "", fmt.Errorf("server.salt is empty")
	}
	return salt, nil
}

func GetServerAutoLoadCasbinRule() time.Duration {
	duration := viper.GetDuration("server.autoLoadCasbinRule")
	if duration == 0 {
		return defaultLoadCasbinTimeDuration
	}
	return duration
}

func GetJwtSecret() (string, error) {
	secret := viper.GetString("jwt.secret")
	if secret == "" {
		return "", fmt.Errorf("jwt.secret is empty")
	}
	return secret, nil
}

func GetJwtIssuer() string {
	issuer := viper.GetString("jwt.issuer")
	if issuer == "" {
		zap.S().Infof("jwt.issuer is empty, set default jwt.issuer: %s", defaultJwtIssuer)
		return defaultJwtIssuer
	}
	return issuer
}

func GetJwtExpirationTime() (time.Duration, error) {
	expireTime := viper.GetDuration("jwt.expireTime")
	if expireTime == 0 {
		expire, err := time.ParseDuration(defaultJwtExpireTime)
		if err != nil {
			return 0, fmt.Errorf("failed to parser jwt.expireTime err: %v", err)
		}
		return expire, nil
	}
	return expireTime, nil
}

func GetMysqlDsn() (dsn string, err error) {
	user := viper.GetString("mysql.username")
	if user == "" {
		return "", fmt.Errorf("mysql.username is empty")
	}
	pas := viper.GetString("mysql.password")
	if pas == "" {
		return "", fmt.Errorf("mysql.password is empty")
	}
	host := viper.GetString("mysql.host")
	if host == "" {
		return "", fmt.Errorf("mysql.host is empty")
	}
	database := viper.GetString("mysql.database")
	if database == "" {
		return "", fmt.Errorf("mysql.database is empty")
	}
	dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local&timeout=10000ms",
		user,
		pas,
		host,
		database,
	)
	return dsn, nil
}

func GetMysqlMaxIdleConns() int {
	maxIdleConns := viper.GetInt("mysql.maxIdleConns")
	if maxIdleConns == 0 {
		return 10
	}
	return maxIdleConns
}

func GetMysqlMaxOpenConns() int {
	maxOpenConns := viper.GetInt("mysql.maxOpenConns")
	if maxOpenConns == 0 {
		return 30
	}
	return maxOpenConns
}

func GetMysqlMaxLifetime() time.Duration {
	maxLifetime := viper.GetDuration("mysql.maxLifetime")
	if maxLifetime == 0 {
		return 30 * time.Minute
	}
	return maxLifetime
}

func GetRedisPoolSize() int {
	poolSize := viper.GetInt("redis.poolSize")
	if poolSize == 0 {
		return 50
	}
	return poolSize
}

func GetRedisMinIdleConns() int {
	minIdleConns := viper.GetInt("redis.minIdleConns")
	if minIdleConns == 0 {
		return 20
	}
	return minIdleConns
}

func GetRedisConnMaxLifetime() time.Duration {
	connMaxLifetime := viper.GetDuration("redis.connMaxLifetime")
	if connMaxLifetime == 0 {
		return 30 * time.Minute
	}
	return connMaxLifetime
}

func GetRedisPassword() (string, error) {
	password := viper.GetString("redis.password")
	if password == "" {
		return "", fmt.Errorf("redis.password is empty")
	}
	return password, nil
}

func GetRedisMasterName() (string, error) {
	masterName := viper.GetString("redis.sentinel.masterName")
	if masterName == "" {
		return "", fmt.Errorf("redis.sentinel.masterName is empty")
	}
	return masterName, nil
}

func GetRedisSentinelPassword() (string, error) {
	sentPassword := viper.GetString("redis.sentinel.password")
	if sentPassword == "" {
		return "", fmt.Errorf("redis.sentinel.password is empty")
	}
	return sentPassword, nil
}

func GetRedisSentinelHosts() ([]string, error) {
	sentinelHosts := viper.GetStringSlice("redis.sentinel.hosts")
	if len(sentinelHosts) == 0 {
		return nil, fmt.Errorf("redis.sentinel.hosts is empty")
	}
	return sentinelHosts, nil
}

func GetRedisHost() (string, error) {
	host := viper.GetString("redis.host")
	if host == "" {
		return "", fmt.Errorf("redis.host is empty")
	}
	return host, nil
}

func GetRedisDB() int {
	return viper.GetInt("redis.db")
}

func GetRedisMode() string {
	return viper.GetString("redis.mode")
}

func GetRedisExpireTime() (time.Duration, error) {
	expireTime := viper.GetDuration("redis.expireTime")
	if expireTime == 0 {
		duration, err := time.ParseDuration(defaultRedisExpireTime)
		if err != nil {
			return 0, fmt.Errorf("failed to parser defaultRedisExpireTime err: %v", err)
		}
		zap.S().Infof("redis.expireTime is empty, set default expireTime: %s", defaultRedisExpireTime)
		return duration, nil
	}

	return expireTime, nil
}

func GetRedisKeyPrefix() (string, error) {
	prefix := viper.GetString("redis.keyPrefix")
	if prefix == "" {
		return "", fmt.Errorf("redis.keyPrefix is empty")
	}
	return prefix, nil
}

func GetOauth2Config() (*oauth2.Config, error) {
	clientID := viper.GetString("oauth2.clientId")
	if clientID == "" {
		return nil, fmt.Errorf("oauth2.clientId is empty")
	}
	clientSecret := viper.GetString("oauth2.clientSecret")
	if clientSecret == "" {
		return nil, fmt.Errorf("oauth2.clientSecret is empty")
	}
	authUrl := viper.GetString("oauth2.authUrl")
	if authUrl == "" {
		return nil, fmt.Errorf("oauth2.authUrl is empty")
	}
	tokenUrl := viper.GetString("oauth2.tokenUrl")
	if tokenUrl == "" {
		return nil, fmt.Errorf("oauth2.tokenUrl is empty")
	}
	redirectUrl := viper.GetString("oauth2.redirectUrl")
	if redirectUrl == "" {
		return nil, fmt.Errorf("oauth2.redirectUrl is empty")
	}
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authUrl,
			TokenURL: tokenUrl,
		},
		RedirectURL: redirectUrl,
	}, nil
}

func GetOauth2State() string {
	state := viper.GetString("oauth2.state")
	if state == "" {
		return "random_state"
	}
	return state
}

func GetOauth2Verifier() string {
	verifier := viper.GetString("oauth2.verifier")
	if verifier == "" {
		return "code_verifier"
	}
	return verifier
}
