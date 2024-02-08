package blind_injection

import "time"

var timeStringers = []timeStringer{
	newMySQL(),
	newMySQLOld(),
	newPgSQL(),
	newMsSQL(),
}

func (b *Blind) TimeBased() (injectable bool, payload string) {
	for i := 0; i < len(b.QParameters); i++ {
		for _, timeStr := range timeStringers {
			if injectable, payload = b.timeCheck(i, timeStr.timeQuery()); injectable {
				return
			}
		}
	}
	return
}

func (b *Blind) timeCheck(parId int, query string) (bool, string) {
	start := time.Now()
	_, payload, err := b.makeRequest(parId, query)
	end := time.Since(start)
	if err != nil {
		return false, ""
	}
	if end.Milliseconds() < (time.Duration(seconds) * time.Second).Milliseconds() {
		return false, ""
	}

	return true, payload
}
