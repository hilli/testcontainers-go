package azurite

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	DefaultImage       = "mcr.microsoft.com/azure-storage/azurite:3.23.0"
	defaultAccountName = "devstoreaccount1"
	defaultAccountKey  = "Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw=="
	BlobPort           = "10000"
	QueryPort          = "10001"
	TablePort          = "10002"
)

// DoltContainer represents the Dolt container type used in the module
type AzuriteContainer struct {
	testcontainers.Container
	AccountName string
	AccountKey  string
}

func WithDefaultCredentials() testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.Env["AZURITE_ACCOUNTS"] = fmt.Sprintf("%s:%s", defaultAccountName, defaultAccountKey)
		// req.Env["AZURITE_ACCOUNT_NAME"] = defaultAccountName
		// req.Env["AZURITE_ACCOUNT_KEY"] = defaultAccountKey
		return nil
	}
}

// WithAccount sets the account name and key for the Azurite container
func WithAccount(name, key string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.Env["AZURITE_ACCOUNTS"] = fmt.Sprintf("%s:%s", name, base64.StdEncoding.EncodeToString([]byte(key)))
		req.Env["AZURITE_ACCOUNT_NAME"] = name
		req.Env["AZURITE_ACCOUNT_KEY"] = base64.StdEncoding.EncodeToString([]byte(key))
		return nil
	}
}

// WithStorageLocation sets the storage location for the Azurite container
// TODO: Implement this function to set the storage location for the Azurite container
func WithStorageLocation(location string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.Env["STORAGE_LOCATION"] = location
		return nil
	}
}

// RunContainer creates an instance of the Azurite container type
func RunContainer(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*AzuriteContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        DefaultImage,
		ExposedPorts: []string{"10000/tcp", "10001/tcp", "10002/tcp"},
		Env: map[string]string{
			"AZURITE_ACCOUNTS":     fmt.Sprintf("%s:%s", defaultAccountName, defaultAccountKey),
			"AZURITE_ACCOUNT_NAME": defaultAccountName,
			"AZURITE_ACCOUNT_KEY":  defaultAccountKey,
		},
		WaitingFor: wait.ForLog("Azurite Blob service is successfully listening at"),
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	// opts = append(opts, WithDefaultCredentials())

	for _, opt := range opts {
		opt.Customize(&genericContainerReq)
	}

	accountName := req.Env["AZURITE_ACCOUNT_NAME"]
	accountKey := req.Env["AZURITE_ACCOUNT_KEY"]

	c, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	if err != nil {
		return nil, err
	}

	return &AzuriteContainer{
		Container:   c,
		AccountName: accountName,
		AccountKey:  accountKey,
	}, nil
}

func (c *AzuriteContainer) ConnectionString(ctx context.Context) (string, error) {
	blobPort, err := c.MappedPort(ctx, "10000/tcp")
	if err != nil {
		return "", err
	}

	queuePort, err := c.MappedPort(ctx, "10001/tcp")
	if err != nil {
		return "", err
	}

	tablePort, err := c.MappedPort(ctx, "10002/tcp")
	if err != nil {
		return "", err
	}

	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}

	connStrFormat := "DefaultEndpointsProtocol=http;AccountName=%s;AccountKey=%s;BlobEndpoint=http://%s:%s/%s;QueueEndpoint=http://%s:%s/%s;TableEndpoint=http://%s:%s/%s;"
	// return fmt.Sprintf("UseDevelopmentStorage=true;DevelopmentStorageProxyUri=http://%s:%s", host, blobPort.Port()), nil
	return fmt.Sprintf(connStrFormat, c.AccountName, c.AccountKey, host, blobPort.Port(), c.AccountName, host, queuePort.Port(), c.AccountName, host, tablePort.Port(), c.AccountName), nil
}
