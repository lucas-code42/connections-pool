package main

import (
	"database/sql"
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestStartConnectionsPool(t *testing.T) {
	StartConnectionsPool()

	assert.NotNil(t, PG_POOL, "PG_POOL can't be nil")
	assert.Equal(t, MAX_POOL, len(PG_POOL.Connections), "The len of pool should be equal to MAX_POOL")

	for _, conn := range PG_POOL.Connections {
		assert.NotNil(t, conn, "each connection cant't be nil")
	}
}

func TestGetFreeConnection(t *testing.T) {
	mockPool := &ConnectionPool{
		Connections: make([]*sql.DB, 3),
		InUse:       make(map[*sql.DB]bool),
		Mutex:       sync.Mutex{},
	}

	for i := 0; i < 3; i++ {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error to create mock database: %v", err)
		}
		mockPool.Connections[i] = mockDB
		mockPool.InUse[mockDB] = false
	}

	// Definir uma conexão como em uso
	mockPool.InUse[mockPool.Connections[1]] = true
	mockPool.InUse[mockPool.Connections[2]] = true

	freeConn := mockPool.GetFreeConnection()
	assert.Equal(t, mockPool.Connections[0], freeConn, "should return the first free connection")
}

func TestReleaseConnection(t *testing.T) {
	mockPool := &ConnectionPool{
		Connections: make([]*sql.DB, 3),
		InUse:       make(map[*sql.DB]bool),
		Mutex:       sync.Mutex{},
	}

	for i := 0; i < 3; i++ {
		// Criar um mock de *sql.DB usando sqlmock
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error to create mock database: %v", err)
		}
		mockPool.Connections[i] = mockDB
		mockPool.InUse[mockDB] = false
	}

	// Definir uma conexão como em uso
	mockPool.InUse[mockPool.Connections[0]] = true

	mockPool.ReleaseConnection(mockPool.Connections[0])
	assert.False(t, mockPool.InUse[mockPool.Connections[0]], "connection should be free")
}
