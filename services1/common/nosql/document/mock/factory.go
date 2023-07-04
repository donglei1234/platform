package mock

import (
	"github.com/donglei1234/platform/services/common/nosql/document"
	"github.com/donglei1234/platform/services/common/nosql/document/mock/internal"
)

type DocumentStoreProvider = internal.DocumentStoreProvider

func NewDocumentStoreProvider() (document.DocumentStoreProvider, error) {
	return internal.NewDocumentStoreProvider()
}
