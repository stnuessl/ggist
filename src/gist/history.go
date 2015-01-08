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
    "encoding/csv"
    "fmt"
    "os"
    "strings"
)

type gistTuple struct {
    id string
    description string
}

type History struct {
    path string
    gists []gistTuple
    file os.File
}

func NewHistory(path string) (*History, error) {
    index := strings.LastIndex(path, "/")
    
    dirs := path[:index]
    
    err := os.MkdirAll(dirs, 0755)
    if err != nil {
        return nil, errors.New("os.MkdirAll() failed with: " + err.Error())
    }
    
    file, err := os.OpenFile(path, os.O_RDWR, 0666)
    if err != nil {
        if !os.IsNotExist(err) {
            return nil, err
        }
        
        file, err = os.Create(path)
        if err != nil {
            return nil, err
        }
    }
    
    history := History{path, make([]gistTuple, 0, 100), *file}
    
    stat, err := history.file.Stat()
    if err != nil {
        return nil, err
    }
    
    size := stat.Size()
    
    if size != 0 {
        buf := make([]byte, size)
        
        n, err := history.file.Read(buf)
        switch {
        case err != nil:
            return nil, err
        case int64(n) != size:
            return nil, errors.New("Failed to read complete file " + path)
        }
        
        reader := bytes.NewReader(buf)
        
        all, err := csv.NewReader(reader).ReadAll()
        if err != nil {
            return nil, err
        }
        
        for _, x := range all {
            if len(x) != 2 {
                msg := "Invalid history file. Manually fix or remove " + path
                return nil, errors.New(msg)
            }

            history.gists = append(history.gists, gistTuple{x[0], x[1]})
        }
    }
    
    return &history, nil
}

func (this *History) String() string {
    ret := ""
    
    length := len(this.gists)
    
    for i, x := range this.gists {
        ret += fmt.Sprintf("%4d <> %s : %s\n", length - i, x.id, x.description)
    }
    
    /* Do not return last newline */
    return ret[:len(ret) - 1]
}

func (this *History) GetGistIdAt(i int) string {
    return this.gists[len(this.gists) - i].id
}

func (this *History) AddGist(gist *Gist) {
    if this.gists[len(this.gists) - 1].id != gist.Id {
        
        this.gists = append(this.gists, gistTuple{gist.Id, gist.Description})
        this.file.WriteString(fmt.Sprintf("%s,%s\n", gist.Id, gist.Description))
    }
}
