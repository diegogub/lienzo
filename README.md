Goals
====
Lienzo is a small program to load Go templates in defined endpoints.
It can add data to it templates with a configuration file. 

Could be useful for teams willing to separate  designers and backend concerns.
Designers and frontend developers don't have to run all dbs to work with
templates.

Basic Usage
==========

- Compile it and drop the binary with the configuration .json in your
  template project
- Then run ./lienzo , and visit localhost:8989/{yourendpoint}

Defaults:
========
~~~
-config="lienzo.json"
-dir="."
-port="8989"
-suffix=".html"
~~~
