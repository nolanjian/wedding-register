package request

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/CodisLabs/codis/pkg/utils/log"
	"github.com/PuerkitoBio/goquery"
	"github.com/mahonia"
)

type WebProxy struct {
	cli    *http.Client
	cookie []*http.Cookie
	header http.Header
	thisID string
	time   *time.Time
}

func GetWebProxy() *WebProxy {
	return &WebProxy{
		cli:    &http.Client{},
		cookie: make([]*http.Cookie, 0),
		header: make(http.Header),
	}
}

func (p *WebProxy) SetData(date string) error {
	myTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Error(err)
		return err
	}
	p.time = &myTime
	log.Debug(p.time)
	return nil
}

func (p *WebProxy) getSelDate() string {
	return p.time.Format("2006-01-02")
}

func (p *WebProxy) Excute() error {
	if p == nil {
		err := errors.New("empty pointer")
		log.Error(err)
		return err
	}

	if err := p.FirstRequest(); err != nil {
		log.Error(err)
		return err
	}

	if err := p.SecondRequest(); err != nil {
		log.Error(err)
		return err
	}

	if err := p.ThirdRequest(); err != nil {
		log.Error(err)
		return err
	}

	if err := p.ForthRequest(); err != nil {
		log.Error(err)
		return err
	}

	// if err := p.FifthRequest(); err != nil {
	// 	log.Error(err)
	// 	return err
	// }

	return nil
}

func (p *WebProxy) getCharset(resp *http.Response) string {
	for _, item := range resp.Header["Content-Type"] {
		strTmp := strings.ToLower(item)
		if index := strings.Index(strTmp, `charset=`); index >= 0 {
			strRet := strTmp[index+len(`charset=`):]
			return strings.Trim(strRet, "\t\n ")
		}

	}
	return "utf-8"
}

func (p *WebProxy) storageCookie(resp *http.Response) {
	if _, exist := resp.Header["Set-Cookie"]; exist {
		log.Debug(resp.Cookies())
		log.Debug(p.cookie)
		p.cookie = resp.Cookies()
		log.Debug(p.cookie)
	}
}

func (p *WebProxy) request(method, reqURL string, form url.Values) (string, error) {
	// url Parse
	pURL, err := url.Parse(reqURL)
	if err != nil {
		log.Error(err)
		return "", err
	}

	body := form.Encode()

	// new request
	req := &http.Request{
		Method: strings.ToUpper(method),
		URL:    pURL,
	}

	if form != nil {
		req.Body = ioutil.NopCloser(strings.NewReader(body))
		req.ContentLength = int64(len(body))
		log.Debug(req.ContentLength)
	}

	// set header
	req.Header = p.header

	// restorage cookie
	log.Debug(p.cookie)
	for ii := range p.cookie {
		req.AddCookie(p.cookie[ii])
	}

	// do request
	resp, err := p.cli.Do(req)
	if err != nil {
		log.Error(err)
		return "", err
	}

	// save cookie
	p.storageCookie(resp)

	// set Referer
	//p.header["Referer"] = []string{resp.Request.URL.String()}
	//log.Debug(p.header["Referer"])

	// read page
	charset := p.getCharset(resp)
	dec := mahonia.NewDecoder(charset)
	rd := dec.NewReader(resp.Body)
	data, err := ioutil.ReadAll(rd)
	if err != nil {
		log.Error(err)
		return "", err
	}
	html := string(data)
	log.Debug(html)
	return html, nil
}

func (p *WebProxy) Get(reqURL string) (string, error) {
	return p.request("GET", reqURL, nil)
}

func (p *WebProxy) Post(reqURL string, form url.Values) (string, error) {
	return p.request("POST", reqURL, form)
}

// # FIRST request
//
// GET http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj_01.jsp?id=1 HTTP/1.1
// Host: wsbs.gzmz.gov.cn
// Connection: keep-alive
// Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
// Upgrade-Insecure-Requests: 1
// User-Agent: Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.221 Safari/537.36 SE 2.X MetaSr 1.0
// Accept-Encoding: gzip, deflate, sdch
// Accept-Language: zh-CN,zh;q=0.8
// Cookie: JSESSIONID=43633C0AA40BD42EA65B1715C1CDB4CC; UM_distinctid=15abe2878bf10-00755db2664684-637e7415-100200-15abe2878c4a

func (p *WebProxy) FirstRequest() error {
	html, err := p.Get(`http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj_01.jsp?id=1`)
	if err != nil {
		return err
	}
	log.Info(html)

	return nil
}

