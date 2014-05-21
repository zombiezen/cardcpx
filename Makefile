# Copyright 2014 Google Inc. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

JSCOMPILER=java -jar closure-compiler.jar
ANGULAR=angular-1.2.15
GOTOOL?=go

JS_MIN_LIB=\
	$(ANGULAR)/angular.min.js \
	$(ANGULAR)/angular-cookies.min.js \
	$(ANGULAR)/angular-route.min.js \
	ui/ui-bootstrap.min.js

JS_LIB=\
	$(ANGULAR)/angular.js \
	$(ANGULAR)/angular-cookies.js \
	$(ANGULAR)/angular-route.js \
	ui/ui-bootstrap.min.js

JS_SRC=\
	ui/base.js \
	ui/components/takes/takes.js \
	ui/components/filesize/filesize-filter.js \
	ui/components/filesize/filesize.js \
	ui/components/autofillbtn/autofillbtn.js \
	ui/components/allbox/allbox.js \
	ui/views/navbar/navbar-controller.js \
	ui/views/navbar/navbar.js \
	ui/views/takelist/takelist-controller.js \
	ui/views/takelist/takelist.js \
	ui/views/takeeditor/takeeditor-controller.js \
	ui/views/takeeditor/takeeditor.js \
	ui/views/importer/importer-service.js \
	ui/views/importer/importer-controller.js \
	ui/views/importer/importer.js \
	ui/app.js

all: cardcpx ui/js.js

cardcpx: *.go */*.go
	$(GOTOOL) build -o $@

ui/js.js: $(JS_MIN_LIB) build/compiled_js.js
	cat $^ > $@

build/compiled_js.js: $(JS_SRC) | build
	$(JSCOMPILER) \
	    --angular_pass \
	    --compilation_level=ADVANCED_OPTIMIZATIONS \
	    --closure_entry_point=cardcpx.module \
	    --externs=$(ANGULAR)/externs.js \
	    --generate_exports \
	    --remove_unused_prototype_props_in_externs=false \
	    --export_local_property_definitions \
	    --js_output_file=$@ \
	    --property_renaming=OFF \
	    $(JS_SRC)

build/uncompiled_js.js: $(JS_LIB) $(JS_SRC) | build
	cat $^ > $@

build:
	mkdir $@

clean:
	rm -rf build
	rm -f cardcpx ui/js.js

.PHONY: all clean
