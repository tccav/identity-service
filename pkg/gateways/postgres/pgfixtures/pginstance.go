package pgfixtures

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

const (
	projectRootPath = "./../../../../"
)

var (
	migrationFS fs.FS
	dbTestURL   = "postgres://postgres:changeme@localhost:5432/%s?sslmode=disable"
	once        = sync.Once{}
)

func NewDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	once.Do(func() {
		abs, err := filepath.Abs(projectRootPath)
		require.NoError(t, err)

		migrationFS = os.DirFS(abs)
	})

	dbUrl := fmt.Sprintf(dbTestURL, uuid.NewString())
	log.Printf("test db url: %s", dbUrl)
	u, _ := url.Parse(dbUrl)

	migrator := dbmate.New(u)
	migrator.FS = migrationFS

	migrator.WaitInterval = 5 * time.Second
	migrator.WaitTimeout = 20 * time.Second
	migrator.WaitBefore = true

	err := migrator.CreateAndMigrate()
	require.NoError(t, err)

	pool, err := pgxpool.New(context.Background(), u.String())
	require.NoError(t, err)

	return pool
}
