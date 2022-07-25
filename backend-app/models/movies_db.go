package models

import (
	"context"
	"database/sql"
	"time"
)

const defaultDatabaseTimeout = 3 * time.Second

// DBModel is the database pool.
type DBModel struct {
	DB *sql.DB
}

// Get returns a single movie and error, if any.
func (m *DBModel) Get(id int) (*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDatabaseTimeout)
	defer cancel()

	// get the requested movie
	query := `
	select
		id,
		title,
		description,
		year,
		release_date,
		rating,
		runtime,
		mpaa_rating,
		created_at,
		updated_at
	from
		movies
	where
		id = $1`

	var movie Movie

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.Year,
		&movie.ReleaseDate,
		&movie.Rating,
		&movie.Runtime,
		&movie.MPAARating,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// get movie's genres, if any
	query = `
	select
		mg.id,
		mg.movie_id,
		mg.genre_id,
		mg.created_at,
		mg.updated_at,
		g.id,
		g.genre_name,
		g.created_at,
		g.updated_at
	from
		movies_genres mg
	left join genres g on
		(mg.genre_id = g.id)
	where
		mg.movie_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	genres := make(map[int]string)

	for rows.Next() {
		var mg MovieGenre

		err := rows.Scan(
			&mg.ID,
			&mg.MovieID,
			&mg.GenreID,
			&mg.CreatedAt,
			&mg.UpdatedAt,
			&mg.Genre.ID,
			&mg.Genre.GenreName,
			&mg.Genre.CreatedAt,
			&mg.Genre.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		genres[mg.ID] = mg.Genre.GenreName
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	movie.MovieGenre = genres

	return &movie, nil
}

// All returns all movies and error, if any.
func (m *DBModel) All() ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDatabaseTimeout)
	defer cancel()

	query := `
	select
		id,
		title,
		description,
		year,
		release_date,
		rating,
		runtime,
		mpaa_rating,
		created_at,
		updated_at
	from
		movies
	order by
		title`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*Movie

	for rows.Next() {
		var movie Movie

		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.Year,
			&movie.ReleaseDate,
			&movie.Rating,
			&movie.Runtime,
			&movie.MPAARating,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// get movie's genres, if any
		genreQuery := `
		select
			mg.id,
			mg.movie_id,
			mg.genre_id,
			mg.created_at,
			mg.updated_at,
			g.id,
			g.genre_name,
			g.created_at,
			g.updated_at
		from
			movies_genres mg
		left join genres g on
			(mg.genre_id = g.id)
		where
			mg.movie_id = $1`

		genreRows, err := m.DB.QueryContext(ctx, genreQuery, movie.ID)
		if err != nil {
			return nil, err
		}

		genres := make(map[int]string)

		for genreRows.Next() {
			var mg MovieGenre

			err := genreRows.Scan(
				&mg.ID,
				&mg.MovieID,
				&mg.GenreID,
				&mg.CreatedAt,
				&mg.UpdatedAt,
				&mg.Genre.ID,
				&mg.Genre.GenreName,
				&mg.Genre.CreatedAt,
				&mg.Genre.UpdatedAt,
			)
			if err != nil {
				return nil, err
			}

			genres[mg.ID] = mg.Genre.GenreName
		}

		if err = genreRows.Err(); err != nil {
			return nil, err
		}

		genreRows.Close()

		movie.MovieGenre = genres

		movies = append(movies, &movie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}
