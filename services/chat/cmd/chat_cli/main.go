package main

import (
	"context"
	"io"
	"time"

	"github.com/donglei1234/platform/services/chat/pkg/tests"

	"go.uber.org/zap"

	"github.com/abiosoft/ishell"
	"github.com/donglei1234/platform/services/auth/pkg/auth"
	"github.com/donglei1234/platform/services/chat/pkg/client"
	ish "github.com/donglei1234/platform/services/chat/pkg/ishellex"
	"github.com/spf13/cobra"
)

const (
	defaultChat      = "chat.sr-development.svc.cluster.local:8081"
	defaultAuth      = "auth.sr-development.svc.cluster.local:8081"
	defaultAuthAppId = "test"
	defaultUsername  = "test"
	defaultSecure    = false
)

var options struct {
	chat     string
	auth     string
	appId    string
	username string
	secure   bool
}

func main() {

	rootCmd := &cobra.Command{
		Use:   "chat_cli",
		Short: "Run an interactive chat",
	}

	{
		shell := &cobra.Command{
			Use:   "shell",
			Short: "Run an interactive chat client",
			Run: func(cmd *cobra.Command, args []string) {
				shell()
			},
		}
		rootCmd.AddCommand(shell)
	}

	{
		chatTestSuite := &cobra.Command{
			Use:   "testSuite",
			Short: "Run the chat test suite",
			Run: func(cmd *cobra.Command, args []string) {
				tests.ChatTestSuite(options.auth, options.chat, options.appId)
			},
		}
		rootCmd.AddCommand(chatTestSuite)
	}

	rootCmd.PersistentFlags().StringVar(
		&options.auth,
		"auth",
		defaultAuth,
		"authentication service (<host>:<port>)",
	)
	rootCmd.PersistentFlags().StringVar(
		&options.chat,
		"chat",
		defaultChat,
		"chat service (<host>:<port>)",
	)
	rootCmd.PersistentFlags().StringVar(
		&options.appId,
		"appId",
		defaultAuthAppId,
		"appId for authentication",
	)
	rootCmd.PersistentFlags().StringVar(
		&options.username,
		"username",
		defaultUsername,
		"username for authentication",
	)
	rootCmd.PersistentFlags().BoolVar(
		&options.secure,
		"secure",
		defaultSecure,
		"if provided, connect securely",
	)

	rootCmd.Execute()
}

func shell() {
	sh := ishell.New()

	sh.Println("Chat Interactive Shell")

	sh.AddCmd(&ishell.Cmd{
		Name: "private",
		Help: "run a private client",
		Func: private,
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "world",
		Help: "run a world client",
		Func: world,
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "team",
		Help: "run a team client",
		Func: team,
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "room",
		Help: "run a room client",
		Func: room,
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "party",
		Help: "run a party client",
		Func: party,
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "lobby",
		Help: "run a lobby client",
		Func: lobby,
	})
	sh.Run()
}

// private chat
func private(c *ishell.Context) {
	ish.Info(c, "private func ...")

	cli, stream := authAndSubscribe(c)
	if cli != nil {
		defer cli.Close()
	}

	ish.Info(c, "Enter destination name ...")
	destname := ish.ReadLine(c, "name: ")

	ish.Info(c, "Enter what you want to say: ")

	// goroutine and channel for chat from the server
	recvChan := make(chan pb.ChatResponse, 100)
	recv(recvChan, stream)

	// goroutine and channel for chat we're sending to the server
	sendChan := make(chan pb.ChatRequest, 100)
	send("private", options.username, options.username, destname, destname, sendChan, stream, c)

	// chat forever
	for {
		loop(sendChan, recvChan, stream, c)
	}
}

// world chat
func world(c *ishell.Context) {
	ish.Info(c, "world func ...")

	cli, stream := authAndSubscribe(c)
	if cli != nil {
		defer cli.Close()
	}

	// goroutine and channel for chat from the server
	recvChan := make(chan pb.ChatResponse, 100)
	recv(recvChan, stream)

	// goroutine and channel for chat we're sending to the server
	sendChan := make(chan pb.ChatRequest, 100)
	send("world", options.username, options.username, "", "", sendChan, stream, c)

	// chat forever
	for {
		loop(sendChan, recvChan, stream, c)
	}
}

