package main

// (c) Harry Nelsen 2021

type Meta struct {
    Name string
    Arguments int
    InputTypes []int
    OutputType int
}

type Command struct {
    Name string
    Arguments []string
}

func (Input *Meta) Reset() {
    Input.Name = "nop"
    Input.Arguments = 0
    Input.InputTypes = nil
    Input.OutputType = 0
}

func (Input *Command) Reset() {
    Input.Name = "nop"
    Input.Arguments = nil
}

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
    if In_int(CheckWith.InputTypes, 99) != -1 {   // Special variables can only be output
        return false
    }

    if CheckWith.OutputType == -1 { // Commands cannot output an unknown variable
        return false
    }

    if Input.Name != CheckWith.Name {   // Check if the names match between the meta of a command and the input command are the same
        return false
    }
    
    if len(Input.Arguments) != CheckWith.Arguments {    // Check if the argument counts match between the meta of a command and the input command are the same
        return false
    }
    
    var VarType int // Prepare the output of the var detection
    
    for Index, Arg := range Input.Arguments {
    
        VarType = CMD_VarType(Arg)  // detect var type
        
        if VarType == 30 {  // If var type is another variable
            if M_VariableExist(Arg) {   // Check if variable exists
                VarType = M_VariableType(Arg)   // Get var type from the variable
            } else {    // The variable doesn't exist
                return false
            }
        }
        
        if VarType == -1 {  // See if var type is unknown
            return false
        }
            
        if CheckWith.InputTypes[Index] != 40 {  // If variable is not anything
            if CheckWith.InputTypes[Index] == 12 {   // I did this so that we could use different types, if we did !(=4 & In) then if either of those tripped, it would return
                if In_int([]int{10, 11}, VarType) == -1 {   // If var type is num, see if it's either int or flo
                    return false
                }
            } else if VarType != CheckWith.InputTypes[Index] {  // If it's not a num, see if the var types match
                return false
            }
        }
    }
    
    return true // Everything passed, we good
}
