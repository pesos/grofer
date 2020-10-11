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

//ErrorMsg displays cute error ASCI Art for PID Errors
func ErrorMsg() {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(7) //gives a pseudorandom number in the range [0, n) (n not included).
	switch num {
	case 0:
		errorDoggo()
	case 1:
		errorCatto()
	case 2:
		errorBunny()
	case 3:
		errorDolphy()
	case 4:
		errorOwl()
	case 5:
		errorBeaver()
	case 6:
		errorElephant()
	}
}

func errorDoggo() {
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

func errorOwl() {
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

func errorCatto() {
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

func errorDolphy() {
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

func errorBunny() {
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

func errorBeaver() {
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

func errorElephant() {
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
