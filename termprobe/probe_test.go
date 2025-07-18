package termprobe

import (
	"strings"
	"testing"

	"github.com/tiendc/gofn"
)

// XTVERSION
// DCS > | text ST
var (
	Empty = []byte{}

	Basic7    = gofn.Concat(DCS7, []byte{'>', '|', 't', 'e', 's', 't'}, ST7)
	Basic8    = []byte{DCS8, '>', '|', 't', 'e', 's', 't', ST8}
	BasicMix  = gofn.Concat(DCS7, []byte{'>', '|', 't', 'e', 's', 't', ST8})
	BasicMix2 = gofn.Concat([]byte{DCS8, '>', '|', 't', 'e', 's', 't'}, ST7)

	MissingFullPrefix7    = gofn.Concat([]byte{'>', '|', 't', 'e', 's', 't'}, ST7)
	MissingPartialPrefix7 = gofn.Concat([]byte{DCS7[0], '>', '|', 't', 'e', 's', 't'}, ST7)
	MissingPrefix8        = []byte{'>', '|', 't', 'e', 's', 't', ST8}

	MissingMid71     = gofn.Concat(DCS7, []byte{'>', 't', 'e', 's', 't'}, ST7)
	MissingMid72     = gofn.Concat(DCS7, []byte{'|', 't', 'e', 's', 't'}, ST7)
	MissingMid73     = gofn.Concat(DCS7, []byte{'t', 'e', 's', 't'}, ST7)
	MissingMid81     = []byte{DCS8, '>', 't', 'e', 's', 't', ST8}
	MissingMid82     = []byte{DCS8, '|', 't', 'e', 's', 't', ST8}
	MissingMid83     = []byte{DCS8, 't', 'e', 's', 't', ST8}
	MissingMidMix781 = gofn.Concat(DCS7, []byte{'>', 't', 'e', 's', 't', ST8})
	MissingMidMix782 = gofn.Concat(DCS7, []byte{'|', 't', 'e', 's', 't', ST8})
	MissingMidMix783 = gofn.Concat(DCS7, []byte{'t', 'e', 's', 't', ST8})
	MissingMidMix871 = gofn.Concat([]byte{DCS8, '>', 't', 'e', 's', 't'}, ST7)
	MissingMidMix872 = gofn.Concat([]byte{DCS8, '|', 't', 'e', 's', 't'}, ST7)
	MissingMidMix873 = gofn.Concat([]byte{DCS8, 't', 'e', 's', 't'}, ST7)

	MissingText7     = gofn.Concat(DCS7, []byte{'>', '|'}, ST7)
	MissingText8     = []byte{DCS8, '>', '|', ST8}
	MissingMixText78 = gofn.Concat(DCS7, []byte{'>', '|', ST8})
	MissingMixText87 = gofn.Concat([]byte{DCS8, '>', '|'}, ST7)

	MissingFullSuffix7    = gofn.Concat(DCS7, []byte{'>', '|', 't', 'e', 's', 't'})
	MissingPartialSuffix7 = gofn.Concat(DCS7, []byte{'>', '|', 't', 'e', 's', 't', ST7[0]})
	MissingSuffix8        = []byte{DCS8, '>', '|', 't', 'e', 's', 't'}
)

