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
package main

import (
	"context"
	"flag"
	"log"

	// "google.golang.org/grpc"
	clientPool "github.com/AmaliMatharaarachchi/APKAgent/apk-agent/internal/client/client-pool"
	apiProtos "github.com/AmaliMatharaarachchi/APKAgent/apk-agent/internal/client/grpc/api"
	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", "localhost:8765", "the address to connect to")
	// see https://github.com/grpc/grpc/blob/master/doc/service_config.md to know more about service config
	retryPolicy = `{
		"methodConfig": [{
		  "name": [{"service": "wso2.agent.api.APIService"}],
		  "waitForReady": true,
		  "retryPolicy": {
			  "MaxAttempts": 5,
			  "InitialBackoff": "1s",
			  "MaxBackoff": "1000s",
			  "BackoffMultiplier": 2.0,
			  "RetryableStatusCodes": [ "UNAVAILABLE" ]
		  }
		}]}`
)

func main() {
	pool, err := clientPool.Init("localhost:8765", 5, 3, []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock()}, 
		clientPool.RetryPolicy{
			MaxAttempts : 500,
			BackOffInMilliSeconds : 1000,
			RetryableStatuses : []string{},
		})
	if err != nil {
		log.Println("init client pool error", err)
		return
	}
	clientConn, err := pool.GetConnection();
	if err != nil {
		log.Println("init client pool error", err)
		return
	}
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	client := apiProtos.NewAPIServiceClient(clientConn);
	// defer cancel()

	response, err := pool.ExecuteGRPCCall(clientConn, func () (interface{}, error) {
		return client.CreateAPI(context.Background(), &apiProtos.API{
			ApiUUID: "test",
			ApiVersion: "1.0.0",
		})
	})
	if err != nil {
		log.Fatalf("Create API failed: %v", err)
	}
	log.Printf("%q",response)
	
	
}
