Feature: Run command
  As a newcomer to lpn
  I want to be able to run the images managed by the tool

  Scenario Outline: Run command
    When I run `lpn run <type> -t <tag>`
    Then the output should contain:
    """
    The container [<container-name>] has been run sucessfully
    """
    And the exit status should be 0
    And I run `lpn rm`

  Examples:
    | type    | tag | container-name |
    | commerce | latest | liferay-portal-nook |
    | nightly | latest | liferay-portal-nook |
    | release | latest | liferay-portal-nook |