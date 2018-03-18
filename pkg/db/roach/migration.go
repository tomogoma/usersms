package roach

import (
	"fmt"
	"github.com/tomogoma/crdb"
	"github.com/tomogoma/go-typed-errors"
)

func (r *Roach) migrate(fromVersion, toVersion int) error {

	var err error
	r.db, err = crdb.TryConnect(r.dsn, r.db)
	if err != nil {
		return fmt.Errorf("connect to db: %v", err)
	}

	if fromVersion == 0 && toVersion == 1 {
		if err := r.migrate0To1(); err != nil {
			return err
		}
		return r.setRunningVersionCurrent()
	}

	return errors.New("not supported")
}

func (r *Roach) migrate0To1() error {
	q := `
		ALTER TABLE ` + TblUsers + `
			ALTER ` + ColName + ` DROP NOT NULL,
			ALTER ` + ColGender + ` DROP NOT NULL
	`
	_, err := r.db.Exec(q)
	if err != nil {
		return fmt.Errorf("migrate %s table: %v", TblUsers, err)
	}
	return nil
}
