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
package config

import (
	"encoding/json"
	"os"
)

const (
	DefaultEndpoint        = "push.secberus.io:7744"
	DefaultCredentialsFile = "$HOME/.s6s/config"
	CredentialsFileEnvVar  = "S6S_CREDENTIALS_FILE"
)

type Config struct {
	Endpoint        string `json:"endpoint"`
	X509Certificate []byte `json:"x509_certificate"`
	PrivateKey      []byte `json:"private_key"`
	CaBundle        []byte `json:"ca_bundle"`
}

func Load() (*Config, error) {

	var credsFile string
	credsFile, ok := os.LookupEnv(CredentialsFileEnvVar)
	if !ok {
		credsFile = DefaultCredentialsFile
	}

	raw, err := os.ReadFile(os.ExpandEnv(credsFile))
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	cfg.Endpoint = DefaultEndpoint

	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
