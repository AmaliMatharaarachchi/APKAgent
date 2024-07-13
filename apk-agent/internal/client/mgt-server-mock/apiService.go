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

package service

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	apiProtos "github.com/AmaliMatharaarachchi/APKAgent/apk-agent/internal/client/grpc/api"
)

type aPIService struct {
	apiProtos.UnimplementedAPIServiceServer
}

func NewApiService() *aPIService {
	return &aPIService{}
}

func (s *aPIService) CreateAPI(ctx context.Context, api *apiProtos.API ) (*apiProtos.Response, error) {
	log.Printf("%q", api);
	// No feature was found, return an unnamed feature
	// return &apiProtos.Response{Result : true}, nil
	return nil, status.Errorf(codes.Unavailable, "maybeFailRequest: failing it")
}

func (s *aPIService) UpdateAPI(ctx context.Context, api *apiProtos.API ) (*apiProtos.Response, error) {
	log.Printf("%q", api);
	// No feature was found, return an unnamed feature
	return &apiProtos.Response{Result : true}, nil
}

func (s *aPIService) DeleteAPI(ctx context.Context, api *apiProtos.API ) (*apiProtos.Response, error) {
	log.Printf("%q", api);
	// No feature was found, return an unnamed feature
	return &apiProtos.Response{Result : true}, nil
}
