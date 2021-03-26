package main

var HGF_V_Major = 1
var HGF_V_Minor = 1
var HGF_V_Micro = 0

// HGF version 1.1.0
// (c) Harry Nelsen 2021

var echo = Meta{Name: "echo", Arguments: 1, InputTypes: []int{1}, OutputType: 0}
var rem  = Meta{Name: "rem",  Arguments: 1, InputTypes: []int{7}, OutputType: 0}
var test = Meta{Name: "test", Arguments: 0, InputTypes: nil,      OutputType: 0}
