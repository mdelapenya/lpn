Feature: Version command
  As a newcomer to lpn
  I want to be able to check the version

  Scenario: Version command
    When I run `lpn version`
    Then the exit status should be 0
    And the output should contain:
    """
    0.3.0
    """