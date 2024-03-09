package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/evalphobia/logrus_sentry"
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"golang.org/x/crypto/ssh"
	"gopkg.in/boj/redistore.v1"

	"utils"
)

// postgreSQL database connection settings.
const (
	PGDBHost   = "localhost"
	PGDBName   = "dittofi"
	PGDBUser   = "dittofi"
	PGDBPass   = "dittofi"
	PGDBSchema = "app_2187"
)

const (
	IntegrityConstraintViolationClass pq.ErrorClass = "23"
	IntegrityConstraintViolation      pq.ErrorCode  = "23000"
	RestrictViolation                 pq.ErrorCode  = "23001"
	NotNullViolation                  pq.ErrorCode  = "23502"
	ForeignKeyViolation               pq.ErrorCode  = "23503"
	UniqueViolation                   pq.ErrorCode  = "23505"
	CheckViolation                    pq.ErrorCode  = "23514"
	ExclusionViolation                pq.ErrorCode  = "23P01"
	UndefinedColumn                   pq.ErrorCode  = "42703"
)

// send email with SMTP settings.
const (
	SMTPHost     = "172.17.0.1"
	SMTPPort     = 7587
	SMTPUsername = "app2187"
	SMTPPassword = "1672758222303704941"
)

// flag values to indicate application runtime mode.
const (
	ProductionMode = "production"
	TestingMode    = "testing"
)

// Sentry (https://sentry.io/) 3rd party application monitoring service.
// Data Source Name (DSN) value to send event logs.
const sentryDSN = ""

// Redis (https://redis.io/) open source in-memory data store settings.
// Set up to manage login sessions with HTTP headers.
// Set up to manage pagination data.
const (
	RediStoreMaxIdle = 1
	RediStoreNetwork = "tcp"
	RediStoreAddress = ":6379"

	RediStorePassword = ""

	RediStoreAuthenticationKey = "abc123"

	RediStoreEncryptionKey = ""

	LoginSessionName      = "session-key"
	SessionValueUserIDKey = "user_id"
	InternalHeaderUserID  = "ih_user_id"

	RedisMaxActive    = 8
	RedisPaginationDB = 1

	PaginateTimeLimitSeconds    = 3600
	PaginateDataMaxWaitSeconds  = 10
	PaginateDataFieldNumTotal   = "num"
	PaginateDataFieldNumPerPage = "num_per_page"
	PaginateDataError           = "error"
	PaginateDataIsDB            = "is_db"

	// database specific pagination
	PaginationSchema        = "pagination"
	PaginationTableNameFmt  = "p-%s"
	PaginationTableRowIdCol = "row_id_123abc"
)

// directory to store files.
const FileSystemRoot = "/srv/data"

// base URL to append attachments for access
const UploadsBaseURL = "dittofi.com/2187/uploads"

// store runtime mode value.
var runtimeMode string

// store http server listening port.
var port string

// store http server host.
var host string

// store access to postgreSQL database.
var pg *sqlx.DB

// store access to log system.
var log *logrus.Logger

// store access to logged in sessions.
var loginSessions *redistore.RediStore

// store access to web and email templates.
var templates *template.Template

// store access to filesystem.
var fs *FileSystem

// store access to paginated data.
var paginationDataStorage *redis.Pool

var Currency1Variable = "usd"

var BookingApprovedMessage1Variable = "Congrats! You have approved a booking!"

var ApplicationUrl1Variable = "https://dittofi.com/2187#/"

var StripeExpressConnect1Variable = "https://connect.stripe.com/express/oauth/authorize?response_type=code&client_id=ca_M4PB8yaPIcafDgQL1z2izVR1boFI465U"

var Approved1Variable = "approved"

var StripeSecretKey1Variable = "Bearer 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"

var StripePaymentMode1Variable = "payment"

var ApplicationFormUrlEncoded1Variable = "application/x-www-form-urlencoded"

var ListingUpdated1Variable = "Congrats! You have updated your listing!"

var YouHaveANewMessage1Variable = "You have 1 new message! Click \"activities\" in your dashboard to view!"

var StripeCheckOut1Variable = "https://api.stripe.com/v1/checkout/sessions"

// function setting up and running http server.
func main() {
	flag.StringVar(&port, "port", "8000", "the port the http server listens on")
	flag.StringVar(&host, "host", "0.0.0.0", "the host the http server listens on")
	flag.StringVar(&runtimeMode, "mode", TestingMode, "'testing' disables logging to sentry, 'production' enables logging to sentry")
	flag.Parse()

	go func() {
		cronClient := cron.New(cron.WithParser(cron.NewParser(cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor)), cron.WithLocation(time.UTC))
		cronClient.Run()
	}()

	log = setUpLogger()
	pg = connectPostgres()
	templates = setupTemplates()
	loginSessions = configRediStore()
	defer loginSessions.Close()
	fs = NewFileSystem(FileSystemRoot)
	paginationDataStorage = configPaginationStorage()

	webSocketsMap = make(map[interface{}]*WebsocketConnection)
	webSocketsLock = &sync.RWMutex{}

	router := configRouter()

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.UseHandler(router)
	n.Run(fmt.Sprintf("%s:%s", host, port))
}

