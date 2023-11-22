# gorenjskabike_api
Show data from gorenjska api, regarding ebike borrowing status.

Local webpage calls throug node.js poxy API for latest data.

## Data source
https://trzic.musiclab.si/openapi/swagger#/

More specific, from this url:
https://trzic.musiclab.si/api/gorenjskabike?page=1&size=1000

## Requirements
Node.js must be installed on your system. This is because, I am trying to evade CORS (Cross-Origin Resource Sharing) error when I call API directly.

You should have these modules installed on your system:
**npm install express node-fetch**

## Requirements for Go Lang version
Use command in Termina/Command Prompt, to get nececary gorilla mux library
**go get github.com/gorilla/mux**

And to build it, use command: **go build main.go**
I have provided go.mod, use directory, If you need help with building code, I have read this page:
https://go.dev/doc/tutorial/create-module

Running server on Linux anbd Mac: **./main** or just run **main.exe** on Windows




