package hooked

import (
	"database/sql"
	"errors"
)

const (
	ActionFollow  = "follow"
	ActionLove    = "love"
	ActionRead    = "read"
	ActionWrite   = "write"
	ActionComment = "comment"
)

type User struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	ID        string `json:"_id"`
}

type Story struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	ID     string `json:"_id"`
}

type Activity struct {
	Action string `json:"action"`
	Date   string `json:"date"`
	Actor  string `json:"actor"`
	User2  string `json:"user2"`
	ID     string `json:"_id"`
	Story  string `json:"story"`
}

func (a *Activity) Validate(db *sql.DB) error {

	switch a.Action {
	case ActionFollow, ActionLove, ActionRead, ActionWrite, ActionComment:
	default:
		return errors.New(
			"Must provide a supported action: follow, love, read, write, comment.")
	}

	if a.Actor == "" {
		return errors.New("Must provide an actor.")
	}

	_, err := getUser(db, a.Actor)
	if err != nil {
		return err
	}
	if a.Action == ActionFollow && a.User2 == "" {
		return errors.New("Must provide a user to follow.")
	}
	if a.User2 != "" {
		_, err = getUser(db, a.User2)
		if err != nil {
			return err
		}
	}
	if a.Story != "" {
		_, err := getStory(db, a.Story)
		if err != nil {
			return err
		}
	}
	return nil
}

func getUser(db *sql.DB, id string) (*User, error) {
	var firstname string
	var lastname string
	err := db.QueryRow(
		"SELECT firstname, lastname FROM users WHERE sid = $1", id).Scan(
		&firstname, &lastname)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName: firstname,
		LastName:  lastname,
		ID:        id,
	}, nil
}

func getStory(db *sql.DB, id string) (*Story, error) {
	var title string
	var author string
	err := db.QueryRow(
		"SELECT title, author FROM users WHERE sid = $1", id).Scan(
		&title, &author)
	if err != nil {
		return nil, err
	}
	return &Story{
		Title:  title,
		Author: author, // This is still a User ID.
		ID:     id,
	}, nil
}
