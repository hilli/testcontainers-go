package azurite_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/azurite"
)

func TestAzurite(t *testing.T) {
	ctx := context.Background()
	container, err := azurite.RunContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	t.Run("WithAccount", func(t *testing.T) {
		c, err := azurite.RunContainer(ctx,
			azurite.WithAccount("testaccount", "testkey"),
		)
		if err != nil {
			t.Fatal(err)
		}
		containerInspection, err := c.Inspect(ctx)
		if err != nil {
			t.Fatal(err)
		}

		ok := false
		for _, env := range containerInspection.Config.Env {
			if env == "AZURITE_ACCOUNTS=testaccount:dGVzdGtleQ==" {
				ok = true
				break
			}
		}
		if !ok {
			t.Error("expected account name and key to be set correctly")
		}

	})

	t.Run("ConnectionStringWithDefaultAccount", func(t *testing.T) {
		// connectionString {
		connectionStr, err := container.ConnectionString(ctx)
		// }
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("Connection string:", connectionStr)

		if connectionStr == "" {
			t.Error("expected connection string to be set")
		}

		expectedRegEx := `DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://localhost:\d+/devstoreaccount1;QueueEndpoint=http://localhost:\d+/devstoreaccount1;TableEndpoint=http://localhost:\d+/devstoreaccount1;`
		match, err := regexp.MatchString(expectedRegEx, connectionStr)
		if err != nil {
			t.Fatal(err)
		}

		if !match {
			t.Errorf("expected connection string to match regexp; got %s", connectionStr)
		}
	})
}
