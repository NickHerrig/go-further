# go-futher
building the let's go futher by alex edwards open movie database api

# Project Layout
- `cmd/api` - app specific code for api(server code,  read/write http req, authentication)
- `internal ` - packages imported by the api(database interaction, data validation, sending emails)
- `migrations` - tbd
- `bin` - tbd
- `remote` - tbd
- `go.mod` - project dependencies, versions, and module path
- `Makefile` - recipes for automation(auditing, builds, migrations)


# Interesting Learnings

## Sending JSON Responses
- json.MarshalIndent() (Slower and more memory, but might be worth it for readability)
- struct tagging for json ('-', and  omitempty)
- Enveloping json responses (using  map[string]interface{})
- Error handling - http has consts for common http [methods](https://pkg.go.dev/net/http#pkg-constants) 
- Custom JSON Encoding implementing the MarshalJSON interface

## Parsing JSON Requests
- Any json k/v pairs which can't be ampped to struct fields will be silently ignored
- json.Unmarshal vs json.Decoder (Unmarshal uses 80% more memory and is a tiny bit slower)
- MaxBytesReader is important for malicous or accidently DDOS requests.
- Custom JSON Decoding implementing the UnmarshalJSON interface
- Validation of incoming JSON requests is pretty dang important.

## Database Setup and Config
pgx: https://github.com/jackc/pgx/blob/master/pgxpool/pool.go#L16
- MaxOpenConns, MaxIdleConns, ConnMaxLifetime, ConnMaxIdleTime.
- Explicitly set MaxOpenConns value below hardlimits set by the database(postgres=100).

## SQL Migrations 
Example Migrations:
`migrate create -seq -ext=.sql -dir=./migrations add_movies_check_constraints`
- `-seq` uses sequential numbering instead of default unix time
- `-ext` uses the .sql extension
- `-dir` places up and down files in the migrations folder

applying migrations:
`migrate -path=./migrations -database=DB_DSN up`

check migration version:
`migrate -path=./migrations -database=DB_DSN version`

migrate up or down specific versions:
`migrate -path=./migrations -database=$DB_DSN goto 1`

migrate down a specific number of migraations:
`migrate -path=./migrations -database=$DB_DSN down 1`

Fixing SQL Errors:
if an error happens, you must manually roll-back th partially applied migration.
Once completed, you must force the version number..
`migrate -path=./migrations -database=$DB_DSN force {version}`

## CRUD Operations
Create, Read, Update, Delete ops were discussed in this chapter. 
This chapter was fairly straightforward. 
There are a few differences with utilizing pgxpool.{DB-Command}()
rather than utilizing the /lib/pg and database/sql packages.

## Advanced CRUD Ops
Partial updates with Patch can be accomplished by utilizing pointers on structs to 
test the difference between nil and the go type zero-value.
JSON items with null values will be ignored and will remain unchanged.
It's important to think through potential race conditions on endpoints that update
values in the database. return conflicts when caught.

## Filtering, Sorting, and Pagination
- [r.URL.Query()](https://pkg.go.dev/net/url#URL.Query) returns a map of query params
- Full Text Search is used to search columns for specific text.
- Indexing for query performance.
- `ORDER BY` help with sorting of records. 
- `LIMIT` and `OFFSET` help with pagination of records.
- `count(*) OVER()` expression results in the filtered record count being included as the first value in each row.

## Rate Limiting

- Common middleware structure

```go
func (app *application) exampleMiddleware(next http.Handler) http.Handler {
    // Any code here will run only once, when we wrap something with the middleware. 
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Any code here will run for every request that the middleware handles.
        next.ServeHTTP(w, r) 
    })
}
```	
- rate limiting with package: golang.org/x/time/rate
- this implementation would not work for an HA setup.

## Signals
- greping for processes `pgrep -l {pattern}`
- send a kill signal `pkill -SIGKILL {process}`
- send a TERM signal `pkill -SIGTERM {process}`
