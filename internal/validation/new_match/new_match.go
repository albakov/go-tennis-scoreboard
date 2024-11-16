package new_match

import (
	"fmt"
	"github.com/albakov/go-tennis-scoreboard/internal/controller"
	"net/http"
	"regexp"
)

type RequestNewMatch struct {
	r            *http.Request
	fields       map[string]string
	errorMessage string
}

func NewMatch(r *http.Request, fields map[string]string) *RequestNewMatch {
	return &RequestNewMatch{
		r:      r,
		fields: fields,
	}
}

func (cc *RequestNewMatch) Validate() {
	for field := range cc.fields {
		v := cc.r.FormValue(field)
		if v == "" {
			cc.errorMessage = fmt.Sprintf(controller.MessageFieldEmpty, field)

			return
		}

		matched, err := regexp.Match(`[A-Za-z]\. [A-Za-z]+`, []byte(v))
		if !matched || err != nil {
			cc.errorMessage = fmt.Sprintf(controller.MessageFieldDoesntMatchThePattern, field)

			return
		}

		cc.fields[field] = v
	}
}

func (cc *RequestNewMatch) IsValid() bool {
	return cc.errorMessage == ""
}

func (cc *RequestNewMatch) ErrorMessage() string {
	return cc.errorMessage
}

func (cc *RequestNewMatch) Field(field string) string {
	return cc.fields[field]
}
