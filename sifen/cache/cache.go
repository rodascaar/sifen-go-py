package cache

import (
	"sync"
	"time"

	"github.com/rodascaar/sifen-go-py/sifen/response"
)

// ============================================================================
// Cache de Consultas SIFEN
// ============================================================================

// CacheConfig contiene la configuración del caché
type CacheConfig struct {
	// TTL por defecto para entradas de caché
	DefaultTTL time.Duration
	// TTL específico para consultas RUC
	RUCTTL time.Duration
	// TTL específico para consultas DE
	DETTL time.Duration
	// Tamaño máximo del caché (0 = sin límite)
	MaxSize int
	// Habilitar limpieza automática de entradas expiradas
	EnableAutoCleanup bool
	// Intervalo de limpieza automática
	CleanupInterval time.Duration
}

// DefaultCacheConfig retorna configuración por defecto
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		DefaultTTL:        5 * time.Minute,
		RUCTTL:            30 * time.Minute, // RUC cambia poco, cachear más tiempo
		DETTL:             10 * time.Minute, // DE puede cambiar estado
		MaxSize:           1000,
		EnableAutoCleanup: true,
		CleanupInterval:   time.Minute,
	}
}

// ============================================================================
// Entrada de Caché
// ============================================================================

type cacheEntry struct {
	value      interface{}
	expiration time.Time
	hitCount   int
}

func (e *cacheEntry) isExpired() bool {
	return time.Now().After(e.expiration)
}

// ============================================================================
// Cache Principal
// ============================================================================

// Cache implementa un caché in-memory con TTL
type Cache struct {
	config  CacheConfig
	entries map[string]*cacheEntry
	mu      sync.RWMutex
	stats   CacheStats
	stopCh  chan struct{}
}

// CacheStats contiene estadísticas del caché
type CacheStats struct {
	Hits        int64
	Misses      int64
	Evictions   int64
	Size        int
	OldestEntry time.Time
}

// NewCache crea un nuevo caché
func NewCache(config CacheConfig) *Cache {
	c := &Cache{
		config:  config,
		entries: make(map[string]*cacheEntry),
		stopCh:  make(chan struct{}),
	}

	if config.EnableAutoCleanup && config.CleanupInterval > 0 {
		go c.cleanupLoop()
	}

	return c
}

// ============================================================================
// Operaciones de Caché Genéricas
// ============================================================================

// Set almacena un valor con TTL personalizado
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Verificar límite de tamaño
	if c.config.MaxSize > 0 && len(c.entries) >= c.config.MaxSize {
		c.evictOldest()
	}

	c.entries[key] = &cacheEntry{
		value:      value,
		expiration: time.Now().Add(ttl),
		hitCount:   0,
	}
}

// Get obtiene un valor del caché
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	entry, exists := c.entries[key]
	c.mu.RUnlock()

	if !exists {
		c.mu.Lock()
		c.stats.Misses++
		c.mu.Unlock()
		return nil, false
	}

	if entry.isExpired() {
		c.mu.Lock()
		delete(c.entries, key)
		c.stats.Misses++
		c.mu.Unlock()
		return nil, false
	}

	c.mu.Lock()
	entry.hitCount++
	c.stats.Hits++
	c.mu.Unlock()

	return entry.value, true
}

// Delete elimina una entrada del caché
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Clear limpia todo el caché
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*cacheEntry)
}

// Stats retorna estadísticas del caché
func (c *Cache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	stats := c.stats
	stats.Size = len(c.entries)
	return stats
}

// Close detiene el cleanup loop y libera recursos
func (c *Cache) Close() {
	close(c.stopCh)
}

// ============================================================================
// Caché Específico para RUC
// ============================================================================

// RUCCache especializado para consultas RUC
type RUCCache struct {
	cache *Cache
}

// NewRUCCache crea un caché especializado para RUC
func NewRUCCache(config CacheConfig) *RUCCache {
	return &RUCCache{cache: NewCache(config)}
}

