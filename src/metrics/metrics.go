package metrics

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math"
	"net/http"
	"os"
	// "os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"tomr/models"
	"tomr/src/utils"

	"github.com/shurcooL/githubv4"
	"github.com/hhatto/gocloc"
)

// var compatibleLicenses = [...]string{"MIT", "LGPLv2.1", "Expat", "X11", "MPL-2.0", "Mozilla Public", "Artistic License 2", "GPLv2", "GPLv3"}
const RegexLicense = `(?i)MIT|LGPLv2.1|Expat|X11|MPL-2.0|Mozilla Public|Artistic License 2|GPLv2|GPLv3`

var ReadMePath = ""

func getBusFactor(jsonRes map[string]interface{}) float64 {

	var disabled float32
	var forking float32
	var visibility float32

	// Collected data from the "web_commit_signoff_required" aspect
	if jsonRes["web_commit_signoff_required"].(bool) {
		disabled = .0
	} else {
		disabled = 0.2
	}

	// Collected data from the "allow_forking" aspect
	if jsonRes["allow_forking"].(bool) {
		forking = 0.2
	} else {
		forking = 0.4
	}

	// Collected data from the "visibility" aspect
	if jsonRes["visibility"].(string) == "public" {
		visibility = .4
	} else {
		visibility = .2
	}

	// Returning weighted sum
	return float64(disabled + forking + visibility)
}

func getCorrectness(jsonRes map[string]interface{}) float64 {

	ownerType := 0.0
	webCommit := 0.0

	// Collecting data from API
	stargazers := jsonRes["stargazers_count"].(float64)
	forksNum := jsonRes["forks_count"].(float64)

	// Analysis of owner type
	owner_map := jsonRes["owner"].(map[string]interface{})
	if owner_map["type"].(string) == "Organization" {
		ownerType = .15
	} else {
		ownerType = .07
	}

	// Analysis of web_commit_signoff_required
	if jsonRes["web_commit_signoff_required"].(bool) {
		webCommit = .1
	} else {
		webCommit = 0.05
	}

	// Assigning weights to stargazers
	if stargazers >= 10000 {
		stargazers = 0.25
	} else if stargazers >= 1000 {
		stargazers = 0.2
	} else if stargazers >= 500 {
		stargazers = 0.15
	} else {
		stargazers = 0.05
	}

	// Assigning weights to forks
	if forksNum >= 10000 {
		forksNum = 0.35
	} else if forksNum >= 1000 {
		forksNum = 0.3
	} else if forksNum >= 100 {
		forksNum = 0.2
	} else if forksNum >= 50 {
		forksNum = 0.15
	} else if forksNum >= 25 {
		forksNum = 0.1
	} else {
		forksNum = 0.05
	}

	total := math.Max(0.1, ownerType+webCommit+stargazers+forksNum)

	return float64(total)
}

func getTotalCommentsGraphQL(jsonRes map[string]interface{}, client *http.Client) int {
	owner_map := jsonRes["owner"].(map[string]interface{})

	var Data struct {
		Viewer struct {
			Login string
		}
		Repository struct {
			CommitComments struct {
				TotalCount int
			}
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner": githubv4.String(owner_map["login"].(string)),
		"name":  githubv4.String(jsonRes["name"].(string)),
	}

	graphQLClient := githubv4.NewClient(client)
	err := graphQLClient.Query(context.Background(), &Data, variables)
	if err != nil {
		Data.Repository.CommitComments.TotalCount = 0
	}
	return Data.Repository.CommitComments.TotalCount
}

func getRampUp(jsonRes map[string]interface{}, client *http.Client) float64 {
	wiki := 0.25
	pages := 0.25
	discussions := 0.0

	// Collecting pertinent data from GITHUB API
	if jsonRes["has_wiki"].(bool) {
		wiki = .45
	}

	if jsonRes["has_pages"].(bool) {
		pages = .5
	}

	if jsonRes["has_discussions"].(bool) {
		discussions = .45
	}

	//fmt.Print("sum", wiki+pages+discussions)

	var commentsScore float32
	totalComments := getTotalCommentsGraphQL(jsonRes, client)

	// Socring comments count based on different ranges of comments
	if totalComments >= 0 && totalComments <= 10 {
		commentsScore = 0.3
	} else if totalComments <= 50 {
		commentsScore = 0.3
	} else if totalComments <= 100 {
		commentsScore = 0.35
	} else if commentsScore <= 400 {
		commentsScore = 0.4
	} else {
		commentsScore = 0.45
	}

	// Returning weighted sum of aspects
	return math.Min(float64(wiki+pages+discussions+float64(commentsScore)), 1.0)
}

func getResponsiveMaintainer(jsonRes map[string]interface{}) float64 {

	var private float32

	// Getting information for last update
	updatedAt := jsonRes["updated_at"].(string)
	if jsonRes["private"].(bool) {
		private = .1
	} else {
		private = .05
	}

	// Parsing the update data
	updateDateList := strings.Split(updatedAt, "-")
	yearStr := updateDateList[0]
	monthStr := updateDateList[1]

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		panic(err)
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		panic(err)
	}
	monthObj := time.Month(month)

	// Arbitrarily taken from the 15 of the month
	t1 := time.Date(year, monthObj, 15, 0, 0, 0, 0, time.UTC)
	t2 := time.Now()
	diff := t2.Sub(t1)

	var updatedLast float32

	// Scoring the update data based on time ranges
	if 0 < diff.Seconds() && diff.Seconds() <= 604800 { // 7 days timeline
		updatedLast = .45
	} else if diff.Seconds() <= 15720000 { // 1/2 a year timeline
		updatedLast = 0.3
	} else if diff.Seconds() <= 15720000*2 { // 1 year timeline
		updatedLast = 0.2
	} else if diff.Seconds() <= 15720000*2*2 { //2 years timeline
		updatedLast = 0.1
	} else {
		updatedLast = 0
	}

	// Acquring additional data from GITHUB API
	hasIssues := jsonRes["has_issues"].(bool)

	openIssues := jsonRes["open_issues"].(float64)

	issuesScore := 0.1

	if hasIssues {
		issuesScore = 0.35 * math.Min(1, openIssues/350)
	}

	archivedStatus := jsonRes["archived"].(bool)
	archivedScore := 0.05

	if !archivedStatus {
		archivedScore = 0.2
	}

	// Returning weighted sum of aspects
	totalValue := math.Min(float64(private+updatedLast+float32(issuesScore)+float32(archivedScore)), 1.0)
	return totalValue
}

