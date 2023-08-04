package api

import (
	"errors"
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/cloudfogtech/sync-up/internal/data"
	"github.com/cloudfogtech/sync-up/internal/exp"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestServer_GetLogAll(t *testing.T) {
	c := &gin.Context{}
	db := &data.DB{}
	s := &Server{
		db: db,
	}
	Convey("TestServer_GetLogAll", t, func() {
		Convey("page parse error", func() {
			countHandleRequestInvalid := 0
			patches := ApplyFunc(exp.HandleRequestInvalid, func(_ error) {
				countHandleRequestInvalid++
				panic(errors.New("test"))
			})
			defer patches.Reset()

			patches.ApplyMethod(c, "DefaultQuery", func(_ *gin.Context, key, defaultValue string) string {
				return "test"
			})
			So(func() { s.GetLogAll(c) }, ShouldPanic)
			So(countHandleRequestInvalid, ShouldEqual, 1)
		})
		Convey("GetAllLogList error", func() {
			countHandleRequestInvalid := 0
			patches := ApplyFunc(exp.HandleRequestInvalid, func(_ error) {
				countHandleRequestInvalid++
				panic(errors.New("test"))
			})
			defer patches.Reset()
			patches.ApplyMethod(c, "DefaultQuery", func(_ *gin.Context, key, defaultValue string) string {
				return "1"
			})
			patches.ApplyMethod(db, "GetAllLogList", func(_ *data.DB, page int) ([]data.Log, int64, error) {
				return []data.Log{}, int64(0), errors.New("GetAllLogList error")
			})
			So(func() { s.GetLogAll(c) }, ShouldPanic)
			So(countHandleRequestInvalid, ShouldEqual, 1)
		})
		Convey("success", func() {
			countHandleRequestInvalid := 0
			countJson := 0
			patches := ApplyFunc(exp.HandleRequestInvalid, func(_ error) {
				countHandleRequestInvalid++
				panic(errors.New("test"))
			})
			defer patches.Reset()
			patches.ApplyMethod(c, "DefaultQuery", func(_ *gin.Context, key, defaultValue string) string {
				return "1"
			})
			patches.ApplyMethod(db, "GetAllLogList", func(_ *data.DB, page int) ([]data.Log, int64, error) {
				return []data.Log{}, int64(0), nil
			})
			patches.ApplyMethod(c, "JSON", func(_ *gin.Context, _ int, _ any) {
				countJson++
			})
			So(func() { s.GetLogAll(c) }, ShouldNotPanic)
			So(countHandleRequestInvalid, ShouldEqual, 0)
			So(countJson, ShouldEqual, 1)
		})
	})
}

func TestServer_GetLogList(t *testing.T) {
	c := &gin.Context{}
	db := &data.DB{}
	s := &Server{
		db: db,
	}
	Convey("TestServer_GetLogList", t, func() {
		Convey("page parse error", func() {
			countHandleRequestInvalid := 0
			patches := ApplyFunc(exp.HandleRequestInvalid, func(_ error) {
				countHandleRequestInvalid++
				panic(errors.New("test"))
			})
			defer patches.Reset()

			patches.ApplyMethod(c, "DefaultQuery", func(_ *gin.Context, key, defaultValue string) string {
				return "test"
			})
			So(func() { s.GetLogList(c) }, ShouldPanic)
			So(countHandleRequestInvalid, ShouldEqual, 1)
		})
		Convey("GetLogListByServiceId error", func() {
			countHandleRequestInvalid := 0
			patches := ApplyFunc(exp.HandleRequestInvalid, func(_ error) {
				countHandleRequestInvalid++
				panic(errors.New("test"))
			})
			defer patches.Reset()
			patches.ApplyMethod(c, "DefaultQuery", func(_ *gin.Context, key, defaultValue string) string {
				return "1"
			})
			patches.ApplyMethod(c, "Param", func(_ *gin.Context, _ string) string {
				return "test"
			})
			patches.ApplyMethod(db, "GetLogListByServiceId", func(_ *data.DB, _ string, _ int) ([]data.Log, int64, error) {
				return []data.Log{}, int64(0), errors.New("GetAllLogList error")
			})
			So(func() { s.GetLogList(c) }, ShouldPanic)
			So(countHandleRequestInvalid, ShouldEqual, 1)
		})
		Convey("success", func() {
			countHandleRequestInvalid := 0
			countJson := 0
			patches := ApplyFunc(exp.HandleRequestInvalid, func(_ error) {
				countHandleRequestInvalid++
				panic(errors.New("test"))
			})
			defer patches.Reset()
			patches.ApplyMethod(c, "DefaultQuery", func(_ *gin.Context, key, defaultValue string) string {
				return "1"
			})
			patches.ApplyMethod(c, "Param", func(_ *gin.Context, _ string) string {
				return "test"
			})
			patches.ApplyMethod(db, "GetLogListByServiceId", func(_ *data.DB, _ string, _ int) ([]data.Log, int64, error) {
				return []data.Log{}, int64(0), nil
			})
			patches.ApplyMethod(c, "JSON", func(_ *gin.Context, _ int, _ any) {
				countJson++
			})
			So(func() { s.GetLogList(c) }, ShouldNotPanic)
			So(countHandleRequestInvalid, ShouldEqual, 0)
			So(countJson, ShouldEqual, 1)
		})
	})
}

func TestServer_GetLogDetail(t *testing.T) {
	c := &gin.Context{}
	db := &data.DB{}
	s := &Server{
		db: db,
	}
	Convey("TestServer_GetLogDetail", t, func() {
		Convey("GetLog error", func() {
			countHandleRequestInvalid := 0
			patches := ApplyFunc(exp.HandleRequestInvalid, func(_ error) {
				countHandleRequestInvalid++
				panic(errors.New("test"))
			})
			defer patches.Reset()
			patches.ApplyMethodSeq(c, "Param", []OutputCell{
				{Values: Params{"1"}},
				{Values: Params{"2"}},
			})
			patches.ApplyMethodSeq(db, "GetLog", []OutputCell{
				{Values: Params{nil, errors.New("test")}},
			})
			So(func() { s.GetLogDetail(c) }, ShouldPanic)
			So(countHandleRequestInvalid, ShouldEqual, 1)
		})
		Convey("success", func() {
			countHandleRequestInvalid := 0
			countJson := 0
			patches := ApplyFunc(exp.HandleRequestInvalid, func(_ error) {
				countHandleRequestInvalid++
				panic(errors.New("test"))
			})
			defer patches.Reset()
			patches.ApplyMethodSeq(c, "Param", []OutputCell{
				{Values: Params{"1"}},
				{Values: Params{"2"}},
			})
			patches.ApplyMethodSeq(db, "GetLog", []OutputCell{
				{Values: Params{data.Log{}, nil}},
			})
			patches.ApplyMethod(c, "JSON", func(_ *gin.Context, _ int, _ any) {
				countJson++
			})
			So(func() { s.GetLogDetail(c) }, ShouldNotPanic)
			So(countHandleRequestInvalid, ShouldEqual, 0)
			So(countJson, ShouldEqual, 1)
		})
	})
}
