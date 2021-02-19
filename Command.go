package main

import (
    "fmt"
    "strings"
    "strconv"
)

//////////////////////////////////////////
// Command Obj
//////////////////////////////////////////

type Command struct {
    Name string
    Arguments []string
}

func (Input *Command) Reset() {
    Input.Name = "nop"
    Input.Arguments = nil
}

//////////////////////////////////////////

//////////////////////////////////////////
// Command tools
//////////////////////////////////////////

func CMD_Convert(Input string) (Command) {  // Convert command with '{}' to obj
    var Output Command
    Output.Reset()
    
    if (Usage(Input, '{') != 1 || Usage(Input, '}') != 1) {
        Output.Name = Input
        Output.Arguments = nil
        return Output
    }
    
    L_First := Find(Input, '{', 1)
    L_Last := Find(Input, '}', 1)
    
    ToSplit := Input[L_First + 1:L_Last]
    
    Output.Name = Input[0:L_First]
    Output.Arguments = Split(ToSplit, ",")
    
    return Output
}

func CMD_FindStr(Input string) (string) {   // Find strings
    L_First := Find(Input, '"', 1)
    L_Last := Find(Input, '"', 2)
    
    RealString := Input[L_First + 1:L_Last]
    
    return RealString
}

func CMD_RestoreStr(Input string, Source []string) (string) {   // Convert string map to strings
    Output := Input // Idk if this will change anything, but it makes me feel better

    for i := 0; i < len(Source); i ++ { // Loop through everything in the string map
        Output = strings.Replace(Output, fmt.Sprintf("STR@%d", i), fmt.Sprintf("\"%s\"", Source[i]), -1)   // Replace
    }
    
    return Output
}

func CMD_ReadLine(Input string, GlobalStr []string) ([]Command, []string) {
    var Output []Command
    
    Input = strings.Replace(Input, " ,", ",", -1)   // Make it able to be separated by commas
    Input = strings.Replace(Input, ", ", ",", -1)
    
    if Usage(Input, '"') % 2 != 0 { // Make sure there is an even number of quotes
        return Output, GlobalStr
    }
    
    StrCount := 0   // String count for the string map
    var ToReplace string    // String that will be replaced to its place in the string map
    
    for (Usage(Input, '"') / 2) != 0 {  // Repeat until no strings left
        ToReplace = CMD_FindStr(Input)  // Locate a string
        GlobalStr = append(GlobalStr, ToReplace)    // Add the string to the string map
        Input = strings.Replace(Input, "\"" + ToReplace + "\"", fmt.Sprintf("STR@%d", StrCount), -1)    // Replace the string with the location in the string map
        StrCount ++ // Advance the string count
    }
    
    Input = StripSpace(Input)   // This will allow us to have \n in our strings THEN remove any unwanted endings
    Input = strings.ToLower(Input)  // Since all the strings have been removed, we can make them lowercase now
    Input = strings.Replace(Input, "str@", "STR@", -1)  // Need to make these uppercase again
    
    SplitCommands := strings.Split(Input, " ")  // Commands will be split by spaces, args split by commas
    for i := 0; i < len(SplitCommands); i ++ {  // Convert command strings to commands
        Output = append(Output, CMD_Convert(SplitCommands[i]))  // Add command to output
    }
    
    return Output, GlobalStr
}

func CMD_VarType(Input string) (int) {  // Get type of variable
    if IsInt(Input) {
        return 2
    } else if IsFloat(Input) {
        return 3
    }
    
    if len(Input) < 4 { // This is for things with letters under 4 characters
        if (Outliers(VariableAllowedChars, Input) == 0 || (Input == "!" || Input == "_")) {
            return 6
        }
    } else if Input[:4] == "STR@" { // For stuff in the string map
        return 1
    } else {    // Everything
        if Outliers(VariableAllowedChars, Input) == 0 {
            return 6
        }
    }
    
    return -1
}

func CMD_IsStrRef(Input string) (bool) {    // See if variable refers to the string map
    return (CMD_VarType(Input) == 1)
}

func CMD_StrInRange(Input string, Source []string) (bool) { // See if reference is in the string map
    if CMD_IsStrRef(Input) {
        StrNum, Error := strconv.Atoi(Input[4:])
        if Error != nil {
            return false
        }
        
        if (StrNum < len(Source) && StrNum >= 0) {
            return true
        }
    }
    
    return false
}

//////////////////////////////////////////
