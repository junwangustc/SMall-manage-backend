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

type Commodity struct {
	C_Id            int       `json:"c_id"`
	C_current_price float64   `json:"c_current_price"`
	C_datetime      time.Time `json:"c_datetime"`
	C_detail        string    `json:"c_detail"`
	C_img           string    `json:"c_img"`
	C_name          string    `json:"c_name"`
	C_other         string    `json:"c_other"`
	C_primary_price float64   `json:"c_primary_price"`
	C_sales_num     float64   `json:"c_sales_num"`
	C_status        string    `json:"c_status"`
	C_stock         float64   `json:"c_stock"`
	M_id            int       `json:"m_id"`
	M_name          string    `json:"m_name"`
}

//获取该表所有的纪录的数目
func GetCountCommoditys(c *gin.Context) {
	var count int64
	err := db.QueryRow("select count(*) from commodity").Scan(&count)
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
func GetCommoditysByPage(c *gin.Context) {
	var (
		commodity  Commodity
		commoditys []Commodity
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
	rows, err := db.Query("select commodity.c_id,menu_classify.m_name,commodity.c_current_price,commodity.c_datetime,commodity.c_detail,commodity.c_img,commodity.c_name,commodity.c_other,commodity.c_primary_price,commodity.c_sales_num,commodity.c_status,commodity.c_stock ,commodity.m_id from commodity left join menu_classify on commodity.m_id=menu_classify.m_id;")
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	for rows.Next() {
		count++
		var m_name interface{}
		if count >= (page_start) && count <= (page_end) {
			err = rows.Scan(&commodity.C_Id, &m_name, &commodity.C_current_price, &commodity.C_datetime, &commodity.C_detail, &commodity.C_img, &commodity.C_name, &commodity.C_other, &commodity.C_primary_price, &commodity.C_sales_num, &commodity.C_status, &commodity.C_stock, &commodity.M_id)
			if err != nil {
				log.Println(err.Error())
				rows.Close()
				SetErrorRespones(c, err.Error())
				return
			}
			var M_name []uint8
			var M_name_str string
			var ok bool
			if M_name, ok = m_name.([]uint8); !ok {
				//M_name = "无效分类"
				M_name_str = "无效分类"
			} else {
				M_name_str = string(M_name)
			}
			commodity.M_name = M_name_str
			commoditys = append(commoditys, commodity)
		}
	}
	defer rows.Close()
	c.JSON(http.StatusOK, gin.H{
		"status": httpOK,
		"error":  "",
		"result": commoditys,
		"count":  len(commoditys),
	})

}

//获取某一个id的数据
func GetCommodity(c *gin.Context) {
	var (
		commodity Commodity
	)
	id := c.Param("id")
	row := db.QueryRow("select c_id,c_current_price,c_datetime,c_detail,c_img,c_name,c_other,c_primary_price,c_sales_num,c_status,c_stock,m_id from commodity where id = ?;", id)
	err := row.Scan(&commodity.C_Id, &commodity.C_current_price, &commodity.C_datetime, &commodity.C_detail, &commodity.C_img, &commodity.C_name, &commodity.C_other, &commodity.C_primary_price, &commodity.C_sales_num, &commodity.C_status, &commodity.C_stock, &commodity.M_id)
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
			"result": commodity,
			"count":  1,
		})

	}
}

//获取所有纪录
func GetCommoditys(c *gin.Context) {
	var (
		commodity  Commodity
		commoditys []Commodity
	)
	rows, err := db.Query("select c_id,c_current_price,c_datetime,c_detail,c_img,c_name,c_other,c_primary_price,c_sales_num,c_status,c_stock,m_id from commodity;")
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	for rows.Next() {
		err = rows.Scan(&commodity.C_Id, &commodity.C_current_price, &commodity.C_datetime, &commodity.C_detail, &commodity.C_img, &commodity.C_name, &commodity.C_other, &commodity.C_primary_price, &commodity.C_sales_num, &commodity.C_status, &commodity.C_stock, &commodity.M_id)
		commoditys = append(commoditys, commodity)
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
		"result": commoditys,
		"count":  len(commoditys),
	})

}

