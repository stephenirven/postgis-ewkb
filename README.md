Encoding / decoding PostGIS EWKB (Extended Well-Known Binary) formats in Go. 

I couldn't find much documentation for the encoding for these types, so it's mainly 
derived from experimentation.

Covers all the PostGIS types I could find in XY, XYZ, XYM, XYZM dimensions

* Point (Type 1)
* Line String (Type 2)
* Polygon (Type 3)
* Multi Point (Type 4)
* Multi Line String (Type 5)
* Multi Polygon (Type 6)
* Geometry Collection (Type 7)
* Circular String (Type 8)
* Compound Curve (Type 9)
* Curve Polygon (Type 10)
* Multi Curve (Type 11)
* Multi Surface (Type 12)

* Polyhedral Surface (Type 15)
* TIN (Type 16)
* Triangle (Type 17)

* Linear Ring (internal implementation structure for other types)







