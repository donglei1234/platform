package v0

const (
	// Already a buddy.
	IsBuddy int32 = 0
	// Need to be identified by self.
	IsInvited int32 = 1
)

type Buddy struct {
	// The Name for a specific buddy.
	Name string
	// The State for buddy.
	State int32
	// Timestamp (UTC) when request add buddy was send.
	ReqTime int64
	// Request add buddy text.
	ReqInfo string
	// Timestamp (UTC) when  add buddies success.
	AckTime int64
	// The Remark for buddy.
	Remark string
}

func NewBuddy(name string, state int32, text string, time int64) *Buddy {
	if state == IsBuddy {
		return &Buddy{
			Name:    name,
			State:   state,
			Remark:  text,
			AckTime: time,
		}
	} else if state == IsInvited {
		return &Buddy{
			Name:    name,
			State:   state,
			ReqTime: time,
			ReqInfo: text,
		}
	}
	return &Buddy{}
}

type BuddyQueue struct {
	SchemaVersion int
	FixupVersion  int
	ProfileName   string
	Buddies       []*Buddy
}

// Add a new buddy instance to the queue.
func (bq *BuddyQueue) AddBuddy(f *Buddy) {
	bq.Buddies = append(bq.Buddies, f)
}

func (bq *BuddyQueue) Delete(name string) {
	for i, n := range bq.Buddies {
		if n.Name == name {
			copy(bq.Buddies[i:], bq.Buddies[i+1:])
			bq.Buddies[len(bq.Buddies)-1] = nil
			bq.Buddies = bq.Buddies[:len(bq.Buddies)-1]
			break
		}
	}

}

func (bq *BuddyQueue) FindByName(name string) (Buddy, bool) {
	for _, n := range bq.Buddies {
		if n.Name == name {
			return *n, true
		}
	}

	return Buddy{}, false
}

// Update buddy info, does not contain remark.
func (bq *BuddyQueue) UpdateState(name string, state int32, ackTime int64) {
	for _, n := range bq.Buddies {
		if n.Name == name {
			n.State = state
			n.AckTime = ackTime
			return
		}
	}
}

func (bq *BuddyQueue) UpdateRemark(name string, remark string) {
	for _, n := range bq.Buddies {
		if n.Name == name {
			n.Remark = remark
			return
		}
	}
}
