package main

import (
	"github.com/asveg/category/common"
	"github.com/asveg/category/domain/repository"
	service2 "github.com/asveg/category/domain/service"
	"github.com/asveg/category/handler"
	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"strconv"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	category "github.com/asveg/category/proto/category"
)


func main() {
	//配置中心，prefix是consoul中配置mysql的目录
	consulConfig,err := common.GetConsulConfig("192.168.10.168",8500,"/micro/config")
	if err !=nil {
		log.Error(err)
	}
	//注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"192.168.10.168:8500",
		}
	})

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.category"),
		micro.Version("latest"),
		//这里设置地址和需要暴露的端口
		micro.Address("0.0.0.0:8082"),
		//添加consul 作为注册中心
		micro.Registry(consulRegistry),
	)

	//获取mysql配置,路径中不带前缀
	mysqlInfo := common.GetMysqlFromConsul(consulConfig,"mysql")

	//连接数据库
	intPort := mysqlInfo.Port
	strPort:= strconv.FormatInt(intPort,10)
	db,err := gorm.Open("mysql",mysqlInfo.User+":"+mysqlInfo.Pwd+"@("+mysqlInfo.Host+":"+strPort+")/"+mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err !=nil {
		log.Error(err)
	}
	defer db.Close()
	//禁止复表
	db.SingularTable(true)


	//rp := repository.NewCategoryRepository(db)
	//rp.InitTable()

	// Initialise service
	service.Init()

	categoryDataService := service2.NewCategoryDataService(repository.NewCategoryRepository(db))

	 err = category.RegisterCategoryHandler(service.Server(),&handler.Category{CategoryDataService:categoryDataService})
	 if  err != nil {
	 	log.Error(err)
	 }

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
