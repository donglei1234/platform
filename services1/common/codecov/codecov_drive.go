package codecov

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

const encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOnline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func EncodeCredentials(credentialPath string) (encodedCredentials string, err error) {
	//if credentials, err := file.ReadFile(credentialPath); err != nil {
	//	return "", err
	//} else {
	//	encodedCredentials = base64.StdEncoding.EncodeToString([]byte(credentials))
	//}
	return encodedCredentials, nil
}

func UpdateDrive(sheetID string, csvFile string, credentials string) (err error) {
	var decodedCreds []byte
	if decodedCreds, err = base64.StdEncoding.DecodeString(credentials); err != nil {
		return errors.Errorf("Could not decode Google Sheets secret: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(decodedCreds, sheets.DriveFileScope)
	if err != nil {
		return errors.Errorf("Unable to parse client secret to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		return errors.Errorf("Unable to retrieve Drive client: %v", err)
	}

	// Upload CSV and convert to Spreadsheet
	filename := csvFile // File you want to upload

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer file.Close()

	var sheet *sheets.Spreadsheet
	if sheetID == "" {
		createSpreadsheet := srv.Spreadsheets.Create(&sheets.Spreadsheet{})
		if sheet, err = createSpreadsheet.Do(); err != nil {
			log.Fatalf("Error: %v", err)
		}
	} else {
		getSpreadsheet := srv.Spreadsheets.Get(sheetID)
		if sheet, err = getSpreadsheet.Do(); err != nil {
			log.Fatalf("Error: %v", err)
		}
		writeRange := "A:B"
		reader := csv.NewReader(file)
		record, err := reader.ReadAll()
		if err != nil {
			fmt.Println("Error", err)
		}

		all := [][]interface{}{}

		for _, value := range record {
			all = append(all, []interface{}{value[0], value[1]})
		}

		rb := &sheets.BatchUpdateValuesRequest{
			ValueInputOption: "USER_ENTERED",
		}
		rb.Data = append(rb.Data, &sheets.ValueRange{
			Range:  writeRange,
			Values: all,
		})

		if _, err := srv.Spreadsheets.Values.BatchUpdate(
			sheet.SpreadsheetId,
			rb,
		).Do(); err != nil {
			log.Fatalf("Error: %v", err)
		} else {
			fmt.Println("Spreadsheet is available to view at", sheet.SpreadsheetUrl)
		}
	}
	return nil
}
