package session

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
)

type Session struct {
	CookieLifeTime string `json:"cookie_lifetime"`
	CookiePersist  string `json:"cookie_persist"`
	CookieName     string `json:"cookie_name"`
	CookieDomain   string `json:"cookie_domain"`
	CookieSecure   string `json:"cookie_secure"`
	SessionType    string `json:"session_type"`
}

func (s *Session) InitSessions() *scs.SessionManager {
	var persist, secure bool

	//how long should session last?
	minutes, err := strconv.Atoi(s.CookieLifeTime)

	if err != nil {
		minutes = 60
	}

	//should cookies persist
	if strings.ToLower(s.CookiePersist) == "true" {
		persist = true
	}

	//must cookies secure
	if strings.ToLower(s.CookieSecure) == "true" {
		persist = true
	}

	//create session
	session := scs.New()
	session.Lifetime = time.Duration(minutes) * time.Minute
	session.Cookie.Persist = persist
	session.Cookie.Name = s.CookieName
	session.Cookie.Secure = secure
	session.Cookie.Domain = s.CookieDomain
	session.Cookie.SameSite = http.SameSiteLaxMode

	//which session stored
	switch strings.ToLower(s.SessionType) {
	case "redis":

	case "mysql", "mariadb":

	case "postgres", "postgresql":

	default:
		//cookies
	}

	return session
}
