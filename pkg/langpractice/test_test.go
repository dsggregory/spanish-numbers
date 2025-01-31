package langpractice

import (
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPlayResponse(t *testing.T) {
	Convey("Test play", t, func() {
		c := &LangPractice{PlayTimeout: time.Second * 3}
		fp, err := os.Open("../../testdata/resp.json")
		So(err, ShouldBeNil)
		defer fp.Close()

		res, err := c.parseResponse(fp)
		So(err, ShouldBeNil)
		So(res.Number, ShouldEqual, 380)
		err = c.PlayResponse(res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
