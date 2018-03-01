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
	"net/http"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Paginator(t *testing.T) {
	Convey("Basic logics", t, func() {
		p := New(20)
		So(len(p.Pages()), ShouldEqual, 2)
	})

	Convey("Custom logics", t, func() {
		p := Custom(&Config{PageSize: 10, Current: 2, LinkedCount: 3}, 23)
		So(len(p.Pages()), ShouldEqual, 3)
		So(p.TotalPages(), ShouldEqual, 3)
		So(p.IsFirst(), ShouldBeFalse)
		So(p.HasPrevious(), ShouldBeTrue)
		So(p.Previous(), ShouldEqual, 1)
		So(p.HasNext(), ShouldBeTrue)
		So(p.Next(), ShouldEqual, 3)
		So(p.IsLast(), ShouldBeFalse)

		Convey("LinkedCount 小于等于0", func() {
			p := Custom(&Config{PageSize: 10, Current: 2, LinkedCount: 0}, 23)
			So(len(p.Pages()), ShouldEqual, 0)

			p = Custom(&Config{PageSize: 10, Current: 2, LinkedCount: 0}, 5)
			So(len(p.Pages()), ShouldEqual, 0)
			So(p.Current(), ShouldEqual, 1)
		})

		Convey("LinkedCount 等于1 或者 TotalPages 等于1", func() {
			p := Custom(&Config{PageSize: 10, Current: 2, LinkedCount: 1}, 23)
			So(len(p.Pages()), ShouldEqual, 1)

			p = Custom(&Config{PageSize: 10, Current: 2, LinkedCount: 1}, 5)
			So(len(p.Pages()), ShouldEqual, 1)
		})

		Convey("TotalPages 大于 LinkedCount", func() {
			Convey("LinkedCount 是奇数", func() {
				Print("\n")
				p := Custom(&Config{PageSize: 10, Current: 1, LinkedCount: 3}, 63)
				for _, page := range p.Pages() {
					Printf("%v ", page)
				}
				Print("\n")

				p = Custom(&Config{PageSize: 10, Current: 3, LinkedCount: 3}, 63)
				for _, page := range p.Pages() {
					Printf("%v ", page)
				}
				Print("\n")

				p = Custom(&Config{PageSize: 10, Current: 6, LinkedCount: 3}, 63)
				for _, page := range p.Pages() {
					Printf("%v ", page)
				}
				Print("\n")
			})

			Convey("LinkedCount 是偶数", func() {
				Print("\n")
				p := Custom(&Config{PageSize: 10, Current: 1, LinkedCount: 4}, 63)
				for _, page := range p.Pages() {
					Printf("%v ", page)
				}
				Print("\n")

				p = Custom(&Config{PageSize: 10, Current: 4, LinkedCount: 4}, 63)
				for _, page := range p.Pages() {
					Printf("%v ", page)
				}
				Print("\n")

				p = Custom(&Config{PageSize: 10, Current: 6, LinkedCount: 4}, 63)
				for _, page := range p.Pages() {
					Printf("%v ", page)
				}
				Print("\n")
			})
		})
	})

	Convey("Url logics", t, func() {
		p := Custom(&Config{PageSize: 10, Current: 6, LinkedCount: 4}, 63)
		r := &http.Request{URL: &url.URL{Path: "http://www.baidu.com", RawQuery: "a=b&b=c&page=2"}}
		p = p.Request(r)
		Printf("FristURL: %s \n", p.FristURL())
		Printf("PreviousURL: %s \n", p.PreviousURL())
		Printf("CurrentURL: %s \n", p.CurrentURL())
		Printf("NextURL: %s \n", p.NextURL())
		Printf("LastURL: %s \n", p.LastURL())
		for _, page := range p.Pages() {
			Printf("%v %s \n", page, p.path(page.num))
		}

		for _, page := range p.PageURLs() {
			Printf("%d %v %s \n", page.Num(), page.IsCurrent(), page.Path())
		}

		Printf("%s \n", p.PageTemp())
	})
}
