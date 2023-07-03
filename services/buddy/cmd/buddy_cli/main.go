package main

import (
	"context"
	"errors"
	"os"

	"github.com/abiosoft/ishell"
	aClient "github.com/donglei1234/platform/services/auth/pkg/client"
	"github.com/donglei1234/platform/services/buddy/pkg/client"
	"github.com/donglei1234/platform/services/buddy/pkg/tests"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	defaultBuddy     = "sr-development.dev.spacerouter.net:8081"
	defaultAuth      = "sr-development.dev.spacerouter.net:8081"
	defaultAuthAppId = "test"
	defaultUsername  = "test"
	defaultSecure    = false
)

var options struct {
	buddy    string
	auth     string
	appId    string
	username string
	secure   bool
}

func main() {

	rootCmd := &cobra.Command{
		Use:   "buddy_cli",
		Short: "Run a buddy CLI",
	}

	rootCmd.PersistentFlags().StringVar(&options.buddy, "buddy", defaultBuddy, "buddy service (<host>:<port>)")
	rootCmd.PersistentFlags().StringVar(&options.auth, "auth", defaultAuth, "authentication service (<host>:<port>)")
	rootCmd.PersistentFlags().StringVar(&options.appId, "appId", defaultAuthAppId, "appId for authentication")
	rootCmd.PersistentFlags().StringVar(&options.username, "username", defaultUsername, "username for authentication")
	rootCmd.PersistentFlags().BoolVar(&options.secure, "secure", defaultSecure, "if provided, connect securely")

	{
		shell := &cobra.Command{
			Use:   "shell",
			Short: "Run an interactive buddy service client",
			Run: func(cmd *cobra.Command, args []string) {
				shell()
			},
		}
		rootCmd.AddCommand(shell)
	}
	{
		buddyTestSuite := &cobra.Command{
			Use:   "buddyTestSuite",
			Short: "Run a buddy service test suite",
			Long:  "buddyTestSuite runs a test suite which adds and removes buddies from a buddy queue.",
			Run: func(cmd *cobra.Command, args []string) {
				if err := tests.BuddyServiceSuite(options.auth, options.buddy); err != nil {
					return
				}
			},
		}
		rootCmd.AddCommand(buddyTestSuite)
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

func authenticate(c *ishell.Context) (token string, err error) {
	info(c, "Authenticating...")
	if authClient, e := aClient.NewPublicClient(zap.L(), options.auth, options.secure); e != nil {
		err = e
	} else if t, e := authClient.AuthenticateWithUsername(
		context.Background(),
		options.appId,
		options.username,
	); e != nil {
		err = e
	} else {
		infof(c, "Authentication token: %s\n", t)
		token = t
	}

	return
}

func shell() {
	sh := ishell.New()

	sh.Println("buddy Interactive Shell")

	sh.AddCmd(&ishell.Cmd{
		Name: "get-private",
		Help: "run a get buddies client",
		Func: getPrivate,
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "watch-private",
		Help: "run a get buddies client",
		Func: watchPrivate,
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "add",
		Help: "run a add client",
		Func: add,
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "remove",
		Help: "run a remove client",
		Func: remove,
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "get",
		Help: "run a get buddies client",
		Func: get,
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "watch",
		Help: "run a watch buddies client",
		Func: watch,
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "reply",
		Help: "run a reply buddies client",
		Func: reply,
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "remark",
		Help: "run a remark buddies client",
		Func: remark,
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "updateSettings",
		Help: "run a client to update buddy settings",
		Func: updateSettings,
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "addBlockedUser",
		Help: "run a add blocked user client",
		Func: addBlockedUser,
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "removeBlockedUser",
		Help: "run a remove blocked user client",
		Func: removeBlockedUser,
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "addToRecentMet",
		Help: "run a add to recent met client",
		Func: addToRecentMet,
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "getBlockedList",
		Help: "run a get blocked list private client",
		Func: getBlockedList,
	})
	sh.Run()
}

func add(c *ishell.Context) {
	info(c, "add func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	authToken, err := authenticate(c)
	if err != nil {
		die(c, err)
	}

	info(c, "Connecting to buddy service ...")
	cli, err := client.NewPublicClient(options.buddy, options.secure)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add buddy info ...")
	buddyName := readLine(c, "BuddyName: ")
	text := readLine(c, "RequestText: ")

	err = cli.AddBuddy(
		auth.ContextWithToken(context.Background(), authToken),
		buddyName,
		text,
	)

	if err != nil {
		die(c, err)
	}
}

func remove(c *ishell.Context) {
	info(c, "remove func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	authToken, err := authenticate(c)
	if err != nil {
		die(c, err)
	}

	info(c, "Connecting to buddy service ...")
	cli, err := client.NewPublicClient(options.buddy, options.secure)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter remove buddy info ...")
	buddyName := readLine(c, "BuddyName: ")

	err = cli.RemoveBuddy(
		auth.ContextWithToken(context.Background(), authToken),
		buddyName,
	)

	if err != nil {
		die(c, err)
	}
}

func getPrivate(c *ishell.Context) {
	info(c, "private get func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	authToken, err := authenticate(c)
	if err != nil {
		die(c, err)
	}

	info(c, "Connecting to buddy private service ...")
	cli, err := client.NewPrivateClient(options.buddy, options.secure)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	buddies, err := cli.GetProfileBuddies(
		auth.ContextWithToken(context.Background(), authToken),
		options.username,
		options.appId,
	)
	if err != nil {
		die(c, err)
	} else {
		infof(c, "buddy data receive success...\nbuddy list is %+v", buddies)
	}
}

func watchPrivate(c *ishell.Context) {
	info(c, "private watch func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	authToken, err := authenticate(c)
	if err != nil {
		die(c, err)
	}

	info(c, "Connecting to buddy private service ...")
	cli, err := client.NewPrivateClient(options.buddy, options.secure)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	stream, err := cli.WatchProfileBuddies(
		auth.ContextWithToken(context.Background(), authToken),
		options.username,
		options.appId,
	)

	if err != nil {
		die(c, err)
	}

	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			info(c, "Watch session terminated...")
			switch ctx.Err() {
			case context.Canceled:
				info(c, "Watch session canceled client side")
			case context.DeadlineExceeded:
				info(c, "Watch session deadline exceeded")
			}
			return
		default:
		}
		info(c, "Waiting for buddies change...\n")

		ros, err := stream.Recv()
		if err != nil {
			die(c, err)
		}

		infof(c, "buddies change data is %s", string(ros.Data))
	}
}

func get(c *ishell.Context) {
	info(c, "get func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	authToken, err := authenticate(c)
	if err != nil {
		die(c, err)
	}

	info(c, "Connecting to buddy service ...")
	cli, err := client.NewPublicClient(options.buddy, options.secure)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	buddies, err := cli.GetBuddies(
		auth.ContextWithToken(context.Background(), authToken),
	)

	if err != nil {
		die(c, err)
	} else {
		infof(c, "buddy data receive success...\nbuddy list is %+v", buddies)
	}

}

func watch(c *ishell.Context) {
	info(c, "watch func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	authToken, err := authenticate(c)
	if err != nil {
		die(c, err)
	}

	info(c, "Connecting to buddy service ...")
	cli, err := client.NewPublicClient(options.buddy, options.secure)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	stream, err := cli.WatchBuddies(
		auth.ContextWithToken(context.Background(), authToken),
	)

	if err != nil {
		die(c, err)
	}

	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			info(c, "Watch session terminated...")
			switch ctx.Err() {
			case context.Canceled:
				info(c, "Watch session canceled client side")
			case context.DeadlineExceeded:
				info(c, "Watch session deadline exceeded")
			}
			return
		default:
		}
		info(c, "Waiting for buddies change...\n")

		ros, err := stream.Recv()
		if err != nil {
			die(c, err)
		}

		infof(c, "buddies change data is %s", string(ros.Data))
	}
}

func reply(c *ishell.Context) {
	info(c, "reply func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	authToken, err := authenticate(c)
	if err != nil {
		die(c, err)
	}

	info(c, "Connecting to buddy service ...")
	cli, err := client.NewPublicClient(options.buddy, options.secure)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter reply buddy name ...")
	buddyName := readLine(c, "BuddyName: ")
	info(c, "Choice reply")

	choice := c.MultiChoice(
		[]string{
			"reject buddy invitation",
			"agree buddy invitation",
		},
		"Select agree or reject invitation",
	)

	var response bool
	switch choice {
	case 0:
		response = false
	case 1:
		response = true
	default:
		die(c, errors.New("Unrecognized response type chosen"))
	}

	err = cli.ReplyAddBuddy(
		auth.ContextWithToken(context.Background(), authToken),
		buddyName,
		response,
	)

	if err != nil {
		die(c, err)
	}
}

func remark(c *ishell.Context) {
	info(c, "remark func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	authToken, err := authenticate(c)
	if err != nil {
		die(c, err)
	}

	info(c, "Connecting to buddy service ...")
	cli, err := client.NewPublicClient(options.buddy, options.secure)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter remark buddy name ...")
	buddyName := readLine(c, "BuddyName: ")
	info(c, "Enter remark text...")
	remark := readLine(c, "Remark: ")

	err = cli.Remark(
		auth.ContextWithToken(context.Background(), authToken),
		buddyName,
		remark,
	)
	if err != nil {
		die(c, err)
	}
}

func updateSettings(c *ishell.Context) {
	info(c, "updateSettings func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	authToken, err := authenticate(c)
	if err != nil {
		die(c, err)
	}

	info(c, "Connecting to buddy service ...")
	cli, err := client.NewPublicClient(options.buddy, options.secure)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter the settings ...")
	AllowToBeAdded := readLine(c, "AllowToBeAdded(true/false): ") == "true"

	err = cli.UpdateBuddySettings(
		auth.ContextWithToken(context.Background(), authToken),
		AllowToBeAdded,
	)
	if err != nil {
		die(c, err)
	}
}

func addBlockedUser(c *ishell.Context) {
	info(c, "addBlockedUser func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	authToken, err := authenticate(c)
	if err != nil {
		die(c, err)
	}

	info(c, "Connecting to buddy service ...")
	cli, err := client.NewPublicClient(options.buddy, options.secure)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter the blocked user id ...")
	id := readLine(c, "BlockedUser:")

	err = cli.AddBlockedProfiles(
		auth.ContextWithToken(context.Background(), authToken),
		id,
	)
	if err != nil {
		die(c, err)
	}
}

func removeBlockedUser(c *ishell.Context) {
	info(c, "removeBlockedUser func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	authToken, err := authenticate(c)
	if err != nil {
		die(c, err)
	}

	info(c, "Connecting to buddy service ...")
	cli, err := client.NewPublicClient(options.buddy, options.secure)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter the blocked user id ...")
	id := readLine(c, "BlockedUser:")

	err = cli.RemoveBlockedProfiles(
		auth.ContextWithToken(context.Background(), authToken),
		id,
	)
	if err != nil {
		die(c, err)
	}
}

func addToRecentMet(c *ishell.Context) {
	info(c, "addToRecentMet func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	authToken, err := authenticate(c)
	if err != nil {
		die(c, err)
	}

	info(c, "Connecting to buddy service ...")
	cli, err := client.NewPublicClient(options.buddy, options.secure)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter the recent profile id ...")
	id := readLine(c, "recentProfileId:")

	err = cli.AddToRecentMet(
		auth.ContextWithToken(context.Background(), authToken),
		id,
	)
	if err != nil {
		die(c, err)
	}
}

func getBlockedList(c *ishell.Context) {
	info(c, "getBlockedList func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	authToken, err := authenticate(c)
	if err != nil {
		die(c, err)
	}

	info(c, "Connecting to buddy service ...")
	cli, err := client.NewPrivateClient(options.buddy, options.secure)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	if resp, err := cli.GetProfileBlockedList(
		auth.ContextWithToken(context.Background(), authToken),
		options.username,
		options.appId,
	); err != nil {
		die(c, err)
	} else {
		infof(c, "blocked list data receive success...\nblocked list is %+v", resp.Profiles)
	}

}
