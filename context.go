package bootstrap

import (
	"context"
	"sync"

	"github.com/kalandramo/lulu"

	v1 "github.com/kalandramo/lulu-bootstrap/conf/gen/go/bootstrap/v1"
)

// Context holds the application lifecycle state created by [Bootstrap] or [BootstrapWithContext].
//
// It is safe to read fields from multiple goroutines after Bootstrap returns.
type Context struct {
	cfg    *v1.BootstrapConfig
	app    *lulu.App
	cancel context.CancelFunc


	storages      map[string]any
	caches        map[string]any
	databases     map[string]any

	cleanupOnce sync.Once
	cleanup     func()
}

// newContext creates a Context from the Bootstrap results.
func newContext(cfg *v1.BootstrapConfig, app *lulu.App, storages map[string]any,  caches map[string]any,  databases map[string]any, cleanup func(), cancel context.CancelFunc) *Context {
	return &Context{
		cfg:           cfg,
		app:           app,
		storages:      storages,
		caches:        caches,
		databases:     databases,
		cleanup:       cleanup,
		cancel:        cancel,
	}
}

// Config returns the loaded bootstrap configuration.
func (c *Context) Config() *v1.BootstrapConfig { return c.cfg }

// App returns the underlying [*wind.App].
func (c *Context) App() *lulu.App { return c.app }

// Cancel triggers graceful shutdown (idempotent).
func (c *Context) Cancel() {
	if c.cancel != nil {
		c.cancel()
	}
}

// Cleanup releases all resources. It is safe to call multiple times.
func (c *Context) Cleanup() {
	if c == nil {
		return
	}
	c.cleanupOnce.Do(func() {
		if c.cleanup != nil {
			c.cleanup()
		}
	})
}

// Storage returns the storage client instance for the given type name (e.g.
// [StorageTypeMinio], [StorageTypeS3]).
// Returns nil if no storage with that name was configured.
//
// The caller should type-assert the result to the concrete storage type:
//
//	s, ok := ctx.Storage(bootstrap.StorageTypeMinio).(*minioPlugin.Storage)
//	if ok { /* use s for PutObject/GetObject */ }
func (c *Context) Storage(name string) any {
	if c == nil || c.storages == nil {
		return nil
	}
	return c.storages[name]
}

// Storages returns all storage instances as a map keyed by type name.
// Returns nil if no storage was configured.
func (c *Context) Storages() map[string]any {
	if c == nil {
		return nil
	}
	return c.storages
}

// Cache returns the cache instance for the given type name (e.g.
// [CacheTypeLocal], [CacheTypeRedis]).
// Returns nil if no cache with that name was configured.
//
// The caller should type-assert the result to the concrete cache type:
//
//	c, ok := ctx.Cache(bootstrap.CacheTypeLocal).(*local.Cache)
//	if ok { /* use c for Get/Set operations */ }
func (c *Context) Cache(name string) any {
	if c == nil || c.caches == nil {
		return nil
	}
	return c.caches[name]
}

// Caches returns all cache instances as a map keyed by type name.
// Returns nil if no cache was configured.
func (c *Context) Caches() map[string]any {
	if c == nil {
		return nil
	}
	return c.caches
}

// Database returns the database client instance for the given type name (e.g.
// [DatabaseTypeGorm], [DatabaseTypeMongodb]).
// Returns nil if no database client with that name was configured.
func (c *Context) Database(name string) any {
	if c == nil || c.databases == nil {
		return nil
	}
	return c.databases[name]
}

// Databases returns all database client instances as a map keyed by type name.
// Returns nil if no database was configured.
func (c *Context) Databases() map[string]any {
	if c == nil {
		return nil
	}
	return c.databases
}
