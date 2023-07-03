package latest

import "time"

type RecentProfile struct {
	ID      string
	AddTime int64
}

func (bq *BuddyQueue) AddRecentProfiles(uid ...string) {
	if len(uid) <= 0 {
		return
	}
	profiles := make([]*RecentProfile, len(uid))
	for k, v := range uid {
		// player self will be filtered
		if v == bq.Uid {
			continue
		}
		profiles[k] = &RecentProfile{
			ID:      v,
			AddTime: time.Now().Unix(),
		}
	}
	profiles = bq.removeDuplicatesInOrder(profiles)
	profilesLen := len(profiles)
	if profilesLen >= RecentMetCountMax {
		bq.RecentMet = profiles[profilesLen-RecentMetCountMax:]
	} else {
		bq.RecentMet = bq.removeDuplicatesInOrder(append(profiles, bq.RecentMet...))
		if len(bq.RecentMet) > RecentMetCountMax {
			bq.RecentMet = bq.RecentMet[:RecentMetCountMax]
		}
	}
}

func (bq *BuddyQueue) DeleteRecentProfiles(ids ...string) {
	for _, v := range ids {
		for i := 0; i < len(bq.RecentMet); i++ {
			if bq.RecentMet[i].ID == v {
				bq.RecentMet = append(bq.RecentMet[:i], bq.RecentMet[i+1:]...)
				i--
			}
		}
	}
}

func (bq *BuddyQueue) removeDuplicatesInOrder(profiles []*RecentProfile) []*RecentProfile {
	if len(profiles) <= 0 {
		return nil
	}
	checks := make(map[string]bool)
	result := make([]*RecentProfile, 0)
	for _, v := range profiles {
		if ok := checks[v.ID]; !ok {
			checks[v.ID] = true
			result = append(result, v)
		}
	}
	return result
}
