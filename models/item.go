package models

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/Dparty/common/utils"
	abstract "github.com/Dparty/dao/abstract"
	"github.com/Dparty/dao/restaurant"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func NewItem(entity restaurant.Item) Item {
	return Item{entity: entity}
}

type Item struct {
	entity restaurant.Item
}

func (i *Item) Save() *Item {
	itemRepository.Save(&i.entity)
	return i
}

func (i *Item) SetImage(url string) *Item {
	i.entity.Images = []string{url}
	return i
}

func (i Item) ID() uint {
	return i.entity.ID()
}

func (i Item) Name() string {
	return i.entity.Name
}

func (i Item) Entity() restaurant.Item {
	return i.entity
}

func (i *Item) Update(name string, pricing int64, attributes restaurant.Attributes, images, tags []string, printers []uint) (*Item, error) {
	return i, nil
}

func (i Item) Delete() bool {
	ctx := itemRepository.Delete(&i.entity)
	return ctx.RowsAffected != 0
}

func (i Item) Owner() abstract.Owner {
	return i.entity.Owner()
}

func (i *Item) UploadImage(file *multipart.FileHeader) string {
	imageId := utils.GenerteId()
	path := "items/" + utils.UintToString(imageId)
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", Bucket, CosClient.Region))
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  CosClient.SecretID,
			SecretKey: CosClient.SecretKey,
		},
	})
	f, _ := file.Open()
	client.Object.Put(context.Background(), path, f,
		&cos.ObjectPutOptions{
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				ContentType: file.Header.Get("content-type"),
			},
		})
	url := fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", Bucket, CosClient.Region, path)
	i.SetImage(url)
	i.Save()
	return url
}
