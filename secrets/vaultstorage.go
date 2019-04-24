package secrets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	consul "github.com/hashicorp/consul/api"
	vault "github.com/hashicorp/vault/api"
	"github.com/juju/errors"
	"github.com/tidwall/gjson"
)

type (
	ConsulClient = consul.Client
	VaultClient  = vault.Client
	SecretData   = map[string]interface{}
)

type VaultStorageConfig struct {
	Address string `json:"address"`
	Approle string `json:"approle"`
	Path    string `json:"path"`
}

type VaultStorage struct {
	cfg    *vault.Config
	vault  *VaultClient
	prefix string
}

func NewVaultStorage(cfg *VaultStorageConfig) (*VaultStorage, error) {
	vaultCfg := vault.DefaultConfig()
	if cfg.Address != "" {
		vaultCfg.Address = cfg.Address
	}

	client, err := vault.NewClient(vaultCfg)
	if err != nil {
		return nil, errors.Annotate(err, "Failed to create Vault client")
	}

	stg := &VaultStorage{
		cfg:    vaultCfg,
		vault:  client,
		prefix: cfg.Path,
	}

	if cfg.Approle != "" {
		if err := stg.authenticate(cfg.Approle, nil); err != nil {
			return nil, err
		}
	}

	return stg, nil
}

func (s *VaultStorage) authenticate(approle string, consulClient *ConsulClient) error {
	if consulClient == nil {
		consulCfg := consul.DefaultConfig()
		c, err := consul.NewClient(consulCfg)
		if err != nil {
			return err
		}
		consulClient = c
	}

	kv := consulClient.KV()

	roleID, _, err := kv.Get(fmt.Sprintf("approle/%s/role_id", approle), nil)
	if err != nil {
		return err
	}

	secretID, _, err := kv.Get(fmt.Sprintf("approle/%s/secret_id", approle), nil)
	if err != nil {
		return err
	}

	if roleID == nil || secretID == nil {
		return errors.Errorf("approle '%s' not found", approle)
	}

	request := s.vault.NewRequest("POST", "/v1/auth/approle/login")
	if err := request.SetJSONBody(map[string]string{
		"role_id":   string(roleID.Value),
		"secret_id": string(secretID.Value),
	}); err != nil {
		return err
	}

	response, err := s.vault.RawRequest(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	token := gjson.GetBytes(raw, "auth.client_token").String()
	s.vault.SetToken(token)

	return nil
}

func (s *VaultStorage) ReadSecret(path string) (SecretData, error) {
	secret, err := s.vault.Logical().Read("secret/" + s.prefix + path)
	if err != nil {
		return nil, err
	}

	if secret != nil {
		return secret.Data, nil
	}

	return nil, nil
}

func (s *VaultStorage) ReadStringValue(path string) (string, error) {
	data, err := s.ReadSecret(path)
	if err != nil {
		return "", err
	}

	value, ok := data["value"]
	if !ok {
		return "", errors.NotFoundf("Value for key %s not found", path)
	}

	str, ok := value.(string)
	if !ok {
		return "", errors.NotFoundf("Value for key %s is not string", path)
	}

	return str, nil
}

func (s *VaultStorage) ReadStringArray(path string) ([]string, error) {
	str, err := s.ReadStringValue(path)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0)
	json.Unmarshal([]byte(str), &keys)

	return keys, nil
}
