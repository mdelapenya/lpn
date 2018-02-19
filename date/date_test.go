package date

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vjeantet/jodaTime"
)

func TestCurrentDate(t *testing.T) {
	assert := assert.New(t)
	var now = time.Now()

	year := jodaTime.Format("YYYY", now)
	month := jodaTime.Format("MM", now)
	day := jodaTime.Format("dd", now)

	assert.Equal(year+month+day, CurrentDate, "Date not properly formed.")
}
