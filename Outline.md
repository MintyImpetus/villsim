# An Outline Of What Is To Come

This is not an organized list of planned features, just some ideas I wrote down at one point.

## Planned Features

Users can see the recent achievements, and the stats of the other player, as if it is the reputation of their group. Users can also manipulate other people's view of them, by sending out messengers to specific areas, and if they interact with another group, it changes to reflect what the person said.

Make some equation to find the chance they accept the new info, something like the difference between current and new, multiplied by something, multiplied by the charisma level of the messenger, probably added later, to find the percent.

This can combine with the fact that villages or locations can be "subscribed" to, whether through a spy, or raising relationship levels. The level of the spy or relationship dictates the info in the messages.
Maybe they shouldn't be spies. How about:
Spies and just normal messengers, scouts or normal workers somehow allied with you, will tell you obvious stuff, like enemies approaching, but maybe more secret stuff, added later could only be accessible by spies.

If stuff like inventions and the tech system is going to spread through interaction with other places, it seems that each place is going to need an algorithm that is capable of playing the game to an extent, so that inventions and interactions do not always have to be initiated by the players, especially if the scale does end up being many many places, too many to populate with just players.

### Title Screen

The title screen, for the terminal version, will open with an ascii art VillSim, with the version number below it, similar to the Doom Emacs splash screen.
It will allow you to input the port and ip of the server to connect to in a way more elegant then just a line by line readLine().
The title screen should probably be moved to inside the bubbletea View() function, so that as the user types, it can update the title screen to contain the ip and port below the logo, while being centered.
edit: The title screen does not need to be moved inside bubbletea. The ReadString() command has all that is required for an easy to understand ip and port input area.
If I want more fancy features in the title screen, it may be wise to use bubbletea for it, as well as https://github.com/charmbracelet/bubbles, for pre-built elements.

### Soldiers
Add the ability to buy barracks, that produce soldiers. Make it so you can send them to attack a village, hopefully win. Make it so you now control that village. Make it so you can somehow switch bases. Make it so the opponent can attack a base.
General:
Make it so you can just place soldiers on a village, defending from enemies, but not actually attacking, or have them attack.

### Trading
You are never given a direct knowledge of the number of a material in existence.
You can offer items to bid, but have to spread knowledge of the bid yourself, maybe some spreads naturally, if it is particularily exciting (possible feature).
The movements of people moving materials between places or people can be witnessed, giving you more knowledge of the market.
Materials can stolen from locations you have gotten.

## Metrics

Average village sizes:

Villages:
Tiny: 100
small: 500
medium: 1 - 2 k
Large: < 10k

Popularity:
Number of people travelling per turn (four hours)
Tiny < 1
Small < 4
Medium < 10
large < 50

Td = number of turns since event

Eventfullness = scale
1 - 3
Could expand later
1 boring. New average person
2 notable. New ship
3 major. bandits seen

May make it average popularity for people walking through villages in future

May make popularity a measure of the percent of villages population. At the same time this better mimics the fact that the number of people is the chance travellers are told, so I don't know.

Equation is:
Chance = Eventfullness * Td^Eventfullness * popularity * sizeofothervillage / distance / 10000 + Td / 2

Popularity = average population of the cities / 200

Average walking speed is 2km an hour or 0.5 a turn

Given that popularity isn't a percentage, it makes sense to not care about your village population, as newsworthiness is what effects whether people will spread it, once they are already there.
