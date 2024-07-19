package data

import (
	"time"
)

const (
	AcceptedStatus  = "accepted"
	ResendingStatus = "resending"
	ProcessedStatus = "processed"
)

type Form struct {
	ID        string    `db:"id"`
	Nullifier string    `db:"nullifier"`
	Status    string    `db:"status"`
	Name      string    `db:"name"`
	Surname   string    `db:"surname"`
	IDNum     string    `db:"id_num"`
	Birthday  string    `db:"birthday"`
	Citizen   string    `db:"citizen"`
	Visited   string    `db:"visited"`
	Purpose   string    `db:"purpose"`
	Country   string    `db:"country"`
	City      string    `db:"city"`
	Address   string    `db:"address"`
	Postal    string    `db:"postal"`
	Phone     string    `db:"phone"`
	Email     string    `db:"email"`
	Image     *string   `db:"image"`
	CreatedAt time.Time `db:"created_at"`
}

type FormsQ interface {
	New() FormsQ
	Insert(*Form) (string, error)
	Update(status string) error

	Select() ([]*Form, error)
	Get() (*Form, error)
	// last returns the most recent form
	Last() (*Form, error)
	Limit(uint64) FormsQ

	FilterByID(ids ...string) FormsQ
	FilterByNullifier(nullifier string) FormsQ
	FilterByStatus(status string) FormsQ
}
