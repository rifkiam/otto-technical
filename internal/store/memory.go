package store

import (
	"be_test/internal/model"
	"fmt"
	"sync"
	"time"
)

// MemoryItemStore adalah penyimpanan in-memory sederhana.
type MemoryItemStore struct {
    items   map[string]model.Item
    mu      sync.RWMutex
    count   int
}

func NewMemoryItemStore() *MemoryItemStore {
    return &MemoryItemStore{items: make(map[string]model.Item)}
}

func (m *MemoryItemStore) Create(name string) (model.Item, error) {
    m.mu.Lock()
    defer m.mu.Unlock()

    m.count++
    id := fmt.Sprintf("item-%d", m.count)
    item := model.Item{
        ID:   id,
        Name: name,
        Done: false,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    m.items[id] = item
    return item, nil
}

func (m *MemoryItemStore) Get(id string) (model.Item, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    item, exists := m.items[id]
    if !exists {
        return model.Item{}, ErrNotFound
    }

    return item, nil
}

func (m *MemoryItemStore) List() ([]model.Item, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    items := make([]model.Item, 0, len(m.items))
    for _, item := range m.items {
        items = append(items, item)
    }
    
    return items, nil
}

func (m *MemoryItemStore) Update(id string, name *string, done *bool) (model.Item, error) {
    m.mu.Lock()
    defer m.mu.Unlock()

    item, exists := m.items[id]
    if !exists {
        return model.Item{}, ErrNotFound
    }

    if name != nil {
        item.Name = *name
    }
    if done != nil {
        item.Done = *done
    }
    item.UpdatedAt = time.Now()

    m.items[id] = item
    return item, nil
}

func (m *MemoryItemStore) Delete(id string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    _, exists := m.items[id]
    if !exists {
        return ErrNotFound
    }

    delete(m.items, id)
    return nil
}