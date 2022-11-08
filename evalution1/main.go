package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	//"os"
	//"path/filepath"

	//"github.com/joho/godotenv"

	"github.com/kelseyhightower/envconfig"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Post struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

/*type RawConfig struct {
IP   string `envconfig:"MYSQL_IP"`
Host string `envconfig:"MYSQL_HOST"`
Port string `envconfig:"MYSQL_PORT,default=3306"`
//TLP      string `envconfig:"MYSQL_TLP,optional"`
User     string `envconfig:"MYSQL_USER"`
Password string `envconfig:"MYSQL_PASSWORD"`
Database string `envconfig:"MYSQL_DATABASE"`*/

type KafkaConfig struct {
	User     string `envconfig:"KAFKA_USER" required:"true"`
	Password string `envconfig:"KAFKA_PASSWORD" required:"true"`
	//Port     string `envconfig:"KAFKA_PORT" required:"true"`
	//IP       string `envconfig:"KAFKA_IP" required:"true"`
	Database string `envconfig:"KAFKA_DATABASE" required:"true"`
}

type ORMConfig struct {
	MySQL mysql.Config
	GORM  gorm.Config
}

var db *sql.DB
var err error

func main() {

	//var ormCfg ORMConfig
	/*b, err = sql.Open("mysql", "root:alanjino@tcp(127.0.0.1:3306)/films")

	//dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s", rawCfg.User, rawCfg.Password,rawCfg.Host, rawCfg.Database)

	//dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s",rawCfg.User ,rawCfg.Password , rawCfg.Host +  rawCfg.Database  "?charset=utf8")
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s", rawCfg.User, rawCfg.Password, "tcp", rawCfg.IP, rawCfg.Port, rawCfg.Database)
	fmt.Printf("%v", dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	if err = db.Ping(); err != nil {
		panic(err.Error())
	}
	// ORM layer
	ormCfg.MySQL.Conn = db
	/*gormDialector := mysql.New(ormCfg.MySQL)
	gormCfg := ormCfg.GORM
	gormDB, e := gorm.Open(gormDialector, &gormCfg)
	if e != nil {
		return nil, e
	}

	close()*/
	//NewMySQLClient(rawCfg, ormCfg)

	kafkaConfig := &KafkaConfig{}
	err = envconfig.Process("KAFKA", kafkaConfig)
	if err != nil {
		fmt.Printf("kafka config parameter error: %v", err)
	}
	fmt.Printf("kafka config: %+v\n", kafkaConfig)
	//dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s", rawCfg.User, rawCfg.Password, "tcp", rawCfg.IP, rawCfg.Port, rawCfg.Database)
	//dbHost := os.Getenv("KAFKA_PORT")
	dbUsername := os.Getenv("KAFKA_USER")
	dbPassword := os.Getenv("KAFKA_PASSWORD")
	dbname := os.Getenv("KAFKA_DATABASE")
	//dbIP := os.Getenv("KAFKA_IP")
	dsn := dbUsername + ":" + dbPassword + "@tcp(127.0.0.1:3306)/" + dbname
	fmt.Printf("%v", dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	//fmt.Printf("kafka config: %+v\n", kafkaConfig)
	router := mux.NewRouter()
	router.HandleFunc("/", get).Methods("GET")
	router.HandleFunc("/user", getPosts).Methods("GET")
	router.HandleFunc("/useradd", createPost).Methods("POST")
	router.HandleFunc("/user/{id}", getPost).Methods("GET")
	//router.HandleFunc("/health", HealthCheckHandler)
	router.HandleFunc("/health", healthApi)
	http.ListenAndServe(":8000", router)
}
func get(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "welocme")
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var posts []Post
	result, err := db.Query("SELECT id, title from posts")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var post Post
		err := result.Scan(&post.ID, &post.Title)
		if err != nil {
			panic(err.Error())
		}
		posts = append(posts, post)
	}
	json.NewEncoder(w).Encode(posts)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stmt, err := db.Prepare("INSERT INTO posts(title) VALUES(?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	title := keyVal["title"]
	_, err = stmt.Exec(title)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New post was created")
}

func getPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT id, title FROM posts WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var post Post
	for result.Next() {
		err := result.Scan(&post.ID, &post.Title)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(post)
}

/*func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"responce": 200}`)

}*/

func healthApi(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://localhost:8000/user")
	if err != nil {
		log.Fatalf("HTTP GET request failed, %v\n", err)
	}
	fmt.Fprintf(w, "<h1>Health check is done  %v</h1>", resp.Status)
	defer resp.Body.Close()
	fmt.Println("Response status:", resp.Status)
	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Body read failed: %v\n", err)
	}
}

/*func NewMySQLClient(rawCfg RawConfig, ormCfg ORMConfig) (*gorm.DB, error) {
	// Raw layer
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s", rawCfg.User, rawCfg.Password, "tcp", rawCfg.IP, rawCfg.Port, rawCfg.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	// ORM layer
	ormCfg.MySQL.Conn = db
	gormDialector := mysql.New(ormCfg.MySQL)
	gormCfg := ormCfg.GORM
	gormDB, e := gorm.Open(gormDialector, &gormCfg)
	if e != nil {
		return nil, e
	}

	return gormDB, nil
}*/
