package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"encoding/json"
	"encoding/base64"
	"net/http"
	"net/url"
	"io/ioutil"
	"bytes"
	"fmt"
	"os"
)

type Service struct {
		Credentials struct {
			Password string `json:"password"`
			URL      string `json:"url"`
			Username string `json:"username"`
		} `json:"credentials"`
		Label string   `json:"label"`
		Name  string   `json:"name"`
		Plan  string   `json:"plan"`
		Tags  []string `json:"tags"`
}

type WatsonResponse struct {
	Question struct {
		Answers []struct {
			Confidence float64 `json:"confidence"`
			ID         float64 `json:"id"`
			Pipeline   string  `json:"pipeline"`
			Text       string  `json:"text"`
		} `json:"answers"`
		Category           string        `json:"category"`
		ErrorNotifications []interface{} `json:"errorNotifications"`
		EvidenceRequest    struct {
			Items   float64 `json:"items"`
			Profile string  `json:"profile"`
		} `json:"evidenceRequest"`
		Evidencelist []struct {
			Copyright   string `json:"copyright"`
			Document    string `json:"document"`
			ID          string `json:"id"`
			MetadataMap struct {
				DOCNO        string `json:"DOCNO"`
				Abstract     string `json:"abstract"`
				CorpusName   string `json:"corpusName"`
				Deepqaid     string `json:"deepqaid"`
				Description  string `json:"description"`
				FileName     string `json:"fileName"`
				Originalfile string `json:"originalfile"`
				Title        string `json:"title"`
			} `json:"metadataMap"`
			TermsOfUse string `json:"termsOfUse"`
			Text       string `json:"text"`
			Title      string `json:"title"`
			Value      string `json:"value"`
		} `json:"evidencelist"`
		Focuslist []struct {
			Value string `json:"value"`
		} `json:"focuslist"`
		FormattedAnswer bool    `json:"formattedAnswer"`
		ID              string  `json:"id"`
		Items           float64 `json:"items"`
		Latlist         []struct {
			Value string `json:"value"`
		} `json:"latlist"`
		Passthru   string `json:"passthru"`
		Pipelineid string `json:"pipelineid"`
		Qclasslist []struct {
			Value string `json:"value"`
		} `json:"qclasslist"`
		QuestionText string `json:"questionText"`
		Status       string `json:"status"`
		SynonymList  []struct {
			Lemma        string `json:"lemma"`
			PartOfSpeech string `json:"partOfSpeech"`
			SynSet       []struct {
				Name    string `json:"name"`
				Synonym []struct {
					IsChosen bool    `json:"isChosen"`
					Value    string  `json:"value"`
					Weight   float64 `json:"weight"`
				} `json:"synonym"`
			} `json:"synSet"`
			Value string `json:"value"`
		} `json:"synonymList"`
	} `json:"question"`
}

type QuestionResponse struct{
	QuestionResponse WatsonResponse
	Error string
}


var(
	service_url = "<service_url>"
	service_username = "<service_username"
	service_password = "<service_password>"
	auth = ""
)


func main() {

	servicesEnv := os.Getenv("VCAP_SERVICES");
   	appInfo := []byte(servicesEnv)

   	if len(appInfo) > 0{
   		var services map[string][]Service
		err := json.Unmarshal(appInfo, &services)
	
		if err != nil {
        	fmt.Printf("error %+v", err.Error())
        	//panic(err)
    	}else{
    		tempArrayServices := services["question_and_answer"];
    		if len(tempArrayServices) > 0{
    			tempCredentials := tempArrayServices[0].Credentials
    			service_url = tempCredentials.URL
    			service_username = tempCredentials.Username
    			service_password = tempCredentials.Password
    		}
   		}
    
    	authCredentials := fmt.Sprintf("%s:%s", service_username, service_password)
    	auth = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(authCredentials)))
   	}
    

	m := martini.Classic()

	StaticOptions := martini.StaticOptions{Prefix: "public"}
	m.Use(martini.Static("public", StaticOptions))

	m.Use(render.Renderer(render.Options{
		Directory:  "templates",       // Specify what path to load the templates from.
		Layout:     "layout",          // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl"}, // Specify extensions to load for templates.
		Charset:    "UTF-8",           // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true,
	}))

	m.Use(martini.Recovery())

	m.Get("/", IndexRouter)
	m.Post("/", PostReq)

	m.Run()

}

func IndexRouter(r render.Render) {
	r.HTML(200, "home/index", nil)
}

func PostReq(req *http.Request, r render.Render) {
	responseObject := new(QuestionResponse)

	parts := fmt.Sprintf("%s/v1/question/%s", service_url, req.FormValue("dataset"))

 	u, err := url.Parse(parts)
    if err != nil {
        panic(err)
    }

    jsonString := fmt.Sprintf(`{"question": {"evidenceRequest": {"items": 5 },"questionText": "%s"}}`, req.FormValue("questionText"));   

	questionRequest, err := http.NewRequest("POST", parts,bytes.NewBuffer([]byte(jsonString)))
	questionRequest.Host = u.Host
	questionRequest.Header.Set("Content-Type", "application/json")
	questionRequest.Header.Set("Accept", "application/json")
	questionRequest.Header.Set("X-synctimeout", "30")
	questionRequest.Header.Set("Authorization", auth)

	client := &http.Client{}
    resp, err := client.Do(questionRequest)
    if err != nil {
    	responseObject.Error = string(err.Error())
    	r.HTML(200, "home/index", responseObject)   
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)

    body, _ := ioutil.ReadAll(resp.Body)
    questionResponseBody := []byte(body)
	var questionAnswers []WatsonResponse
	questionAnswers = make([]WatsonResponse,0)

	errParseJson := json.Unmarshal(questionResponseBody, &questionAnswers)
	if errParseJson != nil || len(questionAnswers) <= 0{
		responseObject.Error = "Houston, we got a problem"
		r.HTML(200, "home/index", responseObject)
	}else{
		responseObject.QuestionResponse = questionAnswers[0];
		r.HTML(200, "home/index", responseObject)
	}
}
