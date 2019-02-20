package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/kamaln7/klein/alias"
	"github.com/kamaln7/klein/alias/alphanumeric"
	"github.com/kamaln7/klein/alias/memorable"
	"github.com/kamaln7/klein/auth"
	"github.com/kamaln7/klein/auth/httpbasic"
	"github.com/kamaln7/klein/auth/statickey"
	"github.com/kamaln7/klein/auth/unauthenticated"
	"github.com/kamaln7/klein/server"
	"github.com/kamaln7/klein/storage"
	"github.com/kamaln7/klein/storage/bolt"
	"github.com/kamaln7/klein/storage/file"
	"github.com/kamaln7/klein/storage/postgresql"
	"github.com/kamaln7/klein/storage/redis"
	"github.com/kamaln7/klein/storage/spaces"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "klein",
	Short: "klein is a minimalist URL shortener.",
	Long:  "klein is a minimalist URL shortener.",
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.New(os.Stdout, "[klein] ", log.Ldate|log.Ltime)

		// 404
		notFoundHTML := []byte("404 not found")
		notFoundPath := viper.GetString("error-template")
		if notFoundPath != "" {
			var err error
			notFoundHTML, err = ioutil.ReadFile(notFoundPath)
			if err != nil {
				logger.Fatal(err)
				return
			}
		}

		// auth
		var authProvider auth.Provider
		switch viper.GetString("auth.driver") {
		case "none":
			authProvider = unauthenticated.New()
		case "basic":
			username := viper.GetString("auth.basic.username")
			password := viper.GetString("auth.basic.password")
			if username == "" || password == "" {
				logger.Fatalf("You need to provide a username and password in order to use basic auth")
			}

			authProvider = httpbasic.New(&httpbasic.Config{
				Username: username,
				Password: password,
			})
		case "key":
			key := viper.GetString("auth.key")
			if key == "" {
				logger.Fatalf("You need to provide an auth key in order to use key auth")
			}
			authProvider = statickey.New(&statickey.Config{
				Key: key,
			})
		default:
			logger.Fatal("invalid auth driver")
		}

		// storage
		var storageProvider storage.Provider
		switch viper.GetString("storage.driver") {
		case "file":
			storageProvider = file.New(&file.Config{
				Path: viper.GetString("storage.file.path"),
			})
		case "boltdb":
			var err error
			storageProvider, err = bolt.New(&bolt.Config{
				Path: viper.GetString("storage.boltdb.path"),
			})

			if err != nil {
				logger.Fatalf("could not open bolt database: %s\n", err.Error())
			}
		case "redis":
			var err error
			storageProvider, err = redis.New(&redis.Config{
				Address: viper.GetString("storage.redis.address"),
				Auth:    viper.GetString("storage.redis.auth"),
				DB:      viper.GetInt("storage.redis.db"),
			})

			if err != nil {
				logger.Fatalf("could not open redis database: %s\n", err.Error())
			}
		case "spaces.stateful":
			accessKey := viper.GetString("storage.spaces.access_key")
			secretKey := viper.GetString("storage.spaces.secret_key")
			region := viper.GetString("storage.spaces.region")
			space := viper.GetString("storage.spaces.space")

			if accessKey == "" || secretKey == "" || region == "" || space == "" {
				logger.Fatalf("You need to provide an access key, secret key, region and space to use the spaces storage backend")
			}

			var err error
			storageProvider, err = spaces.New(&spaces.Config{
				AccessKey: accessKey,
				SecretKey: secretKey,
				Region:    region,
				Space:     space,
				Path:      viper.GetString("storage.spaces.stateful.path"),
			})

			if err != nil {
				logger.Fatalf("could not connect to spaces: %s\n", err.Error())
			}
		case "sql.pg":
			var err error
			storageProvider, err = postgresql.New(&postgresql.Config{
				Host:     viper.GetString("storage.sql.pg.host"),
				Port:     viper.GetInt32("storage.sql.pg.port"),
				User:     viper.GetString("storage.sql.pg.user"),
				Password: viper.GetString("storage.sql.pg.password"),
				Database: viper.GetString("storage.sql.pg.database"),
				Table:    viper.GetString("storage.sql.pg.table"),
				SSLMode:  viper.GetString("storage.sql.pg.sslmode"),
			})

			if err != nil {
				logger.Fatalf("could not connect to postgresql: %s\n", err.Error())
			}
		default:
			logger.Fatal("invalid storage driver")
		}

		// alias
		var aliasProvider alias.Provider
		switch viper.GetString("alias.driver") {
		case "alphanumeric":
			var err error
			aliasProvider, err = alphanumeric.New(&alphanumeric.Config{
				Length: viper.GetInt("alias.alphanumeric.length"),
				Alpha:  viper.GetBool("alias.alphanumeric.alpha"),
				Num:    viper.GetBool("alias.alphanumeric.num"),
			})

			if err != nil {
				logger.Fatalf("could not select alphanumeric alias: %s\n", err.Error())
			}
		case "memorable":
			aliasProvider = memorable.New(&memorable.Config{
				Length: viper.GetInt("alias.memorable.length"),
			})
		default:
			logger.Fatal("invalid alias driver")
		}

		// url
		publicURL := viper.GetString("url")
		if publicURL == "" {
			publicURL = fmt.Sprintf("http://%s/", viper.GetString("listen"))
		}

		// klein
		k := server.New(&server.Config{
			Alias:   aliasProvider,
			Auth:    authProvider,
			Storage: storageProvider,
			Log:     logger,

			ListenAddr:   viper.GetString("listen"),
			RootURL:      viper.GetString("root"),
			PublicURL:    publicURL,
			NotFoundHTML: notFoundHTML,
		})

		k.Serve()
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	// General options
	rootCmd.PersistentFlags().String("error-template", "", "path to error template")
	rootCmd.PersistentFlags().String("url", "", "path to public facing url")
	rootCmd.PersistentFlags().String("listen", "127.0.0.1:5556", "listen address")
	rootCmd.PersistentFlags().String("root", "", "root redirect")

	// Alias options
	rootCmd.PersistentFlags().String("alias.driver", "alphanumeric", "what alias generation to use (alphanumeric, memorable)")

	rootCmd.PersistentFlags().Int("alias.alphanumeric.length", 5, "alphanumeric code length")
	rootCmd.PersistentFlags().Bool("alias.alphanumeric.alpha", true, "use letters in code")
	rootCmd.PersistentFlags().Bool("alias.alphanumeric.num", true, "use numbers in code")

	rootCmd.PersistentFlags().Int("alias.memorable.length", 3, "memorable word count")

	// Auth options
	rootCmd.PersistentFlags().String("auth.driver", "none", "what auth backend to use (basic, key, none)")

	rootCmd.PersistentFlags().String("auth.key", "", "upload API key")

	rootCmd.PersistentFlags().String("auth.basic.username", "", "username for HTTP basic auth")
	rootCmd.PersistentFlags().String("auth.basic.password", "", "password for HTTP basic auth")

	// Storage options
	rootCmd.PersistentFlags().String("storage.driver", "file", "what storage backend to use (file, boltdb, redis, spaces.stateful, sql.pg)")

	rootCmd.PersistentFlags().String("storage.file.path", "urls", "path to use for file store")

	rootCmd.PersistentFlags().String("storage.boltdb.path", "bolt.db", "path to use for bolt db")

	rootCmd.PersistentFlags().String("storage.redis.address", "127.0.0.1:6379", "address:port of redis instance")
	rootCmd.PersistentFlags().String("storage.redis.auth", "", "password to access redis")
	rootCmd.PersistentFlags().Int("storage.redis.db", 0, "db to select within redis")

	rootCmd.PersistentFlags().String("storage.spaces.access-key", "", "access key for spaces")
	rootCmd.PersistentFlags().String("storage.spaces.secret-key", "", "secret key for spaces")
	rootCmd.PersistentFlags().String("storage.spaces.region", "", "region for spaces")
	rootCmd.PersistentFlags().String("storage.spaces.space", "", "space to use")

	rootCmd.PersistentFlags().String("storage.spaces.stateful.path", "klein.json", "path of the file in spaces (spaces.stateful driver)")

	rootCmd.PersistentFlags().String("storage.sql.pg.host", "localhost", "postgresql host")
	rootCmd.PersistentFlags().Int32("storage.sql.pg.port", 5432, "postgresql port")
	rootCmd.PersistentFlags().String("storage.sql.pg.user", "klein", "postgresql user")
	rootCmd.PersistentFlags().String("storage.sql.pg.password", "secret", "postgresql password")
	rootCmd.PersistentFlags().String("storage.sql.pg.database", "klein", "postgresql database")
	rootCmd.PersistentFlags().String("storage.sql.pg.table", "klein", "postgresql table")
	rootCmd.PersistentFlags().String("storage.sql.pg.sslmode", "prefer", "postgresql sslmode")

	viper.BindPFlags(rootCmd.PersistentFlags())
}

func initConfig() {
	viper.SetEnvPrefix("klein")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
