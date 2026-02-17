package netserver

import "strings"

const greetingBanner = `&p             |   |   ---   ^  ---   ^         |
             | \ |    |   |-|  |   |-|       -*-
             |   |  |_|   | |  |   | |        | 
&b           .: &Y*                &b.
             &p__&b .   &Y*    &b.          .: &Y*    &p__
            (())                           (())
           ((()))  __          &b.      &p__  ((()))
          (((())))(())     &b.&W         &p(())(((())))
         ((((()))))())) _         _ ((()((((()))))
           &W|____|&p((())))()   &C,   &p()(((()))&W|____|
           |_[]_| |__|&p((())&W__&zA__&p((())&W|__| |_[]_|
          _|    |_|[]|_|_|I-I-I-I|_|_|[]|_|    |_
         |-|    |-|  |-|||-I-I-I-|||-|  |-|    |-|
        (|-|    |-|  |-| |I-I-I-I| |-|  |-|    |-|)
       ((|-| __ |-|  |-| |-I-I-I-| |-|  |-| __ |-|))
       ()|-|_XX_|-|__|T|_|[T]-[T]|_|T|__|-|_XX_|-()
&g   ^-^^-^^-^^-^    ^-^  """   &O/   \   &g"""  ^-^    ^-^^-^^-^^-^
&w
Welcome to Njata Reborn!

The original Njata was written by Chieftain and Zoie and was based on 
SmaugFUSS 1.9 code by Samson and various members of the SMAUG community.
SMAUG 1.4 code by Thoric (Derek Snider) with Altrag, Blodkai, Haus, Narn,
Scryn, Swordbearer, Tricops, Gorog, Rennard, Grishnakh, Fireblade and Nivek.
MERC 2.1 code by Hatchet, Furey, and Kahn.
DikuMUD code by: Hans Staerfeldt, Katja Nyboe, Tom Madsen, 
Michael Seifert && Sebastian Hammer

Most enthusiastic greetings to you, traveler! You have stumbled across a
world called NJATA. It is our hope that you come to think of this place as a
welcome hearth to hang your hat, whether you hail from Colista, Lavada,
Dori, Ajax or Dogala. Feel free to have a look around or to join us
as a member of [&CCreation&w]. If you have any questions, please ask for
help once you have logged in!`

func WriteBanner(session *Session) {
    lines := strings.Split(greetingBanner, "\n")
    for _, line := range lines {
        session.WriteLine(line)
    }
    session.WriteLine("")
}
