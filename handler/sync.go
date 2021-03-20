package handler

import (
	"net/http"
	"sync"
)

func SyncWrite(h http.Handler, lock *sync.RWMutex) http.Handler { //nolint: interfacer
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lock.Lock()
		defer lock.Unlock()
		h.ServeHTTP(w, r)
	})
}

func SyncRead(h http.Handler, lock *sync.RWMutex) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lock.RLock()
		defer lock.RUnlock()
		h.ServeHTTP(w, r)
	})
}
