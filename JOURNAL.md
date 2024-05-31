
# TODO

# JOURNAL

## 2024-05-30 17:39

```go
// Define a simple home Handler function which writes a byte slice
// containing "Hello from Snippetbox" as a response body
func handleHome(w http.ResponseWriter, r *http.Request) {
	// Write() accepts only []byte as ‘most neutral’ message type
	w.Write([]byte(`Hello from Snippetbox!`))
}

func main() {
	// use the http.NewServeMux() constructor to initialize a new servemux (router),
	// then register the home() function as handler for the `/` endpoint.
	mux := http.NewServeMux()
	// This is how it's done in go 1.22+
	mux.HandleFunc(`GET /`, handleHome)

	// Use the http.ListenAndServe() function as web serving unit. It accepts two parameters:
	//   - the URL (which will be `localhost:3000` here)
	//   - the router we just created.
	// If the webserver returns an error, we handle it using log.Fatal() to log the error and exit.
	// Note that any error returned by http.ListenAndServe() is non-nil!
	log.Println("starting server at port :3000")
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatalf("Uh oh! %s", err)
	}
}
```

Hier wurde ein einfacher Webserver mit den neuen Möglichkeiten von _go 1.22_
implementiert. Es läuft.

Entscheidend ist `mux.HandleFunc()`. Sie akzeptiert eine Route und einen
Endpoint handler. Der Endpoint Handler muss einen `http.ResponseWriter()` und
einen `*http.Request`-Referenz als Parameter akzeptieren. Nur dann wird er als
Endpoint Handler erkannt.

## 2024-05-30

Hab das Projekt noch einmal neu gestartet; dieses Mal soll es mit ausführlichen
Anmerkungen gemacht werden. Alles von Grund auf. „Dieses Mal richtig.“
