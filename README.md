# cool-api

### Stack

- Go 1.8 and PostgreSQL. 
- Docker 

### Installation

- `git clone` this repository
- Install Docker for your operating system from here:

https://www.docker.com/docker-mac

or

https://www.docker.com/docker-windows

The free community edition should be fine

- Once docker is installed, `cd` to _this_ directory and run

`docker-compose up -d`

This should download the necessary docker containers and run the app within a couple of minutes. I inserted a 5-second delay in startup in the `docker-compose.yml` file, see the `sleep 5`, so that the DB starts up fine the very first time. It shouldn't be necessary after that.

- To see logs, do  `docker-compose logs -f api`. The logs log the push notifications as well as other events.
- On startup, we always clear the tables and reload the fixtures to start with an empty slate. It takes about 5-7 seconds to load the fixtures on my laptop.
- To restart the api, `docker-compose restart api`
- To turn it all off, `docker-compose stop`

### Test cases

Some test cases:

```
curl -X POST http://localhost:8086/activity -d '{"action": "read", "actor": "5952930ecc35c8923cca380b", "story": "595294f7fa7c74cd2fe4c33c"}'

Should create 7 notifications (user appears as followed 8 times in the activities.json, but one is a duplicate)
```

```
curl -X POST http://localhost:8086/activity -d '{"action": "follow", "actor": "5952930e4d5ffaf83c757e3d", "user2": "5952930ecc35c8923cca380b"}'

Now there's someone else following the user above

```

```
curl -X POST http://localhost:8086/activity -d '{"action": "love", "actor": "5952930ecc35c8923cca380b", "story": "595294f753e68032ca1feb71"}'

Should now send 8 notifications

```

```
curl -X GET http://localhost:8086/user/5952930e4d5ffaf83c757e3d/notifications

Follower gets the "love" notification above

```


### Future improvements

Improvements (not yet implemented):

- Should make a table for each action type
- Automated tests
- Comment should probably go to a comment table, and accept an additional comment ID, etc.