func checkLicense(readMe []byte) float64 {
	licenseCompatibility := 0.0

	selected := ""

	// fmt.Print("fist" , firstIdx, "last", lastIdx)
	firstIdx := bytes.Index(readMe, []byte("license"))

	if firstIdx == -1 {
		return 0
	}
	firstIdx -= 15

	lastIdx := len(readMe) - 1

	if (lastIdx - firstIdx) >= 200 {
		selected = string((readMe)[firstIdx : firstIdx+200])
		//fmt.Printf(string(selected))
	} else {
		selected = string((readMe)[firstIdx:lastIdx])
		//fmt.Printf(string(selected))
	}
	matched, err := regexp.MatchString(RegexLicense, selected)
	if err != nil {
		fmt.Println("Error matching regex in checkLicense")
		return licenseCompatibility
	}
	if matched {
		licenseCompatibility = 1.0
	}
	return float64(licenseCompatibility)
}

func checkLicenseFromReadMe(directory string) float64 {
	fmt.Printf("checkLicenseFromReadMe(%s)", directory)
	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		maxDepth := 5
		if err != nil {
			return err
		}
		if d.IsDir() && strings.Count(path, string(os.PathSeparator)) > maxDepth {

			return fs.SkipDir
		} else {
			matched, _ := regexp.MatchString(`(?i)readme`, path)
			if matched {
				if filepath.Ext(path) == ".md" {
					ReadMePath = path
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Print(err)
		return 0.0
	}
	//fmt.Print(ReadMePath)
	readMe, err := os.ReadFile(ReadMePath)
	if err != nil {
		log.Fatal(err)
		return 0.0
	}

	return checkLicense(readMe)
}

func getLicenseScore(license string, pkgDir string, readMe *[]byte) float64 {
	// fmt.Printf("getLicenseScore(%s, %s, nil = %t)\n", license, pkgDir, readMe == nil)
	licenseCompatibility := 0.0
	if license != "" {
		matched, err := regexp.MatchString(RegexLicense, license)
		if err != nil {
			fmt.Println("Error matchine regex in getLicenseScore")
		} else if matched {
			licenseCompatibility = 1.0
		}
	} else if readMe != nil {
		licenseCompatibility = checkLicense(*readMe)
	} else {
		// fmt.Print("here 1")
		pkgJsonFile, err := os.Open(path.Join(pkgDir, "package.json"))
		if err != nil {
			fmt.Printf("Can't find package.json in %s\n", pkgDir)
			return 0
		}
		// fmt.Print("here 2")
		dec := json.NewDecoder(pkgJsonFile)
		// fmt.Print("here 3")
		type elem struct {
			License string `json:"license"`
		}
		var jsonElem elem
		for {
			if err := dec.Decode(&jsonElem); err == io.EOF {
				break
			}
		}
		pkgJsonFile.Close()
		//fmt.Println("license: ", jsonElem.License)
		matched, _ := regexp.MatchString(RegexLicense, jsonElem.License)
		if matched {
			licenseCompatibility = 1.0
		}
		if licenseCompatibility == 0 {
			licenseCompatibility = checkLicenseFromReadMe(pkgDir)
		}
	}
	//fmt.Print(licenseCompatibility)
	return licenseCompatibility
}

func getGoodPinningPractices(url string, client *http.Client) float64 {
	// Get the owner and name of repo
	var depUrl string
	var resp map[string]interface{}
	var pinned float64

	// Getting response from dependency sbom url
	depUrl = url + "/dependency-graph/sbom"
	resp, err := utils.GetDataFromGithub(client, depUrl)
	if err != nil {
		log.Fatalf("Failed to compute Pinning Practices Score for %s\n", url)
		return -1
	}

	// Getting list of packages
	packages := resp["sbom"].(map[string]interface{})["packages"].([]interface{})

	// Iterating through and counting major + minor pinned packages
	pinned = 0
	for i := 0; i < len(packages); i++ {
		// fmt.Println(packages[i].(map[string]interface{})["versionInfo"].(string))
		version := packages[i].(map[string]interface{})["versionInfo"].(string)

		if (!strings.Contains(version, "^")) && (strings.Count(version, ".") >= 2) || strings.Contains(version, "~") {
			pinned = pinned + 1
		}
	}

	pinned = pinned / float64(len(packages))

	if pinned > 1.0 {
		pinned = 1.0
	}
	if pinned < 0.0 || math.IsNaN(pinned) {
		pinned = 0.0
	}

	return pinned
}

func getTotalLines(directory string) int {
	languages := gocloc.NewDefinedLanguages()
	options := gocloc.NewClocOptions()
	paths := []string{
		"src/metrics/temp",
	}

	processor := gocloc.NewProcessor(languages, options)
	result, err := processor.Analyze(paths)
	if err != nil {
		fmt.Printf("gocloc fail. error: %v\n", err)
		return -1
	}

	// for _ , item := range result.Files {
	// 	fmt.Println(item)
	// }
	// fmt.Println(result.Files)
	// fmt.Printf("%+v", result.Total.Code)
	// total := -1
	// cmd := exec.Command("go run gocloc-master/cmd/gocloc/main.go", directory)

	// var stdout, stderr bytes.Buffer
	// cmd.Stdout = &stdout
	// cmd.Stderr = &stderr

	// err = cmd.Run()
	// if err != nil {
	// 	log.Fatal(err)
	// 	log.Fatal("Failed to run cloc command")
	// 	return total
	// } else if stderr.Len() != 0 {
	// 	log.Fatalf("Stderr: %s", stderr.String())
	// }

	// _, clocOut, found := strings.Cut(stdout.String(), "TOTAL")
	// if found {
	// 	re := regexp.MustCompile(`[0-9]+`)
	// 	lines := re.FindAllString(clocOut, -1)
	// 	blank, _ := strconv.Atoi(lines[1])
	// 	comment, _ := strconv.Atoi(lines[2])
	// 	code, _ := strconv.Atoi(lines[3])
	// 	total = blank + comment + code
	// }
	// return total
	fmt.Println(result.Total.Code + result.Total.Comments + result.Total.Blanks + result.Total.Total)
	return int(result.Total.Code + result.Total.Comments + result.Total.Blanks + result.Total.Total)
}

func getGoodEngineeringProcess(url string, client *http.Client, pkgDir string) float64 {
	var resp []map[string]interface{}
	var resp2 map[string]interface{}
	var err error

	// Getting URL of closed PRs
	prUrl := url + "/pulls?state=closed"

	// Getting array of closed PRs
	resp = utils.GetPRs(client, prUrl)

	// Iterating through closed PRs and adding lines introduced - deleted
	sum := 0.0
	for i := 0; i < len(resp); i++ {
		link := resp[i]["_links"].(map[string]interface{})["self"].(map[string]interface{})["href"].(string)
		resp2, err = utils.GetDataFromGithub(client, string(link))
		if err != nil {
			log.Fatalf("Failed to compute score for Engineering Process for %s\n", url)
			return -1
		}
		sum = sum + resp2["additions"].(float64)
		sum = sum - resp2["deletions"].(float64)
	}

	// Getting the total lines of code from original link
	resp2, err = utils.GetDataFromGithub(client, url)
	if err != nil {
		log.Fatalf("Failed to compute score for Engineering Process for %s\n", url)
		return -1
	}

	//fmt.Println(pkgDir)

	total := getTotalLines(pkgDir)

	sum = sum / float64(total)

	if sum > 1.0 {
		sum = 1
	}
	if sum < 0.0 {
		sum = 0.0
	}

	return sum
}

func getNetScore(r models.PackageRating) float64 {
	netScore := ((40 * r.Correctness) + (35 * r.BusFactor) + (30 * r.ResponsiveMaintainer) + (30 * r.RampUp) +
		(25 * r.LicenseScore) + (15 * r.GoodEngineeringProcess) + (10 * r.GoodPinningPractice)) / 185
	return netScore
}
