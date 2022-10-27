/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"context"
	"runtime/debug"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

func connect(addr string, dynamicLink bool, ops ...grpc.DialOption) (*grpc.ClientConn, error) {
	if dynamicLink == true {
		return grpc.Dial(addr, ops...)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	return grpc.DialContext(ctx, addr, ops...)
}

func (pool *Pool) init() error {
	len := cap(pool.connections) - len(pool.connections)
	addr := pool.ServerAddr
	dynamicLink := pool.DynamicLink
	ops := pool.dialOptions

	for i := 1; i <= len; i++ {
		client, err := connect(addr, dynamicLink, ops...)
		if err != nil {
			return err
		}
		pool.connections <- client
	}

	return nil
}

func (pool *Pool) count(add int64) {
	atomic.AddInt64(&pool.counter, add)
}

// close Recycling available links
func (pool *Pool) close(conn *grpc.ClientConn) {
	// double check
	if conn == nil {
		return
	}

	go func() {
		defer func() {
			if err := recover(); nil != err {
				debug.PrintStack()
			}
		}()
		detect, _ := passivate(conn)
		if detect && pool.ChannelStat {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			select {
			case pool.connections <- conn:
			case <-ctx.Done():
				destroy(conn)
			}
		}
		pool.count(-1)
	}()
}

// destroy tears down the ClientConn and all underlying connections.
func destroy(conn *grpc.ClientConn) error {
	return conn.Close()
}

type condition = int

const (
	// Ready Can be used
	Ready condition = iota
	// Put Not available. Maybe later.
	Put
	// Destroy Failure occurs and cannot be restored
	Destroy
)

// passivate Action before releasing the resource
func passivate(conn *grpc.ClientConn) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if conn.WaitForStateChange(ctx, connectivity.Ready) && conn.WaitForStateChange(ctx, connectivity.Shutdown) && conn.WaitForStateChange(ctx, connectivity.Idle) {
		return true, nil
	}

	return false, destroy(conn)
}

// activate Action taken after getting the resource
func activate(conn *grpc.ClientConn) int {
	stat := conn.GetState()
	switch {
	case stat == connectivity.Ready:
		return Ready
	case stat == connectivity.Shutdown:
		return Destroy
	default:
		return Put
	}
}