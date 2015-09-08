package main

import (
    "fmt"
    "github.com/lealife/test/db"
    "gopkg.in/mgo.v2/bson"
)


func main() {
    fmt.Println("life");
    db.Init("localhost", "leanote", "", "")

    notes := []map[string]interface{}{}
    db.ListByQ(db.Notes, bson.M{}, &notes);
    fmt.Println(notes)
}