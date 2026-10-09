package testutils

import (
	"context"
	"fmt"
	"go-api-template/internal/config"
	"go-api-template/internal/database"
	"go-api-template/migrations"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/peterldowns/pgtestdb"
	"github.com/peterldowns/pgtestdb/migrators/goosemigrator"
)

func NewTestDB(t pgtestdb.TB) *database.DB {
	t.Helper()

	cfg, err := config.InitConfig()
	if err != nil {
		t.Fatalf("fail to get test env: %v", err)
		return nil
	}

	primaryPoolConfig, err := pgxpool.ParseConfig(cfg.Database.PrimaryDSN)
	if err != nil {
		t.Fatalf("fail to parse primary DSN: %v", err)
		return nil
	}

	conf := pgtestdb.Config{
		DriverName: "pgx",
		User:       primaryPoolConfig.ConnConfig.User,
		Password:   primaryPoolConfig.ConnConfig.Password,
		Host:       primaryPoolConfig.ConnConfig.Host,
		Port:       fmt.Sprintf("%d", primaryPoolConfig.ConnConfig.Port),
		Database:   primaryPoolConfig.ConnConfig.Database,
		Options:    "sslmode=disable",
	}

	migrator := goosemigrator.New(".", goosemigrator.WithFS(migrations.FS))
	testDBConf := pgtestdb.Custom(t, conf, migrator)
	primaryTestDSN := testDBConf.URL()

	replicaPoolConfig, err := pgxpool.ParseConfig(cfg.Database.ReplicaDSN)
	if err != nil {
		t.Fatalf("fail to parse replica DSN: %v", err)
		return nil
	}

	replicaTestDSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		testDBConf.User,
		testDBConf.Password,
		testDBConf.Host,
		replicaPoolConfig.ConnConfig.Port,
		testDBConf.Database,
	)

	ctx := context.Background()

	db, err := database.NewDB(ctx, config.Database{
		PrimaryDSN: primaryTestDSN,
		ReplicaDSN: replicaTestDSN,
	})
	if err != nil {
		t.Fatalf("fail to create DB client: %v", err)
		return nil
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func LoadFixtures(t pgtestdb.TB, db *database.DB, sqlFile string) {
	t.Helper()
	ctx := context.Background()
	file, err := os.ReadFile(sqlFile)
	if err != nil {
		t.Fatalf("fail to open fixture file: %v", err)
		return
	}

	tx, err := db.Primary.Begin(ctx)
	if err != nil {
		t.Fatalf("fail to start transaction: %v", err)
		return
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	for q := range strings.SplitSeq(string(file), ";") {
		q = strings.TrimSpace(q)
		if q == "" {
			continue
		}

		if _, err := tx.Exec(ctx, q); err != nil {
			t.Fatalf("fail to execute fixture statement from file %s: %v", sqlFile, err)
			return
		}
	}

	if err = tx.Commit(ctx); err != nil {
		t.Fatalf("fail to commit transaction: %v", err)
		return
	}
}
