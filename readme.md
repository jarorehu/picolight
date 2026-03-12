https://www.root.cz/clanky/vyuziti-tinygo-pri-programovani-raspberry-pi-pico-od-gpio-az-k-pwm/#k05


tinygo build -o output.uf2 -target=pico soubor.go




výstup na minicom
sudo dmesg
sudo minicom -D /dev/ttyACM0

pro spuštění minicom bez sudo:
udo usermod -a -G dialout $USER

Ctrl+A a následně Q (Quit)
Ctrl+A, Z: Zobrazí nápovědu.
Ctrl+S: Pozastaví výstup.

