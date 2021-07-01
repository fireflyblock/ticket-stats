package client

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/filecoin-project/lotus/api"
	lotusClient "github.com/filecoin-project/lotus/api/client"
)

var Client api.FullNode
var SignClient api.FullNode

func CreateLotusClient() error {
	host := os.Getenv("LOTUS_HOST")
	var err error
	requestHeader := http.Header{}
	requestHeader.Add("Content-Type", "application/json")
	Client, _, err = lotusClient.NewFullNodeRPCV1(context.Background(), host, requestHeader)
	if err != nil {
		return err
	}
	return nil
}
func CreateLotusSignClient() error {
	host := os.Getenv("LOTUS_HOST")
	token := os.Getenv("LOTUS_SIGN_TOKEN")

	var err error
	requestHeader := http.Header{}
	requestHeader.Add("Content-Type", "application/json")
	tokenHeader := fmt.Sprintf("Bearer %s", token)
	requestHeader.Set("Authorization", tokenHeader)
	SignClient, _, err = lotusClient.NewFullNodeRPCV1(context.Background(), host, requestHeader)
	if err != nil {
		return err
	}
	return nil
}
