package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
    "fmt"
	"path/filepath"
	"github.com/webview/webview"
)

var (
	AppDir       string
	windowWidth  = 1214
	windowHeight = 600
	debug        bool
	Version      string
	Build        string
	Commit       string
	DataDir      string
	AppUrl       = "https://abandonware.club/"
	GameList     []GameItem
	settings     Settings
	DEBUG        = 0
	RUNTIMEDEBUG = "no"
)

func SaveSettings() {
	tmpfile := AppDir + string(os.PathSeparator) + "legacybest.json"
	r, err := json.Marshal(settings)
	if err == nil {
		err = ioutil.WriteFile(tmpfile, r, 0644)
		if err == nil {
			WriteLog(true, "Saved program data.. to %s\n", tmpfile)
		}
	}
}

func init() {
        if RUNTIMEDEBUG == "yes" {
        DEBUG = 1
        fmt.Printf("Debug is built-in and enabled!\n")
        }
	var err error
	tmpdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	AppDir = filepath.FromSlash(tmpdir)
	var file, err1 = os.Create(AppDir + string(os.PathSeparator) + "app.log")
	if err1 != nil {
		panic(err1)
	}
	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
	Log.Println("LogFile : " + AppDir + string(os.PathSeparator) + "app.log")

	DataDir = AppDir + string(os.PathSeparator) + "data"
	err = os.MkdirAll(DataDir, os.ModePerm)
	if err != nil {
		WriteLog(true, "Unable to create data directory: %s\n", err)
	}
	settingsfile := AppDir + string(os.PathSeparator) + "legacybest.json"
	if fileExists(settingsfile) {
		jsonFile, err := os.Open(settingsfile)
		// if we os.Open returns an error then handle it
		if err != nil {
			WriteLog(true, "Unable to open system settings: %s\n", err)
			return
		}
		WriteLog(true, "Successfully Opened program settings snapshot data\n")
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)
		err = json.Unmarshal(byteValue, &settings)
		if err == nil {
			WriteLog(true, "Program settings load success!\n")
		}
	}
}

func main() {
	debug := false
	if DEBUG > 0 {
		debug = true
	}
	LoadGameList()
	go app()
	// create a web view
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("LegacyBest")
	w.SetSize(windowWidth, windowHeight, webview.HintNone)
	w.Navigate("http://127.0.0.1:12346/public/html/index.html")
	w.Run()

}
