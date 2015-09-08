package db

import (
    "fmt"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "strings"
    "time"
)

// Init mgo and the common DAO

// 数据连接
var Session *mgo.Session

var DBName string;

// 各个表的Collection对象
var Notebooks *mgo.Collection
var Notes *mgo.Collection

// 初始化时连接数据库
func Init(host, dbname, username, password string) {
    dialInfo := &mgo.DialInfo{
        Addrs: strings.Split(host, ","),
        Username: username,
        Password: password,
        Source: dbname,
        Timeout:   5 * time.Second,
    }

    DBName = dbname

    var err error
    // Session, err = mgo.Dial(url)
    fmt.Println("db connectting...");
    Session, err = mgo.DialWithInfo(dialInfo)
    if err != nil {
        panic(err)
    }
    fmt.Println("db connected");

    // Optional. Switch the session to a monotonic behavior.
    Session.SetMode(mgo.Monotonic, true)

    // notebook
    Notebooks = Session.DB(dbname).C("notebooks")

    // notes
    Notes = Session.DB(dbname).C("notes")
}

// common DAO
// 公用方法

//----------------------

func Insert(collection *mgo.Collection, i interface{}) bool {
    err := collection.Insert(i)
    return Err(err)
}

//----------------------

// 适合一条记录全部更新
func Update(collection *mgo.Collection, query interface{}, i interface{}) bool {
    err := collection.Update(query, i)
    return Err(err)
}
func Upsert(collection *mgo.Collection, query interface{}, i interface{}) bool {
    _, err := collection.Upsert(query, i)
    return Err(err)
}
func UpdateAll(collection *mgo.Collection, query interface{}, i interface{}) bool {
    _, err := collection.UpdateAll(query, i)
    return Err(err)
}
func UpdateByIdAndUserId(collection *mgo.Collection, id, userId string, i interface{}) bool {
    err := collection.Update(GetIdAndUserIdQ(id, userId), i)
    return Err(err)
}

func UpdateByIdAndUserId2(collection *mgo.Collection, id, userId bson.ObjectId, i interface{}) bool {
    err := collection.Update(GetIdAndUserIdBsonQ(id, userId), i)
    return Err(err)
}
func UpdateByIdAndUserIdField(collection *mgo.Collection, id, userId, field string, value interface{}) bool {
    return UpdateByIdAndUserId(collection, id, userId, bson.M{"$set": bson.M{field: value}})
}
func UpdateByIdAndUserIdMap(collection *mgo.Collection, id, userId string, v bson.M) bool {
    return UpdateByIdAndUserId(collection, id, userId, bson.M{"$set": v})
}

func UpdateByIdAndUserIdField2(collection *mgo.Collection, id, userId bson.ObjectId, field string, value interface{}) bool {
    return UpdateByIdAndUserId2(collection, id, userId, bson.M{"$set": bson.M{field: value}})
}
func UpdateByIdAndUserIdMap2(collection *mgo.Collection, id, userId bson.ObjectId, v bson.M) bool {
    return UpdateByIdAndUserId2(collection, id, userId, bson.M{"$set": v})
}

//
func UpdateByQField(collection *mgo.Collection, q interface{}, field string, value interface{}) bool {
    _, err := collection.UpdateAll(q, bson.M{"$set": bson.M{field: value}})
    return Err(err)
}
func UpdateByQI(collection *mgo.Collection, q interface{}, v interface{}) bool {
    _, err := collection.UpdateAll(q, bson.M{"$set": v})
    return Err(err)
}

// 查询条件和值
func UpdateByQMap(collection *mgo.Collection, q interface{}, v interface{}) bool {
    _, err := collection.UpdateAll(q, bson.M{"$set": v})
    return Err(err)
}

//------------------------

