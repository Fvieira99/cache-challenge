package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Cache struct {
	Users map[int]*User
}

func (c *Cache) Read(id int) *User {
	return c.Users[id]
}

func (c *Cache) Write(id int, u *User) interface{} {
	if len(c.Users) == 100 {
		fmt.Println("Cache has already reached it`s limit of 100 users")
		return nil
	}
	c.Users[id] = u
	return nil
}

type Db struct {
	Users   map[int]*User
	Queries int
}

func (db *Db) Seed() {
	for i := 0; i <= 100; i++ {
		db.Users[i] = &User{
			Id:       i,
			Username: fmt.Sprintf("User%d", i),
		}
	}
}

func (db *Db) FindById(id int) *User {
	return db.Users[id]
}

type User struct {
	Id       int
	Username string
}

type Server struct {
	*Db
	*Cache
}

func NewServer() *Server {
	db := &Db{
		Users:   make(map[int]*User),
		Queries: 0,
	}
	cache := &Cache{
		Users: make(map[int]*User),
	}

	db.Seed()
	return &Server{
		Db:    db,
		Cache: cache,
	}
}

func (s *Server) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	strId := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(strId)

	cacheUser := s.Cache.Read(id)

	if cacheUser != nil {
		json.NewEncoder(w).Encode(cacheUser)
		return
	}

	dbUser := s.Db.FindById(id)
	s.Db.Queries++

	if dbUser == nil {
		panic("User Not Found")
	}

	s.Cache.Write(id, dbUser)
	json.NewEncoder(w).Encode(dbUser)
}

func main() {
	// s := NewServer()

}
