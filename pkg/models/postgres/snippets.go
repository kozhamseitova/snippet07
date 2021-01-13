package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"se07/pkg/models"
	"strconv"
	"time"
)

const (
	insertSql                 = "INSERT INTO snippets (title,content,created,expires) VALUES ($1,$2,$3,$4) RETURNING id"
	getSnippetById            = "SELECT id, title, content, created, expires FROM snippets where id=$1 AND expires > now()"
	getLastTenCreatedSnippets = "SELECT id, title, content, created, expires FROM snippets WHERE expires > now() ORDER BY created DESC LIMIT 10"
)

type SnippetModel struct {
	Pool *pgxpool.Pool
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	var id int
	days, err := strconv.Atoi(expires)
	if err != nil {
		return 0, err
	}
	row := m.Pool.QueryRow(context.Background(), insertSql, title, content, time.Now(), time.Now().AddDate(0, 0, days))
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	s := &models.Snippet{}
	err := m.Pool.QueryRow(context.Background(), getSnippetById, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	snippets := []*models.Snippet{}
	rows, err := m.Pool.Query(context.Background(), getLastTenCreatedSnippets)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
