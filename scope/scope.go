package scope

import (
	"fmt"
	"github.com/2908755265/gorm-auto-query/constant"
	"gorm.io/gorm"
)

func BuildScope(p interface{}) (func(*gorm.DB) *gorm.DB, error) {
	conditions, err := parseObject(p)
	if err != nil {
		return nil, err
	}
	return func(db *gorm.DB) *gorm.DB {
		for _, c := range conditions {
			fmt.Printf("%v\n", *c)
			switch c.op {
			case constant.EqStr:
				db = db.Where(fmt.Sprintf("`%s` %s ?", c.key, constant.Eq), c.val)
			case constant.LtStr:
				db = db.Where(fmt.Sprintf("`%s` %s ?", c.key, constant.Lt), c.val)
			case constant.LteStr:
				db = db.Where(fmt.Sprintf("`%s` %s ?", c.key, constant.Lte), c.val)
			case constant.GtStr:
				db = db.Where(fmt.Sprintf("`%s` %s ?", c.key, constant.Gt), c.val)
			case constant.GteStr:
				db = db.Where(fmt.Sprintf("`%s` %s ?", c.key, constant.Gte), c.val)
			case constant.NeqStr:
				db = db.Where(fmt.Sprintf("`%s` %s ?", c.key, constant.Neq), c.val)
			case constant.LikeStr:
				db = db.Where(fmt.Sprintf("`%s` %s ?", c.key, constant.Like), fmt.Sprintf("%%%s%%", c.val.(string)))
			case constant.LLikeStr:
				db = db.Where(fmt.Sprintf("`%s` %s ?", c.key, constant.Like), fmt.Sprintf("%%%s", c.val.(string)))
			case constant.RLikeStr:
				db = db.Where(fmt.Sprintf("`%s` %s ?", c.key, constant.Like), fmt.Sprintf("%s%%", c.val.(string)))
			case constant.InStr:
				db = db.Where(fmt.Sprintf("`%s` %s ?", c.key, constant.In), c.val)
			case constant.IsNullStr:
				db = db.Where(fmt.Sprintf("`%s` %s", c.key, constant.IsNull))
			case constant.IsNotNullStr:
				db = db.Where(fmt.Sprintf("`%s` %s", c.key, constant.IsNotNull))
			}
		}
		return db
	}, nil
}
