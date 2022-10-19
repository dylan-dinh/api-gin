package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dylan-dinh/api-gin/pkg/model"
	"github.com/dylan-dinh/api-gin/pkg/test"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostEngine(t *testing.T) {
	Convey("Testing post handler", t, func() {
		Convey("If body is well formatted", func() {
			db, _ := test.SetUpDbTest()
			ts := httptest.NewServer(NewRouter(db))
			defer ts.Close()
			defer test.ClearDB(db)
			tc := ts.Client()

			site := model.Site{
				Id:       1000,
				Name:     "site1",
				MaxPower: 100,
			}
			So(db.Insert(&site), ShouldBeNil)

			body := &model.Engine{
				SiteId:        site.Id,
				Name:          "foo",
				Type:          "furnace",
				RatedCapacity: 100,
			}

			jsonEncoded, err := json.Marshal(body)
			So(err, ShouldBeNil)

			r, err := http.NewRequest(http.MethodPost,
				ts.URL+"/api/engines", bytes.NewBuffer(jsonEncoded))
			So(err, ShouldBeNil)

			res, err := tc.Do(r)
			So(err, ShouldBeNil)
			defer res.Body.Close()

			So(res.StatusCode, ShouldEqual, http.StatusCreated)
		})

		Convey("If body has many wrong fields", func() {
			db, _ := test.SetUpDbTest()
			ts := httptest.NewServer(NewRouter(db))
			defer ts.Close()
			tc := ts.Client()

			body := &model.Engine{
				SiteId:        10,
				Name:          "foo",
				Type:          "foo",
				RatedCapacity: 0,
			}

			jsonEncoded, err := json.Marshal(body)
			So(err, ShouldBeNil)

			r, err := http.NewRequest(http.MethodPost,
				ts.URL+"/api/engines", bytes.NewBuffer(jsonEncoded))
			So(err, ShouldBeNil)

			res, err := tc.Do(r)
			So(err, ShouldBeNil)
			defer res.Body.Close()

			So(res.StatusCode, ShouldEqual, http.StatusBadRequest)
		})

		Convey("If body has missing fields ", func() {
			db, _ := test.SetUpDbTest()
			ts := httptest.NewServer(NewRouter(db))
			defer ts.Close()
			tc := ts.Client()

			body := &model.Engine{
				Name:          "foo",
				RatedCapacity: 100,
			}

			jsonEncoded, err := json.Marshal(body)
			So(err, ShouldBeNil)

			r, err := http.NewRequest(http.MethodPost,
				ts.URL+"/api/engines", bytes.NewBuffer(jsonEncoded))
			So(err, ShouldBeNil)

			res, err := tc.Do(r)
			So(err, ShouldBeNil)
			defer res.Body.Close()

			So(res.StatusCode, ShouldEqual, http.StatusBadRequest)
		})

	})
}

func TestGetByIdEngine(t *testing.T) {
	db, _ := test.SetUpDbTest()
	ts := httptest.NewServer(NewRouter(db))
	defer ts.Close()
	tc := ts.Client()

	Convey("Testing get handler", t, func() {
		site := model.Site{
			Id:       1000,
			Name:     "site1",
			MaxPower: 100,
		}
		So(db.Insert(&site), ShouldBeNil)
		Convey("If id is known", func() {
			body := model.Engine{
				SiteId:        site.Id,
				Name:          "foo",
				Type:          "furnace",
				RatedCapacity: 100,
			}
			So(db.Insert(&body), ShouldBeNil)

			r, err := http.NewRequest(http.MethodGet,
				ts.URL+fmt.Sprintf("/api/engines/%d", body.Id), nil)
			So(err, ShouldBeNil)

			res, err := tc.Do(r)
			So(err, ShouldBeNil)
			defer res.Body.Close()

			So(res.StatusCode, ShouldEqual, http.StatusOK)
		})

		Convey("If id is unknown", func() {
			body := model.Engine{
				SiteId:        site.Id,
				Name:          "foo",
				Type:          "furnace",
				RatedCapacity: 100,
			}
			So(db.Insert(&body), ShouldBeNil)

			r, err := http.NewRequest(http.MethodGet,
				ts.URL+fmt.Sprintf("/api/engines/43543"), nil)
			So(err, ShouldBeNil)

			res, err := tc.Do(r)
			So(err, ShouldBeNil)
			defer res.Body.Close()

			So(res.StatusCode, ShouldEqual, http.StatusNotFound)
		})

	})
}

