import os

DefaultMakefile = """#---------------------------------------------------------------------------------
.SUFFIXES:
#---------------------------------------------------------------------------------

ifeq ($(strip $(DEVKITARM)),)
$(error "Please set DEVKITARM in your environment. export DEVKITARM=<path to>devkitARM")
endif

include $(DEVKITARM)/ds_rules

#---------------------------------------------------------------------------------
# TARGET is the name of the output
# BUILD is the directory where object files & intermediate files will be placed
# SOURCES is a list of directories containing source code
# INCLUDES is a list of directories containing extra header files
# MAXMOD_SOUNDBANK contains a directory of music and sound effect files
#---------------------------------------------------------------------------------
TARGET		:=	$(shell basename $(CURDIR))
BUILD		:=	build
SOURCES		:=	source
DATA		:=	data  
INCLUDES	:=	include

#---------------------------------------------------------------------------------
# options for code generation
#---------------------------------------------------------------------------------
ARCH	:=	-mthumb -mthumb-interwork -march=armv5te -mtune=arm946e-s

CFLAGS	:=	-g -Wall -O2\
 		 -fomit-frame-pointer\
		-ffast-math \
		$(ARCH)

CFLAGS	+=	$(INCLUDE) -DARM9
CXXFLAGS	:= $(CFLAGS) -fno-rtti -fno-exceptions

ASFLAGS	:=	-g $(ARCH)
LDFLAGS	=	-specs=ds_arm9.specs -g $(ARCH) -Wl,-Map,$(notdir $*.map)

#---------------------------------------------------------------------------------
# any extra libraries we wish to link with the project (order is important)
#---------------------------------------------------------------------------------
LIBS	:= 	-lfat -lnds9
 
 
#---------------------------------------------------------------------------------
# list of directories containing libraries, this must be the top level containing
# include and lib
#---------------------------------------------------------------------------------
LIBDIRS	:=	$(LIBNDS)
 
#---------------------------------------------------------------------------------
# no real need to edit anything past this point unless you need to add additional
# rules for different file extensions
#---------------------------------------------------------------------------------
ifneq ($(BUILD),$(notdir $(CURDIR)))
#---------------------------------------------------------------------------------

export OUTPUT	:=	$(CURDIR)/$(TARGET)

export VPATH	:=	$(foreach dir,$(SOURCES),$(CURDIR)/$(dir)) \
					$(foreach dir,$(DATA),$(CURDIR)/$(dir))

export DEPSDIR	:=	$(CURDIR)/$(BUILD)

CFILES		:=	$(foreach dir,$(SOURCES),$(notdir $(wildcard $(dir)/*.c)))
CPPFILES	:=	$(foreach dir,$(SOURCES),$(notdir $(wildcard $(dir)/*.cpp)))
SFILES		:=	$(foreach dir,$(SOURCES),$(notdir $(wildcard $(dir)/*.s)))
BINFILES	:=	$(foreach dir,$(DATA),$(notdir $(wildcard $(dir)/*.*)))
 
#---------------------------------------------------------------------------------
# use CXX for linking C++ projects, CC for standard C
#---------------------------------------------------------------------------------
ifeq ($(strip $(CPPFILES)),)
#---------------------------------------------------------------------------------
	export LD	:=	$(CC)
#---------------------------------------------------------------------------------
else
#---------------------------------------------------------------------------------
	export LD	:=	$(CXX)
#---------------------------------------------------------------------------------
endif
#---------------------------------------------------------------------------------

export OFILES	:=	$(addsuffix .o,$(BINFILES)) \
			$(CPPFILES:.cpp=.o) $(CFILES:.c=.o) $(SFILES:.s=.o)
 
export INCLUDE	:=	$(foreach dir,$(INCLUDES),-I$(CURDIR)/$(dir)) \
			$(foreach dir,$(LIBDIRS),-I$(dir)/include) \
			$(foreach dir,$(LIBDIRS),-I$(dir)/include) \
			-I$(CURDIR)/$(BUILD)
 
export LIBPATHS	:=	$(foreach dir,$(LIBDIRS),-L$(dir)/lib)
 
.PHONY: $(BUILD) clean
 
#---------------------------------------------------------------------------------
$(BUILD):
	@[ -d $@ ] || mkdir -p $@
	@$(MAKE) --no-print-directory -C $(BUILD) -f $(CURDIR)/Makefile
 
#---------------------------------------------------------------------------------
clean:
	@echo clean ...
	@rm -fr $(BUILD) $(TARGET).elf $(TARGET).nds

#---------------------------------------------------------------------------------
else
 
#---------------------------------------------------------------------------------
# main targets
#---------------------------------------------------------------------------------
$(OUTPUT).nds	: 	$(OUTPUT).elf
$(OUTPUT).elf	:	$(OFILES)
 
#---------------------------------------------------------------------------------
%.bin.o	:	%.bin
#---------------------------------------------------------------------------------
	@echo $(notdir $<)
	$(bin2o)
 
-include $(DEPSDIR)/*.d
 
#---------------------------------------------------------------------------------------
endif
#---------------------------------------------------------------------------------------
"""

