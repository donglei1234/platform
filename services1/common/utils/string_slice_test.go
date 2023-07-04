package utils

import "testing"

func TestStringSlice_Contains(t *testing.T) {
	ss := make(StringSlice, 0, 3)
	ss = append(ss, "one", "two", "three")

	if !ss.Contains("one") {
		t.Error("Test of present element \"one\" failed.")
	}

	if ss.Contains("four") {
		t.Error("Test of absent element \"four\" failed.")
	}
}

func TestStringSlice_Equals(t *testing.T) {
	ss := make(StringSlice, 0, 3)
	ss = append(ss, "one", "two", "three")
	ss2 := make(StringSlice, 0, 2)
	ss2 = append(ss2, "one", "two")

	if ss.Equals(ss2) {
		t.Error("Equals returns true for two StringSlices of different lengths.")
	}

	ss3 := make(StringSlice, 0, 3)
	ss3 = append(ss3, "one", "two", "four")

	if ss.Equals(ss3) {
		t.Error("Equals returns true for two StringSlices of equal lengths but different contents.")
	}

	ss4 := ss

	if !ss.Equals(ss4) {
		t.Error("Equals returns false for two StringSlices of equal lengths and contents.")
	}
}

func TestStringSlice_FilterExclude(t *testing.T) {
	ss := make(StringSlice, 0, 3)
	ss = append(ss, "one", "two", "three")

	if ss2, found := ss.FilterExclude("one"); !found {
		t.Error("Present element not found when attempting to filter it.")
	} else {
		if len(ss2) != 2 {
			t.Error("Incorrect length of filtered StringSlice - expected 2, got ", len(ss2))
		}
	}

	if ss2, found := ss.FilterExclude("four"); found {
		t.Error("Absent element found when attempting to filter it.")
	} else {
		if !ss2.Equals(ss) {
			t.Error("Result of no-op filter exclude is not equal to original StringSlice.")
		}
	}
}
