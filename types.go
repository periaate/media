package media

import "strings"

var v = struct{}{}

const (
	Notmedia = iota
	StdImage
	Image
	Video
	Audio
)

var stdSupported = map[string]struct{}{
	".jpg":  v,
	".jpeg": v,
	".jfif": v,
	".png":  v,
	".gif":  v,
}

var imageTypes = map[string]struct{}{
	".jpg":  v,
	".jpeg": v,
	".jfif": v,
	".png":  v,
	".bmp":  v,
	".tiff": v,
	".webp": v,
	".gif":  v,
}

var videoTypes = map[string]struct{}{
	".mp4":  v,
	".webm": v,
	".mkv":  v,
	".avi":  v,
	".mov":  v,
}

var audioTypes = map[string]struct{}{
	".mp3":  v,
	".ogg":  v,
	".flac": v,
	".wav":  v,
	".opus": v,
}

func MediaType(ext string) int {
	ext = strings.ToLower(ext)
	if _, ok := stdSupported[ext]; ok {
		return StdImage
	}
	if _, ok := imageTypes[ext]; ok {
		return Image
	}
	if _, ok := videoTypes[ext]; ok {
		return Video
	}
	if _, ok := audioTypes[ext]; ok {
		return Audio
	}
	return Notmedia
}
