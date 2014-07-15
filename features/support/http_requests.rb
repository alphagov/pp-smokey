require 'net/http'
require 'rest_client'
require 'cgi'

def get_request(url, options = {}, &block)
  do_http_request(url, :get, options, &block)
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
  if options[:http_headers]
    options[:http_headers].each do |name, value|
      headers[name] = value
    end
  end

  RestClient::Request.new(
    url: url,
    method: method,
    user: user,
    password: password,
    headers: headers,
    payload: options[:payload],
    open_timeout: 10,
    timeout: 10,
    max_redirects: 0
  ).execute(&block)


  rescue RestClient::Exception => e
    e.response
  end
