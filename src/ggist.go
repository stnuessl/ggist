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

package main

import (
    "gist"
    "errors"
    "fmt"
    "io/ioutil"
    "os"
    "strings"
    "time"
    "util"
)

func stdinIsPipe() (bool, error) {
    stat, err := os.Stdin.Stat()
    if err != nil {
        return false, errors.New("os.Stdin.Mode() failed with: " + err.Error())
    }
    
    return stat.Mode() & os.ModeCharDevice == 0, nil
}

func makeSimpleGist(api *gist.GistAPI, 
                    desc string, 
                    public bool,
                    fileName string) (*gist.Gist, error) {

    
    data, err := ioutil.ReadAll(os.Stdin)
    if err != nil {
        return nil, errors.New("ioutil.ReadAll(): " + err.Error())
    }
    
    if len(data) == 0 {
        fmt.Printf("Nothing to read on stdin - done...\n")
        os.Exit(0)
    }
    
    info := gist.SimpleGistInfo{}
    info.Description = ensureValidDescription(desc)
    info.Public      = public
    info.FileName    = ensureValidFileName(fileName)
    info.Data        = data
    
    gist, err := api.CreateSimpleGist(&info)
    if err != nil {
        err = errors.New("gist.GistAPI.createSimpleGist(): " + err.Error())
        return nil, err
    }
    
    return gist, nil
}

func makeGist(api *gist.GistAPI, 
              desc string, 
              public bool, 
              files *[]string) (*gist.Gist, error) {
    info := gist.GistInfo{}
    info.Description = ensureValidDescription(desc)    
    info.Public      = public
    info.Files       = *files
    
    return api.CreateGist(&info)
}

func ensureValidDescription(desc string) string {
    if len(desc) > 0 {
        return desc
    }
    now := time.Now().Format(time.RFC850)
    
    return fmt.Sprintf("ggist upload dated from %s", now)
}

func ensureValidFileName(fileName string) string {
    if len(fileName) > 0 {
        return fileName
    }
    
    return "unknown"
}

func checkFiles(files []string) ([]string, error) {
    checked := make([]string, 0, len(files))
    
    for _, x := range files {
        stat, err := os.Stat(x)
        
        switch {
        case err != nil:
            return nil, err
        case !stat.Mode().IsRegular():
            return nil, errors.New(x + " is not a regular file")
        default:
            checked = append(checked, x)
        }
    }
    
    return checked, nil
}

func printReceivedGist(gist *gist.Gist, lines bool) {
    msg := "\n" + 
            "Gist        : %s\n" +
            "Url         : %s\n" +
            "Description : %s\n"
            
    fmt.Printf(msg, gist.Id, gist.Url, gist.Description)

    msg = "\n" + 
          "File     : %s\n" +
          "Language : %s\n" +
          "------------------------------------------------------------------\n"
    
    for key, val := range gist.Files {
        fmt.Printf(msg, key, val.Language)
        
        if lines {
            for i, x := range strings.Split(val.Content, "\n") {
                fmt.Printf("%4d | %s\n", i + 1, x)
            }
        } else {
            fmt.Printf("%s\n", val.Content)
        }
    }
}

func printUploadedGist(gist *gist.Gist, verbose bool) {
    
    if verbose {
        msg := "Created gist %s\n" +
               "  Url         : %s\n" +
               "  Public      : %s\n" +
               "  Description : %s\n" +
               "  Files       :\n"
        
        fmt.Printf(msg, gist.Id, gist.Url, gist.Public, gist.Description)
        
        for key, _ := range gist.Files {
            fmt.Printf("    %s\n", key)
        }
    } else {
        fmt.Printf("Created gist %s: %s\n", gist.Id, gist.Url)
    } 
}

func main() {
    descDesc        := "Add a description when uploading a gist."
    descName        := "Set a filename; Useful when uploading from stdin."
    descFiles       := "Set files to upload as gist."
    descIndex       := "Get Gists with Index i from history."
    descHelp        := "Print this help message."
    descGet         := "Download specified gists."
    descLineN       := "Print line numbers in source files."
    descHist        := "Print your gist history"
    descVerb        := "Print more information about gists if possible."
    descUsers       := "Retrieve gists from a user."
    
    var desc string
    var fileName string
    var files []string
    var gets []string
    var help bool
    var history bool
    var lineNum bool
    var private bool
    var verbose bool
    var index []int
    var users []string
    
    home := os.Getenv("HOME")
    
    if len(home) == 0 {
        util.Error("Unable to find home directory")
        os.Exit(1)
    }
    
    gistHistory, err := gist.NewHistory(home + "/.config/ggist/history")
    if err != nil {
        util.Error(err)
        os.Exit(1)
    }
    
    options := []util.Option{
        &util.OptStr    { "description,d",  descDesc,  &desc       },
        &util.OptMulStr { "files,f",        descFiles, &files      },
        &util.OptBool   { "help",           descHelp,  &help       },
        &util.OptMulStr { "get,g",          descGet,   &gets       },
        &util.OptBool   { "line-numbers,l", descLineN, &lineNum    },
        &util.OptBool   { "history,h",      descHist,  &history    },
        &util.OptMulInt { "index,i",        descIndex, &index      },
        &util.OptStr    { "file-name,n",    descName,  &fileName   },
        &util.OptBool   { "verbose,v",      descVerb,  &verbose    },
        &util.OptMulStr { "user,u",         descUsers, &users      },
    }
    
    err = util.ParseCommandLine(options, os.Args[1:])
    if err != nil {
        util.Error(err)
        os.Exit(1)
    }

    if help {
        util.PrintCommandHelp(options)
        os.Exit(0)
    }
    
    valid_files, err := checkFiles(files)
    if err != nil {
        util.Error("Invalid file: " + err.Error())
        os.Exit(1)
    }
    
    api := gist.NewGistAPI()
    var gist *gist.Gist
    
    isPipe, err := stdinIsPipe()
    if err != nil {
        util.Warning("unable to read from stdin")
    }
    
    switch {
    case len(valid_files) > 0:
        gist, err = makeGist(api, desc, !private, &valid_files)
    case isPipe:
        gist, err = makeSimpleGist(api, desc, !private, fileName)
    }
    
    if err != nil {
        util.Error(err)
        os.Exit(1)
    } else if gist != nil {
        printUploadedGist(gist, verbose)
        
        gistHistory.AddGist(gist)
    }
    
    for _, x := range index {
        id := gistHistory.GetGistIdAt(x)
        
        gist, err = api.GetGist(id)
        if err != nil {
            util.Error(err)
            os.Exit(1)
        }
        
        gistHistory.AddGist(gist)
        
        printReceivedGist(gist, lineNum)
    }
    
    for _, x := range gets {
        gist, err = api.GetGist(x)
        if err != nil {
            util.Error(err)
            os.Exit(1)
        }
        
        gistHistory.AddGist(gist)
        
        printReceivedGist(gist, lineNum)
    }
    
    for _, x := range users {
        gists, err := api.GetUsersGists(x)
        if err != nil {
            util.Error(err)
            os.Exit(1)
        }
        
        for _, y := range gists {
            fmt.Printf("Gist Id: %s\n", y.Id)
        }
    }
    
    if history {
        fmt.Printf("Gist History:\n%s\n", gistHistory)
    }
}
