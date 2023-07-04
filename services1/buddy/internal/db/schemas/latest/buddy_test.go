package latest

import (
	"reflect"
	"testing"
)

const (
	reqTime = 1560176920
	ackTime = 1560176921
)

func TestNewBuddy(t *testing.T) {
	type testCase struct {
		name          string
		state         MemberType
		text          string
		time          int64
		expectedBuddy Buddy
	}
	testCases := []testCase{
		{
			name:  "testBuddy",
			state: MemberTypeBuddy,
			text:  "testBuddy - Is Buddy",
			time:  ackTime,
			expectedBuddy: Buddy{
				Name:    "testBuddy",
				State:   MemberTypeBuddy,
				Remark:  "testBuddy - Is Buddy",
				AckTime: ackTime,
				ReqTime: 0,
				ReqInfo: "",
			},
		},
		{
			name:  "testInvited",
			state: MemberTypeInviter,
			text:  "testInvited - Is Invited",
			time:  reqTime,
			expectedBuddy: Buddy{
				Name:    "testInvited",
				State:   MemberTypeInviter,
				Remark:  "",
				AckTime: 0,
				ReqTime: reqTime,
				ReqInfo: "testInvited - Is Invited",
			},
		},
		{
			name:  "♥♦♣♠",
			state: MemberTypeBuddy,
			text:  "♥♦♣♠ - Is Buddy",
			time:  ackTime,
			expectedBuddy: Buddy{
				Name:    "♥♦♣♠",
				State:   MemberTypeBuddy,
				Remark:  "♥♦♣♠ - Is Buddy",
				AckTime: ackTime,
				ReqTime: 0,
				ReqInfo: "",
			},
		},
		{
			name:  "亜美",
			state: MemberTypeBuddy,
			text:  "亜美 - Is Buddy",
			time:  ackTime,
			expectedBuddy: Buddy{
				Name:    "亜美",
				State:   MemberTypeBuddy,
				Remark:  "亜美 - Is Buddy",
				AckTime: ackTime,
				ReqTime: 0,
				ReqInfo: "",
			},
		},
		{
			name:  "invalidStateBuddy",
			state: 2,
			text:  "Invalid State Buddy",
			time:  0,
			expectedBuddy: Buddy{
				Name:    "",
				State:   0,
				Remark:  "",
				AckTime: 0,
				ReqTime: 0,
				ReqInfo: "",
			},
		},
	}
	for i, tc := range testCases {
		buddy := NewBuddy(tc.name, tc.state, tc.text, tc.time)
		if *buddy != tc.expectedBuddy {
			t.Errorf("Test #%d: Received Buddy %+v; expected Buddy %+v", i+1, *buddy, tc.expectedBuddy)
		}
	}
}

func TestBuddyQueue_AddBuddy(t *testing.T) {
	type testCase struct {
		buddy *Buddy
	}

	testCases := []testCase{
		{
			&Buddy{
				Name:    "testBuddy",
				State:   MemberTypeInviter,
				ReqTime: reqTime,
				ReqInfo: "testBuddy Info",
				AckTime: ackTime,
				Remark:  "testBuddy Remark",
			},
		},
		{
			&Buddy{
				Name:    "testInvited",
				State:   MemberTypeInviter,
				ReqTime: reqTime,
				ReqInfo: "testInvited Info",
				AckTime: ackTime,
				Remark:  "testInvited Remark",
			},
		},
		{
			&Buddy{
				Name:    "♥♦♣♠",
				State:   MemberTypeInviter,
				ReqTime: reqTime,
				ReqInfo: "♥♦♣♠ Info",
				AckTime: ackTime,
				Remark:  "♥♦♣♠ Remark",
			},
		},
		{
			&Buddy{
				Name:    "亜美",
				State:   MemberTypeBuddy,
				ReqTime: reqTime,
				ReqInfo: "亜美 Info",
				AckTime: ackTime,
				Remark:  "亜美 Remark",
			},
		},
	}

	bq := BuddyQueue{
		Uid:     "testProfile",
		Buddies: make([]*Buddy, 0),
	}

	for i, tc := range testCases {
		lenBeforeAdd := len(bq.Buddies)
		bq.AddBuddy(tc.buddy)
		if len(bq.Buddies) != lenBeforeAdd+1 {
			t.Error("Buddy was not successfully added to buddy queue in test case #", i+1)
		}
	}
}

