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
	Convey("Basic logics", t, func() {
		p := New(20)
		So(len(p.Pages()), ShouldEqual, 2)
	})

	Convey("Custom logics", t, func() {
		p := Custom(&Config{10, 2, 3}, 23)
		So(len(p.Pages()), ShouldEqual, 3)
		So(p.TotalPages(), ShouldEqual,3)
		So(p.IsFirst(), ShouldBeFalse)
		So(p.HasPrevious(), ShouldBeTrue)
		So(p.Previous(), ShouldEqual,1)
		So(p.HasNext(), ShouldBeTrue)
		So(p.Next(), ShouldEqual, 3)
		So(p.IsLast(), ShouldBeFalse)

		Convey("LinkedCount 小于等于0",  func() {
			p := Custom(&Config{10, 2, 0}, 23)
			So(len(p.Pages()), ShouldEqual, 0)

			p = Custom(&Config{10, 2, 0}, 5)
			So(len(p.Pages()), ShouldEqual, 0)
			So(p.Current(), ShouldEqual, 1)
		})

		Convey("LinkedCount 等于1 或者 TotalPages 等于1",  func() {
			p := Custom(&Config{10, 2, 1}, 23)
			So(len(p.Pages()), ShouldEqual, 1)

			p = Custom(&Config{10, 2, 1}, 5)
			So(len(p.Pages()), ShouldEqual, 1)
		})

		Convey("TotalPages 大于 LinkedCount",  func() {
			Convey("LinkedCount 是奇数",  func() {
				Print("\n")
				p := Custom(&Config{10, 1, 3}, 63)
				for _, page := range p.Pages(){
					Printf("%v ", page)
				}
				Print("\n")

				p = Custom(&Config{10, 3, 3}, 63)
				for _, page := range p.Pages(){
					Printf("%v ", page)
				}
				Print("\n")

				p = Custom(&Config{10, 6, 3}, 63)
				for _, page := range p.Pages(){
					Printf("%v ", page)
				}
				Print("\n")
			})

			Convey("LinkedCount 是偶数",  func() {
				Print("\n")
				p := Custom(&Config{10, 1, 4}, 63)
				for _, page := range p.Pages(){
					Printf("%v ", page)
				}
				Print("\n")

				p = Custom(&Config{10, 4, 4}, 63)
				for _, page := range p.Pages(){
					Printf("%v ", page)
				}
				Print("\n")

				p = Custom(&Config{10, 6, 4}, 63)
				for _, page := range p.Pages(){
					Printf("%v ", page)
				}
				Print("\n")
			})
		})
	})
}
