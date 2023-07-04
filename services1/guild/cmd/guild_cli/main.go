package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/abiosoft/ishell"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	pb "github.com/donglei1234/platform/services/guild/generated/grpc/go/guild/api"
	"github.com/donglei1234/platform/services/guild/pkg/client"
)

const (
	DefaultGuildUrl = "localhost:8081"
	DefaultUsername = "test"
)

var options struct {
	auth     string
	username string
}

var md = metadata.Pairs(
	"token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDUxNjkzMDUsInVpZCI6IjQyMjI3NzY0NzcifQ.322Plw1b1VzCYKiVw5DA3uw3_5tM1X_J1_bV8o_CU9A",
)
var ctx = metadata.NewOutgoingContext(context.Background(), md)

func main() {
	rootCmd := &cobra.Command{
		Use:   "guild_cli",
		Short: "Run a guild CLI",
	}

	rootCmd.PersistentFlags().StringVar(&options.auth, "guild", DefaultGuildUrl,
		"guild service (<host>:<port>)")
	rootCmd.PersistentFlags().StringVar(&options.username, "username", DefaultUsername, "username for authentication")
	{
		shell := &cobra.Command{
			Use:   "shell",
			Short: "Run an interactive guild service client",
			Run: func(cmd *cobra.Command, args []string) {
				shell()
			},
		}
		rootCmd.AddCommand(shell)
	}
	_ = rootCmd.Execute()
}

// ReadLine print tip and wait input.
func readLine(c *ishell.Context, val string) string {
	c.Print(val)
	return c.ReadLine()
}

// Info print text and wrap line.
func info(c *ishell.Context, msg string) {
	c.Printf("%s\n", msg)
}

// Infof printf prints to output using string format.
func infof(c *ishell.Context, format string, vals ...interface{}) {
	c.Printf(format, vals...)
}

// Warn printf warning message like "Warning: ..."
func warn(c *ishell.Context, err error) {
	c.Printf("Warning: %s\n", err)
}

// Die printf error message and exit.
func die(c *ishell.Context, err error) {
	c.Printf("Error: %s\n", err)
	os.Exit(1)
}

// Done printf done message and exit.
func done(c *ishell.Context, msg string) {
	c.Printf("Done: %s\n", msg)
	os.Exit(0)
}

func shell() {
	sh := ishell.New()

	sh.Println("guild Interactive Shell")

	sh.AddCmd(&ishell.Cmd{
		Name: "JoinGuild",
		Func: JoinGuild,
		Help: "join a guild",
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "CreateGuild",
		Func: CreateGuild,
		Help: "create a guild",
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "ModifyGuild",
		Func: ModifyGuild,
		Help: "modify Guild info",
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "SearchGuild",
		Func: SearchGuild,
		Help: "search guilds or get guild list",
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "DeleteGuild",
		Func: DeleteGuild,
		Help: "delete the guild",
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "ChangeMemberGuild",
		Func: ChangeMemberGuild,
		Help: "change the member profile level",
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "GetMember",
		Func: GetMember,
		Help: "get user or guild list",
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "Apply",
		Func: Apply,
		Help: "application",
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "Reply",
		Func: Reply,
		Help: "agree or disagree the application",
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "GetApply",
		Func: GetApply,
		Help: "get guild's apply list",
	})
	sh.Run()
}

func JoinGuild(c *ishell.Context) {
	info(c, "add func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewGuildClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")
	appid := readLine(c, "appId:")
	userid := readLine(c, "UserId:")
	guildid := readLine(c, "guildid:")
	profileid := readLine(c, "profileid:")
	guildattr := readLine(c, "guildattribute:")
	err = cli.JoinGuild(ctx, appid, userid, guildid, profileid, guildattr)
	if err != nil {
		die(c, err)
	} else {
		fmt.Println("Join success")
	}
}

func CreateGuild(c *ishell.Context) {
	info(c, "add func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewGuildClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")

	userid := readLine(c, "userid:")
	name := readLine(c, "guildname:")
	notice := readLine(c, "notice:")
	icon := readLine(c, "icon:")
	attribute := readLine(c, "guildattribute:")
	//mode := readLine(c, "create or modify:")
	appid := readLine(c, "appId:")
	profileid := readLine(c, "profileid:")
	idx, err0 := cli.CreateGuild(ctx, "", name, notice, icon, attribute,
		appid, userid, profileid)
	if err0 != nil {
		die(c, err0)
	} else {
		fmt.Println("create success and guildId is:", idx)
	}
}

