package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/packr/v2"
)

func app() {
	box := packr.New("public", "../public")
	mux := http.NewServeMux()
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(box)))
	mux.HandleFunc("/home", home)
	mux.HandleFunc("/key", captureKeys)
	mux.HandleFunc("/open", OpenItem)
	mux.HandleFunc("/sideload", Sideload)
	mux.HandleFunc("/fav", FavItem)
	mux.HandleFunc("/gamelist", ListGames)
	mux.HandleFunc("/genreslist", GameGenres)
	mux.HandleFunc("/modeslist", GameModes)
	mux.HandleFunc("/gamedetails", GameDetails)
	server := &http.Server{
		Addr:    "127.0.0.1:12346",
		Handler: mux,
	}
	server.ListenAndServe()
}

// capture keyboard events
func captureKeys(w http.ResponseWriter, r *http.Request) {
	ev := r.FormValue("event")
	// what to react to when the game is over
	if ev == "81" { // q
		os.Exit(0)
	}
	w.Header().Set("Cache-Control", "no-cache")
}

func GameDetails(w http.ResponseWriter, r *http.Request) {
	uid := r.FormValue("uid")
	if DEBUG > 0 {
		WriteLog(true, "Request GameDetails with uid: %s\n", uid)
	}
	geimas, err := FindGame(uid)
	if IsInstalled(uid) {
		geimas.Installed = true
	}
	if IsFav(uid) {
		geimas.Fav = true
	}
	if settings.Experimental {
		geimas.Experimental = true
	}

	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		val, err := json.Marshal(geimas)
		if err == nil {
			if DEBUG > 0 {
				WriteLog(true, "GameDetails Func to frontend: %s\n", val)
			}
			w.Write(val)
		}
	}

}

func GameGenres(w http.ResponseWriter, r *http.Request) {
	url := AppUrl + "api.php?list=genres"
	res, err := http.Get(url)
	if err == nil {
		body, err := ioutil.ReadAll(res.Body)
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if DEBUG > 0 {
				WriteLog(true, "GameGenres Func to frontend: %s\n", body)
			}
			w.Write(body)
		}
	}
}

func GameModes(w http.ResponseWriter, r *http.Request) {
	url := AppUrl + "api.php?list=modes"
	res, err := http.Get(url)
	if err == nil {
		body, err := ioutil.ReadAll(res.Body)
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if DEBUG > 0 {
				WriteLog(true, "GameModes Func to frontend: %s\n", body)
			}
			w.Write(body)
		}
	}
}

// fav item
func FavItem(w http.ResponseWriter, r *http.Request) {
	ev := r.FormValue("item")
	var favas bool

	if IsFav(ev) {
		UnmarkFav(ev)
		favas = false
	} else {
		MarkFav(ev)
		favas = true
	}
	WriteLog(true, "Got fav instruction on: %s returning: %v\n", ev, favas)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// the standard list (without any params)
	type Resp struct {
		Fav bool
	}
	respas := Resp{Fav: favas}
	val, err := json.Marshal(respas)
	if err == nil {
		w.Write(val)
	}

}

// sideload item
func Sideload(w http.ResponseWriter, r *http.Request) {
	ev := r.FormValue("item")
	tipas := 0
	var name string
	for _, val := range GameList {
		if val.Uid == ev {
			name = StripFileName(val.Name)
		}
	}
	WriteLog(true, "Downloading game for sideload...\n")
	tmpfile := DataDir + string(os.PathSeparator) + name + ".zip"
	fileUrl := AppUrl + "/download.php?file=" + ev
	err := DownloadFile(tmpfile, fileUrl)
	if err != nil {
		WriteLog(true, "Unable to download game: %s, err: %s", ev, err)
	}
	err = TransferToRaspberry(Sshopts{Host: "192.168.1.1", User: "pi", Password: "", Port: 22}, tmpfile)
	if err == nil {
		fmt.Printf("Transfer was completed!\n")
		tipas = 1
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// the standard list (without any params)
	type Resp struct {
		Type int
	}
	respas := Resp{Type: tipas}
	val, err := json.Marshal(respas)
	if err == nil {
		w.Write(val)
	}
}

// open item
func OpenItem(w http.ResponseWriter, r *http.Request) {
	ev := r.FormValue("item")
	Installed := false
	tipas := 0
	for _, val := range GameList {
		if val.Uid == ev {
			if IsInstalled(ev) {
				Installed = true
			}
		}
	}
	if Installed {
		tipas = 1
		RunGame(ev)
	} else {
		tipas = 2
		InstallGame(ev)
	}
	WriteLog(true, "Got open instruction on: %s\n", ev)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// the standard list (without any params)
	type Resp struct {
		Type int
	}
	respas := Resp{Type: tipas}
	val, err := json.Marshal(respas)
	if err == nil {
		w.Write(val)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(AppDir + filepath.FromSlash("/public/html/home.html"))
	WriteLog(true, "Loading home page\n")
	// start generating frames in a new goroutine
	t.Execute(w, 1000)
}

func ListGames(w http.ResponseWriter, r *http.Request) {
	mode := r.FormValue("mode")
	key := r.FormValue("key")
        if DEBUG > 0 {
	WriteLog(true, "ListGames issued with params: %s -> %s\n", mode, key)
        }
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// the standard list (without any params)
	if mode == "0" && key == "" {
		val, err := json.Marshal(GameList)
		if err == nil {
                        if DEBUG > 0 {
			WriteLog(true, "2ListGames without params to frontend: %s\n", val)
                        }
			w.Write(val)

		}
	} else if mode == "0" && key != "" {
		gmlist := []GameItem{}
		for _, gm := range GameList {
			if ContainsI(gm.Name, key) {
				gmlist = append(gmlist, gm)
			}
		}
		val, err := json.Marshal(gmlist)
		if err == nil {
                        if DEBUG > 0 {
			WriteLog(true, "ListGames using search return to frontend: %s\n", val)
			}
			w.Write(val)
		}
	} else if mode == "1" && key != "" {
		gmlist, err := FindByGenre(key)
		if err == nil {
			val, err := json.Marshal(gmlist)
			if err == nil {
				WriteLog(true, "ListGames using find games by genre return to frontend: %s\n", val)
				w.Write(val)
			}
		}
	} else if mode == "2" && key != "" {
		gmlist, err := FindByMode(key)
		if err == nil {
			val, err := json.Marshal(gmlist)
			if err == nil {
				WriteLog(true, "ListGames using find games by mode return to frontend: %s\n", val)
				w.Write(val)
			}
		}
	} else if mode == "3" && key == "" {
		gmlist := []GameItem{}
		for _, gm := range GameList {
			if IsFav(gm.Uid) {
				gmlist = append(gmlist, gm)
			}
		}
		val, err := json.Marshal(gmlist)
		if err == nil {
			WriteLog(true, "ListGames using favlist return to frontend: %s\n", val)
			w.Write(val)
		}
	}
}

func ExecuteDosBox(filename string) {
	SaveSettings()
	app := AppDir + string(os.PathSeparator) + "DOSBox"
	params := []string{}
	params = append(params, filename)
	if settings.DosBoxExitOnGame {
		params = append(params, "-noconsoleflag")
		params = append(params, "-exit")
	}
	//cmd := exec.Command(app, arg0, arg1, arg2, arg3)
	cmd := exec.Command(app, params...)
	stdout, err := cmd.Output()

	if err != nil {
		WriteLog(true, err.Error())
		return
	}

	fmt.Printf(string(stdout))
}
