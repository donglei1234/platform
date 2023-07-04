package dupblock

import (
	"io"

	"github.com/donglei1234/platform/services/common/jsonx/query"
)

func ApplyDUPBlock(data []byte, dupBlocks []byte) ([]byte, error) {
	cmd := Command{}
	var opts []query.Option

	reader, err := NewTextReader(WithBytes(dupBlocks))
	if err != nil {
		return nil, err
	}

	for {
		if err := reader.Read(&cmd); err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		} else {
			switch cmd.Action {
			case ActionSetKey:
				return nil, ErrUnhandledAction
			case ActionSet:
				opts = append(opts, query.WithSetField([]byte(cmd.To), cmd.Value))
			case ActionInsert:
				opts = append(opts, query.WithArrayInsert([]byte(cmd.To), cmd.Value))
			case ActionIncrement:
				opts = append(opts, query.WithIncrement([]byte(cmd.To), cmd.Delta))
			case ActionPushFront:
				opts = append(opts, query.WithArrayPushFront([]byte(cmd.To), cmd.Value))
			case ActionPushBack:
				opts = append(opts, query.WithArrayPushBack([]byte(cmd.To), cmd.Value))
			case ActionAddUnique:
				opts = append(opts, query.WithArrayUnique([]byte(cmd.To), cmd.Value))
			case ActionDelete:
				opts = append(opts, query.WithDelete([]byte(cmd.To)))
			case ActionCopy:
				opts = append(opts, query.WithCopy([]byte(cmd.From), []byte(cmd.To)))
			case ActionMove:
				opts = append(opts, query.WithMove([]byte(cmd.From), []byte(cmd.To)))
			case ActionSwap:
				opts = append(opts, query.WithSwap([]byte(cmd.From), []byte(cmd.To)))
			case ActionUndefined:
				return nil, ErrUnknownDUPAction
			default:
				return nil, ErrUnhandledAction
			}
		}
	}

	opts = append(opts, query.WithInput(data))
	output, err := query.Execute(opts...)
	if err != nil {
		return nil, err
	}

	return output, nil
}
