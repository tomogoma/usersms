package user

type User struct {
	ID        string
	Name      string
	ICEPhone  string
	Gender    string
	AvatarURL string
	Bio       string
	Rating    float32
	NumRaters int64
}

type Int64Update struct {
	IsUpdating bool
	NewValue   int64
}

type StringUpdate struct {
	IsUpdating bool
	NewValue   string
}

type Float32Update struct {
	IsUpdating bool
	NewValue   float32
}

type UserUpdate struct {
	Name      StringUpdate
	ICEPhone  StringUpdate
	Gender    StringUpdate
	AvatarURL StringUpdate
	Bio       StringUpdate
}
