package registry

import (
	"grpc-mafia/registry/db"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// POST /users/:login - create/update user
// GET /users/:login - get user
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
		EndWithError(c, err)
		return
	}

	if err := c.MustBindWith(&info, binding.FormMultipart); err != nil {
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
			EndWithError(c, err)
			return
		}
	}

	new, err := Server.db.UpdateUser(user)
	if err != nil {
		EndWithError(c, err)
		return
	}

	c.JSON(200, new)
}

func getUserHandler(c *gin.Context) {
	info := userInfo{}

	if err := c.BindUri(&info); err != nil {
		EndWithError(c, err)
		return
	}

	user, err := Server.db.GetUserOrCreateDefault(info.Login)
	if err != nil {
		EndWithError(c, err)
		return
	}

	c.JSON(200, user)
}

func deleteUserHandler(c *gin.Context) {
	info := userInfo{}

	if err := c.BindUri(&info); err != nil {
		EndWithError(c, err)
		return
	}

	user, err := Server.db.GetUserOrCreateDefault(info.Login)
	if err != nil {
		EndWithError(c, err)
		return
	}

	Server.db.DeleteUser(info.Login)

	c.JSON(200, user)
}

func registerUsersRoutes(r *gin.Engine) {
	r.POST("/users/:login", updateUserInfoHandler)
	r.GET("/users/:login", getUserHandler)
	r.DELETE("/users/:login", deleteUserHandler)
}
