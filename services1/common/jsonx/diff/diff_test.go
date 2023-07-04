package diff

import (
	"testing"

	"github.com/donglei1234/platform/services/common/dupblock"
)

func TestDiff(t *testing.T) {
	var testCases = []struct {
		left         string
		right        string
		includePaths []string
	}{
		// test case 1 -- change root
		{left: `"test"`, right: `{"value":true}`, includePaths: []string{}},

		// test case 2 -- add
		{left: `{}`, right: `{"value":true}`, includePaths: []string{}},

		// test case 3 -- remove
		{left: `{"value":true}`, right: `{}`, includePaths: []string{}},

		// test case 4 -- replace
		{left: `{"value":true}`, right: `{"value":false}`, includePaths: []string{}},

		// test case 5 -- move
		{left: `{"obj1":{"id":"obj1"}}`, right: `{"obj2":{"id":"obj1"}}`, includePaths: []string{}},

		// test case 6 -- swap
		{left: `{"obj1":{"id1":"obj1"}, "obj2":{"id2":"obj2"}}`, right: `{"obj1":{"id2":"obj2"}, "obj2":{"id1":"obj1"}}`, includePaths: []string{}},

		// test case 7 -- swap array
		{left: `[{"id":1}, {"id":2}]`, right: `[{"id":2}, {"id":1}, {"id":2}]`, includePaths: []string{}},

		// test case 8 -- copy
		{left: `{"obj1":{"id":"obj1"}}`, right: `{"obj1":{"id":"obj1"}, "obj2":{"id":"obj1"}, "obj3":{"id":"obj1"}}`, includePaths: []string{}},

		// test case 9 -- insert array
		{left: `[{"id":1}, {"id":2}, {"id":1}]`, right: `[{"id":2}, {"id":1}, {"id":2}]`, includePaths: []string{}},

		// test case 10 -- filter paths
		{left: `{"obj1":{"id1":"obj1"}, "obj2":{"id2":"obj2"}}`, right: `{"obj1":{"id1":"new obj1"}, "obj2":{"id2":"new obj2"}}`, includePaths: []string{"obj2.id2"}},
	}
	for i, tc := range testCases {
		if leftDocument, err := NewDocument([]byte(tc.left)); err != nil {
			t.Error("Error encountered while creating a new document in TestDiff() test case #:", i+1, ":", err)
		} else {
			if rightDocument, err := NewDocument([]byte(tc.right)); err != nil {
				t.Error("Error encountered while creating a new document in TestDiff() test case #:", i+1, ":", err)
			} else {
				NewDiff(&leftDocument, &rightDocument, WithIncludePaths(tc.includePaths))
			}
		}
	}
}