// # SECOND request
//
// POST http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj_02.jsp HTTP/1.1
// Host: wsbs.gzmz.gov.cn
// Connection: keep-alive
// Content-Length: 4
// Cache-Control: max-age=0
// Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
// Origin: http://wsbs.gzmz.gov.cn
// Upgrade-Insecure-Requests: 1
// User-Agent: Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.221 Safari/537.36 SE 2.X MetaSr 1.0
// Content-Type: application/x-www-form-urlencoded
// Referer: http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj_01.jsp?id=1
// Accept-Encoding: gzip, deflate
// Accept-Language: zh-CN,zh;q=0.8
// Cookie: JSESSIONID=43633C0AA40BD42EA65B1715C1CDB4CC; UM_distinctid=15abe2878bf10-00755db2664684-637e7415-100200-15abe2878c4a
//
// nd=1

func (p *WebProxy) SecondRequest() error {
	form := make(url.Values)
	form["nd"] = []string{url.QueryEscape("1")}
	html, err := p.Post(`http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj_02.jsp`, form)
	if err != nil {
		return err
	}
	log.Info(html)

	return nil
}

// # THIRD request

// # 3.1
//
// GET http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj02_return.jsp?class_code=440000&info=mshi HTTP/1.1
// Host: wsbs.gzmz.gov.cn
// Connection: keep-alive
// Accept: */*
// X-Requested-With: XMLHttpRequest
// User-Agent: Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.221 Safari/537.36 SE 2.X MetaSr 1.0
// Referer: http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj_02.jsp
// Accept-Encoding: gzip, deflate, sdch
// Accept-Language: zh-CN,zh;q=0.8
// Cookie: JSESSIONID=43633C0AA40BD42EA65B1715C1CDB4CC; UM_distinctid=15abe2878bf10-00755db2664684-637e7415-100200-15abe2878c4a

// # 3.2
//
// GET http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj02_return.jsp?class_code=440000&info=wshi HTTP/1.1
// Host: wsbs.gzmz.gov.cn
// Connection: keep-alive
// Accept: */*
// X-Requested-With: XMLHttpRequest
// User-Agent: Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.221 Safari/537.36 SE 2.X MetaSr 1.0
// Referer: http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj_02.jsp
// Accept-Encoding: gzip, deflate, sdch
// Accept-Language: zh-CN,zh;q=0.8
// Cookie: JSESSIONID=43633C0AA40BD42EA65B1715C1CDB4CC; UM_distinctid=15abe2878bf10-00755db2664684-637e7415-100200-15abe2878c4a

// # 3.3
//
// GET http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj02_return.jsp?class_code=440000&info=wshi HTTP/1.1
// Host: wsbs.gzmz.gov.cn
// Connection: keep-alive
// Accept: */*
// X-Requested-With: XMLHttpRequest
// User-Agent: Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.221 Safari/537.36 SE 2.X MetaSr 1.0
// Referer: http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj_02.jsp
// Accept-Encoding: gzip, deflate, sdch
// Accept-Language: zh-CN,zh;q=0.8
// Cookie: JSESSIONID=43633C0AA40BD42EA65B1715C1CDB4CC; UM_distinctid=15abe2878bf10-00755db2664684-637e7415-100200-15abe2878c4a

// # 3.4
//
// GET http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj02_return.jsp?class_code=440100&info=wqu HTTP/1.1
// Host: wsbs.gzmz.gov.cn
// Connection: keep-alive
// Accept: */*
// X-Requested-With: XMLHttpRequest
// User-Agent: Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.221 Safari/537.36 SE 2.X MetaSr 1.0
// Referer: http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj_02.jsp
// Accept-Encoding: gzip, deflate, sdch
// Accept-Language: zh-CN,zh;q=0.8
// Cookie: JSESSIONID=43633C0AA40BD42EA65B1715C1CDB4CC; UM_distinctid=15abe2878bf10-00755db2664684-637e7415-100200-15abe2878c4a

// # 3.5
//
// POST http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj_03.jsp HTTP/1.1
// Host: wsbs.gzmz.gov.cn
// Connection: keep-alive
// Content-Length: 252
// Cache-Control: max-age=0
// Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
// Origin: http://wsbs.gzmz.gov.cn
// Upgrade-Insecure-Requests: 1
// User-Agent: Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.221 Safari/537.36 SE 2.X MetaSr 1.0
// Content-Type: application/x-www-form-urlencoded
// Referer: http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj_02.jsp
// Accept-Encoding: gzip, deflate
// Accept-Language: zh-CN,zh;q=0.8
// Cookie: JSESSIONID=43633C0AA40BD42EA65B1715C1CDB4CC; UM_distinctid=15abe2878bf10-00755db2664684-637e7415-100200-15abe2878c4a
//
// man=%E5%B9%BF%E4%B8%9C%E7%9C%81%E5%B9%BF%E5%B7%9E%E5%B8%82%E6%B5%B7%E7%8F%A0%E5%8C%BA&nd=1&woman=%E5%B9%BF%E4%B8%9C%E7%9C%81%E5%B9%BF%E5%B7%9E%E5%B8%82%E8%B6%8A%E7%A7%80%E5%8C%BA&msheng=440000&mshi=440100&mqu=440105&wsheng=440000&wshi=440100&wqu=440104
//
//	Form:
//			|	man			|	广东省广州市海珠区		|
//			|	woman		|	广东省广州市越秀区		|
//			|	msheng		|	440000					|
//			|	mshi		|	440100					|
//			|	mqu			|	440105					|
//			|	wsheng		|	440000					|
//			|	wshi		|	440100					|
//			|	wqu			|	440104					|
//			|	nd			|	1						|

