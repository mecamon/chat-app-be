package config

import "flag"

var app App

type App struct {
	Port                                           string
	IsProd                                         bool
	DBUser, DBUserPassword, DBHost, DBPort, DBName string
	RecoverHostAndPath                             string
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
	flag.Parse()
}

func GetConfig() *App {
	return &app
}
