package render

import (
	"errors"
	"reflect"
	"testing"

	"github.com/guitarkeegan/shop-cli/internal/shopify"
)

func TestParseSetting(t *testing.T) {
	tests := []struct {
		in      string
		want    shopify.Setting
		wantErr error
	}{
		{in: "text:title:Title", want: shopify.Setting{ID: "title", Type: "text", Label: "Title"}},
		{in: "  image_picker : image : Background Image ", want: shopify.Setting{ID: "image", Type: "image_picker", Label: "Background Image"}},
		{in: "text:time:Time: 5pm", want: shopify.Setting{ID: "time", Type: "text", Label: "Time: 5pm"}},
		{in: "text:title", wantErr: ErrInvalidSetting},
		{in: "title", wantErr: ErrInvalidSetting},
		{in: "txt:title:Title", wantErr: ErrInvalidSetting},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := ParseSetting(tt.in)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestParseSettingDefaults(t *testing.T) {
	rng, err := ParseSetting("range:size:Size")
	if err != nil {
		t.Fatal(err)
	}
	if *rng.Min != 0 || *rng.Max != 100 || *rng.Step != 1 {
		t.Errorf("range defaults = %d..%d step %d, want 0..100 step 1", *rng.Min, *rng.Max, *rng.Step)
	}

	header, err := ParseSetting("header:ignored:Typography")
	if err != nil {
		t.Fatal(err)
	}
	if header.ID != "" || header.Label != "" || *header.Content != "Typography" {
		t.Errorf("header = %+v, want only content %q", header, "Typography")
	}
}

func TestParseSettings(t *testing.T) {
	tests := []struct {
		in      string
		wantIDs []string
		wantErr error
	}{
		{in: "", wantIDs: nil},
		{in: " , ", wantIDs: nil},
		{in: "text:a:A,url:b:B", wantIDs: []string{"a", "b"}},
		{in: "text:a:A,,url:b:B,", wantIDs: []string{"a", "b"}},
		{in: "text:a:A,nope", wantErr: ErrInvalidSetting},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := ParseSettings(tt.in)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			var ids []string
			for _, s := range got {
				ids = append(ids, s.ID)
			}
			if !reflect.DeepEqual(ids, tt.wantIDs) {
				t.Errorf("ids = %v, want %v", ids, tt.wantIDs)
			}
		})
	}
}

func TestParseBlocks(t *testing.T) {
	got, err := ParseBlocks([]string{"Image Slide,image_picker:image:Image", "Spacer"})
	if err != nil {
		t.Fatal(err)
	}
	want := []shopify.SectionBlock{
		{Name: "Image Slide", Type: "image-slide", Settings: []shopify.Setting{{ID: "image", Type: "image_picker", Label: "Image"}}},
		{Name: "Spacer", Type: "spacer"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}

	for _, bad := range []string{"", " ,text:a:A", "Slide,nope"} {
		if _, err := ParseBlocks([]string{bad}); err == nil {
			t.Errorf("ParseBlocks(%q) succeeded, want error", bad)
		}
	}
	if _, err := ParseBlocks([]string{"!!!"}); !errors.Is(err, ErrInvalidBlock) {
		t.Errorf("ParseBlocks(%q) err = %v, want %v", "!!!", err, ErrInvalidBlock)
	}
}

func TestParseSnippetParams(t *testing.T) {
	tests := []struct {
		in   string
		want []shopify.Param
	}{
		{in: "", want: nil},
		{in: "title", want: []shopify.Param{{Name: "title", Type: shopify.StringParam, Required: true}}},
		{in: "title, image?:object , count:number", want: []shopify.Param{
			{Name: "title", Type: shopify.StringParam, Required: true},
			{Name: "image", Type: shopify.ObjectParam},
			{Name: "count", Type: shopify.NumberParam, Required: true},
		}},
		{in: "flag?:boolean,,x:unknown", want: []shopify.Param{
			{Name: "flag", Type: shopify.BooleanParam},
			{Name: "x", Type: shopify.StringParam, Required: true},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := ParseSnippetParams(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}
