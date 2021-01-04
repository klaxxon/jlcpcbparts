#!/bin/bash
wget https://jlcpcb.com/componentSearch/uploadComponentInfo -O parts.xls
ssconvert parts.xls parts.csv

