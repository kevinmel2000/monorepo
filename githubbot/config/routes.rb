Rails.application.routes.draw do
  # For details on the DSL available within this file, see http://guides.rubyonrails.org/routing.html

  post '/payload', to: 'github_event#payload'
  get '/ping', to: 'github_event#ping'
end