PythonImport = """import keyboard
import pygame
import sys
from pygame.locals import *

SpriteList = [[], []]
ColourList = []
CPU_Equ = ["X", "Z", "S", "A", "Q", "W", "=", "-", "up", "down", "left", "right"]
Scale = 1

ColourFinalised = False

def Sync(Equivalent, Scaling):
    global Scale
    global DispScale
    if type(Scaling) != int:
        print("The scale factor isn't an int, you'd be better off fixing that.\\n")
        Screen = pygame.display.set_mode((ScreenSizeX, ScreenSizeY), 0, 32)
        InitScreen = Screen
        Scale = 1
        pygame.display.set_caption('Window on PC with 1x scaling.')
    else:
        Screen = pygame.display.set_mode(((ScreenSizeX * Scaling), (ScreenSizeY * Scaling)), 0, 32)
        InitScreen = Screen
        Scale = Scaling
        pygame.display.set_caption('Window on PC with %sx scaling.' % str(Scaling))
    if type(Equivalent) == tuple or type(Equivalent) == list:
        if len(Equivalent) == 12:
            global CPU_Equ
            CPU_Equ = Equivalent
        else:
            print("The button length doesn't match, you may want to fix it.\\n")
    else:
        print("The button isn't an array, you may want to fix that.\\n")

def CheckKey(KeyName):
    DS_Buttons = ["KEY_A", "KEY_B", "KEY_X", "KEY_Y", "KEY_L", "KEY_R", "KEY_START", "KEY_SELECT", "KEY_UP", "KEY_DOWN", "KEY_LEFT", "KEY_RIGHT"]
    Output = CPU_Equ[DS_Buttons.index(KeyName)]
    try:
        ToReturn = keyboard.is_pressed(Output.lower())
        return ToReturn
    except OSError:
        print("If this is running on Mac OS X, you may need to run this as an admin to have proper outputs.\\n")
        return False

def MakeSprite(SpriteName, Size):
    SpriteList[0].append(SpriteName)
    SpriteList[1].append(Size)

def MakeColour(Colour, ColourNumber):
    if type(ColourList) == tuple and ColourFinalised != False:
        print("The colour list is finalised already, not even going to try.\\n")
    try:
        if type(Colour) != tuple or type(Colour) != list:
            if len(Colour) == 3:
                if len(ColourList) - 1 < ColourNumber:
                    while len(ColourList) - 1 < ColourNumber:
                        ColourList.extend("X")
                ColourList[ColourNumber] = Colour
            else:
                print("Colour is in the incorrect format, consider fixing this.\\n")
        else:
            print("Colour number isn't an array, you might want to fix this.\\n")
    except:
        if type(Colour) == tuple:
            print("Couldn't add another colour, it is most likely finalised.\\n")
        else:
            print("The colour couldn't be added for an unknown reason.\\n")

def FinaliseColours(List):
    ColourList = tuple(List)
    ColourFinalised = True

def CheckExit():
    for event in pygame.event.get():
        if event.type == QUIT:
            pygame.quit()
            sys.exit()

def SpriteDraw(Name, Location, Colour, HiddenMode):
    if HiddenMode != True:
        if type(Colour) != int:
            print("The colour isn't an int, you probably want to change that.\\n")
        elif len(ColourList) - 1 < Colour:
            print("The colour doesn't exist yet, you might want to change that.\\n")
        elif type(ColourList[Colour]) == tuple or type(ColourList[Colour]) == list and len(ColourList[Colour]) == 3:
            if type(Location) not in [tuple, list]:
                print("Location notation is wrong, you may want to fix this.\\n")
            else:
                if type(Location[0]) != int or type(Location[1]) != int or len(Location) != 2:
                    print("Location notation is wrong, you probably want to fix this.\\n")
                else:
                    if Name not in SpriteList[0]:
                        print("Sprite probably doesn't exist yet, you may want to check your spelling.\\n")
                    else:
                        Size = SpriteList[1][SpriteList[0].index(Name)]
                        Size_X = int((Size.split("x"))[0])
                        Size_Y = int((Size.split("x"))[1])
                        X = Location[0]
                        Y = Location[1]
                        pygame.draw.rect(Screen, ColourList[Colour], (((X * Scale), (Y * Scale)), ((Size_X * Scale), (Size_Y * Scale))))
            

def RefreshScreen():
    pygame.display.update()
    fpsClock.tick(FPS)
    

# DS Settings
ScreenSizeX = 256
ScreenSizeY = 192

if __name__ != "__main__":
    fpsClock = pygame.time.Clock()
    pygame.display.init()
    Screen = pygame.display.set_mode(((ScreenSizeX * Scale), (ScreenSizeY * Scale)), 0, 32)
    InitScreen = Screen
    pygame.display.set_caption('Window on PC with %sx scaling.' % str(Scale))
    FPS = 60


"""

