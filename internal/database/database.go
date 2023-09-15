package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.Mutex
}
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

type Chirp struct {
	Body string `json:"body"`
	ID   int    `json:"id"`
}

type User struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
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

func (db *DB) CreateUser(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1
	user := User{email, id}
	dbStructure.Users[id] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
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

func (db *DB) createDB() error {
	dbStructure := DBStructure{Chirps: make(map[int]Chirp), Users: make(map[int]User)}
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

	dbStructure := DBStructure{make(map[int]Chirp), make(map[int]User)}

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
