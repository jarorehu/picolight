## Zdroje pro projekt
https://www.root.cz/clanky/vyuziti-tinygo-pri-programovani-raspberry-pi-pico-od-gpio-az-k-pwm/#k05

## build
tinygo build -o output.uf2 -target=pico soubor.go

## Popis
Použito je Raspberry PI pico, součástky jsou rozmístěné na dev desce a pospojované podle pozice po obou stranách desky. Součástky jsou rozmístěné tak, aby tlačítka a informativní LED byly na jedné straně a zbytek součástek na straně druhé. Napájení z 12V zdroje, pomocí LM7805 je odvozeno 5V pro napájení pico. 12V používá externí rampa s LED a informativní LED na desce.

### součástky
- Raspberry PI pico, napájené piny, do desky je nasouvané 
- LM7805 + cap 10u, 100u
- Spínací mosfet BUZ71A https://bg-electronics.de/datenblaetter/Transistoren/BUZ71A.pdf
- předřazené spínání NPN BC547 https://www.onsemi.com/download/data-sheet/pdf/bc550-d.pdf
- rezistory 10k pro omezení produ na pinu a pro pullup na mosfet base 
- 4x informativní LED, předřadné rezistory 1k
- konektory, mikro tlačítka

Pico řídí 4 kanály - R/G/B/W v několika stupních. Tlačítka +/- jsou společná, dalšími čtyřmi se volí jednotlivé barvy (short press), nebo přepíná režim (long press)

Aplikace je napsaná v Tinygo

## minicom
### výstup na minicom
sudo dmesg
sudo minicom -D /dev/ttyACM0

### pro spuštění minicom bez sudo:
sudo usermod -a -G dialout $USER

### minicom v terminálu
Ctrl+A a následně Q (Quit)
Ctrl+A, Z: Zobrazí nápovědu.
Ctrl+S: Pozastaví výstup.

