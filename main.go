/*
 * Copyright 2024 Secberus, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package main

import (
	"fmt"
	"os"

	"google.golang.org/grpc"

	service "github.com/secberus/go-push-api/service/v1/push"

	"github.com/secberus/push-api-cli/cli"
	"github.com/secberus/push-api-cli/config"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		fatal(fmt.Errorf("failed to load configuration: %w", err))
	}

	tlsCreds, err := config.Credentials(cfg)
	if err != nil {
		fatal(fmt.Errorf("failed to load client credentials: %w", err))
	}

	conn, err := grpc.NewClient(cfg.Endpoint, grpc.WithTransportCredentials(tlsCreds))
	if err != nil {
		fatal(fmt.Errorf("failed to create client: %w", err))
	}

	client := service.NewPushServiceClient(conn)
	cmd := cli.New(client)
	cmd.Execute()
}

func fatal(err error) {
	fmt.Printf("%s\n", err)
	os.Exit(1)
}
