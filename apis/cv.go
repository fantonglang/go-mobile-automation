package apis

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"strings"
	"time"

	m "github.com/fantonglang/go-mobile-automation"
)

type CvMixIn struct {
	m.IOperation
	launchCmd   string
	activateCmd string
	dumpsysCmd  string
}

func NewCvMixIn(ops m.IOperation) *CvMixIn {
	return &CvMixIn{
		IOperation:  ops,
		launchCmd:   `am start -n "cn.amghok.opencvhelper/.MainActivity" -a android.intent.action.MAIN -c android.intent.category.LAUNCHER`,
		activateCmd: `am startservice -a ACTIVATE "cn.amghok.opencvhelper/.NetworkingService"`,
		dumpsysCmd:  `dumpsys activity services cn.amghok.opencvhelper`,
	}
}

type CvTemplateMatchingResult struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type CvContourRectResult struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type CvContourCircleResult struct {
	CenterX int `json:"center_x"`
	CenterY int `json:"center_y"`
	Radius  int `json:"radius"`
}

func (cv *CvMixIn) testConnectivity() error {
	req := cv.GetCvRequest()
	_, _, err := req.GetWithTimeout("/ping", time.Second)
	if err == nil {
		return nil
	}
	dumpsysOut, err := cv.Shell(cv.dumpsysCmd)
	if err != nil {
		return err
	}
	if strings.HasSuffix(strings.Trim(dumpsysOut, " \n"), "(nothing)") {
		_, err = cv.Shell(cv.launchCmd)
		if err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
		_, _, err = req.GetWithTimeout("/ping", time.Second)
		return err
	}
	_, err = cv.Shell(cv.activateCmd)
	if err != nil {
		return err
	}
	time.Sleep(time.Second)
	_, _, err = req.GetWithTimeout("/ping", time.Second)
	return err
}

func (cv *CvMixIn) preparaPostCvInput(imageData []byte, definitionData []byte) ([]byte, string) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fwImage, err := w.CreateFormFile("image", "screenshot.png")
	if err != nil {
		return nil, ""
	}
	_, err = fwImage.Write(imageData)
	if err != nil {
		return nil, ""
	}
	fwDefinition, err := w.CreateFormFile("definition", "definition.json")
	if err != nil {
		return nil, ""
	}
	_, err = fwDefinition.Write(definitionData)
	if err != nil {
		return nil, ""
	}
	w.Close()
	return b.Bytes(), w.FormDataContentType()
}

func (cv *CvMixIn) CvTemplateMatching(imageData []byte, definitionData []byte) *CvTemplateMatchingResult {
	err := cv.testConnectivity()
	if err != nil {
		return nil
	}
	inputBytes, contentType := cv.preparaPostCvInput(imageData, definitionData)
	if inputBytes == nil {
		return nil
	}
	bytes, _, err := cv.GetCvRequest().PostWithTimeout("/cv", inputBytes, contentType, 2*time.Second)
	if err != nil {
		return nil
	}
	result := new(CvTemplateMatchingResult)
	err = json.Unmarshal(bytes, result)
	if err != nil {
		return nil
	}
	return result
}

func (cv *CvMixIn) CvContourRect(imageData []byte, definitionData []byte) *CvContourRectResult {
	err := cv.testConnectivity()
	if err != nil {
		return nil
	}
	inputBytes, contentType := cv.preparaPostCvInput(imageData, definitionData)
	if inputBytes == nil {
		return nil
	}
	bytes, _, err := cv.GetCvRequest().PostWithTimeout("/cv", inputBytes, contentType, 2*time.Second)
	if err != nil {
		return nil
	}
	result := new(CvContourRectResult)
	err = json.Unmarshal(bytes, result)
	if err != nil {
		return nil
	}
	return result
}

func (cv *CvMixIn) CvContourCircle(imageData []byte, definitionData []byte) *CvContourCircleResult {
	err := cv.testConnectivity()
	if err != nil {
		return nil
	}
	inputBytes, contentType := cv.preparaPostCvInput(imageData, definitionData)
	if inputBytes == nil {
		return nil
	}
	bytes, _, err := cv.GetCvRequest().PostWithTimeout("/cv", inputBytes, contentType, 2*time.Second)
	if err != nil {
		return nil
	}
	result := new(CvContourCircleResult)
	err = json.Unmarshal(bytes, result)
	if err != nil {
		return nil
	}
	return result
}
