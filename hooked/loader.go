package hooked

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"

	"github.com/satori/go.uuid"
)

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
	db.Exec("DELETE from followers")
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
            INSERT INTO stories (sid, title, author_id)
            VALUES ($1, $2, $3)
        `)
	for _, story := range stories {
		stmt.Exec(story.ID, story.Title, story.Author)
	}
	tx.Commit()

	tx, _ = db.Begin()
	stmt, _ = tx.Prepare(`
            INSERT INTO activities (sid, action, date, actor_id, user2_id)
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
	preComputeFollowers(db, activities)
}

// Pre-calculate notifications from initial fixtures. Since all activities
// in fixtures are just `follow`, the initial notifications in this table
// will be very similar to the notifications in the activities table.
func preComputeNotifications(db *sql.DB, activities []Activity) {
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(`
        INSERT INTO notifications (id, notified_id, actor_id, action, story_id, date)
        VALUES ($1, $2, $3, $4, $5, $6)
    `)
	// Add notification to followed user (user2)
	for _, activity := range activities {
		stmt.Exec(uuid.NewV4(),
			activity.User2, activity.Actor, activity.Action, nil, activity.Date)
	}
	tx.Commit()
}

func preComputeFollowers(db *sql.DB, activities []Activity) {
	tx, _ := db.Begin()
	// The test fixtures have at least one duplicate follow event,
	// so add the DO NOTHING below :)
	stmt, _ := tx.Prepare(`
        INSERT INTO followers (user_id, follower_id)
        VALUES ($1, $2)
        ON CONFLICT DO NOTHING
    `)

	// Add notification to followed user (user2)
	for _, activity := range activities {
		stmt.Exec(activity.User2, activity.Actor)
	}
	tx.Commit()
}
