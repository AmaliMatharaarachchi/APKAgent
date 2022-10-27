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
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
)

var (
	errAcquire = errors.New("acquire connection timed out, you can fix this error by setting the overflow cap or increasing the maximum capacity of the cap")
	errTimeout = errors.New("connect timed out, check the address configuration or network status")
)

//CloseFunc should defer
type CloseFunc func()

func NewRpcClientPool(ops ...Option) (*Pool, error) {
	pool := &Pool{
		MaxCap:         10,
		AcquireTimeout: 3 * time.Second,
		DynamicLink:    false,
		OverflowCap:    true,
		dialOptions:    []grpc.DialOption{grpc.WithInsecure()},
		lock:           &sync.Mutex{},
		counter:        0,
		ChannelStat:    true,
	}

	// Loop through each option
	for _, opt := range ops {
		// Call the option giving the instantiated
		opt(pool)
	}

	pool.connections = make(chan *grpc.ClientConn, pool.MaxCap)

	if err := pool.init(); err != nil {
		return nil, err
	}

	return pool, nil
}

func (pool *Pool) Acquire() (*grpc.ClientConn, CloseFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pool.AcquireTimeout)
	defer cancel()

	for {
		select {
		case clientConn := <-pool.connections:
			con := activate(clientConn)
			switch con {
			case Ready:
				pool.count(1)
				return clientConn, func() { pool.close(clientConn) }, nil
			case Put:
				pool.close(clientConn)
				continue
			default:
				pool.count(-1)
				destroy(clientConn)
				continue
			}
		case <-ctx.Done():
			return nil, nil, errAcquire
		default:
			if pool.OverflowCap == false && pool.counter >= pool.MaxCap {
				continue
			} else {
				addr := pool.ServerAddr
				dynamicLink := pool.DynamicLink
				ops := append(pool.dialOptions, grpc.WithBlock())
				clientConn, err := connect(addr, dynamicLink, ops...)
				if err != nil {
					if err == context.DeadlineExceeded {
						return nil, nil, errTimeout
					}
					return nil, nil, err
				}
				pool.count(1)
				return clientConn, func() { pool.close(clientConn) }, nil
			}
		}
	}
}

// GetStat Return to the use of resources in the pool
func (pool *Pool) GetStat() (used int64, surplus int) {
	return atomic.LoadInt64(&pool.counter), len(pool.connections)
}

// ClearPool Disconnect the link before exit the program
func (pool *Pool) ClearPool() {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	pool.ChannelStat = false
	close(pool.connections)

	for client := range pool.connections {
		destroy(client)
	}
}
