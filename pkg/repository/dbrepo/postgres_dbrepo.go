package dbrepo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sdblg/meme/pkg/models"
)

// PostgresDBRepo is the struct used to wrap our database connection pool, so that we
// can easily swap out a real database for a test database, or move to another database
// entirely, as long as the thing being swapped implements all of the functions in the type
// repository.DatabaseRepo.
type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 3

// Connection returns underlying connection pool.
func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

// AllMemes returns a slice of memes, sorted by name. If the optional parameter genre
// is supplied, then only all memes for a particular genre is returned.
func (m *PostgresDBRepo) AllMemes() ([]*models.Meme, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	where := ""

	query := fmt.Sprintf(`
		select
			id, title, release_date, runtime,
			mpaa_rating, description, coalesce(image, ''),
			created_at, updated_at
		from
			memes %s
		order by
			title
	`, where)

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memes []*models.Meme

	for rows.Next() {
		var meme models.Meme
		err := rows.Scan(
			&meme.ID,
			&meme.Title,
			&meme.ReleaseDate,
			&meme.RunTime,
			&meme.MPAARating,
			&meme.Description,
			&meme.Image,
			&meme.CreatedAt,
			&meme.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		memes = append(memes, &meme)
	}

	return memes, nil
}

// OneMeme returns a single meme and associated categories, if any.
func (m *PostgresDBRepo) OneMeme(id int) (*models.Meme, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, title, release_date, runtime, mpaa_rating, 
		description, coalesce(image, ''), created_at, updated_at
		from memes where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var meme models.Meme

	err := row.Scan(
		&meme.ID,
		&meme.Title,
		&meme.ReleaseDate,
		&meme.RunTime,
		&meme.MPAARating,
		&meme.Description,
		&meme.Image,
		&meme.CreatedAt,
		&meme.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// get categories, if any
	query = `select g.id, g.genre from memes_genres mg
		left join categories g on (mg.genre_id = g.id)
		where mg.meme_id = $1
		order by g.genre`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Genre
	for rows.Next() {
		var g models.Genre
		err := rows.Scan(
			&g.ID,
			&g.Genre,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, &g)
	}

	meme.Genres = categories

	return &meme, err
}

// GetUserByEmail returns one use, by email.
func (m *PostgresDBRepo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password,
			created_at, updated_at from users where email = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByID returns one use, by ID.
func (m *PostgresDBRepo) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password,
			created_at, updated_at from users where id = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// InsertMeme inserts one meme into the database.
func (m *PostgresDBRepo) InsertMeme(meme models.Meme) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into memes (title, description, release_date, runtime,
			mpaa_rating, created_at, updated_at, image)
			values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`

	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		meme.Title,
		meme.Description,
		meme.ReleaseDate,
		meme.RunTime,
		meme.MPAARating,
		meme.CreatedAt,
		meme.UpdatedAt,
		meme.Image,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// UpdateMeme updates one meme in the database.
func (m *PostgresDBRepo) UpdateMeme(meme models.Meme) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update memes set title = $1, description = $2, release_date = $3,
				runtime = $4, mpaa_rating = $5,
				updated_at = $6, image = $7 where id = $8`

	_, err := m.DB.ExecContext(ctx, stmt,
		meme.Title,
		meme.Description,
		meme.ReleaseDate,
		meme.RunTime,
		meme.MPAARating,
		meme.UpdatedAt,
		meme.Image,
		meme.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteMeme deletes one meme, by id.
func (m *PostgresDBRepo) DeleteMeme(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from memes where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}
