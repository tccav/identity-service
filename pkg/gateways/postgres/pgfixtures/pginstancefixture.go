package pgfixtures

import (
	"context"
	"fmt"
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
	dbTestEnv       = "DB_TEST_URL"
)

var (
	dbTestURL = "postgres://postgres:changeme@localhost:5432/%s?sslmode=disable"
	once      = sync.Once{}
)

func NewDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	abs, err := filepath.Abs(projectRootPath)
	require.NoError(t, err)

	fs := os.DirFS(abs)

	once.Do(func() {
		if dURL := os.Getenv(dbTestEnv); dURL != "" {
			log.Printf("setting db test url template: %s", dURL)
			dbTestURL = dURL
		}
	})

	dbUrl := fmt.Sprintf(dbTestURL, uuid.NewString())
	log.Printf("test db url: %s", dbUrl)
	u, _ := url.Parse(dbUrl)

	migrator := dbmate.New(u)
	migrator.FS = fs

	migrator.WaitInterval = 5 * time.Second
	migrator.WaitTimeout = 20 * time.Second
	migrator.WaitBefore = true

	err = migrator.CreateAndMigrate()
	require.NoError(t, err)

	pool, err := pgxpool.New(context.Background(), u.String())
	require.NoError(t, err)

	return pool
}
