require 'net/http'

module Library
    class CircleCI
        CIRCLECI_ENDPOINT = "https://circleci.com/api/v1.1"
        CIRCLECI_USERNAME = Rails.application.config.circleci_user_name

        # :vcs-type params
        # supported vcs-type: github or bitbucket
        CIRCLECI_VCS_TYPE = Rails.application.config.circleci_vcs_type
        CIRCLECI_PROJECT_NAME = Rails.application.config.circleci_project
        CIRCLECI_TOKEN = Rails.application.config.circleci_token
    
    
        # CircleCI trigger new build API https://circleci.com/docs/api/v1-reference/#new-build
        # path: /project/:vcs-type/:username/:project?circle-token=:token
        # {
        #     "tag": "v0.1", // optional
        #     "parallel": 2, //optional, default null
        #     "build_parameters": { // optional
        #       "RUN_EXTRA_TESTS": "true"
        #     }
        # }
        def self.trigger_build(params)
            
            url = "%{endpoint}/project/%{vcs}/%{username}/%{project}?circle-token=%{token}" % {
                endpoint: CIRCLECI_ENDPOINT,
                vcs: CIRCLECI_VCS_TYPE,
                username: CIRCLECI_USERNAME,
                project: CIRCLECI_PROJECT_NAME,
                token: CIRCLECI_TOKEN
            }


            uri = URI(url)
            req = Net::HTTP::Post.new(uri, 'Content-Type' => 'application/json')
            req.body = params.to_json

            http = Net::HTTP.new(uri.host, uri.port)
            http.use_ssl = true
            http.verify_mode = OpenSSL::SSL::VERIFY_NONE

            puts req.inspect

            res = http.request(req)
            return res
            
        end
    

        # CircleCI trigger new build branch API https://circleci.com/docs/api/v1-reference/#new-build-branch
        # path: /project/:vcs-type/:username/:project/tree/:branch?circle-token=:token
        # {
        #     "parallel": 2, //optional, default null
        #     "revision": "f1baeb913288519dd9a942499cef2873f5b1c2bf" // optional
        #     "build_parameters": { // optional
        #       "RUN_EXTRA_TESTS": "true"
        #     }
        # }
        def self.trigger_build_branch(branch, params)
            url = "%{endpoint}/project/%{vcs}/%{username}/%{project}/tree/%{branch}?circle-token=%{token}" % {
                endpoint: CIRCLECI_ENDPOINT,
                vcs: CIRCLECI_VCS_TYPE,
                username: CIRCLECI_USERNAME,
                project: CIRCLECI_PROJECT_NAME,
                branch: branch,
                token: CIRCLECI_TOKEN
            }

            uri = URI(url)
            req = Net::HTTP::Post.new(uri, 'Content-Type' => 'application/json')
            req.body = params.to_json

            http = Net::HTTP.new(uri.host, uri.port)
            http.use_ssl = true
            http.verify_mode = OpenSSL::SSL::VERIFY_NONE

            res = http.request(req)
    
            return res

        end

    
        def self.get_configuration() 
            return {
                :username => CIRCLECI_USERNAME,
                :vcs => CIRCLECI_VCS_TYPE,
                :project => CIRCLECI_PROJECT_NAME,
                :token => CIRCLECI_TOKEN
            }
        end
    end

end