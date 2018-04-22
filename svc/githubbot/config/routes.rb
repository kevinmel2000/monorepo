Rails.application.routes.draw do
  get 'pages/home'

  # For details on the DSL available within this file, see http://guides.rubyonrails.org/routing.html

  root to: 'pages#home'

  get '/test_x', to: "github_event#test_extract"
  get '/config_ausodfjiasj10010', to: "github_event#test_config_circleci"
  post '/payload', to: 'github_event#payload'
  get '/ping', to: 'github_event#ping'
end
