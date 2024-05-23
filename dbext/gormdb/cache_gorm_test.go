package gormdb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
	"time"
)

var (
	cacheUsersAppIdThirdUserIdPrefix = "cache:users:appId:thirdUserId:"
	cacheUsersPKPrefix               = "cache:users:id:"
)

type Users struct {
	Id                  int64          `db:"id" json:"id,omitempty"` // id
	PersonMsgCurrentSeq int64          `db:"person_msg_current_seq" json:"personMsgCurrentSeq,omitempty"`
	PersonMsgReadSeq    int64          `db:"person_msg_read_seq" json:"personMsgReadSeq,omitempty"`
	PublicCurrentSeq    int64          `db:"public_current_seq" json:",omitempty"`
	PublicReadSeq       int64          `db:"public_read_seq" json:"publicReadSeq,omitempty"`
	ThirdUserId         int64          `db:"third_user_id" json:"thirdUserId,omitempty"`             // 第三方用户id
	Nick                string         `db:"nick" json:"nick,omitempty"`                             // 昵称，最多50个字符
	FaceUrl             string         `db:"face_url" json:"faceUrl,omitempty"`                      // 头像
	Pwd                 string         `db:"pwd" json:"pwd,omitempty"`                               // 密码
	Phone               sql.NullString `db:"phone" json:"phone,omitempty"`                           // 手机号
	AppId               int64          `db:"app_id" json:"appId,omitempty"`                          // 应用id
	Account             sql.NullString `db:"account" json:"account,omitempty"`                       // 账户名称
	OfflineTime         sql.NullTime   `db:"offline_time" json:"offlineTime,omitempty"`              // 离线时间
	OfflinePlatformId   int64          `db:"offline_platform_id" json:"offlinePlatformId,omitempty"` // 离线平台id，0:无数据，1:iOS，2:Android，3:win，4:Linux，5:pcWeb，6:h5
	Status              int64          `db:"status" json:"status,omitempty"`                         // 当前状态,0:正常，1:禁用
	OnlineStatus        int64          `db:"OnlineStatus" json:"onlineStatus,omitempty"`             // 在线状态,0:离线，1:在线 2:隐身
	CreatedAt           time.Time      `db:"created_at" json:"createdAt,omitempty"`                  // 创建时间
	UpdatedAt           time.Time      `db:"updated_at" json:"updatedAt,omitempty"`                  // 更新时间
	OnlinePlatforms     string         `db:"online_platforms" json:"onlinePlatforms,omitempty"`      // 在线平台
}

func TestCacheGormDB_QueryCtx(t *testing.T) {
	appId := 1
	thirdUserID := 639412975567962112
	//	测试通过唯一键查询
	usersAppIdAccountKey := fmt.Sprintf("%s%v:%v", cacheUsersAppIdThirdUserIdPrefix, appId, thirdUserID)
	rdb := redis.NewClient(&redis.Options{
		Addr:     ":6379",
		Password: "root",
	})
	cacheGormDB := MustNewCacheGormDB[Users, int64](Config{
		DSN:    "root:root@tcp(127.0.0.1:3306)/im_server?charset=utf8mb4&parseTime=True&loc=Local",
		DBType: DBTYPE_MySQL,
		GormConfig: gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
		Rdb:               rdb,
		NotFoundExpireSec: 60 * 60,
		CacheExpireSec:    2000,
		RandSec:           200,
	})

	var user Users
	err := cacheGormDB.QueryOneCtx(context.Background(), &user, usersAppIdAccountKey, func(ctx context.Context, p *int64, db *gorm.DB) error {
		return db.Model(&Users{}).Select("id").Where("app_id = ? and third_user_id = ?", appId, thirdUserID).Take(p).Error
	}, cacheUsersPKPrefix, func(ctx context.Context, r *Users, db *gorm.DB) error {
		return db.Model(&Users{}).Where("app_id = ? and third_user_id = ?", appId, thirdUserID).Take(r).Error

	})
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("aa:%v", user)

	//wg := sync.WaitGroup{}
	//num := int64(0)
	//for i := 0; i < 1; i++ {
	//	wg.Add(1)
	//	go func() {
	//		defer func() {
	//			atomic.AddInt64(&num, 1)
	//			wg.Done()
	//		}()
	//		var user Users
	//		err := cacheGormDB.QueryOneCtx(context.Background(), &user, usersAppIdAccountKey, func(ctx context.Context, p int64, db *gorm.DB) error {
	//			return db.Model(&Users{}).Select("id").Where("app_id = ? and third_user_id = ?", appId, thirdUserID).Take(&p).Error
	//
	//		}, cacheUsersPKPrefix, func(ctx context.Context, db *gorm.DB) (Users, error) {
	//			err := db.Model(&Users{}).Where("app_id = ? and third_user_id = ?", appId, thirdUserID).Take(&user).Error
	//
	//			return user, err
	//		})
	//		if err != nil {
	//			t.Error(err)
	//			return
	//		}
	//	}()
	//}
	//fmt.Println("---")
	//wg.Wait()
	fmt.Println("所有查询结束")
	//fmt.Println("结束:", num)
	//time.Sleep(time.Hour)

}
