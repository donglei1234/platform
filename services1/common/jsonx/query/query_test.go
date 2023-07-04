package query

import (
	"bytes"
	"strings"
	"testing"

	"github.com/donglei1234/platform/services/common/jsonx"
)

var bigFatJson = []byte(`
	{
		"Version":1,
		"Name":"Default",
		"SelectedLoadout":"F1CXpGdISIaffb2tzJuXNw",
		"Characters":[
			{"ID":"Sw0i5IIlQfiDHJlD3D5/6w","Type":"arch_char_snakebite","XP":{"Earned":0,"Level":0},"Stats":null},
			{"ID":"UVF+DHRmSR6WJ7r720zeHQ","Type":"arch_char_hotrod","XP":{"Earned":0,"Level":0},"Stats":null},
			{"ID":"KtY4l4/lS6aS0TY2i1pnVw","Type":"arch_char_ironhide","XP":{"Earned":0,"Level":0},"Stats":null},
			{"ID":"4N+eLrYDSwOd3UF+4me/dQ","Type":"arch_char_ratchet","XP":{"Earned":0,"Level":0},"Stats":null},
			{"ID":"aCH5ErdTT3WpO2slry3C/A","Type":"arch_char_bumblebee","XP":{"Earned":0,"Level":0},"Stats":null},
			{"ID":"yZr+YD0CStG4lziRLw/Yhw","Type":"arch_char_shockwave","XP":{"Earned":0,"Level":0},"Stats":null},
			{"ID":"xnEqPhkgSty+5jcssTeafA","Type":"arch_char_flying","XP":{"Earned":0,"Level":0},"Stats":null},
			{"ID":"nurV5dyVRTSczi+ldS1qMw","Type":"arch_char_soundwave","XP":{"Earned":0,"Level":0},"Stats":null}
		],
		"Weapons":[
			{"ID":"ynoLtQ3sSWSr/LMX2sRtwg","Type":"arch_weap_shock_carbine","XP":{"Earned":0,"Level":0}},
			{"ID":"AQotA74EQ1Gfa1ksglKosg","Type":"arch_weap_sword","XP":{"Earned":0,"Level":0}},
			{"ID":"tQKsPTOKTYut/Kx/rEvK1Q","Type":"arch_weap_enervating_charge_blaster","XP":{"Earned":0,"Level":0}},
			{"ID":"G5CgC0U6QXKDkS7xTMTfNQ","Type":"arch_weap_mag_grenade_launcher","XP":{"Earned":0,"Level":0}},
			{"ID":"jBTb/knSRpWh0dubaAuKnw","Type":"arch_weap_energy_cannon","XP":{"Earned":0,"Level":0}},
			{"ID":"6VT8C97ERHakeeKz97tVCA","Type":"arch_weap_gravitic_maul","XP":{"Earned":0,"Level":0}},
			{"ID":"tEOWLnjSTrGGWsxWLyYqEA","Type":"arch_weap_energy_gauntlets","XP":{"Earned":0,"Level":0}},
			{"ID":"yGWxtSA3Smm2w5Hia8wAvQ","Type":"arch_weap_twin_thermal_swords","XP":{"Earned":0,"Level":0}},
			{"ID":"RktXYR3/SGyXYTmu1gt8Hg","Type":"arch_weap_double_barrel_shotgun","XP":{"Earned":0,"Level":0}},
			{"ID":"LpSjS2yzREetJsBumleUZw","Type":"arch_weap_plasma_thrower","XP":{"Earned":0,"Level":0}}
		],
		"Chips":[
			{"ID":"PUcjYjEDSteFw5R9gd4ISw","Type":"arch_chip_precisionmultiplier","Level":1},
			{"ID":"ZXXKxQItQVu7OiI83Aj5/A","Type":"arch_chip_damagethermal","Level":1},
			{"ID":"FowZ+L27RqCOpXssmDSfJA","Type":"arch_chip_shieldrechargerate","Level":1},
			{"ID":"qW/qq3BhSGyXZ8NmV9sKvA","Type":"arch_chip_damagenervating","Level":1},
			{"ID":"QLYCCuADTlipKjUG+tWbGQ","Type":"arch_chip_abilitycritchance","Level":1},
			{"ID":"krjVeMjkTiWYRlQjadqCqQ","Type":"arch_chip_critchanceondamage","Level":1},
			{"ID":"X5FyK53cSbWSZNBCmRThrw","Type":"arch_chip_clipsize","Level":1},
			{"ID":"PyXZUyFZS0mxa133hBXMWQ","Type":"arch_chip_meleeattackspeed","Level":1},
			{"ID":"chcLk88DTCKXC1n76LDIgQ","Type":"arch_chip_vehiclemaxspeed","Level":1},
			{"ID":"B2QUQl3qQCSKcUij2TRBGQ","Type":"arch_chip_abilityrechargeondamage","Level":1},
			{"ID":"eh2Nox0LRJ6R5LiHsLpryA","Type":"arch_chip_damagemagnetic","Level":1},
			{"ID":"gs3V4D+bQraFkp4iApQGDg","Type":"arch_chip_firerate","Level":1},
			{"ID":"5zalE3QvRWuXVYz2yUz7eA","Type":"arch_chip_healthondamage","Level":1},
			{"ID":"RS4Y8jqZR5yqUhFB50IMSw","Type":"arch_chip_health","Level":1},
			{"ID":"r2cs2Y2AReKY1gdmCH48tg","Type":"arch_chip_critchance","Level":1},
			{"ID":"wJwwrKOiTZCIRTjjkAGDbw","Type":"arch_chip_weaponprocrating","Level":1},
			{"ID":"Bgcp2wBqS+qy1mgy6V0+aw","Type":"arch_chip_damagecryo","Level":1},
			{"ID":"QyaRGHrwRcS2zWyyHMzczg","Type":"arch_chip_damagekinetic","Level":1},
			{"ID":"HLdb1bitSdGzIZIrMrpl2g","Type":"arch_chip_damageamount","Level":1},
			{"ID":"1Bep5at+QQKolp53o7CGXA","Type":"arch_chip_abilitydamageamount","Level":1},
			{"ID":"qTch9FcrTUaqHamCd/YaYA","Type":"arch_chip_critmultiplier","Level":1},
			{"ID":"sCmLHufnQ9Kp+frBVBKoMg","Type":"arch_chip_abilityprocrating","Level":1},
			{"ID":"9IG5HARuRYaAQ/JWIr4edQ","Type":"arch_chip_abilityrechargerate","Level":1},
			{"ID":"R1B55cxiRoCNFJNSIeVfPA","Type":"arch_chip_damageelectric","Level":1},
			{"ID":"wMLrhk2qSKqp376EmxGCGw","Type":"arch_chip_abilitycritmultiplier","Level":1},
			{"ID":"NUfda+v1SYWYXW1t+mQRhQ","Type":"arch_chip_shieldsize","Level":1}
		],
		"Loadouts":[
			{"ID":"F1CXpGdISIaffb2tzJuXNw","Character":{"ID":"aCH5ErdTT3WpO2slry3C/A","Chips":[]},"Weapons":[{"ID":"ynoLtQ3sSWSr/LMX2sRtwg","Chips":[]},{"ID":"G5CgC0U6QXKDkS7xTMTfNQ","Chips":[]},{"ID":"tEOWLnjSTrGGWsxWLyYqEA","Chips":[]}],"Level":0},
			{"ID":"1fvQKbKASYaFebQMu7pAjg","Character":{"ID":"4N+eLrYDSwOd3UF+4me/dQ","Chips":[]},"Weapons":[{"ID":"jBTb/knSRpWh0dubaAuKnw","Chips":[]},{"ID":"tQKsPTOKTYut/Kx/rEvK1Q","Chips":[]},{"ID":"yGWxtSA3Smm2w5Hia8wAvQ","Chips":[]}],"Level":0},
			{"ID":"LP2RB9ksTsCF7F4lu3dOhQ","Character":{"ID":"yZr+YD0CStG4lziRLw/Yhw","Chips":[]},"Weapons":[{"ID":"RktXYR3/SGyXYTmu1gt8Hg","Chips":[]},{"ID":"LpSjS2yzREetJsBumleUZw","Chips":[]},{"ID":"6VT8C97ERHakeeKz97tVCA","Chips":[]}],"Level":0},
			{"ID":"2ejRD/73QJ2AXt8gt+ZM2w","Character":{"ID":"UVF+DHRmSR6WJ7r720zeHQ","Chips":[]},"Weapons":[{"ID":"ynoLtQ3sSWSr/LMX2sRtwg","Chips":[]},{"ID":"G5CgC0U6QXKDkS7xTMTfNQ","Chips":[]},{"ID":"AQotA74EQ1Gfa1ksglKosg","Chips":[]}],"Level":0},
			{"ID":"u5HH+4pVTSub2HQ4K7MXTQ","Character":{"ID":"Sw0i5IIlQfiDHJlD3D5/6w","Chips":[]},"Weapons":[{"ID":"ynoLtQ3sSWSr/LMX2sRtwg","Chips":[]},{"ID":"G5CgC0U6QXKDkS7xTMTfNQ","Chips":[]},{"ID":"AQotA74EQ1Gfa1ksglKosg","Chips":[]}],"Level":0},
			{"ID":"eGQcfQT8SU2kE7y11yKZAA","Character":{"ID":"KtY4l4/lS6aS0TY2i1pnVw","Chips":[]},"Weapons":[{"ID":"ynoLtQ3sSWSr/LMX2sRtwg","Chips":[]},{"ID":"G5CgC0U6QXKDkS7xTMTfNQ","Chips":[]},{"ID":"AQotA74EQ1Gfa1ksglKosg","Chips":[]}],"Level":0},
			{"ID":"HltoRcNbQxiYcanV0KbrlA","Character":{"ID":"xnEqPhkgSty+5jcssTeafA","Chips":[]},"Weapons":[{"ID":"ynoLtQ3sSWSr/LMX2sRtwg","Chips":[]},{"ID":"G5CgC0U6QXKDkS7xTMTfNQ","Chips":[]},{"ID":"AQotA74EQ1Gfa1ksglKosg","Chips":[]}],"Level":0},
			{"ID":"GxPTxRfDTLGqZJ1D4q1o9g","Character":{"ID":"nurV5dyVRTSczi+ldS1qMw","Chips":[]},"Weapons":[{"ID":"ynoLtQ3sSWSr/LMX2sRtwg","Chips":[]},{"ID":"G5CgC0U6QXKDkS7xTMTfNQ","Chips":[]},{"ID":"AQotA74EQ1Gfa1ksglKosg","Chips":[]}],"Level":0}
		],
		"Resources":null
	}`)

