package dupblock

type Action struct {
	ID   int
	Text string
}

var (
	ActionUndefined = Action{0, "???"}
	ActionSet       = Action{1, "set"}
	ActionInsert    = Action{2, "ins"}
	ActionIncrement = Action{3, "inc"}
	ActionPushFront = Action{4, "puf"}
	ActionPushBack  = Action{5, "pub"}
	ActionAddUnique = Action{6, "uni"}
	ActionDelete    = Action{7, "del"}
	ActionCopy      = Action{8, "cpy"}
	ActionMove      = Action{9, "mov"}
	ActionSwap      = Action{10, "swp"}
	ActionSetKey    = Action{11, "key"}
)

var textToAction = map[string]Action{
	ActionSet.Text:       ActionSet,
	ActionInsert.Text:    ActionInsert,
	ActionIncrement.Text: ActionIncrement,
	ActionPushFront.Text: ActionPushFront,
	ActionPushBack.Text:  ActionPushBack,
	ActionAddUnique.Text: ActionAddUnique,
	ActionDelete.Text:    ActionDelete,
	ActionCopy.Text:      ActionCopy,
	ActionMove.Text:      ActionMove,
	ActionSwap.Text:      ActionSwap,
	ActionSetKey.Text:    ActionSetKey,
}

func ActionFromText(text string) Action {
	if t, ok := textToAction[text]; ok {
		return t
	} else {
		return ActionUndefined
	}
}
