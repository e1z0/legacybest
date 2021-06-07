package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly"
	"github.com/jmoiron/sqlx"
)

var (
	base_url            = "https://www.myabandonware.com"
	crawl_url           = "https://www.myabandonware.com/browse/name/"
	mysqldb             *sqlx.DB
	MySQL_USER          = "test"
	MySQL_PASS          = "test"
	MySQL_HOST          = "localhost"
	MySQL_PORT          = "3306"
	MySQL_DBAH          = "test"
	user_agent          = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:72.0; Minver18) Gecko/20100101 Firefox/70.0 p4r|{4G3NcY"
	expiration          = time.Now().Add(365 * 24 * time.Hour)
	cookie              = http.Cookie{Name: "PHPSESSID", Value: "f2581197f062f049b90461834c9f86be", Expires: expiration}
	GamesAlreadyScraped []string
)

type Game struct {
	Name        string
	Year        string
	Platform    string
	Genre       string
	Tag         string
	Publisher   string
	Developer   string
	Description string
	Rating      string
	Filename    string
	ImgUrl      []string
	OrigUrl     string
}

func StripFileName(fname string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(fname, "")
}

func dbinit() {

	if mysqldb == nil {
		db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", MySQL_USER, MySQL_PASS, MySQL_HOST, MySQL_PORT, MySQL_DBAH))
		if err != nil {
			fmt.Printf("Got error then tried to contact to mysql server %s\n", err)
			return
		}
		mysqldb = db
		return
	}

	err := mysqldb.Ping()
	if err != nil {
		mysqldb = nil
		fmt.Printf("Unable to ping the MySQL connection, it's lost!! Reconnecting...\n")
		db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", MySQL_USER, MySQL_PASS, MySQL_HOST, MySQL_PORT, MySQL_DBAH))
		if err != nil {
			fmt.Printf("Got error then tried to contact to mysql server %s\n", err)
			return
		}
		fmt.Printf("Reconnection was successfull!\n")
		mysqldb = db
	}
	return
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func DownloadFile(filepath string, url string, referer string) error {
	if fileExists(filepath) {
		return nil
	}

	// Get the data
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", user_agent)
	req.Header.Set("Referer", referer)
	req.AddCookie(&cookie)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
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
	if !strings.Contains(url, "screenshots") {
		fmt.Printf("%s -> Download ok\n", url)
	}
	return err
}

func Cleardb() {
	fmt.Printf("Clearing the database!\n")
	_, err := mysqldb.Exec("DELETE FROM uploads_images")
	if err != nil {
		fmt.Printf("Unable to delete the uploads_images table! %s\n", err)
		return
	}
	_, err = mysqldb.Exec("DELETE FROM uploads")
	if err != nil {
		fmt.Printf("Unable to delete the uploads table! %s\n", err)
		return
	}
	fmt.Printf("All tables now are clean!\n")
}

func WriteError(eruoras string) {
	f, err := os.OpenFile("error.log", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = fmt.Fprintln(f, eruoras)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func CurlVisit(url string, save string) {
	if url != "" {

		out, err := os.Create(save)
		if err != nil {
			fmt.Printf("Unable to create file: %s err: %s\n", save, err)
			return
		}
		defer out.Close()

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("error while getting url %s err: %s\n", url, err)
			// handle err
		}
		req.Header.Set("Host", "www.myabandonware.com")
		req.Header.Set("Referer", "https://www.myabandonware.com/game/dungeon-master-n0")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.11; rv:78.0) Gecko/20100101 Firefox/78.0")

		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "PHPSESSID", Value: "b6ac5a1e9bce1b0d27c655b85fefe3d3", Expires: expiration}
		cookie2 := http.Cookie{Name: "__eu", Value: "all", Expires: expiration}
		req.AddCookie(&cookie)
		req.AddCookie(&cookie2)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// handle err
			fmt.Printf("got error: %s", err)
			return
		}
		defer resp.Body.Close()
		fmt.Printf("test\n")
		// Check server response
		if resp.StatusCode != http.StatusOK {
			fmt.Errorf("bad status: %s", resp.Status)
			return
		}
		fmt.Printf("test2\n")
		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		// if resp.StatusCode == http.StatusOK {
		// 	bodyBytes, err := ioutil.ReadAll(resp.Body)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	out.Write(bodyBytes)
		// }

		fmt.Printf("test3\n")
		if err != nil {
			fmt.Printf("error downloading file: %s\n", err)
			return
		}
		fmt.Printf("File downloaded to: %s\n", save)

	}
}

func IfExistsInArr(arr []string, key string) bool {
	for _, x := range arr {
		if x == key {
			return true
		}
	}
	return false
}

