# go-mssqldb-auth-krb5

This package provides an implementation of auth.Auth from https://github.com/denisenkom/go-mssqldb using the gokrb5/v8 package at https://github.com/jcmturner/gokrb5/

An example usage:

``` golang
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	krb5 "github.com/bet365/go-mssqldb-auth-krb5"
	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/jcmturner/gokrb5/v8/client"
	"github.com/jcmturner/gokrb5/v8/config"
)

func main() {

	var (
		username         string
		password         string
		realm            string
		configFile       string
		connectionString string
	)

	flag.StringVar(&username, "username", "", "Username")
	flag.StringVar(&password, "password", "", "Password")
	flag.StringVar(&realm, "realm", "", "Realm")
	flag.StringVar(&configFile, "krb5-config", "", "Path to krb5.conf")
	flag.StringVar(&connectionString, "connString", "", "Connection string")

	flag.Parse()

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// Login using our credentials via gokrb5, returning a client
	krbClient, err := LoginUsername(username, password, realm, string(data))
	if err != nil {
		log.Fatal(err)
	}

	defer krbClient.Destroy()

	// create a new auth.Provider around the kerberos client
	provider := krb5.NewAuthProvider(krbClient)

	// pass the provider to mssql to override the default authentication mechanism
	mssql.SetIntegratedAuthenticationProvider(provider)

	// connect to sql.
	// when the connection is opened it will use the krb5 Auth Provider created above.
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	var value int
	err = db.QueryRowContext(ctx, "select 1234").Scan(&value)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(value)
}

func LoginUsername(userName, password, realm, configString string) (*client.Client, error) {

	cfg, err := config.NewFromString(configString)
	if err != nil {
		return nil, err
	}

	// configure as needed
	cfg.LibDefaults.DNSLookupKDC = true
	cfg.LibDefaults.UDPPreferenceLimit = 1

	krbClient := client.NewWithPassword(userName, realm, password, cfg, client.DisablePAFXFAST(true))

	if err := krbClient.Login(); err != nil {
		return nil, err
	}

	return krbClient, nil
}
```