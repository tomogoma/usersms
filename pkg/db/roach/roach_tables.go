package roach

const (
	// Database definition version
	Version = 0

	// Table names
	TblConfigurations = "configurations"
	TblAPIKeys        = "api_keys"
	TblUsers          = "users"
	TblRatings        = "ratings"

	// DB Table Columns
	ColID          = "ID"
	ColCreated     = "created"
	ColLastUpdated = "last_updated"
	ColUserID      = "user_id"
	ColKey         = "key"
	ColValue       = "value"
	ColForSection  = "for_section"
	ColForUserID   = "for_user_id"
	ColByUserID    = "by_user_id"
	ColComment     = "comment"
	ColRating      = "rating"
	ColName        = "name"
	ColICEPhone    = "ice_phone"
	ColGender      = "gender"
	ColAvatarURL   = "avatar_url"
	ColBio         = "bio"
	ColNumRaters   = "num_raters"

	// CREATE TABLE DESCRIPTIONS
	TblDescConfigurations = `
	CREATE TABLE IF NOT EXISTS ` + TblConfigurations + ` (
		` + ColKey + ` VARCHAR(56) PRIMARY KEY NOT NULL CHECK (` + ColKey + ` != ''),
		` + ColValue + ` BYTEA NOT NULL CHECK (` + ColValue + ` != ''),
		` + ColCreated + ` TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		` + ColLastUpdated + ` TIMESTAMPTZ NOT NULL
	);
	`
	TblDescAPIKeys = `
	CREATE TABLE IF NOT EXISTS ` + TblAPIKeys + ` (
		` + ColID + ` SERIAL PRIMARY KEY NOT NULL CHECK (` + ColID + `>0),
		` + ColUserID + ` INTEGER NOT NULL,
		` + ColKey + ` VARCHAR(256) NOT NULL CHECK ( LENGTH(` + ColKey + `) >= 56 ),
		` + ColCreated + ` TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		` + ColLastUpdated + ` TIMESTAMPTZ NOT NULL
	);
	`

	TblDescUsers = `
	CREATE TABLE IF NOT EXISTS ` + TblUsers + ` (
		` + ColID + ` VARCHAR(56) PRIMARY KEY CHECK (` + ColID + ` != ''),
		` + ColName + ` VARCHAR(256) NOT NULL CHECK (` + ColName + ` != ''),
		` + ColGender + ` VARCHAR(16) PRIMARY KEY CHECK (` + ColGender + ` IN ('MALE', 'FEMALE', 'OTHER')),
		` + ColICEPhone + ` VARCHAR(24),
		` + ColAvatarURL + ` VARCHAR(256),
		` + ColBio + ` TEXT,
		` + ColRating + ` REAL,
		` + ColNumRaters + ` INT,
		` + ColCreated + ` TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		` + ColLastUpdated + ` TIMESTAMPTZ NOT NULL
	);
	`

	TblDescRatings = `
	CREATE TABLE IF NOT EXISTS ` + TblRatings + ` (
		` + ColID + ` VARCHAR(56) PRIMARY KEY CHECK (` + ColID + ` != ''),
		` + ColForUserID + ` VARCHAR(56) NOT NULL REFERENCES ` + TblUsers + ` (` + ColID + `),
		` + ColByUserID + ` VARCHAR(56) NOT NULL REFERENCES ` + TblUsers + ` (` + ColID + `),
		` + ColForSection + ` VARCHAR(256) NOT NULL CHECK (` + ColForSection + ` != ''),
		` + ColRating + ` INT NOT NULL CHECK (` + ColRating + ` >= 1 AND ` + ColRating + ` <= 5),
		` + ColComment + ` TEXT,
		` + ColCreated + ` TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		` + ColLastUpdated + ` TIMESTAMPTZ NOT NULL
	);
	`
)

// AllTableDescs lists all CREATE TABLE DESCRIPTIONS in order of dependency
// (tables with foreign key references listed after parent table descriptions).
var AllTableDescs = []string{
	TblDescConfigurations,
	TblDescAPIKeys,
	TblDescUsers,
	TblDescRatings,
}

// AllTableNames lists all table names in order of dependency
// (tables with foreign key references listed after parent table descriptions).
var AllTableNames = []string{
	TblConfigurations,
	TblAPIKeys,
	TblUsers,
	TblRatings,
}
