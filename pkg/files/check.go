package files

import (
	"crypto/sha1"
	"encoding/base32"
	"net/http"
)

var enc = base32.NewEncoding("qetzk24ycrvb35oswagf7xdjplh6nmiu").WithPadding(base32.NoPadding)

func GenFileName(b []byte) string {
	sum := sha1.Sum(b)
	return enc.EncodeToString(sum[:])
}

func CheckImage(b []byte) (ext string) {
	mt := http.DetectContentType(b)
	if mt == "image/jpeg" || mt == "image/png" || mt == "image/gif" {
		ext = mt[6:]
	}
	return
}

func CheckAudio(b []byte) (ext string) {
	mt := http.DetectContentType(b)
	switch mt {
	case "audio/mpeg":
		ext = "mp3"
	case "audio/wave":
		ext = "wav"
	case "application/ogg":
		ext = "oga"
	}
	return
}

func CheckVideo(b []byte) (ext string) {
	mt := http.DetectContentType(b)
	switch mt {
	case "video/mp4":
		ext = "mp4"
	case "video/webm":
		ext = "webm"
	case "application/ogg":
		ext = "ogv"
	}
	return
}
