package main

import (
	"os"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	t "github.com/brimstone/go-twitter"
	"github.com/brimstone/logger"
	twitter "github.com/dghubble/go-twitter/twitter"
)

const tpl = `### Twitter Lists
## <a href="https://twitter.com/i/lists/{{ .ID }}">{{ .Name }}</a>
<table>
{{range .Members}}<tr><td><a href="https://twitter.com/{{ .ScreenName }}"><img src="{{ .ProfileImage }}"></a></td><td>
<b><a href="https://twitter.com/{{ .ScreenName }}">@{{ .ScreenName }}</a> ({{ .Name }})</b><br />
<ul>
<li>{{ if .LastTweet }}Last Tweet: {{ .LastTweet }}{{else}}<i>Protected</i>{{end}}</li>
<li>{{ .Description }}</li>
</ul>
</td></tr>
{{end}}
</table>
`

type Member struct {
	Description  string
	ID           int64
	Name         string
	ProfileImage string
	ScreenName   string
	Status       *twitter.Tweet
	LastTweet    string
}

func main() {
	log := logger.New()
	now := time.Now()

	client, _, err := t.NewClient(t.Tokens{
		ConsumerKey:    os.Getenv("CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("CONSUMER_SECRET"),
		AccessToken:    os.Getenv("ACCESS_TOKEN"),
		AccessSecret:   os.Getenv("ACCESS_SECRET"),
	})
	if err != nil {
		panic(err)
	}
	listsResp, _, err := client.Lists.List(nil)
	if err != nil {
		panic(err)
	}

	lists := make(map[string]int64)
	for _, list := range listsResp {
		lists[list.Name] = list.ID
		log.Debug("list",
			log.Field("name", list.Name),
			log.Field("ID", list.ID),
		)
	}

	membersResp, _, err := client.Lists.Members(&twitter.ListsMembersParams{
		ListID: lists["security"],
		Count:  1000,
	})

	log.Debug("members",
		log.Field("count", len(membersResp.Users)),
	)

	members := []Member{}
	for _, member := range membersResp.Users {
		if member.Entities != nil {
			for _, u := range member.Entities.Description.Urls {
				member.Description = strings.ReplaceAll(member.Description, u.URL, u.ExpandedURL)
			}
		}
		m := Member{
			Description:  member.Description,
			ID:           member.ID,
			Name:         member.Name,
			ProfileImage: strings.Replace(member.ProfileImageURLHttps, "_normal", "_200x200", 1),
			ScreenName:   member.ScreenName,
			Status:       member.Status,
		}
		if member.Status != nil {
			c, err := time.Parse("Mon Jan 2 15:04:05 -0700 2006", member.Status.CreatedAt)
			if err == nil {
				thisyear, thisweek := now.ISOWeek()
				thatyear, thatweek := c.ISOWeek()
				// TODO Does this handle year wrap around right?
				if thisyear != thatyear {
					thisweek += 52
				}

				// Years
				if now.Year() != c.Year() {
					m.LastTweet = strconv.Itoa(c.Year())
					// Months
					// TODO might error on first day of month
				} else if now.Month() != c.Month() {
					m.LastTweet = c.Month().String()
					// Weeks
				} else if thisweek-thatweek > 1 {
					m.LastTweet = strconv.Itoa(thisweek-thatweek) + " weeks ago"
				} else if thisweek-thatweek == 1 {
					m.LastTweet = "last week"
					// Days
				} else if now.YearDay()-c.YearDay() > 1 {
					m.LastTweet = strconv.Itoa(now.YearDay()-c.YearDay()) + " days ago"
				} else if now.YearDay()-c.YearDay() == 1 {
					m.LastTweet = "yesterday"
				} else {
					m.LastTweet = "today"
				}
			}
		}
		log.Debug("member",
			log.Field("name", member.Name),
			log.Field("ID", member.ID),
			log.Field("screenname", member.ScreenName),
			//log.Field("profileImage", member.ProfileImageURLHttps),
		)
		members = append(members, m)
	}

	sort.Slice(members, func(i, j int) bool {
		return members[i].ID > members[j].ID
	})

	t, err := template.New("readme").Parse(tpl)
	if err != nil {
		panic(err)
	}

	data := struct {
		ID          int64
		Name        string
		Members     []Member
		LastUpdated time.Time
	}{
		ID:          lists["security"],
		Name:        "security",
		Members:     members,
		LastUpdated: now,
	}

	f, err := os.Create("README.md")
	err = t.Execute(f, data)
	if err != nil {
		panic(err)
	}

}
