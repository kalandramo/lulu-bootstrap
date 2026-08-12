package bootstrap

// Server type 常量，用于注册 Builder 时的 key。
const (
	ServerTypeHTTP         = "http"
)

// Config type 常量，用于注册 ConfigAction 时的 key。
const (
	ConfigTypeFile       = "file"
	ConfigTypeNacos      = "nacos"
)

// Logger type 常量。
const (
	LoggerTypeSlog       = "slog"
)

// Tracer type 常量。
const (
	TracerTypeOTLP = "otlp"
)

// Metrics type 常量。
const (
	MetricsTypePrometheus = "prometheus"
	MetricsTypeOTLP       = "otlp"
)

// Storage type 常量，用于注册 StorageBuilder 时的 key。
const (
	StorageTypeMinio = "minio"
	StorageTypeS3    = "s3"
)

// Cache type 常量，用于注册 CacheBuilder 时的 key。
const (
	CacheTypeLocal = "local"
	CacheTypeRedis = "redis"
)

// Database type 常量，用于注册 DatabaseBuilder 时的 key。
const (
	DatabaseTypeGorm          = "gorm"
	DatabaseTypeElasticsearch = "elasticsearch"
)
