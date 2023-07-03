package main

import (
	"context"
	"fmt"
	"github.com/abiosoft/ishell"
	pb2 "github.com/donglei1234/platform/services/condition/gen/condition/api"
	"github.com/donglei1234/platform/services/condition/pkg/client"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
	"os"
)

const (
	defaultHost     = "localhost:8081"
	defaultUsername = "test"
)

var options struct {
	host     string
	username string
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "cond_cli",
		Short: "Run a condition CLI",
	}
	rootCmd.PersistentFlags().StringVar(&options.host, "host", defaultHost, "default host")
	rootCmd.PersistentFlags().StringVar(&options.username, "username", defaultUsername, "username for authentication")
	{
		shell := &cobra.Command{
			Use:   "shell",
			Short: "Run an interactive condition service client",
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

	sh.Println("condition Interactive Shell")

	sh.AddCmd(&ishell.Cmd{
		Name:    "watch",
		Help:    "Watch",
		Aliases: []string{"w"},
		Func:    Watch,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:    "Register",
		Help:    "Register",
		Aliases: []string{"r"},
		Func:    Register,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:    "Unregister",
		Help:    "Unregister",
		Aliases: []string{"ur"},
		Func:    Unregister,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:    "Update",
		Help:    "Update",
		Aliases: []string{"u"},
		Func:    Update,
	})
	sh.Run()
}

func Unregister(c *ishell.Context) {
	info(c, "Unregister func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to condition service ...")
	cli, err := client.NewConditionClient(options.host, false)
	if err != nil {
		die(c, err)
	}

	info(c, "Enter token info ...")
	token := readLine(c, "token: ")
	md := metadata.Pairs("token", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	err = cli.Unregister(ctx, &pb2.UnRegisterRequest{
		Conditions: []*pb2.Condition{
			{
				OwnerId: 100,
				Type:    2,
			},
		},
	})
	if err != nil {
		die(c, err)
	}
	done(c, "Unregister success")
}

func Update(c *ishell.Context) {
	info(c, "Update func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to condition service ...")
	cli, err := client.NewConditionClient(options.host, false)
	if err != nil {
		die(c, err)
	}

	info(c, "Enter token info ...")
	token := readLine(c, "token: ")
	md := metadata.Pairs("token", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	err = cli.Update(ctx, &pb2.UpdateRequest{
		Update: []*pb2.Condition{{
			OwnerId:  100,
			Type:     2,
			Params:   []int32{0, 0, 0},
			Progress: 100,
		},
		},
	})
	if err != nil {
		die(c, err)
	}
	done(c, "Update success")
}

func Watch(c *ishell.Context) {
	info(c, "Watch func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to condition service ...")
	cli, err := client.NewConditionClient(options.host, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter token info ...")
	token := readLine(c, "token: ")
	md := metadata.Pairs("token", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	resp, err := cli.Watch(ctx)
	if err != nil {
		die(c, err)
	}
	for {
		res, err := resp.Recv()
		if err != nil {
			die(c, err)
		}
		fmt.Println("res:", res)
	}
}

func Register(c *ishell.Context) {
	info(c, "Register func ...")

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to condition service ...")
	cli, err := client.NewConditionClient(options.host, false)
	if err != nil {
		die(c, err)
	}

	info(c, "Enter token info ...")
	token := readLine(c, "token: ")
	md := metadata.Pairs("token", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	err = cli.Register(ctx, &pb2.RegisterRequest{
		Conditions: []*pb2.Condition{
			{
				OwnerId: 100,
				Type:    2,
				Params:  []int32{1, 2, 3},
			},
		},
	})
	if err != nil {
		die(c, err)
	}
	done(c, "Register success")
}