func TestWriteDUPBlock(t *testing.T) {
	var testCases = []struct {
		left         string
		right        string
		includePaths []string
	}{
		// test case 1 -- change root
		{left: `"test"`, right: `{"value":true}`, includePaths: []string{}},

		// test case 2 -- add
		{left: `{}`, right: `{"value":true}`, includePaths: []string{}},

		// test case 3 -- remove
		{left: `{"value":true}`, right: `{}`, includePaths: []string{}},

		// test case 4 -- replace
		{left: `{"value":true}`, right: `{"value":false}`, includePaths: []string{}},

		// test case 5 -- move
		{left: `{"obj1":{"id":"obj1"}}`, right: `{"obj2":{"id":"obj1"}}`, includePaths: []string{}},

		// test case 6 -- swap
		{left: `{"obj1":{"id1":"obj1"}, "obj2":{"id2":"obj2"}}`, right: `{"obj1":{"id2":"obj2"}, "obj2":{"id1":"obj1"}}`, includePaths: []string{}},

		// test case 7 -- swap array
		{left: `[{"id":1}, {"id":2}]`, right: `[{"id":2}, {"id":1}, {"id":2}]`, includePaths: []string{}},

		// test case 8 -- copy
		{left: `{"obj1":{"id":"obj1"}}`, right: `{"obj1":{"id":"obj1"}, "obj2":{"id":"obj1"}, "obj3":{"id":"obj1"}}`, includePaths: []string{}},

		// test case 9 -- insert array
		{left: `[{"id":1}, {"id":2}, {"id":1}]`, right: `[{"id":2}, {"id":1}, {"id":2}]`, includePaths: []string{}},

		// test case 10 -- filter paths
		{left: `{"obj1":{"id1":"obj1"}, "obj2":{"id2":"obj2"}}`, right: `{"obj1":{"id1":"new obj1"}, "obj2":{"id2":"new obj2"}}`, includePaths: []string{"obj2.id2"}},
	}
	for i, tc := range testCases {
		if leftDocument, err := NewDocument([]byte(tc.left)); err != nil {
			t.Error("Error while creating a new document in TestDiff() test case #:", i+1, ":", err)
		} else if rightDocument, err := NewDocument([]byte(tc.right)); err != nil {
			t.Error("Error while creating a new document in TestDiff() test case #:", i+1, ":", err)
		} else if writer, err := dupblock.NewTextWriter(); err != nil {
			t.Error("Error creating dupblock.TextWriter in TestDiff() test case #:", i+1, ":", err)
		} else {
			testDiff := NewDiff(&leftDocument, &rightDocument, WithIncludePaths(tc.includePaths))
			testDiff.WriteDUPBlock(writer)
		}
	}
}

