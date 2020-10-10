/*
Copyright © 2020 The PES Open Source Team pesos@pes.edu

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func ErrorMsg() {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(6)
	switch num {
	case 0:
		ErrorDoggo()
	case 1:
		ErrorCatto()
	case 2:
		ErrorBunny()
	case 3:
		ErrorDolphy()
	case 4:
		ErrorOwl()
	case 5:
		ErrorBeaver()
	case 6:

	}
}

func ErrorDoggo() {
	dog := `
	pid no exist, done doggo a sad
		\
		 \
		/^-----^\
		V  o o  V
		 |  Y  |
		  \ ⌓ /
		  / - \
		  |    \
		  |     \     )
		  || (___\====
	`
	fmt.Println(dog)
}

func ErrorOwl() {
	goobes := `

   /\_/\  The council of wise owls are confused!  /\_/\
  ((@v@))      Please provide a valid PID!       ((@v@))
 ():::::()                                      ():::::()
   VV-VV          /\_/\         /\_/\             VV-VV
                 ((@v@))       ((@v@))
                ():::::()     ():::::()
                  VV-VV         VV-VV
  `
	fmt.Println(goobes)
}

func ErrorCatto() {
	cat := `
Catto says PID is invalid, plis give valid PID
    \
     \
  　██░▀██████████████▀░██
　　█▌▒▒░████████████░▒▒▐█
　　█░▒▒▒░██████████░▒▒▒░█
　　▌░▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒░▐
　　░▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒░
　 ███▀▀▀██▄▒▒▒▒▒▒▒▄██▀▀▀██
　 ██░░░▐█░▀█▒▒▒▒▒█▀░█▌░░░█
　 ▐▌░░░▐▄▌░▐▌▒▒▒▐▌░▐▄▌░░▐▌
　　█░░░▐█▌░░▌▒▒▒▐░░▐█▌░░█
　　▒▀▄▄▄█▄▄▄▌░▄░▐▄▄▄█▄▄▀▒
　　░░░░░░░░░░└┴┘░░░░░░░░░
　　██▄▄░░░░░░░░░░░░░░▄▄██
　　████████▒▒▒▒▒▒████████
　　█▀░░███▒▒░░▒░░▒▀██████
　　█▒░███▒▒╖░░╥░░╓▒▐█████
　　█▒░▀▀▀░░║░░║░░║░░█████
　　██▄▄▄▄▀▀┴┴╚╧╧╝╧╧╝┴┴███
　　██████████████████████`
	fmt.Println(cat)
}

func ErrorDolphy() {
	dolphy := `
                               _.-~  )
                    _..--~~~~,'   ,-/     _
                 .-'. . . .'   ,-','    ,' )
               ,'. . . _   ,--~,-'__..-'  ,'
             ,'. . .  (@)' ---~~~~      ,'
            /. . . . '~~             ,-'
           /. . . . .             ,-'
          ; . . . .  - .        ,'
         : . . . .      \_     /       PID did Dolphy a daze,
        . . . . .          \-.:        Please enter valid PID
       . . . ./  - .          )
      .  . . |  _____..---.._/ _____________
~---~~~~----~~~~             ~~
`
	fmt.Println(dolphy)
}

func ErrorBunny() {
	bunny := `             ,
            /|      __
           / |   ,-~ /
          Y :|  //  /
          | jj /( .^
          >-"~"-v"
         /       Y
        jo  o    |
       ( ~T~     j
        >._-' _./
       /   "~"  |
      Y     _,  |
     /| ;-"~ _  l
    / l/ ,-"~    \
    \//\/      .- \
     Y        /    Y     Bunny couldn't recognise that PID.
     l       I     !      Done bunny a confuse.
     ]\      _\    /"\     Please give bunny a valid PID.
    (" ~----( ~   Y.  )
~~~~~~~~~~~~~~~~~~~~~~~~~
  `

	fmt.Println(bunny)
}

//ErrorBeaver displays a PID Error Message with a beaver as ascii ART
func ErrorBeaver() {
	beaver := `
    /   \          /   \
    \_   \        /  __/
    _\   \      /  /__
    \___  \____/   __/
        \_       _/
          | @ @  \_
          |
        _/     /\
        /o)  (o/\ \_
        \_____/ /
          \____/
  Whoooopsssss, invalid PID. 
  Please enter a valid PID.
  `
	fmt.Println(beaver)
}

//ErrorElephant prints out an elephant to display an error message
func ErrorElephant() {
	elephant := `
                         ____
                    ---'-    \
      .-----------/           \
     /           (         ^  |   __
&   (             \        O  /  / .'
'._/(              '-'  (.   (_.' /
     \                    \     ./
     |    |       |    |/ '._.'
     )   @).____\|  @ |
 .  /    /       (    |           Pawoo. Pawoo. Pawoo!
\|, '_:::\  . ..  '_:::\ ..\).    Plz give Elephant a valid PID.
  `
	fmt.Println(elephant)
}
