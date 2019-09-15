package util

import (
	"sync"
	"time"
)

type IDMap struct {
	sm sync.Map
}

func (m *IDMap) Load(id int64) (time.Time, bool) {
	t, ok := m.sm.Load(id)
	if !ok {
		return time.Time{}, false
	}
	return t.(time.Time), true
}

func (m *IDMap) Store(id int64) {
	m.sm.Store(id, time.Now())
}

func (m *IDMap) LoadOrStore(id int64) (time.Time, bool) {
	t, ok := m.sm.LoadOrStore(id, time.Now())
	if !ok {
		return time.Time{}, false
	}
	return t.(time.Time), true
}

func NewIDMap() *IDMap {
	return &IDMap{}
}
