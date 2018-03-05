Feature: License command
  As a newcomer to lpn
  I want to be able to check the license

  Scenario: License command
    When I run `lpn license`
    Then the exit status should be 0
    And the output should contain:
    """
    Redistributions in binary form must reproduce the above copyright notice
    """