func ModifyGuild(c *ishell.Context) {
	info(c, "add func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewGuildClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")

	guildid := readLine(c, "guildid:")
	appid := readLine(c, "appid:")
	name := readLine(c, "guildname:")
	notice := readLine(c, "notice:")
	icon := readLine(c, "icon:")
	attribute := readLine(c, "guildattribute:")
	err = cli.ModifyGuild(ctx, guildid, name, notice, icon, attribute, appid, "", "")
	if err != nil {
		die(c, err)
	} else {
		fmt.Println("modify success")
	}
}

func SearchGuild(c *ishell.Context) {
	info(c, "add func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewGuildClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")

	appid := readLine(c, "appid:")
	searchinput := readLine(c, "searchinput:")
	strnum := readLine(c, "list number:")
	number, _ := strconv.ParseInt(strnum, 10, 64)
	var response *pb.SearchResponse
	response, err = cli.SearchGuild(ctx, appid, searchinput, number)
	if err != nil {
		die(c, err)
	} else {
		fmt.Println("get guild list")
		fmt.Println(response.Guilds)
	}
}

func DeleteGuild(c *ishell.Context) {
	info(c, "add func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewGuildClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")
	appid := readLine(c, "appid:")
	guildid := readLine(c, "guildid:")
	err = cli.DeleteGuild(ctx, appid, guildid)
	if err != nil {
		die(c, err)
	} else {
		fmt.Println("Delete success")
	}
}

func ChangeMemberGuild(c *ishell.Context) {
	info(c, "add func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewGuildClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")

	var extra_appid string

	appid := readLine(c, "appid:")
	memberid := readLine(c, "the userid you wanna change:")
	profileid := readLine(c, "the level you wanna change to:")
	guildattr := readLine(c, "guild attribute:")
	if profileid == "0" || profileid == "-1" {
		// 涉及会长职级变动
		// 当userid是会长或userid职级目标是会长时
		extra_appid = readLine(c, "when userid is master, extra_id is the candidate:")
	} else {
		extra_appid = ""
	}
	guildid := readLine(c, "guildid:")
	err = cli.ChangeMemberGuild(ctx, appid, memberid, profileid, guildattr, extra_appid, guildid)
	if err != nil {
		die(c, err)
	} else {
		fmt.Println("change member success")
	}
}

func GetMember(c *ishell.Context) {
	info(c, "add func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewGuildClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")
	appid := readLine(c, "appid:")
	idx := readLine(c, "guildid:")
	var response *pb.UserListResponse
	response, err = cli.GetMember(ctx, appid, idx)
	if err != nil {
		die(c, err)
	} else {
		fmt.Println("get list success")
		fmt.Println(response.Users)
	}
}

func Apply(c *ishell.Context) {
	info(c, "add func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewGuildClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")
	appid := readLine(c, "your appid:")
	userid := readLine(c, "userid:")
	guildid := readLine(c, "apply guildid:")
	err = cli.Apply(ctx, appid, userid, guildid)
	if err != nil {
		die(c, err)
	} else {
		fmt.Println("Apply successful")
	}
}

func Reply(c *ishell.Context) {
	info(c, "add func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewGuildClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")
	appid := readLine(c, "the administrator's appid:")
	applyid := readLine(c, "the application's userid:")
	mode_s := readLine(c, "agree or disagree:")
	guildid := readLine(c, "guild id:")
	profileid := readLine(c, "application's level when joining in:")
	guildattr := readLine(c, "guild's attribute:")
	var mode bool
	if mode_s == "agree" {
		mode = true
	} else {
		mode = false
	}
	err = cli.Reply(ctx, appid, applyid, mode, guildid, profileid, guildattr)
	if err != nil {
		die(c, err)
	} else {
		fmt.Println("Reply success")
	}
}

func GetApply(c *ishell.Context) {
	info(c, "add func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewGuildClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")
	appid := readLine(c, "your appid:")
	guildid := readLine(c, "your guildid, cannot be none:")
	var response *pb.GetApplyResponse
	response, err = cli.GetApply(ctx, appid, guildid)
	if err != nil {
		die(c, err)
	} else {
		fmt.Println(response.UserId)
	}
}
