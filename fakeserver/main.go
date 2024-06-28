// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"math/big"
	"net"
	"strings"
	"time"

	configv1 "github.com/illumio/terraform-provider-illumio-cloudsecure/api/illumio/cloud/config/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	AuthorizationHeader = "authorization"
	DefaultToken        = "token1"
	DefaultClientID     = "client_id_1"
	DefaultClientSecret = "client_secret_1"
)

func tokenAuthInterceptor(token string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		authHeader, ok := md[AuthorizationHeader]

		authHeader = strings.Split(authHeader[0], " ")
		if !ok || authHeader[0] != "Bearer" || len(authHeader) != 2 || authHeader[1] != token {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		return handler(ctx, req)
	}
}

func generateSelfSignedCert() (tls.Certificate, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour) // 1 year

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).SetInt64(1<<62))
	if err != nil {
		return tls.Certificate{}, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Illumio-CloudSecure"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	cert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	key := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	tlsCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return tls.Certificate{}, err
	}

	return tlsCert, nil
}

func main() {
	var (
		// debug enables debug logging if true.
		debug bool

		// network is the network used to access the gRPC server.
		network string

		// token is the token used for auth by fake server.
		token string

		// apiEndpoint is the address of the Config API endpoint.
		apiEndpoint string

		// tokenEndpoint is the address of the OAuth 2 Token endpoint.
		tokenEndpoint string
	)

	flag.BoolVar(&debug, "debug", false, "enables debug logging")
	flag.StringVar(&network, "network", "tcp", "network of the address of the gRPC server, e.g., \"tcp\" or \"unix\"")
	flag.StringVar(&token, "token", "", "token used for auth by fake server. If not provided, a default token will be used.")
	flag.StringVar(&apiEndpoint, "apiEndpoint", "127.0.0.1:50123", "address of the Config API endpoint")
	flag.StringVar(&tokenEndpoint, "tokenEndpoint", "127.0.0.1:50124", "address of the OAuth 2 Token endpoint")
	flag.Parse()

	var logger *zap.Logger

	var err error

	if debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if token == "" {
		token = DefaultToken
	}

	if err != nil {
		panic(fmt.Sprintf("failed to configure logger: %s", err))
	}

	listener, err := net.Listen(network, apiEndpoint)
	if err != nil {
		logger.Fatal("failed to open network port", zap.Error(err))
	}

	cert, err := generateSelfSignedCert()
	if err != nil {
		logger.Fatal("failed to generate self-signed cert", zap.Error(err))
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	})
	server := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(tokenAuthInterceptor(token)),
	)
	configv1.RegisterConfigServiceServer(server, NewFakeConfigServer(logger))
	logger.Info("api server listening", zap.String("network", listener.Addr().Network()), zap.String("address", listener.Addr().String()))
	reflection.Register(server)

	go func() {
		startHTTPServer(
			tokenEndpoint,
			logger,
			DefaultClientID,
			DefaultClientSecret,
			token,
			cert,
		)
	}()
	logger.Info("token endpoint listening", zap.String("address", tokenEndpoint))

	if err := server.Serve(listener); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}
}