// team chat
func team(c *ishell.Context) {
	ish.Info(c, "team func ...")

	cli, stream := authAndSubscribe(c)
	if cli != nil {
		defer cli.Close()
	}

	// goroutine and channel for chat from the server
	recvChan := make(chan pb.ChatResponse, 100)
	recv(recvChan, stream)

	// goroutine and channel for chat we're sending to the server
	sendChan := make(chan pb.ChatRequest, 100)
	send("team", options.username, options.username, "01", "", sendChan, stream, c)

	// chat forever
	for {
		loop(sendChan, recvChan, stream, c)
	}
}

// room chat
func room(c *ishell.Context) {
	ish.Info(c, "room func ...")

	cli, stream := authAndSubscribe(c)
	if cli != nil {
		defer cli.Close()
	}

	// goroutine and channel for chat from the server
	recvChan := make(chan pb.ChatResponse, 100)
	recv(recvChan, stream)

	// goroutine and channel for chat we're sending to the server
	sendChan := make(chan pb.ChatRequest, 100)
	send("room", options.username, options.username, "01", "", sendChan, stream, c)

	// chat forever
	for {
		loop(sendChan, recvChan, stream, c)
	}
}

// party chat
func party(c *ishell.Context) {
	ish.Info(c, "party func ...")

	cli, stream := authAndSubscribe(c)
	if cli != nil {
		defer cli.Close()
	}

	// goroutine and channel for chat from the server
	recvChan := make(chan pb.ChatResponse, 100)
	recv(recvChan, stream)

	// goroutine and channel for chat we're sending to the server
	sendChan := make(chan pb.ChatRequest, 100)
	send("party", options.username, options.username, "01", "", sendChan, stream, c)

	// chat forever
	for {
		loop(sendChan, recvChan, stream, c)
	}
}

// lobby chat
func lobby(c *ishell.Context) {
	ish.Info(c, "lobby func ...")

	cli, stream := authAndSubscribe(c)
	if cli != nil {
		defer cli.Close()
	}

	// goroutine and channel for chat from the server
	recvChan := make(chan pb.ChatResponse, 100)
	recv(recvChan, stream)

	// goroutine and channel for chat we're sending to the server
	sendChan := make(chan pb.ChatRequest, 100)
	send("lobby", options.username, options.username, "01", "", sendChan, stream, c)

	// chat forever
	for {
		loop(sendChan, recvChan, stream, c)
	}
}
func authAndSubscribe(c *ishell.Context) (cli *client.Client, stream client.ChatStream) {
	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	// authenticate
	authToken, err := authenticate(c)
	if err != nil {
		ish.Die(c, err)
	}

	// connect to the chat service
	ish.Info(c, "Connecting to chat service ...")
	chatClient, err := client.NewClient(options.chat, options.secure)
	if err != nil {
		ish.Die(c, err)
	}
	//defer cli.Close()

	// create a new chat stream
	ish.Info(c, "Creating a new chat stream ...")
	s, err := client.CreateChatStream(chatClient, authToken)
	if err != nil {
		ish.Die(c, err)
	}
	// subscribe channels
	ish.Info(c, "Subscribe channels with "+options.username)
	s.Subscribe(options.username, &pb.Destination{Channel: "world", Target: "", DisplayName: ""})
	s.Subscribe(options.username, &pb.Destination{Channel: "private", Target: "", DisplayName: ""})
	s.Subscribe(options.username, &pb.Destination{Channel: "team", Target: "01", DisplayName: ""})
	s.Subscribe(options.username, &pb.Destination{Channel: "room", Target: "01", DisplayName: ""})
	s.Subscribe(options.username, &pb.Destination{Channel: "party", Target: "01", DisplayName: ""})
	s.Subscribe(options.username, &pb.Destination{Channel: "lobby", Target: "01", DisplayName: ""})

	cli = chatClient
	stream = s
	return
}

func authenticate(c *ishell.Context) (token string, err error) {
	ish.Info(c, "Authenticating...")
	if authClient, e := auth.NewPublicClient(zap.L(), options.auth, options.secure); e != nil {
		err = e
	} else if t, e := authClient.AuthenticateWithUsername(
		context.Background(),
		options.appId,
		options.username,
	); e != nil {
		err = e
	} else {
		ish.Infof(c, "Authentication token: %s\n", t)
		token = t
	}

	return
}

func recv(c chan pb.ChatResponse, s client.ChatStream) {
	// goroutine and channel for chat from the server
	go func() {
		in := pb.ChatResponse{}
		for {
			err := s.RecvIn(&in)
			if err == io.EOF {
				return
			}
			if err == nil {
				c <- in
			}
		}
	}()
}