func CheckIfExists(game_name string) bool {
	// load all game items to array
	if (len(GamesAlreadyScraped)) == 0 {
		dbinit()
		query := fmt.Sprintf("SELECT name from uploads")
		type itemas struct {
			Name string `db:"name"`
		}
		Rows := []itemas{}
		err := mysqldb.Select(&Rows, query)
		if err != nil {
			fmt.Printf("Unable to load gamelist from sql: %s\n", err)
		}
		for _, row := range Rows {
			GamesAlreadyScraped = append(GamesAlreadyScraped, row.Name)
		}

	}
	ats := IfExistsInArr(GamesAlreadyScraped, game_name)
	//fmt.Printf("Exist state of %s is: %v\n", game_name, ats)
	return ats
}

func main() {
	fmt.Printf("This scraper is designed especially for https://www.myabandonware.com !\n")
	file, err := os.Create("dump.csv")
	if err != nil {
		fmt.Printf("Unable to create file: %s\n", err)
		return
	}
	defer file.Close()
	write := csv.NewWriter(file)
	defer write.Flush()
	dbinit()
	clear := flag.Bool("clear", false, "Clears the scraper database")
	var single string
	flag.StringVar(&single, "single", "", "Specify single scrap url for only one game!")
	//var exists string
	//flag.StringVar(&exists, "exists", "", "Specify game name to check if it's already scraped")
	flag.Parse()
	if *clear {
		Cleardb()
		os.Exit(0)
	}
	// if exists != "" {
	// 	CheckIfExists(exists)
	// 	os.Exit(0)
	// }

	c := colly.NewCollector(colly.MaxBodySize(100 * 1024 * 1024)) // 100mb
	//c.AllowURLRevisit = true
	c.IgnoreRobotsTxt = true
	c.UserAgent = user_agent
	c.SetRequestTimeout(10 * time.Second)
	//c.CacheDir = "./cache"
	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   300 * time.Second, // timeout
			KeepAlive: 300 * time.Second, // keepAlive timeout
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,               // Maximum number of idle connections
		IdleConnTimeout:       90 * time.Second,  // Idle connection timeout
		TLSHandshakeTimeout:   100 * time.Second, // TLS handshake timeout
		ExpectContinueTimeout: 100 * time.Second,
	})

	//extensions.Referer(c)

	// Find and visit all links
	c.OnHTML(".name", func(e *colly.HTMLElement) {
		//url := e.ChildAttr("a", "href")
		url := e.Attr("href")
		fmt.Printf("Got game name: %s link: %s\n", e.Text, url)

		if single == "" {
			e.Request.Visit(url)
		}
	})

	// c.OnRequest(func(r *colly.Request) {
	// 	r.Headers.Set("Host", "www.myabandonware.com")
	// 	r.Headers.Set("Connection", "keep-alive")
	// 	r.Headers.Set("Accept", "*/*")
	// 	r.Headers.Set("Origin", "https://www.myabandonware.com/")
	// 	r.Headers.Set("Referer", "https://www.myabandonware.com/media/")
	// 	r.Headers.Set("Accept-Encoding", "gzip, deflate")
	// 	r.Headers.Set("Accept-Language", "en-US;q=0.9")
	// })

	// c.OnHTML("div.buttons", func(e *colly.HTMLElement) {
	// 	fmt.Printf("radau buttonus\n")
	// 	ch := e.DOM.Find("a.download")
	// 	Urlas, _ := ch.Eq(0).Attr("href")
	// 	if strings.Contains(Urlas, "download") {
	// 		fmt.Printf("Download walking ... %s\n", e.Request.AbsoluteURL(Urlas))
	// 		c.Visit(e.Request.AbsoluteURL(Urlas))
	// 	}
	// })

	// c.OnResponse(func(r *colly.Response) {
	// 	fmt.Printf("%s -> Response code: %d\n", r.Request.URL, r.StatusCode)
	// 	fmt.Printf("body: %s\n", r.Body)
	// 	//r.Save(r.FileName())
	// 	if strings.Index(r.Headers.Get("Content-Type"), "application") > -1 {
	// 		fmt.Printf("zip detected WOOOOOOOOOOOOOOOOOO!\n")

	// 		//return
	// 	}

	// })

	c.OnHTML("#content", func(e *colly.HTMLElement) {
		// e.ForEach("div", func(_ int, elem *colly.HTMLElement) {
		// 	fmt.Printf("name %s\n", elem.ChildText("h2"))
		// })

		ch := e.DOM.Find("div.box").ChildrenFiltered("h2")
		name := ch.Eq(0).Text()
		if name != "" {
			Geimas := Game{}
			Geimas.Name = name
			e.ForEach("table.gameInfo tr", func(_ int, el *colly.HTMLElement) {
				if strings.Contains(el.ChildText("th"), "Year") {
					Geimas.Year = el.ChildText("td")
				}
				if strings.Contains(el.ChildText("th"), "Platform") {
					Geimas.Platform = el.ChildText("td")
				}
				if strings.Contains(el.ChildText("th"), "Genre") {
					Geimas.Genre = el.ChildText("td")
				}
				if strings.Contains(el.ChildText("th"), "Theme") {
					Geimas.Tag = el.ChildText("td")
				}
				if strings.Contains(el.ChildText("th"), "Publisher") {
					Geimas.Publisher = el.ChildText("td")
				}
				if strings.Contains(el.ChildText("th"), "Developer") {
					Geimas.Developer = el.ChildText("td")
				}

			})
			Rating := e.DOM.Find("div.gameRated").ChildrenFiltered("span")
			Geimas.Rating = Rating.Eq(0).Text()

			Description := e.DOM.Find("div.gameDescription").Text()
			Geimas.Description = strings.TrimSpace(Description)
			if Geimas.Description == "" {
				// can be improved somehow, don't know yet
				Geimas.Description = e.DOM.Find("div.box").ChildrenFiltered("p").Eq(2).Text()
			}

			if CheckIfExists(Geimas.Name) == false {

				if strings.Contains(Geimas.Platform, "DOS") {
					fmt.Printf("Fetching game %s\n", Geimas.Name)
					// images
					e.ForEach("div.screens div", func(_ int, el *colly.HTMLElement) {
						platforma, _ := el.DOM.Parent().Attr("data-platform")
						item, _ := el.DOM.Find("div.thumb").ChildrenFiltered("a.lb").Attr("href")
						if platforma == "1" && item != "" {
							bname := path.Base(item)
							//fmt.Printf("Downloading image: %s\n", base_url+item)
							DownloadFile("images/"+bname, base_url+item, e.Request.URL.String())
							Geimas.ImgUrl = append(Geimas.ImgUrl, bname)
						}
					})

					Download := e.DOM.Find("div.buttons").ChildrenFiltered("a.download")
					//Geimas.Filename, _ = Download.Eq(0).Attr("href")
					Geimas.OrigUrl = e.Request.URL.String()
					urlas, _ := Download.Eq(0).Attr("href")
					//fmt.Printf("Gavom download url: %s\n", urlas)
					fakefile := "data/" + StripFileName(Geimas.Name) + ".zip"
					Geimas.Filename = StripFileName(Geimas.Name) + ".zip"
					var downerr error
					if urlas == "" {
						fmt.Printf("This game must be purchashed!\n")
						Geimas.Filename = "buy"
						WriteError(fmt.Sprintf("%s game must be purchashed at %s\n", Geimas.Name, e.Request.URL.String()))
						downerr = nil
					} else {
						downerr = DownloadFile(fakefile, e.Request.AbsoluteURL(urlas), e.Request.URL.String())
					}

					//query := fmt.Sprintf("INSERT INTO uploads (name,year,platform,genre,tag,publisher,developer,description,rating,filename,origurl) VALUES('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s')", Geimas.Name, Geimas.Year, Geimas.Platform, Geimas.Genre, Geimas.Tag, Geimas.Publisher, Geimas.Developer, Geimas.Description, Geimas.Rating, Geimas.Filename, Geimas.OrigUrl)
					//res, erras := mysqldb.Exec(query)
					if downerr == nil {
						res, erras := mysqldb.NamedExec(`INSERT INTO uploads (name,year,platform,genre,tag,publisher,developer,description,rating,filename,origurl) VALUES(:name,:year,:platform,:genre,:tag,:publisher,:developer,:description,:rating,:filename,:origurl)`,
							map[string]interface{}{
								"name":        Geimas.Name,
								"year":        Geimas.Year,
								"platform":    Geimas.Platform,
								"genre":       Geimas.Genre,
								"tag":         Geimas.Tag,
								"publisher":   Geimas.Publisher,
								"developer":   Geimas.Developer,
								"description": Geimas.Description,
								"rating":      Geimas.Rating,
								"filename":    Geimas.Filename,
								"origurl":     Geimas.OrigUrl,
							})
						if erras == nil {
							idas, _ := res.LastInsertId()
							for _, val := range Geimas.ImgUrl {
								sqlimg := fmt.Sprintf("INSERT INTO uploads_images (upload_id,url) VALUES('%d','%s')", idas, val)
								_, err := mysqldb.Exec(sqlimg)
								if err != nil {
									fmt.Printf("Error inserting sql record with image: %s err: %s\n", val, err)
								}
							}
						} else {
							fmt.Printf("Unable to add record to mysql: %s\n", erras)
						}
					} else {
						WriteError(fmt.Sprintf("Unable to download game: %s file: %s, game can be faund on: %s", Geimas.Name, urlas, Geimas.OrigUrl))
					}
				} else {
					fmt.Printf("%s is not DOS game, skipping...\n", Geimas.Name)
				}
			} else {
				fmt.Printf("Game %s already exists on the database!\n", Geimas.Name)
			}
		}

	})

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL)
	// })
	if single != "" {
		c.Visit(single)
	} else {
		c.Visit(crawl_url)
	}
	//c.Visit("https://www.myabandonware.com/game/nam-1965-1975-2x7")
}
