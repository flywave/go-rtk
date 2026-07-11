package rtk

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/rwcarlsen/goexif/tiff"
	"trimmer.io/go-xmp/xmp"
)

const (
	xmpPacketMarker = "<?xpacket"
)

func init() {
	exif.RegisterParsers(mknote.All...)
}

func ReadExifXMP(reader io.Reader) (error, map[string]interface{}) {
	if reader == nil {
		return fmt.Errorf("reader is nil"), nil
	}
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return err, nil
	}

	x, err := exif.Decode(bytes.NewReader(body))
	if err != nil {
		return err, nil
	}
	var w = exifWalker{fields: map[string]interface{}{}}
	if err = x.Walk(w); err != nil {
		return err, nil
	}

	if n := bytes.Count(body, []byte(xmpPacketMarker)); n < 2 {
		return fmt.Errorf("error while finding XMP document: found %d xpacket markers, expected at least 2", n), nil
	}
	var xmpIndex = bytes.Index(body, []byte(xmpPacketMarker))

	var d = xmp.NewDocument()
	d.SetDirty()

	if err := xmp.Unmarshal(body[xmpIndex:], d); err != nil {
		return err, nil
	}
	paths, err := d.ListPaths()
	if err != nil {
		return err, nil
	}

	for _, p := range paths {
		w.fields[string(p.Path)] = p.Value
	}
	return nil, w.fields
}

type exifWalker struct {
	fields map[string]interface{}
}

func (w exifWalker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	if tag == nil {
		return nil
	}
	value := tag.String()

	if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
		value = strings.TrimRight(strings.TrimLeft(value, `"`), `"`)
	}

	w.fields[string(name)] = value

	return nil
}
