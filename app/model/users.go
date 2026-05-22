package model

import (
	"ginskeleton/app/global/variable"
	"ginskeleton/app/service/users/token_cache_redis"
	"ginskeleton/app/utils/md5_encrypt"
	"go.uber.org/zap"
	"time"
)

func CreateUserFactory(sqlType string) *UsersModel {
	return &UsersModel{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type UsersModel struct {
	BaseModel
	UserName    string `gorm:"column:user_name" json:"user_name"`
	Pass        string `json:"-"`
	Phone       string `json:"phone"`
	RealName    string `gorm:"column:real_name" json:"real_name"`
	Status      int    `json:"status"`
	Token       string `json:"token"`
	LastLoginIp string `gorm:"column:last_login_ip" json:"last_login_ip"`
}

func (u *UsersModel) TableName() string {
	return "tb_users"
}

func (u *UsersModel) Register(userName, pass, userIp string) bool {
	sql := "INSERT  INTO tb_users(user_name,pass,last_login_ip) SELECT ?,?,? FROM DUAL   WHERE NOT EXISTS (SELECT 1  FROM tb_users WHERE  user_name=?)"
	result := u.Exec(sql, userName, pass, userIp, userName)
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (u *UsersModel) Login(userName string, pass string) *UsersModel {
	sql := "select id, user_name,real_name,pass,phone  from tb_users where  user_name=?  limit 1"
	result := u.Raw(sql, userName).First(u)
	if result.Error == nil {

		if len(u.Pass) > 0 && (u.Pass == md5_encrypt.Base64Md5(pass)) {
			return u
		}
	} else {
		variable.ZapLog.Error("根据账号查询单条记录出错:", zap.Error(result.Error))
	}
	return nil
}

func (u *UsersModel) OauthLoginToken(userId int64, token string, expiresAt int64, clientIp string) bool {
	sql := `
		INSERT   INTO  tb_oauth_access_tokens(fr_user_id,action_name,token,expires_at,client_ip)
		SELECT  ?,'login',? ,?,? FROM DUAL    WHERE   NOT   EXISTS(SELECT  1  FROM  tb_oauth_access_tokens a WHERE  a.fr_user_id=?  AND a.action_name='login' AND a.token=?  )
	`

	if u.Exec(sql, userId, token, time.Unix(expiresAt, 0).Format(variable.DateFormat), clientIp, userId, token).Error == nil {

		if variable.ConfigYml.GetInt("Token.IsCacheToRedis") == 1 {
			go u.ValidTokenCacheToRedis(userId)
		}
		return true
	}
	return false
}

func (u *UsersModel) OauthRefreshConditionCheck(userId int64, oldToken string) bool {

	var oldTokenIsExists int
	sql := "SELECT count(*)  as  counts FROM tb_oauth_access_tokens  WHERE fr_user_id =? and token=? and NOW()<DATE_ADD(expires_at,INTERVAL  ? SECOND)"
	if u.Raw(sql, userId, oldToken, variable.ConfigYml.GetInt64("Token.JwtTokenRefreshAllowSec")).First(&oldTokenIsExists).Error == nil && oldTokenIsExists == 1 {
		return true
	}
	return false
}

func (u *UsersModel) OauthRefreshToken(userId, expiresAt int64, oldToken, newToken, clientIp string) bool {
	sql := "UPDATE   tb_oauth_access_tokens   SET  token=? ,expires_at=?,client_ip=?,updated_at=NOW(),action_name='refresh'  WHERE   fr_user_id=? AND token=?"
	if u.Exec(sql, newToken, time.Unix(expiresAt, 0).Format(variable.DateFormat), clientIp, userId, oldToken).Error == nil {

		if variable.ConfigYml.GetInt("Token.IsCacheToRedis") == 1 {
			go u.ValidTokenCacheToRedis(userId)
		}
		go u.UpdateUserloginInfo(clientIp, userId)
		return true
	}
	return false
}

func (u *UsersModel) UpdateUserloginInfo(last_login_ip string, userId int64) {
	sql := "UPDATE  tb_users   SET  login_times=IFNULL(login_times,0)+1,last_login_ip=?,last_login_time=?  WHERE   id=?  "
	_ = u.Exec(sql, last_login_ip, time.Now().Format(variable.DateFormat), userId)
}

func (u *UsersModel) OauthResetToken(userId int, newPass, clientIp string) bool {

	userItem, err := u.ShowOneItem(userId)
	if userItem != nil && err == nil && userItem.Pass == newPass {
		return true
	} else if userItem != nil {

		if variable.ConfigYml.GetInt("Token.IsCacheToRedis") == 1 {
			go u.DelTokenCacheFromRedis(int64(userId))
		}

		sql := "UPDATE  tb_oauth_access_tokens  SET  revoked=1,updated_at=NOW(),action_name='ResetPass',client_ip=?  WHERE  fr_user_id=?  "
		if u.Exec(sql, clientIp, userId).Error == nil {
			return true
		}
	}
	return false
}

func (u *UsersModel) OauthDestroyToken(userId int) bool {

	sql := "DELETE FROM  tb_oauth_access_tokens WHERE  fr_user_id=?  "

	if u.Exec(sql, userId).Error == nil {
		return true
	}
	return false
}

func (u *UsersModel) OauthCheckTokenIsOk(userId int64, token string) bool {
	sql := "SELECT   token  FROM  `tb_oauth_access_tokens`  WHERE   fr_user_id=?  AND  revoked=0  AND  expires_at>NOW() ORDER  BY  expires_at  DESC , updated_at  DESC  LIMIT ?"
	maxOnlineUsers := variable.ConfigYml.GetInt("Token.JwtTokenOnlineUsers")
	rows, err := u.Raw(sql, userId, maxOnlineUsers).Rows()
	defer func() {

		_ = rows.Close()
	}()
	if err == nil && rows != nil {
		for rows.Next() {
			var tempToken string
			err := rows.Scan(&tempToken)
			if err == nil {
				if tempToken == token {
					return true
				}
			}
		}
	}
	return false
}

func (u *UsersModel) SetTokenInvalid(userId int) bool {
	sql := "delete from  `tb_oauth_access_tokens`  where  `fr_user_id`=?  "
	if u.Exec(sql, userId).Error == nil {
		if u.Exec("update  tb_users  set  status=0 where   id=?", userId).Error == nil {
			return true
		}
	}
	return false
}

func (u *UsersModel) ShowOneItem(userId int) (*UsersModel, error) {
	sql := "SELECT  `id`, `user_name`,`pass`, `real_name`, `phone`, `status` FROM  `tb_users`  WHERE `status`=1 and   id=? LIMIT 1"
	result := u.Raw(sql, userId).First(u)
	if result.Error == nil {
		return u, nil
	} else {
		return nil, result.Error
	}
}

func (u *UsersModel) counts(userName string) (counts int64) {
	sql := "SELECT  count(*) as counts  FROM  tb_users  WHERE status=1 and   user_name like ?"
	if res := u.Raw(sql, "%"+userName+"%").First(&counts); res.Error != nil {
		variable.ZapLog.Error("UsersModel - counts 查询数据条数出错", zap.Error(res.Error))
	}
	return counts
}

func (u *UsersModel) Show(userName string, limitStart, limitItems int) (counts int64, temp []UsersModel) {
	if counts = u.counts(userName); counts > 0 {
		sql := "SELECT  `id`, `user_name`, `real_name`, `phone`,last_login_ip, `status`,created_at,updated_at  FROM  `tb_users`  WHERE `status`=1 and   user_name like ? LIMIT ?,?"
		if res := u.Raw(sql, "%"+userName+"%", limitStart, limitItems).Find(&temp); res.RowsAffected > 0 {
			return counts, temp
		}
	}
	return 0, nil
}

func (u *UsersModel) Store(userName string, pass string, realName string, phone string, remark string) bool {
	sql := "INSERT  INTO tb_users(user_name,pass,real_name,phone,remark) SELECT ?,?,?,?,? FROM DUAL   WHERE NOT EXISTS (SELECT 1  FROM tb_users WHERE  user_name=?)"
	if u.Exec(sql, userName, pass, realName, phone, remark, userName).RowsAffected > 0 {
		return true
	}
	return false
}

func (u *UsersModel) UpdateDataCheckUserNameIsUsed(userId int, userName string) (exists int64) {
	sql := "select count(*) as counts from tb_users where  id!=?  AND user_name=?"
	_ = u.Raw(sql, userId, userName).First(&exists)
	return exists
}

func (u *UsersModel) Update(id int, userName string, pass string, realName string, phone string, remark string, clientIp string) bool {
	sql := "update tb_users set user_name=?,pass=?,real_name=?,phone=?,remark=?  WHERE status=1 AND id=?"
	if u.Exec(sql, userName, pass, realName, phone, remark, id).RowsAffected >= 0 {
		if u.OauthResetToken(id, pass, clientIp) {
			return true
		}
	}
	return false
}

func (u *UsersModel) Destroy(id int) bool {

	if variable.ConfigYml.GetInt("Token.IsCacheToRedis") == 1 {
		go u.DelTokenCacheFromRedis(int64(id))
	}
	if u.Delete(u, id).Error == nil {
		if u.OauthDestroyToken(id) {
			return true
		}
	}
	return false
}

func (u *UsersModel) ValidTokenCacheToRedis(userId int64) {
	tokenCacheRedisFact := token_cache_redis.CreateUsersTokenCacheFactory(userId)
	if tokenCacheRedisFact == nil {
		variable.ZapLog.Error("redis连接失败，请检查配置")
		return
	}
	defer tokenCacheRedisFact.ReleaseRedisConn()

	sql := "SELECT   token,expires_at  FROM  `tb_oauth_access_tokens`  WHERE   fr_user_id=?  AND  revoked=0  AND  expires_at>NOW() ORDER  BY  expires_at  DESC , updated_at  DESC  LIMIT ?"
	maxOnlineUsers := variable.ConfigYml.GetInt("Token.JwtTokenOnlineUsers")
	rows, err := u.Raw(sql, userId, maxOnlineUsers).Rows()
	defer func() {

		_ = rows.Close()
	}()

	var tempToken, expires string
	if err == nil && rows != nil {
		for i := 1; rows.Next(); i++ {
			err = rows.Scan(&tempToken, &expires)
			if err == nil {
				if ts, err := time.ParseInLocation(variable.DateFormat, expires, time.Local); err == nil {
					tokenCacheRedisFact.SetTokenCache(ts.Unix(), tempToken)

					if i == 1 {
						tokenCacheRedisFact.SetUserTokenExpire(ts.Unix())
					}
				} else {
					variable.ZapLog.Error("expires_at 转换位时间戳出错", zap.Error(err))
				}
			}
		}
	}

	tokenCacheRedisFact.DelOverMaxOnlineCache()
}

func (u *UsersModel) DelTokenCacheFromRedis(userId int64) {
	tokenCacheRedisFact := token_cache_redis.CreateUsersTokenCacheFactory(userId)
	if tokenCacheRedisFact == nil {
		variable.ZapLog.Error("redis连接失败，请检查配置")
		return
	}
	tokenCacheRedisFact.ClearUserToken()
	tokenCacheRedisFact.ReleaseRedisConn()
}
