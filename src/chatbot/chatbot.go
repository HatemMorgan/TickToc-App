package chatbot

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	// WelcomeMessage A constant to hold the welcome message
	WelcomeMessage = "Welcome, what do you want to order?"

	// sessions = {
	//   "uuid1" = Session{
	//     "history" = [
	//       "Message 1",
	//       "Message 2",
	//       "Message 3",
	//       ...
	//     ]
	//   },
	//   ...
	// }
	sessions = map[string]Session{}

	processor = defaultProcessor
)

type (
	// Session Holds info about a session
	Session map[string]string

	// JSON Holds a JSON object
	// JSON map[string]interface{}

	// Processor Alias for Process func
	Processor func(session Session, message string) (string, error)
)

func defaultProcessor(session Session, message string) (string, error) {
	// // Make sure the message is unique in history
	// for _, m := range session["history"] {
	// 	if strings.EqualFold(m, message) {
	// 		return "", fmt.Errorf("You've already ordered %s before!", message)
	// 	}
	// }

	// // Add the message in the parsed body to the messages in the session
	// session["history"] = append(session["history"], message)

	// // Form a sentence out of the history in the form Message 1, Message 2, and Message 3
	// words := session["history"]
	// lenght := len(words)
	// wordsForSentence := make([]string, lenght)
	// copy(wordsForSentence, words)
	// if lenght > 1 {
	// 	wordsForSentence[lenght-1] = "and " + wordsForSentence[lenght-1]
	// }
	// sentence := strings.Join(wordsForSentence, ", ")

	return fmt.Sprintf("defaultProcessor is running !", strings.ToLower(" ")), nil
}

// ProcessFunc Sets the processor of the chatbot
func ProcessFunc(p Processor) {
	processor = p
}

//Welcome is called by welcome handler route to generate an UUID
func Welcome() map[string]string {
	// Generate a UUID.
	// bygeeb time stamp unix format and represent it in base 10
	hasher := md5.New()
	hasher.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10)))
	uuid := hex.EncodeToString(hasher.Sum(nil))

	// Create a session for this UUID
	sessions[uuid] = Session{}

	// Write a JSON containg the welcome message and the generated UUID
	return map[string]string{
		"uuid":    uuid,
		"message": WelcomeMessage,
	}
}

//CheckIfAuthenticated checks if the user has a session opened and his uuid is valid
func CheckIfAuthenticated(uuid string) bool {
	// Make sure a session exists for the extracted UUID
	_, sessionFound := sessions[uuid]
	return sessionFound
}

//Chat is called by chat route handler to save message and return the new question
func Chat(uuid string, data map[string]string) (map[string]string, error) {

	// gets user session
	session, _ := sessions[uuid]

	// Process the received message
	message, err := processor(session, data["message"])
	if err != nil {
		return nil, err
	}

	// Write a JSON containg the processed response
	return map[string]string{
		"message": message,
	}, nil

}
