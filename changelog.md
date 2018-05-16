# Change log

## v3.0.1, 9 May 2018

Bug fix: Fixed a hardcoded path separator character.

## 3 Mar 2018

Sorry, but I rewrote some history in the interest of assigning sensible
semantic version numbers.  (This prompted by vgo.)

* I left the tag "v1" on the version of 31 Dec 2014 but set some new tags on
historical versions after that.

* I had previously identified 3 Dec 2016 as a completion date for some
significant breaking changes.  I'm now tagging this as v2.0.0 to recognize
the breaking changes.

* I'm now identifying 23 Jun 2017 as a date when some bug fixes and minor API
additions were complete.  For this I'm bumping the minor version to v2.1.0


## v3.0.0, 3 Mar 2018

This is an attempt to follow the guidelines of the initial presentation of
vgo.  There are no functional changes between v3.0.0 and v2.0.0 but this new
major version number is introduced to distinguish the new "module" structure
of putting major version v3 in a v3 subdirectory.


## v2.1.0, 22 Feb 2018

This release bumps the minor version for the few small changes and bug fixes
since the 3 Dec 2016 v2.0.0 version.  This version remains API compatible with
the 3 Dec 2016 version.


## v2.0.0 3 Dec 2016

Breaking changes, mostly driven by the sexagesimal formatting.  The formatting
routines are moved to an external package now.  They rely on types defined in
yet another external package.  The most significant change to the meeus
packages is that a large number of function parameters and return values
are changed from float64 to one of these externally defined types.  There
are four types for now, Angle, HourAngle, RA, and Time.  An advantage of
float64s is terseness and simplicity.  Defined types though have advantages
of type checking, clarity and readability, and consistent unit representation.

There are many quantities other than these four types and they remain as
float64s for now.


## v1 31 Dec 2014

Status:  No changes in the last year!  89 GitHub stars currently, otherwise
little feedback.  Some feedback suggested changes to improve algorithms or
wrap existing functions in a more convenient API.  I declined these in the
interest of keeping this library as close as practical to the algorithms in
the book.  Otherwise, bug reports or other issues are always welcome.

There were a few minor changes just after the v0.2 tag.  Also I do have a
rewrite of of the angle formatting stuff ready to merge now.  This changes
the angle formatting API so it's time for a new tag.  The v1 tag will reference
this code that has been stable for the last year, including the minor changes
just after v0.2, but before API breaking changes to the angle formatting.

Happy New Year!

## v0.2 2013-10-11

* All chapters now implemented.

Upcoming work will change the API a bit.

## v0.1 2013-09-29

Tag added with the library starting to get some attention.  Existing code is
pretty good shape but there is some work remaining.  Here's what I can think
of off the top of my head:

* Implement remaining chapters.  There are just a few left.

* Change the API in a few places.  As packages from later chapters have
imported earlier ones, I've seen a few little things I'd like to go back and
change.  The biggest change will be for angle formatting.  I'll keep the
basic idea but allow more flexibility, probably at the expense of adding
some API methods.

* Add correctness tests in addition to the examples from the book.  This is
possible in many cases by drawing on sources outside the book.

* Review everything for consistency in style and naming.  I learned a few
things as I went along and some of the earlier packages need updates.
