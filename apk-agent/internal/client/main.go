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
	"fmt"
	"log"
	"net"

	apiProtos "github.com/AmaliMatharaarachchi/APKAgent/apk-agent/internal/client/grpc/api"
	service "github.com/AmaliMatharaarachchi/APKAgent/apk-agent/internal/client/mgt-server-mock"

	"google.golang.org/grpc"
)

func main() {

	s := grpc.NewServer()

	apiService := service.NewApiService();
	apiProtos.RegisterAPIServiceServer(s, apiService)

	tl, err := net.Listen("tcp", "localhost:8765")
	if err != nil {
		log.Fatal(fmt.Println("Error starting tcp listener on port 8765", err))
	}

	// start listening
	s.Serve(tl)
}
