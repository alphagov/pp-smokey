class AssertionException < Exception
end

def assert test, msg = nil
  msg ||= "Failed assertion, no message given."
  unless test then
    msg = msg.call if Proc === msg
    raise AssertionException, msg
  end
  true
end