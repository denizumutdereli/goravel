package gortisan

import "net/http"

func (g *Gortisan) SessionLoad(next http.Handler) http.Handler {
	g.InfoLog.Println("Session load called.")
	return g.Session.LoadAndSave(next)
}
