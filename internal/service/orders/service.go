package orders

import (
	"container/list"
	"sync"
	"time"
)

type cacheEntry struct {
	key      string
	value    Model
	expireAt time.Time
	element  *list.Element
}

type Service struct {
	repository Repository
	transactor Transactor

	mu        sync.RWMutex
	cache     map[string]*cacheEntry
	lru       *list.List
	cacheCap  int
	cacheTTL  time.Duration
	stopClean chan struct{}
}

const (
	defaultCacheCap = 1000
	defaultCacheTTL = 5 * time.Minute
	cleanInterval   = 1 * time.Minute
)

func NewService(
	repository Repository,
	transactor Transactor,
) *Service {
	s := &Service{
		repository: repository,
		transactor: transactor,
		cache:      make(map[string]*cacheEntry),
		lru:        list.New(),
		cacheCap:   defaultCacheCap,
		cacheTTL:   defaultCacheTTL,
		stopClean:  make(chan struct{}),
	}
	go s.backgroundCleaner()
	return s
}

func NewServiceWithCache(repository Repository, transactor Transactor, cap int, ttl time.Duration) *Service {
	s := &Service{
		repository: repository,
		transactor: transactor,
		cache:      make(map[string]*cacheEntry),
		lru:        list.New(),
		cacheCap:   cap,
		cacheTTL:   ttl,
		stopClean:  make(chan struct{}),
	}
	if s.cacheCap <= 0 {
		s.cacheCap = defaultCacheCap
	}
	if s.cacheTTL <= 0 {
		s.cacheTTL = defaultCacheTTL
	}
	go s.backgroundCleaner()
	return s
}

func (s *Service) Close() {
	select {
	case <-s.stopClean:
		return
	default:
		close(s.stopClean)
	}
}

func (s *Service) backgroundCleaner() {
	ticker := time.NewTicker(cleanInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.evictExpired()
		case <-s.stopClean:
			return
		}
	}
}

func (s *Service) evictExpired() {
	now := time.Now()
	s.mu.Lock()
	for e := s.lru.Back(); e != nil; {
		prev := e.Prev()
		ce := e.Value.(*cacheEntry)
		if now.After(ce.expireAt) {
			s.lru.Remove(e)
			delete(s.cache, ce.key)
		}
		e = prev
	}
	s.mu.Unlock()
}

func (s *Service) cachePut(key string, value Model) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ce, ok := s.cache[key]; ok {
		ce.value = value
		ce.expireAt = time.Now().Add(s.cacheTTL)
		s.lru.MoveToFront(ce.element)
		return
	}
	ce := &cacheEntry{key: key, value: value, expireAt: time.Now().Add(s.cacheTTL)}
	ce.element = s.lru.PushFront(ce)
	s.cache[key] = ce
	if s.lru.Len() > s.cacheCap {
		s.evictOldestLocked()
	}
}

func (s *Service) evictOldestLocked() {
	tail := s.lru.Back()
	if tail == nil {
		return
	}
	ce := tail.Value.(*cacheEntry)
	s.lru.Remove(tail)
	delete(s.cache, ce.key)
}

func (s *Service) cacheGet(key string) (Model, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ce, ok := s.cache[key]; ok {
		if time.Now().After(ce.expireAt) {
			s.lru.Remove(ce.element)
			delete(s.cache, key)
			return Model{}, false
		}
		s.lru.MoveToFront(ce.element)
		return ce.value, true
	}
	return Model{}, false
}