// setUpLoggers sets up access to the Sentry or system logger.
func setUpLogger() *logrus.Logger {
	log := logrus.New()

	// log to Sentry
	if runtimeMode == ProductionMode {
		hook, err := logrus_sentry.NewSentryHook(sentryDSN, []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		})

		if err != nil {
			log.Fatal(err)
		} else {
			hook.StacktraceConfiguration.Enable = true
			hook.StacktraceConfiguration.Skip = 0
			hook.StacktraceConfiguration.Context = 2
			log.Hooks.Add(hook)
		}
	}

	return log
}

// connectPostgres sets up access to postgreSQL database.
func connectPostgres() *sqlx.DB {
	if pg != nil {
		// already connected
		return pg
	}

	log.WithFields(logrus.Fields{
		"db_host": PGDBHost,
	}).Info("Connect to PostgreSQL: ...")

	connString := fmt.Sprintf("dbname=%s user=%s password=%s host=%s sslmode=disable search_path=%s",
		PGDBName, PGDBUser, PGDBPass, PGDBHost, PGDBSchema)

	pg = sqlx.MustConnect("postgres", connString)
	pg.Exec(fmt.Sprintf("set search_path='%s'", PGDBSchema))
	pg.SetMaxIdleConns(1)
	pg.SetMaxOpenConns(32)
	log.Info("... Connected to PostgreSQL")

	return pg
}

// setupTemplates sets up access to web and email templates.
func setupTemplates() *template.Template {
	cleanRoot := filepath.Clean("templates")
	pfx := len(cleanRoot) + 1
	root := template.New("")
	err := filepath.Walk(cleanRoot, func(path string, info os.FileInfo, e1 error) error {
		if e1 != nil {
			return e1
		}

		if !info.IsDir() && strings.HasSuffix(path, ".tmp") {
			b, e2 := ioutil.ReadFile(path)
			if e2 != nil {
				return e2
			}

			name := path[pfx:]
			t := root.New(name)
			_, e2 = t.Parse(string(b))
			if e2 != nil {
				return e2
			}
		}

		return nil
	})

	if err != nil {
		return nil
	}

	return root
}

// configRediStore sets up access to Redis for login sessions.
func configRediStore() *redistore.RediStore {
	var keyPairs [][]byte
	if len(RediStoreAuthenticationKey) < 1 {
		log.Fatal("no authentication key set for RediStore")
	} else {
		keyPairs = append(keyPairs, []byte(RediStoreAuthenticationKey))
	}

	if l := len(RediStoreEncryptionKey); l > 0 {
		if l == 16 || l == 24 || l == 32 {
			keyPairs = append(keyPairs, []byte(RediStoreEncryptionKey))
		} else {
			log.Fatal("wrong length for encryption key set for RediStore")
		}
	}

	store, err := redistore.NewRediStore(RediStoreMaxIdle, RediStoreNetwork, RediStoreAddress, RediStorePassword, keyPairs...)
	if err != nil {
		log.Fatal(err)
	}

	store.Options.Path = "/"
	store.Options.SameSite = http.SameSiteNoneMode
	store.Options.Secure = true

	return store
}

// addLoginSession stores session and user id of logged in user into Redis.
func addLoginSession(w *http.ResponseWriter, r *http.Request, userID int) (err error) {
	// Get a session.
	var loginSession *sessions.Session
	loginSession, err = loginSessions.Get(r, LoginSessionName)
	if err != nil {
		return
	}

	// Add a value.
	loginSession.Values[SessionValueUserIDKey] = userID

	// Save.
	err = sessions.Save(r, *w)

	return
}

// removeLoginSession removes stored session and user id from Redis.
func removeLoginSession(w *http.ResponseWriter, r *http.Request) (err error) {
	// Get the session.
	var loginSession *sessions.Session
	loginSession, err = loginSessions.Get(r, LoginSessionName)
	if err != nil {
		return
	}

	// Delete session.
	loginSession.Options.MaxAge = -1

	// Save.
	err = sessions.Save(r, *w)

	return
}

