package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestLoadReadsDotenvFile verifies that Load() merges values from a `.env`.
func TestLoadReadsDotenvFile(t *testing.T) {
	// Make a .env file in TempDir for testing
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(envPath, []byte("STRIDE_BACKEND_APP_PORT=9830\nSTRIDE_BACKEND_DB_PASSWORD=from-dotenv-file\n"), 0o600))

	// Change working directory. And Back to work directory after the test
	wd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Chdir(wd) })
	require.NoError(t, os.Chdir(dir))

	// Make sure there isn't real env
	t.Setenv("STRIDE_BACKEND_APP_PORT", "")
	t.Setenv("STRIDE_BACKEND_DB_PASSWORD", "")
	require.NoError(t, os.Unsetenv("STRIDE_BACKEND_APP_PORT"))
	require.NoError(t, os.Unsetenv("STRIDE_BACKEND_DB_PASSWORD"))

	// Test
	cfg, err := Load()
	require.NoError(t, err)
	require.Equal(t, 9830, cfg.AppPort)
	require.Equal(t, "from-dotenv-file", cfg.DBPassword)
}

func TestRealEnvOverridesDotenv(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(envPath, []byte("STRIDE_BACKEND_APP_MODE=debug"), 0o600))

	wd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Chdir(wd)})
	require.NoError(t, os.Chdir(dir))

	t.Setenv("STRIDE_BACKEND_APP_MODE", "test")

	cfg, err := Load()
	require.NoError(t, err)
	require.Equal(t, Test, cfg.AppMode)
}

func TestDotenvLocalOverridesDotenv(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(envPath, []byte("STRIDE_BACKEND_APP_PORT=9830\nSTRIDE_BACKEND_DB_PASSWORD=from-dotenv"), 0o600))
	envPathLocal := filepath.Join(dir, ".env.local")
	require.NoError(t, os.WriteFile(envPathLocal, []byte("STRIDE_BACKEND_DB_PASSWORD=from-dotenvlocal\n"), 0o600))

	wd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Chdir(wd)})
	require.NoError(t, os.Chdir(dir))

	t.Setenv("STRIDE_BACKEND_APP_PORT", "")
	t.Setenv("STRIDE_BACKEND_DB_PASSWORD", "")
	require.NoError(t, os.Unsetenv("STRIDE_BACKEND_APP_PORT"))
	require.NoError(t, os.Unsetenv("STRIDE_BACKEND_DB_PASSWORD"))

	cfg, err := Load()
	require.NoError(t, err)
	require.Equal(t, "from-dotenvlocal", cfg.DBPassword, ".env.local should override .env")
	require.Equal(t, 9830, cfg.AppPort, "keys only in .env should still apply")
}