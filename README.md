# ttype
Program for practicing touch typing. Written in Go

## Setup
This program requires [GNU Source-highlight](https://www.gnu.org/software/src-highlite/) to highlight files. Follow the instructions specified on the website or try the following:

Install required packages. There might be some more that are required, read the documentation for mor information.
```bash
sudo apt-get install texinfo libboost-all-dev
```

Now download and install GNU Source-highlight
```bash
git clone git://git.savannah.gnu.org/src-highlite.git
cd src-highlite
autoreconf -i
./configure
make
make install
```
