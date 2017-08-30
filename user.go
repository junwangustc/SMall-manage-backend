package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/junwangustc/ustclog"
)

type User struct {
	U_Id       int    `json:"u_id"`
	U_account  string `json:"u_account"`
	U_datetime string `json:"u_datetime"`
	U_level    string `json:"u_level"`
	U_name     string `json:"u_name"`
	U_other    string `json:"u_other"`
	U_psd      string `json:"u_psd"`
	U_score    string `json:"u_score"`
	U_status   string `json:"u_status"`
	U_tel      string `json:"u_tel"`
}

//获取该表所有的纪录的数目
func GetCountUsers(c *gin.Context) {
	var count int64
	err := db.QueryRow("select count(*) from user").Scan(&count)
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": httpOK,
		"error":  "",
		"count":  count,
	})

}

//获取指定页的数据，传入pageid 返回数据
func GetUsersByPage(c *gin.Context) {
	var (
		user  User
		users []User
	)

	id := c.Param("pageid")
	//	page_size := c.Query("page_size")
	page_size := "10"
	pageidInt, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, "请输入一个合法的页面id")
		return
	}
	page_sizeInt, err := strconv.Atoi(page_size)
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, "请输入一个合法的页面大小")
		return
	}

	page_start := (pageidInt-1)*page_sizeInt + 1
	page_end := pageidInt * page_sizeInt
	var count = 0
	rows, err := db.Query("select u_id,u_account,u_datetime,u_level,u_name,u_other,u_psd,u_score,u_status,u_tel from user;")
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	for rows.Next() {
		count++
		if count >= (page_start) && count <= (page_end) {
			err = rows.Scan(&user.U_Id, &user.U_account, &user.U_datetime, &user.U_level, &user.U_name, &user.U_other, &user.U_psd, &user.U_score, &user.U_status, &user.U_tel)
			if err != nil {
				log.Println(err.Error())
				rows.Close()
				SetErrorRespones(c, err.Error())
				return
			}
			users = append(users, user)
		}
	}
	defer rows.Close()
	c.JSON(http.StatusOK, gin.H{
		"status": httpOK,
		"error":  "",
		"result": users,
		"count":  len(users),
	})

}

//获取某一个id的数据
func GetUser(c *gin.Context) {
	var (
		user User
	)
	id := c.Param("id")
	row := db.QueryRow("select u_id,u_account,u_datetime,u_level,u_name,u_other,u_psd,u_score,u_status,u_tel from user where u_id = ?;", id)
	err := row.Scan(&user.U_Id, &user.U_account, &user.U_datetime, &user.U_level, &user.U_name, &user.U_other, &user.U_psd, &user.U_score, &user.U_status, &user.U_tel)
	if err != nil {
		// If no results send null
		c.JSON(http.StatusOK, gin.H{
			"status": httpOK,
			"error":  "",
			"result": nil,
			"count":  0,
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": httpOK,
			"error":  "",
			"result": user,
			"count":  1,
		})

	}
}

//获取所有纪录
func GetUsers(c *gin.Context) {
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
			SetErrorRespones(c, err.Error())
			return
		}
	}
	defer rows.Close()
	c.JSON(http.StatusOK, gin.H{
		"status": httpOK,
		"error":  "",
		"result": users,
		"count":  len(users),
	})

}

//添加一行数据
func PostUser(c *gin.Context) {
	var buffer bytes.Buffer
	u_account := c.PostForm("u_account")
	u_datetime := c.PostForm("u_datetime")
	u_level := c.PostForm("u_level")
	u_name := c.PostForm("u_name")
	u_other := c.PostForm("u_other")
	u_psd := c.PostForm("u_psd")
	u_score := c.PostForm("u_score")
	u_status := c.PostForm("u_status")
	u_tel := c.PostForm("u_tel")

	stmt, err := db.Prepare("insert into user (u_account,u_datetime,u_level,u_name,u_other,u_psd,u_score,u_status,u_tel) values(?,?,?,?,?,?,?,?,?);")
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	_, err = stmt.Exec(u_account, u_datetime, u_level, u_name, u_other, u_psd, u_score, u_status, u_tel)

	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}

	buffer.WriteString(u_account)
	buffer.WriteString("  ")
	buffer.WriteString(u_datetime)
	buffer.WriteString("  ")
	buffer.WriteString(u_level)
	buffer.WriteString("  ")
	buffer.WriteString(u_name)
	buffer.WriteString("  ")
	buffer.WriteString(u_other)
	buffer.WriteString("  ")
	buffer.WriteString(u_psd)
	buffer.WriteString("  ")
	buffer.WriteString(u_score)
	buffer.WriteString("  ")
	buffer.WriteString(u_status)
	buffer.WriteString("  ")
	buffer.WriteString(u_tel)
	buffer.WriteString("  ")

	defer stmt.Close()
	_name := buffer.String()
	c.JSON(http.StatusOK, gin.H{
		"status":  httpOK,
		"error":   "",
		"message": fmt.Sprintf(" %s successfully created", _name),
	})
}

//修改某一行数据
func PutUser(c *gin.Context) {
	var buffer bytes.Buffer
	id := c.Param("id")
	u_account := c.PostForm("u_account")
	u_datetime := c.PostForm("u_datetime")
	u_level := c.PostForm("u_level")
	u_name := c.PostForm("u_name")
	u_other := c.PostForm("u_other")
	u_psd := c.PostForm("u_psd")
	u_score := c.PostForm("u_score")
	u_status := c.PostForm("u_status")
	u_tel := c.PostForm("u_tel")

	stmt, err := db.Prepare("update user set u_account=?,u_datetime=?,u_level=?,u_name=?,u_other=?,u_psd=?,u_score=?,u_status=?,u_tel=? where u_id= ?;")
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	_, err = stmt.Exec(u_account, u_datetime, u_level, u_name, u_other, u_psd, u_score, u_status, u_tel, id)
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}

	// Fastest way to append strings
	buffer.WriteString(u_account)
	buffer.WriteString("  ")
	buffer.WriteString(u_datetime)
	buffer.WriteString("  ")
	buffer.WriteString(u_level)
	buffer.WriteString("  ")
	buffer.WriteString(u_name)
	buffer.WriteString("  ")
	buffer.WriteString(u_other)
	buffer.WriteString("  ")
	buffer.WriteString(u_psd)
	buffer.WriteString("  ")
	buffer.WriteString(u_score)
	buffer.WriteString("  ")
	buffer.WriteString(u_status)
	buffer.WriteString("  ")
	buffer.WriteString(u_tel)
	buffer.WriteString("  ")

	defer stmt.Close()
	_name := buffer.String()
	c.JSON(http.StatusOK, gin.H{
		"status":  httpOK,
		"error":   "",
		"message": fmt.Sprintf("Successfully updated to %s", _name),
	})

}

//删除某一行纪录
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("delete from user where u_id= ?;")
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	_, err = stmt.Exec(id)
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  httpOK,
		"error":   "",
		"message": fmt.Sprintf("Successfully deleted user: %s", id),
	})

}
