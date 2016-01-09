package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type redditResponseJSON struct {
	Data struct {
		After    string      `json:"after"`
		Before   interface{} `json:"before"`
		Children []struct {
			Data struct {
				ApprovedBy          interface{}   `json:"approved_by"`
				Archived            bool          `json:"archived"`
				Author              string        `json:"author"`
				AuthorFlairCSSClass interface{}   `json:"author_flair_css_class"`
				AuthorFlairText     interface{}   `json:"author_flair_text"`
				BannedBy            interface{}   `json:"banned_by"`
				Clicked             bool          `json:"clicked"`
				Created             int           `json:"created"`
				CreatedUtc          int           `json:"created_utc"`
				Distinguished       string        `json:"distinguished"`
				Domain              string        `json:"domain"`
				Downs               int           `json:"downs"`
				Edited              int           `json:"edited"`
				From                interface{}   `json:"from"`
				FromID              interface{}   `json:"from_id"`
				FromKind            interface{}   `json:"from_kind"`
				Gilded              int           `json:"gilded"`
				Hidden              bool          `json:"hidden"`
				HideScore           bool          `json:"hide_score"`
				ID                  string        `json:"id"`
				IsSelf              bool          `json:"is_self"`
				Likes               interface{}   `json:"likes"`
				LinkFlairCSSClass   string        `json:"link_flair_css_class"`
				LinkFlairText       string        `json:"link_flair_text"`
				Locked              bool          `json:"locked"`
				Media               interface{}   `json:"media"`
				MediaEmbed          struct{}      `json:"media_embed"`
				ModReports          []interface{} `json:"mod_reports"`
				Name                string        `json:"name"`
				NumComments         int           `json:"num_comments"`
				NumReports          interface{}   `json:"num_reports"`
				Over18              bool          `json:"over_18"`
				Permalink           string        `json:"permalink"`
				Quarantine          bool          `json:"quarantine"`
				RemovalReason       interface{}   `json:"removal_reason"`
				ReportReasons       interface{}   `json:"report_reasons"`
				Saved               bool          `json:"saved"`
				Score               int           `json:"score"`
				SecureMedia         interface{}   `json:"secure_media"`
				SecureMediaEmbed    struct{}      `json:"secure_media_embed"`
				Selftext            string        `json:"selftext"`
				SelftextHTML        string        `json:"selftext_html"`
				Stickied            bool          `json:"stickied"`
				Subreddit           string        `json:"subreddit"`
				SubredditID         string        `json:"subreddit_id"`
				SuggestedSort       interface{}   `json:"suggested_sort"`
				Thumbnail           string        `json:"thumbnail"`
				Title               string        `json:"title"`
				Ups                 int           `json:"ups"`
				URL                 string        `json:"url"`
				UserReports         []interface{} `json:"user_reports"`
				Visited             bool          `json:"visited"`
			} `json:"data"`
			Kind string `json:"kind"`
		} `json:"children"`
		Modhash string `json:"modhash"`
	} `json:"data"`
	Kind string `json:"kind"`
}

func downloadFromURL(url string) string {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return ""
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return ""
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return ""
	}
	fmt.Println(n, "bytes downloaded.")

	return fileName
}

func checkIfImage(url string) bool {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]

	fmt.Println("downloading : " + fileName + "...")

	ext := []string{".jpeg", ".jpg", ".png"}

	for i := 0; i < len(ext); i++ {
		if strings.Contains(fileName, ext[i]) {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("Enter the name of subreddit:")
	var subreddit string
	fmt.Scanf("%s", &subreddit)

	var limit int
	fmt.Println("Enter number of links to get:")
	fmt.Scanf("%d", &limit)

	var subredditLink string
	subredditLink = "http://www.reddit.com/r"
	subredditLink += "/" + subreddit

	subredditJSONlink := subredditLink + "/.json"
	var lastPost string
	count := 0

	var startTime time.Time
	startTime = time.Now()

	for count < limit {

		fmt.Println(subredditJSONlink)

		r, e := http.Get(subredditJSONlink)
		if e != nil {
			fmt.Println("Error:", e)
		}
		defer r.Body.Close()

		var data redditResponseJSON
		dec := json.NewDecoder(r.Body)
		dec.Decode(&data)

		for _, child := range data.Data.Children {
			t := child.Data
			fmt.Println(t.Author, "   "+t.Name, "     Over 18: "+strconv.FormatBool(t.Over18)+"         ", t.URL)
			lastPost = t.Name

			//Download the image now
			var fileName string
			if checkIfImage(t.URL) {
				fileName = downloadFromURL(t.URL)
			} else {
				continue
			}

			//upload to Cloudinary
			UploadCmd := exec.Command("bash", "-c", "cloudinary -i "+fileName+" up settings.conf")
			Out, err := UploadCmd.Output()
			if err != nil {
				panic(err)
			}

			fmt.Println("\nUploaded to Cloudinary...")

			Cmd := exec.Command("bash", "-c", "cloudinary -i "+fileName+" url settings.conf")
			Link, err := Cmd.Output()
			if err != nil {
				panic(err)
			}
			fmt.Println("Link:", string(Link))
			fmt.Println(string(Out))

		}

		count += 25
		subredditJSONlink = subredditLink + "/.json" + "?count=" + strconv.FormatInt(int64(count), 10) + "&after=" + lastPost

	}

	fmt.Println("")
	fmt.Println("Time taken to download : ", time.Since(startTime))
}