// TODO: ESC mid seq
// TODO: C1 mid seq
func TestXTVersionResponse(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    string
		wantErr string
	}{
		{"Empty response", Empty, "", "XTVERSION response sequence is empty"},

		{"Basic 7bit", Basic7, "test", ""},
		{"Basic 8bit", Basic8, "test", ""},
		{"Basic 7-8", BasicMix, "test", ""},
		{"Basic 8-7", BasicMix2, "test", ""},

		{"MissingText 7", MissingText7, "", ""},
		{"MissingText 8", MissingText8, "", ""},
		{"MissingMixText 78", MissingMixText78, "", ""},
		{"MissingMixText 87", MissingMixText87, "", ""},

		{"MissingMid 71", MissingMid71, "", "XTVERSION response sequence DCS prefix missing"},
		{"MissingMid 72", MissingMid72, "", "XTVERSION response sequence DCS prefix missing"},
		{"MissingMid 73", MissingMid73, "", "XTVERSION response sequence DCS prefix missing"},
		{"MissingMid 81", MissingMid81, "", "XTVERSION response sequence DCS prefix missing"},
		{"MissingMid 82", MissingMid82, "", "XTVERSION response sequence DCS prefix missing"},
		{"MissingMid 83", MissingMid83, "", "XTVERSION response sequence DCS prefix missing"},
		{"MissingMidMix 781", MissingMidMix781, "", "XTVERSION response sequence DCS prefix missing"},
		{"MissingMidMix 782", MissingMidMix782, "", "XTVERSION response sequence DCS prefix missing"},
		{"MissingMidMix 783", MissingMidMix783, "", "XTVERSION response sequence DCS prefix missing"},
		{"MissingMidMix 871", MissingMidMix871, "", "XTVERSION response sequence DCS prefix missing"},
		{"MissingMidMix 872", MissingMidMix872, "", "XTVERSION response sequence DCS prefix missing"},
		{"MissingMidMix 873", MissingMidMix873, "", "XTVERSION response sequence DCS prefix missing"},

		{"MissingFullPrefix 7", MissingFullPrefix7, "", "XTVERSION response sequence DCS prefix missing"},
		{"MissingPartialPrefix 7", MissingPartialPrefix7, "", "XTVERSION response sequence DCS prefix missing"},
		{"MissingPrefix 8", MissingPrefix8, "", "XTVERSION response sequence DCS prefix missing"},

		{"MissingFullSuffix 7", MissingFullSuffix7, "", "XTVERSION response sequence ST termination suffix missing"},
		{"MissingPartialSuffix 7", MissingPartialSuffix7, "", "XTVERSION response sequence ST termination suffix missing"},
		{"MissingSuffix 8", MissingSuffix8, "", "XTVERSION response sequence ST termination suffix missing"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans, err := XTVersionResponse(tt.input)
			if tt.want != "" && ans != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
			if tt.wantErr != "" && !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("got %v, want %s", err, tt.wantErr)
			}
		})
	}
}

