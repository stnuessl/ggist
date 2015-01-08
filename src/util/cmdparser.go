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

package util

import (
    "errors"
    "fmt"
    "strings"
    "strconv"
)

type Option interface {
    GetOptions() (string, string, error)
    GetDescription() string
    Set(val string) error
}

func getOptions(str string) (string, string, error) {
    s := strings.Trim(str, " ")
    
    if strings.Contains(s, ",") {
        opts := strings.Split(s, ",")
        
        if len(opts) > 2 {
            return "", "", errors.New("Invalid option string")
        }
        
        for _, x := range opts {
            x = strings.Trim(x, " ")
        }
        
        len0 := len(opts[0])
        len1 := len(opts[1])
        
        if len0 == 0 && len1 == 0 {
            return "", "", errors.New("Only empty arguments provided")
        }
        
        if len1 < len0 {
            return opts[1], opts[0], nil
        } else {
            return opts[0], opts[1], nil
        }
    }
    
    return "", s, nil
}

type OptInt struct {
    OptStr string
    
    Description string

    Val *int
}


func (this OptInt) GetOptions() (string, string, error) {
    return getOptions(this.OptStr)
}

func (this OptInt) GetDescription() string {
    return this.Description
}

func (this OptInt) Set(s string) error {
    val, err := strconv.Atoi(s)
    if err != nil {
        return err
    }
    
    *this.Val = val
    
    return nil
}

type OptMulInt struct {
    OptStr string
    
    Description string
    
    Val *[]int
}

func (this OptMulInt) GetOptions() (string, string, error) {
    return getOptions(this.OptStr)
}

func (this OptMulInt) GetDescription() string {
    return this.Description
}

func (this OptMulInt) Set(s string) error {
    val, err := strconv.Atoi(s)
    if err != nil {
        return err
    }
    
    *this.Val = append(*this.Val, val)
    
    return nil
}

type OptBool struct {
    OptStr string
    
    Description string
    
    Val *bool
}

func (this OptBool) GetOptions() (string, string, error) {
    return getOptions(this.OptStr)
}

func (this OptBool) GetDescription() string {
    return this.Description
}

func (this OptBool) Set(s string) error {
    *this.Val = true
    
    return nil
}

type OptStr struct {
    OptStr string
    
    Description string
    
    Val *string
}

func (this OptStr) GetOptions() (string, string, error) {
    return getOptions(this.OptStr)
}

func (this OptStr) GetDescription() string {
    return this.Description
}

func (this OptStr) Set(s string) error {
    *this.Val = s
    
    return nil
}

type OptMulStr struct {
    OptStr string
    
    Description string
    
    Val *[]string
}

func (this OptMulStr) GetOptions() (string, string, error) {
    return getOptions(this.OptStr)
}

func (this OptMulStr) GetDescription() string {
    return this.Description
}

func (this OptMulStr) Set(s string) error {
    *this.Val = append(*this.Val, s)
    
    return nil
}

func newArgumentMap(opts []Option) (map[string]Option, error) {
    m := make(map[string]Option, len(opts))
    
    for _, x := range opts {
        
        s, l, err := x.GetOptions()
        if err != nil {
            return nil, err
        }
        
        if len(s) > 0 {
            m[s] = x
        }
        
        if len(l) > 0 {
            m[l]  = x
        }
    }
    
    return m, nil
}

func possibleArgs(argMap map[string]Option, argv[]string) []string {
    i := 0
    for _, x := range argv {
        _, ok := argMap[x]
        if ok {
            break
        }
        i += 1
    }
    return argv[:i]
}

func setOpt(argv []string, i int, argMap[string]Option) (int, error) {
    
    opt, ok := argMap[argv[i]]
    if ok {
        switch opt.(type) {
            case *OptBool:
                opt.Set("")
                continue
        }
        
        list := possibleArgs(argMap, argv[i + 1:])
        i += len(list)
        
        if len(args) == 0 {
            err := "Missing argument for option " + argv[i]
            return 0, errors.New(err)
        }
        
        for _, x := range args {
            err := opt.Set(x)
            if err != nil {
                return 0, err
            }
        }
        
        return i, nil
    } else {
        return 0, errors.New("Unrecognized option: " + argv[i])
    }
}

func ParseCommandLine(opts []Option, argv []string) error {
    argMap, err := newArgumentMap(opts)
    if err != nil {
        return err
    }
    
    for k, _ := range argMap {
        fmt.Printf("Option: %s\n", k)
    }
    
    return nil
    
    argc := len(argv)
    
    for i := 0; i < argc; i++ {
        s := argv[i]
        
        sLen := len(s)
        
        switch {
        case sLen > 1 && s[0] == '-' && s[1] == '-':
            s = s[2:]
            
            i, err = setOpt(argv, i, argMap)
            if err != nil {
                return err
            }
        case sLen > 0 && s[0] == '-':
            s = s[1:]
            
            for _, x := range s {
                i, err = setOpt(argv, i, argMap)
                if err != nil {
                    return err
                }
            }
        default:
            continue
        }
        
        opt, ok := argMap[argv[i]]
        if ok {
            
        } else {
            
        } 
    }
    
    return nil
}

func PrintCommandHelp(opts []Option) {
    for i, x := range opts {
        d := x.GetDescription()
        
        s, l, err := x.GetOptions()
        if err != nil {
            fmt.Printf("Invalid Option at position %d\n", i)
        } else {
            fmt.Printf("  %2s [ %-15s ] %s\n", s, l, d)
        }
    }
}