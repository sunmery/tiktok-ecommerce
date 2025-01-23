package dal

import (
	"github.com/sunmery/tiktok-e-commence/app/cart/biz/dal/mysql"
	"github.com/sunmery/tiktok-e-commence/app/cart/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
