package model

import (
	pb "github.com/donglei1234/platform/services/condition/gen/condition/api"
	"strconv"
)

// Condition is a condition.
type Condition struct {
	Data map[string]*pb.Condition `json:"data" `
}

// NewCondition creates a new condition.
func NewCondition() *Condition {
	return &Condition{
		Data: make(map[string]*pb.Condition),
	}
}

// IsEmpty returns true if the condition is empty.
func (c *Condition) IsEmpty() bool {
	return len(c.Data) == 0
}

// GetConditions returns all condition.
func (c *Condition) GetConditions() map[string]*pb.Condition {
	return c.Data
}

// GetConditionLst returns a list condition.
func (c *Condition) GetConditionLst() []*pb.Condition {
	var conds []*pb.Condition
	for _, v := range c.Data {
		conds = append(conds, v)
	}
	return conds
}

// SetCondition sets a condition.
func (c *Condition) SetCondition(data map[string]*pb.Condition) {
	c.Data = data
}

// AddCondition adds a condition.
func (c *Condition) AddCondition(conds ...*pb.Condition) {
	for _, cond := range conds {
		id := c.makeConditionId(cond.Theme, cond.OwnerId)
		c.checkAndUpdateCondition(cond)
		c.Data[id] = cond

	}
}

// RemoveCondition removes a condition.
func (c *Condition) RemoveCondition(conds ...*pb.Condition) {
	for _, cond := range conds {
		id := c.makeConditionId(cond.Theme, cond.OwnerId)
		delete(c.Data, id)
	}
}

// UpdateCondition replace
func (c *Condition) UpdateCondition(update *pb.Condition) {
	for _, cond := range c.Data {
		if cond.Params[0] == update.Params[0] &&
			cond.Params[1] == update.Params[1] {
			if cond.UpdateStrategy == pb.Condition_STRATEGY_ADD {
				cond.Progress += update.Progress
			} else {
				cond.Progress = update.Progress
			}
			c.checkAndUpdateCondition(cond)
		}
	}
}

// IncrementCondition increments a condition.
func (c *Condition) IncrementCondition(update *pb.Condition) {
	for _, cond := range c.Data {
		if cond.Params[0] == update.Params[0] &&
			cond.Params[1] == update.Params[1] {
			cond.Progress += update.Progress
			c.checkAndUpdateCondition(cond)
		}
	}
}

// checkAndUpdateCondition checks a condition is finish.
func (c *Condition) checkAndUpdateCondition(cond *pb.Condition) {
	if cond.Progress >= cond.Params[2] {
		cond.Progress = cond.Params[2]
		cond.Status = pb.Condition_FINISHED
	}
}

// RemoveFinishedCondition removes finished condition.
func (c *Condition) RemoveFinishedCondition() {
	for key, cond := range c.Data {
		if cond.Status == pb.Condition_FINISHED {
			delete(c.Data, key)
		}
	}
}

// makeConditionId makes a condition id.
func (c *Condition) makeConditionId(theme string, id int32) string {
	idStr := strconv.Itoa(int(id))
	return theme + "_" + idStr
}

func (c *Condition) DeleteCondition(v *pb.Condition) {

}
