package wraphttpauth_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/go-on/wrap-contrib/third-party/wraphttpauth"

	"github.com/abbot/go-http-auth"
	"github.com/go-on/wrap"
)

func secretBasic(user, realm string) string {
	if user == "john" {
		// password is "hello"
		return "$1$dlPL2MqE$oQmn16q49SqdmhenQuNgs1"
	}
	return ""
}

func secretDigest(user, realm string) string {
	if user == "john" {
		// password is "hello"
		return "b98e16cbc3d01734b264adba7baa3bf9"
	}
	return ""
}

func app(rw http.ResponseWriter, req *http.Request) {
	var authReq auth.AuthenticatedRequest
	rw.(wrap.Contexter).Context(&authReq)
	rw.Write([]byte("user " + authReq.Username + " authenticated"))
	return
}

func ExampleBasic() {
	stackBasic := wrap.New(
		context{},
		wraphttpauth.Basic("example.com", secretBasic),
		wrap.HandlerFunc(app),
	)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.SetBasicAuth("john", "hello")
	stackBasic.ServeHTTP(rec, req)

	fmt.Println("-- success --")
	fmt.Printf("code %d\n", rec.Code)
	fmt.Println(rec.Body.String())

	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/", nil)
	stackBasic.ServeHTTP(rec, req)
	fmt.Println("-- fail --")
	fmt.Printf("code %d\n", rec.Code)

	// Output:
	// -- success --
	// code 200
	// user john authenticated
	// -- fail --
	// code 401
}

func ExampleDigest() {

	stackDigest := wrap.New(
		context{},
		wraphttpauth.Digest("example.com", secretDigest),
		wrap.HandlerFunc(app),
	)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	stackDigest.ServeHTTP(rec, req)
	fmt.Println("-- fail --")
	fmt.Printf("code %d\n", rec.Code)

	authServerHeader := digestMap(rec.Header().Get("WWW-Authenticate"))

	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", authHeaderForClient(authServerHeader))

	rec = httptest.NewRecorder()

	stackDigest.ServeHTTP(rec, req)

	fmt.Println("-- success --")
	fmt.Printf("code %d\n", rec.Code)
	fmt.Println(rec.Body.String())

	// Output:
	// -- fail --
	// code 401
	// -- success --
	// code 200
	// user john authenticated

}

func authHeaderForClient(authServerHeader map[string]string) string {
	HA1 := auth.H("john:" + authServerHeader["realm"] + ":hello")
	HA2 := auth.H("GET:/")
	nc, cnonce := "0", "NjE4MTM2"
	response := auth.H(strings.Join([]string{HA1, authServerHeader["nonce"], nc, cnonce, authServerHeader["qop"], HA2}, ":"))

	return fmt.Sprintf(
		`Digest username="%s", realm="%s", nonce="%s", uri="%s", qop="%s", opaque="%s", response="%s", algorithm="MD5",  nc=%s, cnonce="%s"`,
		"john",
		authServerHeader["realm"],
		authServerHeader["nonce"],
		"/",
		authServerHeader["qop"],
		authServerHeader["opaque"],
		response,
		nc,
		cnonce,
	)

}

func digestMap(header string) map[string]string {
	result := map[string]string{}
	idx := strings.Index(header, "Digest")
	digeststr := strings.TrimSpace(header[idx+7:])

	for _, pair := range strings.Split(digeststr, ",") {
		a := strings.SplitN(pair, "=", 2)
		result[strings.TrimSpace(a[0])] = strings.Trim(strings.TrimSpace(a[1]), `"`)
	}
	return result
}
