package main

// (c) Harry Nelsen 2021

import (
    "fmt"
)

var AllowedErrors = -1
var Coloured = true

// Due to an oversight (mainly me neglecting how we'd pipe commands) piping won't be added

func P_I(Count int) (string) {
    Output := ""
    
    for i := 0; i < Count; i ++ {
        Output = Output + "    "
    }
    
    return Output
}

func P_PrintError(Input string, Line int) {
    ErrorSymbol := "[!]"

    if Coloured {
        ErrorSymbol = "\033[1;31m[!]\033[0m"
    }
    
    fmt.Printf("%s Line %d: %s\n", ErrorSymbol, (Line + 1), Input)
}

func P_ParseFile(File []string) (int) {
    var GlobalStr []string
    var Output string
    Errors := 0
    VarOut := 0
    
    if !(CV_Setup()) {
        P_PrintError("Couldn't setup correctly.", -1)
        Errors ++
    }

    for CurrentLine, FileLine := range File {
        if Outliers(" " + "\r\n", FileLine) != 0 {
            LineCommands, GlobalStr := CMD_ReadLine(FileLine, GlobalStr)
            
            for _, Command := range LineCommands {
                VarOut, Output = CV_Convert(Command, GlobalStr)
                if VarOut == 5 {
                    Errors ++
                    P_PrintError(Output, CurrentLine)
                }
                
                if (Errors > AllowedErrors && !(AllowedErrors == -1)) {
                    return Errors
                }
            }
        }
    }
    
    if !(CV_Teardown()) {
        P_PrintError("Couldn't teardown correctly.", len(File))
        Errors ++
    }

    return Errors
}