func TestGetAllEngine(t *testing.T) {
	db, _ := test.SetUpDbTest()
	ts := httptest.NewServer(NewRouter(db))
	defer ts.Close()
	tc := ts.Client()

	Convey("Testing get handler", t, func() {
		site := model.Site{
			Id:       1000,
			Name:     "site1",
			MaxPower: 1000,
		}
		So(db.Insert(&site), ShouldBeNil)

		body := model.Engine{
			SiteId:        site.Id,
			Name:          "foo",
			Type:          "furnace",
			RatedCapacity: 100,
		}

		body1 := model.Engine{
			SiteId:        site.Id,
			Name:          "foo",
			Type:          "furnace",
			RatedCapacity: 100,
		}
		So(db.Insert(&body, &body1), ShouldBeNil)

		r, err := http.NewRequest(http.MethodGet,
			ts.URL+"/api/engines", nil)
		So(err, ShouldBeNil)

		res, err := tc.Do(r)
		So(err, ShouldBeNil)
		defer res.Body.Close()

		So(res.StatusCode, ShouldEqual, http.StatusOK)
	})

}

func TestDeleteEngine(t *testing.T) {
	db, _ := test.SetUpDbTest()
	ts := httptest.NewServer(NewRouter(db))
	defer ts.Close()
	tc := ts.Client()

	Convey("Testing delete handler", t, func() {
		site := model.Site{
			Id:       1000,
			Name:     "site1",
			MaxPower: 1000,
		}
		So(db.Insert(&site), ShouldBeNil)
		Convey("If id is known", func() {
			body := model.Engine{
				SiteId:        site.Id,
				Name:          "foo",
				Type:          "furnace",
				RatedCapacity: 100,
			}

			So(db.Insert(&body), ShouldBeNil)

			r, err := http.NewRequest(http.MethodDelete,
				ts.URL+fmt.Sprintf("/api/engines/%d", body.Id), nil)
			So(err, ShouldBeNil)

			res, err := tc.Do(r)
			So(err, ShouldBeNil)
			defer res.Body.Close()

			So(res.StatusCode, ShouldEqual, http.StatusNoContent)
		})

		Convey("If id is unknown", func() {
			body := model.Engine{
				SiteId:        site.Id,
				Name:          "foo",
				Type:          "furnace",
				RatedCapacity: 100,
			}
			So(db.Insert(&body), ShouldBeNil)

			r, err := http.NewRequest(http.MethodGet,
				ts.URL+"/api/engines/43543", nil)
			So(err, ShouldBeNil)

			res, err := tc.Do(r)
			So(err, ShouldBeNil)
			defer res.Body.Close()

			So(res.StatusCode, ShouldEqual, http.StatusNotFound)
		})
	})
}

func TestUpdateEngine(t *testing.T) {
	db, _ := test.SetUpDbTest()
	ts := httptest.NewServer(NewRouter(db))
	defer ts.Close()
	tc := ts.Client()

	Convey("Testing update handler", t, func() {
		site := model.Site{
			Id:       1000,
			Name:     "site1",
			MaxPower: 1000,
		}
		So(db.Insert(&site), ShouldBeNil)

		site1 := model.Site{
			Id:       999,
			Name:     "site1",
			MaxPower: 1000,
		}
		So(db.Insert(&site, &site1), ShouldBeNil)
		Convey("If id is known", func() {
			body := model.Engine{
				SiteId:        site.Id,
				Name:          "foo",
				Type:          "furnace",
				RatedCapacity: 100,
			}
			So(db.Insert(&body), ShouldBeNil)

			body.SiteId = site1.Id
			body.Name = "changedName"
			body.Type = "compressor"
			body.RatedCapacity = 50

			jsonEncoded, err := json.Marshal(&body)
			So(err, ShouldBeNil)

			r, err := http.NewRequest(http.MethodPut,
				ts.URL+fmt.Sprintf("/api/engines/%d", body.Id), bytes.NewBuffer(jsonEncoded))
			So(err, ShouldBeNil)

			res, err := tc.Do(r)
			So(err, ShouldBeNil)
			defer res.Body.Close()

			So(res.StatusCode, ShouldEqual, http.StatusNoContent)
		})

		Convey("If id is unknown", func() {
			body := model.Engine{
				SiteId:        site.Id,
				Name:          "foo",
				Type:          "furnace",
				RatedCapacity: 100,
			}
			So(db.Insert(&body), ShouldBeNil)

			jsonEncoded, err := json.Marshal(&body)
			So(err, ShouldBeNil)

			r, err := http.NewRequest(http.MethodPut,
				ts.URL+"/api/engines/34534534", bytes.NewBuffer(jsonEncoded))
			So(err, ShouldBeNil)

			res, err := tc.Do(r)
			So(err, ShouldBeNil)
			defer res.Body.Close()

			So(res.StatusCode, ShouldEqual, http.StatusNotFound)
		})
	})
}
