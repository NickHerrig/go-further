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
