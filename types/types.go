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
	MigrationsFolder string
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
	// TODO: Get the salt from env variable.
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

type CreateUserPayload struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"-" validate:"required,gte=3"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"-" validate:"required,gte=3"`
}

type User struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   Password  `json:"-"`
	IsVerified bool      `json:"is_verified"`
	CreatedOn  time.Time `json:"created_on"`
	RoleId     int64     `json:"role_id"`
	Role       Role      `json:"role,omitempty"`
}

type Role struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Editor struct {
	Id     int64 `json:"id"`
	UserId int64 `json:"user_id"`
}
