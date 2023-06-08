package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tccav/identity-service/pkg/domain/entities"
	"github.com/tccav/identity-service/pkg/domain/identities"
)

type StudentsRepository struct {
	conn *pgxpool.Pool
}

func NewStudentsRepository(conn *pgxpool.Pool) StudentsRepository {
	return StudentsRepository{
		conn: conn,
	}
}

func (s StudentsRepository) CreateStudent(ctx context.Context, student entities.Student) error {
	const statement = `
	INSERT INTO students (id, name, secret, birth_date, cpf, email) VALUES (
		$1,
	    $2,
		$3,
	    $4,
	    $5,
	 	$6                                                                       
	)`

	exec, err := s.conn.Exec(ctx, statement, student.ID, student.Name, student.Secret, student.BirthDate, student.CPF, student.Email)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return identities.ErrStudentAlreadyExists
		}
		return err
	}

	if !exec.Insert() {
		return errors.New("student not stored")
	}

	return nil
}

func (s StudentsRepository) GetStudentSecret(ctx context.Context, id string) (string, error) {
	const query = `SELECT secret FROM students WHERE id=$1`

	var secret string
	err := s.conn.QueryRow(ctx, query, id).Scan(&secret)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", identities.ErrStudentNotFound
		}
		return "", err
	}

	return secret, nil
}