func BenchmarkQuery(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := Execute(WithInput(bigFatJson)); err != nil {
			b.Fatal("Error encountered in BenchmarkQuery Execute()", err)
		}
	}
}

func TestQuery(t *testing.T) {
	expectedResult := `"SelectedLoadout":[1,2,3]`
	if output, err := Execute(
		WithInput(bigFatJson),
		WithSetField([]byte(`SelectedLoadout`), []byte(`[1,2,3]`)),
	); err != nil {
		t.Error(err)
	} else if !strings.Contains(string(output), expectedResult) {
		t.Fatalf("SetField was not successfully applied to SelectedLoadout in bigFatJson; expected %v",
			expectedResult)
	}
}

func TestExecute(t *testing.T) {
	type testCase struct {
		path           []byte
		pathFrom       []byte
		pathTo         []byte
		data           []byte
		value          []byte
		optionType     string
		expectedResult []byte
		expectedErr    error
	}
	testCases := []testCase{
		//<editor-fold desc="WithArrayPushBack">
		{
			optionType:     "WithArrayPushBack",
			path:           []byte(`a`),
			data:           []byte(`{"a":[{"b":"c"}]}`),
			value:          []byte(`{"c":"d"}`),
			expectedResult: []byte(`{"a":[{"b":"c"},{"c":"d"}]}`),
			expectedErr:    nil,
		},
		{
			optionType:     "WithArrayPushBack",
			path:           []byte(`.`),
			data:           []byte(`[{"a":"b"}]`),
			value:          []byte(`{"c":"d"}`),
			expectedResult: []byte(`[{"a":"b"},{"c":"d"}]`),
			expectedErr:    nil,
		},
		{
			optionType:  "WithArrayPushBack",
			path:        []byte(`.`),
			data:        []byte(`{"a":{"b":"c"}}`),
			value:       []byte(`{"c":"d"}`),
			expectedErr: ErrInvalidDestination,
		},
		{
			optionType:  "WithArrayPushBack",
			path:        []byte(`.`),
			data:        []byte(`{"a":"b"}`),
			value:       []byte(`{"c":"d"}`),
			expectedErr: ErrInvalidDestination,
		},
		//</editor-fold>
		//<editor-fold desc="WithArrayPushFront">
		{
			optionType:     "WithArrayPushFront",
			path:           []byte(`a`),
			data:           []byte(`{"a":[{"b":"c"}]}`),
			value:          []byte(`{"c":"d"}`),
			expectedResult: []byte(`{"a":[{"c":"d"},{"b":"c"}]}`),
			expectedErr:    nil,
		},
		{
			optionType:     "WithArrayPushFront",
			path:           []byte(`.`),
			data:           []byte(`[{"a":"b"}]`),
			value:          []byte(`{"c":"d"}`),
			expectedResult: []byte(`[{"a":"b"},{"c":"d"}]`),
			expectedErr:    ErrInvalidDestination,
		},
		{
			optionType:  "WithArrayPushFront",
			path:        []byte(`.`),
			data:        []byte(`{"a":{"b":"c"}}`),
			value:       []byte(`{"c":"d"}`),
			expectedErr: ErrInvalidDestination,
		},
		{
			optionType:  "WithArrayPushFront",
			path:        []byte(`.`),
			data:        []byte(`{"a":"b"}`),
			value:       []byte(`{"c":"d"}`),
			expectedErr: ErrInvalidDestination,
		},
		//</editor-fold>
		//<editor-fold desc="WithDelete">
		{
			optionType:     "WithDelete",
			path:           []byte(`b`),
			data:           []byte(`{"a":[{"c":"d"}], "b":[{"e":"f"}]}`),
			expectedResult: []byte(`{"a":[{"c":"d"}]}`),
			expectedErr:    nil,
		},
		{
			optionType:     "WithDelete",
			path:           []byte(`a`),
			data:           []byte(`{"a":[{"c":"d"}], "b":[{"e":"f"}]}`),
			expectedResult: []byte(`{"b":[{"e":"f"}]}`),
			expectedErr:    nil,
		},
		//</editor-fold>
		//<editor-fold desc="WithSwap">
		{
			optionType:     "WithSwap",
			pathFrom:       []byte(`a`),
			pathTo:         []byte(`b`),
			data:           []byte(`{"a":{"c":"d"}, "b":{"e":"f"}}`),
			expectedResult: []byte(`{"a":{"e":"f"},"b":{"c":"d"}}`),
			expectedErr:    nil,
		},
		{
			optionType:     "WithSwap",
			pathFrom:       []byte(`a`),
			pathTo:         []byte(`b`),
			data:           []byte(`{"a":[{"c":"d"}], "b":[{"e":"f"}]}`),
			expectedResult: []byte(`{"a":[{"e":"f"}],"b":[{"c":"d"}]}`),
			expectedErr:    nil,
		},
		//</editor-fold>
		//<editor-fold desc="WithMove">
		{
			optionType:     "WithMove",
			pathFrom:       []byte(`a`),
			pathTo:         []byte(`b`),
			data:           []byte(`{"a":{"c":"d"}, "b":""}`),
			expectedResult: []byte(`{"b":{"c":"d"}}`),
			expectedErr:    nil,
		},
		{
			optionType:     "WithMove",
			pathFrom:       []byte(`a`),
			pathTo:         []byte(`b`),
			data:           []byte(`{"a":[{"c":"d"}], "b":""}`),
			expectedResult: []byte(`{"b":[{"c":"d"}]}`),
			expectedErr:    nil,
		},
		//</editor-fold>
		//<editor-fold desc="WithCopy">
		{
			optionType:     "WithCopy",
			pathFrom:       []byte(`a`),
			pathTo:         []byte(`b`),
			data:           []byte(`{"a":{"c":"d"}, "b":""}`),
			expectedResult: []byte(`{"a":{"c":"d"},"b":{"c":"d"}}`),
			expectedErr:    nil,
		},
		{
			optionType:     "WithCopy",
			pathFrom:       []byte(`b`),
			pathTo:         []byte(`a`),
			data:           []byte(`{"a":"", "b":{"c":"d"}}`),
			expectedResult: []byte(`{"a":{"c":"d"},"b":{"c":"d"}}`),
			expectedErr:    nil,
		},
		{
			optionType:     "WithCopy",
			pathFrom:       []byte(`a`),
			pathTo:         []byte(`b`),
			data:           []byte(`{"a":[{"c":"d"}], "b":""}`),
			expectedResult: []byte(`{"a":[{"c":"d"}],"b":[{"c":"d"}]}`),
			expectedErr:    nil,
		},
		{
			optionType:     "WithCopy",
			pathFrom:       []byte(`b`),
			pathTo:         []byte(`a`),
			data:           []byte(`{"a":"", "b":[{"c":"d"}]}`),
			expectedResult: []byte(`{"a":[{"c":"d"}],"b":[{"c":"d"}]}`),
			expectedErr:    nil,
		},
		{
			optionType:  "WithCopy",
			pathFrom:    []byte(`a`),
			pathTo:      []byte(`.`),
			data:        []byte(`{"a":"", "b":{"c":"d"}}`),
			expectedErr: ErrInvalidDestination,
		},
		{
			optionType:  "WithCopy",
			pathFrom:    []byte(`.`),
			pathTo:      []byte(`b`),
			data:        []byte(`{"a":"", "b":{"c":"d"}}`),
			expectedErr: ErrInvalidDestination,
		},
		//</editor-fold>
		//<editor-fold desc="WithArrayInsert">
		{
			optionType:     "WithArrayInsert",
			path:           []byte(`a[1]`),
			value:          []byte(`{"g":"h"}`),
			data:           []byte(`{"a":[{"b":"c"}],"d":[{"e":"f"}]}`),
			expectedResult: []byte(`{"a":[{"b":"c"},{"g":"h"}],"d":[{"e":"f"}]}`),
			expectedErr:    nil,
		},
		{
			optionType:     "WithArrayInsert",
			path:           []byte(`d[1]`),
			value:          []byte(`{"g":"h"}`),
			data:           []byte(`{"a":[{"b":"c"}],"d":[{"e":"f"}]}`),
			expectedResult: []byte(`{"a":[{"b":"c"}],"d":[{"e":"f"},{"g":"h"}]}`),
			expectedErr:    nil,
		},
		{
			optionType:  "WithArrayInsert",
			path:        []byte(`.`),
			value:       []byte(`{"g":"h"}`),
			data:        []byte(`{"a":[{"b":"c"}],"d":[{"e":"f"}]}`),
			expectedErr: ErrInvalidDestination,
		},
		{
			optionType:  "WithArrayInsert",
			path:        []byte(`a`),
			value:       []byte(`{"g":"h"}`),
			data:        []byte(`{"a":[{"b":"c"}],"d":[{"e":"f"}]}`),
			expectedErr: ErrInvalidDestination,
		},
		{
			optionType:  "WithArrayInsert",
			path:        []byte(`d`),
			value:       []byte(`{"g":"h"}`),
			data:        []byte(`{"a":[{"b":"c"}],"d":[{"e":"f"}]}`),
			expectedErr: ErrInvalidDestination,
		},
		//</editor-fold>
	}

	for i, tc := range testCases {
		var option Option
		switch tc.optionType {
		case "WithArrayPushBack":
			option = withArrayPushBackTest(tc.path, tc.value)
		case "WithArrayPushFront":
			option = withArrayPushFrontTest(tc.path, tc.value)
		case "WithArrayInsert":
			option = withArrayInsertTest(tc.path, tc.value)
		case "WithDelete":
			option = withDeleteTest(tc.path)
		case "WithMove":
			option = withMoveTest(tc.pathFrom, tc.pathTo)
		case "WithCopy":
			option = withCopyTest(tc.pathFrom, tc.pathTo)
		case "WithSwap":
			option = withSwapTest(tc.pathFrom, tc.pathTo)
		default:
			t.Error("Unrecognized option ", tc.optionType, " provided in test case #", i+1)
		}
		if output, err := Execute(
			withInputTest(tc.data),
			option,
		); err != nil {
			if err != tc.expectedErr {
				t.Error("Unexpected error encountered in test case #", i+1, ": received", err, "where",
					tc.expectedErr, "was expected")
			} else if output != nil {
				t.Error("Output was provided when an error was encountered in test case #", i+1)
			}
		} else if tc.expectedErr != nil {
			t.Error("Expected error", tc.expectedErr, "was not encountered where one was expected in test case #", i+1)
		} else if output == nil {
			t.Error("No output was provided when no error was encountered in test case #", i+1)
		} else if bytes.Compare(output, tc.expectedResult) != 0 {
			t.Error("Incorrect output ", string(output),
				" where ", string(tc.expectedResult), " was expected in test case #", i+1)
		}
	}
}

