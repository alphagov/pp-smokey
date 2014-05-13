def admin_base_url
  application_base_url('admin')
end

def app_domain
  ENV["PP_APP_DOMAIN"] || "preview.performance.service.gov.uk"
end

def full_app_domain
  ENV["PP_FULL_APP_DOMAIN"] || ENV["PP_APP_DOMAIN"] || "preview.performance.service.gov.uk"
end

def govuk_app_domain
  ENV["GOVUK_APP_DOMAIN"] || "preview.alphagov.co.uk"
end

def signon_base_url
  application_base_url('signon')
end

def application_base_url(app_name)
  case app_name
  when 'admin' then "https://admin.#{app_domain}/"
  when 'signon' then "https://signon.#{govuk_app_domain}/"
  else
    raise "Application '#{app_name}' not recognised, unable to boot it up"
  end
end

def replace_env_host(url)
  url.gsub(/{PP_APP_DOMAIN}/, app_domain)
     .gsub(/{PP_FULL_APP_DOMAIN}/, full_app_domain)
end

