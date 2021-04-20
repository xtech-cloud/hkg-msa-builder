package model

import (
    "fmt"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"hkg-msa-builder/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	uuid "github.com/satori/go.uuid"
    "github.com/micro/go-micro/v2/logger"
)

var base64Coder = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")

type Conn struct {
    DB *mongo.Database
    client *mongo.Client
    context context.Context
}
var defaultConn *Conn

func Setup() {
    // 设置客户端参数
    uri := fmt.Sprintf("mongodb://%s", config.Schema.Database.MongoDB.Address)
    cliOptions := options.Client().ApplyURI(uri)

    // 设置连接时长
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Schema.Database.MongoDB.Timeout) * time.Second)
    defer cancel()

    // 连接mongodb
    logger.Infof("connect %s ......", uri)
    client, err := mongo.Connect(ctx, cliOptions)
    if nil != err {
        logger.Fatal(err)
    }

    err = client.Ping(ctx, readpref.Primary())
    if nil != err {
        logger.Fatal(err)
    }

    logger.Infof("connect %s success", uri)
    defaultConn = &Conn{
        client: client,
        context: ctx,
        DB: client.Database(config.Schema.Database.MongoDB.DB),
    }
}

func Cancel() {
    if nil != defaultConn {
        defaultConn.client.Disconnect(defaultConn.context)
    }
}

func NewUUID() string {
	guid := uuid.NewV4()
	h := md5.New()
	h.Write(guid.Bytes())
	return hex.EncodeToString(h.Sum(nil))
}

func ToUUID(_content string) string {
	h := md5.New()
	h.Write([]byte(_content))
	return hex.EncodeToString(h.Sum(nil))
}

func MD5(_content string) string {
	h := md5.New()
	h.Write([]byte(_content))
	return hex.EncodeToString(h.Sum(nil))
}

func ToBase64(_content []byte) string {
	return base64Coder.EncodeToString(_content)
}

func NewContext() (context.Context, context.CancelFunc) {
    return context.WithTimeout(context.TODO(), time.Duration(config.Schema.Database.MongoDB.Timeout) * time.Second)
    //return context.WithTimeout(context.Background(), time.Duration(config.Schema.Database.MongoDB.Timeout) * time.Second)
}

