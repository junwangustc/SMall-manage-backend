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

type Order struct {
	O_Id          int
	O_code        string
	O_addr        string
	O_datetime    time.Time
	O_detail      string
	O_other       string
	O_pay_status  string
	O_pay_type    string
	O_receiver    string
	O_status      string
	O_tel         string
	O_total_money float64
	U_id          int
	U_name        string
}

//获取该表所有的纪录的数目
func GetCountOrders(c *gin.Context) {
	var count int64
	err := db.QueryRow("select count(*) from order_new").Scan(&count)
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
func GetOrdersByPage(c *gin.Context) {
	var (
		order  Order
		orders []Order
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
	rows, err := db.Query("select o_id,o_code,o_addr,o_datetime,o_detail,o_other,o_pay_status,o_pay_type,o_receiver,o_status,o_tel,o_total_money,u_id,u_name from order_new;")
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	for rows.Next() {
		count++
		if count >= (page_start) && count <= (page_end) {
			err = rows.Scan(&order.O_Id, &order.O_code, &order.O_addr, &order.O_datetime, &order.O_detail, &order.O_other, &order.O_pay_status, &order.O_pay_type, &order.O_receiver, &order.O_status, &order.O_tel, &order.O_total_money, &order.U_id, &order.U_name)
			if err != nil {
				log.Println(err.Error())
				rows.Close()
				SetErrorRespones(c, err.Error())
				return
			}
			orders = append(orders, order)
		}
	}
	defer rows.Close()
	c.JSON(http.StatusOK, gin.H{
		"status": httpOK,
		"error":  "",
		"result": orders,
		"count":  len(orders),
	})

}

//获取某一个id的数据
func GetOrder(c *gin.Context) {
	var (
		order Order
	)
	id := c.Param("id")
	row := db.QueryRow("select o_id,o_code,o_addr,o_datetime,o_detail,o_other,o_pay_status,o_pay_type,o_receiver,o_status,o_tel,o_total_money,u_id,u_name from order_new where id = ?;", id)
	err := row.Scan(&order.O_Id, &order.O_code, &order.O_addr, &order.O_datetime, &order.O_detail, &order.O_other, &order.O_pay_status, &order.O_pay_type, &order.O_receiver, &order.O_status, &order.O_tel, &order.O_total_money, &order.U_id, &order.U_name)
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
			"result": order,
			"count":  1,
		})

	}
}

//获取所有纪录
func GetOrders(c *gin.Context) {
	var (
		order  Order
		orders []Order
	)
	rows, err := db.Query("select o_id,o_code,o_addr,o_datetime,o_detail,o_other,o_pay_status,o_pay_type,o_receiver,o_status,o_tel,o_total_money,u_id,u_name from order_new;")
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	for rows.Next() {
		err = rows.Scan(&order.O_Id, &order.O_code, &order.O_addr, &order.O_datetime, &order.O_detail, &order.O_other, &order.O_pay_status, &order.O_pay_type, &order.O_receiver, &order.O_status, &order.O_tel, &order.O_total_money, &order.U_id, &order.U_name)
		orders = append(orders, order)
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
		"result": orders,
		"count":  len(orders),
	})

}

//添加一行数据
func PostOrder(c *gin.Context) {
	log.Println("Post order coming------")
	var buffer bytes.Buffer
	o_addr := c.PostForm("o_addr")
	o_code := c.PostForm("o_code")
	o_datetime := c.PostForm("o_datetime")
	o_detail := c.PostForm("o_detail")
	o_other := c.PostForm("o_other")
	o_pay_status := c.PostForm("o_pay_status")
	o_pay_type := c.PostForm("o_pay_type")
	o_receiver := c.PostForm("o_receiver")
	o_status := c.PostForm("o_status")
	o_tel := c.PostForm("o_tel")
	o_total_money := c.PostForm("o_total_money")
	u_id := c.PostForm("u_id")
	u_name := c.PostForm("u_name")

	stmt, err := db.Prepare("insert into order_new (o_addr,o_code,o_datetime,o_detail,o_other,o_pay_status,o_pay_type,o_receiver,o_status,o_tel,o_total_money,u_id,u_name) values(?,?,?,?,?,?,?,?,?,?,?,?,?);")
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	_, err = stmt.Exec(o_addr, o_code, o_datetime, o_detail, o_other, o_pay_status, o_pay_type, o_receiver, o_status, o_tel, o_total_money, u_id, u_name)

	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}

	buffer.WriteString(o_code)
	buffer.WriteString("  ")
	buffer.WriteString(o_addr)
	buffer.WriteString("  ")
	buffer.WriteString(o_datetime)
	buffer.WriteString("  ")
	buffer.WriteString(o_detail)
	buffer.WriteString("  ")
	buffer.WriteString(o_other)
	buffer.WriteString("  ")
	buffer.WriteString(o_pay_status)
	buffer.WriteString("  ")
	buffer.WriteString(o_pay_type)
	buffer.WriteString("  ")
	buffer.WriteString(o_receiver)
	buffer.WriteString("  ")
	buffer.WriteString(o_status)
	buffer.WriteString("  ")
	buffer.WriteString(o_tel)
	buffer.WriteString("  ")
	buffer.WriteString(o_total_money)
	buffer.WriteString("  ")
	buffer.WriteString(u_id)
	buffer.WriteString("  ")
	buffer.WriteString(u_name)
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
func PutOrder(c *gin.Context) {
	var buffer bytes.Buffer
	id := c.Param("id")
	o_addr := c.PostForm("o_addr")
	o_code := c.PostForm("o_code")
	o_datetime := c.PostForm("o_datetime")
	o_detail := c.PostForm("o_detail")
	o_other := c.PostForm("o_other")
	o_pay_status := c.PostForm("o_pay_status")
	o_pay_type := c.PostForm("o_pay_type")
	o_receiver := c.PostForm("o_receiver")
	o_status := c.PostForm("o_status")
	o_tel := c.PostForm("o_tel")
	o_total_money := c.PostForm("o_total_money")
	u_id := c.PostForm("u_id")
	u_name := c.PostForm("u_name")

	stmt, err := db.Prepare("update order_new set o_addr=?,o_code=?,o_datetime=?,o_detail=?,o_other=?,o_pay_status=?,o_pay_type=?,o_receiver=?,o_status=?,o_tel=?,o_total_money=?,u_id=?,u_name=? where o_id= ?;")
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	_, err = stmt.Exec(o_addr, o_code, o_datetime, o_detail, o_other, o_pay_status, o_pay_type, o_receiver, o_status, o_tel, o_total_money, u_id, u_name, id)
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}

	// Fastest way to append strings
	buffer.WriteString(o_code)
	buffer.WriteString("  ")
	buffer.WriteString(o_addr)
	buffer.WriteString("  ")
	buffer.WriteString(o_datetime)
	buffer.WriteString("  ")
	buffer.WriteString(o_detail)
	buffer.WriteString("  ")
	buffer.WriteString(o_other)
	buffer.WriteString("  ")
	buffer.WriteString(o_pay_status)
	buffer.WriteString("  ")
	buffer.WriteString(o_pay_type)
	buffer.WriteString("  ")
	buffer.WriteString(o_receiver)
	buffer.WriteString("  ")
	buffer.WriteString(o_status)
	buffer.WriteString("  ")
	buffer.WriteString(o_tel)
	buffer.WriteString("  ")
	buffer.WriteString(o_total_money)
	buffer.WriteString("  ")
	buffer.WriteString(u_id)
	buffer.WriteString("  ")
	buffer.WriteString(u_name)
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
func DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("delete from order_new where o_id= ?;")
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
