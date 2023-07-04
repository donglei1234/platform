package db

import (
	pb "github.com/donglei1234/platform/services/guild/generated/grpc/go/guild/api"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

const (
	table_guild        = "platform:guild:guilds"
	table_member       = "platform:guild:members"
	table_guild_max_id = "platform:guild:max"
	table_guild_num    = "platform:guild:num"
	Guild_num          = "guild_num"
	Userid             = "user_id"
	Profileid          = "profile_id"
	Guildid            = "guild_id"
	Guildname          = "guild_name"
	Guildnotice        = "guild_notice"
	Guildicon          = "guild_icon"
	Guildattributes    = "guild_attributes"
	Guildapply         = "guild_apply"
)

type Database struct {
	logger      *zap.Logger
	redisClient *redis.Client
}

type GuildM struct {
	Uguildname   string
	Uguildid     string
	Uguildicon   string
	Uguildnotice string
}

type UMember struct {
	UserId    string
	ProfileId string
}

func OpenDatabase(l *zap.Logger, addr string, pwd string) *Database {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd, // no password set
		DB:       0,   // use default DB
	})
	return &Database{
		l,
		client,
	}
}

// 按guildid做模糊搜索
func (d *Database) ScanByUId(table string) []string {
	var cursor uint64 = 0
	var tmp []string
	var table_list []string
	for {
		tmp, cursor, _ = d.redisClient.Scan(cursor, table, 10).Result()
		if len(tmp) > 0 {
			table_list = append(table_list, tmp...)
		}
		if cursor == 0 {
			break
		}
	}
	return table_list
}

// 判断表是否存在
func (d *Database) IsExist(table string) error {
	if exist := d.redisClient.Exists(table).Val(); exist == 0 {
		return NoExistGet
	} else {
		return nil
	}
}

// 新建工会时设置最大id
func (d *Database) GetGuildId(appId string) (string, error) {
	pip := d.redisClient.Pipeline()
	if exist := d.redisClient.Exists(table_guild_max_id + ":" + appId).Val(); exist == 0 {
		pip.HSet(table_guild_max_id+":"+appId, Guildid, "1")
		pip.HSet(table_guild_num+":"+appId, Guild_num, 1)
		pip.Exec()
		return "0", nil
	} else {
		pip.HGet(table_guild_max_id+":"+appId, Guildid).Val()
		pip.HGet(table_guild_num+":"+appId, Guild_num).Val()
		cmders, _ := pip.Exec()
		guildid, _ := cmders[0].(*redis.StringCmd).Result()
		int_id, _ := strconv.Atoi(guildid)
		num_s, _ := cmders[1].(*redis.StringCmd).Result()
		num, _ := strconv.Atoi(num_s)
		pip.HSet(table_guild_num+":"+appId, Guild_num, num+1)
		pip.HSet(table_guild_max_id+":"+appId, Guildid, int_id+1)
		pip.Exec()
		return guildid, nil
	}
}

// 加入工会
func (d *Database) JoinGuild(appId, guildId, profileId, userId, guildAttr string) error {
	user_id := userId
	// 修改Guild表
	tables := d.ScanByUId(table_guild + ":" + appId + ":" + guildId + ":*")
	var table string
	for _, v := range tables {
		if v[len(v)-2:] == ":d" {
			continue
		} else {
			table = v
			break
		}
	}
	if exist := d.redisClient.Exists(table).Val(); exist == 0 {
		d.logger.Error(NoExistJoin.Error())
		return NoExistJoin
	} else {
		d.redisClient.HMSet(table, map[string]interface{}{
			Guildattributes: guildAttr,
		})
	}
	// 修改member表
	d.redisClient.HMSet(table_member+":"+appId+":"+guildId+":"+userId, map[string]interface{}{
		Userid:    user_id,
		Profileid: profileId})
	return nil
}

