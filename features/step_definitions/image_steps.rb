require 'chunky_png'

Then /^the image should be between (\d+) and (\d+) pixels wide$/ do |lower_bound, upper_bound|
  image = ChunkyPNG::Image.from_blob(@response)
  image.width.should > lower_bound.to_i
  image.width.should < upper_bound.to_i
end

Then /^the image should be between (\d+) and (\d+) pixels high$/ do |lower_bound, upper_bound|
  image = ChunkyPNG::Image.from_blob(@response)
  image.height.should > lower_bound.to_i
  image.height.should < upper_bound.to_i
end
