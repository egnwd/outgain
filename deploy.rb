#!/usr/bin/env ruby

require 'platform-api'
require 'excon'

slug_path = ARGV[0]
slug_archive = File.open(slug_path)

heroku = PlatformAPI.connect_oauth(ENV['HEROKU_OAUTH_TOKEN'])

puts %{Creating slug ...}
slug = heroku.slug.create("outgain", :process_types => { :web => "./outgain" })
puts %{Slug created with id "#{slug["id"]}" !}

puts %{Uploading slug archive "#{slug_path}" to "#{slug["blob"]["url"]}" ...}
Excon.put(slug["blob"]["url"], :body => slug_archive)
puts %{Slug uploaded !}

puts %{Releasing slug ...}
heroku.release.create("outgain", :slug => slug["id"])
puts %{Slug released !}
