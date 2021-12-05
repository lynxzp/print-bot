package bot

import "testing"

func shouldBeWrittenPages(wr string, res string) {
	realRes := s.encodePages(wr)
	if res != realRes {
		testingT.Error(`After writing "` + wr + `" pages to state expected to get "` + res + `", but got "` + realRes + `"`)
	}
}

var (
	s        state
	testingT *testing.T
)

func TestEncodePages(t *testing.T) {
	s.shureValid()
	testingT = t

	shouldBeWrittenPages("123", "123")
	shouldBeWrittenPages("1", "1")
	shouldBeWrittenPages("", "")
	shouldBeWrittenPages("1.", "1")
	shouldBeWrittenPages("1fg.4545", "14545")
	shouldBeWrittenPages("1,", "1")
	shouldBeWrittenPages(",1", "1")
	shouldBeWrittenPages("1,1", "1,1")
	shouldBeWrittenPages("12-", "12")
	shouldBeWrittenPages("-12", "12")
	shouldBeWrittenPages("12,12", "12,12")
	shouldBeWrittenPages("1,22,5", "1,22,5")
	shouldBeWrittenPages("1000,12-500", "1000,12-500")
	shouldBeWrittenPages("1000,-500", "1000,500")
}
