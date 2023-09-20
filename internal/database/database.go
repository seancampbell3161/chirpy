package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

type DB struct {
	path string
	mux  *sync.Mutex
}
type DBStructure struct {
	Chirps        map[int]Chirp        `json:"chirps"`
	Users         map[int]User         `json:"users"`
	RevokedTokens map[string]time.Time `json:"revoked_tokens"`
}

type Chirp struct {
	Body string `json:"body"`
	ID   int    `json:"id"`
}

type User struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	ID       int    `json:"id"`
}

func NewDB(path string) (*DB, error) {
	newDB := DB{
		path: path,
		mux:  &sync.Mutex{},
	}

	err := newDB.ensureDB()
	if err != nil {
		return nil, err
	}

	return &newDB, nil
}

func (db *DB) CreateUser(email string, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1
	user := User{password, email, id}
	dbStructure.Users[id] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUserByID(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	for _, user := range dbStructure.Users {
		if user.ID == id {
			return user, nil
		}
	}
	return User{}, errors.New("no matching record for ID")
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, errors.New("no matching record for email")
}

func (db *DB) UpdateUser(userID int, updatedEmail string, updatedPassword string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	for _, user := range dbStructure.Users {
		if user.ID == userID {
			user.Email = updatedEmail
			user.Password = updatedPassword

			dbStructure.Users[userID] = user
			err = db.writeDB(dbStructure)
			if err != nil {
				return User{}, err
			}
			return user, nil
		}
	}
	return User{}, err
}

func (db *DB) CreateChirp(msg string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{msg, id}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	structure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	var chirps []Chirp
	for _, v := range structure.Chirps {
		chirps = append(chirps, Chirp{v.Body, v.ID})
	}

	return chirps, nil
}

func (db *DB) AddRevokedRefreshToken(tokenString string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	dbStructure.RevokedTokens[tokenString] = time.Now().UTC()
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps:        make(map[int]Chirp),
		Users:         make(map[int]User),
		RevokedTokens: make(map[string]time.Time),
	}
	return db.writeDB(dbStructure)
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	fileBytes, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	dbStructure := DBStructure{
		make(map[int]Chirp),
		make(map[int]User),
		make(map[string]time.Time),
	}

	err = json.Unmarshal(fileBytes, &dbStructure)
	if err != nil {
		return DBStructure{}, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := json.Marshal(dbStructure)
	perm := os.FileMode(0600)

	err = os.WriteFile(db.path, data, perm)
	if err != nil {
		return err
	}
	return nil
}
