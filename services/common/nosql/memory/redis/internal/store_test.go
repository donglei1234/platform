package internal

import (
	"testing"

	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/nosql/pkg/tests/common_tests"
)

const (
	RedisUrl = "10.0.1.3:6378"
)

func TestNewMemoryStoreProvider(t *testing.T) {
	testStoreName := "testStore"
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	if provider, err := NewMemoryStoreProvider(RedisUrl, "", logger); err != nil {
		t.Fatal("Error encountered creating a new redis store provider:", err)
	} else {
		if store, err := provider.OpenMemoryStore(testStoreName); err != nil {
			t.Fatal("Error opening redis store", testStoreName, ":", err)
		} else if store.Name() != testStoreName {
			t.Fatal("Created store does not use the provided redis store name.")
		} else if err := common_tests.MemoryStoreCommonTest(logger, store); err != nil {
			t.Fatal(err)
		} else if err := provider.Shutdown(); err != nil {
			t.Fatal("Error encountered shutting down the redis store provider:", err)
		}
	}
}
