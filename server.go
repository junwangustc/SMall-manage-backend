package main

import (
	"database/sql"

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

		//======END  ADD ROUTER
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
