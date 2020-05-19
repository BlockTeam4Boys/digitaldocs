package main

import (
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/BlockTeam4Boys/digitaldocs/internal/config"
	"github.com/BlockTeam4Boys/digitaldocs/internal/middleware"
	"github.com/BlockTeam4Boys/digitaldocs/internal/model"
	"github.com/BlockTeam4Boys/digitaldocs/internal/role"
	sr "github.com/BlockTeam4Boys/digitaldocs/internal/session/repository"
	ur "github.com/BlockTeam4Boys/digitaldocs/internal/user/repository"
)

var (
	DB          *gorm.DB
	RDB         *redis.Client
	SessionRepo *sr.RedisTokenRepository
	UserRepo    *ur.PostgresUserRepository
	JWTInfo     config.JWT
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
	cfgDB := config.Postgres{}
	err := cleanenv.ReadConfig("configs/postgres.yaml", &cfgDB)
	if err != nil {
		panic(err)
	}
	cfgRDB := config.Redis{}
	err = cleanenv.ReadConfig("configs/redis.yaml", &cfgRDB)
	if err != nil {
		panic(err)
	}
	cfgJWT := config.JWT{}
	err = cleanenv.ReadConfig("configs/jwt.yaml", &cfgJWT)
	if err != nil {
		panic(err)
	}
	JWTInfo = cfgJWT
	argsMsg := "host=%v port=%v user=%v dbname=%v password=%v sslmode=disable"
	argsDB := fmt.Sprintf(argsMsg, cfgDB.Host, cfgDB.Port, cfgDB.User, cfgDB.Name, cfgDB.Password)
	DB, err = gorm.Open("postgres", argsDB)
	if err != nil {
		panic(err)
	}
	migrate()
	RDB = redis.NewClient(&redis.Options{
		Addr:     cfgRDB.Addr,
		Password: cfgRDB.Password,
		DB:       cfgRDB.DB,
	})
	err = RDB.Ping().Err()
	if err != nil {
		panic(err)
	}
	SessionRepo = sr.NewRedisSessionRepository(RDB, JWTInfo.AccessToken.TTL, JWTInfo.RefreshToken.TTL)
	UserRepo = ur.NewPostgresUserRepository(DB)
}

func main() {
	defer DB.Close()
	authentication := middleware.NewAuthMiddleware(SessionRepo, JWTInfo.AccessToken.CookieName)
	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))

	api := r.PathPrefix("/api").Subrouter()

	auth := api.Path("/auth")
	auth.HandlerFunc(AuthHandler)

	refresh := api.Path("/refresh")
	refresh.HandlerFunc(RefreshHandler)

	access := api.PathPrefix("/").Subrouter()
	access.HandleFunc("/request", RequestHandler)
	access.Use(authentication.Middleware)

	http.ListenAndServe(":8080", r)
}
