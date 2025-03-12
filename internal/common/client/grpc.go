package client

import (
	"context"
	"github.com/FacundoChan/gorder-v1/common/discovery"
	"github.com/FacundoChan/gorder-v1/common/genproto/stockpb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewStockGRPCClient(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("stock.service-name"))
	if err != nil {
		return nil, func() error { return nil }, err
	}
	if grpcAddr == "" {
		logrus.Warn("no stock service gRPC address found")
	}
	grpcDialOptions, err := grpcDialOpts(grpcAddr)
	if err != nil {
		logrus.Error("grpcDialOpts failed: ", err)
		return nil, func() error { return nil }, err
	}
	conn, err := grpc.NewClient(grpcAddr, grpcDialOptions...)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	return stockpb.NewStockServiceClient(conn), conn.Close, nil
}

func grpcDialOpts(addr string) ([]grpc.DialOption, error) {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}, nil
}
