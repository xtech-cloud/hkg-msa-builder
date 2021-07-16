package handler

import (
	"context"
	"encoding/base64"
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

func (this *Document) Merge(_ctx context.Context, _req *proto.DocumentMergeRequest, _rsp *proto.DocumentMergeResponse) error {
	logger.Infof("Received Document.Merge, req is %v, %v", _req.Name, _req.Label)

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

	//TODO 使用单独的参数传递融合策略（权重、AI等）

	var patterns []Pattern
	err := json.Unmarshal([]byte(_req.Format), &patterns)
	if nil != err {
		_rsp.Status.Code = 2
		_rsp.Status.Message = err.Error()
		return nil
	}

	patternMap := make(map[string]string)
	for _, v := range patterns {
		for _, from := range v.From {
			patternMap[from] = v.To
		}
	}

	output := make(map[string]interface{})

	//TODO 将合并过程单独记录
	for _, text := range _req.Text {
        bytesText , err := base64.StdEncoding.DecodeString(text)
		if err != nil {
			logger.Error(err)
            continue
		}
		var obj map[string]interface{}
		err = json.Unmarshal(bytesText, &obj)
		if err != nil {
			logger.Error(err)
            continue
		}
		for k, v := range obj {
			this.parse(k, v, output, patternMap)
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
	_rsp.Uuid = document.ID
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
			Uuid:      v.ID,
			Name:      v.Name,
			Label:     v.Label,
			Text:      v.Text,
			UpdatedAt: v.UpdatedAt.Unix(),
		}
	}
	return nil
}

func (this *Document) Delete(_ctx context.Context, _req *proto.DocumentDeleteRequest, _rsp *proto.DocumentDeleteResponse) error {
	logger.Infof("Received Document.Delete, req is %v", _req)

	_rsp.Status = &proto.Status{}
	_rsp.Uuid = _req.Uuid
	dao := model.NewDocumentDAO(nil)
	return dao.DeleteOne(_req.Uuid)
}

func (this *Document) BatchDelete(_ctx context.Context, _req *proto.DocumentBatchDeleteRequest, _rsp *proto.DocumentBatchDeleteResponse) error {
	logger.Infof("Received Document.BatchDelete, req is %v", _req)

	_rsp.Status = &proto.Status{}
	_rsp.Uuid = _req.Uuid
	dao := model.NewDocumentDAO(nil)
	return dao.DeleteMany(_req.Uuid)
}

func (this *Document) parse(_k string, _v interface{}, _output map[string]interface{}, _patterns map[string]string) {
	key := _k
	if k, hasKey := _patterns[_k]; hasKey {
		key = k
	}
	switch _v.(type) {
	case map[string]interface{}:
		mapV, ok := _v.(map[string]interface{})
		if ok {
			for k, v := range mapV {
				this.parse(k, v, _output, _patterns)
			}
		}
	case []interface{}:
		aryV, ok := _v.([]interface{})
		if ok {
			var outAry []interface{}
			if _, found := _output[key]; !found {
				// 不存在则创建一个新数组
				outAry = make([]interface{}, 0)
			} else {
				// 存在则转换
				ary, ok := _output[key].([]interface{})
				if ok {
					outAry = ary
				}
			}
			if nil != outAry {
				outAry = append(outAry, aryV...)
				_output[key] = outAry
			}
		}
	case string:
		strV, ok := _v.(string)
		if ok {
			if _, found := _output[key]; found {
				str, ok := _output[key].(string)
				if ok {
					// 合并值不一样, 则追加到冲突列表中
					if str != strV {
						//TODO 处理冲突列表
					}
				}
			} else {
				_output[key] = strV
			}
		}
	default:
		logger.Warn("unknown type")
	}
}
