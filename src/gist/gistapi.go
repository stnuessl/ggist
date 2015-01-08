/*
 * Copyright (C) 2014  Steffen Nüssle
 * ggist - go gist
 *
 * This file is part of ggist.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package gist

import (
    "bytes"
    "errors"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "path"
    "strings"
)

type GistAPI struct {
    client http.Client
}

type GistInfo struct {
    Description string
    Public bool
    Files []string
}

type SimpleGistInfo struct {
    Description string
    Public bool
    FileName string
    Data []byte
}

type Gist struct {
    Url string                  `json:"html_url"`
    Id  string                  `json:"id"`
    Description string          `json:"description"`
    Files       map[string]struct {
        Size     int64                  `json:"size"`
        Language string                 `json:"language"`
        Content  string                 `json:"content"`
    }                           `json:"files"`
    Public bool                 `json:"public"`
}

type file struct {
    Content string              `json:"content"`
}

type localGist struct {
    Description string              `json:"description"`
    Public bool                     `json:"public"`
    Files map[string]file           `json:"files"`
}

func NewGistAPI() *GistAPI {
    return &GistAPI{http.Client{}}
}

func (this *GistAPI) CreateGist(info *GistInfo) (*Gist, error) {
     gist, err := newLocalGist(info.Description, info.Public, &info.Files)
     if err != nil {
         return nil, err
     }
     
     return this.uploadLocalGist(gist)
}

func (this *GistAPI) CreateSimpleGist(info *SimpleGistInfo) (*Gist, error) {
    
    gist := &localGist{info.Description, info.Public, make(map[string]file)}

    gist.Files[info.FileName] = file{string(info.Data)}
    
    return this.uploadLocalGist(gist)
}

func (this *GistAPI) DeleteGist(id string) error {
    id = ensureIsGistId(id)
    
    url := fmt.Sprintf("https://api.github.com/gists/%s", id)
    
    resp, err := this.getResponse("DELETE", url, nil)
    if err != nil {
        return err
    }
    
    /* 204: No Content */
    if resp.StatusCode != 204 {
        return errors.New(resp.Status)
    }
    
    return nil
}

func (this *GistAPI) GetGist(id string) (*Gist, error) {
    id = ensureIsGistId(id)
    
    url := fmt.Sprintf("https://api.github.com/gists/%s", id)
    
    resp, err := this.getResponse("GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    defer resp.Body.Close()
    
    return decodeGist(resp.Body)
}

func (this *GistAPI) GetUsersGists(user string) ([]Gist, error) {
    url := fmt.Sprintf("https://api.github.com/users/%s/gists", user)

    resp, err := this.getResponse("GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    defer resp.Body.Close()
    
    if (resp.StatusCode != 200) {
        return nil, errors.New(resp.Status)
    }
    
    gists := make([]Gist, 10)
    
    err = json.NewDecoder(resp.Body).Decode(&gists)
    if err != nil {
        return nil, errors.New("json.NewDecoder.Decode(): " + err.Error())
    }
    
    return gists, nil
}

func (this *GistAPI) UpdateGist(id string, desc string, files []string) error {
    return nil
}

func (this *GistAPI) getResponse(what string, 
                                 url string, 
                                 data []byte) (*http.Response, error) {
    msg, err := newRequest(what, url, data)
    if err != nil {
        return nil, err
    }
    
    return this.client.Do(msg)
}

func newRequest(what string, url string, data []byte) (*http.Request, error) {
    msg, err := http.NewRequest(what, url, bytes.NewReader(data))
    if err != nil {
        return nil, err
    }
    
    msg.Header.Add("Accept", "application/vnd.github.v3+json")
    
    return msg, nil
}

func newLocalGist(desc string,
                  public bool, 
                  files *[]string) (*localGist, error) {
    if len(*files) == 0 {
        return nil, errors.New("Failed to create gist: no files were passed")
    }

    gist := localGist{desc, public, make(map[string]file)}

    for _, x := range *files {
        data, err := ioutil.ReadFile(x)
        if err != nil {
            return nil, errors.New("ioutil.ReadFile(): " + err.Error())
        }
        
        if len(data) == 0 {
            return nil, errors.New("File " + x + " is empty - abort.")
        }
        
        gist.Files[path.Base(x)] = file{string(data)}
    }

    return &gist, nil
}

func decodeGist(data io.Reader) (*Gist, error) {
    gist := Gist{}
    
    err := json.NewDecoder(data).Decode(&gist)
    if err != nil {
        return nil, errors.New("json.NewDecoder.Decode(): " + err.Error())
    }
    
    return &gist, nil
}

func (this *GistAPI) uploadLocalGist(gist *localGist) (*Gist, error) {
    msg_data, err := json.Marshal(gist)
    if err != nil {
        return nil, errors.New("json.Marshal(): " + err.Error())
    }
    
    url := "https://api.github.com/gists"

    resp, err := this.getResponse("POST", url, msg_data)
    if err != nil {
        return nil, err
    }
    
    defer resp.Body.Close()
    
    /* 201 - Created */
    if resp.StatusCode != 201 {
        /* 422 - Unprocessable Entity */
        if resp.StatusCode == 422 {
            return nil, handleMessageUnprocessableEntity(resp.Body)
        }
        
        return nil, errors.New(fmt.Sprintf("Server returned: %s", resp.Status))
    }
    
    return decodeGist(resp.Body)
}

func handleMessageUnprocessableEntity(data io.Reader) error {
    errorStruct := struct {
        Message string          `json:"message"`
        Errors []struct {
            Resource string            `json:"resource"`
            Field string                `json:"field"`
            Code string                 `json:"code"`
        }                       `json:"errors"`
    } {}
    
    err := json.NewDecoder(data).Decode(&errorStruct)
    if err != nil {
        return errors.New("Decoding error message failed with: " + err.Error())
    }
    
    errMsg := fmt.Sprintf("Server returned: %s\n", errorStruct.Message)
    
    length := len(errorStruct.Errors)
    
    for i, x := range errorStruct.Errors {
        errMsg += fmt.Sprintf("Error (%d / %d):\n", i + 1, length)
        errMsg += fmt.Sprintf("  Resource: %s\n  Field: %s\n  Code: %s\n",
                              x.Resource, x.Field, x.Code)
    }
    
    return errors.New(errMsg)
}

func ensureIsGistId(s string) string {
    if strings.Contains(s, "https://gist.github.com/") {
        return s[strings.LastIndex(s, "/") + 1:]
    }
    
    return s
}