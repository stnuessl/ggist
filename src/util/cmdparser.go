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

func setOpt(opt Option, 
            argMap map[string]Option,
            argv []string, 
            i int) (int, error) {
    
    list := possibleArgs(argMap, argv[i + 1:])
    i += len(list)
    
    switch opt.(type) {
    case *OptBool:
        if len(list) > 0 {
            return 0, errors.New("Invalid argument - none exspected.")
        }
        
        err := opt.Set("")
        if err != nil {
            return 0, err
        }
    case *OptMulStr, *OptMulInt:
        if len(list) == 0 {
            msg := "Invalid argument - at least one exspected."
            return 0, errors.New(msg)
        }
        
        
        for _, x := range list {
            err := opt.Set(x)
            if err != nil {
                return 0, err
            }
        }
        
        i += len(list)
    case *OptInt, *OptStr:
        if len(list) != 1 {
            msg := "Exactly one argument exspected, "
            msg += fmt.Sprintf("%d provided.", len(list))
            return 0, errors.New(msg)
        }
        
        
        err := opt.Set(list[0])
        if err != nil {
            return 0, err
        }
        
        i += 1
    default:
        return 0, errors.New("Invalid option type")
    }
    
    return i, nil
}

func ParseCommandLine(opts []Option, argv []string) ([]string, error) {
    argMap, err := newArgumentMap(opts)
    if err != nil {
        return nil, err
    }
    
    args := make([]string, 0, len(argv) + 4);
    
    for _, x := range argv {
        
        if len(x) > 1 && x[0] == '-' && x[1] == '-' {
            args = append(args, x[2:])
        } else if len(x) > 1 && x[0] == '-' && x[1] != '-' {
            args = append(args, x[1:])
        } else {
            args = append(args, x)
        }
    }
    
    noOpts := make([]string, 0, len(args));
    
    for i := 0; i < len(args); i++ {
        arg := args[i]
        
        opt, ok := argMap[arg]
        if ok {
            j, err := setOpt(opt, argMap, args, i)
            if err != nil {
                return nil, err
            }
            
            i = j
        } else {
            noOpts = append(noOpts, arg)
        }
    }
    
    return noOpts, nil
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