func send(channel string, username string, profileID string, dest string, destDisplayName string, ch chan pb.ChatRequest, s client.ChatStream, c *ishell.Context) {
	go func() {
		for {
			text := ish.ReadLine(c, "say: ")
			if text == "quit" {
				ish.Done(c, "quit")
				return
			} else if text == "unsubscribe" {
				sendUnSubscribeChat(channel, username, dest, destDisplayName, text, ch, c)
			} else {
				sendChat(channel, username, profileID, dest, destDisplayName, text, ch, c)
			}

		}
	}()
}

func sendChat(channel string, username string, profileID string, dest string, destDisplayName string, text string, ch chan pb.ChatRequest, c *ishell.Context) {
	destination := &pb.Destination{Channel: channel, Target: dest, DisplayName: destDisplayName}
	message := &pb.ChatMessage_Message{Name: username, Text: text, Timestamp: time.Now().UnixNano(), ProfileID: profileID}
	out := pb.ChatRequest{
		Kind: &pb.ChatRequest_Message{
			Message: &pb.ChatRequest_Chat{
				Message: &pb.ChatMessage{
					Destination: destination,
					Message:     message,
				},
			},
		},
	}
	ch <- out
}

func sendUnSubscribeChat(channel string, username string, dest string, destDisplayName string, text string, ch chan pb.ChatRequest, c *ishell.Context) {
	destination := &pb.Destination{Channel: channel, Target: dest, DisplayName: destDisplayName}
	out := pb.ChatRequest{
		Kind: &pb.ChatRequest_Unsubscribe{
			Unsubscribe: &pb.ChatRequest_UnSubscribe{
				Name:        username,
				Destination: destination,
			},
		},
	}
	ch <- out
}

func loop(send chan pb.ChatRequest, receive chan pb.ChatResponse, s client.ChatStream, c *ishell.Context) {
	select {
	case message := <-send:
		//ish.Info(c, "send ...")
		err := s.SendIn(&message)
		if err != nil {
			ish.Info(c, err.Error())
		}
	case received := <-receive:
		//ish.Info(c, "receive ...")
		switch received.Kind.(type) {
		case *pb.ChatResponse_Error:
			ish.Info(c, "\n")
			ish.Info(c, time.Unix(0, received.GetError().Timestamp).String())
			out := received.GetError().GetDestination().Channel + ": Error," + received.GetError().Error + "Message: " + received.GetError().String()
			ish.Info(c, out)
		case *pb.ChatResponse_Message:
			handleMessage(c, received.GetMessage())
		default:
		}
	default:

	}
}

func handleMessage(c *ishell.Context, received *pb.ChatMessage) {
	switch received.Destination.Channel {
	case "private":
		ish.Info(c, "\n")
		ish.Info(c, time.Unix(0, received.GetMessage().Timestamp).String())
		out := "private: @" + received.GetMessage().Name + ": " + received.GetMessage().Text
		ish.Info(c, out)
		//ish.Info(c, ">")

	case "world":
		ish.Info(c, "\n")
		ish.Info(c, time.Unix(0, received.Message.Timestamp).String())
		out := "world: @" + received.Message.Name + ": " + received.Message.Text
		ish.Info(c, out)
		//ish.Info(c, ">")
	case "team":
		ish.Info(c, "\n")
		ish.Info(c, time.Unix(0, received.Message.Timestamp).String())
		out := "team: @" + received.Message.Name + ": " + received.Message.Text
		ish.Info(c, out)
		//ish.Info(c, ">")
	case "room":
		ish.Info(c, "\n")
		ish.Info(c, time.Unix(0, received.Message.Timestamp).String())
		out := "room: @" + received.Message.Name + ": " + received.Message.Text
		ish.Info(c, out)
		//ish.Info(c, ">")
	case "party":
		ish.Info(c, "\n")
		ish.Info(c, time.Unix(0, received.Message.Timestamp).String())
		out := "party: @" + received.Message.Name + ": " + received.Message.Text
		ish.Info(c, out)
		//ish.Info(c, ">")
	case "lobby":
		ish.Info(c, "\n")
		ish.Info(c, time.Unix(0, received.Message.Timestamp).String())
		out := "lobby: @" + received.Message.Name + ": " + received.Message.Text
		ish.Info(c, out)
		//ish.Info(c, ">")
	}
}
