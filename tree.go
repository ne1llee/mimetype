package mimetype

import (
	"sync"

	"github.com/gabriel-vasile/mimetype/internal/magic"
)

// mimetype stores the list of MIME types in a tree structure with
// "application/octet-stream" at the root of the hierarchy. The hierarchy
// approach minimizes the number of checks that need to be done on the input
// and allows for more precise results once the base type of file has been
// identified.
//
// root is a detector which passes for any slice of bytes.
// When a detector passes the check, the children detectors
// are tried in order to find a more accurate MIME type.
var root = newMIME("application/octet-stream", "",
	func([]byte, uint32) bool { return true },
	xpm, sevenZ, zip, pdf, fdf, ole, ps, psd, p7s, ogg, png, jpg, jxl, jp2, jpx,
	jpm, gif, webp, exe, elf, ar, tar, xar, bz2, fits, tiff, bmp, ico, mp3, flac,
	midi, ape, musePack, amr, wav, aiff, au, mpeg, quickTime, mqv, mp4, webM,
	threeGP, threeG2, avi, flv, mkv, asf, aac, voc, aMp4, m4a, m3u, m4v, rmvb,
	gzip, class, swf, crx, ttf, woff, woff2, otf, eot, wasm, shx, dbf, dcm, rar,
	djvu, mobi, lit, bpg, sqlite3, dwg, nes, lnk, macho, qcp, icns, heic,
	heicSeq, heif, heifSeq, hdr, mrc, mdb, accdb, zstd, cab, rpm, xz, lzip,
	torrent, cpio, tzif, xcf, pat, gbr, glb,
	// Keep text last because it is the slowest check
	text,
)

// errMIME is returned from Detect functions when err is not nil.
// Detect can return root for erroneous cases, but it needs to lock mu in order to do so.
// errMIME is same as root but it does not require locking.
var errMIME = newMIME("application/octet-stream", "", func([]byte, uint32) bool { return false })

// mu guards access to the root MIME tree. Access to root must be synchonized with this lock.
var mu = &sync.RWMutex{}