func (p *WebProxy) ThirdRequest() error {

	// getURLs := []string{
	// 	`http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj02_return.jsp?class_code=440000&info=mshi`,
	// 	`http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj02_return.jsp?class_code=440100&info=mqu`,
	// 	`http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj02_return.jsp?class_code=440000&info=wshi`,
	// 	`http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj02_return.jsp?class_code=440100&info=wqu`,
	// }

	// for _, url := range getURLs {
	// 	html, err := p.Get(url)
	// 	if err != nil {
	// 		log.Error(err)
	// 		return err
	// 	}
	// 	log.Info(html)
	// }

	form := make(url.Values)
	form["man"] = []string{"广东省广州市海珠区"}
	form["woman"] = []string{"广东省广州市越秀区"}
	form["msheng"] = []string{"440000"}
	form["mshi"] = []string{"440100"}
	form["mqu"] = []string{"440105"}
	form["wsheng"] = []string{"440000"}
	form["wshi"] = []string{"440100"}
	form["wqu"] = []string{"440104"}
	form["nd"] = []string{"1"}

	p.header["Content-Type"] = []string{`application/x-www-form-urlencoded`}
	html, err := p.Post(`http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj_03.jsp`, form)
	if err != nil {
		return err
	}
	log.Info(html)

	return nil
}

// # forth
//
// POST http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj_04.jsp HTTP/1.1
// Host: wsbs.gzmz.gov.cn
// Connection: keep-alive
// Content-Length: 625
// Cache-Control: max-age=0
// Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
// Origin: http://wsbs.gzmz.gov.cn
// Upgrade-Insecure-Requests: 1
// User-Agent: Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.221 Safari/537.36 SE 2.X MetaSr 1.0
// Content-Type: application/x-www-form-urlencoded
// Referer: http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj_03.jsp
// Accept-Encoding: gzip, deflate
// Accept-Language: zh-CN,zh;q=0.8
// Cookie: JSESSIONID=C9E3513BEAA505103642C23589595427; UM_distinctid=15abe2878bf10-00755db2664684-637e7415-100200-15abe2878c4a
//
// seldate=2017-04-15&r_code=4401041&man=%E5%B9%BF%E4%B8%9C%E7%9C%81%E5%B9%BF%E5%B7%9E%E5%B8%82%E6%B5%B7%E7%8F%A0%E5%8C%BA&woman=%E5%B9%BF%E4%B8%9C%E7%9C%81%E5%B9%BF%E5%B7%9E%E5%B8%82%E8%B6%8A%E7%A7%80%E5%8C%BA&mqu=440105&wqu=440104&nd=1&xuanzheyydate=2017%E5%B9%B404%E6%9C%8815%E6%97%A5%C2%A0%C2%A009%3A00-09%3A15&deptname=%E5%B9%BF%E5%B7%9E%E5%B8%82%E8%B6%8A%E7%A7%80%E5%8C%BA%E6%B0%91%E6%94%BF%E5%B1%80%E5%A9%9A%E5%A7%BB%E7%99%BB%E8%AE%B0%E5%A4%84&yytime=09%3A00_2017-04-15_4401041_09%3A15&deptname=%E5%B9%BF%E5%B7%9E%E5%B8%82%E6%B5%B7%E7%8F%A0%E5%8C%BA%E6%B0%91%E6%94%BF%E5%B1%80%E5%A9%9A%E5%A7%BB%E7%99%BB%E8%AE%B0%E5%A4%84

//	|	seldate			|	2017-04-15						|
//	|	r_code			|	4401041							|
//	|	man				|	广东省广州市海珠区				|
//	|	woman			|	广东省广州市越秀区				|
//	|	mqu				|	440105							|
//	|	wqu				|	440104							|
//	|	nd				|	1								|
//	|	xuanzheyydate	|	2017年04月15日  09:00-09:15	|
//	|	deptname		|	广州市越秀区民政局婚姻登记处		|
//	|	yytime			|	09:00_2017-04-15_4401041_09:15	|
//	|	deptname		|	广州市海珠区民政局婚姻登记处		|

