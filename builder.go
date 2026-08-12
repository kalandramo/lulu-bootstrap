// Package bootstrap provides builder registration for all plugin domains.
//
// Each plugin domain (server, config, registry, log, tracer, metrics, broker)
// maintains a map of string-keyed builder functions. The key is a lowercase
// type string (e.g. "consul", "zap") matching the JSON config value.
// Built-in builders are registered via init() in provider packages.
// Users can register custom builders to extend the framework.
package bootstrap

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/kalandramo/lulu/log"
	"github.com/kalandramo/lulu/transport"

	v1 "github.com/kalandramo/lulu-bootstrap/conf/gen/go/bootstrap/v1"
)

// ServerBuilder builds a [transport.Server] from a [Server] config.
type ServerBuilder func(cfg *v1.Server) (transport.Server, error)

// LogBuilder builds a [log.Logger] from a [Logger] config.
type LogBuilder func(cfg *v1.Logger) (log.Logger, func(), error)

// ConfigAction performs config source loading/watching.
type ConfigAction func(ctx context.Context, cfg *v1.Config) (func(), error)

// TracerBuilder builds a tracer provider.
type TracerBuilder func(cfg *v1.Tracer) (interface{}, func(), error)

// MetricsBuilder builds a metrics backend.
type MetricsBuilder func(cfg *v1.Metrics) (func(), error)

// StorageBuilder builds a storage client instance and returns it along with an
// optional cleanup function. The returned instance (any) is the concrete
// storage client that callers can use for object operations.
//
// The type key (e.g. "minio", "s3") is used to look up the instance
// via [Context.Storage] after bootstrap.
type StorageBuilder func(ctx context.Context, cfg *v1.Storage) (any, func(), error)

// CacheBuilder builds a cache instance and returns it along with an optional
// cleanup function.
type CacheBuilder func(ctx context.Context, cfg *v1.Cache) (any, func(), error)

// DatabaseBuilder builds a database client instance and returns it along with
// an optional cleanup function.
type DatabaseBuilder func(ctx context.Context, cfg *v1.Database) (any, func(), error)

// ---- Global registries (string keyed) ----

var (
	mu                   sync.RWMutex
	serverBuilders       = map[string]ServerBuilder{}
	logBuilders          = map[string]LogBuilder{}
	configActions        = map[string]ConfigAction{}
	tracerBuilders       = map[string]TracerBuilder{}
	metricsBuilders      = map[string]MetricsBuilder{}
	storageBuilders      = map[string]StorageBuilder{}
	cacheBuilders        = map[string]CacheBuilder{}
	databaseBuilders     = map[string]DatabaseBuilder{}
)

// ---- Register functions ----

