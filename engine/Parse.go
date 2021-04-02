package main

// (c) Harry Nelsen 2021

import (
    "fmt"
)

var AllowedErrors = -1
var Coloured = true
var Pipes = false

func P_I(Count int) (string) {  // Indent
    Output := ""
    
    for i := 0; i < Count; i ++ {
        Output = Output + "    "
    }
    
    return Output
}

func P_PrintError(Input string, Line int) {
    ErrorSymbol := "[!]"

    if Coloured {
        ErrorSymbol = "\033[1;31m[!]\033[0m"    // I haven't tested whether or not this works on windows
    }
    
    fmt.Printf("%s Line %d: %s\n", ErrorSymbol, (Line + 1), Input)
}

func P_ParseFile(File []string) (int) {
    Errors := 0
    
    if !(CV_Setup()) {
        P_PrintError("Couldn't setup correctly.", -1)
        if AllowedErrors == 1 {
            return 1
        }
    }

    for CurrentLine, FileLine := range File {
        if Outliers(" " + "\r\n", FileLine) != 0 {  // See if line contains characters other then whitespace
            LineCommands, GlobalStr := CMD_ReadLine(FileLine)    // Convert line into an array of commands and modify GlobalStr
            
            Out_Contents := ""
            Out_Type := -1
            
            for _, Command := range LineCommands {  // Read each command
                if (Out_Type != -1 && Pipes) {
                    if len(Command.Arguments) == 1 {    // So uhh, I can't detect the amount to arguments that it can take so this'll be weird
                        Command.Arguments[0] = Out_Contents
                    } else if len(Command.Arguments) != 1 {
                        P_PrintError("Cannot pipe 1 output into this command.", CurrentLine)
                    } else {
                        P_PrintError("Cannot pipe incompatible types.", CurrentLine)
                    }
                }
                
                Out_Type, Out_Contents = CV_Convert(Command, GlobalStr)
                
                if (Out_Type == 20 && Pipes) {  // Add quotes to string if it can't be detected
                    if CMD_VarType(Out_Contents) != 20 {
                        Out_Contents = fmt.Sprintf("\"%s\"", Out_Contents)
                    }
                    
                    if CMD_VarType(Out_Contents) != 20 {
                        Errors ++
                        P_PrintError("Couldn't add quotes to the output.", CurrentLine)
                    }
                }
                
                if Out_Type == 41 { // See if output type an error
                    Errors ++
                    P_PrintError(Out_Contents, CurrentLine)
                } else if (Out_Type != 99 && Out_Type != CMD_VarType(Out_Contents) && Pipes) {
                    Errors ++
                    P_PrintError("Explicit output type and detected type don't match.", CurrentLine)
                }
                
                if Out_Type == 0 {
                    Out_Contents = ""
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
