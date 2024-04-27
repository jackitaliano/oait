package assts

import (
	"github.com/jackitaliano/oait/internal/openai"
)

func CreateAssistant(key string, createdAsst *openai.CreatedAssistant, orgID string) (*openai.AsstObject, error) {
	asst, err := openai.CreateAssistant(key, createdAsst, orgID)

	return asst, err
}
