// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"
	"fmt"
	"net"

	configv1 "github.com/illumio/terraform-provider-illumio-cloudsecure/api/illumio/cloud/config/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	// debug enables debug logging if true
	debug bool

	// network is the network used to access the gRPC server.
	network string

	// address is the gRPC server address.
	address string
)

func init() {
	flag.BoolVar(&debug, "debug", false, "enables debug logging")
	flag.StringVar(&network, "network", "tcp", "network of the address of the gRPC server, e.g., \"tcp\" or \"unix\"")
	flag.StringVar(&address, "address", "127.0.0.1:50123", "address of the gRPC server")
}

func main() {
	flag.Parse()

	var logger *zap.Logger
	var err error
	if debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(fmt.Sprintf("failed to configure logger: %s", err))
	}

	listener, err := net.Listen(network, address)
	if err != nil {
		logger.Fatal("failed to open network port", zap.Error(err))
	}

	server := grpc.NewServer()
	configv1.RegisterConfigServiceServer(server, NewFakeConfigServer(logger))
	logger.Info("server listening", zap.String("network", listener.Addr().Network()), zap.String("address", listener.Addr().String()))
	if err = server.Serve(listener); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}
}
