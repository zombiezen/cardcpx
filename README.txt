This is a rather raw dump of the project as it exists on my machine.

Things of note:

* httputil and netutil are vendored in from Camlistore.  Modifications are
  purely to reduce dependencies on Camlistore.
* bootstrap is vendored in, unmodified.
* angular-ui bootstrap is vendored in <http://angular-ui.github.io/bootstrap/>
  unmodified.
* The closure-compiler is vendored in as a JAR.
* There's a dependency on AngularJS (not included here)

Final distribution would not include the dependencies, but would include a script
that could download the dependencies.
