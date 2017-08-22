package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ClassifyCommodity struct {
	ClassifyName   string      `json:"name"`
	ClassifyMenuID int         `json:"mid"`
	Commoditys     []Commodity `json:"commoditys"`
}

func API_GetClassifyCommoditys(c *gin.Context) {
	var classifyCommoditys = make([]ClassifyCommodity, 0)
	var menu_classifys = make([]Menu_classify, 0)
	var menu_classify Menu_classify
	rows, err := db.Query("select m_id,m_other,m_datetime,m_name,m_status from menu_classify where m_status='上架';")
	if err != nil {
		log.Println(err.Error())
		SetErrorRespones(c, err.Error())
		return
	}
	for rows.Next() {
		err = rows.Scan(&menu_classify.M_Id, &menu_classify.M_other, &menu_classify.M_datetime, &menu_classify.M_name, &menu_classify.M_status)
		menu_classifys = append(menu_classifys, menu_classify)
		if err != nil {
			log.Println(err.Error())
			rows.Close()
			SetErrorRespones(c, err.Error())
			break
		}
	}
	defer rows.Close()
	for _, menu := range menu_classifys {
		var classifyCommodity ClassifyCommodity
		classifyCommodity.ClassifyName = menu.M_name
		classifyCommodity.ClassifyMenuID = menu.M_Id
		if data, err := getCommodityByClassifyMenu(menu.M_Id); err != nil {
			SetErrorRespones(c, err.Error())
			return
		} else {
			classifyCommodity.Commoditys = data
		}
		classifyCommoditys = append(classifyCommoditys, classifyCommodity)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":              httpOK,
		"error":               "",
		"classify_commoditys": classifyCommoditys,
	})

}
func getCommodityByClassifyMenu(m_id int) ([]Commodity, error) {
	var commoditys = make([]Commodity, 0)
	var commodity Commodity
	rows, err := db.Query("select c_id,c_current_price,c_datetime,c_detail,c_img,c_sku,c_name,c_other,c_primary_price,c_sales_num,c_status,c_stock,m_id from commodity where m_id = ? and c_status='上架';", m_id)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(&commodity.C_Id, &commodity.C_current_price, &commodity.C_datetime, &commodity.C_detail, &commodity.C_img, &commodity.C_sku, &commodity.C_name, &commodity.C_other, &commodity.C_primary_price, &commodity.C_sales_num, &commodity.C_status, &commodity.C_stock, &commodity.M_id)

		commoditys = append(commoditys, commodity)
		if err != nil {
			log.Println(err.Error())
			rows.Close()
			return nil, err
		}
	}
	defer rows.Close()
	return commoditys, nil

}
