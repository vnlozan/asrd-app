package pg

import (
	"auth/internal/dto"
	"auth/internal/repo"
	"context"
	"database/sql"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 3

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) repo.IUserStorage {
	return &UserStorage{db: db}
}

func (a *UserStorage) SelectAll(ctx context.Context) ([]*dto.User, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
	from users order by last_name`

	rows, err := a.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*dto.User

	for rows.Next() {
		var user dto.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (a *UserStorage) SelectOneByEmail(ctx context.Context, email string) (*dto.User, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where email = $1`

	var user dto.User
	row := a.db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *UserStorage) SelectOne(ctx context.Context, id int) (*dto.User, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where id = $1`

	var user dto.User
	row := a.db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *UserStorage) UpdateOne(ctx context.Context, user *dto.User) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	stmt := `update users set
		email = $1,
		first_name = $2,
		last_name = $3,
		user_active = $4,
		updated_at = $5
		where id = $6
	`

	_, err := a.db.ExecContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Active,
		time.Now(),
		user.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (a *UserStorage) DeleteOneByID(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := a.db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (a *UserStorage) InsertOne(ctx context.Context, user dto.User) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	var newID int
	stmt := `insert into users (email, first_name, last_name, password, user_active, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7) returning id`

	err = a.db.QueryRowContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
		user.Active,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (a *UserStorage) UpdateOnePassword(ctx context.Context, id int, encryptedPassword string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	// if err != nil {
	// 	return err
	// }

	stmt := `update users set password = $1 where id = $2`
	_, err := a.db.ExecContext(ctx, stmt, encryptedPassword, id)
	if err != nil {
		return err
	}

	return nil
}
