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
	"crypto/tls"
	"crypto/x509"
  "errors"
	"google.golang.org/grpc/credentials"
)

func Credentials(cfg *Config) (credentials.TransportCredentials, error) {

  cert, err := tls.X509KeyPair(cfg.X509Certificate, cfg.RsaKey)
  if err != nil {
    return nil, err
  }

  certPool := x509.NewCertPool()
  if !certPool.AppendCertsFromPEM(cfg.CaBundle) {
    return nil, errors.New("failed to parse CA certificates")
  }

  tlsConfig := &tls.Config{
    Certificates: []tls.Certificate{cert},
    RootCAs:      certPool,
  }
  return credentials.NewTLS(tlsConfig), nil
}
