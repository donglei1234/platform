package db

import (
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/nosql/document"
)

type Document struct {
	document.DocumentStore
	logger *zap.Logger
}

func OpenDatabase(l *zap.Logger, doc document.DocumentStore) *Document {
	return &Document{
		doc,
		l,
	}
}
