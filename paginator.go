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
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Paginator 是分页程序计算的结果
type Paginator struct {
	total int // 总条数

	pageSize    int // 每页条数
	current     int // 当前页页码
	linkedCount int // 可以链接到的页面数量

	pageKey string            // http 请求的 query 页码关键字
	request *http.Request     // 网络请求
	url     string            // http 请求的 url
	params  map[string]string // http 请求的 query 参数集合，不包含 pageKey
}

// Config 默认配置
type Config struct {
	PageSize    int // 每页条数
	Current     int // 当前页页码
	LinkedCount int // 可以链接到的页面数量

	PageKey string        // http query 请求的页码关键字
	Request *http.Request // 网络请求
}

const (
	defaultPageSize    = 10
	defaultCurrent     = 1
	defaultLinkedCount = 3

	defaultPageKey = "page"
)

// Custom 定制默认参数
func Custom(c *Config, total int) *Paginator {
	if c.PageSize <= 0 {
		c.PageSize = defaultPageSize
	}
	p := &Paginator{
		total:       total,
		pageSize:    c.PageSize,
		current:     c.Current,
		linkedCount: c.LinkedCount,

		pageKey: c.PageKey,
		request: nil,
		url:     "",
		params:  map[string]string{},
	}
	if p.current > p.TotalPages() {
		p.current = p.TotalPages()
	}
	if len(p.pageKey) == 0 {
		p.pageKey = defaultPageKey
	}
	return p
}

// New 使用默认参数初始化并返回一个Paginator
func New(total int) *Paginator {
	c := &Config{
		PageSize:    defaultPageSize,
		Current:     defaultCurrent,
		LinkedCount: defaultLinkedCount,
		PageKey:     defaultPageKey,
	}
	return Custom(c, total)
}

// Request 初始化网络请求参数后才可以调用获取指定 page url 等接口
// 如果接口传入了 pageKey query 参数，会根据参数内容更新 current 页码
func (p *Paginator) Request(r *http.Request) *Paginator {
	if r == nil {
		return p
	}
	current := 0
	params := map[string]string{}
	rawQuerys := strings.Split(r.URL.RawQuery, "&")

	for _, query := range rawQuerys {
		param := strings.Split(query, "=")
		if strings.EqualFold(param[0], p.pageKey) {
			current, _ = strconv.Atoi(param[1])
		} else {
			params[param[0]] = param[1]
		}
	}

	p.url = r.URL.Path
	p.params = params
	if current != 0 {
		p.current = current
	}
	return p
}

// IsFirst 如果当前页是第一页返回true
func (p *Paginator) IsFirst() bool {
	return p.current == 1
}

func (p *Paginator) path(pageNum int) string {
	// url 为空认为是没有初始化过 request 方法，不做处理，直接返回
	if len(p.url) == 0 {
		return ""
	}
	params := fmt.Sprintf("%s?", p.url)
	for key, value := range p.params {
		params = fmt.Sprintf("%s%s=%s&", params, key, value)
	}
	return fmt.Sprintf("%s%s=%d", params, p.pageKey, pageNum)
}

// FristURL 返回第一页 url
func (p *Paginator) FristURL() string {
	return p.path(1)
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

// PreviousURL 返回前一页 url
func (p *Paginator) PreviousURL() string {
	return p.path(p.Previous())
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

// NextURL 返回下一页 url
func (p *Paginator) NextURL() string {
	return p.path(p.Next())
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

// LastURL 返回最后一页 url
func (p *Paginator) LastURL() string {
	return p.path(p.TotalPages())
}

// Current 返回当前页页码
func (p *Paginator) Current() int {
	return p.current
}

// CurrentURL 返回当前页 url
func (p *Paginator) CurrentURL() string {
	return p.path(p.Current())
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

	for i := 0; i < maxIdx-offsetIdx+1; i++ {
		pages[i] = &Page{offsetIdx + i, offsetIdx+i == p.current}
	}
	return pages
}

// PageURLs 返回当前页临近几页的页码和 URL 信息
func (p *Paginator) PageURLs() []*PageURL {
	pages := p.Pages()
	pageURLs := make([]*PageURL, len(pages))
	for i := 0; i < len(pages); i++ {
		pageURLs[i] = &PageURL{
			Page:    pages[i],
			pageKey: p.pageKey,
			request: p.request,
			url:     p.url,
			params:  p.params,
		}
	}
	return pageURLs
}

// PageURL 当前页的数据内容
type PageURL struct {
	*Page

	pageKey string            // http 请求的 query 页码关键字
	request *http.Request     // 网络请求
	url     string            // http 请求的 url
	params  map[string]string // http 请求的 query 参数集合，不包含 pageKey
}

// Num 页码
func (p *PageURL) Num() int {
	return p.num
}

// IsCurrent 是否是当前页
func (p *PageURL) IsCurrent() bool {
	return p.isCurrent
}

// Path 当前网页路径地址
func (p *PageURL) Path() string {
	// url 为空认为是没有初始化过 request 方法，不做处理，直接返回
	if len(p.url) == 0 {
		return ""
	}
	params := fmt.Sprintf("%s?", p.url)
	for key, value := range p.params {
		params = fmt.Sprintf("%s%s=%s&", params, key, value)
	}
	return fmt.Sprintf("%s%s=%d", params, p.pageKey, p.num)
}
