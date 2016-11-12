package main

import (
	"chatbot"
	"fmt"
	"log"
	"os"
	"strings"
)

// Autoload environment variables in .env
import _ "github.com/joho/godotenv/autoload"

var x = 8

func chatbotProcess(session chatbot.Session, message string) (string, error) {
	if strings.EqualFold(message, "add") {
		x = 1
	}

	if strings.EqualFold(message, "done") {
		words := session["titles"]
		l := len(words)
		wordsForSentence := make([]string, l)
		copy(wordsForSentence, words)
		if l > 1 {
			wordsForSentence[l-1] = "and " + wordsForSentence[l-1]
		}
		if l == 0 {
			return fmt.Sprintf("%s", "If you want to add events, type 'add'!"), nil
		}
		sentence := strings.Join(wordsForSentence, ", ")
		x = 8
		return fmt.Sprintf("You added %s! If you want to add other events, type 'add'! ", strings.ToLower(sentence)), nil
	}

	switch x {
	case 1:
		x = 2
		return fmt.Sprintf("%s", "Please enter the title of the event"), nil
	case 2:
		session["titles"] = append(session["titles"], message)
		x = 3
		return fmt.Sprintf("%s", "Please enter the date of the event"), nil
	case 3:
		session["dates"] = append(session["date"], message)
		x = 4
		return fmt.Sprintf("%s", "Please enter the timing of the event"), nil
	case 4:
		session["timing"] = append(session["timings"], message)
		x = 5
		return fmt.Sprintf("%s", "Please enter the longitude of the event"), nil
	case 5:
		session["longitudes"] = append(session["longitudes"], message)
		x = 6
		return fmt.Sprintf("%s", "Please enter the latitude of the event"), nil
	case 6:
		session["latitudes"] = append(session["latitudes"], message)
		return fmt.Sprintf("%s", "This event is done! Either type 'add' or 'done'!"), nil

	default:
		return "", fmt.Errorf("%s", "Invalid text!")

	}

	// 	if strings.EqualFold(message, "chatbot") {
	// 		return "", fmt.Errorf("This can't be, I'm the one and only %s!", message)
	// 	}

	// 	return fmt.Sprintf("Hello %s, my name is chatbot. What was yours again?", message), nil
}

func main() {
	// Uncomment the following lines to customize the chatbot
	//chatbot.WelcomeMessage = "What's your name?"
	chatbot.WelcomeMessage = "Tick-tock, Whenever you want to add an event, just type 'add'!"
	chatbot.ProcessFunc(chatbotProcess)

	// Use the PORT environment variable
	port := os.Getenv("PORT")
	// Default to 3000 if no PORT environment variable was defined
	if port == "" {
		port = "3000"
	}

	// Start the server
	fmt.Printf("Listening on port %s...\n", port)
	log.Fatalln(chatbot.Engage(":" + port))
}
