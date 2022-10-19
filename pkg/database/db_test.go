package database

import (
	"github.com/dylan-dinh/api-gin/pkg/model"
	"github.com/dylan-dinh/api-gin/pkg/test"
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
	"testing"
)

func TestDB_InitDb(t *testing.T) {
	db, cf := test.SetUpDbTest()

	Convey("Testing InitDB", t, func() {
		Convey("It should create two table sites and engines", func() {
			tm := db.TableFor(model.Site{})
			So(tm.TableName, ShouldEqual, "sites")

			tm = db.TableFor(model.Engine{})
			So(tm.TableName, ShouldEqual, "engines")
		})

		Convey("If site does not exist it creates it", func() {
			var s model.Site
			So(db.SelectOne(&s, "SELECT * FROM sites WHERE name = $1 AND max_power = $2",
				cf.SiteName, cf.SiteMaxPower), ShouldBeNil)
		})

		Convey("If site name already exist it modifies his power", func() {
			var s model.Site
			So(db.SelectOne(&s, "SELECT * FROM sites WHERE name =$1 AND max_power = $2",
				cf.SiteName, cf.SiteMaxPower), ShouldBeNil)
			So(s.Name, ShouldEqual, cf.SiteName)
			atoi, err := strconv.ParseInt(cf.SiteMaxPower, 10, 64)
			So(err, ShouldBeNil)
			So(s.MaxPower, ShouldEqual, atoi)

			err = db.InitDb()
			So(err, ShouldBeNil)
			cf.SiteMaxPower = "600"
			So(db.SelectOne(&s, "SELECT * FROM sites WHERE name =$1 AND max_power = $2",
				cf.SiteName, cf.SiteMaxPower), ShouldBeNil)
			So(s.MaxPower, ShouldEqual, 600)

		})

		Convey("It should create a trigger of name site_max_power", func() {
			var tgName string
			err := db.SelectOne(&tgName, "SELECT tgname FROM pg_trigger WHERE NOT tgisinternal AND tgrelid = 'engines'::regclass")
			So(err, ShouldBeNil)
			So(tgName, ShouldEqual, "site_max_power")
		})
	})
}