func TestBuddyQueue_Delete(t *testing.T) {
	type testCase struct {
		buddyName string
		isInQueue bool
	}

	testCases := []testCase{
		{
			buddyName: "testBuddy",
			isInQueue: true,
		},
		{
			buddyName: "testInvited",
			isInQueue: true,
		},
		{
			buddyName: "♥♦♣♠",
			isInQueue: false,
		},
		{
			buddyName: "亜美",
			isInQueue: false,
		},
		{
			buddyName: "NotInBuddyQueue",
			isInQueue: false,
		},
	}

	for i, tc := range testCases {
		bq := BuddyQueue{
			Uid: "test",
			Buddies: []*Buddy{
				{
					Name:    "testBuddy",
					State:   MemberTypeBuddy,
					AckTime: ackTime,
					Remark:  "testBuddy Remark",
				},
				{
					Name:    "testInvited",
					State:   MemberTypeInviter,
					ReqTime: reqTime,
					ReqInfo: "testInvited Info",
				},
			},
		}

		var startingLen int
		startingLen = len(bq.Buddies)

		bq.Delete(tc.buddyName)

		var resultingLen int
		resultingLen = len(bq.Buddies)

		if tc.isInQueue {
			if resultingLen == startingLen-1 {
				for _, buddy := range bq.Buddies {
					if buddy.Name == tc.buddyName {
						t.Error("Buddy is still present after deletion in test case #", i+1)
					}
				}
			} else {
				t.Error("Buddy queue contains an unexpected number of buddies after Delete() in test case #", i+1,
					" - expected ", startingLen-1, " got ", resultingLen)
			}
		} else if resultingLen != startingLen {
			t.Error("Buddy queue contains an unexpected number of buddies after trying to Delete() a buddy"+
				" not in the queue in test case #", i+1, "- expected ", startingLen, " got ", resultingLen)
		}
	}
}

func TestBuddyQueue_FindByName(t *testing.T) {
	type testCase struct {
		providedBuddy *Buddy
		findExpected  bool
	}
	testCases := []testCase{
		{
			providedBuddy: &Buddy{
				Name:    "testName1",
				State:   MemberTypeBuddy,
				AckTime: ackTime,
				Remark:  "testName1 Remark",
			},
			findExpected: true,
		},
		{
			providedBuddy: &Buddy{
				Name:    "testName2",
				State:   MemberTypeBuddy,
				AckTime: ackTime,
				Remark:  "testName2 Remark",
			},
			findExpected: true,
		},
		{
			providedBuddy: &Buddy{
				Name:    "",
				State:   MemberTypeBuddy,
				AckTime: ackTime,
				Remark:  " Remark",
			},
			findExpected: true,
		},
		{
			providedBuddy: &Buddy{
				Name:    "testInvited",
				State:   MemberTypeInviter,
				ReqTime: reqTime,
				ReqInfo: "testInvited Info",
			},
			findExpected: true,
		},
		{
			providedBuddy: &Buddy{
				Name:    "♥♦♣♠",
				State:   MemberTypeInviter,
				ReqTime: reqTime,
				ReqInfo: "♥♦♣♠ Info",
			},
			findExpected: false,
		},
		{
			providedBuddy: &Buddy{
				Name:    "亜美",
				State:   MemberTypeBuddy,
				AckTime: ackTime,
				Remark:  "亜美 Remark",
			},
			findExpected: false,
		},

		{
			providedBuddy: &Buddy{
				Name:    "testNameNot",
				State:   -1,
				ReqTime: 0,
				ReqInfo: "notaValidInvited",
				AckTime: 0,
				Remark:  "notAValidBuddy",
			},
			findExpected: false,
		},
	}

	bq := BuddyQueue{
		Uid: "test",
		Buddies: []*Buddy{
			{
				Name:    "testName1",
				State:   MemberTypeBuddy,
				AckTime: ackTime,
				Remark:  "testName1 Remark",
			},
			{
				Name:    "testName2",
				State:   MemberTypeBuddy,
				AckTime: ackTime,
				Remark:  "testName2 Remark",
			},
			{
				Name:    "",
				State:   MemberTypeBuddy,
				AckTime: ackTime,
				Remark:  " Remark",
			},
			{
				Name:    "testInvited",
				State:   MemberTypeInviter,
				ReqTime: reqTime,
				ReqInfo: "testInvited Info",
			},
		},
	}
	for i, tc := range testCases {
		if foundBuddy, found := bq.IsContains(tc.providedBuddy.Name); !found {
			if tc.findExpected {
				t.Error("Could not find the buddy for the provided name ", tc.providedBuddy.Name,
					" in test case #", i+1)
			}
		} else if !tc.findExpected {
			t.Error("Buddy found where one was not expected to be found in test case #", i+1)
		} else {
			if foundBuddy != *tc.providedBuddy {
				t.Errorf("Test #%d: Received Buddy %+v; expected Buddy %+v", i+1, foundBuddy, tc.providedBuddy)
			}
		}
	}
}

