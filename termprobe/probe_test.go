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

	Basic7    = gofn.Concat(DCS, []byte{'>', '|', 't', 'e', 's', 't'}, ST)
	Basic8    = []byte{C1DCS, '>', '|', 't', 'e', 's', 't', C1ST}
	BasicMix  = gofn.Concat(DCS, []byte{'>', '|', 't', 'e', 's', 't', C1ST})
	BasicMix2 = gofn.Concat([]byte{C1DCS, '>', '|', 't', 'e', 's', 't'}, ST)

	MissingFullPrefix7    = gofn.Concat([]byte{'>', '|', 't', 'e', 's', 't'}, ST)
	MissingPartialPrefix7 = gofn.Concat([]byte{DCS[0], '>', '|', 't', 'e', 's', 't'}, ST)
	MissingPrefix8        = []byte{'>', '|', 't', 'e', 's', 't', C1ST}

	MissingMid71     = gofn.Concat(DCS, []byte{'>', 't', 'e', 's', 't'}, ST)
	MissingMid72     = gofn.Concat(DCS, []byte{'|', 't', 'e', 's', 't'}, ST)
	MissingMid73     = gofn.Concat(DCS, []byte{'t', 'e', 's', 't'}, ST)
	MissingMid81     = []byte{C1DCS, '>', 't', 'e', 's', 't', C1ST}
	MissingMid82     = []byte{C1DCS, '|', 't', 'e', 's', 't', C1ST}
	MissingMid83     = []byte{C1DCS, 't', 'e', 's', 't', C1ST}
	MissingMidMix781 = gofn.Concat(DCS, []byte{'>', 't', 'e', 's', 't', C1ST})
	MissingMidMix782 = gofn.Concat(DCS, []byte{'|', 't', 'e', 's', 't', C1ST})
	MissingMidMix783 = gofn.Concat(DCS, []byte{'t', 'e', 's', 't', C1ST})
	MissingMidMix871 = gofn.Concat([]byte{C1DCS, '>', 't', 'e', 's', 't'}, ST)
	MissingMidMix872 = gofn.Concat([]byte{C1DCS, '|', 't', 'e', 's', 't'}, ST)
	MissingMidMix873 = gofn.Concat([]byte{C1DCS, 't', 'e', 's', 't'}, ST)

	PaddedMid7    = gofn.Concat(DCS, []byte{'+', '+', '>', '|', 't', 'e', 's', 't'}, ST)
	PaddedMid8    = []byte{C1DCS, '+', '+', '>', '|', 't', 'e', 's', 't', C1ST}
	PaddedMidMix  = gofn.Concat(DCS, []byte{'+', '>', '|', 't', 'e', 's', 't', C1ST})
	PaddedMidMix2 = gofn.Concat([]byte{C1DCS, '+', '>', '|', 't', 'e', 's', 't'}, ST)

	MissingText7     = gofn.Concat(DCS, []byte{'>', '|'}, ST)
	MissingText8     = []byte{C1DCS, '>', '|', C1ST}
	MissingMixText78 = gofn.Concat(DCS, []byte{'>', '|', C1ST})
	MissingMixText87 = gofn.Concat([]byte{C1DCS, '>', '|'}, ST)

	MissingFullSuffix7    = gofn.Concat(DCS, []byte{'>', '|', 't', 'e', 's', 't'})
	MissingPartialSuffix7 = gofn.Concat(DCS, []byte{'>', '|', 't', 'e', 's', 't', ST[0]})
	MissingSuffix8        = []byte{C1DCS, '>', '|', 't', 'e', 's', 't'}
)

