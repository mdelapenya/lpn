require 'aruba/cucumber'

Aruba.configure do |config|
    config.exit_timeout = 120
    config.io_wait_timeout = 120
    config.startup_wait_time = 1
end
