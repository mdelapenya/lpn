// Copyright (c) 2000-present Liferay, Inc. All rights reserved.
//
// This library is free software; you can redistribute it and/or modify it under
// the terms of the GNU Lesser General Public License as published by the Free
// Software Foundation; either version 2.1 of the License, or (at your option)
// any later version.
//
// This library is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
// FOR A PARTICULAR PURPOSE. See the GNU Lesser General Public License for more
// details.

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

	assert.Equal(year+month+day, CurrentDate(), "Date not properly formed.")
}