// The list of nodes appended to the root node.
var (
	xz   = newMIME("application/x-xz", ".xz", magic.Xz)
	gzip = newMIME("application/gzip", ".gz", magic.Gzip).alias(
		"application/x-gzip", "application/x-gunzip", "application/gzipped",
		"application/gzip-compressed", "application/x-gzip-compressed",
		"gzip/document")
	sevenZ = newMIME("application/x-7z-compressed", ".7z", magic.SevenZ)
	zip    = newMIME("application/zip", ".zip", magic.Zip, xlsx, docx, pptx, epub, jar, odt, ods, odp, odg, odf, odc, sxc).
		alias("application/x-zip", "application/x-zip-compressed")
	tar = newMIME("application/x-tar", ".tar", magic.Tar)
	xar = newMIME("application/x-xar", ".xar", magic.Xar)
	bz2 = newMIME("application/x-bzip2", ".bz2", magic.Bz2)
	pdf = newMIME("application/pdf", ".pdf", magic.Pdf).
		alias("application/x-pdf")
	fdf  = newMIME("application/vnd.fdf", ".fdf", magic.Fdf)
	xlsx = newMIME("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", ".xlsx", magic.Xlsx)
	docx = newMIME("application/vnd.openxmlformats-officedocument.wordprocessingml.document", ".docx", magic.Docx)
	pptx = newMIME("application/vnd.openxmlformats-officedocument.presentationml.presentation", ".pptx", magic.Pptx)
	epub = newMIME("application/epub+zip", ".epub", magic.Epub)
	jar  = newMIME("application/jar", ".jar", magic.Jar)
	ole  = newMIME("application/x-ole-storage", "", magic.Ole, aaf, msg, xls, pub, ppt, doc)
	aaf  = newMIME("application/octet-stream", ".aaf", magic.Aaf)
	doc  = newMIME("application/msword", ".doc", magic.Doc).
		alias("application/vnd.ms-word")
	ppt = newMIME("application/vnd.ms-powerpoint", ".ppt", magic.Ppt).
		alias("application/mspowerpoint")
	pub = newMIME("application/vnd.ms-publisher", ".pub", magic.Pub)
	xls = newMIME("application/vnd.ms-excel", ".xls", magic.Xls).
		alias("application/msexcel")
	msg  = newMIME("application/vnd.ms-outlook", ".msg", magic.Msg)
	ps   = newMIME("application/postscript", ".ps", magic.Ps)
	fits = newMIME("application/fits", ".fits", magic.Fits)
	ogg  = newMIME("application/ogg", ".ogg", magic.Ogg, oggAudio, oggVideo).
		alias("application/x-ogg")
	oggAudio = newMIME("audio/ogg", ".oga", magic.OggAudio)
	oggVideo = newMIME("video/ogg", ".ogv", magic.OggVideo)
	text     = newMIME("text/plain", ".txt", magic.Text, html, svg, xml, php, js, lua, perl, python, json, ndJson, rtf, tcl, csv, tsv, vCard, iCalendar, warc)
	xml      = newMIME("text/xml", ".xml", magic.Xml, rss, atom, x3d, kml, xliff, collada, gml, gpx, tcx, amf, threemf, xfdf, owl2)
	json     = newMIME("application/json", ".json", magic.Json, geoJson, har)
	har      = newMIME("application/json", ".har", magic.HAR)
	csv      = newMIME("text/csv", ".csv", magic.Csv)
	tsv      = newMIME("text/tab-separated-values", ".tsv", magic.Tsv)
	geoJson  = newMIME("application/geo+json", ".geojson", magic.GeoJson)
	ndJson   = newMIME("application/x-ndjson", ".ndjson", magic.NdJson)
	html     = newMIME("text/html", ".html", magic.Html)
	php      = newMIME("text/x-php", ".php", magic.Php)
	rtf      = newMIME("text/rtf", ".rtf", magic.Rtf)
	js       = newMIME("application/javascript", ".js", magic.Js).
			alias("application/x-javascript", "text/javascript")
	lua    = newMIME("text/x-lua", ".lua", magic.Lua)
	perl   = newMIME("text/x-perl", ".pl", magic.Perl)
	python = newMIME("application/x-python", ".py", magic.Python)
	tcl    = newMIME("text/x-tcl", ".tcl", magic.Tcl).
		alias("application/x-tcl")
	vCard     = newMIME("text/vcard", ".vcf", magic.VCard)
	iCalendar = newMIME("text/calendar", ".ics", magic.ICalendar)
	svg       = newMIME("image/svg+xml", ".svg", magic.Svg)
	rss       = newMIME("application/rss+xml", ".rss", magic.Rss).
			alias("text/rss")
	owl2    = newMIME("application/owl+xml", ".owl", magic.Owl2)
	atom    = newMIME("application/atom+xml", ".atom", magic.Atom)
	x3d     = newMIME("model/x3d+xml", ".x3d", magic.X3d)
	kml     = newMIME("application/vnd.google-earth.kml+xml", ".kml", magic.Kml)
	xliff   = newMIME("application/x-xliff+xml", ".xlf", magic.Xliff)
	collada = newMIME("model/vnd.collada+xml", ".dae", magic.Collada)
	gml     = newMIME("application/gml+xml", ".gml", magic.Gml)
	gpx     = newMIME("application/gpx+xml", ".gpx", magic.Gpx)
	tcx     = newMIME("application/vnd.garmin.tcx+xml", ".tcx", magic.Tcx)
	amf     = newMIME("application/x-amf", ".amf", magic.Amf)
	threemf = newMIME("application/vnd.ms-package.3dmanufacturing-3dmodel+xml", ".3mf", magic.Threemf)
	png     = newMIME("image/png", ".png", magic.Png)
	jpg     = newMIME("image/jpeg", ".jpg", magic.Jpg)
	jxl     = newMIME("image/jxl", ".jxl", magic.Jxl)
	jp2     = newMIME("image/jp2", ".jp2", magic.Jp2)
	jpx     = newMIME("image/jpx", ".jpf", magic.Jpx)
	jpm     = newMIME("image/jpm", ".jpm", magic.Jpm).
		alias("video/jpm")
	xpm  = newMIME("image/x-xpixmap", ".xpm", magic.Xpm)
	bpg  = newMIME("image/bpg", ".bpg", magic.Bpg)
	gif  = newMIME("image/gif", ".gif", magic.Gif)
	webp = newMIME("image/webp", ".webp", magic.Webp)
	tiff = newMIME("image/tiff", ".tiff", magic.Tiff)
	bmp  = newMIME("image/bmp", ".bmp", magic.Bmp).
		alias("image/x-bmp", "image/x-ms-bmp")
	ico  = newMIME("image/x-icon", ".ico", magic.Ico)
	icns = newMIME("image/x-icns", ".icns", magic.Icns)
	psd  = newMIME("image/vnd.adobe.photoshop", ".psd", magic.Psd).
		alias("image/x-psd", "application/photoshop")
	heic    = newMIME("image/heic", ".heic", magic.Heic)
	heicSeq = newMIME("image/heic-sequence", ".heic", magic.HeicSequence)
	heif    = newMIME("image/heif", ".heif", magic.Heif)
	heifSeq = newMIME("image/heif-sequence", ".heif", magic.HeifSequence)
	hdr     = newMIME("image/vnd.radiance", ".hdr", magic.Hdr)
	mp3     = newMIME("audio/mpeg", ".mp3", magic.Mp3).
		alias("audio/x-mpeg", "audio/mp3")
	flac = newMIME("audio/flac", ".flac", magic.Flac)
	midi = newMIME("audio/midi", ".midi", magic.Midi).
		alias("audio/mid", "audio/sp-midi", "audio/x-mid", "audio/x-midi")
	ape      = newMIME("audio/ape", ".ape", magic.Ape)
	musePack = newMIME("audio/musepack", ".mpc", magic.MusePack)
	wav      = newMIME("audio/wav", ".wav", magic.Wav).
			alias("audio/x-wav", "audio/vnd.wave", "audio/wave")
	aiff = newMIME("audio/aiff", ".aiff", magic.Aiff)
	au   = newMIME("audio/basic", ".au", magic.Au)
	amr  = newMIME("audio/amr", ".amr", magic.Amr).
		alias("audio/amr-nb")
	aac  = newMIME("audio/aac", ".aac", magic.Aac)
	voc  = newMIME("audio/x-unknown", ".voc", magic.Voc)
	aMp4 = newMIME("audio/mp4", ".mp4", magic.AMp4).
		alias("audio/x-m4a", "audio/x-mp4a")
	m4a = newMIME("audio/x-m4a", ".m4a", magic.M4a)
	m3u = newMIME("application/vnd.apple.mpegurl", ".m3u", magic.M3u).
		alias("audio/mpegurl")
	m4v  = newMIME("video/x-m4v", ".m4v", magic.M4v)
	mp4  = newMIME("video/mp4", ".mp4", magic.Mp4)
	webM = newMIME("video/webm", ".webm", magic.WebM).
		alias("audio/webm")
	mpeg      = newMIME("video/mpeg", ".mpeg", magic.Mpeg)
	quickTime = newMIME("video/quicktime", ".mov", magic.QuickTime)
	mqv       = newMIME("video/quicktime", ".mqv", magic.Mqv)
	threeGP   = newMIME("video/3gpp", ".3gp", magic.ThreeGP).
			alias("video/3gp", "audio/3gpp")
	threeG2 = newMIME("video/3gpp2", ".3g2", magic.ThreeG2).
		alias("video/3g2", "audio/3gpp2")
	avi = newMIME("video/x-msvideo", ".avi", magic.Avi).
		alias("video/avi", "video/msvideo")
	flv = newMIME("video/x-flv", ".flv", magic.Flv)
	mkv = newMIME("video/x-matroska", ".mkv", magic.Mkv)
	asf = newMIME("video/x-ms-asf", ".asf", magic.Asf).
		alias("video/asf", "video/x-ms-wmv")
	rmvb  = newMIME("application/vnd.rn-realmedia-vbr", ".rmvb", magic.Rmvb)
	class = newMIME("application/x-java-applet", ".class", magic.Class)
	swf   = newMIME("application/x-shockwave-flash", ".swf", magic.Swf)
	crx   = newMIME("application/x-chrome-extension", ".crx", magic.Crx)
	ttf   = newMIME("font/ttf", ".ttf", magic.Ttf).
		alias("font/sfnt", "application/x-font-ttf", "application/font-sfnt")
	woff    = newMIME("font/woff", ".woff", magic.Woff)
	woff2   = newMIME("font/woff2", ".woff2", magic.Woff2)
	otf     = newMIME("font/otf", ".otf", magic.Otf)
	eot     = newMIME("application/vnd.ms-fontobject", ".eot", magic.Eot)
	wasm    = newMIME("application/wasm", ".wasm", magic.Wasm)
	shp     = newMIME("application/octet-stream", ".shp", magic.Shp)
	shx     = newMIME("application/octet-stream", ".shx", magic.Shx, shp)
	dbf     = newMIME("application/x-dbf", ".dbf", magic.Dbf)
	exe     = newMIME("application/vnd.microsoft.portable-executable", ".exe", magic.Exe)
	elf     = newMIME("application/x-elf", "", magic.Elf, elfObj, elfExe, elfLib, elfDump)
	elfObj  = newMIME("application/x-object", "", magic.ElfObj)
	elfExe  = newMIME("application/x-executable", "", magic.ElfExe)
	elfLib  = newMIME("application/x-sharedlib", ".so", magic.ElfLib)
	elfDump = newMIME("application/x-coredump", "", magic.ElfDump)
	ar      = newMIME("application/x-archive", ".a", magic.Ar, deb).
		alias("application/x-unix-archive")
	deb = newMIME("application/vnd.debian.binary-package", ".deb", magic.Deb)
	rpm = newMIME("application/x-rpm", ".rpm", magic.Rpm)
	dcm = newMIME("application/dicom", ".dcm", magic.Dcm)
	odt = newMIME("application/vnd.oasis.opendocument.text", ".odt", magic.Odt, ott).
		alias("application/x-vnd.oasis.opendocument.text")
	ott = newMIME("application/vnd.oasis.opendocument.text-template", ".ott", magic.Ott).
		alias("application/x-vnd.oasis.opendocument.text-template")
	ods = newMIME("application/vnd.oasis.opendocument.spreadsheet", ".ods", magic.Ods, ots).
		alias("application/x-vnd.oasis.opendocument.spreadsheet")
	ots = newMIME("application/vnd.oasis.opendocument.spreadsheet-template", ".ots", magic.Ots).
		alias("application/x-vnd.oasis.opendocument.spreadsheet-template")
	odp = newMIME("application/vnd.oasis.opendocument.presentation", ".odp", magic.Odp, otp).
		alias("application/x-vnd.oasis.opendocument.presentation")
	otp = newMIME("application/vnd.oasis.opendocument.presentation-template", ".otp", magic.Otp).
		alias("application/x-vnd.oasis.opendocument.presentation-template")
	odg = newMIME("application/vnd.oasis.opendocument.graphics", ".odg", magic.Odg, otg).
		alias("application/x-vnd.oasis.opendocument.graphics")
	otg = newMIME("application/vnd.oasis.opendocument.graphics-template", ".otg", magic.Otg).
		alias("application/x-vnd.oasis.opendocument.graphics-template")
	odf = newMIME("application/vnd.oasis.opendocument.formula", ".odf", magic.Odf).
		alias("application/x-vnd.oasis.opendocument.formula")
	odc = newMIME("application/vnd.oasis.opendocument.chart", ".odc", magic.Odc).
		alias("application/x-vnd.oasis.opendocument.chart")
	sxc = newMIME("application/vnd.sun.xml.calc", ".sxc", magic.Sxc)
	rar = newMIME("application/x-rar-compressed", ".rar", magic.Rar).
		alias("application/x-rar")
	djvu    = newMIME("image/vnd.djvu", ".djvu", magic.DjVu)
	mobi    = newMIME("application/x-mobipocket-ebook", ".mobi", magic.Mobi)
	lit     = newMIME("application/x-ms-reader", ".lit", magic.Lit)
	sqlite3 = newMIME("application/x-sqlite3", ".sqlite", magic.Sqlite)
	dwg     = newMIME("image/vnd.dwg", ".dwg", magic.Dwg).
		alias("image/x-dwg", "application/acad", "application/x-acad",
			"application/autocad_dwg", "application/dwg", "application/x-dwg",
			"application/x-autocad", "drawing/dwg")
	warc    = newMIME("application/warc", ".warc", magic.Warc)
	nes     = newMIME("application/vnd.nintendo.snes.rom", ".nes", magic.Nes)
	lnk     = newMIME("application/x-ms-shortcut", ".lnk", magic.Lnk)
	macho   = newMIME("application/x-mach-binary", ".macho", magic.MachO)
	qcp     = newMIME("audio/qcelp", ".qcp", magic.Qcp)
	mrc     = newMIME("application/marc", ".mrc", magic.Marc)
	mdb     = newMIME("application/x-msaccess", ".mdb", magic.MsAccessMdb)
	accdb   = newMIME("application/x-msaccess", ".accdb", magic.MsAccessAce)
	zstd    = newMIME("application/zstd", ".zst", magic.Zstd)
	cab     = newMIME("application/vnd.ms-cab-compressed", ".cab", magic.Cab)
	lzip    = newMIME("application/lzip", ".lz", magic.Lzip)
	torrent = newMIME("application/x-bittorrent", ".torrent", magic.Torrent)
	cpio    = newMIME("application/x-cpio", ".cpio", magic.Cpio)
	tzif    = newMIME("application/tzif", "", magic.TzIf)
	p7s     = newMIME("application/pkcs7-signature", ".p7s", magic.P7s)
	xcf     = newMIME("image/x-xcf", ".xcf", magic.Xcf)
	pat     = newMIME("image/x-gimp-pat", ".pat", magic.Pat)
	gbr     = newMIME("image/x-gimp-gbr", ".gbr", magic.Gbr)
	xfdf    = newMIME("application/vnd.adobe.xfdf", ".xfdf", magic.Xfdf)
	glb     = newMIME("model/gltf-binary", ".glb", magic.Glb)
)
