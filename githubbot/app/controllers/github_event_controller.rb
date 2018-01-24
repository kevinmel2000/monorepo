class GithubEventController < ApplicationController

    EVENT_HEADER = "X-Github-Event"

    # payload handles github hooks
    def payload
        
        case request.headers[EVENT_HEADER]
        when "issues"
            event_handler_issues(params[:github_event])
        else 
            render json: {
                :error => "action not supported"
            }, status: 400
        end

    end


    def ping
        msg = {
            :message => "pong"
        }
        render json: msg, content_type: "application/json", status: 200
    end


    private
    # EVENT HANDLERS
    def event_handler_issues(payload)
        render json: payload
    end
end
