package app

import (
	"context"
	"github.com/glennliao/apijson-go"
	"github.com/glennliao/apijson-go/config"
	"github.com/yitter/idgenerator-go/idgen"
	"strconv"
)

const RowKeyIdGen = "idgen"

func initRowKeyGen(a *apijson.ApiJson) {

	var options = idgen.NewIdGeneratorOptions(1)
	options.WorkerIdBitLength = 6
	options.BaseTime = 1591545600000 // 2020-06-08
	idgen.SetIdGenerator(options)

	a.Config().RowKeyGenFunc(RowKeyIdGen, func(ctx context.Context, req *config.RowKeyGenReq, ret *config.RowKeyGenRet) error {

		if req.AccessName == TableGroups {
			if req.Data["group_id"] != nil {
				return nil
			}
		}

		ret.RowKey(strconv.FormatInt(idgen.NextId(), 10))
		return nil
	})
}
