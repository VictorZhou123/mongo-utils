package main

import (
	"context"
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

// 2. modify data of likes
func modifyLikes(owner string, item pictureItem) (err error) {
	item = pictureItem{
		Id:        item.Id,
		Owner:     owner,
		Desc:      item.Desc,
		Style:     item.Style,
		OBSPath:   item.OBSPath,
		Diggs:     []string{},
		DiggCount: 0,
		CreatedAt: item.CreatedAt,
	}

	doc, _ := genDoc(item)
	doc[fieldVersion] = 1

	f := func(ctx context.Context) error {
		_, err := cli.modifyArrayElemWithoutVersion(
			ctx, collectionName, fieldLikes,
			wukongOwnerFilter(owner), wukongIdFilter(item.Id),
			doc, mongoCmdSet,
		)

		return err
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
		for _, like := range v.Likes {
			modifyLikes(k, like)
		}
	}
}
