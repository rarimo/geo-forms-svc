package data

import (
	"database/sql"
	"time"
)

const (
	CreatedStatus   = "created"
	AcceptedStatus  = "accepted"
	ProcessedStatus = "processed"
)

const (
	ColStatus   = "status"
	ColName     = "name"
	ColSurname  = "surname"
	ColIDNum    = "id_num"
	ColBirthday = "birthday"
	ColCitizen  = "citizen"
	ColVisited  = "visited"
	ColPurpose  = "purpose"
	ColCountry  = "country"
	ColCity     = "city"
	ColAddress  = "address"
	ColPostal   = "postal"
	ColPhone    = "phone"
	ColEmail    = "email"
	ColImageURL = "image_url"
)

type Form struct {
	ID        string         `db:"id"`
	Nullifier string         `db:"nullifier"`
	Status    string         `db:"status"`
	Name      string         `db:"name"`
	Surname   string         `db:"surname"`
	IDNum     string         `db:"id_num"`
	Birthday  string         `db:"birthday"`
	Citizen   string         `db:"citizen"`
	Visited   string         `db:"visited"`
	Purpose   string         `db:"purpose"`
	Country   string         `db:"country"`
	City      string         `db:"city"`
	Address   string         `db:"address"`
	Postal    string         `db:"postal"`
	Phone     string         `db:"phone"`
	Email     string         `db:"email"`
	Image     *string        `db:"image"`
	ImageURL  sql.NullString `db:"image_url"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

type FormStatus struct {
	ID         string    `db:"id"`
	Nullifier  string    `db:"nullifier"`
	Status     string    `db:"status"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	NextFormAt time.Time
}

type FormsQ interface {
	New() FormsQ
	Insert(*Form) (*FormStatus, error)
	UpdateStatus(status string) error

	Update(map[string]interface{}) error

	Select() ([]*Form, error)
	Limit(uint64) FormsQ

	Get(id string) (*FormStatus, error)
	Last(nullifier string) (*FormStatus, error)

	FilterByID(ids ...string) FormsQ
	FilterByStatus(status ...string) FormsQ
}
