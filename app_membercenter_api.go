package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetErrorRespones1(c *gin.Context, errorstring string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status": httpError,
		"error":  errorstring,
		"result": nil,
		"count":  0,
	})
}

var SessionMap = make(map[string]int) //保管session的地方
func MakeSession(userName, psd string) string {
	data := []byte(userName + "####" + psd + time.Now().String())
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func GetUidBySession(session string) int {
	if uid, ok := SessionMap[session]; ok {
		return uid
	}

	return -1
}
func DeleteSession(session string) {
	delete(SessionMap, session)
}
func AddSession(session string, uid int) {
	SessionMap[session] = uid
}

//登录api
func API_PostUserLogin(c *gin.Context) {
	//获取参数
	userName := c.PostForm("emaNresu")
	psd := c.PostForm("dsp")
	code := c.PostForm("code")
	var (
		user  User
		users []User
	)
	rows, err := db.Query("select u_id,u_account,u_datetime,u_level,u_name,u_other,u_psd,u_score,u_status,u_tel from user;")
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	for rows.Next() {
		err = rows.Scan(&user.U_Id, &user.U_account, &user.U_datetime, &user.U_level, &user.U_name, &user.U_other, &user.U_psd, &user.U_score, &user.U_status, &user.U_tel)
		users = append(users, user)
		if err != nil {
			log.Println(err.Error())
			rows.Close()
			SetErrorRespones1(c, err.Error())
			return
		}
	}
	defer rows.Close()
	var logined = false
	var session_key = ""
	_ = code
	for _, tmp_user := range users {
		if tmp_user.U_name == userName && tmp_user.U_psd == psd {
			logined = true
			session_key = MakeSession(userName, psd)
			AddSession(session_key, tmp_user.U_Id)
			break
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status":        httpOK,
		"error":         "",
		"result":        logined,
		"local_session": session_key,
	})

}
func API_PostUserLogout(c *gin.Context) {
	local_session := c.PostForm("local_session")
	//把local_session 在redis中注销掉
	DeleteSession(local_session)
	c.JSON(http.StatusOK, gin.H{
		"status": httpOK,
		"error":  "",
	})

}

//注册api
func API_PostRegister(c *gin.Context) {

}

//获取个人信息API
func API_GetUserInfo(c *gin.Context) {
	//获取个人信息
	local_session := c.PostForm("local_session") //session 和用户关联起来
	uid := GetUidBySession(local_session)
	if uid > 0 {
		var (
			user User
		)
		row := db.QueryRow("select u_id,u_account,u_datetime,u_level,u_name,u_other,u_psd,u_score,u_status,u_tel from user where u_id = ?;", uid)
		err := row.Scan(&user.U_Id, &user.U_account, &user.U_datetime, &user.U_level, &user.U_name, &user.U_other, &user.U_psd, &user.U_score, &user.U_status, &user.U_tel)
		if err != nil {
			// If no results send null
			SetErrorRespones1(c, err.Error())
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": httpOK,
				"error":  "",
				"result": user,
			})

		}

	} else {
		SetErrorRespones1(c, "服务器端找不到该session")
	}

}

//修改个人信息信息
func API_PutUserInfo(c *gin.Context) {

}

//获取收获信息
func API_GetUserAddr(c *gin.Context) {

}

//修改收获信息
func API_PutUserAddr(c *gin.Context) {

}

//删除收获信息
func API_PostUserAddr(c *gin.Context) {

}

func API_DeleteUserAddr(C *gin.Context) {

}