// RegisterServerBuilder registers a server builder for the given type string.
func RegisterServerBuilder(typ string, b ServerBuilder) error {
	if typ == "" {
		return fmt.Errorf("bootstrap: type is empty")
	}
	if b == nil {
		return fmt.Errorf("bootstrap: factory is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := serverBuilders[typ]; ok {
		return fmt.Errorf("bootstrap: server builder %q already registered", typ)
	}
	serverBuilders[typ] = b
	return nil
}

// MustRegisterServerBuilder panics on error. Intended for init().
func MustRegisterServerBuilder(typ string, b ServerBuilder) {
	if err := RegisterServerBuilder(typ, b); err != nil {
		panic(err)
	}
}

// RegisterLogBuilder registers a log builder for the given type string.
func RegisterLogBuilder(typ string, b LogBuilder) error {
	if typ == "" {
		return fmt.Errorf("bootstrap: type is empty")
	}
	if b == nil {
		return fmt.Errorf("bootstrap: factory is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := logBuilders[typ]; ok {
		return fmt.Errorf("bootstrap: log builder %q already registered", typ)
	}
	logBuilders[typ] = b
	return nil
}

// MustRegisterLogBuilder panics on error.
func MustRegisterLogBuilder(typ string, b LogBuilder) {
	if err := RegisterLogBuilder(typ, b); err != nil {
		panic(err)
	}
}

// RegisterConfigAction registers a config action for the given type string.
func RegisterConfigAction(typ string, a ConfigAction) error {
	if typ == "" {
		return fmt.Errorf("bootstrap: type is empty")
	}
	if a == nil {
		return fmt.Errorf("bootstrap: factory is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := configActions[typ]; ok {
		return fmt.Errorf("bootstrap: config action %q already registered", typ)
	}
	configActions[typ] = a
	return nil
}

// MustRegisterConfigAction panics on error.
func MustRegisterConfigAction(typ string, a ConfigAction) {
	if err := RegisterConfigAction(typ, a); err != nil {
		panic(err)
	}
}

// RegisterTracerBuilder registers a tracer builder for the given type string.
func RegisterTracerBuilder(typ string, b TracerBuilder) error {
	if typ == "" {
		return fmt.Errorf("bootstrap: type is empty")
	}
	if b == nil {
		return fmt.Errorf("bootstrap: factory is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := tracerBuilders[typ]; ok {
		return fmt.Errorf("bootstrap: tracer builder %q already registered", typ)
	}
	tracerBuilders[typ] = b
	return nil
}

// MustRegisterTracerBuilder panics on error.
func MustRegisterTracerBuilder(typ string, b TracerBuilder) {
	if err := RegisterTracerBuilder(typ, b); err != nil {
		panic(err)
	}
}

// RegisterMetricsBuilder registers a metrics builder for the given type string.
func RegisterMetricsBuilder(typ string, b MetricsBuilder) error {
	if typ == "" {
		return fmt.Errorf("bootstrap: type is empty")
	}
	if b == nil {
		return fmt.Errorf("bootstrap: factory is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := metricsBuilders[typ]; ok {
		return fmt.Errorf("bootstrap: metrics builder %q already registered", typ)
	}
	metricsBuilders[typ] = b
	return nil
}

// MustRegisterMetricsBuilder panics on error.
func MustRegisterMetricsBuilder(typ string, b MetricsBuilder) {
	if err := RegisterMetricsBuilder(typ, b); err != nil {
		panic(err)
	}
}

// RegisterStorageBuilder registers a storage builder for the given type string.
func RegisterStorageBuilder(typ string, b StorageBuilder) error {
	if typ == "" {
		return fmt.Errorf("bootstrap: type is empty")
	}
	if b == nil {
		return fmt.Errorf("bootstrap: factory is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := storageBuilders[typ]; ok {
		return fmt.Errorf("bootstrap: storage builder %q already registered", typ)
	}
	storageBuilders[typ] = b
	return nil
}

// MustRegisterStorageBuilder panics on error.
func MustRegisterStorageBuilder(typ string, b StorageBuilder) {
	if err := RegisterStorageBuilder(typ, b); err != nil {
		panic(err)
	}
}

// RegisterCacheBuilder registers a cache builder for the given type string.
func RegisterCacheBuilder(typ string, b CacheBuilder) error {
	if typ == "" {
		return fmt.Errorf("bootstrap: type is empty")
	}
	if b == nil {
		return fmt.Errorf("bootstrap: factory is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := cacheBuilders[typ]; ok {
		return fmt.Errorf("bootstrap: cache builder %q already registered", typ)
	}
	cacheBuilders[typ] = b
	return nil
}

// MustRegisterCacheBuilder panics on error.
func MustRegisterCacheBuilder(typ string, b CacheBuilder) {
	if err := RegisterCacheBuilder(typ, b); err != nil {
		panic(err)
	}
}

// RegisterDatabaseBuilder registers a database builder for the given type string.
func RegisterDatabaseBuilder(typ string, b DatabaseBuilder) error {
	if typ == "" {
		return fmt.Errorf("bootstrap: type is empty")
	}
	if b == nil {
		return fmt.Errorf("bootstrap: factory is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := databaseBuilders[typ]; ok {
		return fmt.Errorf("bootstrap: database builder %q already registered", typ)
	}
	databaseBuilders[typ] = b
	return nil
}

// MustRegisterDatabaseBuilder panics on error.
func MustRegisterDatabaseBuilder(typ string, b DatabaseBuilder) {
	if err := RegisterDatabaseBuilder(typ, b); err != nil {
		panic(err)
	}
}

// ---- Lookup helpers ----

func getServerBuilder(typ string) (ServerBuilder, error) {
	mu.RLock()
	defer mu.RUnlock()
	b, ok := serverBuilders[typ]
	if !ok {
		return nil, fmt.Errorf("bootstrap: no server builder registered for %q", typ)
	}
	return b, nil
}

func getLogBuilder(typ string) (LogBuilder, error) {
	mu.RLock()
	defer mu.RUnlock()
	b, ok := logBuilders[typ]
	if !ok {
		return nil, fmt.Errorf("bootstrap: no log builder registered for %q", typ)
	}
	return b, nil
}

func getConfigAction(typ string) (ConfigAction, error) {
	mu.RLock()
	defer mu.RUnlock()
	a, ok := configActions[typ]
	if !ok {
		return nil, fmt.Errorf("bootstrap: no config action registered for %q", typ)
	}
	return a, nil
}

func getTracerBuilder(typ string) (TracerBuilder, error) {
	mu.RLock()
	defer mu.RUnlock()
	b, ok := tracerBuilders[typ]
	if !ok {
		return nil, fmt.Errorf("bootstrap: no tracer builder registered for %q", typ)
	}
	return b, nil
}

func getMetricsBuilder(typ string) (MetricsBuilder, error) {
	mu.RLock()
	defer mu.RUnlock()
	b, ok := metricsBuilders[typ]
	if !ok {
		return nil, fmt.Errorf("bootstrap: no metrics builder registered for %q", typ)
	}
	return b, nil
}

func getStorageBuilder(typ string) (StorageBuilder, error) {
	mu.RLock()
	defer mu.RUnlock()
	b, ok := storageBuilders[typ]
	if !ok {
		return nil, fmt.Errorf("bootstrap: no storage builder registered for %q", typ)
	}
	return b, nil
}

func getCacheBuilder(typ string) (CacheBuilder, error) {
	mu.RLock()
	defer mu.RUnlock()
	b, ok := cacheBuilders[typ]
	if !ok {
		return nil, fmt.Errorf("bootstrap: no cache builder registered for %q", typ)
	}
	return b, nil
}

func getDatabaseBuilder(typ string) (DatabaseBuilder, error) {
	mu.RLock()
	defer mu.RUnlock()
	b, ok := databaseBuilders[typ]
	if !ok {
		return nil, fmt.Errorf("bootstrap: no database builder registered for %q", typ)
	}
	return b, nil
}

// ---- Diagnostic helpers ----

// ListServerBuilders returns all registered server type names.
func ListServerBuilders() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(serverBuilders))
	for k := range serverBuilders {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// ListLogBuilders returns all registered log type names.
func ListLogBuilders() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(logBuilders))
	for k := range logBuilders {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// ListConfigActions returns all registered config type names.
func ListConfigActions() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(configActions))
	for k := range configActions {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// ListTracerBuilders returns all registered tracer type names.
func ListTracerBuilders() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(tracerBuilders))
	for k := range tracerBuilders {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// ListMetricsBuilders returns all registered metrics type names.
func ListMetricsBuilders() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(metricsBuilders))
	for k := range metricsBuilders {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// ListStorageBuilders returns all registered storage type names.
func ListStorageBuilders() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(storageBuilders))
	for k := range storageBuilders {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// ListCacheBuilders returns all registered cache type names.
func ListCacheBuilders() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(cacheBuilders))
	for k := range cacheBuilders {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// ListDatabaseBuilders returns all registered database type names.
func ListDatabaseBuilders() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(databaseBuilders))
	for k := range databaseBuilders {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}