func TestBuddyQueue_UpdateState(t *testing.T) {
	type testCase struct {
		buddy *Buddy
	}
	testCases := []testCase{
		{
			&Buddy{
				Name:    "testBuddy",
				State:   MemberTypeBuddy,
				ReqTime: reqTime,
				ReqInfo: "testBuddy Info",
				AckTime: ackTime,
			},
		},
		{
			&Buddy{
				Name:    "testInvited",
				State:   MemberTypeBuddy,
				ReqTime: reqTime,
				ReqInfo: "testInvited Info",
				AckTime: ackTime,
			},
		},
		{
			&Buddy{
				Name:    "♥♦♣♠",
				State:   MemberTypeBuddy,
				ReqTime: reqTime,
				ReqInfo: "♥♦♣♠ Info",
				AckTime: ackTime,
			},
		},
		{
			&Buddy{
				Name:    "亜美",
				State:   MemberTypeBuddy,
				ReqTime: reqTime,
				ReqInfo: "亜美 Info",
				AckTime: ackTime,
			},
		},
	}

	bq := BuddyQueue{
		Uid: "test",
		Buddies: []*Buddy{
			{
				Name:    "testBuddy",
				State:   MemberTypeInviter,
				ReqTime: reqTime,
				ReqInfo: "testBuddy Info",
				AckTime: 0,
				Remark:  "",
			},

			{
				Name:    "testInvited",
				State:   MemberTypeInviter,
				ReqTime: reqTime,
				ReqInfo: "testInvited Info",
				AckTime: 0,
				Remark:  "",
			},

			{
				Name:    "♥♦♣♠",
				State:   MemberTypeInviter,
				ReqTime: reqTime,
				ReqInfo: "♥♦♣♠ Info",
				AckTime: 0,
				Remark:  "",
			},

			{
				Name:    "亜美",
				State:   MemberTypeInviter,
				ReqTime: reqTime,
				ReqInfo: "亜美 Info",
				AckTime: 0,
				Remark:  "",
			},
		},
	}

	for i, tc := range testCases {
		foundBuddy := false
		lenBeforeUpdate := len(bq.Buddies)
		bq.UpdateState(tc.buddy.Name, tc.buddy.State, tc.buddy.AckTime)
		if len(bq.Buddies) != lenBeforeUpdate {
			t.Errorf("Buddy queue length changed where it was not expected to change in test case %v; expected %v,"+
				"got %v", i+1, len(testCases), len(bq.Buddies))
		} else {
			for _, buddy := range bq.Buddies {
				if buddy.Name == tc.buddy.Name {
					foundBuddy = true
					if buddy.State != tc.buddy.State {
						t.Error("Buddy State updated when it was not expected to change in test case #", i+1,
							"expected value is ", tc.buddy.State, " actual value is ", buddy.State)
					}
					if buddy.ReqInfo != tc.buddy.ReqInfo {
						t.Error("Buddy ReqInfo updated when it was not expected to change in test case #", i+1,
							"expected value is ", tc.buddy.ReqInfo, " actual value is ", buddy.ReqInfo)
					}
					if buddy.ReqTime != tc.buddy.ReqTime {
						t.Error("Buddy ReqTime updated when it was not expected to change in test case #", i+1,
							"expected value is ", tc.buddy.ReqTime, " actual value is ", buddy.ReqTime)
					}
					if buddy.Remark != tc.buddy.Remark {
						t.Error("Buddy State updated when it was not expected to change in test case #", i+1,
							"expected value is ", tc.buddy.Remark, " actual value is ", buddy.Remark)
					}
					if buddy.AckTime != tc.buddy.AckTime {
						t.Error("Buddy AckTime updated when it was not expected to change in test case #", i+1,
							"expected value is ", tc.buddy.AckTime, " actual value is ", buddy.AckTime)
					}
				}
			}
			if !foundBuddy {
				t.Error("Buddy was not found in test case #", i+1)
			}
		}
	}
}

