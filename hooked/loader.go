package hooked

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"

	"github.com/satori/go.uuid"
)

// Mainly for loading the initial fixtures.

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
}

func getModels() ([]Activity, []Story, []User) {
	dat, err := ioutil.ReadFile("./fixtures/activities.json")
	if err != nil {
		panic(err)
	}
	var activities []Activity
	err = json.Unmarshal(dat, &activities)
	if err != nil {
		panic(err)
	}

	dat, err = ioutil.ReadFile("./fixtures/stories.json")
	var stories []Story
	err = json.Unmarshal(dat, &stories)
	if err != nil {
		panic(err)
	}

	dat, err = ioutil.ReadFile("./fixtures/users.json")
	var users []User
	err = json.Unmarshal(dat, &users)
	if err != nil {
		panic(err)
	}
	return activities, stories, users
}

// LoadFixtures destructively loads fixtures, wiping the slate every time.
func LoadFixtures(db *sql.DB) {
	activities, stories, users := getModels()
	db.Exec("DELETE from notifications")
	db.Exec("DELETE from activities")
	db.Exec("DELETE from stories")
	db.Exec("DELETE from users")

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(`
            INSERT INTO users (sid, firstname, lastname)
            VALUES ($1, $2, $3)
        `)
	for _, user := range users {
		stmt.Exec(user.ID, user.FirstName, user.LastName)
	}
	tx.Commit()

	tx, _ = db.Begin()
	stmt, _ = tx.Prepare(`
            INSERT INTO stories (sid, title, author)
            VALUES ($1, $2, $3)
        `)
	for _, story := range stories {
		stmt.Exec(story.ID, story.Title, story.Author)
	}
	tx.Commit()

	tx, _ = db.Begin()
	stmt, _ = tx.Prepare(`
            INSERT INTO activities (sid, action, date, actor, user2)
            VALUES ($1, $2, $3, $4, $5)
            `)
	for _, activity := range activities {
		stmt.Exec(activity.ID, activity.Action, activity.Date, activity.Actor,
			activity.User2)
	}
	tx.Commit()

	// Now pre-compute notifications table from existing activities,
	// so we can begin querying the API right away.
	preComputeNotifications(db, activities)
}

// Pre-calculate notifications from initial fixtures. Since all activities
// in fixtures are just `follow`, the initial notifications in this table
// will be very similar to the notifications in the activities table.
func preComputeNotifications(db *sql.DB, activities []Activity) {
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(`
        INSERT INTO notifications (id, notified, actor, action, story)
        VALUES ($1, $2, $3, $4, $5)
    `)
	// Add notification to followed user (user2)
	for _, activity := range activities {
		stmt.Exec(uuid.NewV4(),
			activity.User2, activity.Actor, activity.Action, nil)
	}
	tx.Commit()
}
