package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const languageKey contextKey = "language"

const (
	DefaultLanguage = "fr"
)

var SupportedLanguages = map[string]bool{
	"fr": true,
	"en": true,
}

func Language(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {

		language := DefaultLanguage

		cookie, err := request.Cookie("nmw-language")

		if err == nil && SupportedLanguages[cookie.Value] {
			language = cookie.Value
		}

		ctx := context.WithValue(
			request.Context(),
			languageKey,
			language,
		)

		next.ServeHTTP(response, request.WithContext(ctx))
	})
}

func GetLanguage(request *http.Request) string {
	language, ok := request.Context().Value(languageKey).(string)

	if !ok || !SupportedLanguages[language] {
		return DefaultLanguage
	}

	return language
}
