package pg

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"strconv"

// 	"golang.org/x/crypto/bcrypt"
// )

// var conn = Connect()
// var db = New(conn)

// func (user *User) Create() {
// 	hashedPassword, err := HashPassword(user.Password)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	res, err := db.CreateUser(context.Background(), CreateUserParams{
// 		Username: user.Username,
// 		Password: hashedPassword,
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(res)
// }

// func (user *User) Authenticate() bool {
// 	res, err := db.GetIdUserByUsername(context.Background(), user.Username)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	var hashedPassword string

// 	return CheckPasswordHash(res.Password, hashedPassword)
// }

// // GetUserIdByUsername check if a user exists in database by given username
// func GetUserIdByUsername(username string) (int, error) {
// 	res, err := db.GetIdUserByUsername(context.Background(), username)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	Id, _ := strconv.Atoi(res.ID)

// 	return Id, nil
// }

// //GetUserByID check if a user exists in database and return the user object.
// func GetUsernameById(userId string) (User, error) {
// 	res, err := db.GetUserById(context.Background(), userId)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return User{ID: res.ID, Username: res.Username}, nil
// }

// //HashPassword hashes given password
// func HashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
// 	return string(bytes), err
// }

// //CheckPassword hash compares raw password with it's hashed values
// func CheckPasswordHash(password, hash string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
// 	return err == nil
// }