// TODO: ESC mid seq
// TODO: C1 mid seq
func TestXTVersionResponse(t *testing.T) {
	msgs := map[ErrorCode]string{
		EmptyResponse:                          "XTVERSION response sequence is empty",
		BadResponseMissingPrefixSequence:       "XTVERSION response sequence DCS > | prefix missing",
		BadResponseMissingTerminiationSequence: "XTVERSION response sequence ST termination suffix missing",
	}
	tests := []struct {
		name    string
		input   []byte
		want    string
		wantErr ErrorCode
	}{
		{"Empty response", Empty, "", EmptyResponse},
		{"Basic 7bit", Basic7, "test", ZeroError},
		{"Basic 8bit", Basic8, "test", ZeroError},
		{"Basic 7-8", BasicMix, "test", ZeroError},
		{"Basic 8-7", BasicMix2, "test", ZeroError},
		{"MissingText 7", MissingText7, "", ZeroError},
		{"MissingText 8", MissingText8, "", ZeroError},
		{"MissingMixText 78", MissingMixText78, "", ZeroError},
		{"MissingMixText 87", MissingMixText87, "", ZeroError},
		{"MissingMid 71", MissingMid71, "", BadResponseMissingPrefixSequence},
		{"MissingMid 72", MissingMid72, "", BadResponseMissingPrefixSequence},
		{"MissingMid 73", MissingMid73, "", BadResponseMissingPrefixSequence},
		{"MissingMid 81", MissingMid81, "", BadResponseMissingPrefixSequence},
		{"MissingMid 82", MissingMid82, "", BadResponseMissingPrefixSequence},
		{"MissingMid 83", MissingMid83, "", BadResponseMissingPrefixSequence},
		{"MissingMidMix 781", MissingMidMix781, "", BadResponseMissingPrefixSequence},
		{"MissingMidMix 782", MissingMidMix782, "", BadResponseMissingPrefixSequence},
		{"MissingMidMix 783", MissingMidMix783, "", BadResponseMissingPrefixSequence},
		{"MissingMidMix 871", MissingMidMix871, "", BadResponseMissingPrefixSequence},
		{"MissingMidMix 872", MissingMidMix872, "", BadResponseMissingPrefixSequence},
		{"MissingMidMix 873", MissingMidMix873, "", BadResponseMissingPrefixSequence},
		{"PaddedMid 7", PaddedMid7, "", BadResponseMissingPrefixSequence},
		{"PaddedMid 8", PaddedMid8, "", BadResponseMissingPrefixSequence},
		{"PaddedMidMix 1", PaddedMidMix, "", BadResponseMissingPrefixSequence},
		{"PaddedMidMix 2", PaddedMidMix2, "", BadResponseMissingPrefixSequence},
		{"MissingFullPrefix 7", MissingFullPrefix7, "", BadResponseMissingPrefixSequence},
		{"MissingPartialPrefix 7", MissingPartialPrefix7, "", BadResponseMissingPrefixSequence},
		{"MissingPrefix 8", MissingPrefix8, "", BadResponseMissingPrefixSequence},
		{"MissingFullSuffix 7", MissingFullSuffix7, "", BadResponseMissingTerminiationSequence},
		{"MissingPartialSuffix 7", MissingPartialSuffix7, "", BadResponseMissingTerminiationSequence},
		{"MissingSuffix 8", MissingSuffix8, "", BadResponseMissingTerminiationSequence},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans, err := ParseXTVERSIONResponse(tt.input)

			if err == nil && tt.want != "" && ans != tt.want {
				t.Errorf("got value: %s, expected value: %s", ans, tt.want)
			}

			if err != nil {
				perr := IsProbeErrorOrUnknown(err, XtVersion)
				if tt.wantErr == ZeroError {
					t.Errorf("got error: %x, expected value: %s", err, tt.want)
				} else if tt.wantErr != perr.Code || msgs[tt.wantErr] != perr.Message {
					t.Errorf(
						"got error: %s [msg is: %s], expected error: %s [msg is: %s]",
						perr.Code.String(),
						perr.Message,
						tt.wantErr.String(),
						msgs[tt.wantErr],
					)
				}
			}
		})
	}
}

