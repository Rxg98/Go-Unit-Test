package dao

import (
	"context"
	"fmt"
	"os"
	"testing"

	mgo "coolcar/shared/mongo"
	mongotesting "coolcar/shared/mongo/testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoURI string

func TestResolveAccountID(t *testing.T) {
	c := context.Background()
	mc, err := mongo.Connect(c, options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}

	m := NewMongo(mc.Database("coolcar"))
	fmt.Println("start insert")
	_, err = m.col.InsertMany(c, []interface{}{
		bson.M{
			mgo.IDField: mustObjID("5f7c245ab0361e00ffb9fd6f"),
			openIDField: "openid_1",
		},
		bson.M{
			mgo.IDField: mustObjID("5f7c245ab0361e00ffb9fd70"),
			openIDField: "openid_2",
		},
	})
	if err != nil {
		t.Fatalf("cannot insert initial values: %v", err)
	}
	fmt.Println("insert success")
	m.NewObjID = func() primitive.ObjectID {
		return mustObjID("5f7c245ab0361e00ffb9fd71")
	}
	cases := []struct {
		name   string
		openID string
		want   string
	}{
		{
			name:   "existing_user",
			openID: "openid_1",
			want:   "5f7c245ab0361e00ffb9fd6f",
		},
		{
			name:   "another_existing_user",
			openID: "openid_2",
			want:   "5f7c245ab0361e00ffb9fd70",
		},
		{
			name:   "new_user",
			openID: "openid_3",
			want:   "5f7c245ab0361e00ffb9fd71",
		},
	}
	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			id, err := m.ResolveAccountID(context.Background(), cc.openID)
			if err != nil {
				t.Errorf("cannot resolve account id for %q: %v", cc.openID, err)
			}
			if id != cc.want {
				t.Errorf("resolve account id: want: %q; got: %q", cc.want, id)
			}
		})
	}
	id, err := m.ResolveAccountID(c, "openid_3")
	if err != nil {
		t.Errorf("faild resolve account id for openid_3: %v", err)
	} else {
		want := "5f7c245ab0361e00ffb9fd71"
		if id != want {
			t.Errorf("resolve account id: want: %q; got: %q", want, id)
		}
	}
}
func mustObjID(hex string) primitive.ObjectID {
	objID, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		panic(err)
	}
	return objID
}
func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m, &mongoURI))
}