func TestBuddyQueue_UpdateRemark(t *testing.T) {
	type testCase struct {
		buddy *Buddy
	}
	testCases := []testCase{
		{
			&Buddy{
				Name:    "testBuddy",
				State:   MemberTypeBuddy,
				AckTime: ackTime,
				Remark:  "testBuddy Remark",
			},
		},
	}
	bq := BuddyQueue{
		Uid: "test",
		Buddies: []*Buddy{
			{
				Name:    "testBuddy",
				State:   MemberTypeBuddy,
				ReqTime: 0,
				ReqInfo: "",
				AckTime: ackTime,
				Remark:  "",
			},
		},
	}

	for i, tc := range testCases {
		foundBuddy := false
		lenBeforeUpdate := len(bq.Buddies)
		bq.UpdateRemark(tc.buddy.Name, tc.buddy.Remark)
		if len(bq.Buddies) != lenBeforeUpdate {
			t.Errorf("Buddy queue length changed where it was not expected to change in test case %v; expected %v,"+
				"got %v", i+1, len(testCases), len(bq.Buddies))
		} else {
			for _, buddy := range bq.Buddies {
				if buddy.Name == tc.buddy.Name {
					foundBuddy = true
					if buddy.State != tc.buddy.State {
						t.Error("Buddy State updated when it was not expected to change in test case #", i+1,
							"expected value is ", tc.buddy.State, " actual value is ", buddy.State)
					}
					if buddy.ReqInfo != tc.buddy.ReqInfo {
						t.Error("Buddy ReqInfo updated when it was not expected to change in test case #", i+1,
							"expected value is ", tc.buddy.ReqInfo, " actual value is ", buddy.ReqInfo)
					}
					if buddy.ReqTime != tc.buddy.ReqTime {
						t.Error("Buddy ReqTime updated when it was not expected to change in test case #", i+1,
							"expected value is ", tc.buddy.ReqTime, " actual value is ", buddy.ReqTime)
					}
					if buddy.Remark != tc.buddy.Remark {
						t.Error("Buddy State updated when it was not expected to change in test case #", i+1,
							"expected value is ", tc.buddy.Remark, " actual value is ", buddy.Remark)
					}
					if buddy.AckTime != tc.buddy.AckTime {
						t.Error("Buddy AckTime updated when it was not expected to change in test case #", i+1,
							"expected value is ", tc.buddy.AckTime, " actual value is ", buddy.AckTime)
					}
				}
			}
			if !foundBuddy {
				t.Error("Buddy was not found in test case #", i+1)
			}
		}
	}
}

func TestBuddyQueue_FilterBuddiesByIds(t *testing.T) {
	tests := []struct {
		name   string
		fields []*Buddy
		args   []string
		want   []string
	}{
		{
			"1",
			[]*Buddy{
				{Name: "1"},
				{Name: "2"},
			},
			[]string{"1", "2", "3"},
			[]string{"1", "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bq := &BuddyQueue{
				Buddies: tt.fields,
			}
			if got := bq.GetBuddiesByIds(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterBuddiesByIds() = %v, want %v", got, tt.want)
			}
		})
	}
}