// XTGETTCAP
// DCS 1 + r Pt ST
var (
	GBasic7    = gofn.Concat(DCS, []byte{'1', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)
	GBasic8    = gofn.Concat([]byte{C1DCS, '1', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', C1ST})
	GBasicMix  = gofn.Concat(DCS, []byte{'1', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', C1ST})
	GBasicMix2 = gofn.Concat([]byte{C1DCS, '1', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)

	GIllegalBasic7    = gofn.Concat(DCS, []byte{'0', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)
	GIllegalBasic8    = gofn.Concat([]byte{C1DCS, '0', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', C1ST})
	GIllegalBasicMix  = gofn.Concat(DCS, []byte{'0', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', C1ST})
	GIllegalBasicMix2 = gofn.Concat([]byte{C1DCS, '0', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)

	GMissingFullPrefix7    = gofn.Concat([]byte{'1', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)
	GMissingPartialPrefix7 = gofn.Concat([]byte{DCS[0], '1', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)
	GMissingPrefix8        = gofn.Concat([]byte{'1', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', C1ST})

	GIllegalMissingFullPrefix7    = gofn.Concat([]byte{'0', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)
	GIllegalMissingPartialPrefix7 = gofn.Concat([]byte{DCS[0], '0', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)
	GIllegalMissingPrefix8        = gofn.Concat([]byte{'0', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', C1ST})

	GMissingMid7Status   = gofn.Concat(DCS, []byte{'+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)
	GMissingMid71        = gofn.Concat(DCS, []byte{'1', '+'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)
	GMissingMid72        = gofn.Concat(DCS, []byte{'1', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)
	GMissingMid73        = gofn.Concat(DCS, []byte{'1'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)
	GIllegalMissingMid71 = gofn.Concat(DCS, []byte{'0', '+'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)
	GIllegalMissingMid72 = gofn.Concat(DCS, []byte{'0', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)
	GIllegalMissingMid73 = gofn.Concat(DCS, []byte{'0'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)
	GMissingMid8Status   = gofn.Concat([]byte{C1DCS, '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', C1ST})
	GMissingMid81        = gofn.Concat([]byte{C1DCS, '1', '+'}, TN, []byte{'=', 't', 'e', 's', 't', C1ST})
	GMissingMid82        = gofn.Concat([]byte{C1DCS, '1', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', C1ST})
	GMissingMid83        = gofn.Concat([]byte{C1DCS, '1'}, TN, []byte{'=', 't', 'e', 's', 't', C1ST})
	GIllegalMissingMid81 = gofn.Concat([]byte{C1DCS, '0', '+'}, TN, []byte{'=', 't', 'e', 's', 't', C1ST})
	GIllegalMissingMid82 = gofn.Concat([]byte{C1DCS, '0', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', C1ST})
	GIllegalMissingMid83 = gofn.Concat([]byte{C1DCS, '0'}, TN, []byte{'=', 't', 'e', 's', 't', C1ST})

	GMissingMidMix78Status   = gofn.Concat(DCS, []byte{'+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', C1ST})
	GMissingMidMix781        = gofn.Concat(DCS, []byte{'1', '+'}, TN, []byte{'=', 't', 'e', 's', 't', C1ST})
	GMissingMidMix782        = gofn.Concat(DCS, []byte{'1', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', C1ST})
	GMissingMidMix783        = gofn.Concat(DCS, []byte{'1'}, TN, []byte{'=', 't', 'e', 's', 't', C1ST})
	GIllegalMissingMidMix781 = gofn.Concat(DCS, []byte{'0', '+'}, TN, []byte{'=', 't', 'e', 's', 't', C1ST})
	GIllegalMissingMidMix782 = gofn.Concat(DCS, []byte{'0', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', C1ST})
	GIllegalMissingMidMix783 = gofn.Concat(DCS, []byte{'0'}, TN, []byte{'=', 't', 'e', 's', 't', C1ST})
	GMissingMidMix87Status   = gofn.Concat([]byte{C1DCS, '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)
	GMissingMidMix871        = gofn.Concat([]byte{C1DCS, '1', '+'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)
	GMissingMidMix872        = gofn.Concat([]byte{C1DCS, '1', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)
	GMissingMidMix873        = gofn.Concat([]byte{C1DCS, '1'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)
	GIllegalMissingMidMix871 = gofn.Concat([]byte{C1DCS, '0', '+'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)
	GIllegalMissingMidMix872 = gofn.Concat([]byte{C1DCS, '0', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)
	GIllegalMissingMidMix873 = gofn.Concat([]byte{C1DCS, '0'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)

	GMissingText7            = gofn.Concat(DCS, []byte{'1', '+', 'r'}, TN, ST)
	GMissingText8            = gofn.Concat([]byte{C1DCS, '1', '+', 'r'}, TN, []byte{C1ST})
	GMissingMixText78        = gofn.Concat(DCS, []byte{'1', '+', 'r'}, TN, []byte{C1ST})
	GMissingMixText87        = gofn.Concat([]byte{C1DCS, '1', '+', 'r'}, TN, ST)
	GIllegalMissingText7     = gofn.Concat(DCS, []byte{'0', '+', 'r'}, TN, ST)
	GIllegalMissingText8     = gofn.Concat([]byte{C1DCS, '0', '+', 'r'}, TN, []byte{C1ST})
	GIllegalMissingMixText78 = gofn.Concat(DCS, []byte{'0', '+', 'r'}, TN, []byte{C1ST})
	GIllegalMissingMixText87 = gofn.Concat([]byte{C1DCS, '0', '+', 'r'}, TN, ST)

	GMissingFullSuffix7           = gofn.Concat(DCS, []byte{'1', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'})
	GMissingPartialSuffix7        = gofn.Concat(DCS, []byte{'1', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', ST[0]})
	GMissingSuffix8               = gofn.Concat([]byte{C1DCS, '1', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'})
	GIllegalMissingFullSuffix7    = gofn.Concat(DCS, []byte{'0', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'})
	GIllegalMissingPartialSuffix7 = gofn.Concat(DCS, []byte{'0', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', ST[0]})
	GIllegalMissingSuffix8        = gofn.Concat([]byte{C1DCS, '0', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'})

	GPartialTn7               = gofn.Concat(DCS, []byte{'1', '+', 'r'}, TN[:2], []byte{'=', 't', 'e', 's', 't'}, ST)
	GPartialTn8               = gofn.Concat([]byte{C1DCS, '1', '+', 'r'}, TN[:2], []byte{'=', 't', 'e', 's', 't', C1ST})
	GPartialTnMix             = gofn.Concat(DCS, []byte{'1', '+', 'r'}, TN[:2], []byte{'=', 't', 'e', 's', 't', C1ST})
	GPartialTnMix2            = gofn.Concat([]byte{C1DCS, '1', '+', 'r'}, TN[:2], []byte{'=', 't', 'e', 's', 't'}, ST)
	GMissingEquals7           = gofn.Concat(DCS, []byte{'1', '+', 'r'}, TN, []byte{'t', 'e', 's', 't'}, ST)
	GMissingEquals8           = gofn.Concat([]byte{C1DCS, '1', '+', 'r'}, TN, []byte{'t', 'e', 's', 't', C1ST})
	GMissingEqualsMix         = gofn.Concat(DCS, []byte{'1', '+', 'r'}, TN, []byte{'t', 'e', 's', 't', C1ST})
	GMissingEqualsMix2        = gofn.Concat([]byte{C1DCS, '1', '+', 'r'}, TN, []byte{'t', 'e', 's', 't'}, ST)
	GIllegalPartialTn7        = gofn.Concat(DCS, []byte{'0', '+', 'r'}, TN[:2], []byte{'=', 't', 'e', 's', 't'}, ST)
	GIllegalPartialTn8        = gofn.Concat([]byte{C1DCS, '0', '+', 'r'}, TN[:2], []byte{'=', 't', 'e', 's', 't', C1ST})
	GIllegalPartialTnMix      = gofn.Concat(DCS, []byte{'0', '+', 'r'}, TN[:2], []byte{'=', 't', 'e', 's', 't', C1ST})
	GIllegalPartialTnMix2     = gofn.Concat([]byte{C1DCS, '0', '+', 'r'}, TN[:2], []byte{'=', 't', 'e', 's', 't'}, ST)
	GIllegalMissingEquals7    = gofn.Concat(DCS, []byte{'0', '+', 'r'}, TN, []byte{'t', 'e', 's', 't'}, ST)
	GIllegalMissingEquals8    = gofn.Concat([]byte{C1DCS, '0', '+', 'r'}, TN, []byte{'t', 'e', 's', 't', C1ST})
	GIllegalMissingEqualsMix  = gofn.Concat(DCS, []byte{'0', '+', 'r'}, TN, []byte{'t', 'e', 's', 't', C1ST})
	GIllegalMissingEqualsMix2 = gofn.Concat([]byte{C1DCS, '0', '+', 'r'}, TN, []byte{'t', 'e', 's', 't'}, ST)

	GPadding7    = gofn.Concat(DCS, []byte{'1', 'X', 'X', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)
	GPadding8    = gofn.Concat([]byte{C1DCS, '1', 'X', 'X', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', C1ST})
	GPaddingMix  = gofn.Concat(DCS, []byte{'1', 'X', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't', C1ST})
	GPaddingMix2 = gofn.Concat([]byte{C1DCS, '1', 'X', '+', 'r'}, TN, []byte{'=', 't', 'e', 's', 't'}, ST)
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

		{"MissingFullPrefix7", GMissingFullPrefix7, "", "XTGETTCAP response sequence DCS 1 + r prefix missing"},
		{"MissingPartialPrefix7", GMissingPartialPrefix7, "", "XTGETTCAP response sequence DCS 1 + r prefix missing"},
		{"MissingPrefix8", GMissingPrefix8, "", "XTGETTCAP response sequence DCS 1 + r prefix missing"},

		{"IllegalMissingFullPrefix7", GIllegalMissingFullPrefix7, "", "XTGETTCAP response sequence DCS 1 + r prefix missing"},
		{"IllegalMissingPartialPrefix7", GIllegalMissingPartialPrefix7, "", "XTGETTCAP response sequence DCS 1 + r prefix missing"},
		{"IllegalMissingPrefix8", GIllegalMissingPrefix8, "", "XTGETTCAP response sequence DCS 1 + r prefix missing"},

		{"MissingMid7Status", GMissingMid7Status, "", "XTGETTCAP illegal request"},
		{"MissingMid71", GMissingMid71, "", "XTGETTCAP response sequence DCS 1 + r prefix missing"},
		{"MissingMid72", GMissingMid72, "", "XTGETTCAP response sequence DCS 1 + r prefix missing"},
		{"MissingMid73", GMissingMid73, "", "XTGETTCAP response sequence DCS 1 + r prefix missing"},
		{"IllegalMissingMid71", GIllegalMissingMid71, "", "XTGETTCAP illegal request"},
		{"IllegalMissingMid72", GIllegalMissingMid72, "", "XTGETTCAP illegal request"},
		{"IllegalMissingMid73", GIllegalMissingMid73, "", "XTGETTCAP illegal request"},
		{"MissingMid8Status", GMissingMid8Status, "", "XTGETTCAP illegal request"},
		{"MissingMid81", GMissingMid81, "", "XTGETTCAP response sequence DCS 1 + r prefix missing"},
		{"MissingMid82", GMissingMid82, "", "XTGETTCAP response sequence DCS 1 + r prefix missing"},
		{"MissingMid83", GMissingMid83, "", "XTGETTCAP response sequence DCS 1 + r prefix missing"},
		{"IllegalMissingMid81", GIllegalMissingMid81, "", "XTGETTCAP illegal request"},
		{"IllegalMissingMid82", GIllegalMissingMid82, "", "XTGETTCAP illegal request"},
		{"IllegalMissingMid83", GIllegalMissingMid83, "", "XTGETTCAP illegal request"},

		{"MissingMidMix78Status", GMissingMidMix78Status, "", "XTGETTCAP illegal request"},
		{"MissingMidMix781", GMissingMidMix781, "", "XTGETTCAP response sequence DCS 1 + r prefix missing"},
		{"MissingMidMix782", GMissingMidMix782, "", "XTGETTCAP response sequence DCS 1 + r prefix missing"},
		{"MissingMidMix783", GMissingMidMix783, "", "XTGETTCAP response sequence DCS 1 + r prefix missing"},
		{"IllegalMissingMidMix781", GIllegalMissingMidMix781, "", "XTGETTCAP illegal request"},
		{"IllegalMissingMidMix782", GIllegalMissingMidMix782, "", "XTGETTCAP illegal request"},
		{"IllegalMissingMidMix783", GIllegalMissingMidMix783, "", "XTGETTCAP illegal request"},
		{"MissingMidMix87Status", GMissingMidMix87Status, "", "XTGETTCAP illegal request"},
		{"MissingMidMix871", GMissingMidMix871, "", "XTGETTCAP response sequence DCS 1 + r prefix missing"},
		{"MissingMidMix872", GMissingMidMix872, "", "XTGETTCAP response sequence DCS 1 + r prefix missing"},
		{"MissingMidMix873", GMissingMidMix873, "", "XTGETTCAP response sequence DCS 1 + r prefix missing"},
		{"IllegalMissingMidMix871", GIllegalMissingMidMix871, "", "XTGETTCAP illegal request"},
		{"IllegalMissingMidMix872", GIllegalMissingMidMix872, "", "XTGETTCAP illegal request"},
		{"IllegalMissingMidMix873", GIllegalMissingMidMix873, "", "XTGETTCAP illegal request"},

		{"GPadding7", GPadding7, "", "XTGETTCAP response sequence DCS 1 + r prefix missing"},
		{"GPadding8", GPadding8, "", "XTGETTCAP response sequence DCS 1 + r prefix missing"},
		{"GPaddingMix", GPaddingMix, "", "XTGETTCAP response sequence DCS 1 + r prefix missing"},
		{"GPaddingMix2", GPaddingMix2, "", "XTGETTCAP response sequence DCS 1 + r prefix missing"},

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
			ans, err := ParseXTGETTCAPResponse(tt.input)
			if (tt.want != "" || (tt.want == "" && tt.wantErr == "")) && ans != tt.want {
				t.Errorf("got \"%s\", want \"%s\"", ans, tt.want)
			}
			if tt.wantErr != "" && !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("got %v, want %s", err, tt.wantErr)
			}
		})
	}
}
