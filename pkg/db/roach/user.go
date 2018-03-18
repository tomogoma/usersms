package roach

import (
	"database/sql"
	"fmt"
	"github.com/tomogoma/go-typed-errors"
	"github.com/tomogoma/usersms/pkg/user"
	"strings"
	"time"
)

var allUserCols = ColDesc(ColID, ColName, ColGender, ColICEPhone, ColAvatarURL,
	ColBio, ColRating, ColNumRaters, ColCreated, ColLastUpdated)

func (r *Roach) UpsertUser(uu user.UserUpdate) (*user.User, error) {
	if err := r.InitDBIfNot(); err != nil {
		return nil, err
	}

	// updCols columns and their args/params will only be used during update of
	// existing user.

	updCols, args := addStrUpdate(uu.Name, ColName, "", []interface{}{})
	updCols, args = addStrUpdate(uu.ICEPhone, ColICEPhone, updCols, args)
	updCols, args = addStrUpdate(uu.Gender, ColGender, updCols, args)
	updCols, args = addStrUpdate(uu.AvatarURL, ColAvatarURL, updCols, args)
	updCols, args = addStrUpdate(uu.Bio, ColBio, updCols, args)

	updCols = ColDesc(updCols, ColLastUpdated)
	args = append(args, uu.Time)
	updParams := genParams(len(args))

	// insCols columns and their args/params includes update columns and
	// columns inserted only during inserts.
	insCols := ColDesc(updCols, ColID, ColCreated)
	args = append(args, uu.UserID, uu.Time)
	insParams := genParams(len(args))

	q := `
		INSERT INTO ` + TblUsers + ` (` + insCols + `)
			VALUES (` + insParams + `)
			ON CONFLICT (` + ColID + `) DO
				UPDATE SET (` + updCols + `) = (` + updParams + `)
			RETURNING ` + allUserCols + `
	`
	usr, err := scanUser(r.db.QueryRow(q, args...))
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (r *Roach) UpdateUserRating(userID string, newRating float32, numRaters int64) error {
	if err := r.InitDBIfNot(); err != nil {
		return err
	}

	cols := ColDesc(ColRating, ColNumRaters)
	q := `UPDATE ` + TblUsers + ` SET (` + cols + `) = ($1, $2) WHERE ` + ColID + `=$3`
	res, err := r.db.Exec(q, newRating, numRaters, userID)
	return checkRowsAffected(res, err, 1)
}

func (r *Roach) User(userID string, offsetUpdateDate time.Time) (*user.User, error) {
	if err := r.InitDBIfNot(); err != nil {
		return nil, err
	}

	whereArgs := []interface{}{userID}
	where := ColUserID + "=$1"

	if !offsetUpdateDate.IsZero() {
		whereArgs = append(whereArgs, offsetUpdateDate)
		// We are sure we have a where clause so safe to use the AND operator here.
		where = fmt.Sprintf("%s AND %s > $2", where, ColLastUpdated)
	}

	q := `SELECT ` + allUserCols + ` FROM ` + TblUsers + ` WHERE ` + where
	usr, err := scanUser(r.db.QueryRow(q, whereArgs...))
	if err != nil {
		if sql.ErrNoRows == err {
			return nil, errors.NewNotFound("no user found for provided filters")
		}
		return nil, err
	}

	return usr, nil
}

// addStrUpdate adds the value of su to cols and args if su.IsUpdating.
// It returns the resulting cols, args.
func addStrUpdate(su user.StringUpdate, col, cols string, args []interface{}) (string, []interface{}) {
	if su.IsUpdating {
		cols = ColDesc(cols, col)
		args = append(args, su.NewValue)
	}
	return cols, args
}

func genParams(count int) string {
	params := ""
	for i := 1; i <= count; i++ {
		params = fmt.Sprintf("%s$%d, ", params, i)
	}
	return strings.TrimSuffix(params, ", ")
}

// scanUser extracts a user from s or returns an error if reported by s.
// The column order for s must be same order as allUserCols variable.
func scanUser(s multiScanner) (*user.User, error) {

	ICEPhone := sql.NullString{}
	avatarURL := sql.NullString{}
	bio := sql.NullString{}
	rating := sql.NullFloat64{}
	numRaters := sql.NullInt64{}
	usr := &user.User{}

	err := s.Scan(&usr.ID, &usr.Name, &usr.Gender, &ICEPhone, &avatarURL,
		&bio, &rating, &numRaters, &usr.Created, &usr.LastUpdated)
	if err != nil {
		return nil, err
	}

	usr.ICEPhone = ICEPhone.String
	usr.AvatarURL = avatarURL.String
	usr.Bio = bio.String
	usr.Rating = float32(rating.Float64)
	usr.NumRaters = numRaters.Int64

	return usr, nil
}
