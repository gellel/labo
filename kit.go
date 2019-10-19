package labo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

// Kit is a Nintendo Labo Kit. Kit structs are built from the Nintendo Labo official website.
//
// Kits contain extended information not found on the Nintendo Labo store and may contain
// varying levels of detail depending on the type of page scraped and the category of
// the Nintendo Labo Kit. Non VR kits will not contain the corresponding VR
// information.
type Kit struct {
	BoxImage    *Image      `json:"box_image"`
	IsVR        bool        `json:"is_VR"`
	KitImage    *Image      `json:"kit_image"`
	Projects    []*Project  `json:"projects"`
	Retailers   []*Retailer `json:"retailers"`
	SoftwareBox *Image      `json:"software_box"`
	Status      string      `json:"status"`
	StatusCode  int         `json:"status_code"`
	Toycons     []*Toycon   `json:"toycons"`
	URL         *url.URL    `json:"URL"`
}

var (
	kitFn = [](func(s *goquery.Selection, k *Kit)){
		getKitBoxImage,
		getKitImage,
		getKitProjects,
		getKitRetailers,
		getKitToyCons}
)

// GetKit gets the extended Nintendo Labo Kit information from the official Nintendo Labo website.
//
// The argument Labo must contain a non NIL string Labo.Ref to perform the look-up on the Nintendo Labo website.
// Some products on the Nintendo Labo store do not contain a official Nintendo Labo website
// counterpart and will return a non-nil Kit pointer with the HTTP status code and status fields being
// set to a http.StatusBadRequest. Successful lookups should contain a Labo.StatusCode value of http.StatusOK.
func GetKit(l *Labo) *Kit {
	var (
		doc *goquery.Document
		err error
		ok  bool
		req *http.Request
		res *http.Response
		s   *goquery.Selection

		k = &Kit{
			Status:     http.StatusText(http.StatusBadRequest),
			StatusCode: http.StatusBadRequest}
	)
	ok = (l.Ref != "NIL")
	if !ok {
		return k
	}
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", laboURI, l.Ref), nil)
	ok = (err == nil)
	if !ok {
		return k
	}
	res, err = client.Do(req)
	ok = (err == nil)
	if !ok {
		return k
	}
	k.Status = res.Status
	k.StatusCode = res.StatusCode
	k.URL = req.URL
	ok = (res.StatusCode == http.StatusOK)
	if !ok {
		return k
	}
	doc, err = goquery.NewDocumentFromResponse(res)
	ok = (err == nil)
	if !ok {
		return k
	}
	s = doc.Find(htmlBody)
	ok = (s.Length() > 0)
	if !ok {
		return k
	}
	return newKit(s, k)
}

// MarshalKit marshals a Kit struct into an ordered byte sequence. On error returns an empty byte slice.
func MarshalKit(k *Kit) (b []byte) {
	b, _ = json.Marshal(k)
	return b
}

// getKitBoxImage gets the Nintendo Labo Kits product box image.
func getKitBoxImage(s *goquery.Selection, k *Kit) {
	const (
		CSS string = ".product-hero .hero-content .kit.column .packshot picture img"
	)
	var (
		ok bool
	)
	s = s.Find(CSS)
	ok = (s.Length() > 0)
	if !ok {
		return
	}
	k.BoxImage = newImage(s)
}

func getKitImage(s *goquery.Selection, k *Kit) {
	const (
		CSS string = ".kit-contents > picture:nth-child(1) > img:nth-child(1)"
	)
	var (
		ok bool
	)
	s = s.Find(CSS)
	ok = (s.Length() > 0)
	if !ok {
		return
	}
	k.KitImage = newImage(s)
}

func getKitProjects(s *goquery.Selection, k *Kit) {
	const (
		CSS string = ".main-toycon:nth-child(1) > .toycon-tag"
	)
	var (
		ok bool
	)
	s = s.Find(CSS)
	ok = (s.Length() > 0)
	if !ok {
		return
	}
	k.Projects = newProjects(s)
}

func getKitToyCons(s *goquery.Selection, k *Kit) {
	const (
		CSS string = ".toy-con-area .toy-con"
	)
	var (
		ok bool
	)
	s = s.Find(CSS)
	ok = (s.Length() > 0)
	if !ok {
		return
	}
	k.Toycons = newToycons(s)
}

func getKitRetailers(s *goquery.Selection, k *Kit) {
	const (
		CSS string = ".retailers > ul:nth-child(1) li"
	)
	var (
		ok bool
	)
	s = s.Find(CSS)
	ok = (s.Length() > 0)
	if !ok {
		return
	}
	k.Retailers = newRetailers(s)
}

func getKitSoftwareBox(s *goquery.Selection, k *Kit) {
	const (
		CSS string = ".content:nth-child(2) > div:nth-child(1) > picture:nth-child(1) > img:nth-child(1)"
	)
	var (
		ok bool
	)
	s = s.Find(CSS)
	ok = (s.Length() > 0)
	if !ok {
		return
	}
	k.SoftwareBox = newImage(s)
}

func newKit(s *goquery.Selection, k *Kit) *Kit {
	for _, fn := range kitFn {
		fn(s, k)
	}
	return k
}
