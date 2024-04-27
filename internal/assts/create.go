package assts

import (
	"github.com/jackitaliano/oait/internal/openai"
)

func CreateAssistant(key string, createdAsst *openai.CreatedAssistant, orgId string) (*openai.AsstObject, error) {
	asst, err := openai.CreateAssistant(key, createdAsst, orgId)

	return asst, err
}
