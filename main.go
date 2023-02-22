package main

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionName   = "wukongpicture"
	collectionWuKong = "wukong"
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

	var updateCmd primitive.M
	if len(likes) > 0 {
		updateCmd = bson.M{fieldLikes: likes}
	} else {
		updateCmd = bson.M{fieldLikes: bson.A{}}
	}

	f := func(ctx context.Context) error {
		return cli.updateDoc(
			ctx, collectionName,
			wukongOwnerFilter(owner),
			updateCmd, mongoCmdSet, version,
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

// 3. get raw official publics
func getOfficialPublics() (v []pictureInfo, err error) {
	wukong := new(dWuKong)
	f := func(ctx context.Context) error {
		return cli.getDoc(
			ctx, collectionWuKong,
			bson.M{fieldId: "1"},
			bson.M{fieldPictures: 1},
			&wukong,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	v = wukong.Pictures

	return
}

// 4. insert official publics into MindSpore publics
func insertIntoPublics(owner string, v []pictureInfo) (err error) {

	publics := make([]primitive.M, len(v))

	p := toPublicPictureItemDoc(owner, v)

	for i := range v {
		doc, _ := genDoc(p[i])
		publics[i] = doc
	}

	for _, publicDoc := range publics {
		err = pushDocIntoPublics(owner, publicDoc)
		if err != nil {
			return
		}
	}

	return
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

	v, err := getOfficialPublics()
	if err != nil {
		panic(err)
	}

	err = insertIntoPublics("MindSpore", v)
	if err != nil {
		panic(err)
	}
}

func pushDocIntoPublics(owner string, doc primitive.M) error {
	f := func(ctx context.Context) error {
		return cli.pushArrayElem(
			ctx, collectionName, fieldPublics, wukongOwnerFilter(owner), doc,
		)
	}

	if err := withContext(f); err != nil {
		return err
	}

	return nil
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

func toPublicPictureItemDoc(owner string, items []pictureInfo) (p []pictureItem) {
	p = make([]pictureItem, len(items))
	for i := range items {
		p[i].Id = newId()
		p[i].Owner = owner
		p[i].Desc = items[i].Desc
		p[i].Style = items[i].Style
		p[i].OBSPath = toOBSPath(items[i].Link)
		p[i].Level = 2
		p[i].Diggs = []string{}
		p[i].DiggCount = 0
		p[i].CreatedAt = "2023-02-22"
		p[i].Version = 1
	}

	return
}

func toOBSPath(link string) (obspath string) {
	parts1 := strings.Split(link, "https://big-model-deploy.obs.cn-central-221.ovaijisuan.com/")[1]
	partsArr := strings.Split(parts1, "AI-gallery/gallery/")
	s := partsArr[0] + "AI-gallery/gallery/" + "MindSpore/1677053755/" + partsArr[1]
	return s
}
