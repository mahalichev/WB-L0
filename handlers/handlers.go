package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/mahalichev/WB-L0/config"
)

func GetRoot(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "")
}

func GetOrderHTML(app *config.App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		order_uid := r.URL.Query().Get("id")
		order, ok := app.Cache[order_uid]
		if !ok {
			app.ErrLog.Printf("order with order_uid %s not found", order_uid)
			http.NotFound(w, r)
			return
		}
		data, err := json.MarshalIndent(order, "", "    ")
		if err != nil {
			app.ErrLog.Print(err)
			http.NotFound(w, r)
			return
		}
		renderTemplate(w, string(data))
	}
}

func GetOrderJSON(app *config.App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		order_uid := r.URL.Query().Get("id")
		order, ok := app.Cache[order_uid]
		if !ok {
			app.ErrLog.Printf("Order with order_uid %s not found", order_uid)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(order)
		if err != nil {
			app.ErrLog.Print(err)
			http.NotFound(w, r)
			return
		}
		w.Write(data)
	}
}

func renderTemplate(w http.ResponseWriter, dataStr string) {
	tmpl, err := template.ParseFiles("./static/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct{ JSONData string }{dataStr}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
