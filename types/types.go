package types

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type DatabaseConfig struct {
	User             string
	Password         string
	Name             string
	connectionString string
	MaxOpenConns     int
	MaxIdleConns     int
	MaxIdleTime      string
}

func (c *DatabaseConfig) GetConnectionString() string {
	return c.connectionString
}

func (c *DatabaseConfig) SetConnectionString(connectionString string) {
	if connectionString == "" {
		// generate connection string from config
		connectionString = fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", c.User, c.Password, c.Name)
	}
	c.connectionString = connectionString

}

type Password struct {
	Hashed []byte
	Text   string
}

func (p *Password) HashPassword() error {
	hashedpw, err := bcrypt.GenerateFromPassword([]byte(p.Text), 8)
	if err != nil {
		return err
	}
	p.Hashed = hashedpw

	return nil
}

func (p *Password) ComparePassword() error {
	err := bcrypt.CompareHashAndPassword(p.Hashed, []byte(p.Text))
	if err != nil {
		return err
	}

	return nil
}

type UserPayload struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"-"` // NOTE: I think this should be string since in payload we'll only get string password.
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}

type User struct {
	UserPayload
	Password Password
}
