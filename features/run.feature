Feature: Run command
  As a newcomer to lpn
  I want to be able to run the images managed by the tool

  Scenario Outline: Run command
    Given I run `lpn run <type> -t <tag>`
    Then the output should contain:
    """
    The container [lpn-<type>] has been run successfully
    """
    And the exit status should be 0
    And I run `lpn rm <type>`

  Examples:
    | type    | tag |
    | ce      | 7.0.6-ga7 |
    | dxp     | 7.0.10.8 |
    | nightly | master |
    | release | latest |

  Scenario Outline: Run command with debug enabled
    Given I run `lpn run <type> -t <tag> -d`
    When I run `docker exec lpn-<type> env`
    Then the output should contain:
    """
    <variable>=true
    """
    And I run `lpn rm <type>`

  Examples:
    | type    | tag | variable |
    | ce      | 7.0.6-ga7 | LIFERAY_JPDA_ENABLED |
    | dxp     | 7.0.10.8 | LIFERAY_JPDA_ENABLED |
    | nightly | master | LIFERAY_JPDA_ENABLED |
    | release | latest | DEBUG_MODE |

  Scenario Outline: Run command with failure
    Given I run `docker run -d --name nginx-<type> -p 9999:80 nginx:1.12.2-alpine`
    When I run `lpn run <type> -t <tag> -p 9999`
    Then the output should contain:
    """
    Impossible to run the container [lpn-<type>]
    """
    And the exit status should be 1
    And I run `lpn rm <type>`
    And I run `docker rm -fv nginx-<type>`

  Examples:
    | type    | tag |
    | ce      | 7.0.6-ga7 |
    | dxp     | 7.0.10.8 |
    | nightly | master |
    | release | latest |

  Scenario Outline: Run command with memory enabled for non official images
    Given I run `lpn run <type> -t <tag> -m 1024m`
    When I run `docker exec lpn-<type> ps aux | grep -e "-Xms1024m -Xmx1024m" | wc -l | xargs`
    Then the output should contain:
    """
    1
    """
    And I run `docker exec lpn-<type> env`
    Then the output should contain:
    """
    <variable>=
    """
    And I run `lpn rm <type>`

  Examples:
    | type    | tag | variable |
    | release | latest | JVM_TUNING_MEMORY |

  Scenario Outline: Run command with memory enabled for official images
    Given I run `lpn run <type> -t <tag> -m "-Xms1024m -Xmx1024m"`
    When I run `docker exec lpn-<type> ps aux | grep -e tomcat | xargs`
    Then the output should contain:
    """
    -Xms1024m -Xmx1024m
    """
    And I run `lpn rm <type>`

  Examples:
    | type    | tag | variable |
    | ce      | 7.0.6-ga7 | LIFERAY_JVM_OPTS |
    | dxp     | 7.0.10.8 | LIFERAY_JVM_OPTS |
    | nightly | master | LIFERAY_JVM_OPTS |