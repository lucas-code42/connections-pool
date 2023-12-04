package main

import (
	"database/sql"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	MAX_POOL  = 3
	MAX_TRIES = MAX_POOL * 2
)

type ConnectionPool struct {
	Connections []*sql.DB          // [dbx001, dbx002, dbx003]
	InUse       map[*sql.DB]bool   // {dbx001: false, dbx002: false, dbx003:false}
	Mutex       sync.Mutex
}

// PG_POOL exportable variable
var PG_POOL *ConnectionPool

// StartConnectionsPool create's the struct tha provides a pool of connections
func StartConnectionsPool() {
	pool := &ConnectionPool{
		Connections: make([]*sql.DB, MAX_POOL), // [dbx001, dbx002, dbx003]
		InUse:       make(map[*sql.DB]bool),    // {dbx001: false, dbx002: false, dbx003:false}
		Mutex:       sync.Mutex{},              // mutex
	}

	for i := 0; i < MAX_POOL; i++ {
		db, err := sql.Open("sqlite3", ".foo.db")
		if err != nil {
			log.Fatal(err)
		}
		pool.Connections[i] = db
		pool.InUse[db] = false
	}

	PG_POOL = pool
}

// GetFreeConnection gets a free connection,
// if any connection is not free wait 100 ms until a free connection returns
func (c *ConnectionPool) GetFreeConnection() *sql.DB {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	for i := 0; i < MAX_TRIES; i++ {
		for _, conn := range c.Connections {
			if !c.InUse[conn] {
				c.InUse[conn] = true
				return conn
			}
		}
		time.Sleep(time.Millisecond * 100)
	}

	panic("error to get a free connection")
}

// ReleaseConnection change to false the attr InUse
func (c *ConnectionPool) ReleaseConnection(conn *sql.DB) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if c.InUse[conn] {
		c.InUse[conn] = false
	}
}
