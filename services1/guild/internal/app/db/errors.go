package db

import (
	"errors"
)

var (
	NoExistGet       = errors.New("Table not exist")
	NoExistDel       = errors.New("Delete data when data not exist")
	NoExistJoin      = errors.New("Join a not existed guild")
	NoExistApp       = errors.New("UserId not in apply list")
	NoExistPost      = errors.New("UserId to be posted not existed")
	NoMasterDel      = errors.New("Del a guild not belong to you")
	CreateExistGuild = errors.New("Guild id existed")
	ExistMemberTable = errors.New("Member table is existed")
	UserHaveGuild    = errors.New("User have already had a guild")
	RefuseApp        = errors.New("Refused application")
	RepeatApp        = errors.New("Repeat application")
)
