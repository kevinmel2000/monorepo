class GithubEventController < ApplicationController
    include GithubEventHelper

    LOGISTIC_BOT_USERNAME = "monorepo-bot"
    MONOREPO = "lab46/monorepo"
    EVENT_HEADER = "X-Github-Event"

    def initialize
        # Get 
        @octokit_client = Octokit::Client.new(:netrc => true)
        @octokit_client.login
    end

    # payload handles github hooks
    def payload
        
        case request.headers[EVENT_HEADER]
        when "issues"
            event_handler_issues(params[:github_event])
        when "issue_comment"
            event_handler_issue_comment(params[:github_event])
        when "pull_request"
            event_handler_pull_request(params[:github_event])
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


    def test_config_circleci
        render json: Library::CircleCI.get_configuration()
    end


    private
    # EVENT HANDLERS
    def event_handler_issues(payload)
        

        if payload[:action] != "created" 
            render json: { :message => "not supported action"}, status: 200
            return
        end

        issue_id = payload[:issue][:number]

        message = "["+DateTime.now.strftime("%d/%m/%Y %H:%M")+"]"+" **Bot** acknowledge this issue "
        @octokit_client.add_comment(MONOREPO, issue_id, message)
        render json: payload
    end

    # event_handler_issue_comment handles action for issue_comment github webhook
    def event_handler_issue_comment(payload)

        # If it is not issue_comment - created event then return
        if payload[:action] != "created" 
            render json: { :message => "not supported action"}, status: 200
            return
        end

        issue_id = payload[:issue][:number]
        message = payload[:comment][:body]
        message = message.strip.downcase

        # check if comment by bot
        if payload[:comment][:user][:login] == LOGISTIC_BOT_USERNAME
            render json: {
                :message => "ignoring bot action"
            }, status: 200
            return
        end


        if message.index("/ci-test") == 0
            # check service name
            service = extract_service_name_from_comment(message)

            # trigger circle ci
            begin
                # add comment
                message = "[" + DateTime.now.strftime("%d/%m/%Y %H:%M") + "]**Bot** /ci-test is invoked"
                if service != ""
                    message += ("\nservice: "+ service)
                end
                @octokit_client.add_comment(MONOREPO, issue_id, message)

                response = @octokit_client.pull_request(MONOREPO, issue_id)
                if response[:head][:ref] == nil 
                    raise 'Response not valid'
                end

                # Get remote branch name
                remote_branch = response[:head][:ref]
                
                # Trigger CircleCI test with branch
                res = Library::CircleCI.trigger_build_branch(remote_branch, {
                    :build_parameters => {
                        :LOGISTIC_SERVICE => service
                    }
                })

            rescue
                render json: {
                    :message => "not a PR"
                }, status: 404
                return
            end

            render json: {
                :message => "accepted",
                :inspect => res.inspect
            }, status: 200
            return
        end

        render json: {
            :message => "command not found",
        }, status: 200
        return
    end

    # event_handler_pull_request handles action for pull_request github webhook
    def event_handler_pull_request(payload)
        action = payload[:action]

        if payload[:action] != "opened" && payload[:action] != "edited"
            render json: {
                :message => "action not supported",
                :action => payload[:action]
            }, status: 200
            return
        end

        pr_number = payload[:number]

        #If pull request is not against master then comment and close the pull request
        if payload[:pull_request][:base][:ref] != "master"
            begin
                # Add comment
                message = "[" + DateTime.now.strftime("%d/%m/%Y %H:%M") + "] **Bot** PR not against master, closing"
                @octokit_client.add_comment(MONOREPO, pr_number, message)

                # Close the PR
                @octokit_client.close_pull_request(MONOREPO, pr_number)
            rescue => e
                puts e.inspect
                render json: {
                    :message => "something went wrong at closing PR"
                }, status: 200
                return
            end

            render json: {
                :message => "success close pull request, not master"
            }, status: 200
            return
        end


        # else analyze the pull request
        begin
            # get service information
            services = extract_service_information(payload[:pull_request][:body])
            labels = []
            services.each do |service| 
                labels.push("service:"+service)
            end

            if labels.length > 0
                #add comment
                message = "[" + DateTime.now.strftime("%d/%m/%Y %H:%M") + "] **Bot** try to apply service labels"
                @octokit_client.add_comment(MONOREPO, pr_number, message)

                #add label
                @octokit_client.add_labels_to_an_issue(MONOREPO, pr_number, labels)
            end


        rescue => e
            puts e.inspect
            render json: {
                :message => "something went wrong at applying label"
            }, status: 200
            return
        end

        render json: {
            :message => "success"
        }, status: 200
        return
    end

end
