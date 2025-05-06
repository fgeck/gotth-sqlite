//go:build unittest

package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fgeck/go-register/internal/service/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTempConfigFile(content string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "config")
	if err != nil {
		return "", err
	}
	tmpFile, err := os.Create(tmpDir + "/config.yaml")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(content); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}
func TestLoadConfig(t *testing.T) {
	validConfig := `
app:
  host: localhost
  port: 8080
  jwtSecret: change-m3-@$ap!
  adminUser: admin
  adminPassword: adminpassword
  adminEmail: wasd
db:
  persistence: FILE
  migrationsPath: "../../migrations"
`

	t.Run("successfully loads valid config", func(t *testing.T) {
		configFile, err := createTempConfigFile(validConfig)
		require.NoError(t, err)
		configPath := filepath.Dir(configFile)
		defer os.Remove(configPath)

		os.Setenv("APP_HOST", "127.0.0.1")
		os.Setenv("APP_PORT", "9090")
		os.Setenv("APP_ADMINUSER", "admin")
		os.Setenv("APP_ADMINPASSWORD", "adminpassword")
		os.Setenv("APP_ADMINEMAIL", "adm@test.io")
		os.Setenv("DB_PERSISTENCE", "memory")
		os.Setenv("DB_MIGRATIONSPATH", "./test/migrations")
		defer os.Unsetenv("APP_HOST")
		defer os.Unsetenv("APP_PORT")
		defer os.Unsetenv("APP_ADMINUSER")
		defer os.Unsetenv("APP_ADMINPASSWORD")
		defer os.Unsetenv("APP_ADMINEMAIL")
		defer os.Unsetenv("DB_PERSISTENCE")
		defer os.Unsetenv("DB_MIGRATIONSPATH")

		// Load the config
		loader := config.NewLoader()
		config, err := loader.LoadConfig(configPath)

		// Validate the results
		require.NoError(t, err)
		assert.Equal(t, "127.0.0.1", config.App.Host)
		assert.Equal(t, "9090", config.App.Port)
		assert.Equal(t, "change-m3-@$ap!", config.App.JwtSecret)
		assert.Equal(t, "admin", config.App.AdminUser)
		assert.Equal(t, "adminpassword", config.App.AdminPassword)
		assert.Equal(t, "adm@test.io", config.App.AdminEmail)
		assert.Equal(t, "memory", config.Db.Persistence)
		assert.Equal(t, "./test/migrations", config.Db.MigrationsPath)
	})

	t.Run("fails when config file is missing", func(t *testing.T) {
		loader := config.NewLoader()
		config, err := loader.LoadConfig("nonexistent.yaml")

		require.Error(t, err)
		assert.Nil(t, config)
	})

	t.Run("fails when config file is invalid", func(t *testing.T) {
		invalidConfig := `
app:
  host: localhost
  port: 8080
  jwtSecret: change-m3-@$ap!
db:
  persistence: FILE
  migrationsPath: "../../migrations"
`
		configPath, err := createTempConfigFile(invalidConfig)
		require.NoError(t, err)
		defer os.Remove(configPath)

		loader := config.NewLoader()
		config, err := loader.LoadConfig(configPath)

		require.Error(t, err)
		assert.Nil(t, config)
	})

	t.Run("uses default values when environment variables are not set", func(t *testing.T) {
		configFile, err := createTempConfigFile(validConfig)
		require.NoError(t, err)
		configPath := filepath.Dir(configFile)
		defer os.Remove(configPath)

		// Ensure no environment variables are set
		os.Unsetenv("APP_HOST")
		os.Unsetenv("APP_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")

		// Load the config
		loader := config.NewLoader()
		config, err := loader.LoadConfig(configPath)

		// Validate the results
		require.NoError(t, err)
		assert.Equal(t, "localhost", config.App.Host)
		assert.Equal(t, "8080", config.App.Port)
		assert.Equal(t, "change-m3-@$ap!", config.App.JwtSecret)
		assert.Equal(t, "FILE", config.Db.Persistence)
		assert.Equal(t, "../../migrations", config.Db.MigrationsPath)
	})
}
