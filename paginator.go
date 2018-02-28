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

// Paginator 是分页程序计算的结果
type Paginator struct {
	total int // 总条数

	pageSize    int // 每页条数
	current     int // 当前页页码
	linkedCount int // 可以链接到的页面数量
}

// Config 默认配置
type Config struct {
	PageSize    int // 每页条数
	Current     int // 当前页页码
	LinkedCount int // 可以链接到的页面数量
}

const (
	defaultPageSize    = 10
	defaultCurrent     = 1
	defaultLinkedCount = 3
)

// Custom 定制默认参数
func Custom(c *Config, total int) *Paginator {
	if c.PageSize <= 0 {
		c.PageSize = defaultPageSize
	}
	p := &Paginator{total, c.PageSize, c.Current, c.LinkedCount}
	if p.current > p.TotalPages() {
		p.current = p.TotalPages()
	}
	return p
}

// New 使用默认参数初始化并返回一个Paginator
func New(total int) *Paginator {
	c := &Config{
		PageSize:    defaultPageSize,
		Current:     defaultCurrent,
		LinkedCount: defaultLinkedCount,
	}
	return Custom(c, total)
}

// IsFirst 如果当前页是第一页返回true
func (p *Paginator) IsFirst() bool {
	return p.current == 1
}

// HasPrevious 如果当前页存在前一页则返回true
func (p *Paginator) HasPrevious() bool {
	return p.current > 1
}

// Previous 如果存在前一页返回前一页页码，否则返回当前页码
func (p *Paginator) Previous() int {
	if !p.HasPrevious() {
		return p.current
	}
	return p.current - 1
}

// HasNext 如果当前页存在后一页则返回true
func (p *Paginator) HasNext() bool {
	return p.total > p.current*p.pageSize
}

// Next 如果存在后一页返回前一页页码，否则返回当前页码
func (p *Paginator) Next() int {
	if !p.HasNext() {
		return p.current
	}
	return p.current + 1
}

// IsLast 如果当前页是最后一页返回true
func (p *Paginator) IsLast() bool {
	if p.total == 0 {
		return true
	}
	return p.total > (p.current-1)*p.pageSize && !p.HasNext()
}

// Total 返回数据总量
func (p *Paginator) Total() int {
	return p.total
}

// TotalPages 返回最后一页页码
func (p *Paginator) TotalPages() int {
	if p.total == 0 {
		return 1
	}
	if p.total%p.pageSize == 0 {
		return p.total / p.pageSize
	}
	return p.total/p.pageSize + 1
}

// Current 返回当前页页码
func (p *Paginator) Current() int {
	return p.current
}

// PageSize 返回每页条数
func (p *Paginator) PageSize() int {
	return p.pageSize
}

// Page 当前页的数据内容
type Page struct {
	num       int
	isCurrent bool
}

// Num 页码
func (p *Page) Num() int {
	return p.num
}

// IsCurrent 是否是当前页
func (p *Page) IsCurrent() bool {
	return p.isCurrent
}

func getMiddleIdx(linkedCount int) int {
	if linkedCount%2 == 0 {
		return linkedCount / 2
	}
	return linkedCount/2 + 1
}

// Pages 返回当前页临近的几页的页码信息
func (p *Paginator) Pages() []*Page {
	if p.linkedCount <= 0 {
		return []*Page{}
	}

	// 只返回当前页
	if p.linkedCount == 1 || p.TotalPages() == 1 {
		return []*Page{{p.current, true}}
	}

	// 总条数小于要返回的条数，就返回全部页码信息
	if p.TotalPages() <= p.linkedCount {
		pages := make([]*Page, p.TotalPages())
		for i := range pages {
			pages[i] = &Page{i + 1, i+1 == p.current}
		}
		return pages
	}

	linkedRadius := p.linkedCount / 2

	// 如果 linkedCount 是奇数， current 前后页数相等
	previousCount, nextCount := linkedRadius, linkedRadius
	// 如果 linkedCount 是偶数， current 前面页数比后面多1
	if p.linkedCount%2 == 0 {
		nextCount--
	}

	// 如果 current<=previousCount 那么需要的页面就是从1到 linkedCount
	// 如果 current>previousCount 并且 current>=TotalPages-nextCount 那么需要的页面就是从 TotalPages-linkedCount+1 到 TotalPages
	// 其余情况就是从 current-previousCount 到 current+nextCount

	pages := make([]*Page, p.linkedCount)
	offsetIdx, maxIdx := 1, 1
	if p.current <= previousCount {
		offsetIdx = 1
		maxIdx = p.linkedCount
	} else if p.current > previousCount && p.current >= p.TotalPages()-nextCount {
		offsetIdx = p.TotalPages() - p.linkedCount + 1
		maxIdx = p.TotalPages()
	} else {
		offsetIdx = p.current - previousCount
		maxIdx = p.current + nextCount
	}

	for i := 0; i < maxIdx - offsetIdx + 1; i++ {
		pages[i] = &Page{offsetIdx + i, offsetIdx + i == p.current}
	}
	return pages
}
