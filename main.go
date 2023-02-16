package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionName = "wukongpicture"
)

// 1. get data
func getWuKongPictures() (m map[string]dWuKongPicture) {
	var v []dWuKongPicture
	m = map[string]dWuKongPicture{}
	f := func(ctx context.Context) error {
		return cli.getDocs(
			ctx, collectionName, nil,
			nil, &v,
		)
	}

	if err := withContext(f); err != nil {
		return
	}

	for _, p := range v {
		m[p.Owner] = p
	}

	return
}

// 2. create likes and intert array data
func createLikes(owner string, items []pictureItem, version int) (err error) {
	var likes []primitive.M
	for _, item := range items {
		like, _ := toLikePictureItemDoc(owner, &item)
		likes = append(likes, like)
	}

	f := func(ctx context.Context) error {
		return cli.updateDoc(
			ctx, collectionName,
			wukongOwnerFilter(owner),
			bson.M{fieldLikes: likes}, mongoCmdSet, version,
		)
	}

	return withContext(f)
}

// 3. create publics without data
func createPublics(owner string, version int) (err error) {
	f := func(ctx context.Context) error {
		return cli.updateDoc(
			ctx, collectionName,
			wukongOwnerFilter(owner),
			bson.M{fieldPublics: bson.A{}}, mongoCmdSet, version,
		)
	}

	return withContext(f)
}

func main() {
	// 修改下方的conn和dbName
	err := Initialize("mongodb://xiheAdmin:Password11@127.0.0.1:27017/xihe-server-test?timeoutMS=5000", "xihe-server-test")
	if err != nil {
		panic(err)
	}

	m := getWuKongPictures()

	for k, v := range m {
		createLikes(k, v.Items, v.Version)
		createPublics(k, v.Version+1)
	}
}

func toLikePictureItemDoc(owner string, item *pictureItem) (primitive.M, error) {
	like := pictureItem{
		Id:        item.Id,
		Owner:     owner,
		Desc:      item.Desc,
		Style:     item.Style,
		OBSPath:   item.OBSPath,
		Diggs:     []string{},
		DiggCount: 0,
		CreatedAt: item.CreatedAt,
	}
	return genDoc(like)
}
