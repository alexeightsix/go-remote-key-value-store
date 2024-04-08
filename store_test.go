package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func StoreFactory(driver int) (*store, error) {
	tmp := os.TempDir()
	file, _ := os.CreateTemp(tmp, "test.db")
	store, err := NewStore(driver, file)
	return store, err
}

func TestArrayStoreHas(t *testing.T) {
	store, _ := StoreFactory(STORE_TYPE_ARRAY)
	assert.False(t, store.has("testing"))
	store.set("testing", "array")
	assert.True(t, store.has("testing"))
}

func TestArrayStoreSet(t *testing.T) {
	store, _ := StoreFactory(STORE_TYPE_ARRAY)
	err := store.set("testing", "array")
	assert.True(t, store.driver.has("testing"))
	assert.Nil(t, err)
}

func TestArrayStoreGet(t *testing.T) {
	store, _ := StoreFactory(STORE_TYPE_ARRAY)
	store.set("testing", "array")
	res, err := store.driver.get("testing")
	assert.Equal(t, res.value, "array")
	assert.Nil(t, err)

	_, err = store.driver.get("something")
	assert.ErrorContains(t, err, "ERROR_KEY_NOT_FOUND")
}

func TestArrayStoreDelete(t *testing.T) {
	store, _ := StoreFactory(STORE_TYPE_ARRAY)
	store.set("testing", "array")
	res, err := store.driver.delete("testing")
	assert.True(t, res)
	assert.Nil(t, err)

	res, err = store.driver.delete("something")
	assert.ErrorContains(t, err, "ERROR_KEY_NOT_FOUND")
}

func TestMapStoreHas(t *testing.T) {
	store, _ := StoreFactory(STORE_TYPE_MAP)
	assert.False(t, store.has("testing"))
	store.set("testing", "array")
	assert.True(t, store.has("testing"))
}

func TestMapStoreSet(t *testing.T) {
	store, _ := StoreFactory(STORE_TYPE_MAP)
	err := store.set("testing", "array")
	assert.True(t, store.driver.has("testing"))
	assert.Nil(t, err)
}

func TestMapStoreGet(t *testing.T) {
	store, _ := StoreFactory(STORE_TYPE_MAP)
	store.set("testing", "array")
	res, err := store.driver.get("testing")
	assert.Equal(t, res.value, "array")
	assert.Nil(t, err)

	_, err = store.driver.get("something")
	assert.ErrorContains(t, err, "ERROR_KEY_NOT_FOUND")
}

func TestMapStoreDelete(t *testing.T) {
	store, _ := StoreFactory(STORE_TYPE_MAP)
	store.set("testing", "array")
	res, err := store.driver.delete("testing")
	assert.True(t, res)
	assert.Nil(t, err)

	res, err = store.driver.delete("something")
	assert.ErrorContains(t, err, "ERROR_KEY_NOT_FOUND")
}

func TestInvalidStore(t *testing.T) {
	_, err := StoreFactory(5)
	assert.ErrorContains(t, err, "ERROR_INVALID_STORE_DRIVER")
}