// 创建或修改工会信息
func (d *Database) Cre_Mod_Guild(appId, profileId, userId, guildId, name, attribute, icon, notice, mode string) (string, error) {
	if mode == "create" {
		// 创建guild表
		var guildid string
		guildid, _ = d.GetGuildId(appId)
		if exist := d.redisClient.Exists(table_guild + ":" + appId + ":" + guildid +
			":" + name).Val(); exist == 0 {
			d.redisClient.HMSet(table_guild+":"+appId+":"+guildid+
				":"+name, map[string]interface{}{
				Guildid:         guildid,
				Guildname:       name,
				Guildnotice:     notice,
				Guildicon:       icon,
				Guildattributes: attribute,
				Guildapply:      "",
			})
		} else {
			d.logger.Error(CreateExistGuild.Error())
		}
		// 创建member表
		if exist := d.redisClient.Exists(table_member + ":" + appId + ":" + guildid + ":" + userId).Val(); exist == 0 {
			d.redisClient.HMSet(table_member+":"+appId+":"+guildid+":"+userId, map[string]interface{}{
				Userid:    userId,
				Profileid: profileId,
			})
		} else {
			d.logger.Error(ExistMemberTable.Error())
		}
		return guildid, nil
	} else {
		// 仅修改guild信息
		tables := d.ScanByUId(table_guild + ":" + appId + ":" + guildId + ":*")
		var table string
		for _, v := range tables {
			if v[len(v)-2:] == ":d" {
				continue
			} else {
				table = v
				break
			}
		}
		info := d.redisClient.HGetAll(table).Val()
		var tmp map[string]interface{}
		tmp = make(map[string]interface{})
		for k, v := range info {
			tmp[k] = v
		}
		//更新表名中的guildName字段
		pip := d.redisClient.Pipeline()
		pip.Del(table)
		pip.HMSet(table_guild+":"+appId+":"+guildId+":"+
			name, tmp)
		//覆盖掉原来的信息
		pip.HMSet(table_guild+":"+appId+":"+guildId+":"+
			name, map[string]interface{}{
			Guildid:         guildId,
			Guildname:       name,
			Guildnotice:     notice,
			Guildicon:       icon,
			Guildattributes: attribute,
		})
		pip.Exec()
		return "", nil
	}
}

// 获取工会信息，被SearchGuild调用
func (d *Database) GetGuildMessage(table string) (*GuildM, error) {
	if exist := d.redisClient.Exists(table).Val(); exist == 0 {
		return nil, nil
	} else {
		tmp := d.redisClient.HGetAll(table).Val()
		res := GuildM{Uguildid: tmp[Guildid],
			Uguildname:   tmp[Guildname],
			Uguildicon:   tmp[Guildicon],
			Uguildnotice: tmp[Guildnotice]}
		return &res, nil
	}
}

// 查询工会、若无查询内容则返回前num个工会的list
func (d *Database) SearchGuild(appId, searchInput string, number int64) ([]*GuildM, error) {
	var message_list []*GuildM
	if len(searchInput) != 0 {
		table := table_guild + ":" + appId + ":" + "*" + searchInput + "*"
		table_list := d.ScanByUId(table)
		for _, t := range table_list {
			if t[len(t)-2:] == ":d" {
				continue
			} else {
				guild_message, _ := d.GetGuildMessage(t)
				if guild_message == nil {
					continue
				} else {
					message_list = append(message_list, guild_message)
				}
			}
		}
	} else {
		guild_num_s := d.redisClient.HGet(table_guild_num+":"+appId, Guild_num).Val()
		guild_num, _ := strconv.Atoi(guild_num_s)
		guild_max_ids := d.redisClient.HGet(table_guild_max_id+":"+appId, Guildid).Val()
		guild_max_id, _ := strconv.Atoi(guild_max_ids)
		strnum := strconv.FormatInt(number, 10)
		num, _ := strconv.Atoi(strnum)
		if guild_num <= num {
			// 返回全部
			for i := 0; i <= guild_max_id; i++ {
				guild_message, _ := d.GetGuildMessage(table_guild + ":" + appId + ":" + strconv.Itoa(i))
				if guild_message != nil {
					message_list = append(message_list, guild_message)
				}
			}
		} else {
			rand.Seed(time.Now().Unix())
			for len(message_list) < num {
				guild_message, _ := d.GetGuildMessage(table_guild + ":" + appId + ":" + strconv.Itoa(rand.Intn(guild_max_id)))
				if guild_message != nil {
					message_list = append(message_list, guild_message)
				}
			}
		}
	}
	return message_list, nil
}

// 退出工会
func (d *Database) QuitGuild(appId, guildId, userId, guildAttr string) error {
	//var score int
	// 更新guild表
	tables := d.ScanByUId(table_guild + ":" + appId + ":" + guildId + ":*")
	var table string
	for _, v := range tables {
		if v[len(v)-2:] == ":d" {
			continue
		} else {
			table = v
			break
		}
	}
	d.redisClient.HMSet(table, map[string]interface{}{
		//Guildscore: score,
		Guildattributes: guildAttr,
	})
	// 更新member表
	d.redisClient.Del(table_member + ":" + appId + ":" + guildId + ":" + userId)
	return nil
}

