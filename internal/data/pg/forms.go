package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/rarimo/geo-forms-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	formsTable        = "forms"
	formsStatusFields = "id, nullifier, status, created_at, updated_at"
)

type formsQ struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
	updater  squirrel.UpdateBuilder
}

func NewForms(db *pgdb.DB) data.FormsQ {
	return &formsQ{
		db:       db,
		selector: squirrel.Select("*").From(formsTable),
		updater:  squirrel.Update(formsTable),
	}
}

func (q *formsQ) New() data.FormsQ {
	return NewForms(q.db)
}

func (q *formsQ) Insert(form *data.Form) (*data.FormStatus, error) {
	var res data.FormStatus

	stmt := squirrel.Insert(formsTable).SetMap(map[string]interface{}{
		"nullifier": form.Nullifier,
		"status":    form.Status,
		"name":      form.Name,
		"surname":   form.Surname,
		"id_num":    form.IDNum,
		"birthday":  form.Birthday,
		"citizen":   form.Citizen,
		"visited":   form.Visited,
		"purpose":   form.Purpose,
		"country":   form.Country,
		"city":      form.City,
		"address":   form.Address,
		"postal":    form.Postal,
		"phone":     form.Phone,
		"email":     form.Email,
		"image":     form.Image,
		"image_url": form.ImageURL,
	}).Suffix("RETURNING id, nullifier, status, created_at, updated_at")

	if err := q.db.Get(&res, stmt); err != nil {
		return nil, fmt.Errorf("insert form: %w", err)
	}

	return &res, nil
}

func (q *formsQ) Update(status string) error {
	if err := q.db.Exec(q.updater.Set("status", status)); err != nil {
		return fmt.Errorf("update forms: %w", err)
	}

	return nil
}

func (q *formsQ) Select() ([]*data.Form, error) {
	var res []*data.Form

	if err := q.db.Select(&res, q.selector); err != nil {
		return nil, fmt.Errorf("select forms: %w", err)
	}

	return res, nil
}

func (q *formsQ) Get(id string) (*data.FormStatus, error) {
	var res data.FormStatus

	stmt := squirrel.Select(formsStatusFields).From(formsTable).Where(squirrel.Eq{"id": id})
	if err := q.db.Get(&res, stmt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get form status by id=%s: %w", id, err)
	}

	return &res, nil
}

func (q *formsQ) Last(nullifier string) (*data.FormStatus, error) {
	var res data.FormStatus

	stmt := squirrel.Select(formsStatusFields).From(formsTable).Where(squirrel.Eq{"nullifier": nullifier}).OrderBy("created_at DESC")
	if err := q.db.Get(&res, stmt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get last form status by nullifier=%s: %w", nullifier, err)
	}

	return &res, nil
}

func (q *formsQ) Limit(limit uint64) data.FormsQ {
	q.selector = q.selector.Limit(limit)
	return q
}

func (q *formsQ) FilterByID(ids ...string) data.FormsQ {
	return q.applyCondition(squirrel.Eq{"id": ids})
}

func (q *formsQ) FilterByNullifier(nullifier string) data.FormsQ {
	return q.applyCondition(squirrel.Eq{"nullifier": nullifier})
}

func (q *formsQ) FilterByStatus(status string) data.FormsQ {
	return q.applyCondition(squirrel.Eq{"status": status})
}

func (q *formsQ) applyCondition(cond squirrel.Sqlizer) data.FormsQ {
	q.selector = q.selector.Where(cond)
	q.updater = q.updater.Where(cond)
	return q
}
