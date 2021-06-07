package main

// Game - The main object in the program
type Game struct {
	UID     string `json:"uid"`
	Name    string
	Ratings struct {
		DownloadedTimes int
		ViewedTimes     int
		Likes           int
		Dislikes        int
		EditorsRating   int
	}
	Description string
	Brief       string
	Producer    string
	Publisher   string
	Year        string
	//Platform    string
	Modes  []string // single player, multi ipx
	Genres []string
	Tags   string
	// for later implementation
	// emulation   struct {
	// 	RunArgs         []string
	// 	RunExecutable   string
	// 	SetupExecutable string
	// 	Engine          int // dosbox, scummvm etc...
	// }
	PackageVersion int
	VideoUrls      []string
	ImgUrls        []string
	// CoverPic       string
	// Icon           string
	Installed    bool
	Experimental bool
	Fav          bool
}

type GameItem struct {
	Uid   string `json:"uid"`
	Name  string `json:"name"`
	Year  string `json:"year"`
	Exe   string `json:"exe"`
	Brief string `json:"brief"`
}

// Settings - Structure for the program settings
type Settings struct {
	windowWidth      int
	windowHeight     int
	DosBoxExitOnGame bool
	Experimental     bool
	InstallList      []string
	FavList          []string
}
