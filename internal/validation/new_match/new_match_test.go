package new_match

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func getRequestForTest(form url.Values) *http.Request {
	return &http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Scheme:      "http",
			Opaque:      "",
			User:        nil,
			Host:        "localhost",
			Path:        "/new-match",
			RawPath:     "",
			OmitHost:    false,
			ForceQuery:  false,
			RawQuery:    "",
			Fragment:    "",
			RawFragment: "",
		},
		Body:             io.NopCloser(strings.NewReader(form.Encode())),
		GetBody:          nil,
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Host:             "localhost",
		Form:             form,
		PostForm:         nil,
		MultipartForm:    nil,
		Trailer:          nil,
		RemoteAddr:       "",
		RequestURI:       "",
	}
}

func TestEmptyFields(t *testing.T) {
	form := url.Values{}
	form.Add("playerOne", "")
	form.Add("playerTwo", "")

	validator := NewMatch(getRequestForTest(form), map[string]string{"playerOne": "", "playerTwo": ""})
	validator.Validate()

	if validator.IsValid() {
		t.Errorf("IsValid() should return false")
	}
}

func TestWrongPattern(t *testing.T) {
	form := url.Values{}
	form.Add("playerOne", "F. Alonso")
	form.Add("playerTwo", "John@Doe.com")

	validator := NewMatch(getRequestForTest(form), map[string]string{"playerOne": "", "playerTwo": ""})
	validator.Validate()

	if validator.IsValid() {
		t.Errorf("IsValid() should return false")
	}
}

func TestFieldsOk(t *testing.T) {
	form := url.Values{}
	form.Add("playerOne", "F. Alonso")
	form.Add("playerTwo", "J. Doe")

	validator := NewMatch(getRequestForTest(form), map[string]string{"playerOne": "", "playerTwo": ""})
	validator.Validate()

	if !validator.IsValid() {
		t.Errorf("IsValid() should return true")
	}
}