func withInputTest(data interface{}) Option {
	return func(o *options) {
		switch dt := data.(type) {
		case string:
			o.input = []byte(dt)
		case []byte:
			o.input = dt
		case map[string]interface{}:
			o.input, _ = jsonx.Marshal(data)
		}
	}
}

func withArrayInsertTest(path []byte, value []byte) Option {
	return func(o *options) {
		o.actions = append(o.actions, func(root interface{}) error {
			if nodes, err := pathToProperty(root, path); err != nil {
				return err
			} else if len(nodes) < 2 {
				return ErrInvalidDestination
			} else if v, err := jsonx.ParseAndReturn(value); err != nil {
				return err
			} else {
				n := nodes[0]
				n1 := nodes[1]
				switch p := n1.Parent.(type) {
				case map[string]interface{}:
					if n1.Key == "" || n.Index == -1 {
						return ErrInvalidDestination
					} else {
						p[n1.Key] = append(p[n1.Key].([]interface{}), v)
						copy(p[n1.Key].([]interface{})[n.Index+1:], p[n1.Key].([]interface{})[n.Index:])
						p[n1.Key].([]interface{})[n.Index] = v
					}
				default:
					return ErrInvalidDestination
				}
			}
			return nil
		})
	}
}

func withArrayPushFrontTest(path []byte, value []byte) Option {
	return func(o *options) {
		o.actions = append(o.actions, func(root interface{}) error {
			if nodes, err := pathToProperty(root, path); err != nil {
				return err
			} else if v, err := jsonx.ParseAndReturn(value); err != nil {
				return err
			} else {
				n := nodes[0]
				switch p := n.Parent.(type) {
				case map[string]interface{}:
					if n.Key == "" {
						return ErrInvalidDestination
					} else {
						p[n.Key] = append([]interface{}{v}, p[n.Key].([]interface{})...)
					}
				default:
					return ErrInvalidDestination
				}
			}
			return nil
		})
	}
}

