package media

var v = struct{}{}

// set implementation
var mediaTypes = map[string]struct{}{
	// Images
	// Jpeg types
	".jpg":  v,
	".jpeg": v,
	".jfif": v,

	".png":  v,
	".bmp":  v,
	".tiff": v,
	".webp": v,
	".gif":  v,

	// Video types
	".mp4":  v,
	".webm": v,
	".mkv":  v,
	".avi":  v,
	".mov":  v,

	// Audio types
	".mp3":  v,
	".ogg":  v,
	".flac": v,
	".wav":  v,
	".opus": v,
}

func IsMedia(ext string) bool {
	_, ok := mediaTypes[ext]
	return ok
}