// XTGETTCAP
// DCS 1 + r Pt ST
var (
	GBasic7    = gofn.Concat(DCS7, []byte{'1', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST7)
	GBasic8    = gofn.Concat([]byte{DCS8, '1', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', ST8})
	GBasicMix  = gofn.Concat(DCS7, []byte{'1', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', ST8})
	GBasicMix2 = gofn.Concat([]byte{DCS8, '1', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST7)

	GIllegalBasic7    = gofn.Concat(DCS7, []byte{'0', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST7)
	GIllegalBasic8    = gofn.Concat([]byte{DCS8, '0', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', ST8})
	GIllegalBasicMix  = gofn.Concat(DCS7, []byte{'0', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', ST8})
	GIllegalBasicMix2 = gofn.Concat([]byte{DCS8, '0', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST7)

	GMissingFullPrefix7    = gofn.Concat([]byte{'1', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST7)
	GMissingPartialPrefix7 = gofn.Concat([]byte{DCS7[0], '1', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST7)
	GMissingPrefix8        = gofn.Concat([]byte{'1', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', ST8})

	GIllegalMissingFullPrefix7    = gofn.Concat([]byte{'0', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST7)
	GIllegalMissingPartialPrefix7 = gofn.Concat([]byte{DCS7[0], '0', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST7)
	GIllegalMissingPrefix8        = gofn.Concat([]byte{'0', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', ST8})

	GMissingMid7Status   = gofn.Concat(DCS7, []byte{'+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST7)
	GMissingMid71        = gofn.Concat(DCS7, []byte{'1', '+'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST7)
	GMissingMid72        = gofn.Concat(DCS7, []byte{'1', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST7)
	GMissingMid73        = gofn.Concat(DCS7, []byte{'1'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST7)
	GIllegalMissingMid71 = gofn.Concat(DCS7, []byte{'0', '+'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST7)
	GIllegalMissingMid72 = gofn.Concat(DCS7, []byte{'0', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST7)
	GIllegalMissingMid73 = gofn.Concat(DCS7, []byte{'0'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST7)
	GMissingMid8Status   = gofn.Concat([]byte{DCS8, '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', ST8})
	GMissingMid81        = gofn.Concat([]byte{DCS8, '1', '+'}, TN, []byte{'=', 't', 'e', 's', 't', ST8})
	GMissingMid82        = gofn.Concat([]byte{DCS8, '1', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', ST8})
	GMissingMid83        = gofn.Concat([]byte{DCS8, '1'}, TN, []byte{'=', 't', 'e', 's', 't', ST8})
	GIllegalMissingMid81 = gofn.Concat([]byte{DCS8, '0', '+'}, TN, []byte{'=', 't', 'e', 's', 't', ST8})
	GIllegalMissingMid82 = gofn.Concat([]byte{DCS8, '0', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', ST8})
	GIllegalMissingMid83 = gofn.Concat([]byte{DCS8, '0'}, TN, []byte{'=', 't', 'e', 's', 't', ST8})

	GMissingMidMix78Status   = gofn.Concat(DCS7, []byte{'+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', ST8})
	GMissingMidMix781        = gofn.Concat(DCS7, []byte{'1', '+'}, TN, []byte{'=', 't', 'e', 's', 't', ST8})
	GMissingMidMix782        = gofn.Concat(DCS7, []byte{'1', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', ST8})
	GMissingMidMix783        = gofn.Concat(DCS7, []byte{'1'}, TN, []byte{'=', 't', 'e', 's', 't', ST8})
	GIllegalMissingMidMix781 = gofn.Concat(DCS7, []byte{'0', '+'}, TN, []byte{'=', 't', 'e', 's', 't', ST8})
	GIllegalMissingMidMix782 = gofn.Concat(DCS7, []byte{'0', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', ST8})
	GIllegalMissingMidMix783 = gofn.Concat(DCS7, []byte{'0'}, TN, []byte{'=', 't', 'e', 's', 't', ST8})
	GMissingMidMix87Status   = gofn.Concat([]byte{DCS8, '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST7)
	GMissingMidMix871        = gofn.Concat([]byte{DCS8, '1', '+'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST7)
	GMissingMidMix872        = gofn.Concat([]byte{DCS8, '1', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST7)
	GMissingMidMix873        = gofn.Concat([]byte{DCS8, '1'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST7)
	GIllegalMissingMidMix871 = gofn.Concat([]byte{DCS8, '0', '+'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST7)
	GIllegalMissingMidMix872 = gofn.Concat([]byte{DCS8, '0', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST7)
	GIllegalMissingMidMix873 = gofn.Concat([]byte{DCS8, '0'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST7)

	GMissingText7            = gofn.Concat(DCS7, []byte{'1', '+', 'r'}, TN, ST7)
	GMissingText8            = gofn.Concat([]byte{DCS8, '1', '+', 'r'}, TN, []byte{ST8})
	GMissingMixText78        = gofn.Concat(DCS7, []byte{'1', '+', 'r'}, TN, []byte{ST8})
	GMissingMixText87        = gofn.Concat([]byte{DCS8, '1', '+', 'r'}, TN, ST7)
	GIllegalMissingText7     = gofn.Concat(DCS7, []byte{'0', '+', 'r'}, TN, ST7)
	GIllegalMissingText8     = gofn.Concat([]byte{DCS8, '0', '+', 'r'}, TN, []byte{ST8})
	GIllegalMissingMixText78 = gofn.Concat(DCS7, []byte{'0', '+', 'r'}, TN, []byte{ST8})
	GIllegalMissingMixText87 = gofn.Concat([]byte{DCS8, '0', '+', 'r'}, TN, ST7)

	GMissingFullSuffix7           = gofn.Concat(DCS7, []byte{'1', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'})
	GMissingPartialSuffix7        = gofn.Concat(DCS7, []byte{'1', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', ST7[0]})
	GMissingSuffix8               = gofn.Concat([]byte{DCS8, '1', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'})
	GIllegalMissingFullSuffix7    = gofn.Concat(DCS7, []byte{'0', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'})
	GIllegalMissingPartialSuffix7 = gofn.Concat(DCS7, []byte{'0', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', ST7[0]})
	GIllegalMissingSuffix8        = gofn.Concat([]byte{DCS8, '0', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'})

	GPartialTn7               = gofn.Concat(DCS7, []byte{'1', '+', 'r'}, TN[:2], []byte{'=', 't', 'e', 's', 't'}, ST7)
	GPartialTn8               = gofn.Concat([]byte{DCS8, '1', '+', 'r'}, TN[:2], []byte{'=', 't', 'e', 's', 't', ST8})
	GPartialTnMix             = gofn.Concat(DCS7, []byte{'1', '+', 'r'}, TN[:2], []byte{'=', 't', 'e', 's', 't', ST8})
	GPartialTnMix2            = gofn.Concat([]byte{DCS8, '1', '+', 'r'}, TN[:2], []byte{'=', 't', 'e', 's', 't'}, ST7)
	GMissingEquals7           = gofn.Concat(DCS7, []byte{'1', '+', 'r'}, TN, []byte{'t', 'e', 's', 't'}, ST7)
	GMissingEquals8           = gofn.Concat([]byte{DCS8, '1', '+', 'r'}, TN, []byte{'t', 'e', 's', 't', ST8})
	GMissingEqualsMix         = gofn.Concat(DCS7, []byte{'1', '+', 'r'}, TN, []byte{'t', 'e', 's', 't', ST8})
	GMissingEqualsMix2        = gofn.Concat([]byte{DCS8, '1', '+', 'r'}, TN, []byte{'t', 'e', 's', 't'}, ST7)
	GIllegalPartialTn7        = gofn.Concat(DCS7, []byte{'0', '+', 'r'}, TN[:2], []byte{'=', 't', 'e', 's', 't'}, ST7)
	GIllegalPartialTn8        = gofn.Concat([]byte{DCS8, '0', '+', 'r'}, TN[:2], []byte{'=', 't', 'e', 's', 't', ST8})
	GIllegalPartialTnMix      = gofn.Concat(DCS7, []byte{'0', '+', 'r'}, TN[:2], []byte{'=', 't', 'e', 's', 't', ST8})
	GIllegalPartialTnMix2     = gofn.Concat([]byte{DCS8, '0', '+', 'r'}, TN[:2], []byte{'=', 't', 'e', 's', 't'}, ST7)
	GIllegalMissingEquals7    = gofn.Concat(DCS7, []byte{'0', '+', 'r'}, TN, []byte{'t', 'e', 's', 't'}, ST7)
	GIllegalMissingEquals8    = gofn.Concat([]byte{DCS8, '0', '+', 'r'}, TN, []byte{'t', 'e', 's', 't', ST8})
	GIllegalMissingEqualsMix  = gofn.Concat(DCS7, []byte{'0', '+', 'r'}, TN, []byte{'t', 'e', 's', 't', ST8})
	GIllegalMissingEqualsMix2 = gofn.Concat([]byte{DCS8, '0', '+', 'r'}, TN, []byte{'t', 'e', 's', 't'}, ST7)
)

func TestXGetTcapResponse(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    string
		wantErr string
	}{
		{"Empty response", Empty, "", "XTGETTCAP response sequence is empty"},

		{"Basic7", GBasic7, "test", ""},
		{"Basic8", GBasic8, "test", ""},
		{"BasicMix", GBasicMix, "test", ""},
		{"BasicMix2", GBasicMix2, "test", ""},

		{"IllegalBasic7", GIllegalBasic7, "", "XTGETTCAP illegal request"},
		{"IllegalBasic8", GIllegalBasic8, "", "XTGETTCAP illegal request"},
		{"IllegalBasicMix", GIllegalBasicMix, "", "XTGETTCAP illegal request"},
		{"IllegalBasicMix2", GIllegalBasicMix2, "", "XTGETTCAP illegal request"},

		{"MissingFullPrefix7", GMissingFullPrefix7, "", "XTGETTCAP response sequence DCS prefix missing"},
		{"MissingPartialPrefix7", GMissingPartialPrefix7, "", "XTGETTCAP response sequence DCS prefix missing"},
		{"MissingPrefix8", GMissingPrefix8, "", "XTGETTCAP response sequence DCS prefix missing"},

		{"IllegalMissingFullPrefix7", GIllegalMissingFullPrefix7, "", "XTGETTCAP response sequence DCS prefix missing"},
		{"IllegalMissingPartialPrefix7", GIllegalMissingPartialPrefix7, "", "XTGETTCAP response sequence DCS prefix missing"},
		{"IllegalMissingPrefix8", GIllegalMissingPrefix8, "", "XTGETTCAP response sequence DCS prefix missing"},

		{"MissingMid7Status", GMissingMid7Status, "", "XTGETTCAP illegal request"},
		{"MissingMid71", GMissingMid71, "", "XTGETTCAP response sequence DCS prefix missing"},
		{"MissingMid72", GMissingMid72, "", "XTGETTCAP response sequence DCS prefix missing"},
		{"MissingMid73", GMissingMid73, "", "XTGETTCAP response sequence DCS prefix missing"},
		{"IllegalMissingMid71", GIllegalMissingMid71, "", "XTGETTCAP illegal request"},
		{"IllegalMissingMid72", GIllegalMissingMid72, "", "XTGETTCAP illegal request"},
		{"IllegalMissingMid73", GIllegalMissingMid73, "", "XTGETTCAP illegal request"},
		{"MissingMid8Status", GMissingMid8Status, "", "XTGETTCAP illegal request"},
		{"MissingMid81", GMissingMid81, "", "XTGETTCAP response sequence DCS prefix missing"},
		{"MissingMid82", GMissingMid82, "", "XTGETTCAP response sequence DCS prefix missing"},
		{"MissingMid83", GMissingMid83, "", "XTGETTCAP response sequence DCS prefix missing"},
		{"IllegalMissingMid81", GIllegalMissingMid81, "", "XTGETTCAP illegal request"},
		{"IllegalMissingMid82", GIllegalMissingMid82, "", "XTGETTCAP illegal request"},
		{"IllegalMissingMid83", GIllegalMissingMid83, "", "XTGETTCAP illegal request"},

		{"MissingMidMix78Status", GMissingMidMix78Status, "", "XTGETTCAP illegal request"},
		{"MissingMidMix781", GMissingMidMix781, "", "XTGETTCAP response sequence DCS prefix missing"},
		{"MissingMidMix782", GMissingMidMix782, "", "XTGETTCAP response sequence DCS prefix missing"},
		{"MissingMidMix783", GMissingMidMix783, "", "XTGETTCAP response sequence DCS prefix missing"},
		{"IllegalMissingMidMix781", GIllegalMissingMidMix781, "", "XTGETTCAP illegal request"},
		{"IllegalMissingMidMix782", GIllegalMissingMidMix782, "", "XTGETTCAP illegal request"},
		{"IllegalMissingMidMix783", GIllegalMissingMidMix783, "", "XTGETTCAP illegal request"},
		{"MissingMidMix87Status", GMissingMidMix87Status, "", "XTGETTCAP illegal request"},
		{"MissingMidMix871", GMissingMidMix871, "", "XTGETTCAP response sequence DCS prefix missing"},
		{"MissingMidMix872", GMissingMidMix872, "", "XTGETTCAP response sequence DCS prefix missing"},
		{"MissingMidMix873", GMissingMidMix873, "", "XTGETTCAP response sequence DCS prefix missing"},
		{"IllegalMissingMidMix871", GIllegalMissingMidMix871, "", "XTGETTCAP illegal request"},
		{"IllegalMissingMidMix872", GIllegalMissingMidMix872, "", "XTGETTCAP illegal request"},
		{"IllegalMissingMidMix873", GIllegalMissingMidMix873, "", "XTGETTCAP illegal request"},

		{"MissingText7", GMissingText7, "", ""},
		{"MissingText8", GMissingText8, "", ""},
		{"MissingMixText78", GMissingMixText78, "", ""},
		{"MissingMixText87", GMissingMixText87, "", ""},
		{"IllegalMissingText7", GIllegalMissingText7, "", "XTGETTCAP illegal request"},
		{"IllegalMissingText8", GIllegalMissingText8, "", "XTGETTCAP illegal request"},
		{"IllegalMissingMixText78", GIllegalMissingMixText78, "", "XTGETTCAP illegal request"},
		{"IllegalMissingMixText87", GIllegalMissingMixText87, "", "XTGETTCAP illegal request"},

		{"MissingFullSuffix7", GMissingFullSuffix7, "", "XTGETTCAP response sequence ST termination suffix missing"},
		{"MissingPartialSuffix7", GMissingPartialSuffix7, "", "XTGETTCAP response sequence ST termination suffix missing"},
		{"MissingSuffix8", GMissingSuffix8, "", "XTGETTCAP response sequence ST termination suffix missing"},
		{"IllegalMissingFullSuffix7", GIllegalMissingFullSuffix7, "", "XTGETTCAP illegal request"},
		{"IllegalMissingPartialSuffix7", GIllegalMissingPartialSuffix7, "", "XTGETTCAP illegal request"},
		{"IllegalMissingSuffix8", GIllegalMissingSuffix8, "", "XTGETTCAP illegal request"},

		{"PartialTn7", GPartialTn7, "", "XTGETTCAP malformed terminal capability body"},
		{"PartialTn8", GPartialTn8, "", "XTGETTCAP malformed terminal capability body"},
		{"PartialTnMix", GPartialTnMix, "", "XTGETTCAP malformed terminal capability body"},
		{"PartialTnMix2", GPartialTnMix2, "", "XTGETTCAP malformed terminal capability body"},
		{"MissingEquals7", GMissingEquals7, "", "XTGETTCAP malformed terminal capability body"},
		{"MissingEquals8", GMissingEquals8, "", "XTGETTCAP malformed terminal capability body"},
		{"MissingEqualsMix", GMissingEqualsMix, "", "XTGETTCAP malformed terminal capability body"},
		{"MissingEqualsMix2", GMissingEqualsMix2, "", "XTGETTCAP malformed terminal capability body"},
		{"IllegalPartialTn7", GIllegalPartialTn7, "", "XTGETTCAP illegal request"},
		{"IllegalPartialTn8", GIllegalPartialTn8, "", "XTGETTCAP illegal request"},
		{"IllegalPartialTnMix", GIllegalPartialTnMix, "", "XTGETTCAP illegal request"},
		{"IllegalPartialTnMix2", GIllegalPartialTnMix2, "", "XTGETTCAP illegal request"},
		{"IllegalMissingEquals7", GIllegalMissingEquals7, "", "XTGETTCAP illegal request"},
		{"IllegalMissingEquals8", GIllegalMissingEquals8, "", "XTGETTCAP illegal request"},
		{"IllegalMissingEqualsMix", GIllegalMissingEqualsMix, "", "XTGETTCAP illegal request"},
		{"IllegalMissingEqualsMix2", GIllegalMissingEqualsMix2, "", "XTGETTCAP illegal request"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans, err := XTGetTcapResponse(tt.input)
			if (tt.want != "" || (tt.want == "" && tt.wantErr == "")) && ans != tt.want {
				t.Errorf("got \"%s\", want \"%s\"", ans, tt.want)
			}
			if tt.wantErr != "" && !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("got %v, want %s", err, tt.wantErr)
			}
		})
	}
}