def ConvertRGB(ColourHex, Mode15=False):
    if type(ColourHex) == str and len(ColourHex) == 6:
        Part1 = ColourHex[:2]
        Part2 = ColourHex[2:-2]
        Part3 = ColourHex[-2:]
        Dec1 = int(Part1, 16)
        Dec2 = int(Part2, 16)
        Dec3 = int(Part3, 16)
        if Mode15 != False:
            Const = (31 / 255)
        else:
            Const = 1
        Output1 = (int(Const * Dec1))
        Output2 = (int(Const * Dec2))
        Output3 = (int(Const * Dec3))
        if Mode15 != False:
            return ('RGB15(%s, %s, %s)' % (Output1, Output2, Output3))
        else:
            return ('(%s, %s, %s)' % (Output1, Output2, Output3))
    else:
        if Mode15 != False:
            return 'RBG15(31,31,31)'
        else:
            return ('(255, 255, 255)')

def CheckSize(Size):
    Powers = [(2**3), (2**4), (2**5), (2**6)]
    if type(Size) == str:
        try:
            if int((Size.split('x'))[0]) in Powers and int((Size.split('x'))[1]) in Powers:
                if int((Size.split('x'))[0]) <= 64 and int((Size.split('x'))[1]) <= 64:
                    return True
                else:
                    return 'Size is larger then 64, consider fixing this.\n'
            else:
                return 'Size is not a power of 2, consider fixing this.\n'
        except:
            return 'Size is not an integer, consider fixing this.\n'
    else:
        return "Size isn't even in the correct format, you'd probably want to fix this.\n"

def MakeIf(Statement, Not=False, PythonMode=True):
    Statement = Statement.lower()
    Statement = Statement.replace('==', '=')
    Statement = Statement.replace('!=', '=')
    
    if PythonMode != True:
        Statement = Statement.replace('or', '||')
        Statement = Statement.replace('and', '&&')
    if Not == False:
        Statement = Statement.replace('=', '==')
        Statement = Statement.replace('oppo', '!=')
        return Statement
    else:
        Statement = Statement.replace('=', '!=')
        Statement = Statement.replace('oppo', '==')
        return Statement

