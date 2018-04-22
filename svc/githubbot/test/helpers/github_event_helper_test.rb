require 'test_helper'

class GithubEventHelperTest < ActionView::TestCase 
    test 'test extract service list' do

        test_cases = []
        test_cases.push({
            :message => %q(
services: ongkirapp, addressapp
),
            :expected => ["ongkirapp", "addressapp"]
        })

        test_cases.each do |tc|
            result = extract_service_list(tc[:message])
            assert_equal result.sort, tc[:expected].sort
        end
    end

    test 'extract service from message body' do 
        test_cases = []
        test_cases.push({
            :message => %q(
[T00001] Changes About Something

changes:
 - wow something
 - wow something else

services: ongkirapp, addressapp
            ),
            :expected => ["ongkirapp", "addressapp"]
        })

        test_cases.each do |tc| 
            result = extract_service_information(tc[:message])
            assert_equal result.sort, tc[:expected].sort
        end

    end

    test 'extract service name from comment' do
        test_cases = []
        test_cases.push({
            :message => "/ci-test service ongkirapp",
            :expected => "ongkirapp"
        })

        test_cases.push({
            :message => " /ci-test  service  ongkirapp ",
            :expected => "ongkirapp"
        })

        test_cases.each do |tc| 
            service = extract_service_name_from_comment(tc[:message])
            assert_equal service, tc[:expected]
        end
    end

end