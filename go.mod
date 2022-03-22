module github.com/bet365/go-mssqldb-auth-krb5

go 1.15

require (
	github.com/denisenkom/go-mssqldb v0.0.0-00010101000000-000000000000
	github.com/jcmturner/gokrb5/v8 v8.4.2
)

replace github.com/denisenkom/go-mssqldb => github.com/bet365/go-mssqldb v0.12.1-0.20220322133645-279444b5898e
