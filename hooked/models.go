package hooked

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/satori/go.uuid"
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

// Generate a 24-character ID. The fixtures use 24-character IDs so
// let's truncate a UUID for now and hope that's enough.
func genID() string {
	str := uuid.NewV4().String()
	str = strings.Replace(str, "-", "", -1)
	return str[:24]
}

const HookedRFC = "2006-01-02T15:04:05.000Z07:00"

// Return now as a RFC-formatted time string in Zulu.
func now() string {
	return time.Now().Format(HookedRFC)
}

// Validate validates the activity, checking for various heuristics,
// prior to saving it to the database.
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

func saveActivity(db *sql.DB, a *Activity) error {
	// Go is a bit verbose
	var story, user2 interface{}
	if a.Story == "" {
		story = nil
	} else {
		story = a.Story
	}
	if a.User2 == "" {
		user2 = nil
	} else {
		user2 = a.User2
	}
	//
	_, err := db.Query(`
        INSERT into activities (sid, action, date, actor_id, user2_id, story_id)
        VALUES($1, $2, $3, $4, $5, $6)
    `, genID(), a.Action, now(), a.Actor, user2, story)
	return err
}

func createNotifications(db *sql.DB, a *Activity) error {
	// Generate the notifications.
	/*
	   - user follows another user
	       - send push notification to followed user
	       - add notification to followed user
	   - user reads a story
	       - send push notification to story author
	       - add notification to actor’s followers
	   - user loves a story
	       - send push notification to story author
	       - add notification to actor’s followers
	   - user writes a story
	       - send push notification to all followers
	       - add notification to all followers
	   - user comments on a story
	       - send push notification to story’s author
	       - add notification to actor’s followers
	*/
	var err error
	switch a.Action {
	case ActionFollow:
		// Add notification to followed user.
		_, err = db.Query(`
	           INSERT into notifications (id, notified_id, actor_id, action, date)
               VALUES ($1, $2, $3, $4, $5)
        `, uuid.NewV4(), a.User2, a.Actor, a.Action, now())
		if err != nil {
			return err
		}
	case ActionRead, ActionLove /* 😍 */, ActionWrite, ActionComment:
		// Add notification to actor's followers.
		followers, err := getFollowerIDs(db, a.Actor)
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] Action=%v, adding %v notifications to followers",
			a.Action, len(followers))
		for _, followerID := range followers {
			_, err = db.Query(`
                   INSERT into notifications
                   (id, notified_id, actor_id, action, date, story_id)
                   VALUES ($1, $2, $3, $4, $5, $6)
            `, uuid.NewV4(), followerID, a.Actor, a.Action, now(), a.Story)
			if err != nil {
				return err
			}
		}

	}
	return err
}

// Save saves the activity to the database, and it also creates the
// notification object(s).
func (a *Activity) Save(db *sql.DB) error {
	// First, save the activity to the database.
	tx, _ := db.Begin()

	err := saveActivity(db, a)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = createNotifications(db, a)
	if err != nil {
		tx.Rollback()
		return err
	}
	// Save as a transaction.
	err = tx.Commit()
	return err
}

func (a *Activity) Notify() {
	// This goroutine, and the fact that the API in general uses goroutines
	// for HTTP requests, allows the API to scale more easily. We don't block
	// until all push notifications are delivered, instead we hand them off
	// in a goroutine and exit. This could be a separate microservice
	// or job queue later on.
	go func() {

	}()
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
		"SELECT title, author_id FROM stories WHERE sid = $1", id).Scan(
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

func getFollowerIDs(db *sql.DB, id string) ([]string, error) {
	ids := []string{}
	rows, err := db.Query(
		"SELECT follower_id FROM followers WHERE user_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
