package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	DocumentCollectionName = "hkg_builder_document"
)

type Document struct {
	ID        string    `bson:"_id"`
	Name      string    `bason:"name"`
	Label     []string  `bason:"label"`
	Text      string    `bason:"text"`
	UpdatedAt time.Time `bason:"updatedAt"`
}

type DocumentDAO struct {
	conn *Conn
}

func NewDocumentDAO(_conn *Conn) *DocumentDAO {
	if nil == _conn {
		return &DocumentDAO{
			conn: defaultConn,
		}
	} else {
		return &DocumentDAO{
			conn: _conn,
		}
	}
}

func (this *DocumentDAO) UpsertOne(_doc *Document) (_err error) {
	_err = nil

	ctx, cancel := NewContext()
	defer cancel()

	filter := bson.D{{"_id", _doc.ID}}
	update := bson.D{
		{"$set", bson.D{
			{"name", _doc.Name},
			{"label", _doc.Label},
			{"text", _doc.Text},
			{"updatedAt", _doc.UpdatedAt},
		}},
	}

	upsert := true
	options := &options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := this.conn.DB.Collection(DocumentCollectionName).UpdateOne(ctx, filter, update, options)
	_err = err
	return
}

func (this *DocumentDAO) Count() (int64, error) {
	ctx, cancel := NewContext()
	defer cancel()
	count, err := this.conn.DB.Collection(DocumentCollectionName).EstimatedDocumentCount(ctx)
	return count, err
}

func (this *DocumentDAO) List(_offset int64, _count int64) ([]*Document, error) {
	ctx, cancel := NewContext()
	defer cancel()

	filter := bson.D{}
	// 1: 倒叙  -1：正序
	sort := bson.D{{"updatedAt", -1}}

	findOptions := options.Find()
	findOptions.SetSort(sort)
	findOptions.SetSkip(_offset)
	findOptions.SetLimit(_count)

	cur, err := this.conn.DB.Collection(DocumentCollectionName).Find(ctx, filter, findOptions)
	if nil != err {
		return make([]*Document, 0), err
	}
	defer cur.Close(ctx)

	var ary []*Document
	for cur.Next(ctx) {
		var document Document
		err = cur.Decode(&document)
		if nil != err {
			return make([]*Document, 0), err
		}
		ary = append(ary, &document)
	}
	return ary, nil
}

func (this *DocumentDAO) UpdateOne(_doc *Document) error {
	ctx, cancel := NewContext()
	defer cancel()

	filter := bson.D{{"_id", _doc.ID}}
	update := bson.D{
		{"$set", bson.D{
			{"text", _doc.Text},
			{"updatedAt", _doc.UpdatedAt},
		}},
	}
	_, err := this.conn.DB.Collection(DocumentCollectionName).UpdateOne(ctx, filter, update)
	if nil != err {
		return err
	}
	return nil
}

func (this *DocumentDAO) FindOne(_id string) (*Document, error) {
	ctx, cancel := NewContext()
	defer cancel()

	filter := bson.D{{"_id", _id}}
	res := this.conn.DB.Collection(DocumentCollectionName).FindOne(ctx, filter)
	if res.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}
	var document Document
	err := res.Decode(&document)
	return &document, err
}

func (this *DocumentDAO) DeleteOne(_id string) (error) {
	ctx, cancel := NewContext()
	defer cancel()

	filter := bson.D{{"_id", _id}}
	_, err := this.conn.DB.Collection(DocumentCollectionName).DeleteOne(ctx, filter)
	return err
}

func (this *DocumentDAO) DeleteMany(_id []string) (error) {
	ctx, cancel := NewContext()
	defer cancel()

    filter := bson.M{"_id": bson.M{"$in": _id}}
	_, err := this.conn.DB.Collection(DocumentCollectionName).DeleteMany(ctx, filter)
	return err
}
