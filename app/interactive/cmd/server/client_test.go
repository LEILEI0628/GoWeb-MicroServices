package main

import (
	"context"
	itrAPI "github.com/LEILEI0628/GoWeb-MicroServices/api/interactive/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

func TestGRPCClient(t *testing.T) {
	cc, err := grpc.Dial("localhost:8000",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	client := itrAPI.NewInteractiveServiceClient(cc)
	resp, err := client.Get(context.Background(), &itrAPI.GetRequest{
		Biz:   "test",
		BizId: 2,
		Uid:   345,
	})
	require.NoError(t, err)
	t.Log(resp.Intr)
}
