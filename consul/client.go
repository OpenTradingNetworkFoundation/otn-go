package consul

import (
	"github.com/hashicorp/consul/api"
)

func NewClient() (*api.Client, error) {
	return api.NewClient(api.DefaultConfig())
}
