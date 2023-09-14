package database

import (
	"encoding/json"
	"errors"
	"math/rand"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.Mutex
}
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
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

func (db *DB) CreateChirp(msg string) (Chirp, error) {
	perm := os.FileMode(0777)
	randNum := rand.Intn(1000)
	chirp := Chirp{randNum, msg}

	data, err := json.Marshal(chirp)
	if err != nil {
		return Chirp{}, err
	}
	err = os.WriteFile(db.path, []byte(data), perm)
	if err != nil {
		return Chirp{}, errors.New("Error creating Chirp")
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
		chirps = append(chirps, Chirp{v.ID, v.Body})
	}

	return chirps, nil
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{Chirps: make(map[int]Chirp)}
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

	dbStructure := DBStructure{}

	err = json.Unmarshal(fileBytes, &dbStructure)
	if err != nil {
		return DBStructure{}, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Lock()

	data, err := json.Marshal(dbStructure)
	perm := os.FileMode(0600)

	err = os.WriteFile(db.path, data, perm)
	if err != nil {
		return err
	}
	return nil
}
