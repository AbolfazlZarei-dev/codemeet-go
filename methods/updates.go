package methods

import (
	"context"

	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

type UpdatesMethods struct {
	parent *Messages
}

func newUpdatesMethods(m *Messages) *UpdatesMethods { return &UpdatesMethods{parent: m} }

type GetUpdatesParams struct {
	Offset         int      `json:"offset,omitempty"`
	Limit          int      `json:"limit,omitempty"`
	Timeout        int      `json:"timeout,omitempty"`
	AllowedUpdates []string `json:"allowed_updates,omitempty"`
}

func (u *UpdatesMethods) Get(ctx context.Context, p *GetUpdatesParams) ([]models.Update, error) {
	var updates []models.Update
	err := u.parent.doWithRetry(ctx, "getUpdates", p, &updates)
	return updates, err
}
