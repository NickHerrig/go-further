package data

import "time"

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Passowrd  password  `json:"-"`
	Activated bool      `json:"activated:`
	Version   int       `json:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	return false, nil
}
