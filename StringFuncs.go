package main

// String tools that I often use
// (c) Harry Nelsen 2020

import (
    "strings"
)

func In(LookingIn string, LookingFor rune) (int) {
    for Index, Element := range LookingIn {
        if Element == LookingFor {
            return Index
        }
    }
    
    return -1
}

func In_string(LookingIn []string, LookingFor string) (int) {
    for Index, Element := range LookingIn {
        if Element == LookingFor {
            return Index
        }
    }
    
    return -1
}

func In_int(LookingIn []int, LookingFor int) (int) {
    for Index, Element := range LookingIn {
        if Element == LookingFor {
            return Index
        }
    }
    
    return -1
}

func Usage_strings(LookingIn []string, LookingFor string) (int) {
    Counter := 0

    for _, Element := range LookingIn {
        if Element == LookingFor {
            Counter ++
        }
    }
    
    return Counter
}

func Usage(LookingIn string, LookingFor rune) (int) {
    Counter := 0

    for _, Element := range LookingIn {
        if Element == LookingFor {
            Counter ++
        }
    }
    
    return Counter
}

func Find(LookingIn string, LookingFor rune, Time int) (int) {
    Counter := 0
    
    for Index, Element := range LookingIn {
        if Element == LookingFor {
            Counter ++
            
            if Counter == Time {
                return Index
            }
        }
    }
    
    return -1
}

func Outliers(Source string, Input string) (int) {
    Counter := 0

    for _, Element := range Input {
        if In(Source, Element) == -1 {
            Counter ++
        }
    }
    
    return Counter
}

func RemoveSpace(Source string) (string) {
    return strings.Replace(Source, " ", "", -1)
}

func StripSpace(Source string) (string) {
    return strings.TrimSpace(Source)
}

func Split(Source string, Input string) ([]string) {
    return strings.Split(Source, Input)
}

func IsInt(Source string) (bool) {
    return Outliers("0123456789", Source) == 0
}

func IsFloat(Source string) (bool) {
    return (Outliers("0123456789.", Source) == 0 && Usage(Source, '.') <= 1)
}

func Contains(Source string, LookingFor string) (bool) {
    return strings.Contains(Source, LookingFor)
}

func Flatten(Source []string) (string) {
    Output := ""
    for _, Item := range Source {
        Output = Output + Item
    }
    
    return Output
}

