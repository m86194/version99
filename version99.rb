# 
# Version 99.0 Does Not Exist
# 
# 
# Released under the MIT License
# Copyright (c) 2008 Erik van Oosten http://www.day-to-day-stuff.blogspot.com/
# 
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
# 
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
# 
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.
#
# Version 1.0 2007-10-10 First public version.
# Version 1.1 2007-11-06 Jar is not completely empty as some tools don't like that.
# Version 1.2 2008-02-09 Jar now has a valid manifest.
# Version 1.3 2009-07-24 Supports artifacts with group id that contains dots.
# Version 2.0 2009-10-17 Removed metadata support, some proxies seem to barf on it.
#

require 'digest/sha1'
require 'digest/md5'

Camping.goes :Version99

index_html = "<html>
<head>
<title>Version 99 Does Not Exist</title>
</head>
<body>
<h1>Version 99 Does Not Exist</h1>
<p>Please see <a href="http://day-to-day-stuff.blogspot.com/2007/10/announcement-version-99-does-not-exist.html">my blog</a> to read why I created Version 99 Does Not Exist and its predecessor no-commons-logging.</p>
<p>Version 99 Does Not Exist emulates a Maven 2 repository and serves empty jars for any valid package that has version number <i>99.0-does-not-exist</i>. It also generates poms, <span style="text-decoration: line-through">metadata files</span> (removed since 2.0) and of course the appropriate hashes.</p>
<p>For example the following links will give an <a href="http://no-commons-logging.zapto.org/mvn2/commons-logging/commons-logging/99.0-does-not-exist/commons-logging-99.0-does-not-exist.jar">empty jar</a>, its <a href="http://no-commons-logging.zapto.org/mvn2/commons-logging/commons-logging/99.0-does-not-exist/commons-logging-99.0-does-not-exist.pom">pom</a> and the <a href="http://no-commons-logging.zapto.org/mvn2/commons-logging/commons-logging/maven-metadata.xml"><span style="text-decoration: line-through">maven metadata</span></a> for commons-logging.</p>
<p><a href="version99.rb">Vesion 99 Does Not Exist source code</a> (rb file, 4Kb, MIT license). See <a href="http://day-to-day-stuff.blogspot.com/2007/10/announcement-version-99-does-not-exist.html">my blog</a> for installation instructions.</p>
</body>
</html>
"

module Version99
  # An (almost) empty uncompressed jar, coded in Base64
  # Really empty jar: EMPTY_JAR = "UEsFBgAAAAAAAAAAAAAAAAAAAAAAAA==".unpack("m*").first
  EMPTY_JAR = "UEsDBAoAAAAAAME+SDiyfwLuGQAAABkAAAAUAAQATUVUQS1JTkYvTUFOSUZFU1QuTUb+ygAATWFuaWZlc3QtVmVyc2lvbjogMS4wDQoNClBLAQIKAAoAAAAAAME+SDiyfwLuGQAAABkAAAAUAAQAAAAAAAAAAAAAAAAAAABNRVRBLUlORi9NQU5JRkVTVC5NRv7KAABQSwUGAAAAAAEAAQBGAAAATwAAAAAA".unpack("m*").first

  # A URL for more information.
  INFO_URL = "http://day-to-day-stuff.blogspot.com/2007/10/announcement-version-99-does-not-exist.html"

  module Controllers
    class Jar < R '/mvn2/(.+)/([^/]+)/99\.0-does-not-exist/\2-99\.0-does-not-exist\.jar(\.sha1|\.md5)?'
      def get(*args)
        group_id, artifact_id, hash = *args
        @headers['Content-Type'] = hash ? "text/plain" : "application/java-archive"
        digest_by_hash(EMPTY_JAR, hash)
      end
    end
    
    class Pom < R '/mvn2/(.+)/([^/]+)/99\.0-does-not-exist/\2-99\.0-does-not-exist\.pom(\.sha1|\.md5)?'
      def get(*args)
        group_id, artifact_id, hash = *args
        @headers['Content-Type'] = hash ? "text/plain" : "text/xml"
        p = pom s2d(group_id), artifact_id
        digest_by_hash(p, hash)
      end
    
      private
      def pom(group_id, artifact_id)
        "<?xml version=\"1.0\" encoding=\"iso-8859-1\"?><project><modelVersion>4.0.0</modelVersion><groupId>#{group_id}</groupId><artifactId>#{artifact_id}</artifactId><name>#{artifact_id}</name><version>99.0-does-not-exist</version><description>This is a generated pom. Version 99.0-does-not-exist is a dummy implementation with actually does nothing and has no dependencies. VERSION 99.0-does-not-exist IS NOT IN ANY WAY AFFILIATED WITH THE ORIGINAL DEVELOPERS of #{group_id}.</description><url>#{INFO_URL}</url></project>"
      end
    end
 
#    class Metadata < R '/mvn2/(.+)/([^/]+)/maven-metadata.xml(\.sha1|\.md5)?'
#      def get(*args)
#        group_id, artifact_id, hash = *args
#        @headers['Content-Type'] = hash ? "text/plain" : "text/xml"
#        m = metadata s2d(group_id), artifact_id
#        digest_by_hash(m, hash)
#      end
#    
#      private
#      def metadata(group_id, artifact_id)
#        last_updated = File.stat(__FILE__).mtime.strftime("%Y%m%d%H%M%S01")
#        "<?xml version=\"1.0\" encoding=\"iso-8859-1\"?><metadata><groupId>#{group_id}</groupId><artifactId>#{artifact_id}</artifactId><version>99.0-does-not-exist</version><versioning><versions><version>99.0-does-not-exist</version></versions><lastUpdated>#{last_updated}</lastUpdated></versioning></metadata>"
#      end
#    end

    class NotFound
      def get(p)
        @status = 404
        html do
          h1 "Version 99 Does Not Exist (Error 404)"
          h2 "Not Found: #{p}"
          p do
            span do
              a "Version 99 Does Not Exist", :href => INFO_URL
            end
            span " is a virtual Maven2 repository. It generates jars and poms for for any artifact with version '99.0-does-not-exist'."
          end
        end
      end
    end

    private
    def digest_by_hash(str, hash)
      hash ? ((hash == ".sha1") ? Digest::SHA1.hexdigest(str) : Digest::MD5.hexdigest(str)) : str
    end
    
    def s2d(str)
      str.gsub('/', '.')
    end
  end
end

