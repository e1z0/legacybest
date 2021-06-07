package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

var (
	Log *log.Logger
)

func WriteLog(stdout bool, text string, a ...interface{}) {
if DEBUG > 0 {
	Log.Printf(text, a...)
	if stdout {
		fmt.Printf(text, a...)
	}
}

}

func LoadGameList() {
	url := AppUrl + "api.php?list=games"
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	json.Unmarshal(body, &GameList)
}

func ContainsI(a string, b string) bool {
	return strings.Contains(
		strings.ToLower(a),
		strings.ToLower(b),
	)
}

func FindByGenre(genre string) ([]GameItem, error) {
	GmLst := []GameItem{}
	url := AppUrl + "api.php?findbygenre=" + genre
	res, err := http.Get(url)
	if err == nil {
		body, err := ioutil.ReadAll(res.Body)
		if err == nil {
			json.Unmarshal(body, &GmLst)
			return GmLst, nil
		}
	}
	return GmLst, err
}

func FindByMode(genre string) ([]GameItem, error) {
	GmLst := []GameItem{}
	url := AppUrl + "api.php?findbymode=" + genre
	res, err := http.Get(url)
	if err == nil {
		body, err := ioutil.ReadAll(res.Body)
		if err == nil {
			json.Unmarshal(body, &GmLst)
			return GmLst, nil
		}
	}
	return GmLst, err
}

func FindGame(uid string) (Game, error) {
	Game := Game{}
	for _, val := range GameList {
		if val.Uid == uid {
			WriteLog(true, "FindGame: Found game: %s\n", uid)
			// fetch local package.json then if not found fetch the remote api
			url := AppUrl + "api.php?gamedetails=" + uid
			res, err := http.Get(url)
			if err == nil {
				body, err := ioutil.ReadAll(res.Body)
				if err == nil {
					json.Unmarshal(body, &Game)
					return Game, nil
				}
			} else {
				WriteLog(true, "Error while getting game... %s\n", err)
			}
		}
	}
	return Game, nil
}

func RunGame(uid string) {
	WriteLog(true, "Running game %s\n", uid)
	for _, val := range GameList {
		if val.Uid == uid {
			if fileExists(DataDir + string(os.PathSeparator) + uid + string(os.PathSeparator) + filepath.FromSlash(val.Exe)) {
				WriteLog(true, "Found installed app in %s\n", DataDir+string(os.PathSeparator)+uid+string(os.PathSeparator)+filepath.FromSlash(val.Exe))
				MarkInstalled(uid)
				ExecuteDosBox(DataDir + string(os.PathSeparator) + uid + string(os.PathSeparator) + filepath.FromSlash(val.Exe))
			} else {
				WriteLog(true, "We are unable to find app exe at: %s the game installation must be broken, reinstalling...\n", DataDir+string(os.PathSeparator)+uid+string(os.PathSeparator)+filepath.FromSlash(val.Exe))
				UnmarkInstalled(uid)
				InstallGame(uid)
			}
			return
		}
	}
}

func IsInstalled(uid string) bool {
	for _, key := range settings.InstallList {
		if key == uid {
			return true
		}
	}
	return false
}

func IsFav(uid string) bool {
	for _, key := range settings.FavList {
		if key == uid {
			return true
		}
	}
	return false
}

func MarkFav(uid string) {
	if IsFav(uid) == false {
		settings.FavList = append(settings.FavList, uid)
	}
}
func UnmarkFav(uid string) {
	if IsFav(uid) {
		for i, key := range settings.FavList {
			fmt.Printf("Element found %s at %d\n", key, i)
			if key == uid {
				copy(settings.FavList[i:], settings.FavList[i+1:])
				settings.FavList[len(settings.FavList)-1] = ""
				settings.FavList = settings.FavList[:len(settings.FavList)-1]
				break
			}
		}
		//fmt.Printf("Found fav element %d to remove from list\n", element)

		fmt.Printf("Got list %v\n", settings.FavList)
		return
	}
}

func MarkInstalled(uid string) {
	if IsInstalled(uid) == false {
		settings.InstallList = append(settings.InstallList, uid)
	}
}

func UnmarkInstalled(uid string) bool {
	element := 0
	for i, key := range settings.InstallList {
		if key == uid {
			element = i
		}
	}
	if element > 0 {
		settings.InstallList[element] = settings.InstallList[len(settings.InstallList)-1]
		settings.InstallList = settings.InstallList[:len(settings.InstallList)-1]
	}
	return false
}

func InstallGame(uid string) {
	WriteLog(true, "Installing game %s\n", uid)
	tmpfile := DataDir + string(os.PathSeparator) + uuid.New().String() + ".zip"
	fileUrl := AppUrl + "/download.php?file=" + uid
	err := DownloadFile(tmpfile, fileUrl)
	if err != nil {
		WriteLog(true, "Unable to download game: %s, err: %s", uid, err)
	}
	if DEBUG > 0 {
		WriteLog(true, "Game: %s downloaded to: %s\n", fileUrl, tmpfile)
	}
	// we should extract it there :D
	files, err := Unzip(tmpfile, DataDir+string(os.PathSeparator)+uid+string(os.PathSeparator))
	if err != nil {
		WriteLog(true, "Unable to unzip archive %s\n", tmpfile)
		return
	}
	if DEBUG > 0 {
		WriteLog(true, "Unzipped: %s\n", strings.Join(files, "\n"))
	}
	if fileExists(tmpfile) {
		_ = os.Remove(tmpfile)
	}
	RunGame(uid)
}

func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func StripFileName(fname string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(fname, "")
}