func TestChildInference(t *testing.T) {
	// TODO: @SNICHOLS: this is a quick and dirty test... really should refactor to something more general

	left := `{
	"ProfileName": "khartsook",
	"LastClaimedLoginBonusTime": 0,
	"Notifications": [
		{
			"ID": "MbWMcb7KSeubbKod9ibdwg",
			"Archetype": "arch_notify_quest_ready_to_complete",
			"SendTime": 1581536311,
			"ReadTime": 0,
			"Metadata": {
				"Metadata": {
					"AcceptanceTime": "1570810863",
					"QuestID": "arch_quest_test_kill5drones"
				}
			}
		},
		{
			"ID": "HUfaX60YQRCy14IIBNLg5Q",
			"Archetype": "arch_notify_quest_ready_to_complete",
			"SendTime": 1581536399,
			"ReadTime": 0,
			"Metadata": {
				"Metadata": {
					"AcceptanceTime": "1579194298",
					"QuestID": "arch_quest_test_acquireitems003"
				}
			}
		},
		{
			"ID": "ta7z6sZ1Sk2pyRJtBGNyjg",
			"Archetype": "arch_notify_quest_ready_to_complete",
			"SendTime": 1581536415,
			"ReadTime": 0,
			"Metadata": {
				"Metadata": {
					"AcceptanceTime": "1570810863",
					"QuestID": "arch_quest_test_useweapon001"
				}
			}
		},
		{
			"ID": "6OUKXYxdQgm/E/2HiX5u4w",
			"Archetype": "arch_notify_quest_ready_to_complete",
			"SendTime": 1581536439,
			"ReadTime": 0,
			"Metadata": {
				"Metadata": {
					"AcceptanceTime": "1576866133",
					"QuestID": "arch_quest_unlock_zone1sector2"
				}
			}
		},
		{
			"ID": "uc+JuKW0SqGHV2HatlHOmg",
			"Archetype": "arch_notify_quest_ready_to_complete",
			"SendTime": 1581536447,
			"ReadTime": 0,
			"Metadata": {
				"Metadata": {
					"AcceptanceTime": "1576866133",
					"QuestID": "arch_quest_unlock_zone1sector3"
				}
			}
		},
		{
			"ID": "fC32oeoGTlOZvWc1pxy35w",
			"Archetype": "arch_notify_mission_completed",
			"SendTime": 1581536482,
			"ReadTime": 0,
			"Metadata": {
				"Metadata": {
					"CompletionCount": "2",
					"QuestID": "arch_mission_test_qa_dangerroom"
				}
			}
		},
		{
			"ID": "y9FGdPDWQnCOAw9W3CaWVg",
			"Archetype": "arch_notify_level_up",
			"SendTime": 1581536501,
			"ReadTime": 0,
			"Metadata": {
				"Metadata": {
					"CharacterArchetypeId": "arch_invalid",
					"CharacterInstanceId": "arch_invalid",
					"WeaponArchetypeId0": "arch_weap_role_energy_cannon_t1",
					"WeaponArchetypeId1": "arch_weap_sword",
					"WeaponArchetypeId2": "arch_invalid",
					"WeaponInstanceId": "P7kXqM8xRBGFl4L0dAH7Ew",
					"WeaponInstanceId2": "arch_invalid",
					"WeaponNewLevel0": "4",
					"WeaponNewLevel1": "4",
					"WeaponNewXp0": "500",
					"WeaponNewXp1": "500",
					"WeaponOldLevel0": "0",
					"WeaponOldLevel1": "0",
					"WeaponOldXp0": "0",
					"WeaponOldXp1": "0"
				}
			}
		},
		{
			"ID": "NOF4Z9lfTX6t01JtxaU1Aw",
			"Archetype": "arch_notify_quest_accepted",
			"SendTime": 1581536543,
			"ReadTime": 0,
			"Metadata": {
				"Metadata": {
					"AcceptanceTime": "1576685516",
					"QuestID": "arch_quest_test_reach_level_6"
				}
			}
		},
		{
			"ID": "IP+VIhuvSLee0AxEvMNn4g",
			"Archetype": "arch_notify_quest_accepted",
			"SendTime": 1581536550,
			"ReadTime": 0,
			"Metadata": {
				"Metadata": {
					"AcceptanceTime": "1576685516",
					"QuestID": "arch_quest_test_own50knano"
				}
			}
		},
		{
			"ID": "SzC3xNpCSsuzWNt+0m3X3g",
			"Archetype": "arch_notify_quest_accepted",
			"SendTime": 1581536561,
			"ReadTime": 0,
			"Metadata": {
				"Metadata": {
					"AcceptanceTime": "1576685516",
					"QuestID": "arch_quest_kill_WithElemDmg01"
				}
			}
		}
	]
}`

	right := `{
	"ProfileName": "khartsook",
	"LastClaimedLoginBonusTime": 0,
	"Notifications": [
		{
			"ID": "MbWMcb7KSeubbKod9ibdwg",
			"Archetype": "arch_notify_quest_ready_to_complete",
			"SendTime": 1581536311,
			"ReadTime": 0,
			"Metadata": {
				"Metadata": {
					"AcceptanceTime": "1570810863",
					"QuestID": "arch_quest_test_kill5drones"
				}
			}
		},
		{
			"ID": "HUfaX60YQRCy14IIBNLg5Q",
			"Archetype": "arch_notify_quest_ready_to_complete",
			"SendTime": 1581536399,
			"ReadTime": 0,
			"Metadata": {
				"Metadata": {
					"AcceptanceTime": "1579194298",
					"QuestID": "arch_quest_test_acquireitems003"
				}
			}
		},
		{
			"ID": "ta7z6sZ1Sk2pyRJtBGNyjg",
			"Archetype": "arch_notify_quest_ready_to_complete",
			"SendTime": 1581536415,
			"ReadTime": 0,
			"Metadata": {
				"Metadata": {
					"AcceptanceTime": "1570810863",
					"QuestID": "arch_quest_test_useweapon001"
				}
			}
		},
		{
			"ID": "6OUKXYxdQgm/E/2HiX5u4w",
			"Archetype": "arch_notify_quest_ready_to_complete",
			"SendTime": 1581536439,
			"ReadTime": 0,
			"Metadata": {
				"Metadata": {
					"AcceptanceTime": "1576866133",
					"QuestID": "arch_quest_unlock_zone1sector2"
				}
			}
		},
		{
			"ID": "uc+JuKW0SqGHV2HatlHOmg",
			"Archetype": "arch_notify_quest_ready_to_complete",
			"SendTime": 1581536447,
			"ReadTime": 0,
			"Metadata": {
				"Metadata": {
					"AcceptanceTime": "1576866133",
					"QuestID": "arch_quest_unlock_zone1sector3"
				}
			}
		},
		{
			"ID": "y9FGdPDWQnCOAw9W3CaWVg",
			"Archetype": "arch_notify_level_up",
			"SendTime": 1581536501,
			"ReadTime": 0,
			"Metadata": {
				"Metadata": {
					"CharacterArchetypeId": "arch_invalid",
					"CharacterInstanceId": "arch_invalid",
					"WeaponArchetypeId0": "arch_weap_role_energy_cannon_t1",
					"WeaponArchetypeId1": "arch_weap_sword",
					"WeaponArchetypeId2": "arch_invalid",
					"WeaponInstanceId": "P7kXqM8xRBGFl4L0dAH7Ew",
					"WeaponInstanceId2": "arch_invalid",
					"WeaponNewLevel0": "4",
					"WeaponNewLevel1": "4",
					"WeaponNewXp0": "500",
					"WeaponNewXp1": "500",
					"WeaponOldLevel0": "0",
					"WeaponOldLevel1": "0",
					"WeaponOldXp0": "0",
					"WeaponOldXp1": "0"
				}
			}
		}
	]
}`

	expected := `diffs: 24
0000: set Notifications[5].Metadata.Metadata.CharacterArchetypeId to "arch_invalid"
0001: set Notifications[5].Metadata.Metadata.CharacterInstanceId to "arch_invalid"
0002: set Notifications[5].Metadata.Metadata.WeaponArchetypeId0 to "arch_weap_role_energy_cannon_t1"
0003: set Notifications[5].Metadata.Metadata.WeaponArchetypeId1 to "arch_weap_sword"
0004: set Notifications[5].Metadata.Metadata.WeaponArchetypeId2 to "arch_invalid"
0005: set Notifications[5].Metadata.Metadata.WeaponInstanceId to "P7kXqM8xRBGFl4L0dAH7Ew"
0006: set Notifications[5].Metadata.Metadata.WeaponInstanceId2 to "arch_invalid"
0007: set Notifications[5].Metadata.Metadata.WeaponNewLevel0 to "4"
0008: set Notifications[5].Metadata.Metadata.WeaponNewLevel1 to "4"
0009: set Notifications[5].Metadata.Metadata.WeaponNewXp0 to "500"
000a: set Notifications[5].Metadata.Metadata.WeaponNewXp1 to "500"
000b: set Notifications[5].Metadata.Metadata.WeaponOldLevel0 to "0"
000c: set Notifications[5].Metadata.Metadata.WeaponOldLevel1 to "0"
000d: set Notifications[5].Metadata.Metadata.WeaponOldXp0 to "0"
000e: set Notifications[5].Metadata.Metadata.WeaponOldXp1 to "0"
000f: set Notifications[5].ID to "y9FGdPDWQnCOAw9W3CaWVg"
0010: set Notifications[5].Archetype to "arch_notify_level_up"
0011: set Notifications[5].SendTime to 1581536501
0012: remove Notifications[9] 0 0
0013: remove Notifications[8] 0 0
0014: remove Notifications[7] 0 0
0015: remove Notifications[6] 0 0
0016: remove Notifications[5].Metadata.Metadata.QuestID 0 0
0017: remove Notifications[5].Metadata.Metadata.CompletionCount 0 0
`

	if ld, err := NewDocument([]byte(left)); err != nil {
		t.Fatal(err)
	} else if rd, err := NewDocument([]byte(right)); err != nil {
		t.Fatal(err)
	} else {
		d := NewDiff(&ld, &rd)
		if d.DebugString() != expected {
			t.Fatal("diff not expected")
		}
	}
}
