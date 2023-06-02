package registry

import (
	"grpc-mafia/registry/db"
	"mime/multipart"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	zlog "github.com/rs/zerolog/log"
)

// POST /users/:login - create/update user
// GET /users/:login - get user
// GET /users/?logins=1,2,3
// DELETE /users/:login - delete user

type userInfo struct {
	Login  string                `uri:"login" json:"login"`
	Avatar *multipart.FileHeader `form:"avatar" json:"avatar"`
	Gender string                `form:"gender" json:"gender"`
	Mail   string                `form:"mail" json:"mail"`
}

func updateUserInfoHandler(c *gin.Context) {
	info := userInfo{}

	if err := c.BindUri(&info); err != nil {
		zlog.Error().Err(err).Msg("bind")
		EndWithError(c, err)
		return
	}

	if err := c.MustBindWith(&info, binding.FormMultipart); err != nil {
		zlog.Error().Err(err).Msg("bind")
		EndWithError(c, err)
		return
	}

	user := db.User{
		Login:  info.Login,
		Gender: info.Gender,
		Mail:   info.Mail,
	}

	if info.Avatar != nil {
		user.AvatarFilename = GenRandomName()

		if err := SaveAvatar(info.Avatar, user.AvatarFilename); err != nil {
			zlog.Error().Err(err).Msg("save avatar")
			EndWithError(c, err)
			return
		}
	}

	new, err := Server.db.UpdateUser(user)
	if err != nil {
		zlog.Error().Err(err).Msg("update user")
		EndWithError(c, err)
		return
	}

	c.JSON(200, new)
}

func getUserHandler(c *gin.Context) {
	info := userInfo{}

	if err := c.BindUri(&info); err != nil {
		zlog.Error().Err(err).Msg("bind")
		EndWithError(c, err)
		return
	}

	user, err := Server.db.GetUser(info.Login)
	if err != nil {
		zlog.Error().Err(err).Msg("get user")
		EndWithError(c, err)
		return
	}

	c.JSON(200, user)
}

func getUsersListHandler(c *gin.Context) {
	cgi_logins := c.Query("logins")

	if len(cgi_logins) == 0 {
		users, err := Server.db.GetAllUsers()
		if err != nil {
			zlog.Error().Err(err).Msg("get all users")
			EndWithError(c, err)
		} else {
			c.JSON(200, users)
		}
		return
	}

	logins := strings.Split(cgi_logins, ",")
	users := make([]*db.User, 0)

	for _, login := range logins {
		user, err := Server.db.GetUser(login)
		if err != nil {
			zlog.Error().Err(err).Msg("get user")
			EndWithError(c, err)
			return
		}

		users = append(users, user)
	}

	c.JSON(200, users)
}

func deleteUserHandler(c *gin.Context) {
	info := userInfo{}

	if err := c.BindUri(&info); err != nil {
		zlog.Error().Err(err).Msg("bind")
		EndWithError(c, err)
		return
	}

	user, err := Server.db.GetUser(info.Login)
	if err != nil {
		zlog.Error().Err(err).Msg("get user")
		EndWithError(c, err)
		return
	}

	Server.db.DeleteUser(info.Login)

	c.JSON(200, user)
}

func registerUsersRoutes(r *gin.Engine) {
	r.POST("/users/:login", updateUserInfoHandler)
	r.GET("/users/:login", getUserHandler)
	r.GET("/users", getUsersListHandler)
	r.DELETE("/users/:login", deleteUserHandler)
}