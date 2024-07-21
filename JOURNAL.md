
# TODO

- [ ] (Long-term) Die selbstgebackene Validierung auf _Validator_ umstellen.
  Hier ein
  [Artikel](https://thedevelopercafe.com/articles/payload-validation-in-go-with-validator-626594a58cf6).

# JOURNAL

<!-- ## 2024-07-XX XX:XX -->

## 2024-07-21 13:39

Hier beginnt eine neue Zeitrechnung! Ich habe begonnen, das Projekt von
Standard-Go-Templates auf [Templ](https://templ.guide) umzustellen.

Hier der entscheidende Code:

```go
func (app *Application) handlePing(w http.ResponseWriter, r *http.Request) {
	main := tpl.Ping()
	page := tpl.Base(`Ping`, main, `2024`)
	// return templ.Handler(component)
	page.Render(r.Context(), w)
}
```

Also: Als erstes wird das (fertig gerenderte) _Ping_ Component erstellt.
Dann wird es an das Layout-Component _Base_ weitergereicht. Und zum Schluss
wird das fertige Gesamtpaket gerendert und an den _ResponseWriter_ zum
Verschicken weitergereicht.

## 2024-07-18 19:26

Hab ein wenig mit _Cypress_ herumgespielt und das hier herausgefunden:

In der Dokumentation von _Cypress_ wird empfohlen, so etwas wie hier als
Markierung zu benutzen.

```html
<h2 data-test="ping">Ping Ping</h2>
```

Da haben sie die Rechnung ohne die _Go Template Engine_ gemacht – die löscht
nämlich solche Kinkerlitzchen, und die Tests laufen ins Leere. Da ist es dann
doch besser, sich auf die gute alte `#id` zu verlassen.

```html
<h2 id="ping">Ping Ping</h2>
```

## 2024-07-18 09:35

Habe für das End-To-End-Testing _Cypress_ installiert und bin hellauf begeistert!

## 2024-07-13 19:31

Noch ein Nachtrag zum letzten Eintrag: Ich habe hier ein Muster gefunden, das man zu einem Snippet ausarbeiten könnte:

```git
func TestSecureHeaders(t *testing.T) {
	// >>> snippet tHandlerstubs
	// initialize a Response object and a http request stub
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, `/`, nil)
	if err != nil {
		t.Fatal(err)
	}
	// <<< end snippet

	// >>> snippet tMWnext
	// Create a mock HTTP handler that we can pass to our secureHeaders()
	// middleware which writes a 200 OK status code and an 'OK.' response body
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `OK.`)
	})
	secureHeaders(next).ServeHTTP(rr, r)
	// <<< end snippet


	// >>> snippet tResResult
	// get the results of the 'request'
	rs := rr.Result()
	defer rs.Body.Close()
	// <<< end snippet

	// test #01
	expected := `default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com`
	assert.Equal(t, rs.Header.Get(`Content-Security-Policy`), expected)
	// test #02
	expected = `origin-when-cross-origin`
	assert.Equal(t, rs.Header.Get(`Referrer-Policy`), expected)
	// test #03
	expected = `nosniff`
	assert.Equal(t, rs.Header.Get(`X-content-Type-Options`), expected)
	// test #04
	expected = `deny`
	assert.Equal(t, rs.Header.Get(`X-Frame-Options`), expected)
	// test #05
	expected = `0`
	assert.Equal(t, rs.Header.Get(`X-XSS-Protection`), expected)

	//------------------------------------------------------------
	// test whether secureHeaders has passed the staff properly to 'next'
	//------------------------------------------------------------

	// Check the Status code
	assert.Equal(t, rs.StatusCode, http.StatusOK)


	// >>> snippet tResBody
	// use IO to read the request body and save it in the 'body' variable as
	// []bytes Slice
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	// remove all leading and trailing whitespace
	bytes.TrimSpace(body)
	// <<< end snippet

	// test the body
	assert.Equal(t, string(body), `OK.`)
}
```

Das ist jetzt erst mal nur für die Schublade. Ganz im Sinne von „YAGNI until PITA“ in _Obsidian._

## 2024-07-13 08:38

Die folgenden Anmerkungen beziehen sich auf die `TestPing()`-Funktion in
`cmd/web/handlers_test.go`

Dieses Mal geht es um das Unit Testing von Middleware und Response Handlern.
Zum Glück liefert uns die _stdlib_ Werkzeuge dafür.  

Das wichtigste Werkzeug ist das _Recorder_-Werkzeug aus dem
`net/http/httptest`-Paket. Dieses Werkzeug liefert eine Menge _Props and
Methods,_ mit denen man die _Response_ auslesen und testen kann.

Neu ist hier auch die `http.NewRequest()`-Methode, mit der sich eine _Request_
auch ohne Browser aus dem Hut zaubern lässt. Als _Mock, Stub_ oder _Fake_ –
oder wie auch immer man das nennen will.

## 2024-07-12 16:48

Als nächstes haben wir uns einen eigenen Helfer geschrieben. So sieht er aus:

```go
package assert

import "testing"

func Equal[T comparable](t *testing.T, actual, expected T) {
	// marks this function as Test Helper Function
	t.Helper()

	if actual != expected {
		t.Errorf("'%v' should be '%v'", actual, expected)
	}
}
```

Da gibt es doch einige Sachen anzumerken:

1. `[T comparable]` -- bedeutet: „... für alle Typen (`T`), die sich als Typen
   vergleichen lassen.“ Diese Syntax ist eine sog. _Template_-Syntax. Statt `T`
   kann jeder Datentyp verwendet werden, der sich unmittelbar mit sich selbst
   vergleichen lässt. Einen `string` kann man mit einem `string` vergleichen, eine
   `int` mit einer `int`, ein `bool` mit einem `bool`, und sogar eine `int` mit
   einer `int64`.
2. `*testing.T` bezieht sich auf den Typ `T` im `testing`-Package – __nicht__ auf
   einen Template-Platzhalter!
3. `expected T` bezieht sich wieder auf den Template-Platzhalter von vorhin.
4. `t.Helper()` markiert diese Funktion als _testing helper._ Es wird gleich klar
   werden, was das bedeutet.

Und so haben wir diese Hilfsfunktion eingesetzt:

```go
type humanDateTest struct {
	name	string
	timestamp time.Time
	expected	string
}

func TestHumanDate(t *testing.T) {
	// test cases as 'table'
	// -1-
	testCases := []humanDateTest{
		{
			name:      `UTC`,
			timestamp: time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
			expected:  `2022-03-17 10:15`,
		},
	}

	for _, tc := range testCases {
		// -2-
		t.Run(tc.name, func(t *testing.T) {
			hd := humanDate(tc.timestamp)
			// -3-
			assert.Equal(t, hd, tc.expected)
		})
	}
}
```

#### Anmerkungen

1. Wie man sieht, haben wir hier mit `testCases` eine sog. _Test-Tabelle_
   definiert: Jedes Glied in dieser Tabelle liefert die Daten für alles,
   was später für ‚seinen‘ Test nötig ist.
0. Die Tests werden in der `for`-Schleife durchlaufen. `t.Run()` akzeptiert
   zwei Argumente: Den Namen, der eindeutig benennt, welcher Test z.B. einen
   Fehler auslöste, und eine anonyme Funktion, die genau die gleiche Signatur
   hat wie die einfache Test-Funktion von oben.
0. `assert.Equal()` ist die Helfer-Funktion von eben. Da sie eine 
   Helfer-Funktion ist, brauchen wir hier nicht mehr anzugeben, wie sie mit
   dem Test und mit einem Fehlerfall umzugehen hat. Das haben wir ja
   dort an Ort und Stelle festgelegt.



## 2024-07-12 10:41

### TESTING IN GO

OK, wir machen jetzt Golang Testing. In seiner einfachsten Form sieht ein Test so aus:

Hier erst einmal die zu testende Funktion:

```go
func humanDate(t time.Time) string {
	// time is equal to the Epoch
	if t.IsZero() {
		return ""
	}
	// convert t to UTC before formatting it
	return t.UTC().Format(`2006-01-02 03:04`)
}
```

Und hier der einfachste Test:

```go
// file: helpers_test.go

func TestHumanDate(t *testing.T) {
	// test cases as 'table'
	hd := humanDate(2022, 3, 17, 10, 15, 0, 0, time.UTC)
	expected := `2022-03-17 10:15`

	if hd != expected {
		t.Errorf(`TestHumanDate: %q should be %q`, hd, tt.expected)
	}
}
```

Man sieht: Das Go-Testing-Framework hat keinerlei Helfer-Funktion wie z.B.
[Jest](https://jestjs.io); Alles ist Handarbeit. Lediglich `t` muss man
erklären. `t` ist die Schnittstelle zur _Testing Engine._

Einige Regeln gibt es auch noch zu beachten.

1. Tests werden von _Go_ nur dann erkannt, wenn sie in Dateien stehen, auf die
   das Muster `*_test.go` passt.
2. Jeder Test muss mit `Test*`beginnen.
3. Jede Test-Funktion muss `*testing.T` als Argument akzeptieren.

## 2024-07-07 15:53

Jetzt ist auch die Sache mit dem CSRF-Abwehr im Kasten. Läuft! Einzelheiten im Buch!

## 2024-07-07 09:13

Die Autorisierung trägt jetzt alle Früchte: Auch das Erschleichen von ‘Create
Snippet’ ist jetzt unterbunden, dieses Mal mit einer neuen Middleware (siehe
Commit)

## 2024-07-07 09:13

Die Autorisierung trägt die ersten Früchte: Die Navbar hat ihre Scenarios des
‘Authorization’-Features implementiert.

## 2024-07-06 20:29

Wir kommen jetzt zum großen Thema __Autorisierung.__ Wir werden es genau so
aufdröseln wie wir es bei größeren Projekten immer machen – mit _Features._

```gherkin
Feature: Authorization

	As a logged-in user
	I want to be authorized
	In order to have exclusive functionality available for me.

# implemented
Scenario: create new snippet 
	Given I am logged-in
	Then I see a link 'Create Snippet' in the navbar
	When I click on the link
	Then I am being taken to the 'Create Snippet' page

Scenario: Sneaking into Creating a new snippet
	Given I am logged out
	When I try to use '/new/snippet' to sneak into the 'Create Snippet' form
	Then I am redirected to the 'Login' page.
	
# implemented
Scenario: Studying navbar as logged-in user
	Given I am logged in
	When I take a look at the navbar
	Then I see links to 'Home', 'Create Snippet' and 'Logout'

# implemented
Scenario: Studying navbar as logged-out user
	Given I am logged-out
	When I take a look at the navbar
	Then I see links to 'Home', 'Log In' and 'Sign Up'
```

## 2024-07-06 20:14

Logout-Feature implementiert. Sollte mir in _Obsidian_ eine Übersicht über die
verschiedenen Webapp-Techniken zusammenstellen und mit Beispielen aus diesem
Projekt unterlegen. So langsam verliere ich die Orientierung, und das ist ein
schlimmes Zeichen.

## 2024-07-05 22:08

Habe gerade den kompletten Login-Abschnitt aus dem Buch fertig gemacht. War
eigentlich gar nicht so schwer -- aber es wird lange dauern, bis ich hier
Routine habe.

Viel zu erklären gibt es da nicht. Es geht im Grunde um dieses Feature:

```gherkin
Feature: Login

	As a guest user
	I want to log in
	In order to have more rights and possibilities using this site

	Scenario: Opening the Login Form
		When I open the Login Form
		Then I see a text field to enter my email
		And I see a password field to enter my password
		And I see a button to submit the form to the backend.

	Scenario: Logging in with an empty email field
		Given I have filled out the fields
		When I submit with an empty email field
		Then the form is sent back to me
		And I am told the email field cannot be empty.

	Scenario: Logging in with an invalid email
		Given I have filled out the fields
		When I submit with an invalid email address
		Then the form is sent back to me
		And I am told the email field must contain a valid email address

	Scenario: Logging in with an empty password field
		Given I have filled out the fields
		When I submit with an empty password field
		Then the form is sent back to me
		And I am told the password field cannot be empty.

	Scenario: Logging in with an invalid password
		Given I have filled out the fields
		When I submit with an invalid passport 
		Then the form is sent back to me
		And I am told the password field must contain at least 8 characters

	Scenario: Logging in with an email address not found in the DB
		Given I have filled out the fields
		When I submit with both valid email and password data
		And the email is unknown to the database
		Then the form is sent back to me
		And I am told that email or password is incorrect

	Scenario: Successful Login
		Given I have filled out the fields
		When I submit with both valid email and password data
		And the email is known to the database
		And the password is recognized to be correct,
		Then I am redirected to the 'Create Snippet' page
		And I see no links for 'Login' or 'Sign Up' in the navbar
		And instead I see a link for 'Logout' in the navbar.
```

Folgende _Specs_ wurden dafür implementiert:

1. `Authenticate()` (Hilfsfunktion): Hier findet die eigentliche
   Authentifizierung statt. Die Kommentare erklären, was im einzelnen läuft.
2. `handleLogin()`. Auch hier erklärt sich vieles selbst. Aber eine Sache
   verdient es, noch besonders erläutert zu werden: `RenewToken()`. Hier wird
   _nur_ der Token erneuert, der die Session definiert; alle anderen Daten in der
   Session-Datenbank bleiben wie sie sind.


## 2024-07-01 17:21

Habe das Signup-Formular erfolgreich implementiert. Beim Validieren gab es aber eine wichtige Lektion zu lernen:

> _Für jedes Formularfeld sollte es nur eine Validierung und nur eine Fehlermeldung geben.

So wie `validator.Validator` und seine Methoden im Moment implementiert ist,

## 2024-06-29 17:45

Hab mir eine _„Under-Construction“_-Seite gegönnt. Brauche nur abschreiben, was
schon vorhanden war, und anpassen.

Besonders hervorzuheben war das „Nachrüsten“ von `templateData`:

```go
type templateData struct {
	CurrentYear int
	Flash       string
	Form        any
	Message     template.HTML   // <= 
	Snippet     TplSnippet
	Snippets    []TplSnippet
}
```

`Message` musste hier als `template.HTML` deklariert werden, damit ich später
HTML aufnehmen kann, das _nicht_ automatisch ‘escaped’ wird.

## 2024-06-28 10:48

Habe einen SigSegFault-Error-Marathon hinter mir und dabei eine sehr, sehr wichtige Lektion gelernt:

> [!abstract] 
> Für SQLc und für Session Management brauchen wir zwei verschiedene
> SQLite3-Datenbanken (nein nicht nur unterschiedliche Tabellen!)



## 2024-06-27: 21:58

Lobend erwähnen sollte man auch, dass ich die Ablaufzeit von Snippets nicht
mehr in der Datenbank berechnen lasse, sondern in _Go._ Ab jetzt wird das Feld
`ends` in der Datenbank nicht mehr berechnet, sondern einfach der Wert von
`expires` zurückgegeben und dann in der Hilfsfunktion `snippetExpired()`
darüber entschieden, ob ein Snippet abgelaufen ist oder nicht.

```go
return created.AddDate(0, 0, timeoutMap[expires]).Before(now)
```
Dieser Code erledigt alles

## 2024-06-27 21:45

Ich musste heute einen Fehler im `handleHome()`-Handler beheben: Bei der
Anzeige der Snippets kam es zu „blinden“ Leerzeilen in der Tabelle. Das war
richtig Kacke.

Was war das Problem? Diese Zeile in `app.RawSnippetsToTpl()`:

```go
lenRS := cap(rs)
lenRS = len(rs)
var tsp = make([]TplSnippet, lenRS)
```

Damit wurden 5 (aktuelle Anzahl der Snippet-Einträge in der Datenbank) oder gar
8 (Kapazität des []rs-Slices) Plätze für `tsp` reserviert – obwohl am Ende nur
2 gebraucht wurden, nämlich die zwei, die gegenwärtig nicht abgelaufen waren.

Dieser Code hier macht es jetzt richtig:

```go
var tsp = make([]TplSnippet, 0)

for _, r in range rs {
	// -> KEIN Leereintrag!
	if r.created.Valid {
		// ...
		tsp = append(tsp, TplSnippet{ /* ... */ }
	}
}
```

So bleiben die übrig, die tatsächlich angezeigt werden sollen.


## 2024-06-26 16:06

Ich habe jetzt das _SCS_-Paket von Alex Edwards für Session-Management
heruntergeladen und eingebaut. Zwei Dinge sind dabei besonders interessant:

#### Die Einbindung in das _Chi_-Framework

Die erfolgt mit Hilfe einer
[Group](https://go-chi.io/#/pages/routing?id=routing-groups). Im Code sieht das
so aus:

```go
// file: ./cmd/web/routes.go
func (app *Application) Routes() *chi.Mux {
	// ...
	
	// define a new subgroup with its own sub-router 'r'
	mux.Group(func(r chi.Router) {
		// middleware for this group
		r.Use(app.sessionManager.LoadAndSave)

		// routes
		r.Get(`/`, app.handleHome)
		// Endpoints with handlers as app methods
		r.Get(`/urlquery`, app.handleUrlQuery)
		r.Get(`/snippets`, app.handleSnippetList)
		r.Get(`/snippets/{id}`, app.handleSingleSnippetView)
		r.Get(`/new/snippet`, app.handleNewSnippetForm)
		r.Post(`/create/snippet`, app.handleNewSnippet)
	})

	// ...
}
```

_Chi_ erstellt hier hinter den Kulissen einen Sub-Router mit eigenen Routen und
einer eigenen Middleware für diese Routen (zusätzlich zu den allgemeinen
Middlewares), und bindet diesen Sub-Router hinterher in `mux` ein. Im Code ist
Was war das Problem? Diese Zeile in `app.RawSnippetsToTpl()`:

```go
lenRS := cap(rs)
lenRS = len(rs)
var tsp = make([]TplSnippet, lenRS)
```

Damit wurden 5 (aktuelle Anzahl der Snippet-Einträge in der Datenbank) oder gar
8 (Kapazität des []rs-Slices) Plätze für `tsp` reserviert – obwohl am Ende nur
2 gebraucht wurden, nämlich die zwei, die gegenwärtig nicht abgelaufen waren.

Dieser Code hier macht es jetzt richtig:

```go
var tsp = make([]TplSnippet, 0)

for _, r in range rs {
	// -> KEIN Leereintrag!
	if r.created.Valid {
		// ...
		tsp = append(tsp, TplSnippet{ /* ... */ }
	}
}
```

So bleiben die übrig, die tatsächlich angezeigt werden sollen.


## 2024-06-26 16:06

Ich habe jetzt das _SCS_-Paket von Alex Edwards für Session-Management
heruntergeladen und eingebaut. Zwei Dinge sind dabei besonders interessant:

#### Die Einbindung in das _Chi_-Framework

Die erfolgt mit Hilfe einer
[Group](https://go-chi.io/#/pages/routing?id=routing-groups). Im Code sieht das
so aus:

```go
// file: ./cmd/web/routes.go
func (app *Application) Routes() *chi.Mux {
	// ...
	
	// define a new subgroup with its own sub-router 'r'
	mux.Group(func(r chi.Router) {
		// middleware for this group
		r.Use(app.sessionManager.LoadAndSave)

		// routes
		r.Get(`/`, app.handleHome)
		// Endpoints with handlers as app methods
		r.Get(`/urlquery`, app.handleUrlQuery)
		r.Get(`/snippets`, app.handleSnippetList)
		r.Get(`/snippets/{id}`, app.handleSingleSnippetView)
		r.Get(`/new/snippet`, app.handleNewSnippetForm)
		r.Post(`/create/snippet`, app.handleNewSnippet)
	})

	// ...
}
```

_Chi_ erstellt hier hinter den Kulissen einen Sub-Router mit eigenen Routen und
einer eigenen Middleware für diese Routen (zusätzlich zu den allgemeinen
Middlewares), und bindet diesen Sub-Router hinterher in `mux` ein. Im Code ist
		// Endpoints with handlers as app methods
		r.Get(`/urlquery`, app.handleUrlQuery)
		r.Get(`/snippets`, app.handleSnippetList)
		r.Get(`/snippets/{id}`, app.handleSingleSnippetView)
		r.Get(`/new/snippet`, app.handleNewSnippetForm)
		r.Post(`/create/snippet`, app.handleNewSnippet)
	})

	// ...
}
```

_Chi_ erstellt hier hinter den Kulissen einen Sub-Router mit eigenen Routen und
einer eigenen Middleware für diese Routen (zusätzlich zu den allgemeinen
Middlewares), und bindet diesen Sub-Router hinterher in `mux` ein. Im Code ist
dieser Sub-Router als `r chi.Router` zu sehen; deshalb müssen auch alle
betroffenen Routen und die Middleware mit `r.` eingeleitet werden, sonst gibt
es eine _Panic!_

## 2024-06-25 16:28

Wir müssen uns mit einem Spezialfall von Error abgeben: Dass das Zielobjekt für
den Formular-Parser keine gültige Referenz auf ein gültiges Objekt ist; in
diesem Fall hätten wir richtig Scheiße gebaut! Zu diesem Zwecke und Behufe
haben wir `app.decodePostForm()` als Helfer geschrieben; falls dieser Fall
eintritt, gibt es eine _Panic,_ andernfalls einen normalen `error` und im
besten Fall ein `nil` als `error`.

## 2024-06-25 15:45

Das nächste Schmankerl – ein automatischer Formular-Parser! Dafür haben wir
`github.com/go-playground/form/v4@v4` installiert. Das hilft uns,
Formulareinträge nicht mehr mühsam „von Hand“ einlesen zu müssen, sondern
einfach mit Hilfe von _struct tags_ arbeiten zu können, so wie wir das aus der
Arbeit mit JSON kennen:

```go
type SnippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             string `form:"expires"`
	validator.Validator `form:"-"`
}
```

Die letzte Zeile mit `form:"-"` sagt dem `form`-Decoder, dass er
`validator.Validator` nicht in seine Arbeit mit einbeziehen soll.


## 2024-06-25 10:35

Das mit dem Validator hat jetzt auch in der Umsetzung geklappt. Siehe die
Änderungen in `app.handleNewSnippet`. Allerdings gibt es da noch eine Sache,
die neu für uns ist:

```go
type SnippetCreateForm struct {
	Title   string
	Content string
	Expires string
	validator.Validator
}
```

In der letzten Zeile wird einfach der Typ `validator.Validator` eingefügt,
__ohne__ ihm einen Schlüssel wie `Expires` zuzuweisen. Was bedeutet das? Es
bedeutet, dass eine Variable vom Typ `SnippetCreateForm` _auf alle Methoden und
alle Attribute von_ `validator.Validator` _zugreifen kann, als wären sie eigens
für_ `SnippetCreateForm` _geschrieben!_ Deshalb funktioniert Code wie hier:

```go
form := SnippetCreateForm{ /* ... */ }

form.CheckField(
	form.WithinRange(4, 20, form.Title),
	`Title`,
	`Title entry must be between 4 and 30 characters long.`,
)

// [...]

if !form.Valid() {
	data := app.buildTemplateData()
	data.Form = form
	// => in the template it's still going to be `.Form.FieldErrors.Title`
	app.Render(w, http.StatusUnprocessableEntity, `createSnippet.go.html`, data)
	return
}

```

## 2024-06-25 07:35

### Ein eigener Validator

In diesem Commit backen wir uns einen eigenen Validator – vorwiegend zu
Lernzwecken, damit wir verstehen, wie so ein Pakete funktioniert. Denn so wie
unser kleiner Validator funktioniert auch das
[Validator](https://thedevelopercafe.com/articles/payload-validation-in-go-with-validator-626594a58cf6)-Paket.

Zeichnen wir doch mal nach, wie es aufgebaut ist:

1. Im Zentrum steht der _Carrier_, `validator`. An ihm hängen alle Methoden,
   die für das `validator`-Paket wichtig sind.
2. Die `CheckField()`-Methode überprüft, ob ein Formulareintrag oder ein
   JSON-Eintrag gültig ist. Falls nicht, wird in `validator.FieldErrors` ein
   neuer Eintrag verbucht. Das folgende Beispiel erstellt einen Validator (`val`),
   überprüft `title` und verbucht den neuen Eintrag `val.FieldErrors['title']` mit
   `must be between 4 and 20 characters long`, falls der Test in
   `val.WithinRange()` durchfällt [Ausrufezeichen!].

```go
val = &validator.validator
title := "Currywurst"
key := `Title`
msg := `must be between 4 and 20 characters long`,

val.CheckField(val.WithinRange(4, 20, title), key, msg)
```

3. Alle anderen öffentlichen Methoden dienen der Wert-Überprüfung wie oben
   beschrieben.



## 2024-06-23 16:41

Hab jetzt eine selbstgebackene Validierung hinbekommen. Einfach durch logisches
Nachdenken und Lesen in der
[Validator](https://thedevelopercafe.com/articles/payload-validation-in-go-with-validator-626594a58cf6)-Dokumentation.

## 2024-06-22 19:19

Ich habe heute ganz anders gearbeitet. Im _Tai-Chi_-Modus. Bedeutet: Ich war
maximal konzentriert, war innerlich maximal entspannt und habe Fehler dadurch
vermieden, dass ich auf Langsamkeit gesetzt habe. All das, was die Arbeit so
belastend gemacht hat, war damit verschwunden. Und ich habe dadurch keine Zeit
verloren. Weil ich keine Zeit für Fehlersuche und Fehlerkorrektur aufbringen
musste.

Zur Sache: Ich habe ein Formular für die Erstellung eines neuen Snippets
erstellt. Nur das Formular. War mit den ganzen Vorarbeiten, was das Template
anging, ein Klacks!

```go
// A handler function to show a "Create Snippet" form
func (app Application) handleNewSnippetForm(w http.ResponseWriter, r *http.Request) {
	data := app.buildTemplateData()
	app.Render(w, http.StatusOK, `createSnippet.go.html`, data)
}
```

## 2024-06-20: 20:30

Habe das Projekt statt auf _httprouter_ auf [Chi](https://go-chi.io/)
umgestellt. Ohne Alex´ Hilfe, nur mit Hilfe der offiziellen Dokumentation. Hab
3 Stunden dafür eingeplant, am Ende war alles nach 30 min fertig. Geil!

## 2024-06-19 18:08

Anstelle des Extra-Pakets, das Alex Edwards vorgeschlagen hat, habe ich hier etwas eigenes
aufgebaut, das ich mir von [diesem Video](https://www.youtube.com/watch?v=H7tbjKFSg58) (Ab 07:43) geklaut habe.

```go
// file: ./cmd/web/middleware.go

// new type representing a middleware function
type Middleware func(http.Handler) http.Handler

func createMdwChain(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		// building a 'triangle' of nested middleware functions.
		// we are moving bottom-up in the 'xs' list that has been passed to us!
		for i := len(xs) - 1; i >= 0; i-- {
			// xs[i] is the current Middleware function in the list
			x := xs[i]
			// x(next) is the return value of the current Middleware function.
			// This is what creates the nested function calls!
			next = x(next)
		}
		return next
	}
}
```

Das funktioniert, weil die _Go Runtime_ beim Aufruf von `createMdwChain()` nur
den „äußersten Rahmen“, also nur die Closure zurückgibt	– ohne sie auszuführen.
Dafür „weiß“ die Closure an dieser Stelle schon ganz genau, woraus sich `xs`
zusammensetzt. Erst wenn diese „scharf gemachte Closure“ dann tatsächlich
aufgerufen wird, rattert ihre `for`-Schleife durch, erstellt das Paket aus den
eingeschachtelten Middleware-Funktionen und gibt dieses Paket dann zurück –
woraufhin `app.routes()` es umgehend an `main()` weiterreicht:

```go
func (app *Application) Routes() http.Handler {
	// [ ... do the routing stuff]

	// build and equip the closure, and then save its reference address in 'mwChain' 
	// => mwChain is now an executable function of type Middleware!
	mwChain := createMdwChain(
		app.recoverPanic,
		app.logRequest,
		secureHeaders,
	)
	
	// pass 'mux' to the closure, execute it and return its return value to main()
	return mwChain(mux)
}
```

## 2024-06-19 16:10

Der Standard-Weg in _Go,_ um mit einer _panic_ umzugehen, ist für ein
Web-Projekt ziemlich Kacke. Einfach abnippeln und das Frontend mit einer leeren
Response abspeisen – das kann es nicht sein für eine serviceorientierte App.

Hier die Middleware, die dieses Problem lösen soll:

```go
// middleware for well-formed death after panic
func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 'defer' guarantees that this func() will always be called,
		// even after a panic event.
		// -1-
		defer func() {
			// is there a panic to recover from? Well, in that case ...
			// -2-
			if err := recover(); err != nil {
				w.Header().Set(`Connection`, `close`)
				app.ServerError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
```

#### Anmerkungen:

1. `defer` plus ausgeführte Funktion sorgt dafür, dass die Funktion __auf jeden Fall__ ausgeführt wird, auch wenn eine _panic_ ausgeworfen wurde! Hier in diesem Fall ist das eine _anonyme_ Funktion, in der eine _panic_ mit Hilfe von `recover()` abgefangen und in einen geordneten `http.InternalServerError` umgewandelt wird, mit allem, was dazugehört. `defer` verlangt, dass diese Funktion im Code ausgeführt wird; deshalb die beiden `()` direkt hinter der Definition.
0. Wie gesagt: `recover()` dient dazu, die _panic_ abzufangen und in einen geordneten `http.InternalServerError` umzuwandeln.

## 2024-06-18 14:39

Hab jetzt auch die nächste Middleware implementiert. War mit dem neuen Snippet ein Klax!
Ohne die ganze Doku und das Update des Snippetswäre das in 10 min fertig geworden.

## 2024-06-18 13:24

Premiere: Wir haben die erste _Middleware_ eingefügt – um mit Hilfe von Headern
die Sicherheit im Browser zu erhöhen. Hat uns sechs Tage gekostet, weil ich
wichtige Vorarbeiten zu leisten hatte, z.B. wie man selbst geschriebene
_Snippets_ in Neovim-Kickstart erfolgreich integriert. Aber das ist nun geschafft.

Die Middleware findet sich in der Datei `cmd/web/middleware.go`, beim Einbau in
der `routes()`-Funktion gibt es allerdings einige Dinge für den Umgang mit
Middleware zu beachten:

```go
// -1-
func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	// [ setting up mux ]

	// including 'mux' into the general 'secureHeaders()' middleware
	// -2-
	return secureHeaders(mux)
}
```
#### Anmerkungen

1. Wenn eine allgemeine Middleware zum Einsatz kommt, müssen wir den
   Rückgabetyp von `app.Routes()` von `*http.ServeMux` nach `http.Handler` ändern.; 
   in der Sache ändert das nichts, weil http.Handler ein _Interface_ ist und 
   `*http.ServeMux` dieses Interface implementiert hat.
0. Wie wir in _Obsidian_ schon geklärt haben, besteht das Wesen von Middleware
   darin, den „nächsten Facharbeiter“ als Argument „einzuklammern“. Hier
   klammert `secureHeaders()` das `mux`-Objekt ein, und da es das komplette 
   `mux`-Objekt einklammert, ist es eine allgemeine Middleware.


## 2024-06-12 19:36

Nice to know: das _template_-Package aus der Standard-Lib hat noch ein
besonderes Schmankerl: das `FuncMap`-Objekt. Mit dieser Funktion können wir
eigene Template-Funktionen für unsere Arbeit innerhalb der Templates hinzufügen
– und die funktionionieren dann genauso wie die etablierten Template-Funktionen
wie `printf` oder `len` oder `index` (vgl. Abschnitt 5.2 im Buch)

Hier das Rezept (im Rahmen dieses Projekts)

```go
// define the function we want as template function
// -1-
func tplFoo(arg1 string, num int) string {
	// [...]
}

// initialize a FuncMap Object as a global variable and add the function we
// just defined
// -2-
var tplFunctions = template.FuncMap{ `tplFoo`: tplFoo }


func buildTemplateCache() (map[string]*template.Template, error) {
	// [...]
	// for each page template
	for _, page := range pages {
		// build a new template set from scratch
		// -3-
		ts := template.New(name)
		// add the FuncMap Object from above ...
		// -3-
		ts.Funcs(tplFunctions)
		// ... BEFORE you parse the first template
		ts.ParseFiles(`./ui/html/base.go.html`)
		if err != nil {
			return nil, err
		}
		// [...]
	}
	// [...]
}

// use the function in the template 
// -4-
// {{ tplFoo 'theString' 42 }}

```

#### Anmerkungen

1. Die neue Template-Funktion darf so viele Argumente aufnehmen wie wir wollen,
   aber sie darf nur __einen__ Rückgabewert haben, und der sollte _dringend
   entweder vom Typ `string` oder vom Typ `bool` oder vom Typ `int` sein!
2. `tplFunctions` müssen global sein, damit jede Funktion im Package darauf
   zugreifen kann.
3. `template.Funcs` funktioniert nur bei einem bereits existierenden
   `template`-Objekt, deshalb müssen wir es vorher mit `New()` erstellen.
4. In diesem Beispiel ist der Rückgabewert ein `string`, und deshalb kann er
   direkt im Template eingetragen werden.

Übrigens: Nach dieser Methode wurden alle speziellen Template-Funktionen in
_Hugo_ implementiert.


## 2024-06-12 18:51

In diesem Commit habe ich _Shared Data_ in einem sehr einfachen Fall
implementiert: Es ging darum, das aktuelle Jahr im Footer anzuzeigen.
Änderungen in diesen Dateien waren dafür nötig:

- `cmd/web/templates.go`: hier wurde `templateData` erweitert.
- `cmd/web/helpers.go`: hier wurde eine Factory Function angelegt (`app.buildTemplateData()`)
- `cmd/web/handlers.go`: hier wurde die Factory Function aufgerufen
- `ui/html/base.go.html`: hier wurde in Änderung im Template umgesetzt


## 2024-06-12 17:06

habe die `buildTemplateCache()` jetzt derart erweitert, dass sie automatisch
_alle_ möglichen Partials in den Template-Cache aufnimmt.

Ganz besonders dieser Trick hier war dabei hilfreich:

```go
// [...]

// parse the base template file into a template set
ts, err := template.ParseFiles(`./ui/html/base.go.html`)
if err != nil {
	return nil, err
}

// parse all possible partials files into this very template set to add them there
// notice how the 'old' ts on the right side creates a new enhanced ts
// on the left hand of '='
ts, err = ts.ParseGlob(`./ui/html/partials/*.go.html`)
if err != nil {
	return nil, err
}

// [...]
```

Im oberen Block wird `ts` neu angelegt. Im unteren Block wird _genau dieses_
`ts` benutzt, um die Partials hinzu zu addieren und sich selbst „neu zu
erfinden“. 


## 2024-06-12 11:18

Großes Refactoring: Ich habe den Render-Vorgang der Templates in einer
allgemeinen Hilfsfunktion (`app.Render()`) zusammengefasst. In den Handlern
ließ sich dadurch eine Menge Code einsparen (vgl. die Veränderungen in diesem
Commit).

Folgende Dateien waren betroffen:

- `cmd/web/handlers.go`
- `cmd/web/helpers.go`
- `cmd/web/main.go`
- `cmd/web/templates.go`


## 2024-06-10 17:56

Jetzt zur eigentlichen Doku:

Ich musste noch einen Konvertierungs-Helfer schreiben, der das „Rohergebnis“
der SQL-Abfrage in ein wohlgeformtes Template-Objekt überführt. Ist am Ende
ziemlich einfach geworden.

Die Arbeit mit Templates erfolgt immer nach dem gleichen Schema. Ich habe es
unten beschrieben. Das ist nervig, aber wenn man sich mal daran gewöhnt hat,
ist es OK.

## 2024-06-10 17:40

Hab mich in den totalen Frust gearbeitet – weil ich ein paar entscheidende Denkfehler gemacht habe:

1. Wenn ich im Go-Code „völlig unerklärlicherweise“ gesagt bekomme, dass eine
   Methode oder ein Attribut undefined ist, obwohl ich sie nachweislich definiert habe,
   dann habe ich wahrscheinlich das Objekt mit dem Typen verwechselt:

```go
// WRONG:
// createdTpl = db.GetSnippetRow.Created.Time.Format("2006-01-02 03:04:05")
// CORRECT:
createdTpl = r.Created.Time.Format("2006-01-02 03:04:05")
// yes, 'r' is a 'db.GetSnippetRow', that's correct ...
```


2. Ein riesiges Problem erwuchs mir daraus, dass ich folgenden Fehler gemacht habe:
   im falschen Beispiel ist `http.Handle` nicht mit unserem extra erstellten
   Spezial-Router verdrahtet, sondern mit `http.DefaultMux`! Deshalb wurde z.B.
   die CSS-Datei nicht mehr an den Browser geladen ...

```go
// WRONG
// http.Handle(`/static/`, http.StripPrefix(`/static`, fileServer))
// CORRECT
mux.Handle(`/static/`, http.StripPrefix(`/static`, fileServer))
```


## 2024-06-10 05:23

Es gibt noch eine Ergänzung zu SQLc: die Fehlerbehandlung.

Was passiert, wenn eine Abfrage kein Ergebnis liefert? In diesem Fall gibt SQLc
einen `sql.ErrNoRows` zurück. Auf den können wir entsprechend reagieren:

```go
resultRaw, err := db.Qs.GetSnippet(ctx, idDB)
if err != nil {
	if errors.Is(err, sql.ErrNoRows) {
		// => 404 error
		app.NotFound(w)
	} else {
		// => 500 error
		app.ServerError(w, err)
	}
	return
}
```


## 2024-06-07 07:14

SQLc wurde erfolgreich installiert und integriert; hier ist die Dokumentation
über die Features, die SQLc zur Verfügung stellt.

```go
// parameter type for InsertSnippet()
type InsertSnippetParams struct {
	Title   string
	Content string
	Expires sql.NullString
}

// return type from InsertSnippetParams()
type InsertSnippetRow struct {
	ID    int64
	Title string
}

func (q *Queries) InsertSnippet(ctx context.Context, arg InsertSnippetParams) (InsertSnippetRow, error) {
	// [...]
}

// returns all snippets
func (q *Queries) GetAllSnippets(ctx context.Context) ([]GetAllSnippetsRow, error) {
	// [...]
}

func (q *Queries) GetSnippet(ctx context.Context, id int64) (GetSnippetRow, error) {
	// [...]
}
```

Es ist eigentlich nicht schwer, und es ist auch ganz logisch. Die Typen für die
Leseausgabe haben wir schon geklärt. Wichtig ist nur noch, dass wir den
_context_ mit einbeziehen müssen. Das geht so:

```go
import "context"

func (app *Application) handleSingleSnippetView(w ResponseWriter, r *http.Request) {
	ctx = context.Background()
	rawSnippet, err = db.Qs.GetSnippet(ctx, 2)
	if err != nil {
		// => "Ach Scheiße!"
	}
	// [...]
}
```

## 2024-06-06 18:09

Ich habe jetzt die Sache mit `Sqlc` fürs erste klar gemacht. War schwer genug,
weil ich nicht einfach dem Buch folgen konnte, sondern irgendwo die richtigen
SQL-Abfragen finden, abschreiben und anpassen musste. Sie stehen in
`./internal/schema.sql` und `./internal/query.sql`.

Wie es Sqlc's Art ist, hat es daraus in `./internal/db/' Go-Code generiert, den
wir jetzt einfach nur noch abrufen müssen. Allerdings gibt es noch einige Dinge
dringend zu beachten:

```go
// file: ./internal/db/query.sql.go

// return type for single snippet query
type GetSnippetRow struct {
	ID      int64
	Title   string
	Content string
	Created sql.NullTime
	Ends    interface{}
}

// return type for snippet list query.
type GetAllSnippetsRow struct {
	ID      int64
	Title   string
	Content string
	Created sql.NullTime
	Ends    interface{}
}
```

Das Problem sind die Zeitangaben. `Created` ist als `sql.NullTime` deklariert,
`Ends` sogar als `interface{}`.


#### sql.NullTime

Jaja, _Go_ und die leidigen `NULL`-Einträge in Datenbanken! Wie soll _Go_ mit
sowas umgehen? Das `database/sql`-Paket bietet für diese Fälle `Null`-Typen an.
Im Fall von `NullTime` ist dieser Typ so definiert:

```go
// inside the StdLib database/sql package

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}
```

Das bedeutet für uns in diesem Fall:

1. Wir müssen überprüfen, ob `NullTime.Valid` `true` oder `false` ist:
2. Wenn `true`, gibt es einen Eintrag, und wir kommen mit `NullTime.Time` heran.
3. Wenn `false`, ist dieser Eintrag leer und wir müssen die entsprechenden
   Konsequenzen ziehen.

#### interface{}

Um den Typ `interface{}` in einen String zu überführen (und der Zeitstempel ist
am Ende nur ein String), genügt folgender Trick:

```go
func main() {
	// same format as our 'ends' column
	var endTimeEntry interface{} = "2018-12-13 18:43:00"
	// interface{} -> String 
	endTimeString := fmt.Sprintf("%v", endTimeEntry)
	// str -> time conversion
	endTime, err := time.Parse("2006-01-02 03:04:05", endTimeString)

	if err != nil {
		fmt.Println(`Fuck, man!`, err)
	}

	// use endTime
}
```

Das funktionioniert natürlich nur, weil wir genau wissen, dass wir einen
Zeitstempel in genau diesem Format von unserem Datenbankmodul geliefert
bekommen!


## 2024-06-04 12:10

Habe alles, was mit `http.ServeMux` zu tun hat, aus `main()` in die Funktion
`app.Routes()` ausgelagert. (Nebenbei habe ich jetzt auch _GoPls_ in Neovim
integriert. Der Programmier-Spaß hat sich damit verzehnfacht! Läuft!)

## 2024-06-04 10:28

Wir haben einige Helfer für Fehlerbehandlung implementiert. Die Doku steht vollständig in
`./cmd/web/helpers.go`, und zwar für `ServerError()`, `ClientError()` und `NotFound()`.

## 2024-06-05 08:45

Die Sache mit der _dependency injection_ funktioniert nur solange wie alle
Handler im _main_-Package sich befinden. Wenn sie auf mehrere Packages verteilt
sind (dumm eigentlich; es reicht, sie auf mehrere Dateien im _main_-Package zu
verteilen), braucht es neue Lösungen.

_Let’s Go_ hat dafür auf S. 146 eine Lösung parat.

## 2024-06-04 19:09

### Implementierung der dependency injection

Als erstes braucht es einen Datentyp, in dem alle globalen Status-Informationen vermerkt sind:

```go
// file: ./cmd/web/main.go

// Define an 'Application' struct to hold global status for the application.
// For now we will only include fields for the two custom loggers, but this one
// will grow and grow and grow over the course of the project.
type Application struct {
	ErrLog  *log.Logger
	InfoLog *log.Logger
}
// ...
```

Dann muss ein Objekt dieses Datentyps in `main()` eingeführt werden: 


```go
// file: ./cmd/web/main.go

func main() {
	// ...

	// introduce the infolog and the errLog instances
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// introduce the app Object in order to grant access to the global
	// application state.
	app := &Application{
		ErrLog:  errLog,
		InfoLog: infoLog,
	}
	
	// ...
}
```

Als nächstes müssen sämtliche Handler als `app`-Methoden umgeschrieben werden:

```go
// file: ./cmd/web/handlers.go

func (app *Application) handleHome(w http.ResponseWriter, r *http.Request) {
	// ...
}

func (app *Application) handleNewSnippet(w http.ResponseWriter, r *http.Request) {
	// ...
}

// etc.
```

Dann müssen auch die Routen auf die neuen Realitäten eingestellt werden:

```go
// file: ./cmd/web/main.go

func main() {
	// ...

	// Endpoints with handlers as app methods
	mux.HandleFunc(`GET /`, app.handleHome)
	mux.HandleFunc(`GET /urlquery`, app.handleUrlQuery)

	mux.HandleFunc(`GET /snippets/{id}`, app.handleSingleSnippetView)
	mux.HandleFunc(`POST /snippets/new`, app.handleNewSnippet)

	// ...
}
```

Und schließlich gilt es, auch den Code in den Handlern anzupassen. Beispiel:

```go
// file: ./cmd/web/handlers.go

func (app *Application) handleHome(w ResponseWriter, r *http.Request) {
	// ...

	ts, err := template.ParseFiles(templates...)
	if err != nil {
		// Because 'handleHome()' is now an 'app' method, it can access its fields,
		// including the error logger. ⇒ We use this logger now instead of the standard logger.
		app.ErrLog.Println(err.Error())
		http.Error(w, `Internal Server Error has occured!`, http.StatusInternalServerError)
		return
	}

	// ...
}
```

__YEEEEHAAAW!__

## 2024-06-04 16:00

### Globaler Status

Es gibt einige Dinge, die lassen sich nicht auf eigene Funktionen und eigene
Handler beschränken, z.B., welche wichtigen Header-Informationen in der
_Request_ mitgeliefert wurden, ob ein Benutzer eingeloggt ist oder nicht. Ob er
Administrator ist oder nicht. Ob er gerade im System gesperrt ist oder nicht.
Welche ID er gerade hat. Das muss irgendwie _global_ gelöst werden, und
`main()`, jeder Handler und jede Hilfsfunktion sollte irgendwie Zugriff auf
diese allgemeinen, globalen Informationen bekommen – wenigstens Lesezugriff
sollte sie haben.

Um dieses Problem zu lösen, gibt es mehrere Wege. Frameworks wie
[Echo](https://echo.labstack.com) oder [Fiber](https://docs.gofiber.io) bieten
dafür das _Context_-Objekt an, entweder als
[`echo.Context`](https://echo.labstack.com/docs/context) oder als
[`fiber.Ctx`](https://docs.gofiber.io/api/ctx).

#### Die Puristen-Lösung: `context` (StdLib)

Vanilla-Puristen versuchen dagegen, uns zu diesem Zweck und Behufe auf das
hauseigene [Context](https://pkg.go.dev/context)-Paket einzuschwören. Geht im
Grunde genommen ganz einfach:

```go
import (
	"context"
	"fmt"
)

func doSomething(ctx context.Context, key String) {
	fmt.Printf("doSomething: myKey's value is %s\n", ctx.Value(key))
}

func main() {
	ctx := context.Background()
	// add a new Context value; returns a deep copy of the existing ctx object
	ctx = context.WithValue(ctx, "foo", "Foo Value")
	// use the new value
	doSomething(ctx, "foo")
}
```

Der _value_ ist bei `WithValue()` als `any` deklariert; es kann also alles
Mögliche sein, was in Go gebaut werden kann. Bei großen Objekten empfehlen sich
für den _value_ _Pointer_ auf diese Objekte.

All diesen Lösungen ist gemein, dass sie einen _Vertrag_ anbieten: „Alles, was
du brauchst, findest du in unserem Context-Objekt!“. Ein solcher Vertrag muss
natürlich intensiv dokumentiert sein -- auch der Vertrag, den der
Vanilla-Purist mit Hilfe des _Context_-Menüs aus der _Standard Library_ baut.
Sonst kann keiner der Teamkollegen das Context-Objekt verstehen, geschweige
verwenden!

#### Alex Edward's Lösung: _Dependency injection_

Das Wort ‘injection’ bedeutet, dass etwas „eingespritzt“ wird. Bei Alex läuft
das so, dass jeder Handler eine Methode für das `app`-Objekt wird, in dem dann
alle globalen Daten drinstehen (wieder ein Vertrag, der intensivst dokumentiert 
werden muss).



## 2024-06-04 09:44

Um auch die internen Fehlermeldungen an unseren neuen `errLog` weiterzugeben,
müssen wir `main()` wieder um einen Punkt ergänzen. Wir brauchen eine selbstgebackene 
`http.Server`-Instanz mit eigenen Einstellungen. Der Code erklärt, wie es geht:

```go
func main()
	// ...

	// Initialize a new http.Server instance. We set the Addr and the Handler fields so
	// that the server uses the same network address and and routes as before,
	// and set the 'ErrorLog' field so that the server now uses the custom errLog logger
	// in case a bug lurks its head in this app.
	// - 1- 
	srv := &http.Server{
		Addr:     *port,
		ErrorLog: errLog,
		Handler:  mux,
	}

	infoLog.Printf("starting server at port %s", *port)
	// now call the 'ListenAndServe()' method of our own http.Server version
	err := srv.ListenAndServe()
	if err != nil {
		errLog.Fatalf("Uh oh! %s", err)
	}
}
```

#### Anmerkungen

1. Genau so erschafft man in Go ein „Objekt“: Man nimmt einen bestehenden
   Datentyp (`struct`) und erstellt etwas mit einer eigenen Speicheradresse
   (`&`), dessen Referenz man an eine Variable zurückgibt – so wie das hier bei
   `srv` passiert.


## 2024-06-04 09:04

Hier geht es darum, hausgeschneiderte Log-Botschaften zu ermöglichen. Der
folgende Code zeigt, wie es geht.

```go
func main() {
	// ... 

	// Use log.New() to create a logger for writing information messages.
	// it takes three parameters:
	//    - the destination to write the log to (os.Stdout)
	//    - a string prefix for the message ('INFO\t')
	//    - flags to indicate what additional messages to include (local date
	//      and time). Notice that the flags are connected with the pipe symbol
	//      '|'.
	infoLog := log.New(os.Stdout, `INFO\t`, log.Ldate|log.Ltime)

	// Now we create a logger for writing error messages in the same way, but use
	// os.Stderr as the destination and use log.Lshortfile flag to include the
	// relevant file name and line number
	errLog := log.New(os.Stderr, `ERROR\t`, log.Ldate|log.Ltime|log.Lshortfile)

	// ...

	// Write messages now using the new loggers, instead of the standard logger
	infoLog.Printf("starting server at port %s", *port)
	err := http.ListenAndServe(*port, mux)
	if err != nil {
		errLog.Fatalf("Uh oh! %s", err)
	}
}
```


## 2024-06-04 08:05

Mit Hilfe von _Command Line Flags_ können wir spezielle Einstellungen auf der
Kommandozeile vornehmen und sind nicht von hart einkodierten Einstellungen
abhängig.

Das folgende Beispiel zeigt, wie man mit einer Kommandozeilen-Option den Port
für die Anwendung ändern kann:

```go
func main() {
	// ...

	// Define a new command line flag with the name 'addr' and a default value
	// of ':3000' and a short help text to tell what this flag is doing.
	port := flag.String(`port`, `:3000`, "setting the port number")

	// Now we have to use the flag.Parse() function to parse the command-line flag.
	// This reads in the command line flag value and assigns it to the 'port' variable.
	// We need to call this **before** we use the 'port' variable; otherwise the value
	// will always be ':3000'.
	// If any errors occur, the application will panic.
	flag.Parse()

	// ...

	// The value returned from flag.String() is a pointer to the flag value,
	// not the value itself. So we need to dereference the pointer. To make
	// this work properly, Println() must become Printf()
	log.Printf("starting server at port %s", *port)
	err := http.ListenAndServe(*port, mux)

	// ...
}
```

Ab jetzt können wir die App so wie hier aufrufen, und es wird einen neuen Port geben:

```bash
$ go run cmd/web/ -port=':9999'
```

## 2024-06-03 17:33

Mit folgenden Zeilen können wir statische Dateien wie Bilder, CSS-Dateien oder JavaScript-Dateien
auf unserer Webseite laden:

```go
// file: cmd/web/handlers.go
func handleHome(w http.ResponseWriter, r *http.Request) {
    // ...

	// Create a file server that serves static files out of './ui/static/'. The
	// path here is relative to the project directory root.
	fileServer := http.FileServer(http.Dir(`./ui/static/`))

	// Register the fileServer for all URL paths that start with '/static/'.
	// For matching paths, we strip the '/static' prefix before the request
	// reaches the fileServer.
	http.Handle(`/static/`, http.StripPrefix(`/static`, fileServer))

	// ...
}
```

Alles Wesentliche findet sich in den Kommentaren.



## 2024-06-03 08:40

Wir haben die Templates ein wenig umgestellt, um _partials_ und _Snippets_ möglich zu machen:

```html
<!-- file: ui/web/base.go.html -->
{{ define "base" }}

  <!-- HTML page layout -->

  {{ template "main" . }}

  <!-- More HTML page layout -->
{{ end }}
```

In der angegebenen Datei `ui/web/base.go.html` haben wir also den “base”-Block
definiert. Es ist nachher unerheblich in welcher Datei genau dieser Block
definiert wird. Wichtig ist nur, dass er definiert _ist,_ wenn
`ExecuteTemplate()` im Handler aufgerufen wird.

Innerhalb des “base”-Blocks wird der “main”-Block an genau der bezeichneten
Stelle hineingeladen. Es ist auch hier nachher unerheblich, in welcher Datei
genau dieser Block definiert wird. Wichtig ist auch hier nur, dass er definiert
_ist,_ wenn `ExecuteTemplate()` im Handler aufgerufen wird.

Zum Schluss mussten noch Änderungen im Handler vorgenommen werden

```go
func handleHome(w http.ResponseWriter, r *http.Request) {
	// - 1 -
	templates := []string{
		"./ui/html/base.go.html",
		"./ui/html/pages/home.go.html",
	}

	// - 2 -
	ts, err := template.ParseFiles(templates...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, `Internal Server Error has occured!`, http.StatusInternalServerError)
	}

	// - 3 -
	err = ts.ExecuteTemplate(w, `base`, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, `Template Error. WTF?!!`, 500)
	}
```

#### Anmerkungen

1. Da für den Aufbau dieser Seite mehr als eine Template-Datei benötigt wird,
   sammeln wir die Dateipfade in einem _Slice._
2. Diesen _Slice_ übergeben wir an `ParseFiles()`, und zwar mit drei Punkten
   `templates...`. Dieses Konstrukt entspricht genau dem _Spread Operator_ in
   JavaScript. Der Slice wird damit in seine Elemente aufgespalten, und diese
   Elemente bilden den „Rattenschwanz“ für die `ParseFiles()`-Methode.
3. aus `ts.Execute()` wurde `ts.ExecuteTemplate()`, und `ExecuteTemplate()`
   braucht als zweiten Parameter den Template-Block, der den äußeren Rahmen für
   die geplante HTML-Einheit bildet. In unserem Fall ist das der “base”-Block.

So geht das in der “Vanilla”-Edition mit Go Templates. Wir haben es einmal
erfolgreich durchgespielt. Halleluja!

## 2024-06-03 07:43

Ich hab das erste Golang-Template geladen und vorher das Projekt umorganisiert. Ging 
erstaunlich reibungslos. Die Landing Page hat noch kein CSS, aber das wird sich bald ändern.

Hier der Code für den `handleHome()`-Handler:

```go
func handleHome(w http.ResponseWriter, r *http.Request) {
	// exclude anything but root as endpoint
	if r.URL.Path != `/` {
		http.NotFound(w, r)
		return
	}

	// template.ParseFiles() reads the templates into a template set.
	// If there is an error, we log the detailed error message on the terminal
	// and use the http.Error() function to send a generic 500 server error.
	ts, err := template.ParseFiles(`./ui/html/pages/home.go.html`)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, `Internal Server Error has occured!`, http.StatusInternalServerError)
	}

	// Since we made it here, we use the Execute() method on the template set
	// to write the template content as the response body. The last parameter
	// of Execute() represents any dynamic data we want to pass in; at the
	// moment it will be nil.
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, `Template Error. WTF?!!`, 500)
	}
}
```

#### Anmerkungen:

In Go's _Vanilla Templates_ 

- muss als erstes in jedem Handler ein eigenes _Template Set_ erstellt werden.
  Dieses _Template Set_ muss alle Template-Dateien beinhalten, die für das
  Erstellen dieser besonderen Seite benötigt werden: _base,_ alle _partials,_
  alle _snippets_ – einfach alles, was dazugehört!
- wird danach das _Template Set_ zu einer vollständigen HTML-Einheit (kann auch
  ein HTMX-Snippet sein!) zusammengefasst und rausgeschickt (dafür sorgt der
  `w`-Parameter)
 

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
