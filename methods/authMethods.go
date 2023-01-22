package methods

import (
	"auth-api/db"
	"auth-api/helpers"
	"auth-api/helpers/env"
	"auth-api/mailer"
	"bytes"
	"context"
	"fmt"
	"github.com/ainsleyclark/go-mail/mail"
	"github.com/logrusorgru/aurora/v4"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"strconv"
	"time"
)

type PwReset struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	UserEmail string    `json:"user_email" db:"email"`
	AuthCode  string    `json:"auth_code"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
}

/**
IsValid
Check if a Password reset is expired or already used
*/
func (r *PwReset) IsValid() (bool, string) {
	switch true {
	case r.Used:
		return false, "This reset code has already been used"
	case time.Now().After(r.ExpiresAt):
		return false, "This reset code has expired"
	default:
		return true, ""
	}
}

/**
MarkAsUsed
Set a pw_reset used to true
*/
func (r *PwReset) MarkAsUsed() error {
	updateStmnt := fmt.Sprintf(
		`UPDATE system.pw_resets SET used = '%v' where id = '%s'`,
		true,
		strconv.Itoa(r.ID),
	)

	fmt.Println(aurora.Gray(12, updateStmnt))

	_, err := db.PgConnection.Db.Exec(context.Background(), updateStmnt)
	if err != nil {
		return err
	}

	return nil
}

/**
GetPwReset
*/
func GetPwReset(authCode string) (PwReset, error) {
	var pwResetData PwReset

	queryStmnt := fmt.Sprintf(
		`SELECT pw_resets.id, pw_resets.user_id, users.email, pw_resets.auth_code, pw_resets.expires_at, pw_resets.used
FROM system.pw_resets
JOIN system.users ON pw_resets.user_id = users.id
WHERE pw_resets.auth_code = '%s';`,
		//`SELECT * FROM system.pw_resets WHERE auth_code = '%s'`,
		authCode,
	)

	fmt.Println(aurora.Gray(12, queryStmnt))

	err := db.PgConnection.Db.QueryRow(context.Background(), queryStmnt).Scan(
		&pwResetData.ID,
		&pwResetData.UserID,
		&pwResetData.UserEmail,
		&pwResetData.AuthCode,
		&pwResetData.ExpiresAt,
		&pwResetData.Used,
	)
	if err != nil {
		fmt.Println(err.Error())
		return pwResetData, err
	}

	return pwResetData, nil
}

/**
CreatePwReset
Generate a password reset auth code and insert record to 'pw_resets' table.
*/
func (u *User) CreatePwReset() (string, error) {

	authCode := helpers.TokenizeString(u.Email)
	hoursValid := 4 * time.Hour

	insertStmnt := fmt.Sprintf(`INSERT INTO system.pw_resets (user_id, auth_code, expires_at) VALUES ('%s', '%s', '%s');`,
		strconv.Itoa(u.ID),
		authCode,
		time.Now().Add(hoursValid).Format("2006-01-02 15:04:05"),
	)

	fmt.Println(aurora.Gray(12, insertStmnt))

	_, err := db.PgConnection.Db.Exec(context.Background(), insertStmnt)
	if err != nil {
		return authCode, err
	}
	return authCode, nil
}

/**
SendPwResetEmail
Send a password reset email to the user
*/
func (u *User) SendPwResetEmail(authCode string) (string, error) {

	/** START TEMPLATE **/

	// parse the template file
	emailTemplate, err := template.ParseFiles("pw-reset.html")
	if err != nil {
		return "", err
	}

	// create a map to hold the template data
	emailData := map[string]string{
		"Name":      u.FirstName,
		"AuthCode":  authCode,
		"ActionUrl": env.Get("APPURL") + "/reset-password?auth=" + authCode,
	}

	// create a new bytes buffer
	var mailBuff bytes.Buffer

	// execute the template and write the output to the buffer
	err = emailTemplate.Execute(&mailBuff, emailData)
	if err != nil {
		return "", err
	}

	// convert the buffer to a string
	emailHtmlString := mailBuff.String()

	/** END TEMPLATE **/

	/** START EMAIL TRANSMISSION **/

	_, emailErr := mailer.MailDriver.Send(&mail.Transmission{
		Recipients:  []string{u.Email},
		Subject:     "Password Reset",
		HTML:        emailHtmlString,
		PlainText:   "Someone requested a password reset for this email address. Please visit this link to reset your password: " + env.Get("APPURL") + "/reset-password?auth=" + authCode,
		Attachments: nil,
	})

	if emailErr != nil {
		return "", emailErr
	}

	/** END EMAIL TRANSMISSION **/

	// print the output
	//fmt.Println(emailHtmlString)

	return emailHtmlString, nil
}

/**
CheckPasswordValid
*/
func CheckPasswordValid(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePlainPass := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlainPass)
	if err != nil {
		return false
	}
	return true
}
