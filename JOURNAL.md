
# TODO

# JOURNAL

## 2024-06-02 20:29

### Query Parameters, die Vanilla-Methode

Das folgende Beispiel zeigt, wie man an _Query Paramaters_ herankommt. _Query_
Parameters, wohlgemerkt! Das sind Parameter, die nicht „ganz hinten“ mit `{id}`
angehängt werden und von dem neuen ServeMux sehr gut gemanaged werden; Query
Parameters werden bei einem `GET` mit Hilfe von `?` und `&` hinten an den
Endpunkt drangehängt.

Beispiel: `https://www.youtube.com/watch?v=RNLLPbMThGM&t=39`

Das folgende Beispiel zeigt, wie man damit umgeht:

```go
func handleUrlQuery(w http.ResponseWriter, r *http.Request) {
	// get the id
	rawId := r.URL.Query().Get(`id`)

	// validate the id:
	//    - it must be numerical
	//    - it must be greater than 0
	id, err := strconv.Atoi(rawId)
	if err != nil || id <= 0 {
		http.Error(w, `invalid ID!`, http.StatusBadRequest)
		// http.Error() is only for writing, not for exiting, so ...
		return
	}

	w.Write([]byte(fmt.Sprintf(`You were looking for something with id '%s'`, rawId)))
}
```

So wie wir hier die ID validiert haben, müssen wir dann jeden einzelnen _Query
Parameter_ validieren ...



## 2024-06-02 11:02

Weil es pädagogisch so wertvoll ist, betrachten wir ein bisschen
`net/http`-Verhalten der alten Schule. Da kam nicht automatisch ein
`415`-Fehler, wenn eine falsche HTTP-Methode angewählt wurde; man musste das
alles selber backen. So sah das z.B. aus:

```go
func handleNewSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		// Use the Header.Set() method to add an 'Allow: POST' header
		// to the response header map.
		// The first parameter is the header name, and the second parameter
		// is the header value.
		// Again: Set() will create or update the header we want to be there; nothing appended!
		w.Header.Set(`Allow`, `POST`)
		// w.writeHeader() provides ‘seal and signature’ to the current header material
		// After this, the response header cannot be reversed! It accepts the HTTP
		// status code as an argument
		w.writeHeader(http.StatusMethodNotAllowed)
		w.write([]byte(`Method not allowed!`)
	}
}
```

Alles wesentliche steht in den Kommentaren.

Eines noch: `http.StatusMethodNotAllowed` ist
[hier](https://pkg.go.dev/net/http@go1.22.3#pkg-constants) dokumentiert!

Bleibt noch die Ausgabe mit _httpie:_

```bash
$ http GET ':3000/snippets/new'

HTTP/1.1 405 Method Not Allowed
Allow: POST
Content-Length: 19
Content-Type: text/plain; charset=utf-8
Date: Sun, 02 Jun 2024 14:38:34 GMT

Method not allowed!
```

### Alternativ: http.Error()

Mit folgendem Trick hätten wir die Sache von oben noch vereinfachen können:

```go
func handleNewSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header.Set(`Allow`, `POST`)
		// http.Error() is just a wrapper for the exact thing we did above!
		http.Error(w, `Method not allowed!`, http.StatusMethodNotAllowed)
	}
}

```


## 2024-06-02 09:21

### Response-Header hinzufügen

Mit Go-Bordmitteln gibt es dafür zwei Methoden:

```go
func (p Proxy) handleXyz(w http.ResponseWriter, r *http.Request) {
    // add a value to an existing header
	r.Header.Add("Header01", "head 01")
	// create a new / reset an existing header
	r.Header.Set("Header02", "head 02")
	p.Proxy.ServeHTTP(w, r)
}
```

#### Anmerkungen

1. Ich habe für dieses Beispiel irgendwo einen Proxy-Server im System laufen,
   dem ich auf diese Weise Anweisungen gebe, in diesem Fall die versteckte
   Anweisung, die Header beim Client anzuzeigen (sieht man hier nicht).
2. Mit `r.Header.Add()` bleibt der entsprechende Header in der Request
   erhalten; nichts wird ersetzt. Stattdessen wird das Neue „unten angefügt“.
3. Mit `r.Header.Set()` wird entweder ein neuer Header geschaffen, oder ein
   bestehender Header wird überschrieben.

`Add()` entspricht also _append,_ oder dem `>>`-Opeartor in der Shell; `Set()`
bedeutet _create or reset,_ das, was der `>`-Operator in der Shell macht.

```bash
$ curl -i http://localhost:8080 -H 'header01: foo' -H 'header0: bar'
HTTP/1.1 200 OK
Content-Length: 34
Content-Type: text/plain; charset=utf-8
Date: Tue, 17 Jan 2023 20:46:34 GMT

header01: [head1 foo], header02: [bar]
```



## 2024-06-01 18:33

Hier geht es um URL-Parameter und wie man im Handler an sie herankommt

```go
// Add a handler function for viewing a specific snippet
func handleSingleSnippetView(w http.ResponseWriter, r *http.Request) {
	// this is how we get URL variables (see below); 
	// they will always be strings or more exactly, []byte chains
	id := r.PathValue(`id`)
	w.Write([]byte(fmt.Sprintf("Display snippet with ID '%s'", id)))
}

func main() {
	// use the http.NewServeMux() constructor to initialize a new servemux (router),
	// then register the home() function as handler for the `/` endpoint.
	mux := http.NewServeMux()

	// This is an endpoint with URL parameters
	mux.HandleFunc(`GET /snippets/{id}`, handleSingleSnippetView)
}
```

## 2024-06-01 17:53

Wir haben unsere Routes um ein einfaches `POST` erweitert:

```go
// Add a handler function for creating a snippet.
func handleNewSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`Creating a new snippet ...`))
}

func main() {
	// use the http.NewServeMux() constructor to initialize a new servemux (router),
	// then register the home() function as handler for the `/` endpoint.
	mux := http.NewServeMux()

	mux.HandleFunc(`GET /`, handleHome)
	// This is how we add a POST endpoint in go 1.22
	mux.HandleFunc(`POST /new`, handleNewSnippet)
    // ...
}
```

Wichtig für `mux.HandleFunc()` ist, dass nur __ein einziges__ Leerzeichen
zwischen der HTTP-Methode und dem Endpoint steht.

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

<!--
vim: ts=4 sw=4 fdm=indent
-->