func withArrayPushBackTest(path []byte, value []byte) Option {
	return func(o *options) {
		o.actions = append(o.actions, func(root interface{}) error {
			if nodes, err := pathToProperty(root, path); err != nil {
				return err
			} else if v, err := jsonx.ParseAndReturn(value); err != nil {
				return err
			} else {
				n := nodes[0]
				switch p := n.Parent.(type) {
				case map[string]interface{}:
					if n.Key == "" {
						return ErrInvalidDestination
					} else {
						p[n.Key] = append(p[n.Key].([]interface{}), v)
					}
				case *interface{}:
					switch (*p).(type) {
					case []interface{}:
						*p = append((*p).([]interface{}), v)
					default:
						return ErrInvalidDestination
					}
				default:
					return ErrInvalidDestination
				}
			}
			return nil
		})
	}
}

func withDeleteTest(path []byte) Option {
	return func(o *options) {
		o.actions = append(o.actions, func(root interface{}) error {
			if nodes, err := pathToProperty(root, path); err != nil {
				return err
			} else {
				n := nodes[0]
				switch p := n.Parent.(type) {
				case map[string]interface{}:
					if n.Key == "" {
						return ErrInvalidDestination
					} else {
						delete(p, n.Key)
					}
				case []interface{}:
					n1 := nodes[1]
					if n.Index == -1 || n1.Key == "" {
						return ErrInvalidDestination
					} else {
						switch p := n1.Parent.(type) {
						case map[string]interface{}:
							p[n1.Key] = append(p[n1.Key].([]interface{})[:n.Index], p[n1.Key].([]interface{})[n.Index+1:]...)
						case []interface{}:
							p[n1.Index] = append(p[n1.Index].([]interface{})[:n.Index], p[n1.Index].([]interface{})[n.Index+1:]...)
						default:
							return ErrInvalidDestination
						}
					}

				default:
					return ErrInvalidDestination
				}
			}

			return nil
		})
	}
}

