package wirepusher

type User struct {
	Username string // a username
	password []byte // sha1-hashed password for message encryption
	ID       string // wirepusherID
}

// NewUser creates a new user and hashes the (optional) password when it is provided.
func NewUser(userName string, id string, password string) *User {
	usr := &User{Username: userName, ID: id}
	if len(password) > 0 {
		hPass, err := hash(password)
		if err != nil {
			panic(err)
		}
		usr.password = hPass
	}
	return usr
}

// Password returns the first AES_KEY_SIZE bytes of the password hash and converts it to a hex-encoded string
func (u *User) Password() string {
	if len(u.password) == 0 {
		return ""
	}
	return bin2hex(u.password[:AES_KEY_SIZE])
}

// CanEncrypt returns true when a password is provided and the message can therefore be encrypted
func (u *User) CanEncrypt() bool {
	return len(u.password) > 0
}
