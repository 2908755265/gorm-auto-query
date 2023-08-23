package scope

import (
	"git.rtzhtech.cn/rtcloud/common/util"
	"github.com/2908755265/gorm-auto-query/constant"
	"reflect"
	"strings"
)

func parseObject(obj any) ([]*cond, error) {
	result := make([]*cond, 0)

	rt := reflect.TypeOf(obj)
	rv := reflect.ValueOf(obj)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
		rt = rt.Elem()
	}

	for i := 0; i < rt.NumField(); i++ {
		ft := rt.Field(i)
		k, op, ok := getKeyAndOp(ft)
		if !ok {
			continue
		}

		c := &cond{
			key: k,
			op:  op,
		}

		if util.InArray(op, []string{constant.IsNullStr, constant.IsNotNullStr}) {
			result = append(result, c)
			continue
		}

		fv := rv.Field(i)
		if fv.Kind() == reflect.Ptr {
			if fv.IsNil() {
				continue
			}
			fv = fv.Elem()
		}

		if fv.IsZero() {
			continue
		}

		if fv.Kind() >= reflect.Int && fv.Kind() <= reflect.Int64 {
			c.val = fv.Int()
		} else if fv.Kind() >= reflect.Uint && fv.Kind() <= reflect.Uint64 {
			c.val = fv.Uint()
		} else if fv.Kind() >= reflect.Float32 && fv.Kind() <= reflect.Float64 {
			c.val = fv.Float()
		} else if fv.Kind() == reflect.String {
			c.val = fv.String()
		} else if fv.Kind() == reflect.Slice {
			if fv.Len() == 0 {
				continue
			}
			kd := fv.Index(0).Kind()
			if kd >= reflect.Int && kd <= reflect.Int64 {
				s := make([]int64, 0, fv.Len())
				for i := 0; i < fv.Len(); i++ {
					s = append(s, fv.Index(i).Int())
				}
				c.val = s
			} else if kd >= reflect.Uint && kd <= reflect.Uint64 {
				s := make([]uint64, 0, fv.Len())
				for i := 0; i < fv.Len(); i++ {
					s = append(s, fv.Index(i).Uint())
				}
				c.val = s
			} else if kd >= reflect.Float32 && kd <= reflect.Float64 {
				s := make([]float64, 0, fv.Len())
				for i := 0; i < fv.Len(); i++ {
					s = append(s, fv.Index(i).Float())
				}
				c.val = s
			} else if kd == reflect.String {
				s := make([]string, 0, fv.Len())
				for i := 0; i < fv.Len(); i++ {
					s = append(s, fv.Index(i).String())
				}
				c.val = s
			} else {
				continue
			}
		} else {
			continue
		}
		result = append(result, c)
	}

	return result, nil
}

func getKeyAndOp(ft reflect.StructField) (string, string, bool) {
	var key, op string

	tag := ft.Tag
	sp := tag.Get(constant.ScopeKey)
	// 没有 scope tag
	if sp == "" {
		return key, op, false
	}

	temp := make(map[string]string)
	params := strings.Split(sp, constant.Delimiter)
	for _, param := range params {
		kv := strings.Split(param, constant.ParamDlt)
		if len(kv) != 2 {
			return key, op, false
		}
		temp[kv[0]] = kv[1]
	}

	opStr, ok := temp[constant.KeyTagOp]
	if !ok {
		opStr = constant.EqStr
	}
	if !opStrCheck(opStr) {
		return key, op, false
	}
	op = opStr

	key, ok = temp[constant.KegTagName]
	if !ok {
		// 将字段转为蛇形
		key = parseToSnake([]byte(ft.Name))
	}

	return key, op, true
}

func parseToSnake(name []byte) string {
	result := make([]byte, 0)
	for i, b := range name {
		if b >= 'A' && b <= 'Z' {
			if i != 0 {
				result = append(result, '_')
			}
			result = append(result, b+uint8(32))
		} else {
			result = append(result, b)
		}
	}
	return string(result)
}

func opStrCheck(opStr string) bool {
	return util.InArray(opStr, []string{constant.EqStr, constant.NeqStr, constant.GtStr, constant.GteStr, constant.LtStr, constant.LteStr, constant.LLikeStr, constant.RLikeStr, constant.LikeStr, constant.InStr, constant.IsNullStr, constant.IsNotNullStr})
}
