package util

import (
	"fmt"
	"io"

	"github.com/rwcarlsen/goexif/exif"
)

// GetExifData decode exif metadata from img raw data into a map[string]string
// The error appears when the exif library meet error during decoding
func GetExifData(r io.Reader) (map[string]string, error) {
	// Decode the file raw data
	x, err := exif.Decode(r)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]string)
	emptyString := ""

	// Get date/time
	tm, err := x.DateTime()
	if err == nil {
		ret[string(exif.DateTime)] = tm.String()
	} else {
		ret[string(exif.DateTime)] = emptyString
	}
	// Get GPS info
	lat, long, _ := x.LatLong()
	if err == nil {
		ret[string("Location")] = fmt.Sprintf("(%v N,%v E)", lat, long)
	} else {
		ret[string("Location")] = emptyString
	}

	// Get Author info
	artist := getStringField(exif.Artist, x)
	ret[string(exif.Artist)] = artist

	// Get Device Model used
	camModel := getStringField(exif.Model, x)
	ret[string(exif.Model)] = camModel

	// Get Device Model Maker used
	maker := getStringField(exif.Make, x)
	ret[string(exif.Make)] = maker

	// Get Image Description
	imageDescription := getStringField(exif.ImageDescription, x)
	ret[string(exif.ImageDescription)] = imageDescription

	// Get Software used to take the photo
	software := getStringField(exif.Software, x)
	ret[string(exif.Software)] = software

	return ret, nil
}

// getStringField helps to get string val from decoded exif object
func getStringField(field exif.FieldName, ex *exif.Exif) string {
	tag, err := ex.Get(field)
	if err == nil {
		str, err2 := tag.StringVal()
		if err2 == nil {
			return str
		}
	}
	return ""
}
