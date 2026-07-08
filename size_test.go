package humane

import "testing"

func TestSizeFormatter_Format(t *testing.T) {
	cases := []struct {
		bytes int64
		want  string
	}{
		{0, "0 B"},
		{7, "7 B"},
		{999, "999 B"},
		{79992, "80 KB"},   // shared fixture with lambada/scandalous's specs
		{225935, "226 KB"}, // matches a real file's Finder-reported size
		{500000, "500 KB"}, // matches zouk's ByteCountFormatter(.file) output
		{1500000, "1.5 MB"},
		{5_240_000_000, "5.2 GB"}, // 2 significant digits, same rounding as the KB cases above
	}

	f := SizeFormatter{}
	for _, c := range cases {
		if got := f.Format(c.bytes); got != c.want {
			t.Errorf("Format(%d) = %q, want %q", c.bytes, got, c.want)
		}
	}
}
