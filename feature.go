package labo

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Feature is a struct that details a unique feature that the Nintendo Labo ToyCon supports.
type Feature struct {
	Description string
	Icon        *Image
	Title       string
	Video       *Video
}

// NewFeature is a constructor function that instantiates a new Feature struct pointer.
func NewFeature(a, b, c *goquery.Selection) (*Feature, error) {
	var (
		description string
		icon        *Image
		title       string
	)
	for _, s := range []*goquery.Selection{a, b, c} {
		if ok := (s != nil); !ok {
			return nil, fmt.Errorf(errorGoQuerySelectionNil)
		}
		if ok := (s.Length() > 0); !ok {
			return nil, fmt.Errorf(errorGoQuerySelectionEmptyHTMLNodes, s)
		}
	}
	iconSelection := a.Find("button picture img")
	if ok := (iconSelection.Length() > 0); !ok {
		return nil, fmt.Errorf(errorGoQuerySelectionEmptyHTMLNodes, iconSelection)
	}
	icon, err := NewImage(iconSelection)
	if err != nil {
		return nil, err
	}
	descriptionSelection := c.Find(".copy p")
	if ok := (descriptionSelection.Length() > 0); !ok {
		return nil, fmt.Errorf(errorGoQuerySelectionEmptyHTMLNodes, descriptionSelection)
	}
	description = strings.TrimSpace(descriptionSelection.Text())
	titleSelection := c.Find(".header span")
	if ok := (titleSelection.Length() > 0); !ok {
		return nil, fmt.Errorf(errorGoQuerySelectionEmptyHTMLNodes, titleSelection)
	}
	title = strings.ToUpper(titleSelection.Text())
	feature := Feature{
		Description: description,
		Icon:        icon,
		Title:       title}
	return &feature, nil
}
