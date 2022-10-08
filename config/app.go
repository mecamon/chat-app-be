package config

import "flag"

var app App

type App struct {
	Port                                           string
	IsProd                                         bool
	DBUser, DBUserPassword, DBHost, DBPort, DBName string
	RecoverHostAndPath                             string
	EmailAcc, EmailAccPass, EmailHost              string
	EmailPort                                      int
	StorageCloudName                               string
	StorageAPIKey                                  string
	StorageAPISecret                               string
	StorageDirectory                               string
}

func SetConfig() {
	flag.StringVar(&app.Port, "port", ":8080", "port in which the app is running")
	flag.BoolVar(&app.IsProd, "prod", true, "production ready")

	flag.StringVar(&app.DBUser, "dbuser", "", "database user")
	flag.StringVar(&app.DBUserPassword, "dbuserpassword", "", "database user's password")
	flag.StringVar(&app.DBHost, "dbhost", "", "host on which the db running")
	flag.StringVar(&app.DBPort, "dbport", "", "port on which the database is running")
	flag.StringVar(&app.DBName, "dbname", "", "database name")

	flag.StringVar(&app.RecoverHostAndPath, "recoverpasslink", "", "recovery pass link")

	flag.StringVar(&app.EmailAcc, "emailacc", "", "email account to send mails from")
	flag.StringVar(&app.EmailAccPass, "emailpass", "", "email account's password")
	flag.StringVar(&app.EmailHost, "emailhost", "smtp.live.com", "email host")
	flag.IntVar(&app.EmailPort, "emailport", 587, "email port")

	flag.StringVar(&app.StorageCloudName, "storagecloud", "", "storage cloud name")
	flag.StringVar(&app.StorageAPIKey, "storagekey", "", "storage API key")
	flag.StringVar(&app.StorageAPISecret, "storagesecret", "", "storage API secret")
	flag.StringVar(&app.StorageDirectory, "storagedir", "", "storage directory name")
	flag.Parse()
}

func GetConfig() *App {
	return &app
}
