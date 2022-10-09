## Update from 2022

This project was made in the wild west of pre-stdlib context days to explore what was possible. Since then, we've all agreed that putting databases in your context is generally a bad idea. Transactions are arguably OK. Regardless, you probably shouldn't use this.

## db [![GoDoc](https://godoc.org/github.com/guregu/db?status.svg)](https://godoc.org/github.com/guregu/db) 
`import "github.com/guregu/db"`

db is a simple helper for using [x/net/context](https://blog.golang.org/context) with various databases. With db you can give each of your connections a name and shove it in your context. Later you can use that name to retrieve the connection. Feel free to fork this and add your favorite drivers. 

### Example
First we make a context with our DB connection. Then we use [kami](https://github.com/guregu/kami) to set up a web server and pass that context to every request. From within the request, we retrieve the DB connection and send a query.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/guregu/db"
	"github.com/guregu/kami"
	"golang.org/x/net/context"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ctx := context.Background()
	ctx = db.OpenSQL(ctx, "main", "mysql", "root:hunter2@unix(/tmp/mysql.sock)/myCoolDB")
	defer db.Close(ctx) // closes all DB connections
	kami.Context = ctx

	kami.Get("/hello/:name", hello)
	kami.Serve()
}

func hello(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	mainDB := db.SQL(ctx, "main") // *sql.DB
	var greeting string
	mainDB.QueryRow("SELECT content FROM greetings WHERE name = ?", kami.Param(ctx, "name")).Scan(&greeting)
	fmt.Fprint(w, greeting)
}
```

### License
BSD
