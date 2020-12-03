package icons

const (
	File = ""
	Dir  = ""

	LinkFile  = ""
	LinkDir   = ""
	LinkArrow = ""

	Archive = ""
	Audio   = "ﱘ"
	Image   = ""
	Video   = ""

	CLang   = ""
	Clojure = ""
	CPP     = ""
	CSharp  = ""
	Python  = ""
	Shell   = ""
	Subl    = ""
	Win     = ""
	Word    = ""
	Pdf     = ""
	Excel   = ""
	Html    = ""
	Log     = ""
	Jar     = ""
	Xml     = ""
)

// To add new icons just add a new key: value pair here
var Extensions = map[string]string{
	".apk":              "",
	".c":                CLang,
	".h":                CLang,
	".hpp":              CPP,
	".hxx":              CPP,
	".cfg":              "",
	".clj":              Clojure,
	".cljc":             Clojure,
	".cljs":             Clojure,
	".coffee":           "",
	".cc":               CPP,
	".cp":               CPP,
	".cpp":              CPP,
	".cxx":              CPP,
	".cs":               CSharp,
	".csproj":           CSharp,
	".csx":              CSharp,
	".css":              "",
	".d":                "",
	".dart":             "",
	".db":               "",
	".ds_store":         "",
	".go":               "ﳑ",
	".ipynb":            Python,
	".md":               "",
	".py":               Python,
	".pyc":              Python,
	".psd":              "",
	".rs":               "",
	".vue":              "﵂",
	".sln":              "",
	".sql":              "",
	".sqlite3":          "",
	".sublime_keymap":   Subl,
	".sublime_package":  Subl,
	".sublime_settings": Subl,
	".sublime_theme":    Subl,
	".txt":              "",
	".ps1":              Shell,
	".sh":               Shell,
	".shell":            Shell,
	".bat":              Win,
	".exe":              Win,
	".msi":              Win,

	".7z":   Archive,
	".a":    Archive,
	".ar":   Archive,
	".bz2":  Archive,
	".cab":  Archive,
	".cpio": Archive,
	".deb":  Archive,
	".dmg":  Archive,
	".egg":  Archive,
	".gz":   Archive,
	".iso":  Archive,
	".lha":  Archive,
	".mar":  Archive,
	".pak":  Archive,
	".pea":  Archive,
	".rar":  Archive,
	".rpm":  Archive,
	".s7z":  Archive,
	".shar": Archive,
	".tar":  Archive,
	".tbz2": Archive,
	".tgz":  Archive,
	".tlz":  Archive,
	".war":  Archive,
	".whl":  Archive,
	".xpi":  Archive,
	".xz":   Archive,
	".zip":  Archive,
	".zipx": Archive,

	".3dm":  Image,
	".3ds":  Image,
	".ai":   Image,
	".bmp":  Image,
	".dds":  Image,
	".dwg":  Image,
	".dxf":  Image,
	".eps":  Image,
	".gif":  Image,
	".gpx":  Image,
	".jpeg": Image,
	".jpg":  Image,
	".kml":  Image,
	".kmz":  Image,
	".max":  Image,
	".png":  Image,
	".ps":   Image,
	".svg":  Image,
	".tga":  Image,
	".thm":  Image,
	".tif":  Image,
	".tiff": Image,
	".webp": Image,
	".xcf":  Image,
	".icns": Image,

	".aac":  Audio,
	".aiff": Audio,
	".ape":  Audio,
	".au":   Audio,
	".flac": Audio,
	".gsm":  Audio,
	".it":   Audio,
	".m3u":  Audio,
	".m4a":  Audio,
	".mid":  Audio,
	".mp3":  Audio,
	".mpa":  Audio,
	".pls":  Audio,
	".ra":   Audio,
	".s3m":  Audio,
	".sid":  Audio,
	".wav":  Audio,
	".wma":  Audio,
	".xm":   Audio,

	".3g2":   Video,
	".3gp":   Video,
	".aaf":   Video,
	".asf":   Video,
	".avchd": Video,
	".avi":   Video,
	".drc":   Video,
	".flv":   Video,
	".m2v":   Video,
	".m4p":   Video,
	".m4v":   Video,
	".mkv":   Video,
	".mng":   Video,
	".mov":   Video,
	".mp2":   Video,
	".mp4":   Video,
	".mpe":   Video,
	".mpeg":  Video,
	".mpg":   Video,
	".mpv":   Video,
	".mxf":   Video,
	".nsv":   Video,
	".ogg":   Video,
	".ogm":   Video,
	".ogv":   Video,
	".qt":    Video,
	".rm":    Video,
	".rmvb":  Video,
	".roq":   Video,
	".srt":   Video,
	".svi":   Video,
	".vob":   Video,
	".webm":  Video,
	".wmv":   Video,
	".yuv":   Video,

	".jar":  Jar,
	".java": Jar,

	".pdf":  Pdf,
	".docx": Word,
	".doc":  Word,
	".xlsx": Excel,
	".xls":  Excel,
	".csv":  Excel,

	".html": Html,
	".htm":  Html,
	".xml":  Xml,
	".iml":  Xml,

	".log":  Log,
	".json": "",
	".epub": "",
}
