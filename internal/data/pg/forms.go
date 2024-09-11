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
	formsTable = "forms"
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

func (q *formsQ) Insert(form data.Form) error {
	values := map[string]interface{}{
		"id":             form.ID,
		"nullifier":      form.Nullifier,
		"status":         form.Status,
		"name":           form.Name,
		"surname":        form.Surname,
		"id_num":         form.IDNum,
		"birthday":       form.Birthday,
		"citizen":        form.Citizen,
		"visited":        form.Visited,
		"purpose":        form.Purpose,
		"country":        form.Country,
		"city":           form.City,
		"address":        form.Address,
		"postal":         form.Postal,
		"phone":          form.Phone,
		"email":          form.Email,
		"image":          form.Image,
		"passport_image": form.PassportImage,
	}

	if err := q.db.Exec(squirrel.Insert(formsTable).SetMap(values)); err != nil {
		return fmt.Errorf("insert form: %w", err)
	}

	return nil
}

func (q *formsQ) Update(fields map[string]any) error {
	if err := q.db.Exec(q.updater.SetMap(fields)); err != nil {
		return fmt.Errorf("update forms: %w", err)
	}

	return nil
}

func (q *formsQ) Select() ([]data.Form, error) {
	var res []data.Form

	if err := q.db.Select(&res, q.selector); err != nil {
		return nil, fmt.Errorf("select forms: %w", err)
	}

	return res, nil
}

func (q *formsQ) Get() (*data.Form, error) {
	var res data.Form

	if err := q.db.Get(&res, q.selector); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get form by id: %w", err)
	}

	return &res, nil
}

func (q *formsQ) Last() (*data.Form, error) {
	var res data.Form

	stmt := q.selector.OrderBy("created_at DESC")
	if err := q.db.Get(&res, stmt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get last form by nullifier: %w", err)
	}

	return &res, nil
}

func (q *formsQ) FilterByID(ids ...string) data.FormsQ {
	return q.applyCondition(squirrel.Eq{"id": ids})
}

func (q *formsQ) FilterByNullifier(nullifiers ...string) data.FormsQ {
	return q.applyCondition(squirrel.Eq{"nullifier": nullifiers})
}

func (q *formsQ) FilterByStatus(status ...string) data.FormsQ {
	return q.applyCondition(squirrel.Eq{"status": status})
}

func (q *formsQ) applyCondition(cond squirrel.Sqlizer) data.FormsQ {
	q.selector = q.selector.Where(cond)
	q.updater = q.updater.Where(cond)
	return q
}
