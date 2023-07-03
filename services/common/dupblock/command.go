package dupblock

type Command struct {
	Action Action
	To     string
	From   string
	Delta  int
	Value  []byte
}
