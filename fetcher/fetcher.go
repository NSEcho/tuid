package fetcher

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	profileURL = "https://api.twitter.com/1.1/users/show.json"
	tokenEnv   = "TUID_TOKEN"
)

type Fetcher struct {
	client   *http.Client
	ids      []int
	userData map[int]*us
}

func NewFetcher() *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
		userData: make(map[int]*us),
	}
}

type ProfileResp struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

func (f *Fetcher) GetByUsername(username string) (*ProfileResp, error) {
	token, err := getToken()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", profileURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	q := req.URL.Query()
	q.Add("screen_name", username)
	req.URL.RawQuery = q.Encode()

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r ProfileResp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (f *Fetcher) GetByID(id int) (*ProfileResp, error) {
	return f.getProfileFromID(id)
}

type us struct {
	name string
	at   string
}

func (f *Fetcher) Monitor(usersPath string) error {
	users, err := readUsernames(usersPath)
	if err != nil {
		return err
	}

	log.Printf("Read %d users", len(users))

	for _, user := range users {
		u, err := f.GetByUsername(user)
		if err != nil {
			return err
		}
		f.ids = append(f.ids, u.ID)
		f.userData[u.ID] = &us{
			name: u.Name,
			at:   u.ScreenName,
		}
	}

	resch := make(chan *ProfileResp, len(users))

	go func() {
		for u := range resch {
			if u.Name != f.userData[u.ID].name {
				log.Printf("Name changed from %s to %s", f.userData[u.ID].name, u.Name)
				f.userData[u.ID].name = u.Name
			}
			if u.ScreenName != f.userData[u.ID].at {
				log.Printf("Username changed from %s to %s", f.userData[u.ID].at, u.ScreenName)
				f.userData[u.ID].at = u.Name
			}
		}
	}()

	t := time.NewTicker(1 * time.Minute)
	for {
		for _, id := range f.ids {
			go f.check(id, resch)
		}
		<-t.C
	}
}

func (f *Fetcher) check(id int, resch chan<- *ProfileResp) {
	u, err := f.getProfileFromID(id)
	if err != nil {
		return
	}
	resch <- u
}

func (f *Fetcher) getProfileFromID(id int) (*ProfileResp, error) {
	token, err := getToken()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", profileURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	q := req.URL.Query()
	q.Add("user_id", strconv.Itoa(id))
	req.URL.RawQuery = q.Encode()

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r ProfileResp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func readUsernames(userPath string) ([]string, error) {
	f, err := os.Open(userPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	uniqueUsers := make(map[string]struct{})

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		uniqueUsers[line] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var users []string

	for u := range uniqueUsers {
		users = append(users, u)
	}
	return users, nil
}

func getToken() (string, error) {
	token, ok := os.LookupEnv(tokenEnv)
	if !ok {
		return "", errors.New("getToken: TUID_TOKEN not found")
	}
	return token, nil
}
