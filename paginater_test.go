// Copyright 2015 FineJian
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package paginator

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Paginator(t *testing.T) {
	Convey("Custom logics", t, func() {
		c := &Config{10, 2, 3}
		p := Custom(c, 20)
		pages := p.Pages()
		So(len(pages), ShouldEqual, 2)
	})
}
