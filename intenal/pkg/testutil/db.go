package testutil

import (
	"context"
	"database/sql"
	"friend_management/intenal/db"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

type TxDB struct {
	*sql.Tx
}

// BeginTx implements Beginner
func (txdb *TxDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (db.Transactor, error) {
	return txdb, nil
}

// Commit implements Transactor
func (txdb *TxDB) Commit() error {
	return nil
}

// Rollback implements Transactor
// After a call to Commit or Rollback, all operations on the
// transaction fail with ErrTxDone.
func (txdb *TxDB) Rollback() error {
	return nil
}

// WithTxDB calls a callback function with a new transaction that will be rolled
// back, so no data is actually written to the database. This is helpful for
// database-related tests where you don't have to care about tearing down a
// database.
func WithTxDB(t *testing.T, callback func(db.BeginnerExecutor)) {
	t.Helper()

	pgDB, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer pgDB.Close()

	tx, err := pgDB.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	callback(&TxDB{Tx: tx})
}

// LoadTestDataFile load test data from a file
func LoadTestDataFile(t *testing.T, tx db.Executor, filename string) {
	t.Helper()

	body, err := Read(filename)
	require.NoError(t, err)

	_, err = tx.Exec(string(body))
	require.NoError(t, err)
}

// Read reads a file completely. But if the file does not exist, try to find it in the parent directory, [repeat...]
func Read(filename string) ([]byte, error) {
	return ioutil.ReadFile(Abs(filename))
}

// Abs returns absolute path of a file in project directory
func Abs(filename string) string {
	cwd, err := os.Getwd()
	if err != nil {
		return filename
	}

	return _abs(cwd, filename)
}

func _abs(dirname, filename string) string {
	fullpath, err := filepath.Abs(dirname + "/" + filename)
	if err != nil {
		return filename
	}

	if _, err = os.Stat(fullpath); err != nil {
		parentdir := filepath.Dir(dirname)
		if parentdir == "/" {
			return filename
		}
		return _abs(parentdir, filename)
	}

	return fullpath
}
