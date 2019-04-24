package secrets_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/opentradingnetworkfoundation/otn-go/secrets"
)

func TestRead(t *testing.T) {
	a := assert.New(t)
	cfg := &secrets.VaultStorageConfig{Approle: "otn"}
	st, err := secrets.NewVaultStorage(cfg)
	a.NoError(err)

	secret, err := st.ReadSecret("otn/otn-faucet/main/keys")
	a.NoError(err)
	a.NotNil(secret)

	t.Logf("Secret: %#v", secret)

	keys, err := st.ReadStringArray("otn/otn-faucet/main/keys")
	a.NoError(err)
	t.Logf("Keys: %#v", keys)

	a.Fail("OOps")
}