// 删除一条
func Delete(collection *mgo.Collection, q interface{}) bool {
    err := collection.Remove(q)
    return Err(err)
}
func DeleteByIdAndUserId(collection *mgo.Collection, id, userId string) bool {
    err := collection.Remove(GetIdAndUserIdQ(id, userId))
    return Err(err)
}
func DeleteByIdAndUserId2(collection *mgo.Collection, id, userId bson.ObjectId) bool {
    err := collection.Remove(GetIdAndUserIdBsonQ(id, userId))
    return Err(err)
}

// 删除所有
func DeleteAllByIdAndUserId(collection *mgo.Collection, id, userId string) bool {
    _, err := collection.RemoveAll(GetIdAndUserIdQ(id, userId))
    return Err(err)
}
func DeleteAllByIdAndUserId2(collection *mgo.Collection, id, userId bson.ObjectId) bool {
    _, err := collection.RemoveAll(GetIdAndUserIdBsonQ(id, userId))
    return Err(err)
}

func DeleteAll(collection *mgo.Collection, q interface{}) bool {
    _, err := collection.RemoveAll(q)
    return Err(err)
}

//-------------------------

func Get(collection *mgo.Collection, id string, i interface{}) {
    // CheckMongoSessionLostErr()
    collection.FindId(bson.ObjectIdHex(id)).One(i)
}
func Get2(collection *mgo.Collection, id bson.ObjectId, i interface{}) {
    // CheckMongoSessionLostErr()
    collection.FindId(id).One(i)
}

func GetByQ(collection *mgo.Collection, q interface{}, i interface{}) {
    // CheckMongoSessionLostErr()
    collection.Find(q).One(i)
}
func ListByQ(collection *mgo.Collection, q interface{}, i interface{}) {
    // CheckMongoSessionLostErr()
    collection.Find(q).All(i)
}

func ListByQLimit(collection *mgo.Collection, q interface{}, i interface{}, limit int) {
    // CheckMongoSessionLostErr()
    collection.Find(q).Limit(limit).All(i)
}

// 查询某些字段, q是查询条件, fields是字段名列表
func GetByQWithFields(collection *mgo.Collection, q bson.M, fields []string, i interface{}) {
    // CheckMongoSessionLost()
    selector := make(bson.M, len(fields))
    for _, field := range fields {
        selector[field] = true
    }
    collection.Find(q).Select(selector).One(i)
}

// 查询某些字段, q是查询条件, fields是字段名列表
func ListByQWithFields(collection *mgo.Collection, q bson.M, fields []string, i interface{}) {
    selector := make(bson.M, len(fields))
    for _, field := range fields {
        selector[field] = true
    }
    collection.Find(q).Select(selector).All(i)
}
func GetByIdAndUserId(collection *mgo.Collection, id, userId string, i interface{}) {
    collection.Find(GetIdAndUserIdQ(id, userId)).One(i)
}
func GetByIdAndUserId2(collection *mgo.Collection, id, userId bson.ObjectId, i interface{}) {
    collection.Find(GetIdAndUserIdBsonQ(id, userId)).One(i)
}

// 按field去重
func Distinct(collection *mgo.Collection, q bson.M, field string, i interface{}) {
    collection.Find(q).Distinct(field, i)
}

//----------------------

func Count(collection *mgo.Collection, q interface{}) int {
    cnt, err := collection.Find(q).Count()
    if err != nil {
        Err(err)
    }
    return cnt
}

func Has(collection *mgo.Collection, q interface{}) bool {
    if Count(collection, q) > 0 {
        return true
    }
    return false
}

//-----------------

// 得到主键和userId的复合查询条件
func GetIdAndUserIdQ(id, userId string) bson.M {
    return bson.M{"_id": bson.ObjectIdHex(id), "UserId": bson.ObjectIdHex(userId)}
}
func GetIdAndUserIdBsonQ(id, userId bson.ObjectId) bson.M {
    return bson.M{"_id": id, "UserId": userId}
}

// DB处理错误
func Err(err error) bool {
    if err != nil {
        fmt.Println(err)
        // 删除时, 查找
        if err.Error() == "not found" {
            return true
        }
        return false
    }
    return true
}
