package date

import (
	"time"

	"github.com/vjeantet/jodaTime"
)

// CurrentDate represents current date
var CurrentDate = jodaTime.Format("YYYYMMdd", time.Now())
