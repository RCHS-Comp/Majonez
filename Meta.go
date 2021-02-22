package main

// (c) Harry Nelsen 2020

type Meta struct {
    Name string
    Arguments int
    InputTypes []int
    OutputType int
}

func (Input *Meta) Reset() {
    Input.Name = "nop"
    Input.Arguments = 0
    Input.InputTypes = nil
    Input.OutputType = 0
}

//////////////////////////////////////////
// Current variable types:              //
//  0: nothing                          //
//  1: string                           //
//  2: int                              //
//  3: float                            //
//  4: int or float                     //
//  5: error                            //
//  6: probably variable                //
//  7: anything                         //
// -1: unknown                          //
//////////////////////////////////////////

var nop = Meta{Name: "nop", Arguments: 0, InputTypes: nil, OutputType: 0}

var VariableAllowedChars = "abcdefghijklmnopqrstuvwxyz"
var VariableUnavailable = []string{"!", "_"}
var VariableNames = []string{"!", "_"}
var VariableTypes = []int{1, 1}

func (Input *Meta) Fix() {
    if (Input.Arguments == 0 && len(Input.InputTypes) != 0) {
        Input.InputTypes = nil
    }

    if len(Input.InputTypes) != Input.Arguments {
        Input.Arguments = len(Input.InputTypes)
    }
    
    for Index, Arg := range Input.InputTypes {
        if Arg == -1 {
            Input.InputTypes[Index] = 0
        }
    }
    
    if Input.OutputType < 0 {
        Input.OutputType = 0
    }
}

func (CheckWith *Meta) Compliant(Input Command) (bool) {
    if Input.Name != CheckWith.Name {
        return false
    }
    
    if len(Input.Arguments) != CheckWith.Arguments {
        return false
    }
    
    var VarType int
    
    for Index, Arg := range Input.Arguments {
    
        VarType = CMD_VarType(Arg)
        
        if VarType == 6 {
            if M_VariableExist(Arg) {
                VarType = M_VariableType(Arg)
            } else {
                return false
            }
        }
        
        if VarType == -1 {
            return false
        }
            
        if CheckWith.InputTypes[Index] != 7 {
        
            if CheckWith.InputTypes[Index] == 4 {   // I did this so that we could use different types, if we did !(=4 & In) then if either of those tripped, it would return
                if In_int([]int{2,3}, VarType) == -1 {
                    return false
                }
            } else if VarType != CheckWith.InputTypes[Index] {
                return false
            }
        }
    }
    
    return true
}

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
    if (In_string(VariableNames, Name) != -1 || In_string(VariableUnavailable, Name) != -1) {
        return false
    }
    
    if (Type < 0 || Type > 5) {
        return false
    }
    
    if (Outliers(VariableAllowedChars, Name) != 0 || len(Name) < 1) {
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
            return "str"
        case 2:
            return "int"
        case 3:
            return "flo"
        case 4:
            return "num"
        case 5:
            return "err"
        case 6:
            return "var"
        case 7:
            return "any"
    }
    
    return "unk"
}

func M_StrToType(Input string) (int) {
    switch Input {
        case "non":
            return 0
        case "str":
            return 1
        case "int":
            return 2
        case "flo":
            return 3
        case "num":
            return 4
        case "err":
            return 5
        case "var":
            return 6
        case "any":
            return 7
    }
    
    return -1
}
