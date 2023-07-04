package db

import (
	"github.com/donglei1234/platform/services/common/nosql/memory/keys"
	pb "github.com/donglei1234/platform/services/condition/gen/condition/api"
	"strconv"
)

const (
	ConditionTheme = "condition"
)

func makeConditionKey(uid string) (keys.Key, error) {
	return keys.NewKeyFromParts(ConditionTheme, uid)
}
func makeConditionFieldKey(cond *pb.Condition) (keys.Key, error) {
	tpStr := "tp" + strconv.Itoa(int(cond.Type))
	param0Str := "param0" + strconv.Itoa(int(cond.Params[0]))
	param1Str := "param1" + strconv.Itoa(int(cond.Params[1]))
	return keys.NewKeyFromParts(ConditionTheme, tpStr, param0Str, param1Str)
}
