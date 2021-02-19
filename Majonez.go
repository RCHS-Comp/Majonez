package main

// (c) Harry Nelsen 2021

import (
    "fmt"
    "strings"
    "os"
    "bufio"
)

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
    var Errors int

    if len(os.Args) > 1 {
        Errors = P_ParseFile(MJ_ReadFile(os.Args[1]))
    }
    
    fmt.Printf("Finished with %d error(s)\n", Errors)
}
