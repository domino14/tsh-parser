This repository parses [TSH](http://www.poslarchive.com/tsh/doc/all.html) data files and computes standings for players over many tournaments. The point system used for these standings is configurable.

#### Example run script:

DB_MIGRATIONS_PATH=file://../../migrations/ DB_PATH=mgi.db TOURNEY_SCHEMA_PATH=../../cfg/pts_mgi.csv ./tshparser