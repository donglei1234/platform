package tests

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.uber.org/zap"

	auth "github.com/donglei1234/platform/services/auth/pkg/module"
	"github.com/donglei1234/platform/services/chat/internal/app/service/public"
	"github.com/donglei1234/platform/services/chat/pkg/client"
	pb "github.com/donglei1234/platform/services/proto/gen/chat/api"
)

const (
	// these are copied from the non-exported fields in public/chatroom.go
	worldChatType   = "world"
	privateChatType = "private"
	partyChatType   = "party"
	lobbyChatType   = "lobby"

	partyId = "party01"
	lobbyId = "lobby01"
)

type testCase struct {
	username string
	reqText  string
	client   client.ChatStream
}

// chat test suite
func ChatTestSuite(authUrl string, chatUrl string, appId string) {
	testCases := []testCase{
		{
			username: "test01",
			reqText:  "test01 says hello!",
		},
		{
			username: "test02",
			reqText:  "test02 says aloha!",
		},
		{
			username: "test03",
			reqText:  "test03 says bonjour!",
		},
	}

	sharedChatDests := []*pb.Destination{
		chatDest(worldChatType, "", ""),
		chatDest(partyChatType, partyId, ""),
		chatDest(lobbyChatType, lobbyId, ""),
	}

	endCount := len(testCases)
	var endChan = make(chan string, endCount)
	var errChan = make(chan error, 1)

	const stop = "stop"

	var wg sync.WaitGroup
	for i, tc := range testCases {
		wg.Add(1)
		currentUser := tc.username
		log.Println("Executing test case", i+1, "-", currentUser, "...")

		if chatClient, token, err := authenticate(authUrl, chatUrl, tc.username, appId); err != nil {
			errEncountered(fmt.Errorf("unable to authenticate %s %w", currentUser, err))
			return
		} else if s, err := client.CreateChatStream(chatClient, token); err != nil {
			errEncountered(fmt.Errorf("unable to CreateChatStream for %s %w", currentUser, err))
			return
		} else if err := s.Subscribe(testCases[i].username, chatDest(privateChatType, "", "")); err != nil {
			errEncountered(fmt.Errorf("unable to subscribe %s to %s %w", currentUser, privateChatType, err))
			return
		} else if err := func() error {
			// subscribe to all of the common channels
			for _, dest := range sharedChatDests {
				if e := s.Subscribe(testCases[i].username, dest); e != nil {
					return fmt.Errorf("unable to subscribe %s to %q:%q %w",
						currentUser, dest.Channel, dest.Target, err,
					)
				}
			}
			return nil
		}(); err != nil {
			return
		} else {
			testCases[i].client = s
			go func() {
				log.Println(currentUser, "is now chatting ...")
				message := pb.ChatResponse{}
				wg.Done()
				for {
					if err := s.RecvIn(&message); err != nil {
						err := fmt.Errorf("unable to recv message for %s %w", currentUser, err)
						errEncountered(err)
						errChan <- err
						return
					} else if chatErr := message.GetError(); chatErr != nil {
						err := fmt.Errorf("got chat system error for %s %s %w", currentUser, chatErr.Error, public.ErrGeneralFailure)
						errEncountered(err)
						errChan <- err
						return
					} else if message.GetMessage() != nil {
						log.Printf("Chat : %s << %s", currentUser, message.String())
						if message.GetMessage().Message.Text == stop {
							endChan <- currentUser
							return
						}
					}
				}
			}()
		}
	}

	// Allow time for all streams to begin receiving
	wg.Wait()

	for _, tc := range testCases {
		message := func() *pb.ChatMessage_Message {
			return &pb.ChatMessage_Message{
				Name:      tc.username,
				Text:      tc.reqText,
				Timestamp: time.Now().UnixNano(),
				ProfileID: tc.username,
			}
		}

		// message self first
		if err := sendMessage(tc.client, message(), chatDest(privateChatType, tc.username, tc.username)); err != nil {
			errEncountered(err)
			return
		}
		// then message all of the subscribed channels
		for _, dest := range sharedChatDests {
			if err := sendMessage(tc.client, message(), dest); err != nil {
				errEncountered(err)
				return
			}
		}
	}

	// done sending normal messages, wait a tick and send the "stop" message
	time.Sleep(5 * time.Second)

	message := &pb.ChatMessage_Message{
		Name:      testCases[0].username,
		Text:      stop,
		Timestamp: time.Now().UnixNano(),
		ProfileID: testCases[0].username,
	}
	if err := sendMessage(testCases[0].client, message, chatDest(worldChatType, "", "")); err != nil {
		errEncountered(err)
		return
	}

	// and wait for stop messages to be received
	for endCount > 0 {
		select {
		case endUser := <-endChan:
			log.Printf("%s is done ...\n", endUser)
			endCount--
		case <-errChan:
			// bail early because we errored - the message was already logged
			return
		}
	}

	log.Println("- - - - -")
	log.Println("┌────────────────────────────────────┐")
	log.Println("│      Chat Test Result: PASS        │")
	log.Println("└────────────────────────────────────┘")
}

func errEncountered(err error) {
	log.Println("- - - - -")
	log.Println(err.Error())
	log.Println("- - - - -")
	log.Println("┌────────────────────────────────────┐")
	log.Println("│    Chat Test Result: FAIL    	  │")
	log.Println("└────────────────────────────────────┘")
}

func chatDest(channel, target string, targetDisplayName string) *pb.Destination {
	return &pb.Destination{
		Channel:     channel,
		Target:      target,
		DisplayName: targetDisplayName,
	}
}

func sendMessage(client client.ChatStream, msg *pb.ChatMessage_Message, dest *pb.Destination) error {
	return client.SendIn(
		&pb.ChatRequest{
			Kind: &pb.ChatRequest_Message{
				Message: &pb.ChatRequest_Chat{
					Message: &pb.ChatMessage{
						Destination: chatDest(dest.Channel, dest.Target, dest.DisplayName),
						Message:     msg,
					},
				},
			},
		},
	)
}

func authenticate(authUrl string, chatUrl string, username string, appId string) (bClient *client.Client, token string, err error) {
	var authClient auth.PublicClient
	if authClient, err = auth.NewPublicClient(zap.L(), authUrl, false); err != nil {
		return nil, "", err
	} else if token, err = authClient.AuthenticateWithUsername(
		context.Background(),
		appId,
		username,
	); err != nil {
		return nil, "", err
	}
	if bClient, err = client.NewClient(chatUrl, false); err != nil {
		return nil, "", err
	}
	return bClient, token, nil
}
