package handler

import (
	"context"
	"encoding/json"
	"hkg-msa-builder/model"
	"time"

	"github.com/micro/go-micro/v2/logger"

	proto "github.com/xtech-cloud/hkg-msp-builder/proto/builder"
)

type Document struct{}

type Pattern struct {
	From []string `json:"from"`
	To   string   `json:"to"`
}

func (this *Document) Merge(_ctx context.Context, _req *proto.DocumentMergeRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Document.Merge, req is %v,%v", _req.Name, _req.Label)

	_rsp.Status = &proto.Status{}

	if "" == _req.Name {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "name is required"
		return nil
	}

	if 0 == len(_req.Label) {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "label is required"
		return nil
	}

	var format []Pattern
	err := json.Unmarshal([]byte(_req.Format), &format)
	if nil != err {
		_rsp.Status.Code = 2
		_rsp.Status.Message = err.Error()
		return nil
	}

	output := make(map[string]interface{})
	for _, text := range _req.Text {
		var obj map[string]interface{}
		err = json.Unmarshal([]byte(text), &obj)
		for k, v := range obj {
			logger.Trace(k)
			logger.Trace(v)
		}
	}

	textOutput, err := json.Marshal(output)
	if nil != err {
		_rsp.Status.Code = 4
		_rsp.Status.Message = err.Error()
		return nil
	}

	uuid := _req.Name
	for _, label := range _req.Label {
		uuid += label
	}
	document := &model.Document{
		ID:        model.ToUUID(uuid),
		Name:      _req.Name,
		Label:     _req.Label,
		Text:      string(textOutput),
		UpdatedAt: time.Now(),
	}
	dao := model.NewDocumentDAO(nil)
	err = dao.UpsertOne(document)
	return err
}

func (this *Document) List(_ctx context.Context, _req *proto.ListRequest, _rsp *proto.DocumentListResponse) error {
	logger.Infof("Received Document.List, req is %v", _req)

	_rsp.Status = &proto.Status{}
	offset := int64(0)
	if _req.Offset > 0 {
		offset = _req.Offset
	}

	count := int64(50)
	if _req.Count > 0 {
		count = _req.Count
	}

	dao := model.NewDocumentDAO(nil)
	total, err := dao.Count()
	if nil != err {
		return err
	}
	_rsp.Total = total

	ary, err := dao.List(offset, count)
	if nil != err {
		return err
	}

	_rsp.Entity = make([]*proto.DocumentEntity, len(ary))
	for i, v := range ary {
		_rsp.Entity[i] = &proto.DocumentEntity{
			Uuid: v.ID,
			Name: v.Name,
		}
	}
	return nil
}