// getLoginSessionUserID retrieves the user id of logged in user session from Redis.
// Attempts to look for internal HTTP request header containing user id already set.
func getLoginSessionUserID(r *http.Request) (userID int, err error) {
	// attempt to see if internal header set already otherwise attempt to find
	var userIDStr = r.Header.Get(InternalHeaderUserID)
	if len(userIDStr) == 0 {
		// Get a session.

		var loginSession *sessions.Session
		loginSession, err = loginSessions.Get(r, LoginSessionName)
		if err != nil {
			return
		}

		if iUserID, ok := loginSession.Values[SessionValueUserIDKey]; !ok {
			err = fmt.Errorf("no user id found")
			return
		} else if userID, ok = iUserID.(int); !ok {
			err = fmt.Errorf("unexpected user id type not int")
			return
		}
	} else {
		userID, err = strconv.Atoi(userIDStr)
	}

	return
}

// RequireLogin wraps handler to check login before processing request.
func RequireLogin(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getLoginSessionUserID(r)
		if err != nil {
			utils.JSON(w, http.StatusUnauthorized, utils.Response{Error: &utils.ValError{Code: utils.ErrCodeUnAuthorized, Message: err.Error(), Param: "password"}})
			return
		}

		r.Header.Set(InternalHeaderUserID, fmt.Sprint(userID))
		handler(w, r)
	}
}

// manage concurrent access to files.
type FileSystem struct {
	// control concurrent access to FileLocks
	sync.Mutex

	// control concurrent read/write access to files
	FileLocks map[string]*sync.RWMutex
	Root      string

	TempDir string
}

// GetFileLock obtain lock to access file.
func (fs *FileSystem) GetFileLock(path string) (lock *sync.RWMutex) {
	fs.Lock()
	defer fs.Unlock()
	lock, ok := fs.FileLocks[path]
	if !ok {
		lock = &sync.RWMutex{}
		fs.FileLocks[path] = lock
	}

	return
}

// GetFileContentType returns the content type of a file.
func GetFileContentType(f File) (contentType string, err error) {
	// to sniff the content type only the first 512 bytes are used.
	buf := make([]byte, 512)

	_, err = f.Read(buf)
	if err != nil {
		return
	}

	_, err = f.File.Seek(0, io.SeekStart)
	if err != nil {
		return
	}

	contentType = http.DetectContentType(buf)

	return
}

// GetFile returns access to a file.
func (fs *FileSystem) GetFile(path string, flag int) (f File, err error) {
	lock := fs.GetFileLock(path)
	lock.RLock()
	defer lock.RUnlock()

	// generate full path for final store location
	var fullPath = path
	if len(fs.Root) > 0 {
		fullPath = filepath.Join(fs.Root, path)
	}

	// Open file.
	var srcFile *os.File
	srcFile, err = os.OpenFile(fullPath, flag, 0755)
	if err != nil {
		return
	} else {
		f.File = srcFile
	}

	f.FilePath = path

	return
}

// SetFile stores the file on to disk.
func (fs *FileSystem) SetFile(f File, path string, overwrite bool) (err error) {
	lock := fs.GetFileLock(path)
	lock.Lock()
	defer lock.Unlock()

	// generate full path for final store location
	var fullPath = path
	if len(fs.Root) > 0 {
		fullPath = filepath.Join(fs.Root, path)
	}

	// create any missing directories
	err = os.MkdirAll(filepath.Dir(fullPath), 0755)
	if err != nil {
		return
	}

	// set file flags to allow overwriting existing file
	var fileFlag = os.O_WRONLY | os.O_CREATE
	if overwrite {
		fileFlag |= os.O_TRUNC
	} else {
		fileFlag |= os.O_EXCL
	}

	// open/create file
	var file *os.File
	file, err = os.OpenFile(fullPath, fileFlag, 0666)
	if !overwrite && os.IsExist(err) {
		err = fmt.Errorf(`cannot overwrite existing file "%s"`, path)
		return
	} else if err != nil {
		return
	}
	defer file.Close()

	// copy data in temporary file to new file
	if f.File != nil {
		_, err = f.File.Seek(0, 0)
		if err != nil {
			return
		}

		_, err = io.Copy(file, f.File)
		if err != nil {
			return
		}
	}

	return
}

// SetTempFile temporary stores data into a file.
func (fs *FileSystem) SetTempFile(data io.Reader) (f File, err error) {
	var filePrefix string
	if data == nil {
		filePrefix = "copy"
	} else {
		filePrefix = "upload"
	}

	// create temporary file
	var osFile *os.File
	osFile, err = ioutil.TempFile(fs.TempDir, fmt.Sprintf("%s-*", filePrefix))
	if err != nil {
		return
	}

	// copy data if available
	if data != nil {
		_, err = io.Copy(osFile, data)
		if err != nil {
			return
		}

		// reset seek to beginning of file
		_, err = osFile.Seek(0, 0)
		if err != nil {
			return
		}
	}

	f.File = osFile
	_, f.FilePath = filepath.Split(osFile.Name())

	return
}

