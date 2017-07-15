package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func loadFixtures() {
	_, _, users := getModels()
	fmt.Println(users)
}
