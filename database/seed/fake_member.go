package seed

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
)

func createData(db *gorm.DB, valueStrings []string, valueArgs []interface{}) error {
	stmt := fmt.Sprintf("INSERT INTO users(id, email, password, first_name, last_name , is_active, is_verified_email) VALUES %s ON DUPLICATE KEY UPDATE updated_at=now()", strings.Join(valueStrings, ","))

	err := db.Exec(stmt, valueArgs...).Error
	if err != nil {
		return err
	}

	return nil
}

func FakeMember(db *gorm.DB) error {
	valueStrings := []string{}
	valueArgs := []interface{}{}

	for i := 1; i <= 2; i++ {
		password := fmt.Sprintf("password%v", i)
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return errors.Wrap(err, "bcrypt.GenerateFromPassword")
		}

		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?)")

		valueArgs = append(valueArgs, i)
		valueArgs = append(valueArgs, fmt.Sprintf("user%v@example.com", i))
		valueArgs = append(valueArgs, hashed)
		valueArgs = append(valueArgs, fmt.Sprintf("%v", gofakeit.FirstName()))
		valueArgs = append(valueArgs, fmt.Sprintf("%v", gofakeit.LastName()))
		valueArgs = append(valueArgs, 1)
		valueArgs = append(valueArgs, 1)

		if i%100 == 0 && i > 0 {
			err1 := createData(db, valueStrings, valueArgs)
			if err1 != nil {
				continue
			}

			valueStrings = []string{}
			valueArgs = []interface{}{}
		}
	}

	if len(valueStrings) > 0 {
		err1 := createData(db, valueStrings, valueArgs)
		if err1 != nil {
			return err1
		}

		valueStrings = []string{}
		valueArgs = []interface{}{}
	}

	fmt.Println("Fake data successfully")
	return nil
}
