package latest

import (
	"fmt"
	"reflect"
	"testing"
)

var (
	TestPlayers0 = []*RecentProfile{{}}
	TestPlayers1 = []*RecentProfile{{ID: "1"}}
	TestPlayers2 = []*RecentProfile{{ID: "1"}, {ID: "1"}, {ID: "1"}, {ID: "2"}, {ID: "2"}}
	TestPlayers3 = []*RecentProfile{{ID: "1"}, {ID: "1"}, {ID: "2"}, {ID: "2"}, {ID: "3"}}
	TestPlayers4 = []*RecentProfile{{ID: "1"}, {ID: "2"}, {ID: "2"}, {ID: "3"}, {ID: "3"}}
	TestPlayers5 = []*RecentProfile{{ID: "1"}, {ID: "2"}, {ID: "1"}, {ID: "3"}, {ID: "2"}}
	TestPlayers6 = []*RecentProfile{{ID: "1"}, {ID: "2"}, {ID: "3"}}
)
var (
	AddPlayers0 = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14"}
	AddPlayers1 = []string{"3", "2", "3"}
)

func TestRecentProfiles_removeDuplicateInOrder(t *testing.T) {
	type fields struct {
		ID      string
		Players []*RecentProfile
	}
	type args struct {
		in []*RecentProfile
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*RecentProfile
	}{
		{
			name: "0",
			args: args{in: TestPlayers0},
			want: []*RecentProfile{{}},
		},
		{
			name: "1",
			args: args{in: TestPlayers1},
			want: []*RecentProfile{{ID: "1"}},
		},
		{
			name: "2",
			args: args{in: TestPlayers2},
			want: []*RecentProfile{{ID: "1"}, {ID: "2"}},
		},
		{
			name: "3",
			args: args{in: TestPlayers3},
			want: []*RecentProfile{{ID: "1"}, {ID: "2"}, {ID: "3"}},
		},
		{
			name: "4",
			args: args{in: TestPlayers4},
			want: []*RecentProfile{{ID: "1"}, {ID: "2"}, {ID: "3"}},
		},
		{
			name: "5",
			args: args{in: TestPlayers5},
			want: []*RecentProfile{{ID: "1"}, {ID: "2"}, {ID: "3"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := &BuddyQueue{
				Uid:       tt.fields.ID,
				RecentMet: tt.fields.Players,
			}
			if got := rp.removeDuplicatesInOrder(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeDuplicatesInOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecentProfiles_Add(t *testing.T) {
	tests := []struct {
		name   string
		fields BuddyQueue
		args   []string
	}{
		{
			name: "0",
			fields: BuddyQueue{
				Uid:       "0",
				RecentMet: TestPlayers6,
			},
			args: AddPlayers0,
		},
		{
			name: "1",
			fields: BuddyQueue{
				Uid:       "1",
				RecentMet: TestPlayers6,
			},
			args: AddPlayers1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := &BuddyQueue{
				Uid:       tt.fields.Uid,
				RecentMet: tt.fields.RecentMet,
			}
			lastStr := makeStr(rp.RecentMet)
			rp.AddRecentProfiles(tt.args...)
			resultStr := makeStr(rp.RecentMet)
			t.Logf("recent players:%v + players:%v = result:%v",
				lastStr, tt.args, resultStr)
		})
	}
}

func makeStr(in []*RecentProfile) string {
	length := len(in)
	out := "["
	for k, v := range in {
		if k >= length-1 {
			out += fmt.Sprintf("%v", v.ID)
		} else {
			out += fmt.Sprintf("%v ", v.ID)
		}

	}
	return out + "]"
}
