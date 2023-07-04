package main

import (
	"context"
	"os"
	"strconv"

	"github.com/abiosoft/ishell"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/notification/pkg/client"
)

const (
	DefaultNotification = "localhost:8081"
	DefaultUsername     = "test"
	DefaultAppID        = "aow"
	DefaultRegion       = "cn"
	DefaultDeviceType   = "googleplay"
)

var options struct {
	auth     string
	username string
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "auth_cli",
		Short: "Run a auth CLI",
	}

	rootCmd.PersistentFlags().StringVar(&options.auth, "notification", DefaultNotification, "notification service (<host>:<port>)")
	rootCmd.PersistentFlags().StringVar(&options.username, "username", DefaultUsername, "username for authentication")
	{
		shell := &cobra.Command{
			Use:   "shell",
			Short: "Run an interactive notification service client",
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

	sh.Println("notification Interactive Shell")

	sh.AddCmd(&ishell.Cmd{
		Name: "registerArn",
		Help: "registerArn service",
		Func: registerArn,
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "deleteArn",
		Help: "deleteArn service",
		Func: deleteArn,
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "publishMessage",
		Help: "publishMessage service",
		Func: publishMessage,
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "subscribeTopic",
		Help: "subscribeTopic service",
		Func: subscribeTopic,
	})
	sh.Run()
}

func registerArn(c *ishell.Context) {
	info(c, "add func ...")
	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewNotificationClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")
	//authName := readLine(c, "Name: ")

	profileId := readLine(c, "profileId: ")
	deviceToken := readLine(c, "deviceToken: ")
	deviceId := readLine(c, "deviceId: ")
	deviceType := readLine(c, "deviceType: ")
	_, e := cli.RegisterArn(
		context.Background(),
		profileId,
		deviceToken,
		deviceId,
		DefaultRegion,
		DefaultAppID,
		deviceType,
	)
	if e != nil {
		die(c, e)
	}
}

func deleteArn(c *ishell.Context) {
	info(c, "add func ...")
	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewNotificationClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")
	//authName := readLine(c, "Name: ")
	tempDelType, err := strconv.Atoi(readLine(c, "delType: "))
	delType := int32(tempDelType)
	publishId := readLine(c, "publishId: ")
	deviceId := readLine(c, "deviceId: ")
	_, e := cli.DeleteArn(
		context.Background(),
		delType,
		publishId,
		deviceId,
		DefaultAppID,
	)
	if e != nil {
		die(c, e)
	}
}

func publishMessage(c *ishell.Context) {
	info(c, "add func ...")
	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewNotificationClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")
	//authName := readLine(c, "Name: ")
	tempPubType, err := strconv.Atoi(readLine(c, "pubType: "))
	pubType := int32(tempPubType)
	publishId := readLine(c, "publishId: ")
	message := readLine(c, "message: ")
	msg := []byte(message)
	_, e := cli.PublishMessage(
		context.Background(),
		pubType,
		msg,
		publishId,
		DefaultAppID,
	)
	if e != nil {
		die(c, e)
	}

}

func subscribeTopic(c *ishell.Context) {
	info(c, "add func ...")
	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewNotificationClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")
	//authName := readLine(c, "Name: ")

	tempSubType, err := strconv.Atoi(readLine(c, "subType: "))
	subType := int32(tempSubType)
	topicName := readLine(c, "topicName: ")
	profileId := readLine(c, "profileId: ")
	_, e := cli.SubscribeTopic(
		context.Background(),
		subType,
		topicName,
		profileId,
		DefaultAppID,
	)
	if e != nil {
		die(c, e)
	}
}
