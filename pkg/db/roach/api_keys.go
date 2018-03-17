package roach

import (
	"database/sql"

	apiG "github.com/tomogoma/go-api-guard"
	"github.com/tomogoma/go-typed-errors"
	"github.com/tomogoma/usersms/pkg/api"
)

// InsertAPIKey inserts an API key for the userID.
func (r *Roach) InsertAPIKey(userID string, key []byte) (apiG.Key, error) {
	if err := r.InitDBIfNot(); err != nil {
		return nil, err
	}
	k := api.Key{UserID: userID, Val: key}
	insCols := ColDesc(ColUserID, ColKey, ColLastUpdated)
	retCols := ColDesc(ColID, ColCreated, ColLastUpdated)
	q := `
		INSERT INTO ` + TblAPIKeys + ` (` + insCols + `)
			VALUES ($1, $2, CURRENT_TIMESTAMP)
			RETURNING ` + retCols
	err := r.db.QueryRow(q, userID, key).Scan(&k.ID, &k.Created, &k.LastUpdated)
	if err != nil {
		return nil, err
	}
	return k, nil
}

// APIKeyByUserIDVal returns API keys for the provided userID/key combination.
func (r *Roach) APIKeyByUserIDVal(userID string, key []byte) (apiG.Key, error) {
	if err := r.InitDBIfNot(); err != nil {
		return nil, err
	}
	cols := ColDesc(ColID, ColUserID, ColKey, ColCreated, ColLastUpdated)
	q := `
	SELECT ` + cols + `
		FROM ` + TblAPIKeys + `
		WHERE ` + ColUserID + `=$1 AND ` + ColKey + `=$2`
	k := api.Key{}
	err := r.db.QueryRow(q, userID, key).
		Scan(&k.ID, &k.UserID, &k.Val, &k.Created, &k.LastUpdated)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFound("API key not found")
		}
		return nil, err
	}
	return k, nil
}
