package latest

import (
	pb "github.com/donglei1234/platform/services/proto/gen/buddy/api"
	"time"
)

type BlockedProfile struct {
	ID      string
	AddTime int64
}

func (p *BlockedProfile) ToProto() *pb.Blocked {
	return &pb.Blocked{
		Uid:     p.ID,
		AddTime: p.AddTime,
	}
}

func (bq *BuddyQueue) AddBlockedProfiles(ids ...string) {
	for _, v := range ids {
		if v == "" {
			continue
		}
		bq.BlockedProfiles[v] = &BlockedProfile{
			ID:      v,
			AddTime: time.Now().Unix(),
		}
	}
}

func (bq *BuddyQueue) DeleteBlockedProfiles(ids ...string) {
	if len(ids) <= 0 {
		bq.BlockedProfiles = make(map[string]*BlockedProfile)
		return
	}
	for _, v := range ids {
		if v == "" {
			continue
		}
		delete(bq.BlockedProfiles, v)
	}
}

func (bq *BuddyQueue) IsBlocked(id string) bool {
	if _, ok := bq.BlockedProfiles[id]; ok {
		return true
	}
	return false
}

func (bq *BuddyQueue) GetBlockedNum() int32 {
	return int32(len(bq.BlockedProfiles))
}

func (bq *BuddyQueue) FilterBlocked(ids ...string) []string {
	if len(bq.BlockedProfiles) <= 0 || bq.BlockedProfiles == nil {
		return ids
	}
	var blockedProfiles []string

	for _, v := range ids {
		if bq.IsBlocked(v) {
			continue
		}
		blockedProfiles = append(blockedProfiles, v)
	}
	return blockedProfiles
}
