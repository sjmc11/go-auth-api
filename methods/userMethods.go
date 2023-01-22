package methods

import (
	"auth-api/db"
	"auth-api/helpers"
	"context"
	"database/sql"
	"fmt"
	"github.com/logrusorgru/aurora/v4"
	"strconv"
	"strings"
	"time"
)

type CoreUserData struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email,omitempty"`
	Role      int       `json:"role,omitempty"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type User struct {
	CoreUserData
	ID    int            `json:"id,omitempty"`
	Token sql.NullString `json:"token,omitempty" db:"token"`
	// Activity
	LastLogin sql.NullTime `json:"last_login"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

type StrippedUser struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      int    `json:"role"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	LastLogin int64  `json:"last_login"`
}

func (u *User) Strip() StrippedUser {
	// Set the basics
	formattedUser := StrippedUser{
		ID:        u.ID,
		Name:      u.FirstName + " " + u.LastName,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt.Unix(),
	}

	// If has valid `last_login` - format to unix timestamp
	if u.LastLogin.Valid {
		formattedUser.LastLogin = u.LastLogin.Time.Unix()
	}

	// If has valid `updated_at` - format to unix timestamp
	if u.UpdatedAt.Valid {
		formattedUser.UpdatedAt = u.UpdatedAt.Time.Unix()
	}

	return formattedUser
}

func (u *CoreUserData) WriteUser() error {
	//dbConn := db.PgConnection.Db.PgConn()

	insertStmnt := fmt.Sprintf(`INSERT INTO system.users (first_name, last_name, email, password) VALUES ('%s', '%s', '%s', '%s');`,
		helpers.Capitalize(u.FirstName),
		helpers.Capitalize(u.LastName),
		strings.TrimSpace(strings.ToLower(u.Email)),
		helpers.StringToHash(u.Password),
	)

	fmt.Println(aurora.Gray(12, insertStmnt))

	//defer func() {
	//	dbConn.Close(context.Background())
	//}()

	_, err := db.PgConnection.Db.Exec(context.Background(), insertStmnt)
	if err != nil {
		return err
	}
	return nil
}

// GetUserBy /**
func GetUserBy(searchColumn string, searchVal string) (User, error) {
	var userData User

	queryStmnt := fmt.Sprintf(
		`SELECT * FROM system.users WHERE %s = $1`,
		searchColumn,
	)

	err := db.PgConnection.Db.QueryRow(context.Background(), queryStmnt, strings.TrimSpace(searchVal)).Scan(
		&userData.FirstName,
		&userData.LastName,
		&userData.Email,
		&userData.ID,
		&userData.Token,
		&userData.Password,
		&userData.LastLogin,
		&userData.CreatedAt,
		&userData.UpdatedAt,
		&userData.Role,
	)
	if err != nil {
		fmt.Println(err.Error())
		return userData, err
	}

	return userData, nil
}

// GenerateToken
// Set token in DB as well as last_login timestamp

func (u *User) GenerateToken() (string, error) {

	newToken := helpers.TokenizeString(u.Email)

	updateStmnt := fmt.Sprintf(
		`UPDATE system.users SET token = '%s', last_login = '%s' where id = '%s'`,
		newToken,
		time.Now().Format("2006-01-02 15:04:05"),
		strconv.Itoa(u.ID),
	)

	fmt.Println(aurora.Gray(12, updateStmnt))

	_, err := db.PgConnection.Db.Exec(context.Background(), updateStmnt)
	if err != nil {
		return "", err
	}

	return newToken, nil
}

func (u *User) UpdatePassword(password string) error {

	updateStmnt := fmt.Sprintf(
		`UPDATE system.users SET password = '%s', updated_at = '%s' where id = '%s'`,
		helpers.StringToHash(password),
		time.Now().Format("2006-01-02 15:04:05"),
		strconv.Itoa(u.ID),
	)

	fmt.Println(aurora.Gray(12, updateStmnt))

	_, err := db.PgConnection.Db.Exec(context.Background(), updateStmnt)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) UpdateUser() error {

	updateStmnt := fmt.Sprintf(
		`UPDATE system.users SET first_name = '%s', last_name = '%s', role = '%v', updated_at = '%s' where id = '%s'`,
		u.FirstName,
		u.LastName,
		strconv.Itoa(u.Role),
		time.Now().Format("2006-01-02 15:04:05"),
		strconv.Itoa(u.ID),
	)

	fmt.Println(aurora.Gray(12, updateStmnt))

	_, err := db.PgConnection.Db.Exec(context.Background(), updateStmnt)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) DeleteUser() error {

	updateStmnt := fmt.Sprintf(
		`DELETE from system.users where id = '%s'`,
		strconv.Itoa(u.ID),
	)

	fmt.Println(aurora.Gray(12, updateStmnt))

	_, err := db.PgConnection.Db.Exec(context.Background(), updateStmnt)
	if err != nil {
		return err
	}

	return nil
}