// GetRUC obtiene un RUC del caché
func (r *RUCCache) GetRUC(ruc string) (*response.RespuestaConsultaRUC, bool) {
	value, ok := r.cache.Get("ruc:" + ruc)
	if !ok {
		return nil, false
	}
	resp, ok := value.(*response.RespuestaConsultaRUC)
	return resp, ok
}

// SetRUC almacena un RUC en el caché
func (r *RUCCache) SetRUC(ruc string, resp *response.RespuestaConsultaRUC) {
	r.cache.Set("ruc:"+ruc, resp, r.cache.config.RUCTTL)
}

// InvalidateRUC invalida una entrada RUC
func (r *RUCCache) InvalidateRUC(ruc string) {
	r.cache.Delete("ruc:" + ruc)
}

// Stats retorna estadísticas
func (r *RUCCache) Stats() CacheStats {
	return r.cache.Stats()
}

// Close cierra el caché
func (r *RUCCache) Close() {
	r.cache.Close()
}

// ============================================================================
// Caché Específico para DE
// ============================================================================

// DECache especializado para consultas DE
type DECache struct {
	cache *Cache
}

// NewDECache crea un caché especializado para DE
func NewDECache(config CacheConfig) *DECache {
	return &DECache{cache: NewCache(config)}
}

// GetDE obtiene un DE del caché por CDC
func (d *DECache) GetDE(cdc string) (*response.RespuestaConsultaDE, bool) {
	value, ok := d.cache.Get("de:" + cdc)
	if !ok {
		return nil, false
	}
	resp, ok := value.(*response.RespuestaConsultaDE)
	return resp, ok
}

// SetDE almacena un DE en el caché
func (d *DECache) SetDE(cdc string, resp *response.RespuestaConsultaDE) {
	d.cache.Set("de:"+cdc, resp, d.cache.config.DETTL)
}

// InvalidateDE invalida una entrada DE
func (d *DECache) InvalidateDE(cdc string) {
	d.cache.Delete("de:" + cdc)
}

// Stats retorna estadísticas
func (d *DECache) Stats() CacheStats {
	return d.cache.Stats()
}

// Close cierra el caché
func (d *DECache) Close() {
	d.cache.Close()
}

// ============================================================================
// Helpers Internos
// ============================================================================

func (c *Cache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range c.entries {
		if oldestKey == "" || entry.expiration.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.expiration
		}
	}

	if oldestKey != "" {
		delete(c.entries, oldestKey)
		c.stats.Evictions++
	}
}

func (c *Cache) cleanupLoop() {
	ticker := time.NewTicker(c.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopCh:
			return
		}
	}
}

func (c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.entries {
		if now.After(entry.expiration) {
			delete(c.entries, key)
			c.stats.Evictions++
		}
	}
}

// ============================================================================
// SifenCache - Caché unificado para el cliente
// ============================================================================

// SifenCache combina todos los cachés específicos
type SifenCache struct {
	RUC *RUCCache
	DE  *DECache

	// Caché genérico para otros usos
	Generic *Cache
}

// NewSifenCache crea un caché unificado con configuración por defecto
func NewSifenCache() *SifenCache {
	config := DefaultCacheConfig()
	return &SifenCache{
		RUC:     NewRUCCache(config),
		DE:      NewDECache(config),
		Generic: NewCache(config),
	}
}

// NewSifenCacheWithConfig crea un caché unificado con configuración personalizada
func NewSifenCacheWithConfig(config CacheConfig) *SifenCache {
	return &SifenCache{
		RUC:     NewRUCCache(config),
		DE:      NewDECache(config),
		Generic: NewCache(config),
	}
}

// Close cierra todos los cachés
func (s *SifenCache) Close() {
	s.RUC.Close()
	s.DE.Close()
	s.Generic.Close()
}

// ClearAll limpia todos los cachés
func (s *SifenCache) ClearAll() {
	s.RUC.cache.Clear()
	s.DE.cache.Clear()
	s.Generic.Clear()
}

// AllStats retorna estadísticas de todos los cachés
func (s *SifenCache) AllStats() map[string]CacheStats {
	return map[string]CacheStats{
		"RUC":     s.RUC.Stats(),
		"DE":      s.DE.Stats(),
		"Generic": s.Generic.Stats(),
	}
}
