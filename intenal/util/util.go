package util

import (
	"database/sql"
	"io/ioutil"
	"strings"
)

//Contains return true or false
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

//LoadFixture is...
func LoadFixture(dbConnPool *sql.DB, fixturePath string) error {
	if fixturePath != "" {
		input, err := ioutil.ReadFile(fixturePath)
		if err != nil {
			return err
		}
		queries := strings.Split(string(input), ";")
		for _, query := range queries {
			_, err = dbConnPool.Exec(query)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
