class ApplicationController < ActionController::Base
  #protect_from_forgery with: :exception
  # for apis make csrf token null
  protect_from_forgery with: :null_session
end
