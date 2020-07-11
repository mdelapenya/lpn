Feature: License command
  As a newcomer to lpn
  I want to be able to check the license

  Scenario: License command
    When I run `lpn license`
    Then the exit status should be 0
    And the output should contain:
    """
    Copyright (c) 2000-present Liferay, Inc. All rights reserved.
    """