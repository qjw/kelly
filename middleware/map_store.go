// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package middleware

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sync"
	"time"
)

const (
	defaultCleaningPeriod = 15 * time.Minute
)

// A very simple implementation of a key value store (a concurrent safe map)
type mapStore struct {
	*sync.RWMutex
	data    map[string][]byte
	binding FreshnessInformer
}

type FreshnessInformer interface {
	IsFresh() bool
}

type mapStoreConfig struct {
	// The period to clean the store in
	CleaningPeriod time.Duration
}

// Error Type for the key value store
type MapStoreError string

// The Error for Key Value Store
func (err MapStoreError) Error() string {
	return "Throttle Map Store Error: " + string(err)
}

// Set a key
func (s *mapStore) Set(key string, value []byte) error {
	s.Lock()
	s.data[key] = value
	s.Unlock()

	return nil
}

// Delete a key
func (s *mapStore) Delete(key string) {
	s.Lock()
	delete(s.data, key)
	s.Unlock()
}

// Get a key, will return an error if the key does not exist
func (s *mapStore) Get(key string) (value []byte, err error) {
	s.RLock()
	value, ok := s.data[key]
	s.RUnlock()
	if !ok {
		err = MapStoreError("Key " + key + " does not exist")
		return value, err
	} else {
		return value, nil
	}
}

// Read the data into the given binding
func (s *mapStore) Read(key string) (FreshnessInformer, error) {
	byteArray, err := s.Get(key)
	if err != nil {
		return nil, err
	}

	byteBufferString := bytes.NewBuffer(byteArray)
	var arbitraryStructure interface{}
	if err := json.NewDecoder(byteBufferString).Decode(&arbitraryStructure); err != nil {
		return nil, err
	}

	for k, v := range arbitraryStructure.(map[string]interface{}) {
		if field := reflect.ValueOf(s.binding).FieldByName(k); field.IsValid() && field.CanSet() {
			field.Set(reflect.ValueOf(v))
		}
	}

	return s.binding, err
}

// Clean the store from expired values
func (s *mapStore) Clean() {
	for key := range s.data {
		value, err := s.Read(key)
		if err == nil && !value.IsFresh() {
			s.Delete(key)
		} else if err != nil {
			panic(err)
		}
	}
}

// Simple cleanup mechanism, cleaning the store every 15 minutes
func (s *mapStore) CleanEvery(cleaningPeriod time.Duration) {
	c := time.Tick(cleaningPeriod)

	for {
		select {
		case <-c:
			s.Clean()
		}
	}
}

// Returns a simple key value store
func newMapStore(binding FreshnessInformer, options ...*mapStoreConfig) *mapStore {
	s := &mapStore{
		&sync.RWMutex{},
		make(map[string][]byte),
		binding,
	}

	o := newMapStoreOptions(options)

	go s.CleanEvery(o.CleaningPeriod)

	return s
}

// Returns new map store options from defaults and given options
func newMapStoreOptions(options []*mapStoreConfig) *mapStoreConfig {
	o := &mapStoreConfig{
		defaultCleaningPeriod,
	}

	if len(options) == 0 {
		return o
	}

	if options[0].CleaningPeriod != 0 {
		o.CleaningPeriod = options[0].CleaningPeriod
	}

	return o
}