// 更新工会成员职级，会长只能进行转让或转让+退出
func (d *Database) ChangeMemberGuild(appId, guildId, profileId, guildAttr, memberId, extraUserId string) error {
	profileid := profileId //目标profileid
	userId := memberId
	guildid := guildId
	n_profileid := d.redisClient.HGet(table_member+":"+appId+":"+guildid+":"+userId, Profileid).Val() //当前profileid
	if profileid == "-1" {
		//开除userId的会员
		if n_profileid == "0" {
			//会长退出
			d.PostGuild(appId, guildId, extraUserId, memberId, profileId) //转移权限
			d.QuitGuild(appId, guildId, memberId, guildAttr)              //退出工会
		} else {
			d.QuitGuild(appId, guildId, memberId, guildAttr)
		}
	} else {
		if profileid == "0" {
			//会长转让
			d.PostGuild(appId, guildId, extraUserId, memberId, profileId)
		} else {
			//更新guild
			tables := d.ScanByUId(table_guild + ":" + appId + ":" + guildid + "*")
			var table string
			for _, v := range tables {
				if v[len(v)-2:] == ":d" {
					continue
				} else {
					table = v
					break
				}
			}
			pip := d.redisClient.Pipeline()
			pip.HGet(table, Guildname).Val()
			Cmders, _ := pip.Exec()
			guildname, _ := Cmders[0].(*redis.StringCmd).Result()
			d.redisClient.HMSet(table_guild+":"+appId+":"+guildid+":"+guildname, map[string]interface{}{
				Guildattributes: guildAttr,
			})
			//更新member
			pip.HGet(table_member+":"+appId+":"+guildid+":"+memberId, Profileid)
			Cmders, _ = pip.Exec()
			profileid, _ = Cmders[0].(*redis.StringCmd).Result()
			if profileid == "0" {
				if len(extraUserId) == 0 {
					d.logger.Error(NoExistPost.Error())
					return NoExistPost
				}
				err := d.PostGuild(appId, guildId, extraUserId, memberId, profileId)
				if err != nil {
					return err
				}
				pip.Exec()
			} else {
				pip.HMSet(table_member+":"+appId+":"+guildid+":"+memberId, map[string]interface{}{
					Profileid: profileId,
				})
				pip.Exec()
			}
		}
	}
	return nil
}

// 转让工会，先做职位替换
func (d *Database) PostGuild(appId, guildId, extraUserId, memberId, profileId string) error {
	pip := d.redisClient.Pipeline()
	userId := memberId
	toid := extraUserId
	err := d.IsExist(table_member + ":" + appId + ":" + guildId + ":" + userId)
	if err != nil {
		return err
	}
	err = d.IsExist(table_member + ":" + appId + ":" + guildId + ":" + userId)
	if err != nil {
		return err
	}
	pip.HGet(table_member+":"+appId+":"+guildId+":"+userId, Profileid).Val()
	cmders, _ := pip.Exec()
	app_pro, _ := cmders[0].(*redis.StringCmd).Result()

	pip.HGet(table_member+":"+appId+":"+guildId+":"+toid, Profileid).Val()
	cmders, _ = pip.Exec()
	to_pro, _ := cmders[0].(*redis.StringCmd).Result()

	if len(profileId) != 0 {
		to_pro = profileId
	}

	// 修改member表
	table_app := table_member + ":" + appId + ":" + guildId + ":" + userId
	table_to := table_member + ":" + appId + ":" + guildId + ":" + toid

	pip.HMSet(table_app, map[string]interface{}{
		Profileid: to_pro,
	})
	pip.HMSet(table_to, map[string]interface{}{
		Profileid: app_pro,
	})
	pip.Exec()
	return nil
}