// NewFileSystem sets up access to file system.
func NewFileSystem(root string) (fs *FileSystem) {
	if len(root) == 0 || root == "/" {
		log.Fatalf(`invalid file system root "%s" set`, root)
	}

	tempDir, err := ioutil.TempDir("", "temp-file-dir*")
	if err != nil {
		log.Fatal(err)
	}

	fs = &FileSystem{
		Mutex:     sync.Mutex{},
		FileLocks: make(map[string]*sync.RWMutex, 0),
		Root:      root,
		TempDir:   tempDir,
	}

	return
}

// PublicKey generates authentication method from private key for SSH.
func PublicKey(privateKey string) (*ssh.AuthMethod, error) {
	signer, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return nil, err
	}

	authMethod := ssh.PublicKeys(signer)
	return &authMethod, nil
}

//Generates JWT Token for Authentication
func GenerateJwtToken(signingMethod, privateKey string, payload interface{}, ttl int64) (token string, err error) {
	var key interface{}
	if signingMethod == "jwt.SigningMethodRS256" {
		key, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
		if err != nil {
			return "", fmt.Errorf("create: parse key: %w", err)
		}
	} else if signingMethod == "jwt.SigningMethodHS256" {
		key = []byte(privateKey)
	} else {
		return "", fmt.Errorf("create: invalid signing method: %w", err)
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["payload"] = payload                                    // Our custom data.
	claims["exp"] = now.Add(time.Duration(ttl) * time.Hour).Unix() // The expiration time after which the token must be disregarded.
	claims["iat"] = now.Unix()                                     // The time at which the token was issued.
	claims["nbf"] = now.Unix()                                     // The time before which the token must be disregarded.

	if signingMethod == "jwt.SigningMethodRS256" {
		token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	} else if signingMethod == "jwt.SigningMethodHS256" {
		token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(key)
	} else {
		return "", fmt.Errorf("create: invalid signing method: %w", err)
	}

	// To avoid key declared but not used error incase not used in if statements
	_ = key

	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}
	return
}

