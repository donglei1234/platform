package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	pb "github.com/donglei1234/platform/services/proto/gen/auth/api"
	"os"
	"strings"
	"time"

	"github.com/abiosoft/ishell"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	"github.com/donglei1234/platform/services/auth/pkg/client"
)

const (
	DefaultAuth     = "localhost:8081"
	DefaultUsername = "test"
	DefaultAppID    = "tata"

	TestAccessToken = "EAAwNqERruIkBABlSe1ZBPPz6FZB9hbmz1G4dZBmz8t8nAqPA3OvF4sZCCc0s8seZCE4vqcIkgyt5gY4GMhAi6kN3aeLiYaaqMtEWm0QiV7cVPlUuGwQ44ipYNYEoGGk7xM3bjFBdsvmAz8muTZBZBGgbQ0q8dMn3WFO2dEDNZB5Cri7HwCmoVLhuoixFZA3S6y3WhciXLcelK4wZDZD"
	TestIdToken     = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjRiODNmMTgwMjNhODU1NTg3Zjk0MmU3NTEwMjI1MTEyMDg4N2Y3MjUiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiI3NDM5MDY0ODMyNDEtaWRqZmcyNmhkaWFsdXNycmYxMHY3bDNkZGtwZjMyamsuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJhdWQiOiI3NDM5MDY0ODMyNDEtaWRqZmcyNmhkaWFsdXNycmYxMHY3bDNkZGtwZjMyamsuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMDM0MDUyNjkzNjIwODk4ODI2NDMiLCJhdF9oYXNoIjoiaXl4ZVplbDVfOC1qMDlsdGFrODhPUSIsIm5hbWUiOiJ3ZW5oYW8gemhhbyIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vLVJQMjdsOTdfbWNzL0FBQUFBQUFBQUFJL0FBQUFBQUFBQUFBL0FNWnV1Y2xjd0p3NUFYejE0LWVvU3BMX2E2TjNodUZsOXcvczk2LWMvcGhvdG8uanBnIiwiZ2l2ZW5fbmFtZSI6IndlbmhhbyIsImZhbWlseV9uYW1lIjoiemhhbyIsImxvY2FsZSI6InpoLUNOIiwiaWF0IjoxNjAwNzU2NDAyLCJleHAiOjE2MDA3NjAwMDJ9.Z_bQayO4iROSIVbkVV8349yhvy17Dn3M8s1xS5dKryr4Et92c6R4gSiZcSaq-3Jhg9DyyomvicStcxhe1sJuOwYmSuAAwrCI40bhGQsYVVmgsN5CsmtYRTXUds4eLrvXkAWtlqeLsxEOvxCpNCFFuqqohF58NQ28nBbJBHo-sldBObVrq9_4GGfcRgWo1x0lDCDKXaNgixHrJQDNK--E3YPNesqmWiXbRxx98XxmyuHVLVMuFa-8KmEBIoq0oCu7BuosOlbn2J0wfisfGYXUmElkesjC92TPFmWQl-p81hMJmBIqPK-esNkuIpBG2XkLXEvHQYpZDXXcidufD2bhjQ"
)

var options struct {
	auth     string
	username string
}

