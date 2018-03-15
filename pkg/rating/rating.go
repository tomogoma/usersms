package rating

type Rating struct {
	ID          string
	ForUserID   string
	ByUserID    string
	Comment     string
	CreateDate  string
	UpdateDate  string
	Rating      int32
}

type Filter struct {
	ForUserID   string
	ByUserID    string
	Offset      int64
	Count       int32
}
