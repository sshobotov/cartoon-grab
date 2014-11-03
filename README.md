Cartoon grabber
============

Simple tool to grab cartoon image series from resources with page-by-page viewing model. Written in go language and allows to save images to single PDF file.

Supported keys:
- ```-f``` Forces script processing next URLs even if one of imegs returns not valid response status (e.g. 404)
- ```-u``` Initial page URL, script load the page and parse it for image and next page URL
- ```-i``` XPath selector to grab image src from the page, if no images will be found script stops
- ```-l``` XPath selector to grab next page URL, if no URL will be found script stops
