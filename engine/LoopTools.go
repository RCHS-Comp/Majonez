package main

// FruitLoops
// (c) Harry Nelsen 2021

import (
    "fmt"
)

//////////////////////////////////////////
// List possible statements for both loop types
//////////////////////////////////////////

var LOOP_PossibleStatements = []string{"loop_st", "loop_nd", "loop_bk"}
var IF_PossibleStatements = []string{"if", "elseif", "else", "endif"}

//////////////////////////////////////////

//////////////////////////////////////////
// Definitions for the commands
//////////////////////////////////////////

var loop_start = Meta{Name: "loop_st", Arguments: 0, InputTypes: nil, OutputType: 99}
var loop_end = Meta{Name: "loop_nd", Arguments: 0, InputTypes: nil, OutputType: 99}
var loop_break = Meta{Name: "loop_bk", Arguments: 0, InputTypes: nil, OutputType: 99}

var if_statement = Meta{Name: "if", Arguments: 3, InputTypes: []int{42, 40, 40}, OutputType: 99}
var elseif_statement = Meta{Name: "elseif", Arguments: 3, InputTypes: []int{42, 40, 40}, OutputType: 99}
var else_statement = Meta{Name: "else", Arguments: 0, InputTypes: nil, OutputType: 99}
var endif_statement = Meta{Name: "endif", Arguments: 0, InputTypes: nil, OutputType: 99}

//////////////////////////////////////////

//////////////////////////////////////////
// Global variables for loop counts
//////////////////////////////////////////

var ClosedLoops = 0
var OpenLoops = 0

//////////////////////////////////////////

//////////////////////////////////////////
// Functions for parsing the loops
//////////////////////////////////////////

func LOOP_ParseStatement(Statement Command, IndentCount *int) (string, bool) {
    Output := ""

    switch Statement.Name {
        case "loop_st":
            Output = P_I(*IndentCount) + "while True:"
            
            OpenLoops ++
            *IndentCount ++
        case "loop_nd": // In a lot of languages, this should be the same as endif (not all)
            if OpenLoops - (ClosedLoops + 1) < 0 {
                return "Cannot end a statement that hasn't been started.", false
            } else if *IndentCount < 1 {
                return "Indent count too low for this function.", false
            }
        
            // For something like C
            // Output = P_I(*IndentCount - 1) + "# }"
            
            ClosedLoops ++
            *IndentCount --
        case "loop_bk":
            // This one is weird
            return (P_I(*IndentCount) + "break"), true
    }
    
    return Output, true
}

func IF_ParseStatement(Statement Command, IndentCount *int) (string, bool) {
    Output := ""
    
    // Normally, we wouldn't need to get the variable types in the command
    // This is a werid command so we will have to (Mainly because no overloading)
    
    ArgOneType := -1
    ArgTwoType := -1
    ReturnStatement := []string{"", " != ", ""}
    
    if len(Statement.Arguments) >= 3 {  // Logic for comparative statements
        ReturnStatement[0] = Statement.Arguments[1]
        ReturnStatement[2] = Statement.Arguments[2]
        
        if (Statement.Arguments[0] != "!") {
            ReturnStatement[1] = " == "
        }
        
        ArgOneType = CMD_VarType(Statement.Arguments[1])
        ArgTwoType = CMD_VarType(Statement.Arguments[2])
        
        for Index, Arg := range []int{ArgOneType, ArgTwoType} {
            if Arg == 6 {
                VarName := Statement.Arguments[Index + 1]
                if M_VariableExist(VarName) {
                    Arg = M_VariableType(VarName)
                } else {
                    return fmt.Sprintf("Variable '%s' not found.", VarName), false
                }
            }
        }
    }

    switch Statement.Name {
        case "if":
            if ArgOneType != ArgTwoType {
                return fmt.Sprintf("Trying to compare %s with %s.", M_TypeToStr(ArgOneType), M_TypeToStr(ArgTwoType)), false
            }
            
            Output = P_I(*IndentCount) + fmt.Sprintf("if %s:", Flatten(ReturnStatement))
            
            OpenLoops ++
            *IndentCount ++
        case "elseif":
            if OpenLoops - (ClosedLoops + 1) < 0 {
                return "No open statements.", false
            } else if ArgOneType != ArgTwoType {
                return fmt.Sprintf("Trying to compare %s with %s.", M_TypeToStr(ArgOneType), M_TypeToStr(ArgTwoType)), false
            }
            
            Output = P_I(*IndentCount - 1) + fmt.Sprintf("elif %s:", Flatten(ReturnStatement))
            
            OpenLoops ++
            ClosedLoops ++
        case "else":
            if OpenLoops - (ClosedLoops + 1) < 0 {
                return "No open statements.", false
            }
            
            Output = P_I(*IndentCount - 1) + "else:"
            
            OpenLoops ++
            ClosedLoops ++
        case "endif":
            if OpenLoops - (ClosedLoops + 1) < 0 {
                return "Cannot end a statement that hasn't been started.", false
            } else if *IndentCount < 1 {
                return "Indent count too low for this function.", false
            }
        
            // For something like C
            // Output = P_I(*IndentCount - 1) + "# }"
            
            ClosedLoops ++
            *IndentCount --
    }
    
    return Output, true
}

//////////////////////////////////////////

//////////////////////////////////////////
// This makes sure everything is closed
//////////////////////////////////////////

func IF_AllClosed() (bool) {
    return (OpenLoops - ClosedLoops == 0)
}

//////////////////////////////////////////
