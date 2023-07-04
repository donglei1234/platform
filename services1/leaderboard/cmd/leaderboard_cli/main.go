package main

import (
	"context"
	"os"
	"strconv"

	"github.com/abiosoft/ishell"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	"github.com/donglei1234/platform/services/leaderboard/pkg/client"
)

const (
	DefaultLeaderboard = "localhost:8081"
	DefaultUsername    = "test"
)

var options struct {
	auth     string
	username string
}

var md = metadata.Pairs(
	"token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDQ1ODU1MjcsInVpZCI6IjM0NDUxNzAyNTAifQ.3dmpRwxRiKlTDBwe37OOksP3ee02AXw6lqJ_E6nyOvc",
)
var ctx = metadata.NewOutgoingContext(context.Background(), md)

func main() {
	rootCmd := &cobra.Command{
		Use:   "auth_cli",
		Short: "Run a auth CLI",
	}

	rootCmd.PersistentFlags().StringVar(&options.auth, "leaderboard", DefaultLeaderboard,
		"leaderboard service (<host>:<port>)")
	rootCmd.PersistentFlags().StringVar(&options.username, "username", DefaultUsername, "username for authentication")
	{
		shell := &cobra.Command{
			Use:   "shell",
			Short: "Run an interactive leaderboard service client",
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

	sh.Println("leaderboard Interactive Shell")

	sh.AddCmd(&ishell.Cmd{
		Name: "getTopK",
		Func: getTopK,
		Help: "get top k from list",
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "getMToN",
		Func: getMToN,
		Help: "get m to n from list",
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "getIdRank",
		Func: getIdRank,
		Help: "get id rank",
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "updateScore",
		Func: updateScore,
		Help: "update id score",
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "deleteMember",
		Func: deleteMember,
		Help: "delete member from list",
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "newLeaderboard",
		Func: newLeaderboard,
		Help: "new leaderboard",
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "getLeaderboardSize",
		Func: getLeaderboardSize,
		Help: "get list size",
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "resetLeaderboard",
		Func: resetLeaderboard,
		Help: "reset list",
	})
	sh.Run()
}

func getTopK(c *ishell.Context) {
	info(c, "add func ...")
	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewLeaderboardClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")
	//authName := readLine(c, "Name: ")

	appId := readLine(c, "appId: ")
	listId := readLine(c, "listId ")
	k, err := strconv.Atoi(readLine(c, "k: "))
	if err != nil {
		die(c, err)
	}
	list, err := cli.GetTopK(
		ctx,
		appId,
		listId,
		int32(k))
	if err != nil {
		die(c, err)
	}

	for _, v := range list.GetLeaderboard() {
		info(c, "id: "+v.GetId()+" score: "+strconv.Itoa(int(v.GetScore())))
	}
}

func getMToN(c *ishell.Context) {
	info(c, "add func ...")
	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewLeaderboardClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")
	//authName := readLine(c, "Name: ")

	appId := readLine(c, "appId: ")
	listId := readLine(c, "listId ")
	m, err := strconv.Atoi(readLine(c, "m: "))
	if err != nil {
		die(c, err)
	}
	n, err := strconv.Atoi(readLine(c, "n: "))
	if err != nil {
		die(c, err)
	}

	list, err := cli.GetMToN(
		ctx,
		appId,
		listId,
		int32(m),
		int32(n))
	if err != nil {
		die(c, err)
	}

	for _, v := range list.GetLeaderboard() {
		info(c, "id: "+v.GetId()+" score: "+strconv.Itoa(int(v.GetScore())))
	}
}

func getIdRank(c *ishell.Context) {
	info(c, "add func ...")
	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewLeaderboardClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")
	//authName := readLine(c, "Name: ")

	appId := readLine(c, "appId: ")
	listId := readLine(c, "listId ")
	id := readLine(c, "id: ")

	rank, err := cli.GetIdRank(
		ctx,
		appId,
		listId,
		id)
	if err != nil {
		die(c, err)
	}

	info(c, "rank: "+strconv.Itoa(int(rank.GetRank())))
}

func updateScore(c *ishell.Context) {
	info(c, "add func ...")
	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewLeaderboardClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")
	//authName := readLine(c, "Name: ")

	appId := readLine(c, "appId: ")
	listId := readLine(c, "listId ")
	id := readLine(c, "id: ")
	score, err := strconv.Atoi(readLine(c, "score: "))
	if err != nil {
		die(c, err)
	}

	_, err = cli.UpdateScore(
		ctx,
		appId,
		listId,
		id,
		int32(score))
	if err != nil {
		die(c, err)
	}
	info(c, "update success")
}

func deleteMember(c *ishell.Context) {
	info(c, "add func ...")
	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewLeaderboardClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")
	//authName := readLine(c, "Name: ")

	appId := readLine(c, "appId: ")
	listId := readLine(c, "listId ")
	id := readLine(c, "id: ")

	_, err = cli.DeleteMember(
		ctx,
		appId,
		listId,
		id)
	if err != nil {
		die(c, err)
	}
	info(c, "delete member success")
}

func newLeaderboard(c *ishell.Context) {
	info(c, "add func ...")
	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewLeaderboardClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")
	//authName := readLine(c, "Name: ")

	appId := readLine(c, "appId: ")
	listId := readLine(c, "listId ")
	method := readLine(c, "method ")
	order := readLine(c, "order ")
	reset, err := strconv.Atoi(readLine(c, "reset time: "))
	if err != nil {
		die(c, err)
	}
	update, err := strconv.Atoi(readLine(c, "update time: "))
	if err != nil {
		die(c, err)
	}

	_, err = cli.NewLeaderboard(
		ctx,
		appId,
		listId,
		method,
		order,
		int32(reset),
		int32(update))
	if err != nil {
		die(c, err)
	}
	info(c, "new leaderboard success")
}

func getLeaderboardSize(c *ishell.Context) {
	info(c, "add func ...")
	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewLeaderboardClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")
	//authName := readLine(c, "Name: ")

	appId := readLine(c, "appId: ")
	listId := readLine(c, "listId ")

	size, err := cli.GetLeaderBoardSize(
		ctx,
		appId,
		listId)
	if err != nil {
		die(c, err)
	}

	info(c, "size: "+strconv.Itoa(int(size.GetSize())))
}

func resetLeaderboard(c *ishell.Context) {
	info(c, "add func ...")
	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewLeaderboardClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")
	//authName := readLine(c, "Name: ")

	appId := readLine(c, "appId: ")
	listId := readLine(c, "listId ")

	_, err = cli.ResetLeaderboard(
		ctx,
		appId,
		listId)
	if err != nil {
		die(c, err)
	}

	info(c, "reset success")
}