func withCopyTest(pathFrom []byte, pathTo []byte) Option {
	return func(o *options) {
		o.actions = append(o.actions, func(root interface{}) error {
			if nodesFrom, err := pathToProperty(root, pathFrom); err != nil {
				return err
			} else if nodesTo, err := pathToProperty(root, pathTo); err != nil {
				return err
			} else {
				to := nodesTo[0]
				from := nodesFrom[0]
				switch p := to.Parent.(type) {
				case map[string]interface{}:
					switch f := from.Parent.(type) {
					case map[string]interface{}:
						if to.Key == "" || from.Key == "" {
							return ErrInvalidDestination
						} else {
							p[to.Key] = f[from.Key]
						}
					case []interface{}:
						if to.Key == "" || from.Index == -1 {
							return ErrInvalidDestination
						} else {
							p[to.Key] = f[from.Index]
						}
					default:
						return ErrInvalidDestination
					}
				case []interface{}:
					switch f := from.Parent.(type) {
					case map[string]interface{}:
						if to.Index == -1 || from.Key == "" {
							return ErrInvalidDestination
						} else {
							p[to.Index] = f[from.Key]
						}
					case []interface{}:
						if to.Index == -1 || from.Index == -1 {
							return ErrInvalidDestination
						} else {
							p[to.Index] = f[from.Index]
						}
					default:
						return ErrInvalidDestination
					}
				default:
					return ErrInvalidDestination
				}
			}
			return nil
		})
	}
}

