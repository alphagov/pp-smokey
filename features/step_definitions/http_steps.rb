Given /^I have the JSON data (.*)$/ do |json_data|
  @json_data = json_data
  add_header('Content-type', 'application/json')
end

Given /^I have the HTTP header "(.+?):\w*(.*)"$/ do |name, value|
  add_header(name, value)
end

When /^I (GET|POST|PUT|DELETE) (.*)$/ do |method, url|
  url = replace_env_host(url)
  @response = do_http_request(url, method.downcase.to_sym, options) do |response, request, result|
    response
  end
end

Then /^I should receive an HTTP (\d{3})$/ do |status|
  @response.code.to_i.should == status.to_i
end

Then /^I should receive an HTTP (30[12]) redirect to (.*)$/ do |status, url|
  url = replace_env_host(url)
  @response.code.to_i.should == status.to_i
  @response.headers[:location].should == url
end


def headers
  @headers
end

def add_header key, value
  @headers ||= []
  @headers << [key, value] 
end

def options
  options = {}
  options[:http_headers] = headers if headers
  options[:payload] = @json_data if @json_data
  options
end
