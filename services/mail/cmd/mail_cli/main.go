package main

import (
	"context"
	"fmt"
	"github.com/abiosoft/ishell"
	pb "github.com/donglei1234/platform/services/proto/gen/mail/api"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
	"os"

	"github.com/donglei1234/platform/services/mail/pkg/client"
)

const (
	defaultMail     = "localhost:8081"
	defaultUsername = "test"
)

var options struct {
	mail     string
	username string
}

func main() {

	rootCmd := &cobra.Command{
		Use:   "cond_cli",
		Short: "Run a mail CLI",
	}

	rootCmd.PersistentFlags().StringVar(&options.mail, "mail", defaultMail, "mail service (<host>:<port>)")
	rootCmd.PersistentFlags().StringVar(&options.username, "username", defaultUsername, "username for authentication")
	{
		shell := &cobra.Command{
			Use:   "shell",
			Short: "Run an interactive mail service client",
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

	sh.Println("Mail Interactive Shell")

	sh.AddCmd(&ishell.Cmd{
		Name:    "watch",
		Help:    "Watch mail",
		Aliases: []string{"w"},
		Func:    WatchMail,
	})
	sh.AddCmd(&ishell.Cmd{
		Name:    "send",
		Help:    "send mail",
		Aliases: []string{"s"},
		Func:    SendMail,
	})
	sh.Run()
}

func WatchMail(c *ishell.Context) {
	info(c, "WatchMail func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to mail service ...")
	cli, err := client.NewMailClient(options.mail, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter token info ...")
	token := readLine(c, "token: ")
	md := metadata.Pairs("token", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	stream, err := cli.Watch(ctx)
	if err != nil {
		die(c, err)
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			info(c, "Recv")
			die(c, err)
		}
		fmt.Println("res:", res)
	}
}

func SendMail(c *ishell.Context) {
	info(c, "SendMail func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to mail service ...")
	cli, err := client.NewMailClient(options.mail, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter token info ...")
	token := readLine(c, "token: ")
	md := metadata.Pairs("token", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err = cli.SendMail(ctx, &pb.Mail{}, nil)
	if err != nil {
		die(c, err)
	}
}