func withMoveTest(pathFrom []byte, pathTo []byte) Option {
	return func(o *options) {
		o.actions = append(o.actions, func(root interface{}) error {
			if nodesFrom, err := pathToProperty(root, pathFrom); err != nil {
				return err
			} else if nodesTo, err := pathToProperty(root, pathTo); err != nil {
				return err
			} else {
				to := nodesTo[0]
				from := nodesFrom[0]
				switch p := to.Parent.(type) {
				case map[string]interface{}:
					switch f := from.Parent.(type) {
					case map[string]interface{}:
						if to.Key == "" || from.Key == "" {
							return ErrInvalidDestination
						} else {
							p[to.Key] = f[from.Key]
							delete(f, from.Key)
						}
					case []interface{}:
						if to.Key == "" || from.Index == -1 {
							return ErrInvalidDestination
						} else {
							p[to.Key] = f[from.Index]
							f = append(f[:from.Index], f[from.Index+1:]...)
						}
					default:
						return ErrInvalidDestination
					}
				case []interface{}:
					switch f := from.Parent.(type) {
					case map[string]interface{}:
						if to.Index == -1 || from.Key == "" {
							return ErrInvalidDestination
						} else {
							p[to.Index] = f[from.Key]
							delete(f, from.Key)
						}
					case []interface{}:
						if to.Index == -1 || from.Index == -1 {
							return ErrInvalidDestination
						} else {
							p[to.Index] = f[from.Index]
							f = append(f[:from.Index], f[from.Index+1:]...)
						}
					default:
						return ErrInvalidDestination
					}
				default:
					return ErrInvalidDestination
				}
			}
			return nil
		})
	}
}

