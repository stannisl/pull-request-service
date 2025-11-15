package response

import (
	"github.com/stannisl/avito-test/internal/transport/dto"
)

type Error struct {
	Error *dto.Error `json:"error"`
}
