package main

// (c) Harry Nelsen 2021

// Variable types
// -1   -   Unknown type
// 0    -   Nothing
// 1    -   Boolean
// 10   -   Int
// 11   -   Float
// 12   -   Either int or float
// 20   -   String
// 30   -   Variable
// 40   -   Anything
// 41   -   Error
// 42   -   Flag
// 99   -   Special

var nop = Meta{Name: "nop", Arguments: 0, InputTypes: nil, OutputType: 0}   // Meta for nop command

var VariableAllowedChars = "abcdefghijklmnopqrstuvwxyz"
var VariableAllowedTypes = []int{1, 10, 11, 20}
var VariableUnavailable = []string{"!", "_", "?"}
var VariableNames = []string{"!", "_", "?"}
var VariableTypes = []int{42, 42, 42}

func M_VariableExist(Input string) (bool) {
    return (In_string(VariableNames, Input) != -1)
}

func M_VarAvailable(Input string) (bool) {
    return (In_string(VariableUnavailable, Input) == -1)
}

func M_VariableType(Input string) (int) {
    VarType := -1
    Location := In_string(VariableNames, Input)

    if Location != -1 {
        VarType = VariableTypes[Location]
    }
        
    return VarType
}

func M_AddVar(Name string, Type int) (bool) {
    if (In_string(VariableNames, Name) != -1 || In_string(VariableUnavailable, Name) != -1) {   // See if variable doesn't already exist and is valid
        return false
    }
    
    if In_int(VariableAllowedTypes, Type) != -1 {   // See if it's an allowed type
        return false
    }
    
    if (Outliers(VariableAllowedChars, Name) != 0 || len(Name) < 1) {   // See if it contains invalid characters or length is less then 1
        return false
    }
    
    VariableNames = append(VariableNames, Name)
    VariableTypes = append(VariableTypes, Type)
    
    return true
}

func M_TypeToStr(Type int) (string) {
    switch Type {
        case 0:
            return "non"
        case 1:
            return "boo"
        case 10:
            return "int"
        case 11:
            return "flo"
        case 12:
            return "num"
        case 20:
            return "str"
        case 30:
            return "var"
        case 40:
            return "any"
        case 41:
            return "err"
        case 42:
            return "flg"
    }
    
    return "unk"
}

func M_StrToType(Input string) (int) {
    switch Input {
        case "non":
            return 0
        case "boo":
            return 1
        case "int":
            return 10
        case "flo":
            return 11
        case "num":
            return 12
        case "str":
            return 20
        case "var":
            return 30
        case "any":
            return 40
        case "err":
            return 41
        case "flg":
            return 42
    }
    
    return -1
}
