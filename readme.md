# dat2xlsx

A **tele.dat** fájl tartalmát egy megadott dátumtól kezdve egy xlsx fájlba exportálja

### Használata

Nincs szükség telepítésre, bárhonnan indítható.
Indítás parancssorból windows 10 környezetben alapértelmezett paraméterekkel: 
Start menü -> cmd

```cmd

c:\path\to\dat2xlsx.exe

```

Ebben az esetben a **tele.dat** fájlnak a programmal azonos könyvtárban kell lennie 
és az aktuális dátum előtti egy havi rekordokat fogja exportálni a **datexport.xlsx** fájlba.


#### Parancssori kapcsolók

A -h kapcsoló a konzolra írja a további lehetőségeket:

```cmd

c:\path\to\dat2xlsx.exe -h

Usage of C:\path\to\dat2xlsx.exe:
  -begin 
        a megadott yymmdd dátumtól kezdje el az exportálást (default "201001")
  -in 
        a tele.dat fájl helye (default ".\path\to\tele.dat")
  -out 
        az xlsx export fájl helye (default ".\path\to\datexport.xlsx")
```

### Példa egy futtatásra

```cmd

c:\path\to\dat2xlsx.exe -in c:\path\to\tele.dat -out c:\path\to\export.xlsx -begin 201101

```

### cmd file
Lehet készíteni egy start.cmd, vagy start.bat kiterjesztésű fájlt, ami segítségével egy asztali ikonon keresztül is elindítható a program.
A cmd fájl így nézhet ki az előbbi példa alapján:

```cmd

@echo off

set datfile=c:\path\to\tele.dat
set expfile=c:\path\to\export.xlsx
set sdate=201101

c:\path\to\dat2xlsx.exe -in %datfile% -out %expfile% -begin %sdate%

```