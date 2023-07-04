package db

import (
	"testing"
)

func TestDatabase_NewLeaderboard(t *testing.T) {
	db := OpenDatabase(nil, "localhost:6379", "")
	db.NewLeaderboard("ww4", "testList", "SUM", "DESCENDING", 0, 0)
}

func TestDatabase_GetLeaderboardSize(t *testing.T) {
	db := OpenDatabase(nil, "localhost:6379", "")
	r, err := db.GetLeaderboardSize("ww4", "testList")
	t.Log(r, err)
}

func TestDatabase_ResetLeaderboard(t *testing.T) {
	db := OpenDatabase(nil, "localhost:6379", "")
	err := db.ResetLeaderboard("ww4", "testList")
	t.Log(err)
}

func TestDatabase_UpdateScore(t *testing.T) {
	db := OpenDatabase(nil, "localhost:6379", "")
	err := db.UpdateScore("ww4", "testList", "cxm", 100.0)
	t.Log(err)
}

func TestDatabase_DeleteMember(t *testing.T) {
	db := OpenDatabase(nil, "localhost:6379", "")
	err := db.DeleteMember("ww4", "testList", "cxx")
	t.Log(err)
}

func TestDatabase_GetRankById(t *testing.T) {
	db := OpenDatabase(nil, "localhost:6379", "")
	r, err := db.GetRankById("ww4", "testList", "cxm")
	t.Log(r, err)
}

func TestDatabase_GetRankFromMToN(t *testing.T) {
	db := OpenDatabase(nil, "localhost:6379", "")
	r, err := db.GetRankFromMToN("ww4", "testList", 1, 3)
	t.Log(r, err)
}
