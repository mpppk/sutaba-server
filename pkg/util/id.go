package util

import (
	"sync"
	"time"
)

type IDMap struct {
	sm                            sync.Map
	expirationLimitSec            time.Duration
	expirationCheckingIntervalSec time.Duration
	stopExpirationCheckChan       chan struct{}
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

func (m *IDMap) StartExpirationCheck() {
	go func() {
		for {
			LogPrintlnInOneLine("expiration checking will be started after ", m.expirationCheckingIntervalSec)
			select {
			case <-time.After(m.expirationCheckingIntervalSec):
				LogPrintlnInOneLine("start expiration checking")
				m.checkExpiration()
				LogPrintlnInOneLine("finish expiration checking")
				continue
			case <-m.stopExpirationCheckChan:
				return
			}
		}
	}()
}

func (m *IDMap) StopExpirationCheck() {
	m.stopExpirationCheckChan <- struct{}{}
}

func (m *IDMap) checkExpiration() {
	now := time.Now()

	// O(N)
	m.sm.Range(func(idInter interface{}, tInter interface{}) bool {
		m.deleteIfExpired(idInter.(int64), tInter.(time.Time), now)
		return true
	})
}

func (m *IDMap) isExpired(t, now time.Time) bool {
	return now.Sub(t) > m.expirationLimitSec
}

func (m *IDMap) deleteIfExpired(id int64, t, now time.Time) bool {
	if m.isExpired(t, now) {
		LogPrintlnInOneLine("delete", id)
		m.sm.Delete(id)
		return true
	}
	return false
}

func NewIDMap(expirationLimitSec, expirationCheckingIntervalSec time.Duration) *IDMap {
	return &IDMap{
		expirationLimitSec:            expirationLimitSec * time.Second,
		expirationCheckingIntervalSec: expirationCheckingIntervalSec * time.Second,
		stopExpirationCheckChan:       make(chan struct{}),
	}
}
