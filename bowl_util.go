package bowl

import (
	"time"
)

func bowlByID(bowls []*Bowl, id string) (*Bowl, bool) {
	for _, b := range bowls {
		if b.BowlID == id {
			return b, true
		}
	}
	return nil, false
}

func kickoff(season, bowl string) time.Time {
	for _, b := range bowls[season] {
		if b.BowlID == bowl {
			return b.Kickoff
		}
	}
	return time.Time{}
}
