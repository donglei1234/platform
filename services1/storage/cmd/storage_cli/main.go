package main

import (
	"context"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"

	"github.com/abiosoft/ishell"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/storage/pkg/client"
)

const (
	DefaultStorage = "localhost:8081"
)

var options struct {
	storage string
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "storage_cli",
		Short: "Run a storage CLI",
	}

	rootCmd.PersistentFlags().StringVar(&options.storage, "storage", DefaultStorage, "storage service (<host>:<port>)")
	{
		shell := &cobra.Command{
			Use:   "shell",
			Short: "Run an interactive storage service client",
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

	sh.Println("storage Interactive Shell")

	sh.AddCmd(&ishell.Cmd{
		Name: "1",
		Help: "GetObjects service",
		Func: GetObjects,
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "2",
		Help: "Upload service",
		Func: Upload,
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "3",
		Help: "Delete service",
		Func: Delete,
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "4",
		Help: "DownloadForUrl service",
		Func: DownloadForUrl,
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "5",
		Help: "DownloadForItem service",
		Func: DownloadForItem,
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "6",
		Help: "GetObjectACL service",
		Func: GetObjectACL,
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "7",
		Help: "SetObjectACL service",
		Func: SetObjectACL,
	})

	sh.Run()
}

func GetObjects(c *ishell.Context) {
	info(c, "start GetObjects func ...")
	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting service ...")
	cli, err := client.NewStorageClient(zap.NewExample(), options.storage, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()
	profileId := readLine(c, "profileId: ")
	appId := readLine(c, "appId: ")
	_, e := cli.GetFiles(
		context.Background(), profileId, appId,
	)

	if e != nil {
		die(c, err)
	}
	info(c, "end GetObjects func ...")
}

func Upload(c *ishell.Context) {
	info(c, "start Upload func ...")
	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting service ...")
	cli, err := client.NewStorageClient(zap.NewExample(), options.storage, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	path := readLine(c, "path: ")
	key := readLine(c, "key: ")
	appId := readLine(c, "appId: ")
	r, _ := regexp.Compile("[^0-9a-zA-Z_]")
	key = r.ReplaceAllString(key, "_")
	profileId := readLine(c, "profileId: ")
	file, err := os.Open(path)
	all, err := ioutil.ReadAll(file)
	if err != nil {
		die(c, err)
	}
	defer file.Close()
	cli.UploadFile(
		context.Background(),
		all, key, profileId, appId,
	)
	info(c, "end Upload func ...")
}

func Delete(c *ishell.Context) {
	info(c, "start Delete func ...")
	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting service ...")
	cli, err := client.NewStorageClient(zap.NewExample(), options.storage, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	//path := readLine(c, "path: ")
	key := readLine(c, "key: ")
	profileId := readLine(c, "profileId: ")
	appId := readLine(c, "appId: ")
	cli.Delete(
		context.Background(),
		key, profileId, appId,
	)
	info(c, "end Delete func ...")
}

func DownloadForUrl(c *ishell.Context) {
	info(c, "start DownloadForUrl func ...")
	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting service ...")
	cli, err := client.NewStorageClient(zap.NewExample(), options.storage, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	//path := readLine(c, "path: ")
	key := readLine(c, "key: ")
	profileId := readLine(c, "profileId: ")
	appId := readLine(c, "appId: ")
	url, err := cli.DownloadForUrl(
		context.Background(),
		key, profileId,
		appId,
	)
	info(c, "url:"+url.Url)
	info(c, "end DownloadForUrl func ...")
}

func DownloadForItem(c *ishell.Context) {
	info(c, "start DownloadForItem func ...")
	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting service ...")
	cli, err := client.NewStorageClient(zap.NewExample(), options.storage, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	appId := readLine(c, "appId: ")
	key := readLine(c, "key: ")
	profileId := readLine(c, "profileId: ")
	cli.DownloadForItem(
		context.Background(),
		key, profileId, appId,
	)
	info(c, "end DownloadForItem func ...")
}

func GetObjectACL(c *ishell.Context) {
	info(c, "start GetObjectACL func ...")
	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting service ...")
	cli, err := client.NewStorageClient(zap.NewExample(), options.storage, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	key := readLine(c, "key: ")
	profileId := readLine(c, "profileId: ")
	appId := readLine(c, "appId: ")
	acl, _ := cli.GetFileACL(
		context.Background(),
		key, profileId, appId,
	)

	info(c, "  Grantee:   "+acl.Grantee)
	info(c, "  Type:      "+acl.Type)
	info(c, "  Permission:"+acl.Permission)
	info(c, "end GetObjectACL func ...")
}

func SetObjectACL(c *ishell.Context) {
	info(c, "start SetObjectACL func ...")
	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting service ...")
	cli, err := client.NewStorageClient(zap.NewExample(), options.storage, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	aclType := readLine(c, "aclType: ")
	atoi, _ := strconv.Atoi(aclType)
	key := readLine(c, "key: ")
	profileId := readLine(c, "profileId: ")
	appId := readLine(c, "appId: ")

	cli.SetFileACL(
		context.Background(),
		key, profileId, int32(atoi), appId,
	)
	info(c, "end SetObjectACL func ...")
}
