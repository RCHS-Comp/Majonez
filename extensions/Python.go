package main

// Python handler for HGF
// (c) Harry Nelsen 2021

import (
    "fmt"
)

var Peeking = true
var IndentCount = 0

func CV_Setup() (bool) {
    IndentCount = 0
    return true
}

func CV_Teardown() (bool) {
    return IF_AllClosed()
}

func CV_Convert(Input Command, GlobalStr []string) (int, string) {
    Output := ""
    
    switch Input.Name {
        case "echo":
            if !(echo.Compliant(Input)) {
                return 5, "Invalid type. (Only strings supported.)"
            }
            
            Output = fmt.Sprintf("print(%s)", Input.Arguments[0])
        case "test":
            if !(test.Compliant(Input)) {
                return 5, "Invalid usage."
            }
            
            Output = "# Success"
        case "rem":
            if !(rem.Compliant(Input)) {
                return 5, "Invalid usage."
            }
            
            Output = "# " + Input.Arguments[0]
            
        //////////////////////////////////////////
        // Everything past here is handled by FL /
        //////////////////////////////////////////
        
        case "if":
            if !(if_statement.Compliant(Input)) {
                return 5, "Invalid type."
            }
            
            Success := false    // Create here because if we do "Output, Success := ..." it will make Output unable to be returned
            Output, Success = IF_ParseStatement(Input, &IndentCount)
            
            if !(Success) {
                return 5, Output
            }
        case "elseif":
            if !(elseif_statement.Compliant(Input)) {
                return 5, "Invalid type."
            }
            
            Success := false
            Output, Success = IF_ParseStatement(Input, &IndentCount)
            
            if !(Success) {
                return 5, Output
            }
        case "else":
            if !(else_statement.Compliant(Input)) {
                return 5, "Invalid usage."
            }
            
            Success := false
            Output, Success = IF_ParseStatement(Input, &IndentCount)
            
            if !(Success) {
                return 5, Output
            }
        case "endif":
            if !(endif_statement.Compliant(Input)) {
                return 5, "Invalid usage."
            }
            
            Success := false
            Output, Success = IF_ParseStatement(Input, &IndentCount)
            
            if !(Success) {
                return 5, Output
            }
        case "loop_st":
            if !(loop_start.Compliant(Input)) {
                return 5, "Invalid usage."
            }
            
            Success := false
            Output, Success = LOOP_ParseStatement(Input, &IndentCount)
            
            if !(Success) {
                return 5, Output
            }
        case "loop_nd":
            if !(loop_end.Compliant(Input)) {
                return 5, "Invalid usage."
            }
            
            Success := false
            Output, Success = LOOP_ParseStatement(Input, &IndentCount)
            
            if !(Success) {
                return 5, Output
            }
        case "loop_bk":
            if !(loop_break.Compliant(Input)) {
                return 5, "Invalid usage."
            }
            
            Success := false
            Output, Success = LOOP_ParseStatement(Input, &IndentCount)
            
            if !(Success) {
                return 5, Output
            }
        default:
            return 5, "Missing command."
    }
    
    if (In_string(IF_PossibleStatements, Input.Name) == -1 && In_string(LOOP_PossibleStatements, Input.Name) == -1) {
        Output = P_I(IndentCount) + Output
    }
    
    Output = CMD_RestoreStr(Output, GlobalStr)
    
    if Peeking {
        fmt.Println(Output)
    }
    
    return 0, Output
}