// DBPaginationRemover manages removing database paginated data.
func DBPaginationRemover(redisPool *redis.Pool) {
	var removePaginationTables = func() {
		// recover from panics
		defer func() {
			if r := recover(); r != nil {
				log.Println("recover", r)
				debug.PrintStack()
			}
		}()

		var IDs []string
		err := pg.Select(&IDs, fmt.Sprintf(`
			DELETE FROM %s.pagination_logs WHERE now() - ts > ($1 || ' seconds')::INTERVAL RETURNING id
		`, PaginationSchema), PaginateDataMaxWaitSeconds)
		if err != nil {
			log.WithError(err).Errorf("Failed to remove outdated pagination ids from pagination_logs.")
			return
		}

		conn := redisPool.Get()
		defer conn.Close()

		for _, id := range IDs {
			// delete record from redis
			_, err = conn.Do("DEL", id)
			if err != nil {
				log.WithError(err).Errorf("Failed to remove pagination from redis.")
			}

			// delete table from database
			tableName := fmt.Sprintf(PaginationTableNameFmt, id)
			_, err = pg.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s."%s"`, PaginationSchema, tableName))
			if err != nil {
				log.WithError(err).Errorf("Failed to remove outdated pagination table.")
				continue
			}
		}
	}

	// clean every half PaginateTimeLimitSeconds time passed
	ticker := time.NewTicker(PaginateTimeLimitSeconds / 2 * time.Second)
	for {
		select {
		case <-ticker.C:
			go removePaginationTables()
		}
	}
}

// configPaginationStorage sets up access to Redis/database for paginated data.
func configPaginationStorage() *redis.Pool {
	var err error

	// create redis pool
	redisPool := redis.Pool{
		MaxIdle:     RediStoreMaxIdle,
		MaxActive:   RedisMaxActive,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			options := []redis.DialOption{
				redis.DialDatabase(1),
			}

			if len(RediStorePassword) > 0 {
				options = append(options, redis.DialPassword(RediStorePassword))
			}

			return redis.Dial(RediStoreNetwork, RediStoreAddress, options...)
		},
		Wait: true,
	}

	// create pagination schema
	_, err = pg.Exec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s AUTHORIZATION %s`, PaginationSchema, PGDBUser))
	if err != nil {
		log.Fatal(err)
	}

	// create pagination log table to store database pagination ids
	_, err = pg.Exec(fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s.pagination_logs (
            id VARCHAR,
            ts TIMESTAMPTZ NOT NULL DEFAULT now(),
            PRIMARY KEY (id)
        )
    `, PaginationSchema))
	if err != nil {
		log.Fatal(err)
	}

	// start db pagination remover
	go DBPaginationRemover(&redisPool)

	return &redisPool
}

// store data for page of paginated data.
type PaginatedPageData struct {
	Data json.RawMessage `json:"data"`
	// number of elements in Data
	NumElements int `json:"num_elements"`
	// flag to indicate if last page of Data
	IsLast bool `json:"is_last"`
}

// output for paginated data.
type PaginatedData struct {
	PaginatedPageData
	// total number of all elements
	NumTotal int `json:"num_total"`
	// relative to PaginatedData.NumElements, 0 indexed
	ElementIndex int `json:"element_index"`
}

// EncodePageDataFn function signature to build PaginatedPageData
// i is index for start of slice of elements (inclusive)
// j is index for end of slice of elements (exclusive) must handle out of range
type EncodePageDataFn func(i, j int) (PaginatedPageData, error)

// PaginateData stores key with fields for numPerPage, numTotal and pages in redis.
func PaginateData(key string, numPerPage, numTotal int, getEncodedData EncodePageDataFn) (err error) {
	const (
		setNotExistCmd = "HSETNX"
		setCmd         = "HSET"
		expireCmd      = "EXPIRE"
	)

	// get connection to redis
	conn := paginationDataStorage.Get()
	defer conn.Close()

	// set pagination metadata total number of elements
	n, err := redis.Int(conn.Do(setNotExistCmd, key, PaginateDataFieldNumTotal, numTotal))
	if err != nil {
		return
	} else if n == 0 {
		err = fmt.Errorf("key (%s) already exists", key)
		return
	}
	// remove key upon returning with error
	defer func(e *error) {
		if *e != nil {
			conn := paginationDataStorage.Get()
			defer conn.Close()

			if _, err := conn.Do("DEL", key); err != nil {
				log.WithError(err).Errorf("Failed to remove key from redis.")
				return
			}
		}
	}(&err)

	// set pagination metadata number of elements per page and is_db flag
	n, err = redis.Int(conn.Do(setCmd, key, PaginateDataFieldNumPerPage, numPerPage, PaginateDataIsDB, false))
	if err != nil {
		return
	} else if n < 2 {
		err = fmt.Errorf("error setting key (%s) num_per_page/is_db field", key)
		return
	}

	// set the first paginated page
	page, err := getEncodedData(0, numPerPage)
	if err != nil {
		return
	}
	value, err := json.Marshal(page)
	if err != nil {
		return
	}

	pageNum := 0
	n, err = redis.Int(conn.Do(setCmd, key, pageNum, value))
	if err != nil {
		return
	} else if n == 0 {
		err = fmt.Errorf("error setting key (%s) page (%d)", key, pageNum)
		return
	}

	// set key expiration
	n, err = redis.Int(conn.Do(expireCmd, key, PaginateTimeLimitSeconds))
	if err != nil {
		return
	} else if n == 0 {
		err = fmt.Errorf("error setting key (%s) expiry", key)
		return
	}

	// set remaining pages in goroutine
	go func() {
		setErrorFn := func(conn redis.Conn, err error) {
			log.Errorf(err.Error())

			n, err = redis.Int(conn.Do(setCmd, key, PaginateDataError, err.Error()))
			if err != nil {
				log.WithError(err).Errorf("error setting key (%s) error field", key)
			} else if n == 0 {
				log.Errorf("failed to set key (%s) error field", key)
			}
		}

		// recover from panics
		defer func() {
			if r := recover(); r != nil {
				log.Println("recover", r)
				debug.PrintStack()

				// get connection to redis
				conn := paginationDataStorage.Get()
				defer conn.Close()

				setErrorFn(conn, fmt.Errorf("encountered a panic saving pagination data pages"))
			}
		}()

		// get connection to redis
		conn := paginationDataStorage.Get()
		defer conn.Close()

		for i := numPerPage; i < numTotal; i += numPerPage {
			page, err := getEncodedData(i, i+numPerPage)
			if err != nil {
				setErrorFn(conn, err)
				return
			}

			value, err := json.Marshal(page)
			if err != nil {
				setErrorFn(conn, err)
				return
			}

			pageNum++
			n, err = redis.Int(conn.Do(setCmd, key, pageNum, value))
			if err != nil {
				setErrorFn(conn, err)
				return
			} else if n == 0 {
				setErrorFn(conn, fmt.Errorf("error setting key (%s) page (%d)", key, pageNum))
				return
			}
		}
	}()

	return nil
}

// PaginateDBData stores key with fields for numPerPage, numTotal in redis and data in database.
func PaginateDBData(numPerPage int, statement string, args ...interface{}) (key string, err error) {
	const (
		setNotExistCmd = "HSETNX"
		setCmd         = "HSET"
		expireCmd      = "EXPIRE"
	)

	// generate key
	key, err = GeneratePaginationKey()
	if err != nil {
		return
	}

	// get connection to redis
	conn := paginationDataStorage.Get()
	defer conn.Close()

	// set pagination metadata number of elements per page
	n, err := redis.Int(conn.Do(setNotExistCmd, key, PaginateDataFieldNumPerPage, numPerPage))
	if err != nil {
		return
	} else if n == 0 {
		err = fmt.Errorf("key (%s) already exists", key)
		return
	}
	// remove key upon returning with error
	defer func(e *error) {
		if *e != nil {
			conn := paginationDataStorage.Get()
			defer conn.Close()

			if _, err := conn.Do("DEL", key); err != nil {
				log.WithError(err).Errorf("Failed to remove key from redis.")
				return
			}
		}
	}(&err)

	// set key expiration
	n, err = redis.Int(conn.Do(expireCmd, key, PaginateTimeLimitSeconds))
	if err != nil {
		return
	} else if n == 0 {
		err = fmt.Errorf("error setting key (%s) expiry", key)
		return
	}

	tx, err := pg.Beginx()
	if err != nil {
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(fmt.Sprintf(`INSERT INTO %s.pagination_logs(id) VALUES ($1)`, PaginationSchema), key)
	if err != nil {
		return
	}

	// create new pagination table to store the query output and add row number column for indexing
	tableName := fmt.Sprintf(PaginationTableNameFmt, key)
	_, err = tx.Exec(fmt.Sprintf(`
		CREATE UNLOGGED TABLE %s."%s"
		WITH (fillfactor=100)
		AS
			SELECT *, row_number() OVER () AS %s
			FROM (%s) AS data_table
		WITH DATA
	`, PaginationSchema, tableName, PaginationTableRowIdCol, statement), args...)
	if err != nil {
		return
	}

	// make PaginationTableRowIdCol primary key for pagination table
	_, err = tx.Exec(fmt.Sprintf(`
		ALTER TABLE %s."%s" ADD PRIMARY KEY (%s)
	`, PaginationSchema, tableName, PaginationTableRowIdCol))
	if err != nil {
		return
	}

	var numTotal int
	err = tx.Get(&numTotal, fmt.Sprintf(`SELECT count(*) FROM %s."%s"`, PaginationSchema, tableName))
	if err != nil {
		return
	}

	// set pagination metadata total number of elements and is_db flag
	n, err = redis.Int(conn.Do(setCmd, key, PaginateDataFieldNumTotal, numTotal, PaginateDataIsDB, true))
	if err != nil {
		return
	} else if n < 2 {
		err = fmt.Errorf("error setting key (%s) num_per_page/is_db field", key)
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	return
}

// GetPaginatedMetadata retrieves metadata for pagination key.
func GetPaginatedMetadata(connPtr *redis.Conn, key string) (numTotal, numPerPage int, isDB bool, err error) {
	const (
		getMultiCmd = "HMGET"
		getCmd      = "HGET"
	)

	var conn redis.Conn
	if connPtr == nil {
		conn = paginationDataStorage.Get()
		defer conn.Close()
	} else {
		conn = *connPtr
	}

	// get pagination metadata number total and number per page
	is, err := redis.Ints(conn.Do(getMultiCmd, key, PaginateDataFieldNumTotal, PaginateDataFieldNumPerPage))
	if err != nil {
		return
	} else if len(is) < 2 {
		err = fmt.Errorf("invalid number (%d) of fields returned expected at least 2", len(is))
		return
	} else {
		numTotal = is[0]
		numPerPage = is[1]
	}

	// get is db flag
	isDB, err = redis.Bool(conn.Do(getCmd, key, PaginateDataIsDB))
	if err != nil {
		return
	}

	return
}

// GetPaginatedData retrieves a page of data for key containing elementIndex starting from 0.
func GetPaginatedData(key string, elementIndex int) (PaginatedData, error) {
	const getCmd = "HGET"

	// get connection to redis
	conn := paginationDataStorage.Get()
	defer conn.Close()

	// get pagination metadata
	var paginatedData PaginatedData
	numTotal, numPerPage, isDB, err := GetPaginatedMetadata(&conn, key)
	if err != nil {
		return PaginatedData{}, err
	} else if isDB {
		return PaginatedData{}, fmt.Errorf("unsupported data retrieval for data from database")
	} else {
		paginatedData.NumTotal = numTotal
	}

	// get paginated data
	if numPerPage > 0 && paginatedData.NumTotal > 0 {
		// check elementIndex is valid
		if elementIndex >= paginatedData.NumTotal {
			return PaginatedData{}, fmt.Errorf("element index (%d) out of range", elementIndex)
		}

		page := elementIndex / numPerPage

		// attempt to get data since PaginateData goroutine may not be done saving
		var b []byte
		for start := time.Now(); time.Now().Sub(start) < PaginateDataMaxWaitSeconds*time.Second; time.Sleep(time.Second) {
			// check for errors
			var errStr string
			errStr, err = redis.String(conn.Do(getCmd, key, PaginateDataError))
			if err != nil && !errors.Is(err, redis.ErrNil) {
				break
			} else if len(errStr) > 0 {
				err = fmt.Errorf(errStr)
				break
			}

			b, err = redis.Bytes(conn.Do(getCmd, key, page))
			if err == nil {
				break
			}
		}
		if err != nil {
			return PaginatedData{}, err
		}

		err = json.Unmarshal(b, &paginatedData.PaginatedPageData)
		if err != nil {
			return PaginatedData{}, err
		}

		paginatedData.ElementIndex = elementIndex % numPerPage
	} else {
		b, err := redis.Bytes(conn.Do(getCmd, key, 0))
		if err != nil {
			return PaginatedData{}, err
		}

		err = json.Unmarshal(b, &paginatedData.PaginatedPageData)
		if err != nil {
			return PaginatedData{}, err
		}
	}

	return paginatedData, nil
}

// GeneratePaginationKey function to generate a key for pagination data.
func GeneratePaginationKey() (key string, err error) {
	var uuidKey uuid.UUID
	uuidKey, err = uuid.NewV4()
	if err != nil {
		return
	}

	return uuidKey.String(), nil
}

func DecodeJwtToken(tokenString, key, signingMethod string, makePayload func(data []byte) error) (err error) {
	// parse token
	var decodedToken *jwt.Token
	decodedToken, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != signingMethod {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		if token.Method.Alg() == "RS256" {
			key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(key))
			if err != nil {
				return "", fmt.Errorf("validate: parse key: %w", err)
			}
			return key, nil
		} else if token.Method.Alg() == "HS256" {
			return []byte(key), nil
		} else {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// validate token
		return []byte(key), nil
	})

	if err != nil {
		err = fmt.Errorf("validate: parse token: %w", err)
		return
	}

	// validate claims
	claims, ok := decodedToken.Claims.(jwt.MapClaims)
	if !ok || !decodedToken.Valid {
		err = fmt.Errorf("Invalid token")
		return
	}

	// unmarshal claims into payload
	var ba []byte
	ba, err = json.Marshal(claims)
	if err != nil {
		err = fmt.Errorf("validate: marshal claims: %w", err)
		return
	}

	err = makePayload(ba)
	if err != nil {
		err = fmt.Errorf("validate: unmarshal claims: %w", err)
		return
	}

	return
}

// GeneratePassword to generate a password
func GeneratePassword(passwordLength int, symbol, lower, upper, number bool) (string, error) {
	lowerCharSet := "abcdedfghijklmnopqrst"
	upperCharSet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	symbolCharSet := "!@#$%&*"
	numberSet := "0123456789"
	var allCharSet string
	var password strings.Builder

	if lower {
		allCharSet += lowerCharSet
	}
	if upper {
		allCharSet += upperCharSet
	}
	if symbol {
		allCharSet += symbolCharSet
	}
	if number {
		allCharSet += numberSet
	}

	for i := 0; i < passwordLength; i++ {
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(allCharSet))))
		if err != nil {
			return "", err
		}

		random := int(nBig.Int64())
		password.WriteString(string(allCharSet[random]))
	}

	return password.String(), nil
}

// Convert a validator error to a string
func MsgToTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "gte":
		return "Value must be greater than or equal to " + fe.Param()
	case "lte":
		return "Value must be less than or equal to " + fe.Param()
	case "min":
		return "Value must be greater than " + fe.Param()
	case "max":
		return "Value must be less than " + fe.Param()
	case "eq":
		return "Value must be equal to " + fe.Param()
	case "ne":
		return "Value must not be equal to " + fe.Param()
	case "oneof":
		return "Value must be one of " + fe.Param()
	case "contains":
		return "Value must contain " + fe.Param()
	case "url":
		return "Invalid URL"
	case "uuid":
		return "Invalid UUID"
	}
	return fe.Error() // default error
}

func FormatErrorResponse(err error) (response utils.Response, code int) {
	// Set default errors
	code = http.StatusInternalServerError
	valError := &utils.ValError{
		Code:    utils.ErrCodeInternalError,
		Message: err.Error(),
	}

	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case UniqueViolation:
			valError = &utils.ValError{
				Code:    utils.ErrCodeUniqueViolation,
				Param:   pqErr.Constraint,
				Message: err.Error(),
			}
			code = http.StatusConflict
		case ForeignKeyViolation:
			valError = &utils.ValError{
				Code:    utils.ErrCodeForeignKeyViolation,
				Param:   pqErr.Constraint,
				Message: err.Error(),
			}
			code = http.StatusConflict
		case CheckViolation:
			valError = &utils.ValError{
				Code:    utils.ErrCodeCheckViolation,
				Param:   pqErr.Constraint,
				Message: err.Error(),
			}
			code = http.StatusConflict
		case ExclusionViolation:
			valError = &utils.ValError{
				Code:    utils.ErrCodeExclusionViolation,
				Param:   pqErr.Constraint,
				Message: err.Error(),
			}
			code = http.StatusConflict
		case NotNullViolation:
			valError = &utils.ValError{
				Code:    utils.ErrCodeNotNullViolation,
				Param:   pqErr.Constraint,
				Message: err.Error(),
			}
			code = http.StatusConflict
		case RestrictViolation:
			valError = &utils.ValError{
				Code:    utils.ErrCodeRestrictViolation,
				Param:   pqErr.Constraint,
				Message: err.Error(),
			}
			code = http.StatusConflict
		case IntegrityConstraintViolation:
			valError = &utils.ValError{
				Code:    utils.ErrCodeIntegrityConstraintViolation,
				Param:   pqErr.Constraint,
				Message: err.Error(),
			}
			code = http.StatusConflict
		case UndefinedColumn:
			valError = &utils.ValError{
				Code:    utils.ErrCodeUndefinedColumn,
				Param:   pqErr.Constraint,
				Message: err.Error(),
			}
			code = http.StatusInternalServerError
		default:
			valError = &utils.ValError{
				Code:    utils.ErrCodeInternalError,
				Message: err.Error(),
			}
			code = http.StatusInternalServerError
		}
	}

	response = utils.Response{Error: valError}
	return
}

// Handle websocket connections.
type WebsocketConnection struct {
	ws          *websocket.Conn
	wsWriteLock *sync.Mutex
	closed      bool
}

type WSPacket struct {
	Type      string      `json:"type"`
	Message   interface{} `json:"message"`
	Error     *string     `json:"error"`
	TSSending time.Time   `json:"ts_sent"`
}

func makeWSPacket(data interface{}) []byte {
	b, err := json.Marshal(data)
	if err != nil {
		log.Warn("Marshalling websocket packet message.")
		b = marshalError()
	}

	return b
}

func marshalError() []byte {
	now := time.Now().Format(time.RFC3339Nano)
	return []byte(fmt.Sprintf(`{"type": "error", "message": null, "error": "marshal failed", "ts_sent": "%s"}`, now))
}

func sendWebsocketMessage(key interface{}, data interface{}, msgType string) {
	webSocketsLock.RLock()
	ws, ok := webSocketsMap[key]
	webSocketsLock.RUnlock()
	if ok && ws != nil {
		msg := makeWSPacket(WSPacket{Type: msgType, Message: data, TSSending: time.Now()})
		ws.WriteMsg(websocket.TextMessage, msg)
	} else {
		log.Warn("Could not find websocket")
	}
}

func (c *WebsocketConnection) WriteMsg(messageType int, message []byte) {
	var err error

	c.wsWriteLock.Lock()
	if !c.closed {
		err = c.ws.WriteMessage(messageType, message)
	}
	c.wsWriteLock.Unlock()

	if err != nil {
		log.Error("Error writing to websocket.")
	}
}

var webSocketsMap map[interface{}]*WebsocketConnection
var webSocketsLock *sync.RWMutex

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Upgrade(w http.ResponseWriter, r *http.Request, key interface{}) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	defer func() {
		conn.Close()
		webSocketsLock.Lock()
		// remove connection only if currently stored connection matches old connection
		if storedWS, ok := webSocketsMap[key]; ok && storedWS.ws == conn {
			delete(webSocketsMap, key)
		}
		webSocketsLock.Unlock()
	}()

	webSocketsLock.Lock()
	oldConn, ok := webSocketsMap[key]
	ws := &WebsocketConnection{conn, &sync.Mutex{}, false}
	webSocketsMap[key] = ws
	if ok {
		// close conflicting websocket connection
		go func(oc *WebsocketConnection) {
			oc.closed = true
			time.Sleep(3 * time.Second)
			oc.ws.Close()
		}(oldConn)
	}
	webSocketsLock.Unlock()

	for {
		if _, _, err := conn.NextReader(); err != nil {
			break
		}
	}

	return nil
}
