require 'rubygems'
require 'cucumber/rake/task'

Cucumber::Rake::Task.new("test:all",
    "Run all tests") do |t|
  t.cucumber_opts = %w{--format progress -t ~@pending}
end

task :default => "test:all"