def ConvertFromHGF(Script, PythonMode=True, ForceCreate=False):
    TypesAllowed = ['int', 'double', 'float']
    Output = []
    IndentCount = 0
    Setup = False
    Refresh = False
    StartedLoops = 0
    EndedLoops = 0
    Error = 0
    Colours = []
    Sprites = []
    ExistingVars = []
    SpriteSize = []
    KeyCheck = False
    ExtPalette = False
    IfWarning = False
    MathWarning = False
    if ForceCreate == True:
        print('This program will be *forced* to be created no matter how many errors there are.')
        print('It can and will be bound to fail, by doing this you know the repercussions.')
        print("If you're doing this because normally your program won't compile, it might be an error on your part.\n")
    if type(Script) == list or type(Script) == tuple:
        for Line in Script:
            Indent = ('    ' * IndentCount)
            Line = Line.replace(Indent, '')
            Line = Line.strip()
            Line = Line.replace('%indent% ', Indent)
            Line = Line.replace('%indent%', Indent)
            Line = Line.replace('%screenx%', '256')
            Line = Line.replace('%screeny%', '192')
            if (Line[:7].lower()) == 'pymode:' and PythonMode == True:
                Output.append(Line[7:])
            elif (Line[:7].lower()) == 'dsmode:' and PythonMode != True:
                Output.append(Line[7:])
            elif Line.lower() == 'setup':
                if Setup == False:
                    Setup = True
                    if PythonMode == True:
                        Output.append(Indent + 'from RCHS_Comp import *')
                        Output.append(Indent + '# Generated from HGF')
                        Output.append(Indent + 'print("HGF and the HGF converter were\\ncreated by Harry Nelsen\\n")')
                        Output.append(Indent + 'Scale = 1\n')
                        Output.append(Indent + 'DS_Buttons = ["A", "B", "X", "Y", "L", "R", "START", "SELECT", "UP", "DOWN", "LEFT", "RIGHT"]')
                        Output.append(Indent + 'CPU_Equ = ["X", "Z", "S", "A", "Q", "W", "=", "-", "up", "down", "left", "right"]\n')
                        Output.append(Indent + 'Sync(CPU_Equ, Scale)\n')
                    else:
                        Output.append(Indent + '#include <nds.h>')
                        Output.append(Indent + '#include <stdio.h>\n')
                        Output.append(Indent + '// Generated from HGF')
                        Output.append(Indent + 'int main(void) {')
                        IndentCount =+ 1
                        Indent = ('    ' * IndentCount)
                        Output.append('\n' + Indent + 'int i = 0;\n')
                        Output.append(Indent + 'PrintConsole bottomScreen;\n')
                        Output.append(Indent + 'videoSetMode(MODE_0_2D);')
                        Output.append(Indent + 'videoSetModeSub(MODE_0_2D);\n')
                        Output.append(Indent + 'vramSetBankA(VRAM_A_MAIN_SPRITE);')
                        Output.append(Indent + 'vramSetBankC(VRAM_C_SUB_BG);\n')
                        Output.append(Indent + 'vramSetBankF(VRAM_F_LCD);\n')
                        Output.append(Indent + 'oamInit(&oamMain, SpriteMapping_1D_32, true);\n')
                        Output.append(Indent + 'consoleInit(&bottomScreen, 3, BgType_Text4bpp, BgSize_T_256x256, 31, 0, false, true);')
                        Output.append(Indent + 'consoleSelect(&bottomScreen);\n')
                        Output.append(Indent + 'iprintf("HGF and the HGF converter were\\n");')
                        Output.append(Indent + 'iprintf("created by Harry Nelsen\\n\\n");')
                else:
                    print('Setup has already been called, consider removing it.\n')
            elif (Line[:5].lower()) == 'echo ':
                if PythonMode == True:
                    Output.append(Indent + 'print("' +Line[5:] + '")')
                else:
                    Output.append(Indent + 'iprintf("' +Line[5:] + '\\n");')
            elif Line.lower() == 'loop_st':
                StartedLoops += 1
                if PythonMode == True:
                    Output.append(Indent + 'while True:')
                    IndentCount += 1
                else:
                    Output.append(Indent + 'while(1) {')
                    IndentCount += 1
            elif Line.lower() == 'loop_nd':
                EndedLoops += 1
                if PythonMode == True:
                    IndentCount -= 1
                    Output.append('')
                else:
                    IndentCount -= 1
                    Indent = ('    ' * IndentCount)
                    Output.append(Indent + '}')
            elif Line.lower() == 'keycheck':
                if KeyCheck == False:
                    KeyCheck = True
                    if PythonMode == True:
                        Output.append('')
                    else:
                        Output.append('\n' + Indent + 'int keys;')
                        Output.append(Indent + 'scanKeys();')
                        Output.append(Indent + 'keys = keysHeld();\n')
                else:
                    print('Key checking has already been called, consider fixing it.\n')
            elif (Line[:10].lower()) == 'check_key:':
                if KeyCheck == True:
                    if (Line[10:].lower()) in ['a', 'b', 'x', 'y', 'l', 'r', 'select', 'start', 'up', 'down', 'left', 'right']:
                        StartedLoops += 1
                        if PythonMode == True:
                            Output.append(Indent + 'if CheckKey("KEY_' + (Line[10:].upper()) + '"):')
                            IndentCount += 1
                        else:
                            Output.append(Indent + 'if((keys & KEY_' + (Line[10:].upper()) + ')) {')
                            IndentCount += 1
                    else:
                        print("Key doesn't exist, consider fixing it.\n")
                else:
                    print('Key checking has not been called yet, consider fixing it.\n')
            elif Line.lower() == 'refreshscr':
                if Refresh == False:
                    Refresh = True
                    if PythonMode == True:
                        Output.append('\n' + Indent + 'RefreshScreen()\n')
                    else:
                        Output.append('\n' + Indent + 'swiWaitForVBlank();')
                        Output.append(Indent + 'oamUpdate(&oamMain);\n')
                else:
                    print('A refresh sequence has already been implemented, consider fixing this.\n')
            elif Line.lower() == 'exitcheck':
                if KeyCheck != False:
                    if PythonMode == True:
                        Output.append(Indent + 'CheckExit()')
                    else:
                        Output.append(Indent + 'if(keys & KEY_SELECT) break;\n')
                else:
                    Error += 1
                    print('Key checking has not been called yet, consider fixing it.\n')
            elif (Line[:8].lower()) == 'colmake:':
                if ExtPalette == False:
                    if StartedLoops > 0:
                        print('It is suggested to place this before the loop, you might want to fix it.\n')
                    if PythonMode == True:
                        Output.append((Indent + 'MakeColour(%s, %s)') % ((ConvertRGB((Line[8:].lower().split('|'))[1], False)), ((Line[8:].lower().split('|'))[0])))
                        Colours.append((Line[8:].lower().split('|'))[0])
                    else:
                        Output.append((Indent + 'VRAM_F_EXT_SPR_PALETTE[%s][1] = %s;') % (((Line[8:].lower().split('|'))[0]), ConvertRGB((Line[8:].lower().split('|'))[1], True)))
                        Colours.append((Line[8:].lower().split('|'))[0])
                else:
                    print('Palette already ended.\n')
            elif (Line[:8].lower()) == 'sprmake:':
                if StartedLoops > 0:
                    print('It is suggested to place this before the loop, you might want to fix it.\n')
                if PythonMode == True:
                    if CheckSize((Line[8:].lower().split('|'))[1]) == True:
                        Output.append((Indent + 'MakeSprite("%s", "%s")') % (((Line[8:].lower().split('|'))[0]), ((Line[8:].lower().split('|'))[1])))
                        Sprites.append((Line[8:].lower().split('|'))[0])
                        SpriteSize.append((Line[8:].lower().split('|'))[1])
                    else:
                        print(CheckSize((Line[8:].lower().split('|'))[1]))
                else:
                    if CheckSize((Line[8:].lower().split('|'))[1]) == True:
                        Sprites.append((Line[8:].lower().split('|'))[0])
                        SpriteSize.append((Line[8:].lower().split('|'))[1])
                        Output.append(('\n' + Indent + 'u16* %s = oamAllocateGfx(&oamMain, SpriteSize_%s, SpriteColorFormat_256Color);') % (((Line[8:].lower().split('|'))[0]), ((Line[8:].lower().split('|'))[1])))
                        Output.append(('\n' + Indent + 'for(i = 0; i < %s / 2; i++) {') % ((Line[8:].lower().split('|'))[1]).replace('x', ' * '))
                        IndentCount += 1
                        Indent = ('    ' * IndentCount)
                        Output.append((Indent + '%s[i] = 1 | (1 << 8);') % ((Line[8:].lower().split('|'))[0]))
                        IndentCount -= 1
                        Indent = ('    ' * IndentCount)
                        Output.append(Indent + '}\n')
                    else:
                        print(CheckSize((Line[8:].lower().split('|'))[1]))
            elif (Line[:8].lower()) == 'sprdraw:':
                if ExtPalette == True:
                    if ((Line[8:].lower().split('|'))[0]) in Sprites:
                        if (Line[8:].lower().split('|'))[2] in Colours:
                            if PythonMode == True:
                                State = (Line[8:].lower().split('|'))[3]
                                if State in ['false', 'true']:
                                    State = (State[:1]).upper() + (State[1:]).lower()
                                    Output.append((Indent + 'SpriteDraw("%s", (%s, %s), %s, %s)') % ((Line[8:].lower().split('|'))[0], (((Line[8:].lower().split('|'))[1]).split('x'))[0], (((Line[8:].lower().split('|'))[1]).split('x'))[1], (Line[8:].lower().split('|'))[2], State))
                                else:
                                    print('Invalid state, you might want to fix that.\n')
                            else:
                                Output.append('\n' + Indent + 'oamSet(&oamMain,')
                                IndentCount += 1
                                Indent = ('    ' * IndentCount)
                                Output.append((Indent + '%s,') % str(Sprites.index((Line[8:].lower().split('|'))[0])))
                                Output.append((Indent + '%s, %s,') % ((((Line[8:].lower().split('|'))[1]).split('x'))[0], (((Line[8:].lower().split('|'))[1]).split('x'))[1]))
                                Output.append(Indent + '0,')
                                Output.append((Indent + '%s,') % (Line[8:].lower().split('|'))[2])
                                Output.append((Indent + 'SpriteSize_%s,') % SpriteSize[Sprites.index((Line[8:].lower().split('|'))[0])])
                                Output.append(Indent + 'SpriteColorFormat_256Color,')
                                Output.append((Indent + '%s,') % (Line[8:].lower().split('|'))[0])
                                Output.append(Indent + '-1,')
                                Output.append(Indent + 'false,')
                                Output.append((Indent + '%s,') % (Line[8:].lower().split('|'))[3])
                                Output.append(Indent + 'false, false, false')
                                Output.append(Indent + ');')
                                IndentCount -= 1
                        else:
                            print("Colour doesn't exist yet, you probably might want to check that.\n")
                    else:
                        print("Sprite doesn't exist yet, consider checking that.\n")
                else:
                    print("Colours haven't been finalised, therefore sprites cannot be used.\n")
            elif Line.lower() == 'colend':
                if ExtPalette == False:
                    ExtPalette = True
                    if PythonMode == True:
                        Output.append(Indent + 'FinaliseColours(ColourList)')
                    else:
                        Output.append(Indent + 'vramSetBankF(VRAM_F_SPRITE_EXT_PALETTE);')
                else:
                    print('Palette already ended, consider removing it.\n')
            elif (Line[:3].lower()) == 'if ':
                if IfWarning == False:
                    print('Remember that strings look like "String" and variables like Variable.')
                    print('You cannot compare with different types.\n')
                    IfWarning = True
                if '=' in Line or 'oppo' in Line:
                    StartedLoops += 1
                    if PythonMode == True:
                        Output.append((Indent + 'if %s:') % MakeIf(Line[3:], False, True))
                        IndentCount += 1
                    else:
                        Output.append((Indent + 'if(%s) {') % MakeIf(Line[3:], False, False))
                        IndentCount += 1
            elif (Line[:7].lower()) == 'not_if ':
                if IfWarning == False:
                    print('Remember that strings look like "String" and variables like Variable.')
                    print('You cannot compare with different types.\n')
                    IfWarning = True
                if '=' in Line or 'oppo' in Line:
                    StartedLoops += 1
                    if PythonMode == True:
                        Output.append((Indent + 'if %s:') % MakeIf(Line[7:], True, True))
                        IndentCount += 1
                    else:
                        Output.append((Indent + 'if(%s) {') % MakeIf(Line[7:], True, False))
                        IndentCount += 1
                else:
                    print('Equality statement missing, you might want to fix it.\n')
            elif (Line[:7].lower()) == 'newvar ':
                if (Line[:7].lower()) != 'i':
                    if ((Line[7:].lower()).split('|'))[0] in TypesAllowed:
                        if ((Line[7:].lower()).split('|'))[1] not in ExistingVars:
                            if PythonMode == True:
                                Output.append((Indent + '%s = %s') % (((Line[7:].lower()).split('|'))[1], ((Line[7:].lower()).split('|'))[0]))
                                ExistingVars.append(((Line[7:].lower()).split('|'))[1])
                            else:
                                Output.append((Indent + '%s %s;') % (((Line[7:].lower()).split('|'))[0], ((Line[7:].lower()).split('|'))[1]))
                                ExistingVars.append(((Line[7:].lower()).split('|'))[1])
                        else:
                            print('Variable already exists, you might want fix it.\n')
                    else:
                        print('Variable not in the available types, you might want to fix it.\n')
                else:
                    print("The variable name 'i' isn't allowed, you probably want to change it.\n")
            elif (Line[:4].lower()) == 'var ':
                if (Line[4:].lower()) != 'i':
                    if ' = ' in Line:
                        if (((Line[4:].lower()).split(' = '))[0]) in ExistingVars and len((Line[4:].lower()).split(' = ')) == 2:
                            if PythonMode == True:
                                Output.append((Indent + '%s = %s') % (((Line[4:].lower()).split(' = '))[0], ((Line[4:].lower()).split(' = '))[1]))
                            else:
                                Output.append((Indent + '%s = %s;') % (((Line[4:].lower()).split(' = '))[0], ((Line[4:].lower()).split(' = '))[1]))
                        else:
                            print("Variable doesn't exist yet, you might want to fix it.\n")
                    else:
                        print('Equality statement missing, consider fixing it.\n')
                else:
                    print("The variable name 'i' isn't allowed, you might want to change it.\n")
            elif Line.lower() == 'else':
                EndedLoops += 1
                StartedLoops += 1
                IndentCount -= 1
                Indent = ('    ' * IndentCount)
                if PythonMode == True:
                    Output.append(Indent + 'else:')
                else:
                    Output.append(Indent + '} else {')
                IndentCount += 1
            elif Line == '' or Line == '\n':
                pass
            elif Line[:1] == '#':
                if PythonMode == True:
                    Output.append((Indent + '#%s') % Line[1:])
                else:
                    Output.append((Indent + '//%s') % Line[1:])
            elif Line[:2] == '//':
                if PythonMode == True:
                    Output.append((Indent + '#%s') % Line[2:])
                else:
                    Output.append((Indent + '//%s') % Line[2:])
            elif Line[:4] == 'rem ':
                if PythonMode == True:
                    Output.append((Indent + '#%s') % Line[3:])
                else:
                    Output.append((Indent + '//%s') % Line[3:])
            else:
                print('Command not found, it either may be incorrectly spelled or just not exist.')
                print('You may want to check your spelling.\n')
                    
    if IndentCount == 1:
        if PythonMode != True:
            Output.append('\n' + Indent + 'return 0;')
            IndentCount = 0
            Output.append('}')
    elif IndentCount > 1:
        Error += 1
        print('The indentation count seems incorrect, a loop may not be closed.\n')
    if StartedLoops == 0:
        print('No loop detected, the program will end right after you start it.')
        print('Consider adding it.\n')
    if StartedLoops - EndedLoops != 0:
        Error += 1
        if StartedLoops > EndedLoops:
            print('A loop has been started without it being closed, consider fixing it.\n')
        else:
            print('A loop has been ended without being started, consider fixing it.\n')
    if Setup == False:
        Error += 1
        print('Missing setup sequence, consider fixing it.\n')
    if Refresh == False:
        Error += 1
        print('Missing refresh sequence, consider fixing it.\n')
    if len(Sprites) > 128:
        Error += 1
        print('Too many sprites, consider removing a few.\n')
    if len(Colours) > 32:
        print('Too many colours, consider removing a few.\n')
    print('Finished with ' + str(Error) + ' error(s)')
    if Error == 0 or ForceCreate == True:
        return Output
    else:
        return []


