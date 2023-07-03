package latest

import "testing"

func TestBuddyQueue_DeleteBlockedProfiles(t *testing.T) {
	tests := []struct {
		name    string
		fields  []*BlockedProfile
		args    []string
		results []*BlockedProfile
	}{
		{
			"1",
			[]*BlockedProfile{
				{
					ID:      "1",
					AddTime: 0,
				},
				{
					ID:      "2",
					AddTime: 0,
				},
				{
					ID:      "3",
					AddTime: 0,
				},
				{
					ID:      "4",
					AddTime: 0,
				},
			},
			[]string{"1", "3"},
			[]*BlockedProfile{
				{
					ID:      "4",
					AddTime: 0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bq := &BuddyQueue{
				BlockedProfiles: tt.fields,
			}
			bq.DeleteBlockedProfiles(tt.args...)
			t.Logf("blocked list remove %v ", tt.args)
			for _, v := range bq.BlockedProfiles {
				t.Log("result: ID:", v.ID)
			}
		})
	}
}
