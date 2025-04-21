package kvstore

import "sync"

type KVStore struct {
	mu    sync.RWMutex
	store map[string]string
}

func New() *KVStore {
	return &KVStore{
		store: make(map[string]string),
	}
}

func (kv *KVStore) Set(key, value string) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.store[key] = value
}

func (kv *KVStore) Get(key string) (string, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	value, exists := kv.store[key]
	return value, exists
}

func (kv *KVStore) Delete(key string) bool {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	_, exists := kv.store[key]
	if exists {
		delete(kv.store, key)
	}
	return exists
}

func (kv *KVStore) Close() error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.store = nil
	return nil
}
