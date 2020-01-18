package trending

import (
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	DEFAULT_URL   = "https://github.com"
	TRENDING_PATH = "trending"
	TIME_TODAY = "daily"
	TIME_WEEK = "weekly"
	TIME_MONTH = "monthly"
)

type Repository struct {
	Name        string
	Owner       string
	Description string
	Language    string
	Star        uint
	URL         string
}

type trendingClient struct {
	BaseUrlStr string
	Client     *http.Client
}

type Option func(*trendingClient)

func NewClient(ops ...Option) *trendingClient {
  c := &trendingClient{
		BaseUrlStr: DEFAULT_URL,
		Client:     http.DefaultClient,
	}
  for _ ,op := range ops {
    op(c)
  }
	return c
}

func WithHttpClient(client *http.Client) Option {
  return func(c *trendingClient){
    c.Client = client
  }
}

func WithBaseUrlStr(urlStr string) Option {
  return func(c *trendingClient){
    c.BaseUrlStr = urlStr
  }
}

func (c *trendingClient) GetRepository(time, lang string) ([]Repository, error) {
	url, err := c.genTrendingURL(time, lang)
	if err != nil {
		return nil, err
	}
	res, err := c.Client.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}
	repos := []Repository{}
	doc.Find("article.Box-row").Each(func(i int, s *goquery.Selection) {
		repo := Repository{}
		un, exists := s.Find("h1 > a").First().Attr("href")
		if exists {
			strs := strings.Split(un, "/")
			repo.Name = strs[2]
			repo.Owner = strs[1]
			repo.URL = c.BaseUrlStr + un
		}
    repo.Description = strings.TrimSpace(s.Find("p").First().Text())
		repo.Language = s.Find("div.f6.text-gray.mt-2 > span.d-inline-block.ml-0.mr-3 > span:nth-child(2)").First().Text()
		star, _ := strconv.Atoi(
      strings.TrimSpace(
        strings.Replace(s.Find("div.f6.text-gray.mt-2 >  a:nth-child(2)").First().Text(),",","",-1),
      ))
		repo.Star = uint(star)
		repos = append(repos, repo)
	})
	return repos, nil
}

func (c *trendingClient) genTrendingURL(time, lang string) (*url.URL, error) {
	u, err := url.Parse(c.BaseUrlStr)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, TRENDING_PATH, lang)
	q := u.Query()
	if time != "" {
		q.Set("since", time)
	}
	u.RawQuery = q.Encode()
	return u, nil
}
