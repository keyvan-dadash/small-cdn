# Small CDN
This was an experiment to gain insight to how a CDN does the minification process. This project contains auth apis that use JWT authentication, and it also provides apis that get CSS or JS files for the minification process.
Not to mention that it also logs all the minification process' stats, such as memory consumtion and duration of minification. Finally, it stores all the minified files under the /opt/{username} directory.
