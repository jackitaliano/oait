package threads

import (
	"encoding/json"
	"errors"

	"github.com/jackitaliano/oait/internal/openai"
)

func AddMessage(key string, threadID string, createdMessage *openai.CreatedMessage, orgID string) (*openai.Message, error) {
	message, err := openai.AddMessage(key, threadID, createdMessage, orgID)

	if err != nil {
		return nil, err
	}

	return message, nil
}

func CreateMessage(text string, role string) *openai.CreatedMessage {
	message := openai.CreatedMessage{Role: role, Content: text}

	return &message
}

func CreatedMessageToJson(message *openai.CreatedMessage) ([]byte, error) {
	b, err := json.MarshalIndent(message, "", "\t")

	if err != nil {
		err = errors.New("JSON Marshal failed with error: " + err.Error())
		return nil, err
	}

	return b, nil
}
