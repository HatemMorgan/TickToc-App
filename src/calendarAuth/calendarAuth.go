package calendarAuth

import (
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
// func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
// 	cacheFile, err := tokenCacheFile()
// 	if err != nil {
// 		log.Fatalf("Unable to get path to cached credential file. %v", err)
// 	}
// 	tok, err := tokenFromFile(cacheFile)
// 	if err != nil {
// 		tok = getTokenFromWeb(config)
// 		// saveToken(cacheFile, tok)
// 	}
// 	return config.Client(ctx, tok)
// }

//GetAuthURLFromWeb return authentication url
func GetAuthURLFromWeb() (string, error) {
	config, err := getConfig()
	if err != nil {
		return "", err
	}
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	return authURL, nil
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
// func tokenCacheFile() (string, error) {
// 	usr, err := user.Current()
// 	if err != nil {
// 		return "", err
// 	}
// 	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
// 	os.MkdirAll(tokenCacheDir, 0700)
// 	return filepath.Join(tokenCacheDir,
// 		url.QueryEscape("calendar-go-quickstart.json")), err
// }

// // tokenFromFile retrieves a Token from a given file path.
// // It returns the retrieved Token and any read error encountered.
// func tokenFromFile(file string) (*oauth2.Token, error) {
// 	f, err := os.Open(file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	t := &oauth2.Token{}
// 	err = json.NewDecoder(f).Decode(t)
// 	defer f.Close()
// 	return t, err
// }

//GetTokenToBeSaved used to get token for authentication of user pass it insertUser to insert user
func GetTokenToBeSaved(tokenCode string) (*oauth2.Token, error) {
	// fmt.Printf("Saving credential file to: %s\n", file)
	// f, err := os.Create(file)
	// if err != nil {
	// 	log.Fatalf("Unable to cache oauth token: %v", err)
	// }
	// defer f.Close()
	// json.NewEncoder(f).Encode(token)
	config, err := getConfig()

	if err != nil {
		return &oauth2.Token{}, err
	}

	tok, err := config.Exchange(oauth2.NoContext, tokenCode)

	if err != nil {
		return &oauth2.Token{}, fmt.Errorf("Unable to retrieve token from web %v", err)
	}

	return tok, nil

}

func getConfig() (*oauth2.Config, error) {

	b, err := ioutil.ReadFile("./calendarAuth/client_secret.json")
	if err != nil {
		fmt.Println("Unable to read client secret file: ", err)
		return &oauth2.Config{}, fmt.Errorf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	// config, err := google.ConfigFromJSON(b)
	if err != nil {
		fmt.Println("Unable to parse client secret file to config: ", err)
		return &oauth2.Config{}, fmt.Errorf("Unable to parse client secret file to config: %v", err)
	}
	return config, nil
}

// GetCalendarService gives service in order to call calendar apis
func GetCalendarService(token *oauth2.Token) (calendar.Service, error) {
	ctx := context.Background()
	// this should be changed to get path dynamically
	//	b, err := ioutil.ReadFile("/home/hatem/workspaceGO/src/calendarAuth/client_secret.json")
	b, err := ioutil.ReadFile("./calendarAuth/client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/calendar-go-quickstart.json
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	// config, err := google.ConfigFromJSON(b)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := config.Client(ctx, token)

	srv, err := calendar.New(client)

	return *srv, err

}

// func main() {
// 	ctx := context.Background()

// 	b, err := ioutil.ReadFile("client_secret.json")
// 	if err != nil {
// 		log.Fatalf("Unable to read client secret file: %v", err)
// 	}

// 	// If modifying these scopes, delete your previously saved credentials
// 	// at ~/.credentials/calendar-go-quickstart.json
// 	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
// 	if err != nil {
// 		log.Fatalf("Unable to parse client secret file to config: %v", err)
// 	}
// 	client := getClient(ctx, config)

// 	srv, err := calendar.New(client)
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve calendar Client %v", err)
// 	}

// 	t := time.Now().Format(time.RFC3339)
// 	fmt.Println(t + "")
// 	events, err := srv.Events.List("primary").ShowDeleted(false).
// 		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve next ten of the user's events. %v", err)
// 	}

// 	fmt.Println("Upcoming events:")
// 	if len(events.Items) > 0 {
// 		for _, i := range events.Items {
// 			var when string
// 			// If the DateTime is an empty string the Event is an all-day Event.
// 			// So only Date is available.
// 			if i.Start.DateTime != "" {
// 				when = i.Start.DateTime
// 			} else {
// 				when = i.Start.Date
// 			}
// 			fmt.Printf("%s (%s)\n", i.Summary, when)
// 		}
// 	} else {
// 		fmt.Printf("No upcoming events found.\n")
// 	}

// 	listRes, err := srv.CalendarList.List().Fields("items/id").Do()
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve list of calendars: %v", err)
// 	}
// 	for _, v := range listRes.Items {
// 		log.Printf("Calendar ID: %v and description: %d \n", v.Id)
// 	}

// 	log.Println("------------------------------------------ ")

// 	if len(listRes.Items) > 0 {
// 		id := listRes.Items[2].Id
// 		res, err := srv.Events.List(id).Fields("items(updated,summary)", "summary", "nextPageToken").Do()
// 		if err != nil {
// 			log.Fatalf("Unable to retrieve calendar events list: %v", err)
// 		}
// 		for _, v := range res.Items {
// 			log.Printf("Calendar ID %q event: %v: %q\n", id, v.Updated, v.Summary)
// 		}
// 		log.Printf("Calendar ID %q Summary: %v\n", id, res.Summary)
// 		log.Printf("Calendar ID %q next page token: %v\n", id, res.NextPageToken)
// 	}

// }
