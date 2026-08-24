package middleware

import (
	"net/http"
)

func ChangeLanguage(response http.ResponseWriter, request *http.Request) {
	language := request.URL.Query().Get("lang")

	if !SupportedLanguages[language] {
		http.Error(response, "Langue non supportée", http.StatusBadRequest)
		return
	}

	cookie := &http.Cookie{
		Name:     "nmw-language",
		Value:    language,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(response, cookie)

	referer := request.Header.Get("Referer")

	if referer == "" {
		referer = "/"
	}

	http.Redirect(response, request, referer, http.StatusSeeOther)
}
