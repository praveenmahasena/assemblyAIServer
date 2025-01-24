package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	tokenParam           = "TOKEN"
	mediauploadlinkParam = "MEDIAUPLOADLINK"
	transcribelinkParam  = "TRANSCRIBELINK"
)

type TCPResponse struct {
	Status int    `json:"status"`
	Data   string `json:"data"`
}

func NewResponse(s int, d string) TCPResponse {
	return TCPResponse{
		s,
		d,
	}
}

func (t TCPResponse)WriteTo(con net.Conn){
	data,dataErr:=json.Marshal(t)
	if dataErr!=nil{
		con.Write([]byte("Error"))
		return
	}
	con.Write(data)
}

type Res struct {
	Upload_url string `json:"upload_url"`
}

func handle(con net.Conn) {
	defer con.Close()

	s := bufio.NewScanner(con)
	s.Scan()
	fb := s.Bytes()

	link, err := upLoadFile(fb)
	if err != nil {
		r := NewResponse(http.StatusInternalServerError, err.Error())
		r.WriteTo(con)
		return
	}

	transcriptionID, transcriptionIDErr := getTranScriptID(link)
	if transcriptionIDErr != nil {
		r := NewResponse(http.StatusInternalServerError, transcriptionIDErr.Error())
		r.WriteTo(con)
		return
	}
	transcription, transcriptionErr := getTranScript(transcriptionID)
	if transcriptionErr != nil {
		r := NewResponse(http.StatusInternalServerError, transcriptionErr.Error())
		r.WriteTo(con)
		return
	}
	r := NewResponse(http.StatusInternalServerError, transcription)
	r.WriteTo(con)
}

/* write docs for lsp */
func upLoadFile(fb []byte) (string, error) {
	mediaLink := os.Getenv(mediauploadlinkParam)
	token := os.Getenv(tokenParam)
	req, reqErr := http.NewRequest(http.MethodPost, mediaLink, bytes.NewReader(fb))
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type" ,"application/octet-stream")
	fmt.Println(req,reqErr,mediaLink,token)
	r:=struct{
		Upload_url string `json:"upload_url"`
	}{}
	if reqErr != nil {
		return "", fmt.Errorf("error during request for transcript link %v", reqErr)
	}
	 res, resErr := http.DefaultClient.Do(req)
	 if resErr != nil {
	 	return "", fmt.Errorf("error during request for transcript link %v", resErr)
	 }
	defer res.Body.Close()
	if err:=json.NewDecoder(res.Body).Decode(&r);err!=nil{
	 	return "", fmt.Errorf("error decoding transcript link %v", err)
	}
	 return r.Upload_url,nil
}

func getTranScriptID(payloadLink string) (string, error) {
	link := os.Getenv(transcribelinkParam)
	payload := fmt.Sprintf(`{"audio_url": "%v"}`, payloadLink)
	r := bytes.NewBufferString(payload)
	id := struct {
		ID string `json:"id"`
	}{}
	<-time.After(time.Duration(3) * time.Second)
	req, reqErr := http.NewRequest(http.MethodPost, link, r)
	req.Header.Set("Authorization", os.Getenv(tokenParam))
	req.Header.Set("Content-type", "Application/json")
	if reqErr != nil {
		return "", fmt.Errorf("error during sending request to get transcript %v", reqErr)
	}
	res, resErr := http.DefaultClient.Do(req)
	if resErr != nil {
		return "", fmt.Errorf("error during getting response to get transcript %v", resErr)
	}
	if err := json.NewDecoder(res.Body).Decode(&id); err != nil {
		return "", fmt.Errorf("error reading id %v", err)
	}
	return id.ID, nil
}

func getTranScript(id string) (string, error) {
	data := struct {
		Text string `json:"text"`
	}{}
	uri := os.Getenv(transcribelinkParam) + "/" + id
	for {
		<-time.After(time.Duration(10) * time.Second)
		req, reqErr := http.NewRequest(http.MethodGet, uri, nil)
		req.Header.Set("Authorization", os.Getenv(tokenParam))
		if reqErr != nil {
			return "", fmt.Errorf("error during request of transcription %v", reqErr)
		}
		res, resErr := http.DefaultClient.Do(req)
		if resErr != nil {
			return "", fmt.Errorf("error during response to transcription %v", resErr)
		}
		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
			return "", fmt.Errorf("error during decoding transcription %v", err)
		}
		if data.Text != "" {
			return data.Text, nil
		}
	}
}
