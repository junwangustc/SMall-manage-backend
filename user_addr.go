package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/junwangustc/ustclog"
)

type User_addr struct {
	Ua_Id       int       `json:"ua_id"`
	U_id        int       `json:"u_id"`
	Ua_addr     string    `json:"ua_addr"`
	Ua_datetime time.Time `json:"ua_datetime"`
	Ua_name     string    `json:"ua_name"`
	Ua_other    string    `json:"ua_other"`
	Ua_tel      string    `json:"ua_tel"`
}

//获取该表所有的纪录的数目
func GetCountUser_addrs(c *gin.Context) {
	var count int64
	err := db.QueryRow("select count(*) from user_addr").Scan(&count)
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
func GetUser_addrsByPage(c *gin.Context) {
	var (
		user_addr  User_addr
		user_addrs []User_addr
	)

	id := c.Param("pageid")
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
	rows, err := db.Query("select ua_id,u_id,ua_addr,ua_datetime,ua_name,ua_other,ua_tel from user_addr;")
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	for rows.Next() {
		count++
		if count >= (page_start) && count <= (page_end) {
			err = rows.Scan(&user_addr.Ua_Id, &user_addr.U_id, &user_addr.Ua_addr, &user_addr.Ua_datetime, &user_addr.Ua_name, &user_addr.Ua_other, &user_addr.Ua_tel)
			if err != nil {
				log.Println(err.Error())
				rows.Close()
				SetErrorRespones(c, err.Error())
				return
			}
			user_addrs = append(user_addrs, user_addr)
		}
	}
	defer rows.Close()
	c.JSON(http.StatusOK, gin.H{
		"status": httpOK,
		"error":  "",
		"result": user_addrs,
		"count":  len(user_addrs),
	})

}

//获取某一个id的数据
func GetUser_addr(c *gin.Context) {
	var (
		user_addr User_addr
	)
	id := c.Param("id")
	row := db.QueryRow("select ua_id,u_id,ua_addr,ua_datetime,ua_name,ua_other,ua_tel from user_addr where id = ?;", id)
	err := row.Scan(&user_addr.Ua_Id, &user_addr.U_id, &user_addr.Ua_addr, &user_addr.Ua_datetime, &user_addr.Ua_name, &user_addr.Ua_other, &user_addr.Ua_tel)
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
			"result": user_addr,
			"count":  1,
		})

	}
}

//获取所有纪录
func GetUser_addrs(c *gin.Context) {
	var (
		user_addr  User_addr
		user_addrs []User_addr
	)
	rows, err := db.Query("select ua_id,u_id,ua_addr,ua_datetime,ua_name,ua_other,ua_tel from user_addr;")
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	for rows.Next() {
		err = rows.Scan(&user_addr.Ua_Id, &user_addr.U_id, &user_addr.Ua_addr, &user_addr.Ua_datetime, &user_addr.Ua_name, &user_addr.Ua_other, &user_addr.Ua_tel)
		user_addrs = append(user_addrs, user_addr)
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
		"result": user_addrs,
		"count":  len(user_addrs),
	})

}

//添加一行数据
func PostUser_addr(c *gin.Context) {
	var buffer bytes.Buffer
	u_id := c.PostForm("u_id")
	ua_addr := c.PostForm("ua_addr")
	ua_datetime := c.PostForm("ua_datetime")
	ua_name := c.PostForm("ua_name")
	ua_other := c.PostForm("ua_other")
	ua_tel := c.PostForm("ua_tel")

	stmt, err := db.Prepare("insert into user_addr (u_id,ua_addr,ua_datetime,ua_name,ua_other,ua_tel) values(?,?,?,?,?,?);")
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	_, err = stmt.Exec(u_id, ua_addr, ua_datetime, ua_name, ua_other, ua_tel)

	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}

	buffer.WriteString(u_id)
	buffer.WriteString("  ")
	buffer.WriteString(ua_addr)
	buffer.WriteString("  ")
	buffer.WriteString(ua_datetime)
	buffer.WriteString("  ")
	buffer.WriteString(ua_name)
	buffer.WriteString("  ")
	buffer.WriteString(ua_other)
	buffer.WriteString("  ")
	buffer.WriteString(ua_tel)
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
func PutUser_addr(c *gin.Context) {
	var buffer bytes.Buffer
	id := c.Param("id")
	u_id := c.PostForm("u_id")
	ua_addr := c.PostForm("ua_addr")
	ua_datetime := c.PostForm("ua_datetime")
	ua_name := c.PostForm("ua_name")
	ua_other := c.PostForm("ua_other")
	ua_tel := c.PostForm("ua_tel")

	stmt, err := db.Prepare("update user_addr set u_id=?,ua_addr=?,ua_datetime=?,ua_name=?,ua_other=?,ua_tel=? where ua_id= ?;")
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	_, err = stmt.Exec(u_id, ua_addr, ua_datetime, ua_name, ua_other, ua_tel, id)
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}

	// Fastest way to append strings
	buffer.WriteString(u_id)
	buffer.WriteString("  ")
	buffer.WriteString(ua_addr)
	buffer.WriteString("  ")
	buffer.WriteString(ua_datetime)
	buffer.WriteString("  ")
	buffer.WriteString(ua_name)
	buffer.WriteString("  ")
	buffer.WriteString(ua_other)
	buffer.WriteString("  ")
	buffer.WriteString(ua_tel)
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
func DeleteUser_addr(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("delete from user_addr where ua_id= ?;")
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