func (p *WebProxy) ForthRequest() error {
	form := make(url.Values)
	form["seldate"] = []string{"2017-04-22"}
	form["r_code"] = []string{"4401041"}
	form["man"] = []string{"广东省广州市海珠区"}
	form["woman"] = []string{"广东省广州市越秀区"}
	form["mqu"] = []string{"440105"}
	form["wqu"] = []string{"440104"}
	form["nd"] = []string{"1"}
	form["xuanzheyydate"] = []string{"2017年04月22日  09:00-09:15"}
	form["deptname"] = []string{"广州市越秀区民政局婚姻登记处"}
	form["yytime"] = []string{"09:00_2017-04-22_4401041_09:15"}
	form["deptname"] = []string{"广州市海珠区民政局婚姻登记处"}

	p.header["Content-Type"] = []string{`application/x-www-form-urlencoded`}
	html, err := p.Post(`http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj_04.jsp`, form)
	if err != nil {
		return err
	}
	log.Info(html)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Error(err)
		return err
	}
	thisid := doc.Find("#thisid")
	if val, exist := thisid.Attr("value"); exist {
		p.thisID = val
	} else {
		err := errors.New("thisid not found")
		log.Error(err)
		return err
	}

	return nil
}

// # fifth Request
//
// POST http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj_05.jsp HTTP/1.1
// Host: wsbs.gzmz.gov.cn
// Proxy-Connection: keep-alive
// Content-Length: 539
// Cache-Control: max-age=0
// Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
// Origin: http://wsbs.gzmz.gov.cn
// Upgrade-Insecure-Requests: 1
// User-Agent: Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.221 Safari/537.36 SE 2.X MetaSr 1.0
// Content-Type: application/x-www-form-urlencoded
// Referer: http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj_04.jsp
// Accept-Encoding: gzip, deflate
// Accept-Language: zh-CN,zh;q=0.8
// Cookie: JSESSIONID=C9E3513BEAA505103642C23589595427; UM_distinctid=15abe2878bf10-00755db2664684-637e7415-100200-15abe2878c4a
//
//thisid=d3afe0b5d90a4689a7f1dbea7193560f&seldate=2017-04-21&mhuji=%E5%B9%BF%E4%B8%9C%E7%9C%81%E5%B9%BF%E5%B7%9E%E5%B8%82%E6%B5%B7%E7%8F%A0%E5%8C%BA&whuji=%E5%B9%BF%E4%B8%9C%E7%9C%81%E5%B9%BF%E5%B7%9E%E5%B8%82%E8%B6%8A%E7%A7%80%E5%8C%BA&nd=1&str=09%3A30_2017-04-21_4401041_09%3A45&xuanzheyydate=2017%E5%B9%B404%E6%9C%8821%E6%97%A5%C2%A0%C2%A009%3A30-09%3A45&mname=%E7%AE%80%E5%86%A0%E8%85%BE&mphone=13570506413&mcardtype=0&midcard=440105199206065775&wname=%E7%9B%9B%E7%A5%89%E5%90%9B&wphone=15975471716&wcardtype=0&widcard=440102199303104028
//
// thisid	d3afe0b5d90a4689a7f1dbea7193560f
// seldate	2017-04-21
// mhuji	广东省广州市海珠区
// whuji	广东省广州市越秀区
// nd	1
// str	09:30_2017-04-21_4401041_09:45
// xuanzheyydate	2017年04月21日  09:30-09:45
// mname	简冠腾
// mphone	13570506413
// mcardtype	0
// midcard	440105199206065775
// wname	盛祉君
// wphone	15975471716
// wcardtype	0
// widcard	440102199303104028

func (p *WebProxy) FifthRequest() error {
	form := make(url.Values)
	form["thisid"] = []string{p.thisID}
	form["seldate"] = []string{p.getSelDate()}
	form["mhuji"] = []string{"广东省广州市海珠区"}
	form["whuji"] = []string{"广东省广州市越秀区"}
	form["nd"] = []string{"1"}
	form["str"] = []string{"09:30_2017-04-21_4401041_09:45"}
	form["xuanzheyydate"] = []string{"2017年04月21日  09:30-09:45"}
	form["mname"] = []string{"简冠腾"}
	form["mphone"] = []string{"13570506413"}
	form["mcardtype"] = []string{"0"}
	form["midcard"] = []string{"440105199206065775"}
	form["wname"] = []string{"盛祉君"}
	form["wphone"] = []string{"15975471716"}
	form["wcardtype"] = []string{"0"}
	form["widcard"] = []string{"440102199303104028"}

	p.header["Content-Type"] = []string{`application/x-www-form-urlencoded`}
	html, err := p.Post(`http://wsbs.gzmz.gov.cn/gsmpro/web/jhdj_04.jsp`, form)
	if err != nil {
		return err
	}
	log.Info(html)

	return nil
}
