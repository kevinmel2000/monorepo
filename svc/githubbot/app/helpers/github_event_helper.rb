module GithubEventHelper

    def extract_service_information(message)
        lines = message.downcase.split("\n")

        services = []
        lines.each do |line|
            if line.strip.index("services") == 0 
                services = extract_service_list(line)
            end
        end

        return services

    end

    def extract_service_list(message)
        s = message.downcase.split(":")

        is_label = false
        labels = []

        if s.length < 2 
            return labels
        end

        services = s[1].strip.split(",")
        services.each do |svc|
            labels.push(svc.strip)
        end

        return labels
    end

    def extract_service_name_from_comment(message)
        service = ""
        words = message.downcase.strip.split(" ")

        if words.length < 3 
            return service
        end

        if words[0] != "/ci-test"
            return service
        end

        if words[1].strip != "service"
            return service
        end

        return words[2]

    end
end
