# go-mssqldb-auth-krb5

A pure go kerberos authentication provider package for https://github.com/microsoft/go-mssqldb using the gokrb5/v8 package at https://github.com/jcmturner/gokrb5/

In order to use the package, import it alongside the main driver 

``` golang
	_ "github.com/bet365/go-mssqldb-auth-krb5"
	_ "github.com/microsoft/go-mssqldb"
```

It will register itself and become available for use when the connection string parameter "authenticator=krb5" is used.

e.g.

    authenticator=krb5;server=DatabaseServerName;database=DBName;krb5-params.....

The package supports authentication via 3 methods.

* Keytabs - Specify the username, keytab file, the krb5.conf file, and realm.
  
      authenticator=krb5;server=DatabaseServerName;database=DBName;user id=MyUserName;krb5-realm=domain.com;krb5-configfile=/etc/krb5.conf;krb5-keytabfile=~/MyUserName.keytab
  
* Credential Cache - Specify the krb5.conf file path and credential cache file path.

      authenticator=krb5;server=DatabaseServerName;database=DBName;krb5-configfile=/etc/krb5.conf;krb5-keytabcachefile=~/MyUserNameCachedCreds

* Raw credentials - Specity krb5.confg, Username, Password and Realm. 
  
      authenticator=krb5;server=DatabaseServerName;database=DBName;user id=MyUserName;password=MyPassword;krb5-realm=comani.com;krb5-configfile=/etc/krb5.conf;

The parameter names themselves are as follows :

`krb5-configfile`  
path to krb5 configuration file. e.g. /etc/krb5.conf

`krb5-keytabfile`  
path to keytab file.

`krb5-keytabcachefile`  
path to credential cache file.

`krb5-realm`  
domain name for account.

`krb5-dnslookupkdc`  
Optional parameter in all contexts. Set to lookup KDCs in DNS. Boolean. Default is true.

`krb5-udppreferencelimit`  
Optional parameter in all contexts. 1 means to always use tcp. MIT krb5 has a default value of 1465, and it prevents user setting more than 32700. Integer. Default is 1.

An example usage:

``` golang
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	
	_ "github.com/bet365/go-mssqldb-auth-krb5"
	_ "github.com/microsoft/go-mssqldb"	
)

func main() {
	var (
		connectionString string
	)

	flag.StringVar(&connectionString, "connString", "", "Connection string")
	flag.Parse()

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

	sql := "select 1234"

	var value string
	err = db.QueryRowContext(ctx, sql).Scan(&value)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(value)
}
```