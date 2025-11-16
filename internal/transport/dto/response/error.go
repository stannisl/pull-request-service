package response

import "github.com/stannisl/pull-request-service/internal/transport/dto"

type Error struct {
	Error *dto.Error `json:"error"`
}
