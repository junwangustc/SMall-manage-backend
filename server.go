package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/junwangustc/ustclog"
)

var db *sql.DB

const httpOK = "OK"
const httpError = "ERROR"

type Server struct {
	Cfg    *Config
	Db     *sql.DB
	Logger *log.Logger
}

func SetErrorRespones(c *gin.Context, errorstring string) {
	c.JSON(http.StatusOK, gin.H{
		"status": httpError,
		"error":  errorstring,
		"result": nil,
		"count":  0,
	})
}

func (s *Server) Open() error {
	var err error
	s.Db, err = sql.Open("mysql", "root:12345@tcp(127.0.0.1:3306)/bling?parseTime=true")
	if err != nil {
		log.Println(err)
		return err
	}
	err = s.Db.Ping()
	if err != nil {
		log.Println(err)
		return err
	}
	db = s.Db
	go func() {
		router := gin.Default()
		//=======ADD ROUTER
		router.GET("/api/v1/menu_classify/:id", func(c *gin.Context) {
			GetMenu_classify(c)
		})
		router.GET("/api/v1/menu_classifys", func(c *gin.Context) {
			GetMenu_classifys(c)
		})
		router.GET("/api/v1/menu_classifys/total", func(c *gin.Context) {
			GetCountMenu_classifys(c)
		})
		router.GET("/api/v1/menu_classifys/page/:pageid", func(c *gin.Context) {
			GetMenu_classifysByPage(c)
		})
		router.POST("/api/v1/menu_classify", func(c *gin.Context) {
			PostMenu_classify(c)
		})
		router.PUT("/api/v1/menu_classify/:id", func(c *gin.Context) {
			PutMenu_classify(c)
		})
		router.DELETE("/api/v1/menu_classify/:id", func(c *gin.Context) {
			DeleteMenu_classify(c)
		})

		router.GET("/api/v1/flag_classify/:id", func(c *gin.Context) {
			GetFlag_classify(c)
		})
		router.GET("/api/v1/flag_classifys", func(c *gin.Context) {
			GetFlag_classifys(c)
		})
		router.GET("/api/v1/flag_classifys/total", func(c *gin.Context) {
			GetCountFlag_classifys(c)
		})
		router.GET("/api/v1/flag_classifys/page/:pageid", func(c *gin.Context) {
			GetFlag_classifysByPage(c)
		})
		router.POST("/api/v1/flag_classify", func(c *gin.Context) {
			PostFlag_classify(c)
		})
		router.PUT("/api/v1/flag_classify/:id", func(c *gin.Context) {
			PutFlag_classify(c)
		})
		router.DELETE("/api/v1/flag_classify/:id", func(c *gin.Context) {
			DeleteFlag_classify(c)
		})

		router.GET("/api/v1/commodity/:id", func(c *gin.Context) {
			GetCommodity(c)
		})
		router.GET("/api/v1/commoditys", func(c *gin.Context) {
			GetCommoditys(c)
		})
		router.GET("/api/v1/commoditys/total", func(c *gin.Context) {
			GetCountCommoditys(c)
		})
		router.GET("/api/v1/commoditys/page/:pageid", func(c *gin.Context) {
			GetCommoditysByPage(c)
		})
		router.POST("/api/v1/commodity", func(c *gin.Context) {
			PostCommodity(c)
		})
		router.PUT("/api/v1/commodity/:id", func(c *gin.Context) {
			PutCommodity(c)
		})
		router.DELETE("/api/v1/commodity/:id", func(c *gin.Context) {
			DeleteCommodity(c)
		})

		router.GET("/api/v1/order/:id", func(c *gin.Context) {
			GetOrder(c)
		})
		router.GET("/api/v1/orders", func(c *gin.Context) {
			GetOrders(c)
		})
		router.GET("/api/v1/orders/total", func(c *gin.Context) {
			GetCountOrders(c)
		})
		router.GET("/api/v1/orders/page/:pageid", func(c *gin.Context) {
			GetOrdersByPage(c)
		})
		router.POST("/api/v1/order", func(c *gin.Context) {
			PostOrder(c)
		})
		router.PUT("/api/v1/order/:id", func(c *gin.Context) {
			PutOrder(c)
		})
		router.DELETE("/api/v1/order/:id", func(c *gin.Context) {
			DeleteOrder(c)
		})

		router.GET("/api/v1/user/:id", func(c *gin.Context) {
			GetUser(c)
		})
		router.GET("/api/v1/users", func(c *gin.Context) {
			GetUsers(c)
		})
		router.GET("/api/v1/users/total", func(c *gin.Context) {
			GetCountUsers(c)
		})
		router.GET("/api/v1/users/page/:pageid", func(c *gin.Context) {
			GetUsersByPage(c)
		})
		router.POST("/api/v1/user", func(c *gin.Context) {
			PostUser(c)
		})
		router.PUT("/api/v1/user/:id", func(c *gin.Context) {
			PutUser(c)
		})
		router.DELETE("/api/v1/user/:id", func(c *gin.Context) {
			DeleteUser(c)
		})

		router.GET("/api/v1/user_addr/:id", func(c *gin.Context) {
			GetUser_addr(c)
		})
		router.GET("/api/v1/user_addrs", func(c *gin.Context) {
			GetUser_addrs(c)
		})
		router.GET("/api/v1/user_addrs/total", func(c *gin.Context) {
			GetCountUser_addrs(c)
		})
		router.GET("/api/v1/user_addrs/page/:pageid", func(c *gin.Context) {
			GetUser_addrsByPage(c)
		})
		router.POST("/api/v1/user_addr", func(c *gin.Context) {
			PostUser_addr(c)
		})
		router.PUT("/api/v1/user_addr/:id", func(c *gin.Context) {
			PutUser_addr(c)
		})
		router.DELETE("/api/v1/user_addr/:id", func(c *gin.Context) {
			DeleteUser_addr(c)
		})

		//======END  ADD ROUTER

		//===========APP Restful API

		router.GET("/app-api/v1/classify_commoditys", func(c *gin.Context) {
			API_GetClassifyCommoditys(c)
		})
		router.Run(":3000")
	}()
	return nil
}
func (s *Server) Close() {
	if s.Db != nil {
		s.Db.Close()
	}
}

func NewServer(c *Config, l *log.Logger) (srv *Server, err error) {
	s := &Server{Cfg: c, Logger: l, Db: nil}
	return s, nil
}