// 解散工会
func (d *Database) DelGuild(appId, guildId string) error {
	pip := d.redisClient.Pipeline()
	//修改guild表
	tables := d.ScanByUId(table_guild + ":" + appId + ":" + guildId + ":*")
	var table string
	for _, v := range tables {
		if v[len(v)-2:] == ":d" {
			continue
		} else {
			table = v
			break
		}
	}
	guildinfo := d.redisClient.HGetAll(table).Val()
	//暂存到:d后缀下的表
	d.redisClient.HMSet(table+":d", map[string]interface{}{
		Guildid:         guildinfo[Guildid],
		Guildattributes: guildinfo[Guildattributes],
	})
	pip.Del(table_guild + ":" + appId + ":" + guildId + ":" + guildinfo[Guildname])
	//更新工会数量
	pip.HGet(table_guild_num+":"+appId, Guild_num).Val()
	cmders, _ := pip.Exec()
	num_s, _ := cmders[1].(*redis.StringCmd).Result()
	num, _ := strconv.Atoi(num_s)
	pip.HSet(table_guild_num+":"+appId, Guild_num, num-1)
	pip.Exec()
	// 修改member部分
	tables = d.ScanByUId(table_member + ":" + appId + ":" + guildId + ":*")
	for _, table := range tables {
		//删除member表
		d.redisClient.Del(table)
	}
	return nil
}

// 获取工会内的用户信息list
func (d *Database) GetMemberList(appId, Idx string) ([]*UMember, error) {
	var memberlist []*UMember
	keys := d.ScanByUId(table_member + ":" + appId + ":" + Idx + ":*")
	for _, key := range keys {
		tmp := d.redisClient.HGetAll(key).Val()
		mem := UMember{UserId: tmp[Userid], ProfileId: tmp[Profileid]}
		memberlist = append(memberlist, &mem)
	}
	return memberlist, nil
}

// 工会申请
func (d *Database) Apply(appId, userId, guildId string) error {
	pip := d.redisClient.Pipeline()
	isrepeat := false
	tables := d.ScanByUId(table_guild + ":" + appId + ":" + guildId + ":*")
	if len(tables) == 0 {
		d.logger.Error(NoExistApp.Error())
		return NoExistApp
	}
	var table string
	for _, v := range tables {
		if v[len(v)-2:] == ":d" {
			continue
		} else {
			table = v
			break
		}
	}
	apply_info := pip.HGet(table, Guildapply).Val()
	var apply_list []string
	if len(apply_info) > 0 {
		apply_list = strings.Split(apply_info, "\001")
	}
	for _, v := range apply_list {
		if v == userId {
			isrepeat = true
			break
		}
	}
	if !isrepeat {
		if len(apply_list) > 0 {
			apply_list = append(apply_list, userId)
			apply_info = strings.Join(apply_list, "\001")
		} else {
			apply_info = userId
		}
		pip.HSet(table, Guildapply, apply_info)
	} else {
		return RepeatApp
	}

	pip.Exec()
	return nil
}

// 获取工会申请信息
func (d *Database) GetApply(appId, guildId string) (*pb.GetApplyResponse, error) {
	var table string
	tables := d.ScanByUId(table_guild + ":" + appId + ":" + guildId + ":*")
	for _, v := range tables {
		if v[len(v)-2:] == ":d" {
			continue
		} else {
			table = v
			break
		}
	}
	apply_info := d.redisClient.HGet(table, Guildapply).Val()
	apply_list := strings.Split(apply_info, "\001")
	//var response pb.GetApplyResponse
	//response.UserId = apply_list
	return &pb.GetApplyResponse{UserId: apply_list}, nil
}

// 批准或拒绝工会申请
func (d *Database) Reply(appId, applyId, guildId string, mode bool, profileId,
	userId, guildAttr string) error {
	tmp_tables := d.ScanByUId(table_member + ":" + appId + ":*:" + userId)
	var table string
	pip := d.redisClient.Pipeline()
	tables := d.ScanByUId(table_guild + ":" + appId + ":" + guildId + ":*")
	for _, v := range tables {
		if v[len(v)-2:] == ":d" {
			continue
		} else {
			table = v
		}
	}
	pip.HGet(table, Guildapply).Val()
	cmders, _ := pip.Exec()
	apply_info, _ := cmders[0].(*redis.StringCmd).Result()
	var new_apply []string
	apply_list := strings.Split(apply_info, "\001")
	for _, v := range apply_list {
		if v != applyId {
			new_apply = append(new_apply, v)
		}
	}
	// 仅更新guildapply
	if len(new_apply) == 0 {
		apply_info = ""
	} else {
		apply_info = strings.Join(new_apply, "\001")
	}
	if len(new_apply) == len(apply_list) {
		d.logger.Error(NoExistApp.Error())
		return NoExistApp
	}
	pip.HSet(table, Guildapply, apply_info)
	pip.Exec()
	if mode {
		d.JoinGuild(appId, guildId, profileId, userId, guildAttr)
	}
	if len(tmp_tables) > 0 {
		return UserHaveGuild
	}
	if mode {
		return nil
	} else {
		return RefuseApp
	}
}