//添加一行数据
func PostCommodity(c *gin.Context) {
	var buffer bytes.Buffer
	c_current_price := c.PostForm("c_current_price")
	c_datetime := c.PostForm("c_datetime")
	c_detail := c.PostForm("c_detail")
	c_img := c.PostForm("c_img")
	c_name := c.PostForm("c_name")
	c_other := c.PostForm("c_other")
	c_primary_price := c.PostForm("c_primary_price")
	c_sales_num := c.PostForm("c_sales_num")
	c_status := c.PostForm("c_status")
	c_stock := c.PostForm("c_stock")
	m_id := c.PostForm("m_id")

	stmt, err := db.Prepare("insert into commodity (c_current_price,c_datetime,c_detail,c_img,c_name,c_other,c_primary_price,c_sales_num,c_status,c_stock,m_id) values(?,?,?,?,?,?,?,?,?,?,?);")
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	_, err = stmt.Exec(c_current_price, c_datetime, c_detail, c_img, c_name, c_other, c_primary_price, c_sales_num, c_status, c_stock, m_id)

	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}

	buffer.WriteString(c_current_price)
	buffer.WriteString("  ")
	buffer.WriteString(c_datetime)
	buffer.WriteString("  ")
	buffer.WriteString(c_detail)
	buffer.WriteString("  ")
	buffer.WriteString(c_img)
	buffer.WriteString("  ")
	buffer.WriteString(c_name)
	buffer.WriteString("  ")
	buffer.WriteString(c_other)
	buffer.WriteString("  ")
	buffer.WriteString(c_primary_price)
	buffer.WriteString("  ")
	buffer.WriteString(c_sales_num)
	buffer.WriteString("  ")
	buffer.WriteString(c_status)
	buffer.WriteString("  ")
	buffer.WriteString(c_stock)
	buffer.WriteString("  ")
	buffer.WriteString(m_id)
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
func PutCommodity(c *gin.Context) {
	var buffer bytes.Buffer
	id := c.Param("id")
	c_current_price := c.PostForm("c_current_price")
	c_datetime := c.PostForm("c_datetime")
	c_detail := c.PostForm("c_detail")
	c_img := c.PostForm("c_img")
	c_name := c.PostForm("c_name")
	c_other := c.PostForm("c_other")
	c_primary_price := c.PostForm("c_primary_price")
	c_sales_num := c.PostForm("c_sales_num")
	c_status := c.PostForm("c_status")
	c_stock := c.PostForm("c_stock")
	m_id := c.PostForm("m_id")

	stmt, err := db.Prepare("update commodity set c_current_price=?,c_datetime=?,c_detail=?,c_img=?,c_name=?,c_other=?,c_primary_price=?,c_sales_num=?,c_status=?,c_stock=?,m_id=? where c_id= ?;")
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	_, err = stmt.Exec(c_current_price, c_datetime, c_detail, c_img, c_name, c_other, c_primary_price, c_sales_num, c_status, c_stock, m_id, id)
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}

	// Fastest way to append strings
	buffer.WriteString(c_current_price)
	buffer.WriteString("  ")
	buffer.WriteString(c_datetime)
	buffer.WriteString("  ")
	buffer.WriteString(c_detail)
	buffer.WriteString("  ")
	buffer.WriteString(c_img)
	buffer.WriteString("  ")
	buffer.WriteString(c_name)
	buffer.WriteString("  ")
	buffer.WriteString(c_other)
	buffer.WriteString("  ")
	buffer.WriteString(c_primary_price)
	buffer.WriteString("  ")
	buffer.WriteString(c_sales_num)
	buffer.WriteString("  ")
	buffer.WriteString(c_status)
	buffer.WriteString("  ")
	buffer.WriteString(c_stock)
	buffer.WriteString("  ")
	buffer.WriteString(m_id)
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
func DeleteCommodity(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("delete from commodity where c_id= ?;")
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
