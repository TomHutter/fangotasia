# Fangotasia

## Fantasia

In 1984 my first computer game was Fantasia on a CP/M computer, a text adventure written in BASIC. Since then I played a lot of games. But, as you may know, there is always a specific kind of romance for the first game ever played :-) At that time we hat 5 1/4" floppy disks to store data on. But time marched on and soon 5 1/4" drives were replaced by 3" drives and by now even the 3" drives are not much more than a myth. From time to time I searched the Internet for Fantasia without any luck and then considered it "lost forever".

But then this year (in 2020) a colleague of mine found it on [www.c64-wiki.de](https://www.c64-wiki.de/wiki/Fantasia) . It turned out, the game was written for C64 and has been ported to CP/M already then. 

Starting to play Fantasia again on a C64 emulator I soon remembered my experiences back in 1984. "What verbs can I use?" As back in 1984, I had to look them up in the code...

## Fangotasia

This gave me the idea to reprogram Fantasia in GO. On the one hand as a finger exercise and to learn GO programming and on the other hand to add some features as looking up verbs (not in the code).

Getting in contact with the guys from www.c64-wiki.de Klaus informed me, that the name Fantasia is a Trademark. To avoid any legal problems, I decided to rename my GO-Version to Fangotasia - yes, you will find fango in the game now :-)

There is also a map feature now in the game. Of course you have to find it first ;-)

Fantasia was originally written in German. Now you an switch language between English and German (currently). At the beginning when no language is set and by entering `lang` in the game.

## Binaries

Currently you can build an `x86`and an `arm`version by using `make`.

* `make fangotasia.x86`will build a x86 version for the Linux operating system.

* `make fangotasia.arm`will build the arm version for Android

* `make release`will build both versions an pack them with the necessary config files to a tar file. 

  The Linux version runs obviously on Linux. The Android version can be run in Android Termux. Both release packages, with all config files included can be found on [GitHub](https://github.com/TomHutter/fangotasia/releases). 