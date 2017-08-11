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

type Flag_classify struct {
	F_Id       int       `json:"f_id"`
	F_datetime time.Time `json:"f_datetime"`
	F_name     string    `json:"f_name"`
	F_other    string    `json:"f_other"`
	F_status   string    `json:"f_status"`
}

//获取该表所有的纪录的数目
func GetCountFlag_classifys(c *gin.Context) {
	var count int64
	err := db.QueryRow("select count(*) from flag_classify").Scan(&count)
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
func GetFlag_classifysByPage(c *gin.Context) {
	var (
		flag_classify  Flag_classify
		flag_classifys []Flag_classify
	)

	id := c.Param("pageid")
	page_size := c.Query("page_size")
	page_size = "10"
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
	rows, err := db.Query("select f_id,f_datetime,f_name,f_other,f_status from flag_classify;")
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	for rows.Next() {
		count++
		if count >= (page_start) && count <= (page_end) {
			err = rows.Scan(&flag_classify.F_Id, &flag_classify.F_datetime, &flag_classify.F_name, &flag_classify.F_other, &flag_classify.F_status)
			if err != nil {
				log.Println(err.Error())
				rows.Close()
				SetErrorRespones(c, err.Error())
				return
			}
			flag_classifys = append(flag_classifys, flag_classify)
		}
	}
	defer rows.Close()
	c.JSON(http.StatusOK, gin.H{
		"status": httpOK,
		"error":  "",
		"result": flag_classifys,
		"count":  len(flag_classifys),
	})

}

//获取某一个id的数据
func GetFlag_classify(c *gin.Context) {
	var (
		flag_classify Flag_classify
	)
	id := c.Param("id")
	row := db.QueryRow("select f_id,f_datetime,f_name,f_other,f_status from flag_classify where id = ?;", id)
	err := row.Scan(&flag_classify.F_Id, &flag_classify.F_datetime, &flag_classify.F_name, &flag_classify.F_other, &flag_classify.F_status)
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
			"result": flag_classify,
			"count":  1,
		})

	}
}

//获取所有纪录
func GetFlag_classifys(c *gin.Context) {
	var (
		flag_classify  Flag_classify
		flag_classifys []Flag_classify
	)
	rows, err := db.Query("select f_id,f_datetime,f_name,f_other,f_status from flag_classify;")
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	for rows.Next() {
		err = rows.Scan(&flag_classify.F_Id, &flag_classify.F_datetime, &flag_classify.F_name, &flag_classify.F_other, &flag_classify.F_status)
		flag_classifys = append(flag_classifys, flag_classify)
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
		"result": flag_classifys,
		"count":  len(flag_classifys),
	})

}

//添加一行数据
func PostFlag_classify(c *gin.Context) {
	var buffer bytes.Buffer
	f_datetime := c.PostForm("f_datetime")
	f_name := c.PostForm("f_name")
	f_other := c.PostForm("f_other")
	f_status := c.PostForm("f_status")

	stmt, err := db.Prepare("insert into flag_classify (f_datetime,f_name,f_other,f_status) values(?,?,?,?);")
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	_, err = stmt.Exec(f_datetime, f_name, f_other, f_status)

	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}

	buffer.WriteString(f_datetime)
	buffer.WriteString("  ")
	buffer.WriteString(f_name)
	buffer.WriteString("  ")
	buffer.WriteString(f_other)
	buffer.WriteString("  ")
	buffer.WriteString(f_status)
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
func PutFlag_classify(c *gin.Context) {
	var buffer bytes.Buffer
	id := c.Param("id")
	f_datetime := c.PostForm("f_datetime")
	f_name := c.PostForm("f_name")
	f_other := c.PostForm("f_other")
	f_status := c.PostForm("f_status")

	stmt, err := db.Prepare("update flag_classify set f_datetime=?,f_name=?,f_other=?,f_status=? where f_id= ?;")
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	_, err = stmt.Exec(f_datetime, f_name, f_other, f_status, id)
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}

	// Fastest way to append strings
	buffer.WriteString(f_datetime)
	buffer.WriteString("  ")
	buffer.WriteString(f_name)
	buffer.WriteString("  ")
	buffer.WriteString(f_other)
	buffer.WriteString("  ")
	buffer.WriteString(f_status)
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
func DeleteFlag_classify(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("delete from flag_classify where f_id= ?;")
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
