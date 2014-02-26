require 'net/http'
require 'rest_client'
require 'cgi'

def head_request(url, options = {})
  do_http_request(url, :head, options)
end

def get_request(url, options = {})
  do_http_request(url, :get, options)
end

def do_http_request(url, method = :get, options = {}, &block)
  defaults = {
    :auth => true,
  }
  options = defaults.merge(options)

  ip_last_octet = rand(256)
  headers = {
    'User-Agent' => 'Smokey Test / Ruby',
    'Accept' => 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
    'X-Forwarded-For' => "10.0.0.#{ip_last_octet}",
    'True-Client-Ip' => "10.0.0.#{ip_last_octet}",
  }

  started_at = Time.now
  url = options[:cache_bust] ? cache_bust(url) : url
  if options[:auth]
    user     = ENV['AUTH_USERNAME']
    password = ENV['AUTH_PASSWORD']
  end
  if options[:client_auth]
    headers["Authorization"] = "Bearer #{ENV['BEARER_TOKEN']}"
    headers["Accept"] = "application/json"
  end
  if options[:host_header]
    headers["Host"] = options[:host_header]
  end

  RestClient::Request.new(
    url: url,
    method: method,
    user: user,
    password: password,
    headers: headers,
    payload: options[:payload]
  ).execute &block
rescue RestClient::Unauthorized => e
  raise "Unable to fetch '#{url}' due to '#{e.message}'. Maybe you need to set AUTH_USERNAME and AUTH_PASSWORD?"
rescue RestClient::Exception => e
  finished_at = Time.now
  message = ["Unable to fetch '#{url}'"]
  message += ["  Exception: '#{e}'"]
  message += ["  Response headers: #{e.response.headers.inspect if e.response}"]
  message += ["  Response time in seconds: #{finished_at - started_at}"]
  raise message.join("\n")
end