func withSwapTest(pathFrom []byte, pathTo []byte) Option {
	return func(o *options) {
		o.actions = append(o.actions, func(root interface{}) error {
			if nodesFrom, err := pathToProperty(root, pathFrom); err != nil {
				return err
			} else if nodesTo, err := pathToProperty(root, pathTo); err != nil {
				return err
			} else {
				to := nodesTo[0]
				from := nodesFrom[0]
				switch p := to.Parent.(type) {
				case map[string]interface{}:
					switch f := from.Parent.(type) {
					case map[string]interface{}:
						if to.Key == "" || from.Key == "" {
							return ErrInvalidDestination
						} else {
							p[to.Key], f[from.Key] = f[from.Key], p[to.Key]
						}
					case []interface{}:
						if to.Key == "" || from.Index == -1 {
							return ErrInvalidDestination
						} else {
							p[to.Key], f[from.Index] = f[from.Index], p[to.Key]
						}
					default:
						return ErrInvalidDestination
					}
				case []interface{}:
					switch f := from.Parent.(type) {
					case map[string]interface{}:
						if to.Index == -1 || from.Key == "" {
							return ErrInvalidDestination
						} else {
							p[to.Index], f[from.Key] = f[from.Key], p[to.Index]
						}
					case []interface{}:
						if to.Index == -1 || from.Index == -1 {
							return ErrInvalidDestination
						} else {
							p[to.Index], f[from.Index] = f[from.Index], p[to.Index]
						}
					default:
						return ErrInvalidDestination
					}
				default:
					return ErrInvalidDestination
				}
			}
			return nil
		})
	}
}
