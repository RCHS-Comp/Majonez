package main

// (c) Harry Nelsen 2021

import (
    "fmt"
    "strings"
    "os"
    "bufio"
    "runtime"
)

var MJ_V_Major = 1
var MJ_V_Minor = 1
var MJ_V_Micro = 1
var MJ_Year = 2021

var MJ_InfoTriggers = []string{"-h", "--help", "--info", "--about"}

func MJ_Info() {
    fmt.Println("//////////////////////////////////////////")
    fmt.Printf("Majonez [Version PRE-%d.%d.%d]\n", MJ_V_Major, MJ_V_Minor, MJ_V_Micro)
    fmt.Printf("OS: %s, Arch: %s\n", runtime.GOOS, runtime.GOARCH)
    fmt.Printf("(c) Harry Nelsen %d\n", MJ_Year)
    fmt.Println("//////////////////////////////////////////")
    
    fmt.Println("Settings:")
    if AllowedErrors == -1 {
        fmt.Println("Errors allowed: Unlimited")
    } else {
        fmt.Printf("Errors allowed: %d\n", AllowedErrors)
    }
    fmt.Printf("HGF Version: %d.%d.%d\n", HGF_V_Major, HGF_V_Minor, HGF_V_Micro)
    fmt.Println("//////////////////////////////////////////")
    
    fmt.Println("Unavailable variables:")
    for _, VarName := range VariableUnavailable {
        fmt.Println(VarName)
    }
    fmt.Println("//////////////////////////////////////////")
    
    fmt.Println("Usage:")
    fmt.Println("Majonez FILENAME   -   Convert FILENAME")
    fmt.Println("Majonez --about    -   Show this menu")
    fmt.Println("//////////////////////////////////////////")
}

func MJ_WriteFile(ToWrite string) {
    // fmt.Println(ToWrite)
}

func MJ_FileExist(FileName string) (bool) {
    if _, Error := os.Stat(FileName); Error == nil {
        return true
    } else {
        return false
    }
}

func MJ_ReadFile(FileName string) ([]string) {
    if !(MJ_FileExist(FileName)) {
        return nil
    }

    FileObj, Error := os.Open(FileName)
    if Error != nil {
        return nil
    }
    
    defer FileObj.Close()
    var Output []string
    var ToAppend string

    ScanObj := bufio.NewScanner(FileObj)
    
    for ScanObj.Scan() {
        ToAppend = ScanObj.Text()
        ToAppend = strings.TrimSpace(ToAppend)
        Output = append(Output, ToAppend)
    }

    if Error := ScanObj.Err(); Error != nil {
        return nil
    }
    
    return Output
}

func main() {
    Errors := 0

    if len(os.Args) > 1 {
        CMD_Arg := os.Args[1]
    
        if In_string(MJ_InfoTriggers, CMD_Arg) == -1 {
            fmt.Printf("Reading %s...\n", CMD_Arg)
            Errors = P_ParseFile(MJ_ReadFile(CMD_Arg))
        } else {
            MJ_Info()
        }
    } else {
        fmt.Println("No arguments.")
    }
    
    fmt.Printf("\nFinished with %d error(s)\n", Errors)
}
