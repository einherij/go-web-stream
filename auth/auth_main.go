package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"strings"
)

const (
	uRegPath = "./unregusers/"
	regPath  = "./users/"
)

// Person contains usernames and information about users
type Person struct {
	Login    string
	Name     string
	Email    string
	Password string
}

// encryptPassword encrypts password via sha256
func encryptPassword(pass string) (encPass string) {
	hasher := sha256.New()
	hasher.Write([]byte(pass))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

// NewPerson creates new Person struct
func NewPerson(login string, name string, email string, password string) *Person {
	p := Person{login, name, email, encryptPassword(password)}
	return &p
}

// Save saves person to file
func (p *Person) Save() error {
	filename := uRegPath + p.Login + ".txt"
	return ioutil.WriteFile(filename, []byte(p.Name+"\n"+p.Email+"\n"+p.Password), 0600)
}

// loadPerson loades person from file by name
func loadPerson(name string) (p *Person) {
	filename := regPath + name + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	strBody := string(body)
	strBodySplited := strings.Split(strBody, "\n")
	return &Person{name, strBodySplited[0], strBodySplited[1], strBodySplited[2]}
}

// checkPassword checks that user exists and password is coincides
func (p *Person) checkPassword(pass string) bool {
	if encryptPassword(pass) == p.Password {
		return true
	}
	return false
}

// LoginUser check correct login and password
func LoginUser(name string, password string) bool {
	p := loadPerson(name)
	if p == nil {
		return false
	}
	pass := p.checkPassword(password)
	return pass
}
