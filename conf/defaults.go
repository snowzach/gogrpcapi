package conf

import (
	config "github.com/spf13/viper"
)

func init() {

	// Logger Defaults
	config.SetDefault("logger.level", "info")
	config.SetDefault("logger.encoding", "console")
	config.SetDefault("logger.color", true)
	config.SetDefault("logger.dev_mode", true)
	config.SetDefault("logger.disable_caller", false)
	config.SetDefault("logger.disable_stacktrace", true)

	// Pidfile
	config.SetDefault("pidfile", "")

	// Profiler config
	config.SetDefault("profiler.enabled", false)
	config.SetDefault("profiler.host", "")
	config.SetDefault("profiler.port", "6060")

	// Server Configuration
	config.SetDefault("server.host", "")
	config.SetDefault("server.port", "8900")
	config.SetDefault("server.tls", false)
	config.SetDefault("server.devcert", false)
	config.SetDefault("server.certfile", "server.crt")
	config.SetDefault("server.keyfile", "server.key")
	config.SetDefault("server.log_requests", true)
	config.SetDefault("server.profiler_enabled", false)
	config.SetDefault("server.profiler_path", "/debug")
	// GRPC JSON Marshaler Options
	config.SetDefault("server.rest.enums_as_ints", false)
	config.SetDefault("server.rest.emit_defaults", true)
	config.SetDefault("server.rest.orig_names", true)

	// Database Settings
	config.SetDefault("storage.type", "postgres")
	config.SetDefault("storage.username", "postgres")
	config.SetDefault("storage.password", "password")
	config.SetDefault("storage.host", "postgres")
	config.SetDefault("storage.port", 5432)
	config.SetDefault("storage.database", "gogrpcapi")
	config.SetDefault("storage.sslmode", "disable")
	config.SetDefault("storage.retries", 5)
	config.SetDefault("storage.sleep_between_retries", "7s")
	config.SetDefault("storage.max_connections", 80)
	config.SetDefault("storage.wipe_confirm", false)

}
