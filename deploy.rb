#!/usr/bin/env ruby

# Executed on Circle once the build is succesful
# It creates a new slug on Heroku, uploads the corresponding archive,
# and releases it
# The first argument is the name of the app
# The second argument is the path to the slug archive
# The HEROKU_OAUTH_TOKEN environment variable must be set

require 'platform-api'
require 'excon'

app_name = ARGV[0]

slug_path = ARGV[1]
slug_archive = File.open(slug_path)

heroku = PlatformAPI.connect_oauth(ENV['HEROKU_OAUTH_TOKEN'])

puts %{Creating slug ...}
slug = heroku.slug.create(app_name, :process_types => {
    :web => "./start.sh"
})
puts %{Slug created with id "#{slug["id"]}" !}

puts %{Uploading slug archive "#{slug_path}" to "#{slug["blob"]["url"]}" ...}
Excon.put(slug["blob"]["url"], :body => slug_archive)
puts %{Slug uploaded !}

puts %{Releasing slug ...}
heroku.release.create(app_name, :slug => slug["id"])
puts %{Slug released !}
