package response

import "github.com/stannisl/avito-test/internal/domain"

type Error struct {
	Error *domain.Error `json:"error"`
}
