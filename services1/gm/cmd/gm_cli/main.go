package main

import (
	"github.com/abiosoft/ishell"
	"github.com/spf13/cobra"
	"os"
)

const (
	defaultMail     = "localhost:8081"
	defaultUsername = "test"
)

var options struct {
	gm       string
	username string
}

func main() {

	rootCmd := &cobra.Command{
		Use:   "cond_cli",
		Short: "Run a gm CLI",
	}

	rootCmd.PersistentFlags().StringVar(&options.gm, "gm", defaultMail, "gm service (<host>:<port>)")
	rootCmd.PersistentFlags().StringVar(&options.username, "username", defaultUsername, "username for authentication")
	{
		shell := &cobra.Command{
			Use:   "shell",
			Short: "Run an interactive gm service client",
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

	sh.Run()
}
