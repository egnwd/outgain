package controller

import (
	"html/template"
	"log"
	"net/http"

	"github.com/egnwd/outgain/server/achievements"
)

func GetAchievements(staticDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, err := template.New("achievements.tmpl").ParseFiles(staticDir + "/achievements.tmpl")
		if err != nil {
			log.Println(err.Error())
			return
		}

		// Get the username of the authenicated user
		username, err := GetUserName(r)
		if err != nil {
			log.Println(err.Error())
			return
		}

		achmnts := achievements.GetUserAchievements(username)
		if err = t.Execute(w, achmnts); err != nil {
			log.Println(err.Error())
			return
		}
	})
}

func Achievements(staticDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, staticDir+"/achievements.html")
	})
}
