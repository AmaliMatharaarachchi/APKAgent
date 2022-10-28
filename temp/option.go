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
	"sync"
	"time"

	"google.golang.org/grpc"
)

type chanStat = bool

type Pool struct {
	ServerAddr     string
	MaxCap         int64
	AcquireTimeout time.Duration
	DynamicLink    bool
	OverflowCap    bool
	dialOptions    []grpc.DialOption

	lock        *sync.Mutex
	connections chan *grpc.ClientConn
	ChannelStat chanStat
	counter     int64
}

type Option func(*Pool)

func WithMaxCap(num int64) Option {
	return func(pool *Pool) {
		pool.MaxCap = num
	}
}

func WithServerAddr(addr string) Option {
	return func(pool *Pool) {
		pool.ServerAddr = addr
	}
}

func WithAcquireTimeOut(timeout time.Duration) Option {
	return func(pool *Pool) {
		pool.AcquireTimeout = timeout
	}
}

func WithOverflowCap(overflowCap bool) Option {
	return func(pool *Pool) {
		pool.OverflowCap = overflowCap
	}
}

func WithDynamicLink(dynamicLink bool) Option {
	return func(pool *Pool) {
		pool.DynamicLink = dynamicLink
	}
}

func WithDialOption(ops ...grpc.DialOption) Option {
	return func(pool *Pool) {
		pool.dialOptions = append(pool.dialOptions, ops...);
	}
}


