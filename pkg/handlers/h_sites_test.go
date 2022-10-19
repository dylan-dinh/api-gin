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

func TestPostSite(t *testing.T) {
	Convey("Testing post handler", t, func() {
		Convey("If body is well formatted", func() {
			db, _ := test.SetUpDbTest()
			ts := httptest.NewServer(NewRouter(db))
			defer ts.Close()
			defer test.ClearDB(db)
			tc := ts.Client()
			body := &model.Site{
				Name:     "foo",
				MaxPower: 100,
			}

			jsonEncoded, err := json.Marshal(body)
			So(err, ShouldBeNil)

			r, err := http.NewRequest(http.MethodPost,
				ts.URL+"/api/sites", bytes.NewBuffer(jsonEncoded))
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

			body := &model.Site{
				Name:     "",
				MaxPower: -10,
			}

			jsonEncoded, err := json.Marshal(body)
			So(err, ShouldBeNil)

			r, err := http.NewRequest(http.MethodPost,
				ts.URL+"/api/sites", bytes.NewBuffer(jsonEncoded))
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
			body := &model.Site{
				Name: "foo",
			}

			jsonEncoded, err := json.Marshal(body)
			So(err, ShouldBeNil)

			r, err := http.NewRequest(http.MethodPost,
				ts.URL+"/api/sites", bytes.NewBuffer(jsonEncoded))
			So(err, ShouldBeNil)

			res, err := tc.Do(r)
			So(err, ShouldBeNil)
			defer res.Body.Close()

			So(res.StatusCode, ShouldEqual, http.StatusBadRequest)
		})

	})
}

func TestGetByIdSite(t *testing.T) {
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
			r, err := http.NewRequest(http.MethodGet,
				ts.URL+fmt.Sprintf("/api/sites/%d", site.Id), nil)
			So(err, ShouldBeNil)

			res, err := tc.Do(r)
			So(err, ShouldBeNil)
			defer res.Body.Close()

			So(res.StatusCode, ShouldEqual, http.StatusOK)
		})

		Convey("If id is unknown", func() {
			r, err := http.NewRequest(http.MethodGet,
				ts.URL+fmt.Sprintf("/api/sites/43543"), nil)
			So(err, ShouldBeNil)

			res, err := tc.Do(r)
			So(err, ShouldBeNil)
			defer res.Body.Close()

			So(res.StatusCode, ShouldEqual, http.StatusNotFound)
		})

	})
}

func TestGetAllSite(t *testing.T) {
	db, _ := test.SetUpDbTest()
	ts := httptest.NewServer(NewRouter(db))
	defer ts.Close()
	tc := ts.Client()

	Convey("Testing get all handler", t, func() {
		site := model.Site{
			Id:       1000,
			Name:     "site1",
			MaxPower: 1000,
		}

		site1 := model.Site{
			Id:       1000,
			Name:     "site1",
			MaxPower: 1000,
		}
		So(db.Insert(&site, &site1), ShouldBeNil)

		r, err := http.NewRequest(http.MethodGet,
			ts.URL+"/api/sites", nil)
		So(err, ShouldBeNil)

		res, err := tc.Do(r)
		So(err, ShouldBeNil)
		defer res.Body.Close()

		So(res.StatusCode, ShouldEqual, http.StatusOK)
	})

}

func TestDeleteSite(t *testing.T) {
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
			r, err := http.NewRequest(http.MethodDelete,
				ts.URL+fmt.Sprintf("/api/sites/%d", site.Id), nil)
			So(err, ShouldBeNil)

			res, err := tc.Do(r)
			So(err, ShouldBeNil)
			defer res.Body.Close()

			So(res.StatusCode, ShouldEqual, http.StatusNoContent)
		})

		Convey("If id is unknown", func() {
			r, err := http.NewRequest(http.MethodGet,
				ts.URL+"/api/sites/43543", nil)
			So(err, ShouldBeNil)

			res, err := tc.Do(r)
			So(err, ShouldBeNil)
			defer res.Body.Close()

			So(res.StatusCode, ShouldEqual, http.StatusNotFound)
		})
	})
}

func TestUpdateSite(t *testing.T) {
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

		Convey("If id is known", func() {
			site.MaxPower = 10
			site.Name = "changedName"

			jsonEncoded, err := json.Marshal(&site)
			So(err, ShouldBeNil)

			r, err := http.NewRequest(http.MethodPut,
				ts.URL+fmt.Sprintf("/api/sites/%d", site.Id), bytes.NewBuffer(jsonEncoded))
			So(err, ShouldBeNil)

			res, err := tc.Do(r)
			So(err, ShouldBeNil)
			defer res.Body.Close()

			So(res.StatusCode, ShouldEqual, http.StatusNoContent)
		})

		Convey("If id is unknown", func() {
			jsonEncoded, err := json.Marshal(&site)
			So(err, ShouldBeNil)

			r, err := http.NewRequest(http.MethodPut,
				ts.URL+"/api/sites/34534534", bytes.NewBuffer(jsonEncoded))
			So(err, ShouldBeNil)

			res, err := tc.Do(r)
			So(err, ShouldBeNil)
			defer res.Body.Close()

			So(res.StatusCode, ShouldEqual, http.StatusNotFound)
		})

		Convey("If missing fields", func() {
			body := model.Site{
				Id:   site.Id,
				Name: "foo",
			}

			jsonEncoded, err := json.Marshal(&body)
			So(err, ShouldBeNil)

			r, err := http.NewRequest(http.MethodPut,
				ts.URL+fmt.Sprintf("/api/sites/%d", site.Id), bytes.NewBuffer(jsonEncoded))
			So(err, ShouldBeNil)

			res, err := tc.Do(r)
			So(err, ShouldBeNil)
			defer res.Body.Close()

			So(res.StatusCode, ShouldEqual, http.StatusBadRequest)
		})
	})
}
