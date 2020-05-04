package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/BlockTeam4Boys/digitaldocs/internal/config"
	"github.com/BlockTeam4Boys/digitaldocs/internal/model"
	"github.com/BlockTeam4Boys/digitaldocs/internal/role"
)

var (
	DB *gorm.DB
)

func migrate() {
	rolesExists := DB.HasTable(&model.Role{})
	DB.AutoMigrate(&model.Role{}, &model.Organization{})
	DB.AutoMigrate(&model.User{}).
		AddForeignKey("organization_id", "organizations(id)", "RESTRICT", "CASCADE").
		AddForeignKey("role_id", "roles(id)", "RESTRICT", "CASCADE")
	DB.AutoMigrate(&model.OrganizationsAgreements{}).
		AddForeignKey("from_id", "organizations(id)", "CASCADE", "CASCADE").
		AddForeignKey("to_id", "organizations(id)", "CASCADE", "CASCADE")
	if !rolesExists {
		DB.Create(&model.Role{ID: role.Admin, Name: "Admin"})
		DB.Create(&model.Role{ID: role.User, Name: "User"})
	}
}

func init() {
	cfgDB := config.Database{}
	err := cleanenv.ReadConfig("configs/db.yaml", &cfgDB)
	if err != nil {
		panic(err)
	}
	cfgJWT := config.JWT{}
	err = cleanenv.ReadConfig("configs/jwt.yaml", &cfgJWT)
	if err != nil {
		panic(err)
	}
	TokenCookieName = cfgJWT.Token.CookieName
	TokenTTL = cfgJWT.Token.TTL
	TokenRefreshTime = cfgJWT.Token.RefreshTime
	JWTKey = []byte(cfgJWT.Key)
	argsMsg := "host=%v port=%v user=%v dbname=%v password=%v sslmode=disable"
	argsDB := fmt.Sprintf(argsMsg, cfgDB.Host, cfgDB.Port, cfgDB.User, cfgDB.Name, cfgDB.Password)
	DB, err = gorm.Open("postgres", argsDB)
	if err != nil {
		panic(err)
	}
	migrate()
}

func main() {
	defer DB.Close()
	r := mux.NewRouter()
	r.HandleFunc("/api/auth", AuthHandler).Methods("POST")
	r.HandleFunc("/api/refresh", RefreshHandler)
	r.HandleFunc("/api/request", RequestHandler).Methods("POST")
	http.ListenAndServe(":8080", r)
}
