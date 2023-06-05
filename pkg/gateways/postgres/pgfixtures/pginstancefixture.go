package pgfixtures

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

const (
	projectRootPath = "./../../../../"
	pgTestURL       = "postgres://postgres:changeme@localhost:5432/%s?sslmode=disable"
)

func NewDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	abs, err := filepath.Abs(projectRootPath)
	require.NoError(t, err)

	fs := os.DirFS(abs)

	dbUrl := fmt.Sprintf(pgTestURL, uuid.NewString())
	u, _ := url.Parse(dbUrl)

	migrator := dbmate.New(u)
	migrator.FS = fs

	err = migrator.CreateAndMigrate()
	require.NoError(t, err)

	pool, err := pgxpool.New(context.Background(), u.String())
	require.NoError(t, err)

	return pool
}
