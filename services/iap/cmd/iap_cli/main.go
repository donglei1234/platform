package main

import (
	"context"
	"os"

	"github.com/abiosoft/ishell"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	pb "github.com/donglei1234/platform/services/iap/generated/grpc/go/iap/api"
	"github.com/donglei1234/platform/services/iap/pkg/client"
)

const (
	DefaultIAP      = "localhost:8081"
	DefaultUsername = "test"
	AppName         = "com.addictive.empire.clash.conquest"
	Sys             = pb.SYS_ANDROID
	ProductToken    = "nfodcagpheigegkegljclfib.AO-J1OyOazDgFDfvJKTHjzci3e_OjwLsmYZFh3HZoWYMGn16kuwvHsw-4DyM25NMVkOurtVpTjE8w4KOhOwB-fI020vWviu_tjBM8D7V0JEKAx7RXzF6KJiAdUuK7C9mVvhtwdQGelBVbm00ZvoixRKEpjVuleuTaA"
	ProductId       = "battle_pass"
	UserToken       = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjRiODNmMTgwMjNhODU1NTg3Zjk0MmU3NTEwMjI1MTEyMDg4N2Y3MjUiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiI3NDM5MDY0ODMyNDEtaWRqZmcyNmhkaWFsdXNycmYxMHY3bDNkZGtwZjMyamsuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJhdWQiOiI3NDM5MDY0ODMyNDEtaWRqZmcyNmhkaWFsdXNycmYxMHY3bDNkZGtwZjMyamsuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMDM0MDUyNjkzNjIwODk4ODI2NDMiLCJhdF9oYXNoIjoiaXl4ZVplbDVfOC1qMDlsdGFrODhPUSIsIm5hbWUiOiJ3ZW5oYW8gemhhbyIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vLVJQMjdsOTdfbWNzL0FBQUFBQUFBQUFJL0FBQUFBQUFBQUFBL0FNWnV1Y2xjd0p3NUFYejE0LWVvU3BMX2E2TjNodUZsOXcvczk2LWMvcGhvdG8uanBnIiwiZ2l2ZW5fbmFtZSI6IndlbmhhbyIsImZhbWlseV9uYW1lIjoiemhhbyIsImxvY2FsZSI6InpoLUNOIiwiaWF0IjoxNjAwNzU2NDAyLCJleHAiOjE2MDA3NjAwMDJ9.Z_bQayO4iROSIVbkVV8349yhvy17Dn3M8s1xS5dKryr4Et92c6R4gSiZcSaq-3Jhg9DyyomvicStcxhe1sJuOwYmSuAAwrCI40bhGQsYVVmgsN5CsmtYRTXUds4eLrvXkAWtlqeLsxEOvxCpNCFFuqqohF58NQ28nBbJBHo-sldBObVrq9_4GGfcRgWo1x0lDCDKXaNgixHrJQDNK--E3YPNesqmWiXbRxx98XxmyuHVLVMuFa-8KmEBIoq0oCu7BuosOlbn2J0wfisfGYXUmElkesjC92TPFmWQl-p81hMJmBIqPK-esNkuIpBG2XkLXEvHQYpZDXXcidufD2bhjQ"
)

var options struct {
	iap      string
	username string
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "iap_cli",
		Short: "Run a IAP CLI",
	}

	rootCmd.PersistentFlags().StringVar(&options.iap, "iap", DefaultIAP, "iap service (<host>:<port>)")
	rootCmd.PersistentFlags().StringVar(&options.username, "username", DefaultUsername, "username for iap")
	{
		shell := &cobra.Command{
			Use:   "shell",
			Short: "Run an interactive IAP service client",
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

	sh.Println("IAP Interactive Shell")

	sh.AddCmd(&ishell.Cmd{
		Name: "checkIAPToken",
		Help: "checkIAPToken service",
		Func: checkIAPToken,
	})
	sh.Run()
}
func checkIAPToken(c *ishell.Context) {
	info(c, "add func ...")
	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)
	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	info(c, "Connecting to iap service ...")
	cli, err := client.NewIAPClient(zap.NewExample(), options.iap, false)
	if err != nil {
		die(c, err)
	}
	defer cli.Close()

	request := &pb.IAPRequest{
		AppStoreId:   AppName,
		Sys:          Sys,
		ProductToken: ProductToken,
		ProductId:    ProductId,
	}
	md := metadata.Pairs("token", UserToken)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	response, err := cli.CheckIAPToken(ctx, request)

	if err != nil {
		info(c, err.Error())
	} else {
		grep := response.Data
		infof(c, "data: %s\n", grep)

	}

	info(c, "Enter add iap info ...")
	//authName := readLine(c, "Name: ")
}