type tokenPayload struct {
	Exp int    `json:"exp"`
	UID string `json:"uid"`
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "auth_cli",
		Short: "Run a auth CLI",
	}

	rootCmd.PersistentFlags().StringVar(&options.auth, "auth", DefaultAuth, "auth service (<host>:<port>)")
	rootCmd.PersistentFlags().StringVar(&options.username, "username", DefaultUsername, "username for authentication")
	{
		shell := &cobra.Command{
			Use:   "shell",
			Short: "Run an interactive auth service client",
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

	sh.Println("auth Interactive Shell")

	sh.AddCmd(&ishell.Cmd{
		Name: "auth",
		Help: "auth service",
		Func: auth,
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "bindfacebook",
		Help: "bind facebook service",
		Func: bindfacebook,
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "bindgoogle",
		Help: "bind google service",
		Func: bindgoogle,
	})
	sh.AddCmd(&ishell.Cmd{
		Name: "intercept",
		Help: "intercept service",
		Func: intercept,
	})
	sh.Run()
}
func auth(c *ishell.Context) {
	info(c, "add func ...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewAuthClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")
	authName := readLine(c, "Name: ")

	token, e := cli.Auth(
		context.Background(),
		authName,
		DefaultAppID,
	)
	if e != nil {
		die(c, e)
	}
	info(c, "token is "+token.Session.Token)

	_, err = cli.ValidateToken(ctx, token.Session.Token)
	if err != nil {
		die(c, err)
	}
	info(c, "validate is ok !")
}

func bindfacebook(c *ishell.Context) {
	info(c, "add func ...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewAuthClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")
	authName := readLine(c, "Name: ")

	info(c, "1、auth ing...........")
	token, e := cli.Auth(
		context.Background(),
		authName,
		DefaultAppID,
	)
	if e != nil {
		die(c, e)
	}
	info(c, "token is "+token.Session.Token)
	splitToken := strings.Split(token.Session.Token, ".")
	//info(c, "splitToken[1] is "+splitToken[1])
	//decodeBytes, err := base64.StdEncoding.DecodeString(splitToken[1]+"==")
	decodeString, err := base64.RawURLEncoding.DecodeString(splitToken[1])
	if err != nil {
		die(c, err)
	}
	var payload tokenPayload
	err = json.Unmarshal(decodeString, &payload)
	if err != nil {
		die(c, err)
	}
	//info(c,"uid -->"+payload.UID)

	md := metadata.Pairs("token", token.Session.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	info(c, "2、bindFacebook ing...........")
	id, err := cli.Bind(ctx, &pb.BindRequest{
		AppId: DefaultAppID,
		Token: &pb.BindRequest_FacebookToken{
			FacebookToken: TestAccessToken,
		},
	})
	if err != nil {
		die(c, err)
	}
	info(c, "uid :"+id.Uid)
}
func bindgoogle(c *ishell.Context) {
	info(c, "add func ...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewAuthClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")
	authName := readLine(c, "Name: ")

	info(c, "1、auth ing...........")
	token, e := cli.Auth(
		context.Background(),
		authName,
		DefaultAppID,
	)
	if e != nil {
		die(c, e)
	}
	info(c, "token is "+token.Session.Token)
	splitToken := strings.Split(token.Session.Token, ".")
	//info(c, "splitToken[1] is "+splitToken[1])
	//decodeBytes, err := base64.StdEncoding.DecodeString(splitToken[1]+"==")
	decodeString, err := base64.RawURLEncoding.DecodeString(splitToken[1])
	if err != nil {
		die(c, err)
	}
	var payload tokenPayload
	err = json.Unmarshal(decodeString, &payload)
	if err != nil {
		die(c, err)
	}
	//info(c,"uid -->"+payload.UID)

	md := metadata.Pairs("token", token.Session.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	info(c, "2、bindGoogle ing...........")
	id, err := cli.Bind(ctx, &pb.BindRequest{
		AppId: DefaultAppID,
		Token: &pb.BindRequest_GoogleToken{
			GoogleToken: TestIdToken,
		},
	})
	if err != nil {
		die(c, err)
	}
	info(c, "uid :"+id.Uid)
}
func intercept(c *ishell.Context) {
	info(c, "add func ...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)

	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to auth service ...")
	cli, err := client.NewAuthClient(zap.NewExample(), options.auth, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	info(c, "Enter add auth info ...")
	authName := readLine(c, "Name: ")

	token, e := cli.Auth(
		ctx,
		authName,
		DefaultAppID,
	)
	if e != nil {
		die(c, e)
	}
	info(c, "token is "+token.Session.Token)

	md := metadata.Pairs("token", token.Session.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	_, err = cli.ValidateToken(ctx, token.Session.Token)
	if err != nil {
		die(c, err)
	}
	info(c, "validate is ok !")
}
