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
	"github.com/kamaln7/klein/storage/redis"
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
		notFoundPath := viper.GetString("template")
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
		switch viper.GetString("auth") {
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
			logger.Fatal("Invalid auth provider")
		}

		// storage
		var storageProvider storage.Provider
		switch viper.GetString("storage") {
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
		default:
			logger.Fatal("Invalid storage backend")
		}

		// alias

		var aliasProvider alias.Provider
		switch viper.GetString("alias") {
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
			logger.Fatal("Invalid alias generator")
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
	rootCmd.PersistentFlags().String("template", "", "path to error template")
	viper.BindPFlag("template", rootCmd.PersistentFlags().Lookup("template"))

	rootCmd.PersistentFlags().String("url", "", "path to public facing url")
	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))

	rootCmd.PersistentFlags().String("listen", "127.0.0.1:5556", "listen address")
	viper.BindPFlag("listen", rootCmd.PersistentFlags().Lookup("listen"))

	rootCmd.PersistentFlags().String("root", "", "root redirect")
	viper.BindPFlag("root", rootCmd.PersistentFlags().Lookup("root"))

	// Alias options
	rootCmd.PersistentFlags().String("alias", "alphanumeric", "what alias generation to use (alphanumeric, memorable)")
	viper.BindPFlag("alias", rootCmd.PersistentFlags().Lookup("alias"))

	rootCmd.PersistentFlags().Int("alias.alphanumeric.length", 5, "alphanumeric code length")
	viper.BindPFlag("alias.alphanumeric.length", rootCmd.PersistentFlags().Lookup("alias.alphanumeric.length"))

	rootCmd.PersistentFlags().Bool("alias.alphanumeric.alpha", true, "use letters in code")
	viper.BindPFlag("alias.alphanumeric.alpha", rootCmd.PersistentFlags().Lookup("alias.alphanumeric.alpha"))

	rootCmd.PersistentFlags().Bool("alias.alphanumeric.num", true, "use numbers in code")
	viper.BindPFlag("alias.alphanumeric.num", rootCmd.PersistentFlags().Lookup("alias.alphanumeric.num"))

	rootCmd.PersistentFlags().Int("alias.memorable.length", 3, "memorable word count")
	viper.BindPFlag("alias.memorable.length", rootCmd.PersistentFlags().Lookup("alias.memorable.length"))

	// Auth options
	rootCmd.PersistentFlags().String("auth", "none", "what auth backend to use (basic, key, none)")
	viper.BindPFlag("auth", rootCmd.PersistentFlags().Lookup("auth"))

	rootCmd.PersistentFlags().String("auth.key", "", "upload API key")
	viper.BindPFlag("auth.key", rootCmd.PersistentFlags().Lookup("auth.key"))

	rootCmd.PersistentFlags().String("auth.basic.username", "", "username for HTTP basic auth")
	viper.BindPFlag("auth.basic.username", rootCmd.PersistentFlags().Lookup("auth.basic.username"))

	rootCmd.PersistentFlags().String("auth.basic.password", "", "password for HTTP basic auth")
	viper.BindPFlag("auth.basic.password", rootCmd.PersistentFlags().Lookup("auth.basic.password"))

	// Storage options
	rootCmd.PersistentFlags().String("storage", "file", "what storage backend to use (file, boltdb, redis)")
	viper.BindPFlag("storage", rootCmd.PersistentFlags().Lookup("storage"))

	rootCmd.PersistentFlags().String("storage.file.path", "urls", "path to use for file store")
	viper.BindPFlag("storage.file.path", rootCmd.PersistentFlags().Lookup("storage.file.path"))

	rootCmd.PersistentFlags().String("storage.boltdb.path", "bolt.db", "path to use for bolt db")
	viper.BindPFlag("storage.boltdb.path", rootCmd.PersistentFlags().Lookup("storage.boltdb.path"))

	rootCmd.PersistentFlags().String("storage.redis.address", "127.0.0.1:6379", "address:port of redis instance")
	viper.BindPFlag("storage.redis.address", rootCmd.PersistentFlags().Lookup("storage.redis.address"))

	rootCmd.PersistentFlags().String("storage.redis.auth", "", "password to access redis")
	viper.BindPFlag("storage.redis.auth", rootCmd.PersistentFlags().Lookup("storage.redis.auth"))

	rootCmd.PersistentFlags().Int("storage.redis.db", 0, "db to select within redis")
	viper.BindPFlag("storage.redis.db", rootCmd.PersistentFlags().Lookup("storage.redis.db"))
}

func initConfig() {
	viper.SetEnvPrefix("klein")
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
