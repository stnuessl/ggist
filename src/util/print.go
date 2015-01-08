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
    "fmt"
)

func Error(x interface{}) {
    fmt.Printf("** ERROR: %s\n", x)
}

func Warning(x interface{}) {
    fmt.Printf("** WARNING: %s\n", x)
}

func Info(x interface{}) {
    fmt.Printf("** INFO: %s\n", x)
}

func Debug(x interface{}) {
    fmt.Printf("** DEBUG: %s\n", x)
}