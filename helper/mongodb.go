package helper

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"strings"
)

type Mongo struct {
	session     *mgo.Session
	db          string
	collections map[string]*mgo.Collection
}

func NewMongo(dns, mode string, refresh bool) *Mongo {
	pos := strings.LastIndex(dns, "/")
	db := dns[pos+1:]
	if dbPos := strings.LastIndex(db, "?"); dbPos > 0 {
		db = db[0:dbPos]
	}
	session, err := mgo.Dial(dns)
	if err != nil {
		panic("collection mongodb error:" + err.Error())
	}
	return &Mongo{
		session:     session.Copy(),
		db:          db,
		collections: make(map[string]*mgo.Collection),
	}
}

func (m *Mongo) C(collection bson.M) *mgo.Collection {
	colName := collection["name"].(string)
	if _, ok := m.collections[colName]; !ok {
		m.collections[colName] = m.session.DB(m.db).C(colName)
		if index, ok := collection["index"]; ok {
			if indexs, ok := index.([]string); ok {
				for _, v := range indexs {
					indexSlice := strings.Split(v, ",")
					err := m.collections[colName].EnsureIndex(mgo.Index{Key: indexSlice, Unique: false, DropDups: false})
					if err != nil {
						return nil
					}
				}
			}
		}
	}

	if iUnique, ok := collection["unique"]; ok {
		if uniques, ok := iUnique.([]string); ok {
			for _, v := range uniques {
				uniqueSlice := strings.Split(v, ",")
				err := m.collections[colName].EnsureIndex(mgo.Index{Key: uniqueSlice, Unique: true})
				if err != nil {
					return nil
				}
			}

		}
	}
	return m.collections[colName]
}

func (m *Mongo) Close() {
	m.session.Close()
}

func setMode(session *mgo.Session, mode string, refresh bool) {
	switch strings.ToLower(mode) {
	case "eventual":
		session.SetMode(mgo.Eventual, refresh)
	case "monotoic":
		session.SetMode(mgo.Monotonic, refresh)
	case "strong":
		session.SetMode(mgo.Strong, refresh)
	}
}