def ReadFile(FileName, PythonMode=True):
    try:
        FileObj = open(FileName, 'r')
        FileContents = FileObj.readlines()
        FileObj.close()
        return ConvertFromHGF(FileContents, PythonMode)
    except Exception as e:
        print('Catastrophic error: ' + str(e))
        return []

def WriteFile(FileName):
    print('\n===== Python log =====\n\n')
    PythonMode = ReadFile(FileName + '.hgf', True)
    print('\n===== DS log =====\n\n')
    DS = ReadFile(FileName + '.hgf', False)
    if os.path.exists(FileName) and os.path.isfile(FileName) != True:
        if not os.path.exists(FileName + '/source'):
            os.makedirs(FileName + '/source')
    else:
        os.makedirs(FileName + '/source')
    Makefile = open(FileName + '/Makefile', 'w')
    Makefile.write(DefaultMakefile)
    Makefile.close()
    Write = open(FileName + '/source/' + FileName + '.cpp', 'w')
    for Command in DS:
        Write.write(Command + '\n')
    Write.close()
    Write = open(FileName + '.py', 'w')
    for Command in PythonMode:
        Write.write(Command + '\n')
    Write.close()
    Write = open('RCHS_Comp.py', 'w')
    Write.write(PythonImport)
    Write.close()
    print('Done.')
            
    

WriteFile('Gen')
