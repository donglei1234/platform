package db

import (
	"github.com/donglei1234/platform/services/common/nosql/memory/keys"
)

func makeChatMsgKey(uid string) (keys.Key, error) {
	return keys.NewKeyFromParts("chat", "message", "queue", uid)
}
